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
