package modules

import (
	"fmt"

	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
	"github.com/bharath-srinivas/giterm/views"
)

// User represents a github user.
type User struct {
	*views.TextWidget
}

// UserWidget returns a new instance of user widget.
func UserWidget(app *tview.Application, config config.Config) *User {
	widget := views.NewTextView(app, config, true)
	widget.SetWrap(false).SetTitle(string('\U0001F464') + " [green::b]User")
	p := &User{TextWidget: widget}
	go p.Refresh()
	return p
}

// Refresh refreshes the user widget.
func (p *User) Refresh() {
	p.Redraw(p.display)
}

// display renders the user data in a text view.
func (p *User) display() {
	user, _, err := p.Client.Users.Get(p.Context, "")
	if err != nil {
		_, _ = fmt.Fprint(p, "[::b]an error occurred while retrieving user data")
		return
	}

	_, _ = fmt.Fprintf(p, "[gray::b]%s: [white]%s\n", "Name", user.GetName())
	_, _ = fmt.Fprintf(p, "[gray::b]%s: [white]%s\n", "Username", user.GetLogin())
	_, _ = fmt.Fprintf(p, "[gray::b]%s: [white]%s\n", "Company", user.GetCompany())
	_, _ = fmt.Fprintf(p, "[gray::b]%s: [white]%s\n", "Blog", user.GetBlog())
	_, _ = fmt.Fprintf(p, "[gray::b]%s: [white]%s\n", "Location", user.GetLocation())
	_, _ = fmt.Fprintf(p, "[gray::b]%s: [white]%s\n", "Email", user.GetEmail())
	_, _ = fmt.Fprintf(p, "[gray::b]%s: [white]%s\n", "Bio", user.GetBio())
	_, _ = fmt.Fprintf(p, "[gray::b]%s: [white]%d\n", "Public Repos", user.GetPublicRepos())
	_, _ = fmt.Fprintf(p, "[gray::b]%s: [white]%d\n", "Public Gists", user.GetPublicGists())
	_, _ = fmt.Fprintf(p, "[gray::b]%s: [white]%d\n", "Followers", user.GetFollowers())
	_, _ = fmt.Fprintf(p, "[gray::b]%s: [white]%d\n", "Following", user.GetFollowing())
	_, _ = fmt.Fprintf(p, "[gray::b]%s: [white]%d\n", "Owned Private Repos", user.GetOwnedPrivateRepos())
	_, _ = fmt.Fprintf(p, "[gray::b]%s: [white]%d\n", "Total Private Repos", user.GetTotalPrivateRepos())
	_, _ = fmt.Fprintf(p, "[gray::b]%s: [white]%d\n", "Private Gists", user.GetPrivateGists())
}
