package clients

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
	"time"
)

const (
	StatusOKDefaultMessage    = "success"
	ErrorDefaultMessage       = "generic message"
	ServerErrorDefaultMessage = "server error"
	BadRequestDefaultMessage  = "bad request"
)

type Option func(executableError *ApiResponse)

type ApiResponse struct {
	Code           int               `json:"code,omitempty" yaml:"code,omitempty"`
	Msg            string            `json:"msg,omitempty" yaml:"msg,omitempty"`
	AdditionalInfo map[string]string `json:"additional-info,omitempty" yaml:"additional-info,omitempty"`
	Ts             string            `yaml:"timestamp,omitempty" mapstructure:"timestamp,omitempty" json:"timestamp,omitempty"`
}

func (ar *ApiResponse) Error() string {
	var sv strings.Builder
	const sep = " - "

	if ar.Code > 0 {
		sv.WriteString(fmt.Sprintf("code: %d"+sep, ar.Code))
	}

	if ar.Msg != "" {
		sv.WriteString(fmt.Sprintf("message: %s"+sep, ar.Msg))
	}

	if ar.Ts != "" {
		sv.WriteString(fmt.Sprintf("timestamp: %s"+sep, ar.Ts))
	}

	return strings.TrimSuffix(sv.String(), sep)
}

func DeserializeApiResponseFromJson(b []byte) (ApiResponse, error) {
	a := ApiResponse{}
	err := json.Unmarshal(b, &a)
	return a, err
}

func WithCode(c int) Option {
	return func(e *ApiResponse) {
		e.Code = c
	}
}

func WithMessage(m string) Option {
	return func(e *ApiResponse) {
		e.Msg = m
	}
}

func NewInternalServerError(opts ...Option) *ApiResponse {
	err := &ApiResponse{Code: http.StatusInternalServerError, Msg: ErrorDefaultMessage}
	for _, o := range opts {
		o(err)
	}
	return err
}

func NewSuccessResponse(opts ...Option) *ApiResponse {
	err := &ApiResponse{Code: http.StatusOK, Msg: StatusOKDefaultMessage}
	for _, o := range opts {
		o(err)
	}
	return err
}

func NewBadRequestError(opts ...Option) *ApiResponse {
	err := &ApiResponse{Code: http.StatusBadRequest, Msg: BadRequestDefaultMessage}
	for _, o := range opts {
		o(err)
	}
	return err
}

func (ar *ApiResponse) ToJSON() []byte {
	ar.Ts = time.Now().Format(time.RFC3339)
	var b []byte
	var err error

	b, err = json.Marshal(ar)
	if err != nil {
		log.Error().Err(err).Msg("error in marshalling api-error")
		return []byte(`{"msg": "error in marshalling api-error"}`)
	}

	return b
}

func Code(err error) int {
	var resp *ApiResponse
	if errors.As(err, &resp) {
		return resp.Code
	}

	return 0
}
