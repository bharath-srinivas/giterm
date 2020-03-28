package modules

import (
	"fmt"
	"github.com/google/go-github/v29/github"
	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
	"github.com/bharath-srinivas/giterm/views"
)

// user holds the github user information.
var user *github.User

// User represents a github user.
type User struct {
	*views.TextWidget
}

// UserWidget returns a new instance of user widget.
func UserWidget(app *tview.Application, config config.Config) *User {
	widget := views.NewTextView(app, config, true)
	widget.SetWrap(false).
		SetTitle(string('\U0001F464') + " [green::b]User")
	u := &User{TextWidget: widget}
	go u.Refresh()
	return u
}

// Refresh refreshes the user widget.
func (u *User) Refresh() {
	u.Redraw(u.display)
}

// display renders the user data in a text view.
func (u *User) display() {
	var err error
	user, _, err = u.Client.Users.Get(u.Context, "")
	if err != nil {
		_, _ = fmt.Fprint(u, "[::b]an error occurred while retrieving user data")
		return
	}

	_, _ = fmt.Fprintf(u, "[gray::b]%s: [white]%s\n", "Name", user.GetName())
	_, _ = fmt.Fprintf(u, "[gray::b]%s: [white]%s\n", "Username", user.GetLogin())
	_, _ = fmt.Fprintf(u, "[gray::b]%s: [white]%s\n", "Bio", user.GetBio())
	_, _ = fmt.Fprintf(u, "[gray::b]%s: [white]%s\n", "Company", user.GetCompany())
	_, _ = fmt.Fprintf(u, "[gray::b]%s: [white]%s\n", "Location", user.GetLocation())
	_, _ = fmt.Fprintf(u, "[gray::b]%s: [white]%s\n", "Email", user.GetEmail())
	_, _ = fmt.Fprintf(u, "[gray::b]%s: [white]%s\n", "Blog", user.GetBlog())
}

// UserStats represents various stats of a github user.
type UserStats struct {
	*views.TextWidget
}

// UserStatsWidget returns a new instance of user stats widget.
func UserStatsWidget(app *tview.Application, config config.Config) *UserStats {
	widget := views.NewTextView(app, config, true)
	widget.SetWrap(false).
		SetTitle("[green::b]User stats")
	u := &UserStats{TextWidget: widget}
	go u.Refresh()
	return u
}

// Refresh refreshes the user stats widget.
func (u *UserStats) Refresh() {
	u.Redraw(u.display)
}

// display renders the user data in a text view.
func (u *UserStats) display() {
	if user == nil {
		_, _ = fmt.Fprint(u, "[::b]an error occurred while retrieving user data")
		return
	}

	organizations, _, err := u.Client.Organizations.List(u.Context, "", &github.ListOptions{
		Page:    1,
		PerPage: 100,
	})
	if err != nil {
		_, _ = fmt.Fprint(u, "[::b]an error occurred while retrieving user data")
		return
	}

	orgCount := len(organizations)
	totalRepos := user.GetPublicRepos() + user.GetTotalPrivateRepos()
	totalGists := user.GetPublicGists() + user.GetPrivateGists()
	_, _ = fmt.Fprintf(u, "[gray::b]%s: [white]%d\n", "Repositories", totalRepos)
	_, _ = fmt.Fprintf(u, "[gray::b]%s: [white]%d\n", "Gists", totalGists)
	_, _ = fmt.Fprintf(u, "[gray::b]%s: [white]%d\n", "Organizations", orgCount)
	_, _ = fmt.Fprintf(u, "[gray::b]%s: [white]%d\n", "Followers", user.GetFollowers())
	_, _ = fmt.Fprintf(u, "[gray::b]%s: [white]%d\n", "Following", user.GetFollowing())
}
