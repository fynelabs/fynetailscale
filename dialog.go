package fynetailscale

import (
	"context"
	"image/color"
	"io"
	"net/url"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"tailscale.com/client/tailscale"
)

type login struct {
	d      dialog.Dialog
	cancel func()
}

var _ io.Closer = (*login)(nil)

// NewLogin will show a dialog that will allow the user to login to tailscale if necessary.
func NewLogin(ctx context.Context, win fyne.Window, lc *tailscale.LocalClient, done func(succeeded bool)) io.Closer {
	cancellable, cancel := context.WithCancel(ctx)

	connecting := container.NewVBox(layout.NewSpacer(), container.NewBorder(nil, nil, widget.NewLabel("Connecting"), nil, widget.NewProgressBarInfinite()), layout.NewSpacer())
	info, _ := NewQRCode(nil)
	info.Hide()
	minSizeRect := canvas.NewRectangle(color.Transparent)
	minSizeRect.SetMinSize(fyne.NewSize(255, 255))
	content := container.NewMax(minSizeRect, connecting, container.NewCenter(info))

	d := dialog.NewCustom("Login", "Cancel", content, win)
	d.SetOnClosed(func() {
		status, err := lc.Status(context.Background())
		if err != nil {
			done(false)
		} else {
			done(status.BackendState == "Running")
		}
	})
	d.Show()

	go func() {
		displayURL := func(targetURL string) error {
			u, err := url.Parse(targetURL)
			if err != nil {
				d.Hide()
				return err
			}
			err = info.SetURL(u)
			if err != nil {
				d.Hide()
				return err
			}
			info.Show()
			connecting.Hide()
			return nil
		}
		defer cancel()

		oldState := ""

		for {
			select {
			case <-cancellable.Done():
				return
			case <-time.After(100 * time.Millisecond):
				status, err := lc.Status(cancellable)
				if err != nil {
					done(false)
					return
				}

				if oldState == status.BackendState {
					continue
				}

				switch status.BackendState {
				case "Running":
					d.Hide()
					return
				case "NeedsLogin":
					if status.AuthURL == "" {
						continue
					}

					err := displayURL(status.AuthURL)
					if err != nil {
						return
					}
				case "NeedsMachineAuth":
					pref, err := lc.GetPrefs(cancellable)
					if err != nil {
						d.Hide()
						return
					}

					if pref.AdminPageURL() == "" {
						continue
					}

					err = displayURL(pref.AdminPageURL())
					if err != nil {
						return
					}
				}
				oldState = status.BackendState
			}
		}
	}()

	return &login{
		d:      d,
		cancel: cancel,
	}
}

// Close will close the dialog.
func (d *login) Close() error {
	d.cancel()
	return nil
}
