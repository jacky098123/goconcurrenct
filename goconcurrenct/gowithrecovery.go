package goconcurrenct

import (
	"context"
	"errors"
	"sync"
	"time"
)

const (
	defaultName     = "defaultName"
	defaultDuration = 2 * time.Second
	logTag          = "GoWithRecovery"
)

var (
	// ErrConcurrentPanic ...
	ErrConcurrentPanic = errors.New("panic")

	// ErrConcurrentTimeout ...
	ErrConcurrentTimeout = errors.New("timeout")

	// ErrConcurrentContext ...
	ErrConcurrentContext = errors.New("context error")
)

type options struct {
	name            string
	timeoutDuration time.Duration
}

// Option ...
type Option func(*options)

// WithName define the log, stats name for the goroutine
func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

// WithTimeoutDuration defines if need send the goroutine result to stats
func WithTimeoutDuration(duration time.Duration) Option {
	return func(o *options) {
		o.timeoutDuration = duration
	}
}

// GoWithRecovery is a var, design for mock
var GoWithRecovery = goWithRecovery

type command struct {
	err     error
	errOnce *sync.Once
	tags    []string
}

func newCommand() *command {
	c := command{
		err:     nil,
		errOnce: &sync.Once{},
	}

	return &c
}

func (c *command) handleErrOnce(err error) {
	c.errOnce.Do(func() {
		c.err = err
	})
}

// goWithRecovery runs a command pattern, it is sync call, it is a subset of Hystrix.Do
func goWithRecovery(ctx context.Context, trackID string, callFun func() error, ops ...Option) ([]string, error) {
	opt := &options{
		name:            defaultName,
		timeoutDuration: defaultDuration,
	}

	for _, op := range ops {
		op(opt)
	}

	runErrChan := make(chan error, 1)
	cmd := newCommand()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				runErrChan <- ErrConcurrentPanic
			}
			close(runErrChan)
		}()

		runErrChan <- callFun()
	}()

	timer := time.NewTimer(opt.timeoutDuration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		cmd.handleErrOnce(ErrConcurrentContext)
		cmd.tags = append(cmd.tags, "reason:ctxdone", "err:"+ctx.Err().Error())
	case runErr := <-runErrChan:
		if runErr == nil {
			cmd.handleErrOnce(nil)
		} else if runErr == ErrConcurrentPanic {
			cmd.handleErrOnce(runErr)
			cmd.tags = append(cmd.tags, "reason:panic")
		} else {
			cmd.handleErrOnce(runErr)
			cmd.tags = append(cmd.tags, "reason:callFunc")
		}
	case <-timer.C:
		cmd.handleErrOnce(ErrConcurrentTimeout)
		cmd.tags = append(cmd.tags, "reason:timeout")
	}

	return cmd.tags, cmd.err
}
