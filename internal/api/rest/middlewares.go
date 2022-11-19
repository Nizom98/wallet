package rest

import (
	"bytes"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// MiddlewareLog логируем метод, путь, время выполнения запроса.
func (h *Handler) MiddlewareLog(next func(w http.ResponseWriter, req *http.Request)) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var body []byte
		startTime := time.Now()

		defer func() {
			entry := log.WithFields(log.Fields{
				"URI":       req.RequestURI,
				"METHOD":    req.Method,
				"WORK_TIME": time.Since(startTime),
				"BODY":      string(body),
			})
			entry.Debugf("request logging")
		}()

		body, err := io.ReadAll(req.Body)
		if err != nil {
			printError(w, err.Error(), http.StatusBadRequest)
		}

		req.Body = io.NopCloser(bytes.NewBuffer(body))
		next(w, req)
	}

}
