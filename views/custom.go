package views

import (
	"regexp"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
	"github.com/rivo/uniseg"
)

// regular expressions.
var (
	boundaryPattern = regexp.MustCompile(`(([,\.\-:;!\?&#+]|\n)[ \t\f\r]*|([ \t\f\r]+))`)
	colorPattern    = regexp.MustCompile(`\[([a-zA-Z]+|#[0-9a-zA-Z]{6}|\-)?(:([a-zA-Z]+|#[0-9a-zA-Z]{6}|\-)?(:([lbdru]+|\-)?)?)?\]`)
	spacePattern    = regexp.MustCompile(`\s+`)
)

// PinnedItemContent represents the pinned item details.
type PinnedItemContent struct {
	Name        string
	Description string
	Lang        string
	Stars       string
	Archived    string
}

// PinnedItem represents a pinned github item.
type PinnedItem struct {
	*tview.Box
	*PinnedItemContent
}

// NewPinnedItem returns an instance of pinned item with borders.
func NewPinnedItem() *PinnedItem {
	return &PinnedItem{
		Box: tview.NewBox().
			SetBorder(true),
	}
}

// SetContent sets the content of the pinned item.
func (p *PinnedItem) SetContent(content *PinnedItemContent) *PinnedItem {
	p.PinnedItemContent = content
	return p
}

// draw draws the primitive onto the screen.
func (p *PinnedItem) draw(screen tcell.Screen) {
	p.Box.Draw(screen)
	x, y, width, height := p.GetInnerRect()

	name := p.PinnedItemContent.Name
	archived := p.PinnedItemContent.Archived
	description := wordWrap(p.PinnedItemContent.Description, width-2)
	lang := p.PinnedItemContent.Lang
	stars := p.PinnedItemContent.Stars

	centerY := y + height/2
	if len(description) > 1 {
		centerY = centerY - len(description)/2
	}

	if tview.TaggedStringWidth(name) > width-2 {
		yPos := y
		for _, line := range wordWrap(name, width-2) {
			tview.Print(screen, "[white::b]"+line, x+1, yPos, width, height, tview.AlignLeft)
			yPos = yPos + 1
		}
	} else {
		tview.Print(screen, name, x+1, y, width, height, tview.AlignLeft)
	}

	tview.Print(screen, archived, x+width-tview.TaggedStringWidth(archived)-1, y, width, height, tview.AlignRight)
	for _, desc := range description {
		tview.Print(screen, "[white]"+desc, x+1, centerY, width-2, height, tview.AlignLeft)
		centerY = centerY + 1
	}
	tview.Print(screen, lang, x+1, y+height-1, width, height, tview.AlignLeft)
	tview.Print(screen, stars, x+width-tview.TaggedStringWidth(stars)-1, y+height-1, width, height, tview.AlignRight)
}

// PinnedView is box which displays the pinned items within their respective box primitives.
type PinnedView struct {
	*tview.Box

	pinnedItems []*PinnedItem
}

// NewPinnedView returns a new pinned view.
func NewPinnedView() *PinnedView {
	return &PinnedView{
		Box: tview.NewBox(),
	}
}

// SetItems sets the pinned items that can be rendered within the pinned view.
func (p *PinnedView) SetItems(items []*PinnedItem) {
	p.pinnedItems = items
}

// Clear removes all the pinned items from the pinned view.
func (p *PinnedView) Clear() {
	p.pinnedItems = nil
}

// Draw draws the primitive onto the screen.
func (p *PinnedView) Draw(screen tcell.Screen) {
	p.Box.Draw(screen)
	if len(p.pinnedItems) == 0 {
		return
	}

	x, y, width, height := p.Box.GetInnerRect()
	width = width / len(p.pinnedItems)
	for _, pinnedItem := range p.pinnedItems {
		pinnedItem.SetRect(x, y, width, height)
		pinnedItem.draw(screen)
		x = x + width
	}
}

// wordWrap splits the string into multiple lines on spaces or boundaries, each not exceeding the given width and returns
//a slice of strings.
func wordWrap(s string, width int) []string {
	s = colorPattern.ReplaceAllString(s, "")
	if len(s) <= width {
		return []string{s}
	}

	var wrapped []string
	for len(s) > 0 {
		extract := runewidth.Truncate(s, width, "")
		if len(extract) == 0 {
			gr := uniseg.NewGraphemes(s)
			gr.Next()
			_, to := gr.Positions()
			extract = s[:to]
		}
		if len(extract) < len(s) {
			// Add any spaces from the next line.
			if spaces := spacePattern.FindStringIndex(s[len(extract):]); spaces != nil && spaces[0] == 0 {
				extract = s[:len(extract)+spaces[1]]
			}

			// Can we split before the mandatory end?
			matches := boundaryPattern.FindAllStringIndex(extract, -1)
			if len(matches) > 0 {
				// Yes. Let's split there.
				extract = extract[:matches[len(matches)-1][1]]
			}
		}
		wrapped = append(wrapped, extract)
		s = s[len(extract):]
	}
	return wrapped
}
