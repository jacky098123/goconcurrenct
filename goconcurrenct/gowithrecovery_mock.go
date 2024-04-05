package goconcurrenct

import "context"

// SetupMock ...
func SetupMock() (restore func()) {
	original := GoWithRecovery

	GoWithRecovery = goWithRecoveryMock

	restore = func() {
		GoWithRecovery = original
	}
	return
}

// goWithRecoveryMock is a mock for goWithRecovery for UT, should be setup in test.Main func
// this func is using sequential call
func goWithRecoveryMock(ctx context.Context, trackID string, callFun func() error, ops ...Option) ([]string, error) {
	err := callFun()
	return nil, err
}
