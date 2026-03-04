package services

import (
	"context"
	"fmt"
	"hash/fnv"
	"sort"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/workermill-examples/flagdeck/api/internal/models"
)

type EvaluationService struct {
	flagsCollection        *mongo.Collection
	environmentsCollection *mongo.Collection
	segmentsCollection     *mongo.Collection
}

type EvaluationResult struct {
	Value        interface{} `json:"value"`
	Reason       string      `json:"reason"`
	RuleID       *string     `json:"rule_id,omitempty"`
	EvaluationMs int64       `json:"evaluation_ms"`
}

type EvaluationRequest struct {
	FlagKey     string                 `json:"flag_key"`
	Environment string                 `json:"environment"`
	UserID      string                 `json:"user_id"`
	UserContext map[string]interface{} `json:"user_context,omitempty"`
}

type BulkEvaluationRequest struct {
	Environment string                 `json:"environment"`
	UserID      string                 `json:"user_id"`
	UserContext map[string]interface{} `json:"user_context,omitempty"`
	FlagKeys    []string               `json:"flag_keys"`
}

type BulkEvaluationResponse struct {
	Results map[string]EvaluationResult `json:"results"`
}

func NewEvaluationService(flagsCollection, environmentsCollection, segmentsCollection *mongo.Collection) *EvaluationService {
	return &EvaluationService{
		flagsCollection:        flagsCollection,
		environmentsCollection: environmentsCollection,
		segmentsCollection:     segmentsCollection,
	}
}

// EvaluateFlag evaluates a single flag for a user
func (s *EvaluationService) EvaluateFlag(req EvaluationRequest) (*EvaluationResult, error) {
	startTime := time.Now()

	// 1. Check if flag exists and is active
	var flag models.Flag
	err := s.flagsCollection.FindOne(context.Background(), bson.M{"key": req.FlagKey}).Decode(&flag)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &EvaluationResult{
				Value:        nil,
				Reason:       "flag_not_found",
				EvaluationMs: time.Since(startTime).Milliseconds(),
			}, nil
		}
		return nil, err
	}

	// Check if flag is globally active
	if !flag.IsActive {
		return &EvaluationResult{
			Value:        nil,
			Reason:       "flag_disabled",
			EvaluationMs: time.Since(startTime).Milliseconds(),
		}, nil
	}

	// 2. Check environment configuration
	envConfig, exists := flag.Environments[req.Environment]
	if !exists {
		// Check if environment exists
		var env models.Environment
		err := s.environmentsCollection.FindOne(context.Background(), bson.M{"key": req.Environment}).Decode(&env)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return &EvaluationResult{
					Value:        nil,
					Reason:       "environment_not_found",
					EvaluationMs: time.Since(startTime).Milliseconds(),
				}, nil
			}
			return nil, err
		}

		// Environment exists but flag not configured for it
		return &EvaluationResult{
			Value:        nil,
			Reason:       "flag_not_configured_for_environment",
			EvaluationMs: time.Since(startTime).Milliseconds(),
		}, nil
	}

	// Check if flag is enabled in this environment
	if !envConfig.Enabled {
		return &EvaluationResult{
			Value:        envConfig.DefaultValue,
			Reason:       "environment_disabled",
			EvaluationMs: time.Since(startTime).Milliseconds(),
		}, nil
	}

	// 3. Evaluate targeting rules by priority
	if len(envConfig.TargetingRules) > 0 {
		// Sort rules by priority (lower number = higher priority)
		rules := make([]models.TargetingRule, len(envConfig.TargetingRules))
		copy(rules, envConfig.TargetingRules)
		sort.Slice(rules, func(i, j int) bool {
			return rules[i].Priority < rules[j].Priority
		})

		// Evaluate each rule
		for _, rule := range rules {
			matches, err := s.evaluateTargetingRule(rule, req.UserID, req.UserContext)
			if err != nil {
				return nil, err
			}
			if matches {
				ruleID := rule.ID.Hex()
				return &EvaluationResult{
					Value:        rule.Value,
					Reason:       "targeting_rule_match",
					RuleID:       &ruleID,
					EvaluationMs: time.Since(startTime).Milliseconds(),
				}, nil
			}
		}
	}

	// 4. Check rollout percentage using FNV-1a hash
	if envConfig.RolloutPercent < 100 {
		userInRollout := s.isUserInRollout(req.FlagKey, req.UserID, envConfig.RolloutPercent)
		if !userInRollout {
			return &EvaluationResult{
				Value:        envConfig.DefaultValue,
				Reason:       "rollout_excluded",
				EvaluationMs: time.Since(startTime).Milliseconds(),
			}, nil
		}
	}

	// 5. Return the flag value
	return &EvaluationResult{
		Value:        envConfig.DefaultValue,
		Reason:       "default_value",
		EvaluationMs: time.Since(startTime).Milliseconds(),
	}, nil
}

// EvaluateFlags evaluates multiple flags for a user
func (s *EvaluationService) EvaluateFlags(req BulkEvaluationRequest) (*BulkEvaluationResponse, error) {
	results := make(map[string]EvaluationResult)

	for _, flagKey := range req.FlagKeys {
		evalReq := EvaluationRequest{
			FlagKey:     flagKey,
			Environment: req.Environment,
			UserID:      req.UserID,
			UserContext: req.UserContext,
		}

		result, err := s.EvaluateFlag(evalReq)
		if err != nil {
			// On error, return error reason
			results[flagKey] = EvaluationResult{
				Value:  nil,
				Reason: "evaluation_error",
			}
		} else {
			results[flagKey] = *result
		}
	}

	return &BulkEvaluationResponse{
		Results: results,
	}, nil
}

// evaluateTargetingRule evaluates if a targeting rule matches the user
func (s *EvaluationService) evaluateTargetingRule(rule models.TargetingRule, userID string, userContext map[string]interface{}) (bool, error) {
	// All conditions must match (AND logic)
	for _, condition := range rule.Conditions {
		matches, err := s.evaluateCondition(condition, userID, userContext)
		if err != nil {
			return false, err
		}
		if !matches {
			return false, nil
		}
	}
	return true, nil
}

// evaluateCondition evaluates a single condition
func (s *EvaluationService) evaluateCondition(condition models.Condition, userID string, userContext map[string]interface{}) (bool, error) {
	var actualValue interface{}

	// Get the actual value to compare against
	switch condition.Attribute {
	case "user_id":
		actualValue = userID
	case "segment":
		// Special handling for segment membership
		if segmentKey, ok := condition.Value.(string); ok {
			inSegment, err := s.isUserInSegment(userID, userContext, segmentKey)
			if err != nil {
				return false, err
			}
			switch condition.Operator {
			case "in":
				return inSegment, nil
			case "not_in":
				return !inSegment, nil
			default:
				return false, fmt.Errorf("unsupported operator for segment: %s", condition.Operator)
			}
		}
		return false, nil
	default:
		// Look up in user context
		if userContext != nil {
			actualValue = userContext[condition.Attribute]
		}
	}

	// Apply the operator
	return s.applyOperator(actualValue, condition.Operator, condition.Value)
}

// applyOperator applies the comparison operator
func (s *EvaluationService) applyOperator(actualValue interface{}, operator string, expectedValue interface{}) (bool, error) {
	switch operator {
	case "equals":
		return actualValue == expectedValue, nil
	case "not_equals":
		return actualValue != expectedValue, nil
	case "in":
		// Check if actualValue is in the expectedValue slice
		if expectedSlice, ok := expectedValue.([]interface{}); ok {
			for _, item := range expectedSlice {
				if actualValue == item {
					return true, nil
				}
			}
		}
		return false, nil
	case "not_in":
		// Check if actualValue is NOT in the expectedValue slice
		if expectedSlice, ok := expectedValue.([]interface{}); ok {
			for _, item := range expectedSlice {
				if actualValue == item {
					return false, nil
				}
			}
			return true, nil
		}
		return true, nil
	case "contains":
		// String contains check
		if actualStr, ok := actualValue.(string); ok {
			if expectedStr, ok := expectedValue.(string); ok {
				return strings.Contains(actualStr, expectedStr), nil
			}
		}
		return false, nil
	case "starts_with":
		// String starts with check
		if actualStr, ok := actualValue.(string); ok {
			if expectedStr, ok := expectedValue.(string); ok {
				return strings.HasPrefix(actualStr, expectedStr), nil
			}
		}
		return false, nil
	case "ends_with":
		// String ends with check
		if actualStr, ok := actualValue.(string); ok {
			if expectedStr, ok := expectedValue.(string); ok {
				return strings.HasSuffix(actualStr, expectedStr), nil
			}
		}
		return false, nil
	case "greater_than":
		// Numeric comparison
		if actualNum, ok := s.toFloat64(actualValue); ok {
			if expectedNum, ok := s.toFloat64(expectedValue); ok {
				return actualNum > expectedNum, nil
			}
		}
		return false, nil
	case "greater_than_equal":
		// Numeric comparison
		if actualNum, ok := s.toFloat64(actualValue); ok {
			if expectedNum, ok := s.toFloat64(expectedValue); ok {
				return actualNum >= expectedNum, nil
			}
		}
		return false, nil
	case "less_than":
		// Numeric comparison
		if actualNum, ok := s.toFloat64(actualValue); ok {
			if expectedNum, ok := s.toFloat64(expectedValue); ok {
				return actualNum < expectedNum, nil
			}
		}
		return false, nil
	case "less_than_equal":
		// Numeric comparison
		if actualNum, ok := s.toFloat64(actualValue); ok {
			if expectedNum, ok := s.toFloat64(expectedValue); ok {
				return actualNum <= expectedNum, nil
			}
		}
		return false, nil
	default:
		return false, fmt.Errorf("unsupported operator: %s", operator)
	}
}

// toFloat64 converts various numeric types to float64
func (s *EvaluationService) toFloat64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	default:
		return 0, false
	}
}

// isUserInRollout determines if a user should be included in the rollout using FNV-1a hash
func (s *EvaluationService) isUserInRollout(flagKey, userID string, rolloutPercent int) bool {
	if rolloutPercent <= 0 {
		return false
	}
	if rolloutPercent >= 100 {
		return true
	}

	// Create hash input by combining flag key and user ID
	hashInput := flagKey + ":" + userID

	// Calculate FNV-1a hash
	hash := fnv.New32a()
	hash.Write([]byte(hashInput))
	hashValue := hash.Sum32()

	// Convert to percentage (0-99)
	userPercent := int(hashValue % 100)

	// User is in rollout if their hash percentage is less than the rollout percentage
	return userPercent < rolloutPercent
}

// isUserInSegment checks if a user belongs to a specific segment
func (s *EvaluationService) isUserInSegment(userID string, userContext map[string]interface{}, segmentKey string) (bool, error) {
	// Get the segment
	var segment models.Segment
	err := s.segmentsCollection.FindOne(context.Background(), bson.M{"key": segmentKey}).Decode(&segment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil // Segment doesn't exist, user not in it
		}
		return false, err
	}

	// Evaluate all segment rules (AND logic)
	for _, rule := range segment.Rules {
		condition := models.Condition{
			Attribute: rule.Attribute,
			Operator:  rule.Operator,
			Value:     rule.Value,
		}

		matches, err := s.evaluateCondition(condition, userID, userContext)
		if err != nil {
			return false, err
		}
		if !matches {
			return false, nil
		}
	}

	return true, nil
}
