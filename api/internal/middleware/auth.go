package middleware

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/workermill-examples/flagdeck/api/internal/models"
)

type UserContext struct {
	ID    primitive.ObjectID `json:"id"`
	Email string             `json:"email"`
	Name  string             `json:"name"`
	Role  string             `json:"role"`
}

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	Type   string `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

type AuthConfig struct {
	JWTSecret      string
	UserCollection *mongo.Collection
}

func AuthenticateJWT(config AuthConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return NewUnauthorizedError("Authorization header required")
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			return NewUnauthorizedError("Invalid authorization header format")
		}

		tokenString := bearerToken[1]
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, NewUnauthorizedError("Invalid token signing method")
			}
			return []byte(config.JWTSecret), nil
		})

		if err != nil {
			log.Printf("JWT parsing error: %v", err)
			return NewUnauthorizedError("Invalid token")
		}

		if !token.Valid {
			return NewUnauthorizedError("Invalid token")
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			return NewUnauthorizedError("Invalid token claims")
		}

		if claims.Type != "access" {
			return NewUnauthorizedError("Invalid token type")
		}

		userIDObj, err := primitive.ObjectIDFromHex(claims.UserID)
		if err != nil {
			log.Printf("Invalid user ID in token: %v", err)
			return NewUnauthorizedError("Invalid user ID")
		}

		var user models.User
		err = config.UserCollection.FindOne(context.Background(), bson.M{"_id": userIDObj}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return NewUnauthorizedError("User not found")
			}
			log.Printf("Database error while fetching user: %v", err)
			return NewDatabaseError("Failed to verify user")
		}

		userCtx := UserContext{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Role:  user.Role,
		}

		c.Locals("user", userCtx)
		return c.Next()
	}
}

func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(UserContext)
		if !ok {
			return NewUnauthorizedError("User context not found")
		}

		if user.Role != requiredRole && requiredRole != "" {
			return NewForbiddenError("Insufficient permissions")
		}

		return c.Next()
	}
}

func RequireRoles(requiredRoles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(UserContext)
		if !ok {
			return NewUnauthorizedError("User context not found")
		}

		for _, role := range requiredRoles {
			if user.Role == role {
				return c.Next()
			}
		}

		return NewForbiddenError("Insufficient permissions")
	}
}

func GenerateTokens(userID primitive.ObjectID, email, name, role, jwtSecret string) (accessToken, refreshToken string, err error) {
	now := time.Now()

	accessClaims := &JWTClaims{
		UserID: userID.Hex(),
		Email:  email,
		Name:   name,
		Role:   role,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "flagdeck-api",
		},
	}

	refreshClaims := &JWTClaims{
		UserID: userID.Hex(),
		Email:  email,
		Name:   name,
		Role:   role,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "flagdeck-api",
		},
	}

	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessToken, err = accessTokenObj.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", "", err
	}

	refreshToken, err = refreshTokenObj.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func ValidateRefreshToken(tokenString, jwtSecret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, NewUnauthorizedError("Invalid token signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, NewUnauthorizedError("Invalid token")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, NewUnauthorizedError("Invalid token claims")
	}

	if claims.Type != "refresh" {
		return nil, NewUnauthorizedError("Invalid token type")
	}

	return claims, nil
}
