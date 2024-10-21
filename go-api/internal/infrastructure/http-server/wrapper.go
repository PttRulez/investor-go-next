package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/handlers"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/metrics"
	"github.com/pttrulez/investor-go-next/go-api/pkg/logger"
)

func (wr *Wrapper) makeHttpHandler(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			latency := time.Since(start).Seconds()
			wr.mtrcs.WriteLatency(r.Method, r.URL.Path, latency)
		}(time.Now())

		if err := handler(w, r); err != nil {
			var apiErr handlers.APIError
			// в метрики и логи летят 500-е
			if errors.As(err, &apiErr) && apiErr.Code == http.StatusInternalServerError {
				// увеличиваем счётчик 500-х ошибок в метриках, если это нага ApiErr
				wr.mtrcs.IncTotalInternalErrors()

				// логгируем в графану или просто в stdout
				wr.log.Error(err)
			}
		}

		// увеличиваем счётчик запросов
		wr.mtrcs.IncTotalRequests()
	}
}

type Wrapper struct {
	log   *logger.Logger
	mtrcs *metrics.Metrics
}

func NewWrapper(log *logger.Logger, mtrcs *metrics.Metrics) *Wrapper {
	return &Wrapper{
		log:   log,
		mtrcs: mtrcs,
	}
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error
