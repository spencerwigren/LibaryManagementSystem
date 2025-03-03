package tui

import (
	"github.com/rivo/tview"
)

func App() {
	// This lay out is found: https://github.com/rivo/tview/wiki/Grid
	// Will be adapting it for my project.
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	// TODO have this primitive show a title
	menuPrimitive := func(text string) tview.Primitive {
		menuPrim := tview.NewForm().
			AddInputField("Test", "", 0, nil, nil).
			AddButton("Add Media", nil)

		menuPrim.SetTitle(text).SetTitleAlign(tview.AlignCenter)

		return menuPrim
	}

	menu := menuPrimitive("menu")
	main := newPrimitive("Main content")
	sideBar := newPrimitive("Side bar")

	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0, 30).
		SetBorders(true).
		AddItem(newPrimitive("Header"), 0, 0, 1, 3, 0, 0, false).
		AddItem(newPrimitive("Footer"), 2, 0, 1, 3, 0, 0, false)

		// Layout for screens wider than 100 cells.
	grid.AddItem(menu, 1, 0, 1, 1, 0, 100, true).
		AddItem(main, 1, 1, 1, 1, 0, 100, false).
		AddItem(sideBar, 1, 2, 1, 1, 0, 100, false)

	if err := tview.NewApplication().SetRoot(grid, true).SetFocus(grid).Run(); err != nil {
		panic(err)
	}

}
