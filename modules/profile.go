package modules

import (
	"encoding/json"
	"fmt"

	"github.com/rivo/tview"
)

type profile struct {
	Name              string `json:"name,omitempty"`
	UserName          string `json:"login,omitempty"`
	Company           string `json:"company,omitempty"`
	Blog              string `json:"blog,omitempty"`
	Location          string `json:"location,omitempty"`
	Email             string `json:"email,omitempty"`
	Hireable          bool   `json:"hireable,omitempty"`
	Bio               string `json:"bio,omitempty"`
	PublicRepos       int    `json:"public_repos,omitempty"`
	PublicGists       int    `json:"public_gists,omitempty"`
	Followers         int    `json:"followers,omitempty"`
	Following         int    `json:"following,omitempty"`
	OwnedPrivateRepos int    `json:"owned_private_repos,omitempty"`
	TotalPrivateRepos int    `json:"total_private_repos,omitempty"`
	PrivateGists      int    `json:"private_gists,omitempty"`
}

func (g *GitApp) ProfilePage() *Page {
	widget := tview.NewTextView()
	widget.SetBorder(true)
	widget.SetDynamicColors(true)
	widget.SetScrollable(true)
	widget.SetTitle(string('\U0001F642') + " [green::b]Profile")

	user, _, err := g.Client.Users.Get(g.Context, "")
	if err != nil {
		_, _ = fmt.Fprint(widget, "")
		return &Page{}
	}

	var m map[string]interface{}
	userJson, _ := json.Marshal(user)
	_ = json.Unmarshal(userJson, &m)

	var p profile
	_ = json.Unmarshal(userJson, &p)

	_, _ = fmt.Fprintf(widget, "[gray::b]%s: [white]%s\n", "Name", p.Name)
	_, _ = fmt.Fprintf(widget, "[gray::b]%s: [white]%s\n", "Username", p.UserName)
	_, _ = fmt.Fprintf(widget, "[gray::b]%s: [white]%s\n", "Company", p.Company)
	_, _ = fmt.Fprintf(widget, "[gray::b]%s: [white]%s\n", "Blog", p.Blog)
	_, _ = fmt.Fprintf(widget, "[gray::b]%s: [white]%s\n", "Location", p.Location)
	_, _ = fmt.Fprintf(widget, "[gray::b]%s: [white]%s\n", "Email", p.Email)
	_, _ = fmt.Fprintf(widget, "[gray::b]%s: [white]%v\n", "Hireable", p.Hireable)
	_, _ = fmt.Fprintf(widget, "[gray::b]%s: [white]%s\n", "Bio", p.Bio)
	_, _ = fmt.Fprintf(widget, "[gray::b]%s: [white]%d\n", "Public Repos", p.PublicRepos)
	_, _ = fmt.Fprintf(widget, "[gray::b]%s: [white]%d\n", "Public Gists", p.PublicGists)
	_, _ = fmt.Fprintf(widget, "[gray::b]%s: [white]%d\n", "Followers", p.Followers)
	_, _ = fmt.Fprintf(widget, "[gray::b]%s: [white]%d\n", "Following", p.Following)
	_, _ = fmt.Fprintf(widget, "[gray::b]%s: [white]%d\n", "Owned Private Repos", p.OwnedPrivateRepos)
	_, _ = fmt.Fprintf(widget, "[gray::b]%s: [white]%d\n", "Total Private Repos", p.TotalPrivateRepos)
	_, _ = fmt.Fprintf(widget, "[gray::b]%s: [white]%d\n", "Private Gists", p.PrivateGists)

	return &Page{
		Name:            "Profile",
		Parent:          widget,
		ChildComponents: nil,
	}
}
