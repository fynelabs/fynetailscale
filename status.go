package fynetailscale

import (
	"context"
	"time"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"tailscale.com/client/tailscale"
)

// NewStatus will return a widget that will update with the current status of tailscale network connection.
func NewStatus(ctx context.Context, lc *tailscale.LocalClient) *widget.Label {
	return widget.NewLabelWithData(NewStatusBinding(ctx, lc))
}

func NewStatusBinding(ctx context.Context, lc *tailscale.LocalClient) binding.String {
	r := binding.NewString()
	r.Set("Connecting...")

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(100 * time.Millisecond):
				status, err := lc.Status(ctx)
				if err != nil {
					r.Set(err.Error())
					continue
				}
				r.Set(status.BackendState)
			}
		}
	}()
	return r
}
