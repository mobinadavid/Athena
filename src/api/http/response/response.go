package response

import (
	"athena/src/pkg/i18n"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	IsSuccessful bool              `json:"is_successful"`
	RequestUuid  string            `json:"request_uuid"`
	RequestIp    string            `json:"request_ip"`
	StatusCode   int               `json:"status_code"`
	Message      string            `json:"message"`
	Data         map[string]any    `json:"data,omitempty"`
	Errors       map[string]string `json:"errors,omitempty"`
	ErrorCode    int               `json:"error_code,omitempty"`
	Log          bool              `json:"-"`
}

// Builder is the builder for constructing API responses
type Builder struct {
	c        *gin.Context
	response Response
}

// Api initializes a new response builder with default values
func Api(c *gin.Context) *Builder {
	return &Builder{
		c: c,
		response: Response{
			RequestUuid: c.GetString("request-uuid"),
			RequestIp:   c.GetHeader("x-forwarded-for"),
			// Default values.
			IsSuccessful: false,
			StatusCode:   http.StatusBadRequest,
			Message:      i18n.Localize(c.GetString("locale"), "invalid-payload"),
		},
	}
}

// SetStatusCode sets the status code of the response
func (builder *Builder) SetStatusCode(statusCode int) *Builder {
	builder.response.StatusCode = statusCode
	return builder
}

// SetErrorCode sets the status code of the response
func (builder *Builder) SetErrorCode(errorCode int) *Builder {
	builder.response.ErrorCode = errorCode
	return builder
}

// SetMessage sets the message of the response
func (builder *Builder) SetMessage(message string) *Builder {
	builder.response.Message = message
	return builder
}

// SetData sets the data of the response
func (builder *Builder) SetData(data map[string]any) *Builder {
	builder.response.Data = data
	return builder
}

// SetErrors sets the errors of the response
func (builder *Builder) SetErrors(errors map[string]string) *Builder {
	builder.response.Errors = errors
	return builder
}

// SetLog log the response
func (builder *Builder) SetLog() *Builder {
	builder.response.Log = true
	return builder
}

// Send sends the constructed response to the client
func (builder *Builder) Send() {
	builder.response.IsSuccessful = builder.response.StatusCode >= 200 && builder.response.StatusCode < 300
	builder.c.JSON(builder.response.StatusCode, builder.response)
	builder.c.Abort()
}
