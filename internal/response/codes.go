package response

type Code int

const (
	OKCode           Code = 200 // http.StatusOK
	PendingCode      Code = 201 // http.StatusCreated
	BadRequestCode   Code = 400 // http.StatusBadRequest
	UnauthorizedCode Code = 401 // http.StatusUnauthorized
	NotFoundCode     Code = 404 // http.StatusNotFound
	// AlreadyExistsCode           Code = 409 // http.StatusConflict
	// WrongVerificationCode       Code = 410 // http.StatusGone
	// WrongPasswordCode           Code = 411 // http.StatusLengthRequired
	// InvalidTokenCode            Code = 412 // http.StatusPreconditionFailed
	// UnsupportedAvatarFormatCode Code = 415 // http.StatusUnsupportedMediaType
	InternalServerErrorCode Code = 500 // http.StatusInternalServerError
	BadGatewayCode          Code = 504 // http.StatusBadGateway
)

type Response interface {
	WithMessage(message string) Response
	WithPayload(payload any) Response
	GetCode() int
}

type response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Payload any    `json:"payload,omitempty"`
}

func Make(code Code) Response {
	return &response{
		Code:    int(code),
		Message: statusTest(code),
	}
}

func (r *response) WithMessage(message string) Response {
	r.Message = message
	return r
}

func (r *response) WithPayload(payload any) Response {
	r.Payload = payload
	return r
}

func (r *response) GetCode() int {
	return r.Code
}

func statusTest(code Code) string {
	switch code {
	case OKCode:
		return "Success"
	case PendingCode:
		return "In Process"
	case BadRequestCode:
		return "Bad Request"
	case UnauthorizedCode:
		return "Unauthorized"
	case NotFoundCode:
		return "Not Found"
	// case AlreadyExistsCode:
	// 	return "Already Exists"
	// case WrongVerificationCode:
	// 	return "Wrong Verification"
	// case WrongPasswordCode:
	// 	return "Wrong Password"
	// case InvalidTokenCode:
	// 	return "Invalid Token"
	// case UnsupportedAvatarFormatCode:
	// 	return "Unsupported Avatar Format"
	case InternalServerErrorCode:
		return "Internal Server Error"
	case BadGatewayCode:
		return "External Service Error"
	}
	return ""
}
