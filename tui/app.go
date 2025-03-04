package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

func App() {
	// This lay out is found: https://github.com/rivo/tview/wiki/Grid
	// Will be adapting it for my project.

	app := tview.NewApplication()

	// Setting own things
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	// Layout of TUI
	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0, 30).
		SetBorders(true).
		AddItem(newPrimitive("Header"), 0, 0, 1, 3, 0, 0, false).
		AddItem(newPrimitive("Footer"), 2, 0, 1, 3, 0, 0, false)

	// Pop up for adding media to app
	addMediaModal := tview.NewModal().
		AddButtons([]string{"Quit", "Submit"}).
		SetText("Add Media").
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				app.SetRoot(grid, true)
			} else if buttonLabel == "Submit" {
				fmt.Println("Items Submit")
			}
		})

	// Main menu prmitive to interact with TUI
	menuPrimitive := func() tview.Primitive {
		menuPrim := tview.NewForm().
			// TODO create function to seach data base
			AddInputField("Search", "", 0, nil, nil).
			AddButton("Add Media", func() {
				app.SetRoot(addMediaModal, true)

			}).
			AddButton("Quit", func() {
				app.Stop()
			})

		return menuPrim
	}

	// Setting more layout in grid
	menu := menuPrimitive()
	main := newPrimitive("Main content")
	sideBar := newPrimitive("Side bar")

	// adding menu, main, and sidebar to grid to draw.
	grid.AddItem(menu, 1, 0, 1, 1, 0, 100, true).
		AddItem(main, 1, 1, 1, 1, 0, 100, false).
		AddItem(sideBar, 1, 2, 1, 1, 0, 100, false)

	if err := app.SetRoot(grid, true).SetFocus(menu).Run(); err != nil {
		panic(err)
	}

}
