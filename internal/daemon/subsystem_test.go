package daemon

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tofutf/tofutf/internal"
	"github.com/tofutf/tofutf/internal/xslog"
	"golang.org/x/sync/errgroup"
)

func TestSubsystem(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		exclusive bool
	}{
		{"backoff", false},
		{"backoff and wait and lock", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := &Subsystem{
				Name:      tt.name,
				System:    &fakeStartable{},
				Logger:    slog.New(&xslog.NoopHandler{}),
				Exclusive: tt.exclusive,
			}
			if tt.exclusive {
				sub.DB = &fakeWaitAndLock{}
				sub.LockID = internal.Int64(123)
			}
			err := sub.Start(ctx, &errgroup.Group{})
			require.NoError(t, err)
		})
	}
}

type (
	fakeStartable   struct{}
	fakeWaitAndLock struct{}
)

func (f *fakeStartable) Start(ctx context.Context) error {
	return nil
}

func (f *fakeWaitAndLock) WaitAndLock(ctx context.Context, id int64, fn func(context.Context) error) error {
	return fn(ctx)
}
