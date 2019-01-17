package middleware

import (
	"net/http"
	"fmt"
	"github.com/cruisechang/dbServer/handler"
	"encoding/json"
	"io"
	"github.com/cruisechang/dbex"
)

type middleware struct {
	lg         *dbex.Logger
	requestURI string
}

func NewMiddleware(lg *dbex.Logger) (*middleware) {
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
		m.lg.LogFile(dbex.LevelInfo, fmt.Sprintf("requestURI:%s", r.RequestURI))
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

//func (m *middleware) CheckBody(next http.Handler) http.Handler {
//
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//
//		body, err := r.GetBody()
//		if err != nil {
//			m.lg.Log(dbexLog.LevelError, fmt.Sprintf("requestURI=%s, body error=%s\n", r.RequestURI, err.Error()))
//			m.writeError(w, handler.C)
//			return
//		}
//
//		if body == nil {
//			m.lg.Log(dbexLog.LevelError, fmt.Sprintf("requestURI=%s, body =nil\n", r.RequestURI))
//			m.writeError(w, config.CodeError0)
//			return
//		}
//
//		b, err := ioutil.ReadAll(body)
//
//		if err != nil {
//			m.lg.Log(dbexLog.LevelError, fmt.Sprintf("requestURI=%s, parsing body error=%s\n", r.RequestURI, err.Error()))
//			m.writeError(w, config.CodeError0)
//			return
//		}
//
//		if b == nil {
//			m.lg.Log(dbexLog.LevelError, fmt.Sprintf("requestURI=%s, parsing body ==nil\n", r.RequestURI))
//			m.writeError(w, config.CodeError0)
//			return
//		}
//		next.ServeHTTP(w, r)
//	})
//}
