package resp

type ErrorCode int

const (
	ErrorCodeSuccess ErrorCode = iota

	ErrorCodeNotFound ErrorCode = iota + 1000
	ErrorCodeBadRequest
	ErrorCodeUnauthorized
	ErrorCodeForbidden
	ErrorCodeInternalError
)

type CommonResponse struct {
	Message   string      `json:"message"`
	ErrorCode ErrorCode   `json:"error_code"`
	Data      interface{} `json:"data,omitempty"`
}

func Ok(data interface{}) CommonResponse {
	return CommonResponse{
		Message:   "success",
		ErrorCode: ErrorCodeSuccess,
		Data:      data,
	}
}

func NotFound() CommonResponse {
	return Fail(ErrorCodeNotFound, "not found")
}

func BadRequest() CommonResponse {
	return Fail(ErrorCodeBadRequest, "bad request")
}

func Unauthorized() CommonResponse {
	return Fail(ErrorCodeUnauthorized, "unauthorized")
}

func Forbidden() CommonResponse {
	return Fail(ErrorCodeForbidden, "forbidden")
}

func InternalError() CommonResponse {
	return Fail(ErrorCodeInternalError, "internal error")
}

func Fail(errorCode ErrorCode, message string) CommonResponse {
	return CommonResponse{
		Message:   message,
		ErrorCode: errorCode,
	}
}
