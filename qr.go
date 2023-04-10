package fynetailscale

import (
	"image"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	qrcode "github.com/skip2/go-qrcode"
)

// QRCode is a widget that displays a QR code and a link to the URL it represents.
type QRCode struct {
	widget.BaseWidget

	image     *canvas.Image
	loginLink *widget.Hyperlink
}

var _ fyne.Widget = (*QRCode)(nil)

// NewQRCode creates a new QRCode widget.
func NewQRCode(u *url.URL) (*QRCode, error) {
	r := &QRCode{
		image: &canvas.Image{
			FillMode: canvas.ImageFillOriginal,
		},
		loginLink: widget.NewHyperlink("", nil),
	}
	r.BaseWidget.ExtendBaseWidget(r)
	r.setURL(u)

	return r, nil
}

// CreateRenderer implements fyne.Widget.
func (r *QRCode) CreateRenderer() fyne.WidgetRenderer {
	return qrRenderer{container: container.NewBorder(nil, r.loginLink, nil, nil, r.image)}
}

// SetURL sets the URL to display.
func (r *QRCode) SetURL(url *url.URL) error {
	err := r.setURL(url)
	if err != nil {
		return err
	}

	r.Refresh()
	return nil
}

func (r *QRCode) setURL(u *url.URL) error {
	if u == nil {
		r.image.Image = image.NewRGBA(image.Rect(0, 0, 256, 256))
		r.loginLink.Text = ""
		r.loginLink.URL = nil
		return nil
	}

	qr, err := qrcode.New(u.String(), qrcode.Medium)
	if err != nil {
		return err
	}

	r.image.Image = qr.Image(256)
	r.loginLink.Text = "Connect to: " + u.String()
	r.loginLink.URL = u

	return nil
}

type qrRenderer struct {
	container *fyne.Container
}

var _ fyne.WidgetRenderer = (*qrRenderer)(nil)

func (r qrRenderer) Destroy() {
}

func (r qrRenderer) Layout(size fyne.Size) {
	r.container.Resize(size)
}

func (r qrRenderer) MinSize() fyne.Size {
	return r.container.MinSize()
}

func (r qrRenderer) Objects() []fyne.CanvasObject {
	return r.container.Objects
}

func (r qrRenderer) Refresh() {
	r.container.Refresh()
}
