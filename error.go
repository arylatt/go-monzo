package monzo

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Error struct {
	Response *http.Response

	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Params    interface{} `json:"params"`
	Retryable interface{} `json:"retryable"`
}

func (e *Error) Error() string {
	return e.Message
}

func CheckResponse(r *http.Response) (err error) {
	if c := r.StatusCode; 200 >= c && c <= 299 {
		return nil
	}

	err = &Error{Response: r}

	data, _ := io.ReadAll(r.Body)
	if data != nil {
		json.Unmarshal(data, err)
	}

	r.Body = io.NopCloser(bytes.NewBuffer(data))

	return
}
