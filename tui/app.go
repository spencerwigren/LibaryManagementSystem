package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func App() {
	app := tview.NewApplication()
	mainFlex := tview.NewFlex().
		SetDirection(tview.FlexRow)

	helpBox := tview.NewTextView().
		SetDynamicColors(true).
		SetWordWrap(true).
		SetBorder(true).SetTitle("Media")

	// Form for function of app
	inputForm := tview.NewForm().
		AddButton("Add Media", func() {
			addMediaView(app, mainFlex)
		}).
		AddButton("Exit", func() {
			app.Stop()
		})
	inputForm.SetBorder(true).SetTitle("Input Media")

	mainFlex.
		AddItem(helpBox, 0, 1, false).
		AddItem(inputForm, 0, 1, true)

	if err := app.SetRoot(mainFlex, true).Run(); err != nil {
		panic(err)
	}

}

func addMediaView(app *tview.Application, mainFlex *tview.Flex) {

	addMedia := tview.NewForm().
		AddButton("Add Book", func() {
			// TODO convert page number for str to int
			inputField1 := tview.NewInputField().SetLabel("Input Book Title: ")
			inputField2 := tview.NewInputField().SetLabel("Input Number of Pages: ")
			inputField3 := tview.NewInputField().SetLabel("Input Author Name: ")

			backButton := tview.NewButton("Back").SetSelectedFunc(func() {
				mainFlex.Clear()
			})

			addBookForm := tview.NewForm().
				AddFormItem(inputField1).
				AddFormItem(inputField2).
				AddFormItem(inputField3)

			mainFlex.Clear().SetDirection(tview.FlexRow).
				AddItem(addBookForm, 0, 1, true).
				AddItem(backButton, 0, 1, false)
			app.SetFocus(addBookForm)

			inputField1.SetDoneFunc(func(key tcell.Key) {
				if key == tcell.KeyTab {
					app.SetFocus(inputField2)
				}
			})

			inputField2.SetDoneFunc(func(key tcell.Key) {
				if key == tcell.KeyTab {
					app.SetFocus(inputField3)
				}
			})

			inputField3.SetDoneFunc(func(key tcell.Key) {
				if key == tcell.KeyTab {
					app.SetFocus(backButton)
				}
			})
		})

	mainFlex.Clear().SetDirection(tview.FlexRow).
		AddItem(addMedia, 0, 1, true)
	app.SetFocus(addMedia)

}

// func Tui() {
// 	app := tview.NewApplication()

// 	helpBox := tview.NewTextView().
// 		SetDynamicColors(true).
// 		SetWordWrap(true).
// 		SetBorder(true).SetTitle("Media")

// 	// Setting root and nodes of tree
// 	root := tview.NewTreeNode("Base").
// 		AddChild(tview.NewTreeNode("Book").SetReference("book"))

// 		// Creating Tree View
// 	// Create TreeView
// 	treeView := tview.NewTreeView().
// 		SetRoot(root).
// 		SetCurrentNode(root)

// 	// Form for function of app
// 	inputForm := tview.NewForm().
// 		AddButton("Exit", func() {
// 			app.Stop()
// 		}).
// 		SetBorder(true).SetTitle("Input Media")

// 	mainFlex := tview.NewFlex().
// 		AddItem(helpBox, 0, 1, false).
// 		AddItem(inputForm, 0, 1, false).
// 		AddItem(treeView, 0, 1, true)

// 	openBookInputFields := func(nodeName string) {
// 		// Setting input fields for Book
// 		inputBookField1 := tview.NewInputField().SetLabel(nodeName + "Input Title: ")
// 		inputBookField2 := tview.NewInputField().SetLabel(nodeName + "Input Page Number: ")
// 		inputBookField3 := tview.NewInputField().SetLabel(nodeName + "Input Author Name: ")

// 		backButton := tview.NewButton("Back").SetSelectedFunc(func() {
// 			mainFlex.Clear().
// 				// Will need to add all main flex views here
// 				AddItem(helpBox, 0, 1, false).
// 				AddItem(inputForm, 0, 1, false).
// 				AddItem(treeView, 0, 1, true)

// 			app.SetFocus(treeView)

// 		})

// 		// Creating New Layout for input fields
// 		inputFlex := tview.NewFlex().
// 			SetDirection(tview.FlexRow).
// 			AddItem(inputBookField1, 1, 1, true).
// 			AddItem(inputBookField2, 1, 1, true).
// 			AddItem(inputBookField3, 1, 1, true).
// 			AddItem(backButton, 1, 1, true)

// 		// Switch to input view
// 		mainFlex.Clear().
// 			AddItem(inputFlex, 0, 1, true)
// 		app.SetFocus(inputBookField1)

// 		inputBookField1.SetDoneFunc(func(key tcell.Key) {
// 			if key == tcell.KeyTab {
// 				app.SetFocus(inputBookField2)
// 			}
// 		})

// 		inputBookField2.SetDoneFunc(func(key tcell.Key) {
// 			if key == tcell.KeyTab {
// 				app.SetFocus(inputBookField3)
// 			}
// 		})

// 		inputBookField3.SetDoneFunc(func(key tcell.Key) {
// 			if key == tcell.KeyTab {
// 				app.SetFocus(backButton)
// 			}
// 		})

// 	}

// 	treeView.SetSelectedFunc(func(node *tview.TreeNode) {
// 		if ref := node.GetReference(); ref != nil {
// 			nodeName := ref.(string)
// 			openBookInputFields(nodeName)
// 		}
// 	})

// 	if err := app.SetRoot(mainFlex, true).Run(); err != nil {
// 		panic(err)
// 	}

// }
