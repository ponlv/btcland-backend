package middleware

import (
	"os"
	"strings"

	"api/internal/jwt"
	"api/internal/response"
	"api/schema/usercol"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type contextKey string

// AuthMiddleware which authorizes the external client.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authParts := strings.Split(c.GetHeader("authorization"), " ")
		if len(authParts) == 2 {
			authData := authParts[1]

			if err := validateToken(c, authData); err != nil {
				code := response.ErrorResponse(err.Error())
				c.JSON(code.Code, code)
				c.Abort()
				return
			}

			// Get user and set in context
			user, err := CurrentUser(c)
			if err != nil {
				code := response.ErrorResponse(err.Error())
				c.JSON(code.Code, code)
				c.Abort()
				return
			}

			if user != nil {
				c.Set("current_user", user)
			}

			c.Next()
		} else {
			code := response.ErrorResponse("Unauthorized")
			c.JSON(code.Code, code)
			c.Abort()
			return
		}
	}
}

func validateToken(c *gin.Context, token string) error {
	// validate token
	claim, err := jwt.VerifyJWTToken(os.Getenv("KEY_API_KEY"), token)
	if err != nil {
		return err
	}

	// Add claim data to gin context
	c.Set("user_id", claim.UserId)
	return nil
}

func CurrentUser(c *gin.Context) (*usercol.User, error) {
	authParts := strings.Split(c.GetHeader("authorization"), " ")
	if len(authParts) != 2 {
		return nil, nil
	}

	authData := authParts[1]

	if err := validateToken(c, authData); err != nil {
		return nil, nil
	}

	userId := c.GetString("user_id")
	user, err := usercol.FindWithUserID(c.Request.Context(), userId)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	// Check if user is deleted
	if user != nil && user.IsDelete {
		return nil, errors.New("account has been deleted")
	}

	return user, nil
}
