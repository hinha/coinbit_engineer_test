package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Error struct {
	w    http.ResponseWriter
	orig error
	code int
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v", e.orig)
}

func (e *Error) ResponseJson() {
	response, _ := json.Marshal(map[string]interface{}{"error": e.Error(), "code": e.code})

	e.w.Header().Set("Content-Type", "application/json")
	e.w.WriteHeader(e.code)
	e.w.Write(response)
}

func NewErrorf(w http.ResponseWriter, orig error, code int) *Error {
	return &Error{
		w:    w,
		code: code,
		orig: orig,
	}
}

func ResponseJson(w http.ResponseWriter, msg, data interface{}, code int) {
	response, _ := json.Marshal(map[string]interface{}{"code": code, "message": msg, "data": data})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
