package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cruisechang/dbServer/handler"
	"github.com/cruisechang/dbex"
)

type middleware struct {
	lg         *dbex.Logger
	requestURI string
}

func NewMiddleware(lg *dbex.Logger) *middleware {
	return &middleware{
		lg: lg,
	}
}

func (m *middleware) writeError(w http.ResponseWriter, errorCode int) {

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	s := struct {
		Code  int    `json:"code"`
		Count int    `json:"count"`
		Data  string `json:"data"`
	}{
		Code:  errorCode,
		Count: 0,
		Data:  "[{}]",
	}

	b, _ := json.Marshal(s)
	io.WriteString(w, string(b))
}

func (m *middleware) LogRequestURI(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.lg.LogFile(dbex.LevelInfo, fmt.Sprintf("dbServer requestURI=%s", r.RequestURI))
		next.ServeHTTP(w, r)
	})
}

func (m *middleware) CheckHead(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("API-Key") != handler.HeadAPIKey {
			m.lg.Log(dbex.LevelError, fmt.Sprintf("requestURI=%s, checkHead error %s", r.RequestURI, r.Header.Get("API-Key")))
			m.writeError(w, handler.CodeHeaderError)
			return
		}
		next.ServeHTTP(w, r)
	})
}
