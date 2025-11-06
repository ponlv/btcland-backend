package login

import (
	"errors"
	"net/http"
	"os"

	"api/internal/jwt"
	"api/internal/plog"
	"api/internal/response"
	"api/internal/timer"
	"api/internal/utils"
	"api/schema/usercol"
	"api/schema/userdevicecol"
	"api/schema/usersessioncol"
	"api/types"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const JWT_LIFETIME = 7 * 24 * 60 * 60

type LoginRequestData struct {
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	DeviceId    string `json:"device_id"`
	DeviceName  string `json:"device_name"`
	DeviceToken string `json:"device_token"`
	Platform    string `json:"platform"`
}

type LoginResponseData struct {
	AccessToken string     `json:"access_token"`
	User        types.User `json:"user"`
}

func Login() gin.HandlerFunc {
	logger := plog.NewBizLogger("[auth][login]")

	return func(c *gin.Context) {
		// get data from request
		var req LoginRequestData
		err := c.BindJSON(&req)
		if err != nil {
			code := response.ErrorResponse(err.Error())
			c.JSON(code.Code, code)
			c.Abort()
			return
		}

		res, err := doLogin(c, req)
		if err != nil {
			logger.Err(err).Msg("failed to do login")
			code := response.ErrorResponse(err.Error())
			c.JSON(code.Code, code)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, response.SuccessResponse(res))
		return
	}
}

func doLogin(c *gin.Context, data LoginRequestData) (*LoginResponseData, error) {
	var user *usercol.User
	var err error
	ctx := c.Request.Context()

	if utils.IsEmailValid(data.Email) {
		// Email flow
		user, err = usercol.FindWithEmail(ctx, data.Email)
	} else if utils.IsValidPhoneNumber(data.Email) {
		// Phone flow
		_, errConv := utils.ConvertPhoneToInternational(data.Email)
		if errConv != nil {
			return nil, response.ErrorResponse("phone format is invalid")
		}

		user, err = usercol.FindUserByPhone(ctx, data.Email, true)
	} else {
		return nil, response.ErrorResponse("invalid email or phone number")
	}

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// check is delete
	if user.IsDelete == true {
		return nil, errors.New("USER_IS_DELETED")
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		return nil, errors.New("PASSWORD_NOT_MATCH")
	}

	// remove all token
	err = usersessioncol.RemoveAllSession(c.Request.Context(), user.GetIDString())
	if err != nil {
		return nil, err
	}

	// create session
	newSession := &usersessioncol.UserSession{
		CreatedAt:  timer.Now(),
		ActiveAt:   timer.Now(),
		UserId:     user.GetIDString(),
		IsDelete:   false,
		IP:         c.ClientIP(),
		DeviceName: c.GetHeader("User-Agent"),
		Platform:   data.Platform,
	}
	_, err = usersessioncol.Create(c.Request.Context(), newSession)
	if err != nil {
		return nil, err
	}

	err = userdevicecol.DisableAllDevice(c.Request.Context(), user.GetIDString())
	if err != nil {
		return nil, err
	}

	// find device if exists
	device, err := userdevicecol.FindWithDeviceId(c.Request.Context(), user.GetIDString(), data.DeviceId)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	// if device not exists, create new device
	if device == nil {
		newDevice := &userdevicecol.UserDevice{
			CreatedAt:   timer.Now(),
			UpdatedAt:   timer.Now(),
			IP:          c.ClientIP(),
			UserId:      user.GetIDString(),
			DeviceID:    data.DeviceId,
			DeviceName:  data.DeviceName,
			DeviceToken: data.DeviceToken,
			Platform:    data.Platform,
			IsEnable:    true,
			IsCurrent:   true,
			LastLoginAt: timer.Now(),
		}

		_, err = userdevicecol.Create(c.Request.Context(), newDevice)
		if err != nil {
			return nil, err
		}
	} else {

		device.IsCurrent = true
		device.LastLoginAt = timer.Now()

		_, err = userdevicecol.Update(c.Request.Context(), device)
		if err != nil {
			return nil, err
		}
	}

	apiToken, _ := jwt.GenerateJWTToken(
		os.Getenv("KEY_API_KEY"),
		user.GetIDString(),
		"",
		"fanxipan",
		JWT_LIFETIME,
	)

	return &LoginResponseData{
		User: types.User{
			ID:                      user.GetIDString(),
			CreatedAt:               user.CreatedAt,
			UpdatedAt:               user.UpdatedAt,
			CitizenPid:              user.CitizenPid,
			IdCardExpireDate:        user.IdCardExpireDate,
			DateOfIssue:             user.DateOfIssue,
			IssuingAuthority:        user.IssuingAuthority,
			NationalityCode:         user.NationalityCode,
			PermanentAddress:        user.PermanentAddress,
			PermanentVillageCode:    user.PermanentVillageCode,
			PermanentVillageText:    user.PermanentVillageText,
			PermanentDistrictCode:   user.PermanentDistrictCode,
			PermanentDistrictText:   user.PermanentDistrictText,
			PermanentCityCode:       user.PermanentCityCode,
			PermanentCityText:       user.PermanentCityText,
			LivingPlaceAddress:      user.LivingPlaceAddress,
			LivingPlaceVillageCode:  user.LivingPlaceVillageCode,
			LivingPlaceVillageText:  user.LivingPlaceVillageText,
			LivingPlaceDistrictCode: user.LivingPlaceDistrictCode,
			LivingPlaceDistrictText: user.LivingPlaceDistrictText,
			LivingPlaceCityCode:     user.LivingPlaceCityCode,
			LivingPlaceCityText:     user.LivingPlaceCityText,
			IdentifyLevel:           user.IdentifyLevel,
			IsKYC:                   user.IsKYC,
			IsShareInfo:             user.IsShareInfo,
			FullName:                user.FullName,
			Email:                   user.Email,
			UserType:                user.UserType.String(),
		},
		AccessToken: apiToken,
	}, nil
}
