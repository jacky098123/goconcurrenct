package goconcurrenct

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_goWithRecovery(t *testing.T) {
	zero := 0
	scenarios := []struct {
		name            string
		callFun         func() error
		timeoutDuration time.Duration
		expErr          error
	}{
		{
			name: "happy path",
			callFun: func() error {
				time.Sleep(defaultDuration / 2)
				return nil
			},
			timeoutDuration: defaultDuration * 100,
			expErr:          nil,
		},
		{
			name: "panic path",
			callFun: func() error {
				time.Sleep(defaultDuration / 2)
				_ = 10 / zero // raise panic
				return errors.New("return panic error")
			},
			timeoutDuration: defaultDuration * 100,
			expErr:          ErrConcurrentPanic,
		},
		{
			name: "timeout path",
			callFun: func() error {
				time.Sleep(defaultDuration * 2)
				_ = 10 / zero // raise panic
				return errors.New("return panic error")
			},
			timeoutDuration: defaultDuration * 100,
			expErr:          ErrConcurrentTimeout,
		},
		{
			name: "context timeout path",
			callFun: func() error {
				time.Sleep(defaultDuration * 2)
				_ = 10 / zero // raise panic
				return nil
			},
			timeoutDuration: defaultDuration / 2,
			expErr:          ErrConcurrentContext,
		},
	}

	for _, p := range scenarios {
		scenario := p
		t.Run(scenario.name, func(t *testing.T) {
			ctx, _ := context.WithTimeout(context.Background(), scenario.timeoutDuration)
			_, err := goWithRecovery(ctx, scenario.name, scenario.callFun)
			assert.Equal(t, scenario.expErr, err)
		})
	}
}
