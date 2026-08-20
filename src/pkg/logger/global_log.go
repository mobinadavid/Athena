package logger

import (
	"athena/src/api/errs"
	"athena/src/models"
	"athena/src/models/consts"
	"athena/src/pkg/utils"
	"context"
	"errors"
	"os"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// getLoggerWithLevel decides what log level the problem should be logged with
// if error is given, it decides on log level, but if nil given as error, it will be logged in Warn level
func getLoggerWithLevel(err error, callerSkip int) (func(msg string, fields ...zap.Field), *errs.AppError) {
	logger := GetInstance().WithOptions(zap.AddCallerSkip(callerSkip))
	if err == nil {
		return logger.Warn, nil
	}

	var ae *errs.AppError
	if errors.As(err, &ae) {
		switch ae.SeverityLevel {
		case errs.SeverityLevelError:
			return logger.Error, ae
		case errs.SeverityLevelWarn:
			return logger.Warn, ae
		case errs.SeverityLevelInfo:
			return logger.Info, ae
		default:
			return logger.Error, ae
		}
	} else {
		// add more logic as needed
		if utils.CheckError(err, gorm.ErrRecordNotFound) {
			return logger.Warn, nil
		}

		return logger.Error, nil
	}
}

func LogJSONBindError(c *gin.Context, err error) {
	// Retrieve the authentication data from the context
	auth := utils.GetAuthData(c)

	// Prepare log fields
	logFields := []zap.Field{
		zap.Error(err),
	}

	// If user is authenticated, add session info to the log
	if auth.IsAuthenticated {
		logFields = append(logFields,
			zap.Uint(string(consts.OwnerId), auth.OwnerId),
			zap.String(string(consts.OwnerType), auth.OwnerType),
			zap.String("access-token-uuid", auth.TokenUUID),
		)
	}

	// Log the error with the appropriate details
	GetInstance().WithOptions(zap.AddCallerSkip(CallerSkipOne)).Error("Failed to bind JSON body", logFields...)
}

func LogCookieDoesNotExist(c *gin.Context, err error, key string) {
	// Retrieve the authentication data from the context
	auth := utils.GetAuthData(c)

	// Prepare log fields
	logFields := []zap.Field{
		zap.String("key", key),
		zap.Error(errors.New("cookie not exist")),
		zap.Error(err),
	}

	// If user is authenticated, add session info to the log
	if auth.IsAuthenticated {
		logFields = append(logFields,
			zap.Uint(string(consts.OwnerId), auth.OwnerId),
			zap.String(string(consts.OwnerType), auth.OwnerType),
			zap.String("access-token-uuid", auth.TokenUUID),
		)
	}

	// Log the error with the appropriate details
	GetInstance().WithOptions(zap.AddCallerSkip(CallerSkipOne)).Error("Failed to bind JSON body", logFields...)
}

func LogValidationError(c *gin.Context, errors map[string]string) {
	// Retrieve the authentication data from the context
	auth := utils.GetAuthData(c)

	// Prepare log fields
	logFields := []zap.Field{
		zap.Any("errors", errors),
	}

	// If user is authenticated, add session info to the log
	if auth.IsAuthenticated {
		logFields = append(logFields,
			zap.Uint(string(consts.OwnerId), auth.OwnerId),
			zap.String(string(consts.OwnerType), auth.OwnerType),
			zap.String("access-token-uuid", auth.TokenUUID),
		)
	}

	// Log the validation error with the appropriate details
	GetInstance().WithOptions(zap.AddCallerSkip(CallerSkipOne)).Error("Validation fail on refresh access token request", logFields...)
}

func LogQueryBuilderError(c *gin.Context) {
	// Retrieve the authentication data from the context
	auth := utils.GetAuthData(c)

	// Prepare log fields
	logFields := []zap.Field{
		zap.String("reason", "The query builder middleware does not assign to route"),
	}

	// If user is authenticated, add session info to the log
	if auth.IsAuthenticated {
		logFields = append(logFields,
			zap.Uint(string(consts.OwnerId), auth.OwnerId),
			zap.String(string(consts.OwnerType), auth.OwnerType),
			zap.String("access-token-uuid", auth.TokenUUID),
		)
	}

	// Log the error with the appropriate details
	GetInstance().WithOptions(zap.AddCallerSkip(CallerSkipOne)).Error("Cannot find query parameters builder", logFields...)
}

func LogInvalidMultiPart(c *gin.Context, err error) {
	// Retrieve the authentication data from the context
	auth := utils.GetAuthData(c)

	// Prepare log fields
	logFields := []zap.Field{
		zap.String("reason", "the multi part data should be valid"),
		zap.Error(err),
	}

	// If user is authenticated, add session info to the log
	if auth.IsAuthenticated {
		logFields = append(logFields,
			zap.Uint(string(consts.OwnerId), auth.OwnerId),
			zap.String(string(consts.OwnerType), auth.OwnerType),
			zap.String("access-token-uuid", auth.TokenUUID),
		)
	}

	// Log the error with the appropriate details
	GetInstance().WithOptions(zap.AddCallerSkip(CallerSkipOne)).Error("invalid form data", logFields...)
}

func LogBucketMissed(c *gin.Context) {
	// Retrieve the authentication data from the context
	auth := utils.GetAuthData(c)

	// Prepare log fields
	logFields := []zap.Field{
		zap.String("reason", "the bucket or folder is missed"),
	}

	// If user is authenticated, add session info to the log
	if auth.IsAuthenticated {
		logFields = append(logFields,
			zap.Uint(string(consts.OwnerId), auth.OwnerId),
			zap.String(string(consts.OwnerType), auth.OwnerType),
			zap.String("access-token-uuid", auth.TokenUUID),
		)
	}

	// Log the error with the appropriate details
	GetInstance().WithOptions(zap.AddCallerSkip(CallerSkipOne)).Error("failed to get bucket or folder", logFields...)
}

func LogRestyFailed(c *gin.Context, resp *resty.Response) {
	// Retrieve the authentication data from the context
	auth := utils.GetAuthData(c)

	// Prepare log fields
	logFields := []zap.Field{
		zap.String("reason", "return an un success response"),
		zap.Int("status-code", resp.StatusCode()),
		zap.String("response", resp.String()),
	}

	// If user is authenticated, add session info to the log
	if auth.IsAuthenticated {
		logFields = append(logFields,
			zap.Uint(string(consts.OwnerId), auth.OwnerId),
			zap.String(string(consts.OwnerType), auth.OwnerType),
			zap.String("access-token-uuid", auth.TokenUUID),
		)
	}

	// Log the error with the appropriate details
	GetInstance().WithOptions(zap.AddCallerSkip(CallerSkipOne)).Error("failed to get bucket or folder", logFields...)
}

func LogValueNotExist(c *gin.Context, key string) {
	// Retrieve the authentication data from the context
	auth := utils.GetAuthData(c)

	// Prepare log fields
	logFields := []zap.Field{
		zap.String("reason", "the value with given key does not exist"),
		zap.String("key", key),
	}

	// If user is authenticated, add session info to the log
	if auth.IsAuthenticated {
		logFields = append(logFields,
			zap.Uint(string(consts.OwnerId), auth.OwnerId),
			zap.String(string(consts.OwnerType), auth.OwnerType),
			zap.String("access-token-uuid", auth.TokenUUID),
		)
	}

	// Log the error with the appropriate details
	GetInstance().WithOptions(zap.AddCallerSkip(CallerSkipOne)).Error("failed to get query param", logFields...)
}

func LogParseUUIDError(c *gin.Context, err error) {
	// Retrieve the authentication data from the context
	auth := utils.GetAuthData(c)

	// Prepare log fields
	logFields := []zap.Field{
		zap.String("reason", "Not a valid UUID"),
		zap.Error(err),
	}

	// If user is authenticated, add session info to the log
	if auth.IsAuthenticated {
		logFields = append(logFields,
			zap.Uint(string(consts.OwnerId), auth.OwnerId),
			zap.String(string(consts.OwnerType), auth.OwnerType),
			zap.String("access-token-uuid", auth.TokenUUID),
		)
	}

	// Log the error with the appropriate details
	GetInstance().WithOptions(zap.AddCallerSkip(CallerSkipOne)).Error("Cannot parse UUID", logFields...)
}

func LogAToIError(c *gin.Context, err error) {
	// Retrieve the authentication data from the context
	auth := utils.GetAuthData(c)

	// Prepare log fields
	logFields := []zap.Field{
		zap.Error(err),
	}

	// If user is authenticated, add session info to the log
	if auth.IsAuthenticated {
		logFields = append(logFields,
			zap.Uint(string(consts.OwnerId), auth.OwnerId),
			zap.String(string(consts.OwnerType), auth.OwnerType),
			zap.String("access-token-uuid", auth.TokenUUID),
		)
	}

	// Log the error with the appropriate details
	GetInstance().WithOptions(zap.AddCallerSkip(CallerSkipOne)).Error("Cannot cast string to int", logFields...)
}

func LogFailToCast(c *gin.Context, err error, castTo string) {
	// Retrieve the authentication data from the context
	auth := utils.GetAuthData(c)

	// Prepare log fields
	logFields := []zap.Field{
		zap.Error(err),
	}

	// If user is authenticated, add session info to the log
	if auth.IsAuthenticated {
		logFields = append(logFields,
			zap.Uint(string(consts.OwnerId), auth.OwnerId),
			zap.String(string(consts.OwnerType), auth.OwnerType),
			zap.String("access-token-uuid", auth.TokenUUID),
		)
	}

	// Log the error with the appropriate details
	GetInstance().WithOptions(zap.AddCallerSkip(CallerSkipOne)).Error("Failed casting to "+castTo, logFields...)
}

func LogRedirect(c *gin.Context, statusCode int, url string) {
	// Retrieve the authentication data from the context
	auth := utils.GetAuthData(c)

	// Prepare log fields
	logFields := []zap.Field{
		zap.Int("status-code", statusCode),
		zap.String("url", url),
	}

	// If user is authenticated, add session info to the log
	if auth.IsAuthenticated {
		logFields = append(logFields,
			zap.Uint(string(consts.OwnerId), auth.OwnerId),
			zap.String(string(consts.OwnerType), auth.OwnerType),
			zap.String("access-token-uuid", auth.TokenUUID),
		)
	}

	// Log the error with the appropriate details
	GetInstance().WithOptions(zap.AddCallerSkip(CallerSkipOne)).Error("redirecting...", logFields...)
}

func LogService(c *gin.Context, err error, message, service string) {
	// Retrieve the authentication data from the context
	auth := utils.GetAuthData(c)

	// Prepare log fields
	logFields := []zap.Field{
		zap.String("service", service),
		zap.Error(err),
	}

	// If user is authenticated, add session info to the log
	if auth.IsAuthenticated {
		logFields = append(logFields,
			zap.Uint(string(consts.OwnerId), auth.OwnerId),
			zap.String(string(consts.OwnerType), auth.OwnerType),
			zap.String("access-token-uuid", auth.TokenUUID),
		)
	}

	// Log the error with the appropriate details
	GetInstance().WithOptions(zap.AddCallerSkip(CallerSkipOne)).Error(message, logFields...)
}

func LogServiceV2(c *gin.Context, msg string, service interface{}, err error) {
	auth := utils.GetAuthData(c)

	logFields := []zap.Field{}
	if service != nil {
		typ := reflect.TypeOf(service)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		logFields = append(logFields, zap.String("service", typ.Name()))
	}

	reqUuidInterface, ok := c.Get(string(consts.RequestUuid))
	if ok {
		logFields = append(logFields,
			zap.String(string(consts.RequestUuid), reqUuidInterface.(string)),
		)
	}

	if auth.IsAuthenticated {
		logFields = append(logFields,
			zap.String(string(consts.OwnerType), auth.OwnerType),
			zap.Uint(string(consts.OwnerId), auth.OwnerId),
			zap.String("access-token-uuid", auth.TokenUUID),
		)
	}

	if err != nil {
		logFields = append(logFields, zap.Error(err))
	}

	GetInstance().WithOptions(zap.AddCallerSkip(CallerSkipOne)).Error(msg, logFields...)
}

func LogIncomingRequest(c *gin.Context) {
	host, _ := os.Hostname()

	GetInstance().WithOptions(zap.AddCallerSkip(CallerSkipOne)).Info("Incoming request",
		zap.String("type", "request"),
		zap.String("method", c.Request.Method),
		zap.String("url", c.Request.URL.String()),
		zap.String("client_ip", c.ClientIP()),
		zap.String("user_agent", c.Request.UserAgent()),
		zap.String("OS", host),
		zap.String("request-uuid", c.GetString(string(consts.RequestUuid))),
	)
}

func LogOwnerSession(token *models.AccessTokenModel) {
	GetInstance().WithOptions(zap.AddCallerSkip(CallerSkipOne)).Info("User authenticated",
		zap.Uint(string(consts.OwnerId), token.OwnerID),
		zap.String(string(consts.OwnerType), token.OwnerType),
		zap.String("access-token-uuid", token.Uuid.String()),
		zap.String("session-status", "active"),
	)
}

func LogErrorWithFields(msg string, service interface{}, ownerType string, ownerId uint, err error, extraFields ...zap.Field) {
	fields := make([]zap.Field, 0)
	if service != nil {
		typ := reflect.TypeOf(service)
		// if it's a pointer, get element
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		fields = append(fields, zap.String("service", typ.Name()))
	}

	if ownerType != "" {
		fields = append(fields, zap.String(string(consts.OwnerType), ownerType))
	}

	if ownerId > 0 {
		fields = append(fields, zap.Uint(string(consts.OwnerId), ownerId))
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
	}

	fields = append(fields, extraFields...)

	GetInstance().WithOptions(zap.AddCallerSkip(CallerSkipOne)).Error(msg, fields...)
}

func LogErrorWithFieldsV2(ctx context.Context, msg string, service interface{}, err error, extraFields ...zap.Field) {
	fields := make([]zap.Field, 0)
	if service != nil {
		typ := reflect.TypeOf(service)
		// if it's a pointer, get element
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		fields = append(fields, zap.String("service", typ.Name()))
	}

	var reqUuid string
	var ownerType string
	var ownerId uint

	reqUuidInterface := ctx.Value(consts.RequestUuid)
	if reqUuidInterface == nil {
		reqUuid = ""
	} else {
		reqUuid = reqUuidInterface.(string)
	}

	ownerTypeInterface := ctx.Value(consts.OwnerType)
	if ownerTypeInterface == nil {
		ownerType = ""
	} else {
		ownerType = ownerTypeInterface.(string)
	}

	ownerIdInterface := ctx.Value(consts.OwnerId)
	if ownerIdInterface == nil {
		ownerId = uint(0)
	} else {
		ownerId = ownerIdInterface.(uint)
	}

	if reqUuid != "" {
		fields = append(fields, zap.String(string(consts.RequestUuid), reqUuid))
	}

	if ownerType != "" {
		fields = append(fields, zap.String(string(consts.OwnerType), ownerType))
	}

	if ownerId > 0 {
		fields = append(fields, zap.Uint(string(consts.OwnerId), ownerId))
	}

	if ownerId == 0 || ownerType == "" { // values have not been set, search for national-identity-code
		var nid string
		nidInterface := ctx.Value(consts.NationalIdentityCode)
		if nidInterface == nil {
			nid = ""
		} else {
			nid = nidInterface.(string)
		}

		if nid != "" {
			fields = append(fields, zap.String(string(consts.NationalIdentityCode), nid))
		}
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
	}

	fields = append(fields, extraFields...)

	GetInstance().WithOptions(zap.AddCallerSkip(CallerSkipOne)).Error(msg, fields...)
}

func LogInfo(ctx context.Context, msg string, service interface{}, extraFields ...zap.Field) {
	fields := make([]zap.Field, 0)
	if service != nil {
		typ := reflect.TypeOf(service)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		fields = append(fields, zap.String("service", typ.Name()))
	}

	var reqUuid string
	var ownerType string
	var ownerId uint

	reqUuidInterface := ctx.Value(consts.RequestUuid)
	if reqUuidInterface == nil {
		reqUuid = ""
	} else {
		reqUuid = reqUuidInterface.(string)
	}

	ownerTypeInterface := ctx.Value(consts.OwnerType)
	if ownerTypeInterface == nil {
		ownerType = ""
	} else {
		ownerType = ownerTypeInterface.(string)
	}

	ownerIdInterface := ctx.Value(consts.OwnerId)
	if ownerIdInterface == nil {
		ownerId = uint(0)
	} else {
		ownerId = ownerIdInterface.(uint)
	}

	if reqUuid != "" {
		fields = append(fields, zap.String(string(consts.RequestUuid), reqUuid))
	}

	if ownerType != "" {
		fields = append(fields, zap.String(string(consts.OwnerType), ownerType))
	}

	if ownerId > 0 {
		fields = append(fields, zap.Uint(string(consts.OwnerId), ownerId))
	}

	if ownerId == 0 || ownerType == "" { // values have not been set, search for national-identity-code
		var nid string
		nidInterface := ctx.Value(consts.NationalIdentityCode)
		if nidInterface == nil {
			nid = ""
		} else {
			nid = nidInterface.(string)
		}

		if nid != "" {
			fields = append(fields, zap.String(string(consts.NationalIdentityCode), nid))
		}
	}

	fields = append(fields, extraFields...)
	GetInstance().WithOptions(zap.AddCallerSkip(CallerSkipOne)).Info(msg, fields...)
}

// LogErrorDirect logs an error message with optional fields without context. Automatically detects non-critical errors and logs them as warnings.
func LogErrorDirect(msg string, err error, extraFields ...zap.Field) {
	fields := make([]zap.Field, 0)
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	fields = append(fields, extraFields...)

	loggerOfLevel, _ := getLoggerWithLevel(err, CallerSkipOne)
	loggerOfLevel(msg, fields...)
}
