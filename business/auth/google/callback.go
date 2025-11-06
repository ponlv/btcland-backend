package google

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"api/business/auth/login"
	"api/internal/jwt"
	bsonutil "api/internal/mongodb/utils"
	"api/internal/plog"
	"api/internal/response"
	"api/internal/timer"
	"api/schema/usercol"
	"api/schema/userdevicecol"
	"api/schema/usersessioncol"
	"api/services/oauth2"
	"api/services/oauth2/google"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type googleUser struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func Callback() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][auth][google][callback]")

	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Validate OAuth state & code
		code, deviceID, deviceName, browserName, platform, ok := validateOAuthCallbackParams(c, logger)
		if !ok {
			return
		}

		// Exchange code for token & get user info
		userInfo, err := getUserInfoFromGoogle(ctx, code, logger)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
			c.Abort()
			return
		}

		// Find or create user
		user, err := findOrCreateUser(ctx, userInfo, logger)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse("failed to get or create user"))
			c.Abort()
			return
		}

		// Create session
		_, err = createUserSession(c, user.ID.(primitive.ObjectID), deviceName, browserName, platform, logger)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse("session creation failed"))
			c.Abort()
			return
		}

		// Update device info
		err = updateDeviceInfo(ctx, c, user.GetIDString(), deviceID, deviceName, platform, logger)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse("device update failed"))
			c.Abort()
			return
		}

		// Generate JWT & Redirect
		accessToken, _ := jwt.GenerateJWTToken(
			os.Getenv("KEY_API_KEY"),
			user.GetIDString(),
			"dipnet",
			"fanxipan",
			login.JWT_LIFETIME,
		)

		redirectURL := fmt.Sprintf("%s?access_token=%s", google.OAuthConfig.ClientRedirectURL, accessToken)

		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
	}
}

func validateOAuthCallbackParams(c *gin.Context, _ plog.Logger) (
	code,
	deviceID,
	deviceName,
	browserName,
	platform string,
	ok bool,
) {
	state := c.Query("state")
	if state != google.OAuthConfig.OAuthStateString {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("invalid state"))
		c.Abort()
		return
	}

	code = c.Query("code")
	if code == "" {
		err := response.ErrorResponse("code Not Found to provide AccessToken...")
		reason := c.Query("error_reason")
		if reason == "user_denied" {
			err.Message = errors.Wrap(err, "user has denied Permission..").Error()
		}
		c.JSON(http.StatusBadRequest, err)
		c.Abort()
		return
	}

	deviceName = c.DefaultQuery("device_name", "Unknown")
	browserName = c.DefaultQuery("browser_name", "Unknown")
	platform = c.DefaultQuery("platform", "Unknown")
	deviceID = c.Query("device_id")

	return code, deviceID, deviceName, browserName, platform, true
}

func getUserInfoFromGoogle(ctx context.Context, code string, logger plog.Logger) (*googleUser, error) {
	token, err := google.OAuthConfig.Exchange(ctx, code)
	if err != nil {
		logger.Err(err).Msg("token exchange error")
		return nil, errors.New("failed to exchange token")
	}

	client := google.OAuthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
	if err != nil {
		logger.Err(err).Msg("user info error")
		return nil, errors.New("failed to get user info")
	}

	var userInfo googleUser
	if err = json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		logger.Err(err).Msg("failed to parse user info")
		return nil, errors.New("failed to parse user info")
	}

	_ = resp.Body.Close()

	return &userInfo, nil
}

func findOrCreateUser(ctx context.Context, info *googleUser, logger plog.Logger) (*usercol.User, error) {
	filter := bsonutil.BsonAdd(nil, "email", info.Email)

	user, err := usercol.FindWithCondition(ctx, filter, nil)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		logger.Err(err).Msg("database error")
		return nil, err
	}

	if user == nil {
		newUser := &usercol.User{
			Email:    info.Email,
			FullName: info.Name,
			Avatar:   info.Picture,
			OAuthProvider: &usercol.OAuthProvider{
				ProviderID:   info.ID,
				ProviderName: oauth2.Google,
			},
			UserType: usercol.UserTypeGeneral,
		}
		_, err = usercol.Create(ctx, newUser)
		if err != nil {
			logger.Err(err).Msg("failed to create user")
			return nil, err
		}
		user = newUser

	} else {
		// Check if user is deleted
		if user.IsDelete {
			return nil, errors.New("account has been deleted")
		}
	}

	return user, nil
}

func createUserSession(
	c *gin.Context,
	userID primitive.ObjectID,
	deviceName,
	browserName,
	platform string,
	logger plog.Logger,
) (*usersessioncol.UserSession, error) {
	session := &usersessioncol.UserSession{
		CreatedAt:   timer.Now(),
		ActiveAt:    timer.Now(),
		UserId:      userID.Hex(),
		IsDelete:    false,
		BrowserName: browserName,
		IP:          c.ClientIP(),
		DeviceName:  c.GetHeader("User-Agent"),
		Platform:    platform,
	}
	_, err := usersessioncol.Create(c.Request.Context(), session)
	if err != nil {
		logger.Err(err).Msg("create session failed")
		return nil, err
	}

	return session, nil
}

func updateDeviceInfo(ctx context.Context, c *gin.Context, userID, deviceID, deviceName, platform string, logger plog.Logger) error {
	if err := userdevicecol.DisableAllDevice(ctx, userID); err != nil {
		logger.Err(err).Msg("disable all devices failed")
		return err
	}

	device, err := userdevicecol.FindWithDeviceId(ctx, userID, deviceID)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		logger.Err(err).Msg("find device failed")
		return err
	}

	if device == nil {
		newDevice := &userdevicecol.UserDevice{
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			IP:          c.ClientIP(),
			UserId:      userID,
			DeviceID:    deviceID,
			DeviceName:  deviceName,
			Platform:    platform,
			IsEnable:    true,
			IsCurrent:   true,
			LastLoginAt: time.Now(),
		}
		_, err = userdevicecol.Create(ctx, newDevice)
		if err != nil {
			logger.Err(err).Msg("create device failed")
			return err
		}
	} else {
		device.IsCurrent = true
		device.LastLoginAt = timer.Now()
		_, err = userdevicecol.Update(ctx, device)
		if err != nil {
			logger.Err(err).Msg("update device failed")
			return err
		}
	}

	return nil
}
