package fynetailscale

import (
	"context"
	"time"

	"fyne.io/fyne/v2/widget"
	"tailscale.com/client/tailscale"
)

func NewStatus(ctx context.Context, lc *tailscale.LocalClient) *widget.Label {
	r := widget.NewLabel("Connecting...")
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(100 * time.Millisecond):
				status, err := lc.Status(ctx)
				if err != nil {
					r.SetText(err.Error())
					continue
				}
				r.SetText(status.BackendState)
			}
		}
	}()
	return r
}
