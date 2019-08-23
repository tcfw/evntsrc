package apigw

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func loggingWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl := &responseLogger{rw: w, start: time.Now()}
		h.ServeHTTP(rl, r)

		zap.S().Infow(r.Method,
			"url", r.URL.String(),
			"size", rl.size,
			"status", rl.status,
			"took", time.Since(rl.start).String(),
			"ua", r.UserAgent(),
			"ref", r.Referer(),
		)
	})
}

type responseLogger struct {
	rw     http.ResponseWriter
	start  time.Time
	status int
	size   int
}

func (rl *responseLogger) Header() http.Header {
	return rl.rw.Header()
}

func (rl *responseLogger) Write(bytes []byte) (int, error) {
	if rl.status == 0 {
		rl.status = http.StatusOK
	}

	size, err := rl.rw.Write(bytes)

	rl.size += size

	return size, err
}

func (rl *responseLogger) WriteHeader(status int) {
	rl.status = status

	rl.rw.WriteHeader(status)
}

func (rl *responseLogger) Flush() {
	f, ok := rl.rw.(http.Flusher)

	if ok {
		f.Flush()
	}
}
