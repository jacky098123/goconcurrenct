// Copyright (c) 2012-2024 Grabtaxi Holdings PTE LTD (GRAB), All Rights Reserved. NOTICE: All information contained herein
// is, and remains the property of GRAB. The intellectual and technical concepts contained herein are confidential, proprietary
// and controlled by GRAB and may be covered by patents, patents in process, and are protected by trade secret or copyright law.
//
// You are strictly forbidden to copy, download, store (in any medium), transmit, disseminate, adapt or change this material
// in any way unless prior written permission is obtained from GRAB. Access to the source code contained herein is hereby
// forbidden to anyone except current GRAB employees or contractors with binding Confidentiality and Non-disclosure agreements
// explicitly covering such access.
//
// The copyright notice above does not evidence any actual or intended publication or disclosure of this source code,
// which includes information that is confidential and/or proprietary, and is a trade secret, of GRAB.
//
// ANY REPRODUCTION, MODIFICATION, DISTRIBUTION, PUBLIC PERFORMANCE, OR PUBLIC DISPLAY OF OR THROUGH USE OF THIS SOURCE
// CODE WITHOUT THE EXPRESS WRITTEN CONSENT OF GRAB IS STRICTLY PROHIBITED, AND IN VIOLATION OF APPLICABLE LAWS AND
// INTERNATIONAL TREATIES. THE RECEIPT OR POSSESSION OF THIS SOURCE CODE AND/OR RELATED INFORMATION DOES NOT CONVEY
// OR IMPLY ANY RIGHTS TO REPRODUCE, DISCLOSE OR DISTRIBUTE ITS CONTENTS, OR TO MANUFACTURE, USE, OR SELL ANYTHING
// THAT IT MAY DESCRIBE, IN WHOLE OR IN PART.

package gconcurrent

import (
	"context"
	"errors"
	"runtime/debug"
	"sync"
	"time"

	"gitlab.myteksi.net/gophers/go/commons/util/log/logging"
)

/*
 * I have read several implimentation for sync IO call, I combine these to my need
 * Reference:
 *   vendor/gitlab.myteksi.net/gophers/go/commons/util/parallel/gconcurrent/goroutine.go
 *     the Go without timeout, and DD
 *   food/common/concurrencyutils/concurrencyutils.go
 *     got the ideal for extract context
 *   food/food-dataservice/common/goroutine/gowithrecovery.go
 *     good ideal for context, opts
 *   golang.org/x/sync/errgroup
 *     no timeout control
 *   vendor/github.com/myteksi/hystrix-go/hystrix/hystrix.go
 *     too complicated
 */

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
				logging.Error(logTag, "recover trackID: %s, name: %s, recover: %+v", trackID, opt.name, r)
				logging.Error(logTag, "recover trackID: %s, Stack: %+v", trackID, string(debug.Stack()))
			}
			close(runErrChan)
		}()

		runErrChan <- callFun()
	}()

	timer := time.NewTimer(opt.timeoutDuration)
	defer timer.Stop()

	timeBegin := time.Now()

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

	if cmd.err != nil {
		timeEnd := time.Now()
		deadLine, _ := ctx.Deadline()
		logging.Warn(logTag, "trackID: %s, opt: %+v, tags: %+v, duration: %+v, err: %+v, deadLine: %+v",
			trackID, opt, cmd.tags, timeEnd.Sub(timeBegin), cmd.err, deadLine)
	}

	return cmd.tags, cmd.err
}
