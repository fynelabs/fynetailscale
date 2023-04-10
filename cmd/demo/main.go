package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/fynelabs/fynetailscale"
	"tailscale.com/tsnet"
)

func main() {
	s := new(tsnet.Server)
	defer s.Close()

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

	lc, err := s.LocalClient()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a := app.NewWithID("com.fynelabs.tailscale.demo")
	w := a.NewWindow("Tailscale Demo")
	w.Resize(fyne.NewSize(800, 600))

	status := fynetailscale.NewStatus(ctx, lc)

	body := widget.NewEntry()
	body.Disable()

	request := widget.NewEntry()
	request.OnSubmitted = func(text string) {
		cli := s.HTTPClient()
		resp, err := cli.Get(text)
		if err != nil {
			body.SetText(err.Error())
			return
		}
		defer resp.Body.Close()
		remote, err := io.ReadAll(resp.Body)
		if err != nil {
			body.SetText(err.Error())
			return
		}
		body.SetText(string(remote))
	}

	w.SetContent(container.NewBorder(container.NewBorder(nil, nil, widget.NewLabel("URL in tailnet"), nil, request), status, nil, nil, body))

	d := fynetailscale.NewLogin(ctx, w, lc, func(succeeded bool) {
		if succeeded {
			fmt.Println("Connected")
		} else {
			fmt.Println("Failed to connect")
		}
	})
	defer d.Close()

	w.ShowAndRun()
}
