package context

import (
	"encoding/json"
	"github.com/flint92/webapp/resp"
	"io"
	"net/http"
)

type Context struct {
	W          http.ResponseWriter
	R          *http.Request
	PathParams map[string]string
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		W:          w,
		R:          r,
		PathParams: make(map[string]string),
	}
}

func (c *Context) ReadJson(obj interface{}) error {
	body, err := io.ReadAll(c.R.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, obj)
}

func (c *Context) OkJson(data interface{}) error {
	return c.writeJson(http.StatusOK, resp.Ok(data))
}

func (c *Context) ErrorJson(errorCode resp.ErrorCode, errorMessage string) error {
	var httpStatus int
	switch errorCode {
	case resp.ErrorCodeForbidden:
		httpStatus = http.StatusForbidden
	case resp.ErrorCodeUnauthorized:
		httpStatus = http.StatusUnauthorized
	case resp.ErrorCodeBadRequest:
		httpStatus = http.StatusBadRequest
	case resp.ErrorCodeNotFound:
		httpStatus = http.StatusNotFound
	default:
		httpStatus = http.StatusInternalServerError
	}

	return c.writeJson(httpStatus, resp.Fail(errorCode, errorMessage))
}

func (c *Context) ServerErrorJson() error {
	return c.writeJson(http.StatusInternalServerError, resp.InternalError())
}

func (c *Context) BadRequestJson() error {
	return c.writeJson(http.StatusBadRequest, resp.BadRequest())
}

func (c *Context) UnauthorizedJson() error {
	return c.writeJson(http.StatusUnauthorized, resp.Unauthorized())
}

func (c *Context) ForbiddenJson() error {
	return c.writeJson(http.StatusForbidden, resp.Forbidden())
}

func (c *Context) NotFoundJson() error {
	return c.writeJson(http.StatusNotFound, resp.NotFound())
}

func (c *Context) writeJson(code int, resp resp.CommonResponse) error {
	respJson, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	c.W.Header().Set("Content-Type", "application/json")
	c.W.WriteHeader(code)
	_, err = c.W.Write(respJson)
	if err != nil {
		return err
	}

	return nil
}
