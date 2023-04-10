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

type QRHelper struct {
	widget.BaseWidget

	image     *canvas.Image
	loginLink *widget.Hyperlink
}

var _ fyne.Widget = (*QRHelper)(nil)

func NewQRHelper(u *url.URL) (*QRHelper, error) {
	r := &QRHelper{
		image: &canvas.Image{
			FillMode: canvas.ImageFillOriginal,
		},
		loginLink: widget.NewHyperlink("", nil),
	}
	r.BaseWidget.ExtendBaseWidget(r)
	r.setURL(u)

	return r, nil
}

func (r *QRHelper) CreateRenderer() fyne.WidgetRenderer {
	return qrRenderer{container: container.NewBorder(nil, r.loginLink, nil, nil, r.image)}
}

func (r *QRHelper) SetURL(url *url.URL) error {
	err := r.setURL(url)
	if err != nil {
		return err
	}

	r.Refresh()
	return nil
}

func (r *QRHelper) setURL(u *url.URL) error {
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
