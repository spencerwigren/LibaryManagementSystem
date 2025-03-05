package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func App() {
	// This lay out is found: https://github.com/rivo/tview/wiki/Grid
	// Will be adapting it for my project.

	app := tview.NewApplication()
	pages := tview.NewPages()

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

	addBookMedia(pages)
	mediaModal(pages)

	// Main menu prmitive to interact with TUI
	menuPrimitive := func() tview.Primitive {
		menuPrim := tview.NewForm().
			// TODO create function to seach data base
			AddInputField("Search", "", 0, nil, nil).
			AddButton("Add Media", func() {
				// pages.SwitchToPage("addMedia")
				pages.SwitchToPage("mediaModal")

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

	pages.AddPage("mainMenu", grid, true, true)

	if err := app.SetRoot(pages, true).SetFocus(menu).Run(); err != nil {
		panic(err)
	}
}

func mediaModal(pages *tview.Pages) {
	// Pop up for adding media to app
	addMediaModal := tview.NewModal().
		AddButtons([]string{"Add Books", "Back"}).
		SetText("Add Media").
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Add Books" {
				pages.SwitchToPage("addMedia")
			} else if buttonLabel == "Back" {
				pages.SwitchToPage("mainMenu")
			}
		})

	addMediaModal.SetBackgroundColor(tcell.ColorBlack)

	pages.AddPage("mediaModal", addMediaModal, true, false)
}

func addBookMedia(pages *tview.Pages) {
	// TODO make this for the add Book form for addMediaForm
	addMediaForm := tview.NewForm().
		AddInputField("Input Title: ", "", 0, nil, nil).
		AddInputField("Input Page Number: ", "", 0, nil, nil).
		AddInputField("Insert Author Name: ", "", 0, nil, nil).
		AddButton("Submit", func() {
			fmt.Println("Items Submited")
		}).
		AddButton("Back", func() {
			// app.SetRoot(grid, true)
			pages.SwitchToPage("mediaModal")
		})

	// This has to be set out side of the tview.NewForm() for the input fields and buttons to show
	addMediaForm.SetBorder(true).SetTitle("Add Media").SetTitleAlign(tview.AlignCenter)

	// Setting the addMediaForm to be shown
	addMediaFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(addMediaForm, 0, 1, true)

	pages.AddPage("addMedia", addMediaFlex, true, false)
}
