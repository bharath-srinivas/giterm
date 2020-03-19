package modules

import (
	"fmt"

	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
	"github.com/bharath-srinivas/giterm/views"
)

type Profile struct {
	*views.TextWidget
}

func ProfileWidget(app *tview.Application, config config.Config) *Profile {
	widget := views.NewTextView(app, config, true)
	widget.TextView.SetTitle(string('\U0001F642') + " [green::b]Profile")
	p := &Profile{TextWidget: widget}
	go p.Refresh()
	return p
}

func (p *Profile) Refresh() {
	p.Redraw(p.display)
}

func (p *Profile) display() {
	user, _, err := p.Client.Users.Get(p.Context, "")
	if err != nil {
		_, _ = fmt.Fprint(p.TextView, "[::b]an error occurred while retrieving profile")
		return
	}

	_, _ = fmt.Fprintf(p.TextView, "[gray::b]%s: [white]%s\n", "Name", user.GetName())
	_, _ = fmt.Fprintf(p.TextView, "[gray::b]%s: [white]%s\n", "Username", user.GetLogin())
	_, _ = fmt.Fprintf(p.TextView, "[gray::b]%s: [white]%s\n", "Company", user.GetCompany())
	_, _ = fmt.Fprintf(p.TextView, "[gray::b]%s: [white]%s\n", "Blog", user.GetBlog())
	_, _ = fmt.Fprintf(p.TextView, "[gray::b]%s: [white]%s\n", "Location", user.GetLocation())
	_, _ = fmt.Fprintf(p.TextView, "[gray::b]%s: [white]%s\n", "Email", user.GetEmail())
	_, _ = fmt.Fprintf(p.TextView, "[gray::b]%s: [white]%s\n", "Bio", user.GetBio())
	_, _ = fmt.Fprintf(p.TextView, "[gray::b]%s: [white]%d\n", "Public Repos", user.GetPublicRepos())
	_, _ = fmt.Fprintf(p.TextView, "[gray::b]%s: [white]%d\n", "Public Gists", user.GetPublicGists())
	_, _ = fmt.Fprintf(p.TextView, "[gray::b]%s: [white]%d\n", "Followers", user.GetFollowers())
	_, _ = fmt.Fprintf(p.TextView, "[gray::b]%s: [white]%d\n", "Following", user.GetFollowing())
	_, _ = fmt.Fprintf(p.TextView, "[gray::b]%s: [white]%d\n", "Owned Private Repos", user.GetOwnedPrivateRepos())
	_, _ = fmt.Fprintf(p.TextView, "[gray::b]%s: [white]%d\n", "Total Private Repos", user.GetTotalPrivateRepos())
	_, _ = fmt.Fprintf(p.TextView, "[gray::b]%s: [white]%d\n", "Private Gists", user.GetPrivateGists())
}
