package http

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/failsafe-go/failsafe-go"
	"github.com/failsafe-go/failsafe-go/circuitbreaker"
	"github.com/failsafe-go/failsafe-go/retrypolicy"
	"github.com/moises-ba/mb-crypto-mms-api/internal/errors"
	"github.com/moises-ba/mb-crypto-mms-api/internal/log"
)

type Response struct {
	Status   int
	Response []byte
}

type Client interface {
	Get(ctx context.Context, url string) (*Response, errors.ApiError)
	//TODO PUT, PATCH, DELETE, POST...
}

type client struct {
	httpClient *http.Client
	executor   failsafe.Executor[*http.Response]
}

func NewClient() *client {
	handlelf := func(resp *http.Response, err error) bool {
		return err != nil || resp.StatusCode >= 500 || resp.StatusCode == 429
	}

	cb := circuitbreaker.NewBuilder[*http.Response]().
		HandleIf(handlelf).
		WithFailureThreshold(5).
		WithDelay(30 * time.Second).
		WithSuccessThreshold(2).
		Build()

	retry := retrypolicy.NewBuilder[*http.Response]().
		HandleIf(handlelf).
		WithBackoff(500*time.Millisecond, 2*time.Second).
		WithMaxRetries(2).
		Build()

	return &client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		// Retry é a política interna; o circuit breaker a envolve.
		executor: failsafe.With(cb, retry),
	}
}

func (c *client) Get(ctx context.Context, url string) (*Response, errors.ApiError) {
	resp, err := c.executor.WithContext(ctx).GetWithExecution(
		func(exec failsafe.Execution[*http.Response]) (*http.Response, error) {
			req, err := http.NewRequestWithContext(exec.Context(), http.MethodGet, url, nil)
			if err != nil {
				log.Error("request error: "+req.URL.String(), err)
				return nil, errors.NewApiError("error exec calling get", errors.WithError(err), errors.WithKind(errors.Unexpected))
			}

			return c.httpClient.Do(req)
		},
	)

	if err != nil {
		log.Error("request error", err)
		return nil, errors.NewApiError("error calling get", errors.WithError(err), errors.WithKind(errors.Unexpected))
	}
	defer resp.Body.Close()

	if apiErr := evaluateError(resp.StatusCode); apiErr != nil {
		return nil, apiErr
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewApiError("error calon read body", errors.WithError(err), errors.WithKind(errors.Unexpected))
	}

	return &Response{Status: resp.StatusCode, Response: body}, nil
}

func evaluateError(httStatusCode int) errors.ApiError {
	switch httStatusCode {
	case http.StatusNotFound:
		return errors.NewApiError("not found", errors.WithKind(errors.NotFound))
	case http.StatusBadRequest:
		return errors.NewApiError("invalid request", errors.WithKind(errors.Invalid))
	//TODO mapear outros erros...
	default:
		return nil
	}
}
