package middleware

import (
	"net/http"
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

const (
	_initDataKey contextKey = "init-data"
)

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

func RequireKYC() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetString("user_id")
		if userId == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{
				"message": "Unauthorized",
			})
			return
		}

		user, err := usercol.FindWithUserID(c.Request.Context(), userId)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			code := response.ErrorResponse(err.Error())
			c.JSON(code.Code, code)
			c.Abort()
			return
		}

		if user == nil {
			code := response.ErrorResponse("user not found")
			c.JSON(code.Code, code)
			c.Abort()
			return
		}

		if !user.IsKYC {
			code := response.ErrorResponse("need to be finish kyc before continue")
			c.JSON(code.Code, code)
			c.Abort()
			return
		}

		c.Set("current_user", user)

		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exist := c.Get("current_user")
		if !exist {
			code := response.ErrorResponse("Unauthorized")
			c.JSON(code.Code, code)
			c.Abort()
			return
		}

		u, ok := user.(*usercol.User)
		if !ok {
			code := response.ErrorResponse("Unauthorized")
			c.JSON(code.Code, code)
			c.Abort()
			return
		}

		if u.UserType != "ADMIN" {
			code := response.ErrorResponse("Admin access required")
			c.JSON(code.Code, code)
			c.Abort()
			return
		}

		c.Next()
	}
}
