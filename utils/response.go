package utils

type CustomResponseWrapper struct {
	Data    interface{} `json:"data"`
	Error   bool        `json:"error"`
	Message string      `json:"message"`
}

func SuccessResponse(data interface{}, message ...string) CustomResponseWrapper {

	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}

	return CustomResponseWrapper{
		Data:    data,
		Message: msg,
		Error:   false,
	}
}

func ErrorResponse(message string, data ...interface{}) CustomResponseWrapper {

	var d interface{}
	if len(data) > 0 {
		d = data[0]
	}

	return CustomResponseWrapper{
		Data:    d,
		Message: message,
		Error:   true,
	}
}
