package authentication

import (
	"athena/src/api/errs"
	authenticationrequests "athena/src/api/http/requests/authentication"
	"athena/src/api/http/response"
	"athena/src/config"
	"athena/src/models"
	"athena/src/models/consts"
	"athena/src/pkg/logger"
	"athena/src/pkg/validator"
	"athena/src/services/authentication"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	LoginService authentication.ILoginService
	AccessToken  authentication.IAccessTokenService
}

func (controller *LoginController) LoginViaPassword(c *gin.Context) {
	// bind the incoming request to json
	var req authenticationrequests.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.LogJSONBindError(c, err)
		response.Api(c).SetLog().Send()
		return
	}

	// validate request
	if err := validator.Validate(&req, c.GetString("locale")); err != nil {
		logger.LogValidationError(c, err)
		response.Api(c).SetErrors(err).SetLog().Send()
		return
	}

	ctx := context.WithValue(context.Background(), "req", &req)
	ctx = context.WithValue(ctx, "request-ip", c.GetString("request-ip"))
	ctx = context.WithValue(ctx, "request-user-agent", c.GetHeader("User-Agent"))
	nationalCode := req.NationalIdentityCode
	if nationalCode != "" {
		ctx = context.WithValue(ctx, consts.NationalIdentityCode, nationalCode)
	}
	ctx = context.WithValue(ctx, consts.RequestUuid, c.GetString("request-uuid"))
	ctx = context.WithValue(ctx, consts.OwnerType, models.UserRole)

	jwtDTO, err := controller.LoginService.AdminLoginViaPassword(ctx)
	if err != nil {
		logger.LogServiceV2(c, "failed to login via password (admin)", controller, err)
		response.Api(c).SetMessage(err.Error()).SetLog().Send()
		return
	}

	expSeconds, err := strconv.Atoi(config.GetInstance().Get("COOKIE_EXPIRATION"))
	if err != nil {
		logger.LogAToIError(c, err)
		response.Api(c).SetMessage(errs.SomeThingWentWrong.Error()).SetLog().Send()
		return
	}

	expiration := time.Now().Add(time.Duration(expSeconds) * time.Second)

	cookie := &http.Cookie{
		Name:     config.GetInstance().Get("COOKIE_NAME"),
		Value:    jwtDTO.RefreshTokenString,
		Path:     "/",
		Domain:   config.GetInstance().Get("COOKIE_HOST"),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Expires:  expiration,
	}

	http.SetCookie(c.Writer, cookie)
	c.Header("Set-Cookie", cookie.String())

	response.Api(c).
		SetMessage("welcome").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"access_tokens": jwtDTO,
		}).
		SetLog().Send()
}
