package middleware

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"ritchie-server/server"
	"ritchie-server/server/metrics"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type Handler struct {
	Authorization 		server.Constraints
}

func NewMiddlewareHandler(a server.Constraints) Handler {
	return Handler{Authorization: a}
}

func (mh Handler) Filter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ok := true
		var err error
		if !mh.Authorization.ValidatePublicConstraints(r.URL.Path, r.Method) {
			authorizationToken := r.Header.Get(server.AuthorizationHeader)
			organization := r.Header.Get(server.OrganizationHeader)
			ok, err = mh.Authorization.AuthorizationPath(authorizationToken, r.URL.Path, r.Method, organization)
			if err != nil {
				log.Error("Authorization failed ", err)
				w.WriteHeader(http.StatusUnauthorized)
				_, err := w.Write([]byte(fmt.Sprintf("Authorization Failed: %v", err.Error())))
				if err != nil {
					fmt.Sprintln("Error in Write ")
					return
				}
				path:= strings.ReplaceAll(r.URL.Path, ".", "-")
				metrics.Metric(path).With(prometheus.Labels{"code": "401"}).Inc()
				return
			}
		}
		if ok {
			start := time.Now()
			ww := &responseWriter{w, http.StatusOK}
			next.ServeHTTP(ww, r)
			metrics.Metric(strings.ReplaceAll(r.URL.Path, ".", "-")).With(prometheus.Labels{"code": strconv.Itoa(ww.statusCode)}).Inc()
			metrics.LatencyOpsRequest.With(prometheus.Labels{"path": r.URL.Path}).Observe(float64(time.Since(start).Milliseconds()))
		} else {
			w.WriteHeader(http.StatusForbidden)
			_, err := w.Write([]byte("Forbidden "))
			if err != nil {
				fmt.Sprintln("Error in Write ")
				return
			}
		}
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *responseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
