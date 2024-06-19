package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type respData struct {
	size   int
	status int
}

func newRespData() respData {
	return respData{
		size:   0,
		status: 0,
	}
}

type rwWrapper struct {
	http.ResponseWriter
	respData
}

func (rw *rwWrapper) Write(p []byte) (int, error) {
	s, err := rw.ResponseWriter.Write(p)
	rw.size += s
	return s, err
}

func (rw *rwWrapper) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func WithLogging(logger *logrus.Logger) func(next http.Handler) http.Handler {
	mw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			httpMethod := r.Method
			remoteAddr := r.RemoteAddr
			uri := r.RequestURI
			ww := &rwWrapper{
				ResponseWriter: w,
				respData:       newRespData(),
			}
			next.ServeHTTP(ww, r)
			logger.Info(fmt.Sprintf("[%s]  uri=%s  remote_addr=%s status=%d data_send=%d duration=%s",
				httpMethod,
				uri,
				remoteAddr,
				ww.status,
				ww.size,
				time.Since(start)))
		})
	}

	return mw
}
