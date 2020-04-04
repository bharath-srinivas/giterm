package modules

import (
	"bytes"
	"io"

	"github.com/eliukblau/pixterm/pkg/ansimage"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
	"github.com/bharath-srinivas/giterm/views"
)

// Avatar represents the github user's avatar.
type Avatar struct {
	*views.TextWidget
}

// AvatarWidget returns an instance of an avatar widget.
func AvatarWidget(app *tview.Application, config config.Config) *Avatar {
	widget := views.NewTextView(app, config, true)
	a := &Avatar{widget}
	go a.Refresh()
	return a
}

// Refresh refreshes the avatar widget.
func (a *Avatar) Refresh() {
	a.Redraw(a.display)
}

// display renders the avatar within the widget.
func (a *Avatar) display() {
	_, _, width, height := a.TextView.GetInnerRect()
	bg, _ := colorful.Hex("#000000")
	sfx, sfy := 1, 2
	go func() {
		avatar, err := ansimage.NewScaledFromURL(user.Viewer.AvatarUrl, sfy*height, sfx*width, bg, ansimage.ScaleModeResize, ansimage.NoDithering)
		if err != nil {
			return
		}

		writer := tview.ANSIWriter(a)
		reader := bytes.NewBufferString(avatar.Render())
		if _, err := io.Copy(writer, reader); err != nil {
			return
		}
	}()
}
