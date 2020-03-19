package views

import (
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Paginator interface {
	First()
	Last()
	Prev()
	Next()
}

type PageSizer interface {
	SetPageSize(pageSize int)
}

type Pagination struct {
	First *tview.Button
	Last  *tview.Button
	Prev  *tview.Button
	Next  *tview.Button
}

func PaginationWidget(paginator Paginator) *Pagination {
	var pagination Pagination
	pagination.First = createButton(string('\U000000AB')).
		SetSelectedFunc(func() {
			paginator.First()
		})

	pagination.Last = createButton(string('\U000000BB')).
		SetSelectedFunc(func() {
			paginator.Last()
		})

	pagination.Prev = createButton(string('\U000025C4')).
		SetSelectedFunc(func() {
			paginator.Prev()
		})

	pagination.Next = createButton(string('\U000025BA')).
		SetSelectedFunc(func() {
			paginator.Next()
		})
	return &pagination
}

func PageSizeWidget(pageSizer PageSizer) *tview.DropDown {
	pageSizes := tview.NewDropDown().
		SetLabelColor(tcell.ColorWhite).
		SetFieldTextColor(tcell.ColorWhite).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetLabel("Items per page: ").
		SetOptions([]string{"25", "50", "75", "100"}, func(text string, index int) {
			pageSize, _ := strconv.Atoi(text)
			pageSizer.SetPageSize(pageSize)
		}).SetCurrentOption(0)
	return pageSizes
}

func createButton(label string) *tview.Button {
	button := tview.NewButton(label).
		SetLabelColor(tcell.ColorWhite).
		SetLabelColorActivated(tcell.ColorBlack).
		SetBackgroundColorActivated(tcell.ColorWhite)
	button.SetBackgroundColor(tcell.ColorBlack)
	return button
}
