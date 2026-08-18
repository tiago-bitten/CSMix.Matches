// Package httpin holds what every inbound request passes through: the
// middleware chain and the two ways this service answers.
package httpin

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// Problem is RFC 9457, the same shape CSMix.Accounts answers with, so a client
// reads one error format across the platform.
type Problem struct {
	Title  string `json:"title"`
	Status int    `json:"status"`
	Code   string `json:"code"`
	Detail string `json:"detail,omitempty"`
}

func JSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if body == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Error("writing the response body failed", "error", err)
	}
}

func Fail(w http.ResponseWriter, status int, code string) {
	JSON(w, status, Problem{Title: code, Status: status, Code: code})
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
