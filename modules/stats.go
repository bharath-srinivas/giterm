package modules

import (
	"fmt"

	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
	"github.com/bharath-srinivas/giterm/views"
)

// UserStats represents various stats of a github user.
type UserStats struct {
	*views.TextWidget
}

// UserStatsWidget returns a new instance of user stats widget.
func UserStatsWidget(app *tview.Application, config config.Config) *UserStats {
	widget := views.NewTextView(app, config, true)
	widget.SetWrap(false).
		SetTitle("[green::b]User stats")
	u := &UserStats{widget}
	go u.Refresh()
	return u
}

// Refresh refreshes the user stats widget.
func (u *UserStats) Refresh() {
	u.Redraw(u.display)
}

// display renders the user stats data in a text view.
func (u *UserStats) display() {
	if user == nil {
		_, _ = fmt.Fprint(u, "[::b]an error occurred while retrieving user data")
		return
	}

	_, _ = fmt.Fprintf(u, "[gray::b]%s: [white]%d\n", "Repositories", user.Viewer.Repositories.TotalCount)
	_, _ = fmt.Fprintf(u, "[gray::b]%s: [white]%d\n", "Gists", user.Viewer.Gists.TotalCount)
	_, _ = fmt.Fprintf(u, "[gray::b]%s: [white]%d\n", "Organizations", user.Viewer.Organizations.TotalCount)
	_, _ = fmt.Fprintf(u, "[gray::b]%s: [white]%d\n", "Starred", user.Viewer.StarredRepositories.TotalCount)
	_, _ = fmt.Fprintf(u, "[gray::b]%s: [white]%d\n", "Followers", user.Viewer.Followers.TotalCount)
	_, _ = fmt.Fprintf(u, "[gray::b]%s: [white]%d\n", "Following", user.Viewer.Following.TotalCount)
}
