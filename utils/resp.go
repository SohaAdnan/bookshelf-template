package utils

import (
	"net/http"

	"github.com/go-chi/render"
)

// RespBody represents rest api response
type RespBody struct {
	OK         bool        `json:"ok"`
	StatusCode int         `json:"-"`
	Data       interface{} `json:"data,omitempty"`
	Message    string      `json:"message,omitempty"`
}

// Render is used for implementing `render.Renderer` interface
func (res *RespBody) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, res.StatusCode)

	return nil
}

// NewSuccessResp is used for creating success response
func NewSuccessResp(data interface{}) *RespBody {
	return &RespBody{
		OK:         true,
		StatusCode: http.StatusOK,
		Data:       data,
	}
}

func NewErrorResp(statusCode int, message string) *RespBody {
	return &RespBody{
		OK:         false,
		StatusCode: statusCode,
		Message:    message,
	}
}
