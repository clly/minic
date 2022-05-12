package sdk

import (
	"context"
	"errors"
	"sync"
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
