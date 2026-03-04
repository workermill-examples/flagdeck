package handlers

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/workermill-examples/flagdeck/api/internal/middleware"
	"github.com/workermill-examples/flagdeck/api/internal/models"
)

type AuthHandler struct {
	UserCollection *mongo.Collection
	JWTSecret      string
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func NewAuthHandler(userCollection *mongo.Collection, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		UserCollection: userCollection,
		JWTSecret:      jwtSecret,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return middleware.NewValidationError("Invalid request body")
	}

	// Validate required fields
	if req.Email == "" || req.Name == "" || req.Password == "" {
		return middleware.NewValidationError("Email, name, and password are required")
	}

	// Validate email format (basic check)
	if !strings.Contains(req.Email, "@") {
		return middleware.NewValidationError("Invalid email format")
	}

	// Check password length
	if len(req.Password) < 6 {
		return middleware.NewValidationError("Password must be at least 6 characters long")
	}

	// Check if user already exists
	var existingUser models.User
	err := h.UserCollection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		return middleware.NewConflictError("User already exists with this email")
	}
	if err != mongo.ErrNoDocuments {
		log.Printf("Database error while checking existing user: %v", err)
		return middleware.NewDatabaseError("Failed to check user existence")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return middleware.NewInternalError("Failed to process password")
	}

	// Create new user with default role "viewer" (no role field accepted in request)
	now := time.Now()
	user := models.User{
		ID:        primitive.NewObjectID(),
		Email:     req.Email,
		Name:      req.Name,
		Role:      "viewer", // Default role as specified in ticket
		Password:  string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Insert user into database
	_, err = h.UserCollection.InsertOne(context.Background(), user)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return middleware.NewDatabaseError("Failed to create user")
	}

	// Generate tokens for auto-login
	accessToken, refreshToken, err := middleware.GenerateTokens(
		user.ID,
		user.Email,
		user.Name,
		user.Role,
		h.JWTSecret,
	)
	if err != nil {
		log.Printf("Failed to generate tokens: %v", err)
		return middleware.NewInternalError("Failed to generate authentication tokens")
	}

	// Return tokens according to spec
	return c.Status(fiber.StatusCreated).JSON(AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900, // 15 minutes in seconds
		TokenType:    "Bearer",
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return middleware.NewValidationError("Invalid request body")
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" {
		return middleware.NewValidationError("Email and password are required")
	}

	// Find user by email
	var user models.User
	err := h.UserCollection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return middleware.NewUnauthorizedError("Invalid email or password")
		}
		log.Printf("Database error while finding user: %v", err)
		return middleware.NewDatabaseError("Failed to authenticate user")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return middleware.NewUnauthorizedError("Invalid email or password")
	}

	// Generate tokens
	accessToken, refreshToken, err := middleware.GenerateTokens(
		user.ID,
		user.Email,
		user.Name,
		user.Role,
		h.JWTSecret,
	)
	if err != nil {
		log.Printf("Failed to generate tokens: %v", err)
		return middleware.NewInternalError("Failed to generate authentication tokens")
	}

	// Return tokens according to spec
	return c.JSON(AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900, // 15 minutes in seconds
		TokenType:    "Bearer",
	})
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return middleware.NewValidationError("Invalid request body")
	}

	if req.RefreshToken == "" {
		return middleware.NewValidationError("Refresh token is required")
	}

	// Validate refresh token
	claims, err := middleware.ValidateRefreshToken(req.RefreshToken, h.JWTSecret)
	if err != nil {
		return middleware.NewUnauthorizedError("Invalid refresh token")
	}

	// Get user to ensure they still exist and get current info
	userIDObj, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		return middleware.NewUnauthorizedError("Invalid user ID in token")
	}

	var user models.User
	err = h.UserCollection.FindOne(context.Background(), bson.M{"_id": userIDObj}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return middleware.NewUnauthorizedError("User not found")
		}
		log.Printf("Database error while fetching user for refresh: %v", err)
		return middleware.NewDatabaseError("Failed to refresh token")
	}

	// Generate new tokens
	accessToken, refreshToken, err := middleware.GenerateTokens(
		user.ID,
		user.Email,
		user.Name,
		user.Role,
		h.JWTSecret,
	)
	if err != nil {
		log.Printf("Failed to generate tokens: %v", err)
		return middleware.NewInternalError("Failed to generate authentication tokens")
	}

	return c.JSON(AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900, // 15 minutes in seconds
		TokenType:    "Bearer",
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// For stateless JWT, logout is typically handled client-side by removing tokens
	// In a production system, you might want to maintain a blacklist of tokens
	// For now, we'll just return success
	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	// Get user from context (set by JWT middleware)
	user, ok := c.Locals("user").(middleware.UserContext)
	if !ok {
		return middleware.NewUnauthorizedError("User context not found")
	}

	// Get full user details from database to return created_at
	var fullUser models.User
	err := h.UserCollection.FindOne(context.Background(), bson.M{"_id": user.ID}).Decode(&fullUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return middleware.NewUnauthorizedError("User not found")
		}
		log.Printf("Database error while fetching user details: %v", err)
		return middleware.NewDatabaseError("Failed to fetch user details")
	}

	return c.JSON(UserResponse{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: fullUser.CreatedAt,
	})
}
