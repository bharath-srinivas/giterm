package modules

import (
	"fmt"
	"regexp"

	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
	"github.com/bharath-srinivas/giterm/views"
)

// userQuery represents a graphql query.
type userQuery *struct {
	Viewer struct {
		AvatarUrl                string
		Name                     string
		Login                    string
		Bio                      string
		Company                  string
		Location                 string
		Email                    string
		WebsiteUrl               string
		IsDeveloperProgramMember bool
		Followers                struct {
			TotalCount int
		}
		Following struct {
			TotalCount int
		}
		Organizations struct {
			TotalCount int
		}
		StarredRepositories struct {
			TotalCount int
		}
		Repositories struct {
			TotalCount int
		}
		Gists struct {
			TotalCount int
		}
		PinnedItems struct {
			TotalCount int
			Nodes      []struct {
				Repository struct {
					Name          string
					NameWithOwner string
					Description   string
					IsArchived    bool
					Owner         struct {
						Login string
					}
					PrimaryLanguage *struct {
						Name  string
						Color string
					}
					Stargazers struct {
						TotalCount int
					}
				} `graphql:"... on Repository"`
				Gist struct {
					Description string
				} `graphql:"... on Gist"`
			}
		} `graphql:"pinnedItems(first:6)"`
	}
}

// user holds the github user information.
var user userQuery

// User represents a github user.
type User struct {
	*views.TextWidget
}

// UserWidget returns a new instance of user widget.
func UserWidget(app *tview.Application, config config.Config) *User {
	widget := views.NewTextView(app, config, true)
	widget.SetWrap(true).
		SetWordWrap(true).
		SetTitle(string('\U0001F464') + " [green::b]User")
	u := &User{widget}
	go u.Refresh()
	return u
}

// Refresh refreshes the user widget.
func (u *User) Refresh() {
	u.Redraw(u.display)
}

// display renders the user information in a text view.
func (u *User) display() {
	if err := u.GqlClient.Query(u.Context, &user, nil); err != nil {
		_, _ = fmt.Fprint(u, "[::b]an error occurred while retrieving user data")
		return
	}

	bio := user.Viewer.Bio
	r := regexp.MustCompile(`([@A-Z])\b\w+`)
	bio = r.ReplaceAllStringFunc(bio, func(s string) string {
		return "[::b]" + s + "[::-]"
	})

	_, _ = fmt.Fprintf(u, formatText("", user.Viewer.Name, "[::b]", "\n"))
	_, _ = fmt.Fprintf(u, formatText("", user.Viewer.Login, "[gray::d]", "\n\n"))
	_, _ = fmt.Fprintf(u, formatText("", bio, "", "\n\n"))
	_, _ = fmt.Fprintf(u, formatText(string('\U0001F465'), user.Viewer.Company, "[::b]", "\n\n"))
	_, _ = fmt.Fprintf(u, formatText(string('\U0001F4CD'), user.Viewer.Location, "", "\n\n"))
	_, _ = fmt.Fprintf(u, formatText(string('\u2709'), user.Viewer.Email, "", "\n\n"))
	_, _ = fmt.Fprintf(u, formatText(string('\U0001F517'), user.Viewer.WebsiteUrl, "", "\n\n"))
	if user.Viewer.IsDeveloperProgramMember {
		_, _ = fmt.Fprintln(u, "[::b]Developer Program Member")
	}
}

// formatText is a helper function that returns a formatted string with the provided title, prefix and suffix.
func formatText(title, text, prefix, suffix string) string {
	if text == "" {
		return ""
	}
	if title == "" {
		return prefix + text + suffix
	}
	return title + " " + prefix + text + suffix
}
