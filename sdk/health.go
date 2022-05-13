package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

func AlwaysOK(context.Context) error { return nil }

func AlwaysFailing(context.Context) error { return errors.New("oopsie") }

type Healthcheck struct {
	checks     map[string]Healthchecker
	m          *sync.RWMutex
	initalized sync.Once
}

type result struct {
	Name    string `json:"name"`
	Result  int    `json:"result"`
	Message string `json:"message"`
}

func (h *Healthcheck) AddHealthcheck(name string, check Healthchecker) error {
	h.initialize()
	h.m.Lock()
	h.checks[name] = check
	h.m.Unlock()
	return nil
}

func (h *Healthcheck) initialize() {
	h.initalized.Do(func() {
		h.m = &sync.RWMutex{}
		h.checks = make(map[string]Healthchecker)
	})
}

type Healthchecker interface {
	Health(ctx context.Context) error
}

var _ Healthchecker = (*HealthcheckerFunc)(nil)

type HealthcheckerFunc func(context.Context) error

func (f HealthcheckerFunc) Health(ctx context.Context) error {
	return f(ctx)
}

func runCheck(ctx context.Context, name string, check Healthchecker) result {
	err := check.Health(ctx)
	var msg string
	var state int
	if err != nil {
		msg = err.Error()
		state = 1
	}
	return result{
		Name:    name,
		Result:  state,
		Message: msg,
	}
}

// This should be replaced with something on the healthcheck struct
var healthcheckTimeout = 30 * time.Second

func (h *Healthcheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	ctx, cancel := context.WithTimeout(ctx, healthcheckTimeout)
	defer cancel()

	out := make(chan []result)
	runChecks := func(ctx context.Context) {
		results := h.runChecks(ctx)
		out <- results
	}

	var results []result
	go runChecks(ctx)
	select {
	case <-ctx.Done():
		results = <-out
		w.WriteHeader(http.StatusGatewayTimeout)
	case results = <-out:
	}
	for _, result := range results {
		if result.Result != 0 {
			w.WriteHeader(512)
		}
	}
	err := encoder.Encode(results)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to write healthcheck results", err.Error())
	}
}

func (h *Healthcheck) runChecks(ctx context.Context) []result {
	results := make([]result, 0, len(h.checks))
	h.m.RLock()
	for name, check := range h.checks {
		if ctx.Err() != nil {
			results = append(results, result{
				Name:    name,
				Result:  -1,
				Message: "healthcheck execution timed out",
			})
		}
		result := runCheck(ctx, name, check)
		results = append(results, result)
	}
	return results
}
