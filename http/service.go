package http

import (
	"encoding/json"
	"net/http"

	"github.com/shanebailey05/future_backend_homework/app"
)

type Service struct {
	app *app.Service
}

func New() (*Service, error) {
	ret := new(Service)
	var err error
	ret.app, err = app.New()
	if err != nil {
		return nil, err
	}
	return ret, nil
}

type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func (s *Service) errorResponse(code int, errStr string, w http.ResponseWriter, r *http.Request) {
	errResp := &ErrorResponse{
		Code:  code,
		Error: errStr,
	}

	b, err := json.Marshal(&errResp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Could not marshal json error response"))
		return
	}

	w.WriteHeader(code)
	_, _ = w.Write(b)
}
