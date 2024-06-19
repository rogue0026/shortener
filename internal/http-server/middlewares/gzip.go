package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
)

type wWrapper struct {
	http.ResponseWriter
	compressor io.Writer
}

func (ww *wWrapper) Write(p []byte) (int, error) {
	compressedSize, err := ww.compressor.Write(p)
	return compressedSize, err
}

func (ww *wWrapper) WriteHeader(statusCode int) {
	if statusCode < 300 && statusCode >= 200 {
		ww.ResponseWriter.Header().Set("Content-Encoding", "gzip")
	}
	ww.ResponseWriter.WriteHeader(statusCode)
}

func WithCompressing(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		} else {
			encoding := r.Header.Get("Accept-Encoding")
			switch encoding {
			case "":
				next.ServeHTTP(w, r)
			case "gzip":
				c := gzip.NewWriter(w)
				defer c.Close()
				wr := &wWrapper{
					ResponseWriter: w,
					compressor:     c,
				}
				next.ServeHTTP(wr, r)
			}
		}
	}
}
