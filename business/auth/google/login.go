package google

import (
	"net/http"
	"net/url"
	"strings"

	"api/internal/plog"
	"api/internal/response"
	"api/services/oauth2/google"

	"github.com/gin-gonic/gin"
)

func Login() func(c *gin.Context) {
	logger := plog.NewBizLogger("[business][auth][google][login]")
	type req struct {
		DeviceName  string `json:"device_name"`
		BrowserName string `json:"browser_name"`
		Platform    string `json:"platform" binding:"required"`
	}

	return func(c *gin.Context) {
		browserName := c.Query("browser_name")
		platform := c.Query("platform")

		URL, err := url.Parse(google.OAuthConfig.Endpoint.AuthURL)
		if err != nil {
			logger.Err(err).Msg("failed to parse auth URL")
			code := response.ErrorResponse("failed to parse auth URL")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		parameters := url.Values{}
		parameters.Add("client_id", google.OAuthConfig.ClientID)
		parameters.Add("scope", strings.Join(google.OAuthConfig.Scopes, " "))
		parameters.Add("redirect_uri", google.OAuthConfig.RedirectURL)
		parameters.Add("response_type", "code")
		parameters.Add("state", google.OAuthConfig.OAuthStateString)
		parameters.Add("browser_name", browserName)
		parameters.Add("platform", platform)
		URL.RawQuery = parameters.Encode()
		oauthURL := URL.String()

		// Return the OAuth URL instead of redirecting
		c.JSON(http.StatusOK, response.SuccessResponse(map[string]string{
			"oauth_url": oauthURL,
		}))
	}
}
