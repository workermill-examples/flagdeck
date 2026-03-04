package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"github.com/workermill-examples/flagdeck/api/internal/config"
	"github.com/workermill-examples/flagdeck/api/internal/database"
	"github.com/workermill-examples/flagdeck/api/internal/models"
)

func main() {
	log.Println("Starting seed script...")

	// Load configuration
	cfg := config.Load()

	// Connect to MongoDB
	mongodb, err := database.NewMongoDB(cfg.MongodbURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongodb.Close()

	ctx := context.Background()

	// Seed data
	log.Println("Seeding users...")
	userID := seedUsers(ctx, mongodb)

	log.Println("Seeding environments...")
	seedEnvironments(ctx, mongodb)

	log.Println("Seeding API keys...")
	seedApiKeys(ctx, mongodb)

	log.Println("Seeding flags...")
	seedFlags(ctx, mongodb)

	log.Println("Seeding segments...")
	seedSegments(ctx, mongodb)

	log.Println("Seeding experiments...")
	seedExperiments(ctx, mongodb)

	log.Println("Seeding audit log...")
	seedAuditLog(ctx, mongodb, userID)

	log.Println("Seed script completed successfully!")
}

func seedUsers(ctx context.Context, mongodb *database.MongoDB) primitive.ObjectID {
	usersCollection := mongodb.UsersCollection()

	// Hash password for admin user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("demo1234"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	adminUser := models.User{
		ID:        primitive.NewObjectID(),
		Email:     "demo@workermill.com",
		Name:      "Demo Admin",
		Role:      "admin",
		Password:  string(hashedPassword),
		CreatedAt: time.Now().AddDate(0, 0, -14), // 14 days ago
		UpdatedAt: time.Now().AddDate(0, 0, -14),
	}

	// Upsert admin user (using email as unique key)
	filter := bson.M{"email": adminUser.Email}
	update := bson.M{"$setOnInsert": adminUser}
	opts := options.Update().SetUpsert(true)

	result, err := usersCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Fatalf("Failed to upsert admin user: %v", err)
	}

	var userID primitive.ObjectID
	if result.UpsertedID != nil {
		userID = result.UpsertedID.(primitive.ObjectID)
		log.Printf("Created admin user with ID: %v", userID)
	} else {
		// User already exists, get the ID
		var existingUser models.User
		err = usersCollection.FindOne(ctx, filter).Decode(&existingUser)
		if err != nil {
			log.Fatalf("Failed to find existing admin user: %v", err)
		}
		userID = existingUser.ID
		log.Printf("Admin user already exists with ID: %v", userID)
	}

	return userID
}

func seedEnvironments(ctx context.Context, mongodb *database.MongoDB) {
	environmentsCollection := mongodb.EnvironmentsCollection()

	environments := []models.Environment{
		{
			ID:          primitive.NewObjectID(),
			Key:         "development",
			Name:        "Development",
			Description: "Development environment for testing new features",
			Color:       "#10b981", // Green
			CreatedAt:   time.Now().AddDate(0, 0, -14),
			UpdatedAt:   time.Now().AddDate(0, 0, -14),
		},
		{
			ID:          primitive.NewObjectID(),
			Key:         "staging",
			Name:        "Staging",
			Description: "Staging environment for pre-production testing",
			Color:       "#f59e0b", // Yellow/Orange
			CreatedAt:   time.Now().AddDate(0, 0, -13),
			UpdatedAt:   time.Now().AddDate(0, 0, -13),
		},
		{
			ID:          primitive.NewObjectID(),
			Key:         "production",
			Name:        "Production",
			Description: "Live production environment",
			Color:       "#ef4444", // Red
			CreatedAt:   time.Now().AddDate(0, 0, -12),
			UpdatedAt:   time.Now().AddDate(0, 0, -12),
		},
	}

	for _, env := range environments {
		filter := bson.M{"key": env.Key}
		update := bson.M{"$setOnInsert": env}
		opts := options.Update().SetUpsert(true)

		_, err := environmentsCollection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			log.Printf("Failed to upsert environment %s: %v", env.Key, err)
		} else {
			log.Printf("Upserted environment: %s", env.Key)
		}
	}
}

func generateAPIKey() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return "fd_" + hex.EncodeToString(bytes)
}

func seedApiKeys(ctx context.Context, mongodb *database.MongoDB) {
	apiKeysCollection := mongodb.APIKeysCollection()

	// Generate API keys
	devKeyRaw := generateAPIKey()
	prodKeyRaw := generateAPIKey()

	devKeyHash, _ := bcrypt.GenerateFromPassword([]byte(devKeyRaw), bcrypt.DefaultCost)
	prodKeyHash, _ := bcrypt.GenerateFromPassword([]byte(prodKeyRaw), bcrypt.DefaultCost)

	apiKeys := []models.ApiKey{
		{
			ID:          primitive.NewObjectID(),
			Name:        "Development API Key",
			KeyHash:     string(devKeyHash),
			Environment: "development",
			LastUsedAt:  nil,
			CreatedAt:   time.Now().AddDate(0, 0, -10),
			UpdatedAt:   time.Now().AddDate(0, 0, -10),
		},
		{
			ID:          primitive.NewObjectID(),
			Name:        "Production API Key",
			KeyHash:     string(prodKeyHash),
			Environment: "production",
			LastUsedAt:  timePtr(time.Now().AddDate(0, 0, -1)), // Used yesterday
			CreatedAt:   time.Now().AddDate(0, 0, -8),
			UpdatedAt:   time.Now().AddDate(0, 0, -1),
		},
	}

	for i, apiKey := range apiKeys {
		filter := bson.M{"name": apiKey.Name, "environment": apiKey.Environment}
		update := bson.M{"$setOnInsert": apiKey}
		opts := options.Update().SetUpsert(true)

		result, err := apiKeysCollection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			log.Printf("Failed to upsert API key %s: %v", apiKey.Name, err)
		} else if result.UpsertedID != nil {
			var keyName string
			if i == 0 {
				keyName = "Development"
				log.Printf("Created %s API key: %s", keyName, devKeyRaw)
			} else {
				keyName = "Production"
				log.Printf("Created %s API key: %s", keyName, prodKeyRaw)
			}
		} else {
			log.Printf("API key %s already exists", apiKey.Name)
		}
	}
}

func seedFlags(ctx context.Context, mongodb *database.MongoDB) {
	flagsCollection := mongodb.FlagsCollection()

	flags := []models.Flag{
		{
			ID:          primitive.NewObjectID(),
			Key:         "new_dashboard",
			Name:        "New Dashboard",
			Description: "Enable the redesigned dashboard interface",
			IsActive:    true,
			Type:        "boolean",
			Environments: map[string]models.FlagEnvironment{
				"development": {
					Enabled:        true,
					DefaultValue:   true,
					RolloutPercent: 100,
					TargetingRules: []models.TargetingRule{},
				},
				"staging": {
					Enabled:        true,
					DefaultValue:   true,
					RolloutPercent: 50,
					TargetingRules: []models.TargetingRule{},
				},
				"production": {
					Enabled:        true,
					DefaultValue:   false,
					RolloutPercent: 10,
					TargetingRules: []models.TargetingRule{
						{
							ID:       primitive.NewObjectID(),
							Priority: 1,
							Conditions: []models.Condition{
								{
									Attribute: "user_role",
									Operator:  "equals",
									Value:     "admin",
								},
							},
							Value: true,
						},
					},
				},
			},
			CreatedAt: time.Now().AddDate(0, 0, -11),
			UpdatedAt: time.Now().AddDate(0, 0, -2),
		},
		{
			ID:          primitive.NewObjectID(),
			Key:         "api_rate_limit",
			Name:        "API Rate Limit",
			Description: "Rate limit for API endpoints",
			IsActive:    true,
			Type:        "number",
			Environments: map[string]models.FlagEnvironment{
				"development": {
					Enabled:        true,
					DefaultValue:   1000,
					RolloutPercent: 100,
					TargetingRules: []models.TargetingRule{},
				},
				"staging": {
					Enabled:        true,
					DefaultValue:   500,
					RolloutPercent: 100,
					TargetingRules: []models.TargetingRule{},
				},
				"production": {
					Enabled:        true,
					DefaultValue:   100,
					RolloutPercent: 100,
					TargetingRules: []models.TargetingRule{
						{
							ID:       primitive.NewObjectID(),
							Priority: 1,
							Conditions: []models.Condition{
								{
									Attribute: "subscription_tier",
									Operator:  "equals",
									Value:     "premium",
								},
							},
							Value: 500,
						},
					},
				},
			},
			CreatedAt: time.Now().AddDate(0, 0, -9),
			UpdatedAt: time.Now().AddDate(0, 0, -1),
		},
		{
			ID:          primitive.NewObjectID(),
			Key:         "feature_announcements",
			Name:        "Feature Announcements",
			Description: "Show feature announcement banner",
			IsActive:    true,
			Type:        "string",
			Environments: map[string]models.FlagEnvironment{
				"development": {
					Enabled:        true,
					DefaultValue:   "Check out our new AI-powered insights!",
					RolloutPercent: 100,
					TargetingRules: []models.TargetingRule{},
				},
				"staging": {
					Enabled:        true,
					DefaultValue:   "Check out our new AI-powered insights!",
					RolloutPercent: 100,
					TargetingRules: []models.TargetingRule{},
				},
				"production": {
					Enabled:        false,
					DefaultValue:   "",
					RolloutPercent: 0,
					TargetingRules: []models.TargetingRule{},
				},
			},
			CreatedAt: time.Now().AddDate(0, 0, -7),
			UpdatedAt: time.Now().AddDate(0, 0, -7),
		},
		{
			ID:          primitive.NewObjectID(),
			Key:         "experimental_search",
			Name:        "Experimental Search",
			Description: "Enable experimental search algorithm",
			IsActive:    true,
			Type:        "boolean",
			Environments: map[string]models.FlagEnvironment{
				"development": {
					Enabled:        true,
					DefaultValue:   true,
					RolloutPercent: 100,
					TargetingRules: []models.TargetingRule{},
				},
				"staging": {
					Enabled:        true,
					DefaultValue:   false,
					RolloutPercent: 25,
					TargetingRules: []models.TargetingRule{},
				},
				"production": {
					Enabled:        false,
					DefaultValue:   false,
					RolloutPercent: 0,
					TargetingRules: []models.TargetingRule{},
				},
			},
			CreatedAt: time.Now().AddDate(0, 0, -6),
			UpdatedAt: time.Now().AddDate(0, 0, -3),
		},
		{
			ID:          primitive.NewObjectID(),
			Key:         "maintenance_mode",
			Name:        "Maintenance Mode",
			Description: "Enable maintenance mode banner",
			IsActive:    false,
			Type:        "boolean",
			Environments: map[string]models.FlagEnvironment{
				"development": {
					Enabled:        false,
					DefaultValue:   false,
					RolloutPercent: 0,
					TargetingRules: []models.TargetingRule{},
				},
				"staging": {
					Enabled:        false,
					DefaultValue:   false,
					RolloutPercent: 0,
					TargetingRules: []models.TargetingRule{},
				},
				"production": {
					Enabled:        false,
					DefaultValue:   false,
					RolloutPercent: 0,
					TargetingRules: []models.TargetingRule{},
				},
			},
			CreatedAt: time.Now().AddDate(0, 0, -5),
			UpdatedAt: time.Now().AddDate(0, 0, -5),
		},
		{
			ID:          primitive.NewObjectID(),
			Key:         "theme_config",
			Name:        "Theme Configuration",
			Description: "Application theme configuration",
			IsActive:    true,
			Type:        "json",
			Environments: map[string]models.FlagEnvironment{
				"development": {
					Enabled: true,
					DefaultValue: map[string]interface{}{
						"primary_color":   "#3b82f6",
						"secondary_color": "#64748b",
						"dark_mode":       true,
					},
					RolloutPercent: 100,
					TargetingRules: []models.TargetingRule{},
				},
				"staging": {
					Enabled: true,
					DefaultValue: map[string]interface{}{
						"primary_color":   "#3b82f6",
						"secondary_color": "#64748b",
						"dark_mode":       false,
					},
					RolloutPercent: 100,
					TargetingRules: []models.TargetingRule{},
				},
				"production": {
					Enabled: true,
					DefaultValue: map[string]interface{}{
						"primary_color":   "#1f2937",
						"secondary_color": "#6b7280",
						"dark_mode":       false,
					},
					RolloutPercent: 100,
					TargetingRules: []models.TargetingRule{},
				},
			},
			CreatedAt: time.Now().AddDate(0, 0, -4),
			UpdatedAt: time.Now().AddDate(0, 0, -1),
		},
	}

	// Add more flags to reach 10+ requirement
	moreFlags := []models.Flag{
		createBooleanFlag("mobile_app_upsell", "Mobile App Upsell", "Show mobile app download prompt", true, false, 5, time.Now().AddDate(0, 0, -8)),
		createStringFlag("support_chat_widget", "Support Chat Widget", "Customer support chat widget", "enabled", "disabled", time.Now().AddDate(0, 0, -6)),
		createNumberFlag("max_upload_size", "Max Upload Size", "Maximum file upload size in MB", 50, 10, time.Now().AddDate(0, 0, -4)),
		createBooleanFlag("beta_features", "Beta Features", "Enable beta feature access", false, false, 15, time.Now().AddDate(0, 0, -3)),
		createStringFlag("welcome_message", "Welcome Message", "User welcome message", "Welcome to FlagDeck! 🚀", "", time.Now().AddDate(0, 0, -2)),
	}

	allFlags := append(flags, moreFlags...)

	for _, flag := range allFlags {
		filter := bson.M{"key": flag.Key}
		update := bson.M{"$setOnInsert": flag}
		opts := options.Update().SetUpsert(true)

		_, err := flagsCollection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			log.Printf("Failed to upsert flag %s: %v", flag.Key, err)
		} else {
			log.Printf("Upserted flag: %s", flag.Key)
		}
	}
}

func seedSegments(ctx context.Context, mongodb *database.MongoDB) {
	segmentsCollection := mongodb.SegmentsCollection()

	segments := []models.Segment{
		{
			ID:          primitive.NewObjectID(),
			Key:         "premium_users",
			Name:        "Premium Users",
			Description: "Users with premium subscription",
			Rules: []models.SegmentRule{
				{
					Attribute: "subscription_tier",
					Operator:  "equals",
					Value:     "premium",
				},
			},
			CreatedAt: time.Now().AddDate(0, 0, -10),
			UpdatedAt: time.Now().AddDate(0, 0, -10),
		},
		{
			ID:          primitive.NewObjectID(),
			Key:         "beta_testers",
			Name:        "Beta Testers",
			Description: "Users enrolled in beta testing program",
			Rules: []models.SegmentRule{
				{
					Attribute: "beta_tester",
					Operator:  "equals",
					Value:     true,
				},
				{
					Attribute: "account_age_days",
					Operator:  "greater_than",
					Value:     30,
				},
			},
			CreatedAt: time.Now().AddDate(0, 0, -8),
			UpdatedAt: time.Now().AddDate(0, 0, -8),
		},
		{
			ID:          primitive.NewObjectID(),
			Key:         "enterprise_clients",
			Name:        "Enterprise Clients",
			Description: "Large enterprise customers",
			Rules: []models.SegmentRule{
				{
					Attribute: "organization_size",
					Operator:  "greater_than",
					Value:     100,
				},
				{
					Attribute: "subscription_tier",
					Operator:  "in",
					Value:     []string{"enterprise", "custom"},
				},
			},
			CreatedAt: time.Now().AddDate(0, 0, -6),
			UpdatedAt: time.Now().AddDate(0, 0, -6),
		},
	}

	for _, segment := range segments {
		filter := bson.M{"key": segment.Key}
		update := bson.M{"$setOnInsert": segment}
		opts := options.Update().SetUpsert(true)

		_, err := segmentsCollection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			log.Printf("Failed to upsert segment %s: %v", segment.Key, err)
		} else {
			log.Printf("Upserted segment: %s", segment.Key)
		}
	}
}

func seedExperiments(ctx context.Context, mongodb *database.MongoDB) {
	experimentsCollection := mongodb.ExperimentsCollection()

	experiments := []models.Experiment{
		{
			ID:          primitive.NewObjectID(),
			Key:         "dashboard_layout_test",
			Name:        "Dashboard Layout A/B Test",
			Description: "Testing new dashboard layout vs current layout",
			FlagKey:     "new_dashboard",
			Environment: "production",
			Status:      "running",
			StartDate:   timePtr(time.Now().AddDate(0, 0, -7)),
			EndDate:     timePtr(time.Now().AddDate(0, 0, 7)),
			Variants: []models.ExperimentVariant{
				{
					Name:         "control",
					Value:        false,
					TrafficSplit: 50,
				},
				{
					Name:         "new_layout",
					Value:        true,
					TrafficSplit: 50,
				},
			},
			Results: map[string]interface{}{
				"control": map[string]interface{}{
					"impressions": 2340.0,
					"conversions": 156.0,
					"revenue":     4680.50,
				},
				"new_layout": map[string]interface{}{
					"impressions": 2285.0,
					"conversions": 189.0,
					"revenue":     5670.25,
				},
			},
			CreatedAt: time.Now().AddDate(0, 0, -8),
			UpdatedAt: time.Now().AddDate(0, 0, -1),
		},
		{
			ID:          primitive.NewObjectID(),
			Key:         "search_algorithm_test",
			Name:        "Search Algorithm Performance Test",
			Description: "Comparing search response times and relevance",
			FlagKey:     "experimental_search",
			Environment: "staging",
			Status:      "completed",
			StartDate:   timePtr(time.Now().AddDate(0, 0, -14)),
			EndDate:     timePtr(time.Now().AddDate(0, 0, -7)),
			Variants: []models.ExperimentVariant{
				{
					Name:         "current_search",
					Value:        false,
					TrafficSplit: 60,
				},
				{
					Name:         "experimental_search",
					Value:        true,
					TrafficSplit: 40,
				},
			},
			Results: map[string]interface{}{
				"current_search": map[string]interface{}{
					"impressions": 5240.0,
					"conversions": 324.0,
					"revenue":     2160.75,
				},
				"experimental_search": map[string]interface{}{
					"impressions": 3490.0,
					"conversions": 267.0,
					"revenue":     2890.40,
				},
			},
			CreatedAt: time.Now().AddDate(0, 0, -15),
			UpdatedAt: time.Now().AddDate(0, 0, -7),
		},
	}

	for _, experiment := range experiments {
		filter := bson.M{"key": experiment.Key}
		update := bson.M{"$setOnInsert": experiment}
		opts := options.Update().SetUpsert(true)

		_, err := experimentsCollection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			log.Printf("Failed to upsert experiment %s: %v", experiment.Key, err)
		} else {
			log.Printf("Upserted experiment: %s", experiment.Key)
		}
	}
}

func seedAuditLog(ctx context.Context, mongodb *database.MongoDB, userID primitive.ObjectID) {
	auditCollection := mongodb.AuditLogCollection()

	// Generate 60+ audit entries spread across 14 days (business hours)
	// This ensures we exceed the 50+ requirement with realistic activity patterns
	baseTime := time.Now().AddDate(0, 0, -14)

	auditEntries := []models.AuditLogEntry{}

	// Day 1: Initial setup (4 entries)
	day1 := baseTime
	auditEntries = append(auditEntries,
		createAuditEntry(userID, "environment.created", "environment", "development", "Development environment created", day1.Add(9*time.Hour)),
		createAuditEntry(userID, "environment.created", "environment", "staging", "Staging environment created", day1.Add(9*time.Hour+15*time.Minute)),
		createAuditEntry(userID, "environment.created", "environment", "production", "Production environment created", day1.Add(9*time.Hour+30*time.Minute)),
		createAuditEntry(userID, "api_key.created", "api_key", "dev-key-1", "Development API key created", day1.Add(10*time.Hour)),
	)

	// Day 2: Flag creation and initial configuration (6 entries)
	day2 := baseTime.AddDate(0, 0, 1)
	auditEntries = append(auditEntries,
		createAuditEntry(userID, "flag.created", "flag", "new_dashboard", "New Dashboard flag created", day2.Add(9*time.Hour)),
		createAuditEntry(userID, "flag.created", "flag", "api_rate_limit", "API Rate Limit flag created", day2.Add(9*time.Hour+30*time.Minute)),
		createAuditEntry(userID, "flag.updated", "flag", "new_dashboard", "Targeting rules configured for new_dashboard", day2.Add(10*time.Hour)),
		createAuditEntry(userID, "flag.toggled", "flag", "new_dashboard", "Enabled new_dashboard in staging", day2.Add(11*time.Hour)),
		createAuditEntry(userID, "flag.updated", "flag", "api_rate_limit", "Rate limit values configured", day2.Add(14*time.Hour)),
		createAuditEntry(userID, "flag.toggled", "flag", "api_rate_limit", "Enabled api_rate_limit in production", day2.Add(15*time.Hour)),
	)

	// Day 3: Segment and user management (4 entries)
	day3 := baseTime.AddDate(0, 0, 2)
	auditEntries = append(auditEntries,
		createAuditEntry(userID, "segment.created", "segment", "premium_users", "Premium Users segment created", day3.Add(9*time.Hour)),
		createAuditEntry(userID, "segment.created", "segment", "beta_testers", "Beta Testers segment created", day3.Add(10*time.Hour)),
		createAuditEntry(userID, "segment.updated", "segment", "premium_users", "Updated segment rules for premium users", day3.Add(11*time.Hour)),
		createAuditEntry(userID, "api_key.created", "api_key", "staging-key-1", "Staging API key created", day3.Add(14*time.Hour)),
	)

	// Day 4: More flags and experiments (5 entries)
	day4 := baseTime.AddDate(0, 0, 3)
	auditEntries = append(auditEntries,
		createAuditEntry(userID, "flag.created", "flag", "feature_announcements", "Feature Announcements flag created", day4.Add(9*time.Hour)),
		createAuditEntry(userID, "flag.created", "flag", "experimental_search", "Experimental Search flag created", day4.Add(10*time.Hour)),
		createAuditEntry(userID, "experiment.created", "experiment", "dashboard_layout_test", "Dashboard layout A/B test created", day4.Add(11*time.Hour)),
		createAuditEntry(userID, "flag.toggled", "flag", "feature_announcements", "Enabled feature announcements in staging", day4.Add(14*time.Hour)),
		createAuditEntry(userID, "experiment.updated", "experiment", "dashboard_layout_test", "Updated experiment traffic allocation", day4.Add(15*time.Hour)),
	)

	// Day 5: Flag toggles and configuration updates (6 entries)
	day5 := baseTime.AddDate(0, 0, 4)
	auditEntries = append(auditEntries,
		createAuditEntry(userID, "flag.toggled", "flag", "new_dashboard", "Enabled new_dashboard in production with 10% rollout", day5.Add(9*time.Hour)),
		createAuditEntry(userID, "flag.updated", "flag", "api_rate_limit", "Updated rate limit for premium tier", day5.Add(10*time.Hour)),
		createAuditEntry(userID, "flag.toggled", "flag", "experimental_search", "Enabled experimental search in staging", day5.Add(13*time.Hour)),
		createAuditEntry(userID, "segment.updated", "segment", "beta_testers", "Added account age requirement to beta testers", day5.Add(14*time.Hour)),
		createAuditEntry(userID, "flag.updated", "flag", "experimental_search", "Configured targeting rules for search algorithm", day5.Add(15*time.Hour)),
		createAuditEntry(userID, "flag.toggled", "flag", "feature_announcements", "Disabled feature announcements in production", day5.Add(16*time.Hour)),
	)

	// Day 6: Weekend - lighter activity (2 entries)
	day6 := baseTime.AddDate(0, 0, 5)
	auditEntries = append(auditEntries,
		createAuditEntry(userID, "flag.created", "flag", "maintenance_mode", "Maintenance Mode flag created", day6.Add(10*time.Hour)),
		createAuditEntry(userID, "flag.updated", "flag", "maintenance_mode", "Configured maintenance mode settings", day6.Add(11*time.Hour)),
	)

	// Day 8: Heavy activity day (7 entries) - skipping day 7 (weekend)
	day8 := baseTime.AddDate(0, 0, 7)
	auditEntries = append(auditEntries,
		createAuditEntry(userID, "api_key.created", "api_key", "prod-key-1", "Production API key created", day8.Add(9*time.Hour)),
		createAuditEntry(userID, "segment.created", "segment", "enterprise_clients", "Enterprise Clients segment created", day8.Add(10*time.Hour)),
		createAuditEntry(userID, "experiment.updated", "experiment", "dashboard_layout_test", "Updated experiment targeting", day8.Add(11*time.Hour)),
		createAuditEntry(userID, "flag.created", "flag", "theme_config", "Theme Configuration flag created", day8.Add(12*time.Hour)),
		createAuditEntry(userID, "flag.updated", "flag", "theme_config", "Configured theme colors and settings", day8.Add(13*time.Hour)),
		createAuditEntry(userID, "flag.toggled", "flag", "theme_config", "Enabled theme config in all environments", day8.Add(14*time.Hour)),
		createAuditEntry(userID, "experiment.updated", "experiment", "dashboard_layout_test", "Experiment started - traffic split activated", day8.Add(16*time.Hour)),
	)

	// Day 9: Flag management and refinement (5 entries)
	day9 := baseTime.AddDate(0, 0, 8)
	auditEntries = append(auditEntries,
		createAuditEntry(userID, "flag.updated", "flag", "new_dashboard", "Increased rollout to 25% in production", day9.Add(9*time.Hour)),
		createAuditEntry(userID, "flag.updated", "flag", "feature_announcements", "Updated announcement message content", day9.Add(11*time.Hour)),
		createAuditEntry(userID, "flag.toggled", "flag", "maintenance_mode", "Temporarily enabled maintenance mode", day9.Add(12*time.Hour)),
		createAuditEntry(userID, "flag.toggled", "flag", "maintenance_mode", "Disabled maintenance mode", day9.Add(12*time.Hour+30*time.Minute)),
		createAuditEntry(userID, "segment.updated", "segment", "enterprise_clients", "Refined enterprise client criteria", day9.Add(15*time.Hour)),
	)

	// Day 10: Experiment management and new features (6 entries)
	day10 := baseTime.AddDate(0, 0, 9)
	auditEntries = append(auditEntries,
		createAuditEntry(userID, "experiment.created", "experiment", "search_algorithm_test", "Search algorithm performance test created", day10.Add(9*time.Hour)),
		createAuditEntry(userID, "experiment.updated", "experiment", "search_algorithm_test", "Configured search experiment variants", day10.Add(10*time.Hour)),
		createAuditEntry(userID, "flag.updated", "flag", "theme_config", "Updated theme configuration values", day10.Add(11*time.Hour)),
		createAuditEntry(userID, "flag.created", "flag", "mobile_app_upsell", "Mobile App Upsell flag created", day10.Add(13*time.Hour)),
		createAuditEntry(userID, "flag.updated", "flag", "mobile_app_upsell", "Configured mobile upsell targeting", day10.Add(14*time.Hour)),
		createAuditEntry(userID, "flag.toggled", "flag", "mobile_app_upsell", "Enabled mobile upsell in staging", day10.Add(15*time.Hour)),
	)

	// Day 11: Feature additions and testing (5 entries)
	day11 := baseTime.AddDate(0, 0, 10)
	auditEntries = append(auditEntries,
		createAuditEntry(userID, "flag.created", "flag", "support_chat_widget", "Support Chat Widget flag created", day11.Add(9*time.Hour)),
		createAuditEntry(userID, "flag.created", "flag", "max_upload_size", "Max Upload Size flag created", day11.Add(10*time.Hour)),
		createAuditEntry(userID, "flag.updated", "flag", "support_chat_widget", "Configured chat widget settings", day11.Add(11*time.Hour)),
		createAuditEntry(userID, "flag.updated", "flag", "max_upload_size", "Set different upload limits per environment", day11.Add(13*time.Hour)),
		createAuditEntry(userID, "flag.updated", "flag", "new_dashboard", "Updated dashboard feature rollout to 50%", day11.Add(14*time.Hour)),
	)

	// Day 12: Configuration refinements (4 entries)
	day12 := baseTime.AddDate(0, 0, 11)
	auditEntries = append(auditEntries,
		createAuditEntry(userID, "flag.toggled", "flag", "experimental_search", "Enabled experimental search in production with 25% rollout", day12.Add(9*time.Hour)),
		createAuditEntry(userID, "flag.toggled", "flag", "support_chat_widget", "Enabled support chat in staging", day12.Add(11*time.Hour)),
		createAuditEntry(userID, "experiment.updated", "experiment", "search_algorithm_test", "Updated experiment results and metrics", day12.Add(12*time.Hour)),
		createAuditEntry(userID, "experiment.updated", "experiment", "dashboard_layout_test", "Recorded conversion metrics", day12.Add(16*time.Hour)),
	)

	// Day 13: Beta features and optimization (4 entries)
	day13 := baseTime.AddDate(0, 0, 12)
	auditEntries = append(auditEntries,
		createAuditEntry(userID, "flag.created", "flag", "beta_features", "Beta Features flag created", day13.Add(9*time.Hour)),
		createAuditEntry(userID, "flag.updated", "flag", "beta_features", "Configured beta feature access rules", day13.Add(10*time.Hour)),
		createAuditEntry(userID, "segment.updated", "segment", "beta_testers", "Updated beta tester segment criteria", day13.Add(11*time.Hour)),
		createAuditEntry(userID, "flag.toggled", "flag", "feature_announcements", "Re-enabled feature announcements with updated content", day13.Add(15*time.Hour)),
	)

	// Day 14: Recent activity and final touches (6 entries)
	day14 := baseTime.AddDate(0, 0, 13)
	auditEntries = append(auditEntries,
		createAuditEntry(userID, "flag.created", "flag", "welcome_message", "Welcome Message flag created", day14.Add(9*time.Hour)),
		createAuditEntry(userID, "flag.updated", "flag", "welcome_message", "Configured personalized welcome messages", day14.Add(9*time.Hour+30*time.Minute)),
		createAuditEntry(userID, "flag.updated", "flag", "api_rate_limit", "Adjusted rate limits based on usage patterns", day14.Add(10*time.Hour)),
		createAuditEntry(userID, "flag.updated", "flag", "theme_config", "Updated theme colors for better accessibility", day14.Add(11*time.Hour)),
		createAuditEntry(userID, "flag.toggled", "flag", "mobile_app_upsell", "Enabled mobile upsell in production", day14.Add(14*time.Hour)),
		createAuditEntry(userID, "experiment.updated", "experiment", "search_algorithm_test", "Completed search experiment - analyzed results", day14.Add(16*time.Hour)),
	)

	// Add recent entries for today (3 entries)
	recent := time.Now().Add(-3 * time.Hour)
	auditEntries = append(auditEntries,
		createAuditEntry(userID, "experiment.updated", "experiment", "dashboard_layout_test", "Updated experiment metrics and analysis", recent),
		createAuditEntry(userID, "flag.toggled", "flag", "beta_features", "Enabled beta features for qualified users", recent.Add(1*time.Hour)),
		createAuditEntry(userID, "flag.updated", "flag", "new_dashboard", "Increased production rollout to 75% based on positive metrics", recent.Add(2*time.Hour)),
	)

	// Use full upsert (not $setOnInsert) so data gets refreshed on each deploy
	for i, entry := range auditEntries {
		// Create a unique composite filter for each audit entry
		filter := bson.M{
			"user_id":     entry.UserID,
			"resource":    entry.Resource,
			"resource_id": entry.ResourceID,
			"action":      entry.Action,
			"timestamp":   entry.Timestamp,
		}
		// Use $set for full upsert - updates existing entries with fresh data
		update := bson.M{"$set": entry}
		opts := options.Update().SetUpsert(true)

		_, err := auditCollection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			log.Printf("Failed to upsert audit entry %d: %v", i, err)
		}
	}

	log.Printf("Upserted %d audit log entries across 14 days", len(auditEntries))
}

// Helper functions

func timePtr(t time.Time) *time.Time {
	return &t
}

func createBooleanFlag(key, name, description string, devValue, prodValue bool, prodRollout int, createdAt time.Time) models.Flag {
	return models.Flag{
		ID:          primitive.NewObjectID(),
		Key:         key,
		Name:        name,
		Description: description,
		IsActive:    true,
		Type:        "boolean",
		Environments: map[string]models.FlagEnvironment{
			"development": {
				Enabled:        true,
				DefaultValue:   devValue,
				RolloutPercent: 100,
				TargetingRules: []models.TargetingRule{},
			},
			"staging": {
				Enabled:        true,
				DefaultValue:   devValue,
				RolloutPercent: int(math.Max(float64(prodRollout*2), 100)),
				TargetingRules: []models.TargetingRule{},
			},
			"production": {
				Enabled:        prodValue || prodRollout > 0,
				DefaultValue:   prodValue,
				RolloutPercent: prodRollout,
				TargetingRules: []models.TargetingRule{},
			},
		},
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
}

func createStringFlag(key, name, description, devValue, prodValue string, createdAt time.Time) models.Flag {
	return models.Flag{
		ID:          primitive.NewObjectID(),
		Key:         key,
		Name:        name,
		Description: description,
		IsActive:    true,
		Type:        "string",
		Environments: map[string]models.FlagEnvironment{
			"development": {
				Enabled:        true,
				DefaultValue:   devValue,
				RolloutPercent: 100,
				TargetingRules: []models.TargetingRule{},
			},
			"staging": {
				Enabled:        len(devValue) > 0,
				DefaultValue:   devValue,
				RolloutPercent: 100,
				TargetingRules: []models.TargetingRule{},
			},
			"production": {
				Enabled:        len(prodValue) > 0,
				DefaultValue:   prodValue,
				RolloutPercent: 100,
				TargetingRules: []models.TargetingRule{},
			},
		},
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
}

func createNumberFlag(key, name, description string, devValue, prodValue float64, createdAt time.Time) models.Flag {
	return models.Flag{
		ID:          primitive.NewObjectID(),
		Key:         key,
		Name:        name,
		Description: description,
		IsActive:    true,
		Type:        "number",
		Environments: map[string]models.FlagEnvironment{
			"development": {
				Enabled:        true,
				DefaultValue:   devValue,
				RolloutPercent: 100,
				TargetingRules: []models.TargetingRule{},
			},
			"staging": {
				Enabled:        true,
				DefaultValue:   (devValue + prodValue) / 2,
				RolloutPercent: 100,
				TargetingRules: []models.TargetingRule{},
			},
			"production": {
				Enabled:        true,
				DefaultValue:   prodValue,
				RolloutPercent: 100,
				TargetingRules: []models.TargetingRule{},
			},
		},
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
}

func createAuditEntry(userID primitive.ObjectID, action, resource, resourceID, description string, timestamp time.Time) models.AuditLogEntry {
	return models.AuditLogEntry{
		ID:         primitive.NewObjectID(),
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		UserID:     userID,
		UserEmail:  "demo@workermill.com",
		Changes: map[string]interface{}{
			"description": description,
			"timestamp":   timestamp.Format(time.RFC3339),
		},
		Metadata: map[string]interface{}{
			"source":      "seed_script",
			"environment": "seeded_data",
		},
		Timestamp: timestamp,
	}
}
