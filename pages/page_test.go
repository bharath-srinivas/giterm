package pages

import (
	"reflect"
	"testing"

	"github.com/rivo/tview"

	"github.com/bharath-srinivas/giterm/config"
)

func testPage(config config.Config, pageName string) *Page {
	return MakePage(tview.NewApplication(), config, pageName)
}

func testPages() Pages {
	return MakePages(tview.NewApplication(), config.Config{})
}

func TestMakePage(t *testing.T) {
	tests := []struct {
		name     string
		config   config.Config
		pageName string
		expected string
	}{
		{
			name:     "invalid page",
			config:   config.Config{},
			pageName: "",
			expected: "",
		},
		{
			name:     "valid page",
			config:   config.Config{},
			pageName: "feeds",
			expected: "Feeds",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := testPage(tt.config, tt.pageName)
			actual := page.Name
			if actual != tt.expected {
				t.Errorf("MakePage returned an incorrect page, got: %s, want: %s", actual, tt.expected)
			}
		})
	}
}

func TestPage_PrevWidget(t *testing.T) {
	tests := []struct {
		app      *tview.Application
		name     string
		pageName string
		focusIdx int
		expected string
	}{
		{
			app:      tview.NewApplication(),
			name:     "without children",
			pageName: "",
			expected: "tview.Primitive",
		},
		{
			app:      tview.NewApplication(),
			name:     "with children",
			pageName: "feeds",
			focusIdx: 0,
			expected: "tview.Primitive",
		},
		{
			app:      tview.NewApplication(),
			name:     "with children",
			pageName: "feeds",
			focusIdx: 2,
			expected: "tview.Primitive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := testPage(config.Config{}, tt.pageName)
			if tt.focusIdx > 0 {
				tt.app.SetFocus(page.Children[tt.focusIdx])
			}
			primitiveType := reflect.TypeOf((*tview.Primitive)(nil)).Elem()
			actual := reflect.TypeOf(page.PrevWidget())
			if !actual.Implements(primitiveType) {
				t.Errorf("PrevWidget returned an incorrect type, got: %v, want: %v", actual, tt.expected)
			}
		})
	}
}

func TestPage_NextWidget(t *testing.T) {
	tests := []struct {
		name     string
		pageName string
		expected string
	}{
		{
			name:     "without children",
			pageName: "",
			expected: "tview.Primitive",
		},
		{
			name:     "with children",
			pageName: "feeds",
			expected: "tview.Primitive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := testPage(config.Config{}, tt.pageName)
			primitiveType := reflect.TypeOf((*tview.Primitive)(nil)).Elem()
			actual := reflect.TypeOf(page.NextWidget())
			if !actual.Implements(primitiveType) {
				t.Errorf("NextWidget returned an incorrect type, got: %v, want: %v", actual, tt.expected)
			}
		})
	}
}

func TestMakePages(t *testing.T) {
	pages := testPages()
	actual := reflect.TypeOf(pages).String()
	expected := "pages.Pages"
	if actual != expected {
		t.Errorf("MakePages returned an incorrect type, got: %s, want: %s", actual, expected)
	}
}

func TestPages_Get(t *testing.T) {
	tests := []struct {
		name     string
		pageName string
		expected string
	}{
		{
			name:     "invalid page",
			pageName: "",
			expected: "nil",
		},
		{
			name:     "valid page",
			pageName: "Feeds",
			expected: "*pages.Page",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pages := testPages()
			page := pages.Get(tt.pageName)
			actual := reflect.TypeOf(page).String()
			if page == nil {
				actual = "nil"
			}
			if actual != tt.expected {
				t.Errorf("Get returned an incorrect type, got:%s, want: %s", actual, tt.expected)
			}
		})
	}
}

func TestPages_Prev(t *testing.T) {
	tests := []struct {
		name     string
		pageName string
		expected string
	}{
		{
			name:     "invalid page",
			pageName: "",
			expected: "",
		},
		{
			name:     "valid page",
			pageName: "Feeds",
			expected: "Repos",
		},
		{
			name:     "valid page",
			pageName: "Profile",
			expected: "Feeds",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pages := testPages()
			actual := pages.Prev(tt.pageName)
			if actual != tt.expected {
				t.Errorf("Prev returned an incorrect page name, got: %s, want: %s", actual, tt.expected)
			}
		})
	}
}

func TestPages_Next(t *testing.T) {
	tests := []struct {
		name     string
		pageName string
		expected string
	}{
		{
			name:     "invalid page",
			pageName: "",
			expected: "",
		},
		{
			name:     "valid page",
			pageName: "Feeds",
			expected: "Profile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pages := testPages()
			actual := pages.Next(tt.pageName)
			if actual != tt.expected {
				t.Errorf("Next returned an incorrect page name, got: %s, want: %s", actual, tt.expected)
			}
		})
	}
}
