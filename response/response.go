package response

import (
	"time"
)

type JSONResponse struct {
	Data        interface{} `json:"data,omitempty"`
	Message     string      `json:"message,omitempty"`
	Code        string      `json:"code,omitempty"`
	StatusCode  int         `json:"status_code,omitempty"`
	ErrorString string      `json:"error,omitempty"`
	Latency     int64       `json:"latency,omitempty"`
	Result      interface{} `json:"result,omitempty"`
}

func NewJSONResponse() JSONResponse {
	return JSONResponse{}
}

func (r JSONResponse) WithData(data interface{}) JSONResponse {
	r.Data = data
	return r
}

func (r JSONResponse) WithMessage(message string) JSONResponse {
	r.Message = message
	return r
}

func (r JSONResponse) WithCode(code string) JSONResponse {
	r.Code = code
	return r
}

func (r JSONResponse) WithStatusCode(statusCode int) JSONResponse {
	r.StatusCode = statusCode
	return r
}

func (r JSONResponse) WithErrorString(errorStr string) JSONResponse {
	r.ErrorString = errorStr
	return r
}

func (r JSONResponse) WithLatency(start time.Time) JSONResponse {

	end := time.Now()
	r.Latency = end.Sub(start).Microseconds()
	return r
}

func (r JSONResponse) WithResult(result interface{}) JSONResponse {
	r.Result = result
	return r
}
