package tui

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"Libarymanagementsystem/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// func App(db *sql.DB) {
// 	app := tview.NewApplication()
// 	pages := tview.NewPages()

// 	// Function to create simple text-based primitives
// 	newPrimitive := func(text string) tview.Primitive {
// 		return tview.NewTextView().
// 			SetTextAlign(tview.AlignCenter).
// 			SetText(text)
// 	}

// 	// Layout grid
// 	grid := tview.NewGrid().
// 		SetRows(3, 0, 3).
// 		SetColumns(30, 0, 30).
// 		SetBorders(true).
// 		AddItem(newPrimitive("Header"), 0, 0, 1, 3, 0, 0, false).
// 		AddItem(newPrimitive("Footer"), 2, 0, 1, 3, 0, 0, false)

// 	// Input field for searching
// 	searchInput := tview.NewInputField().SetLabel("Search: ")

// 	// Create a single instance of the main TextView to update later**
// 	mainTextView := tview.NewTextView().
// 		SetTextAlign(tview.AlignCenter).
// 		SetWordWrap(true).
// 		SetText("Results will appear here")

// 	// Function to process search results and update the main text view
// 	updateMainView := func(searchText string) {
// 		if searchText == "" {
// 			return
// 		}

// 		// Run search in a Goroutine
// 		go func() {
// 			searchResults, err := utils.SearchTables(db, searchText)
// 			if err != nil {
// 				log.Printf("Error searching: %v", err)
// 				return
// 			}

// 			// Convert results to string
// 			resultsText := mainPrimitiveResults(searchResults)

// 			// **Update the UI safely**
// 			app.QueueUpdateDraw(func() {
// 				mainTextView.SetText(resultsText)
// 				searchInput.SetText("") // Clear input field
// 			})
// 		}()
// 	}

// 	// **Menu form that handles user input**
// 	menu := tview.NewForm().
// 		AddFormItem(searchInput).
// 		AddButton("Search", func() {
// 			updateMainView(searchInput.GetText()) // Call update function
// 		}).
// 		AddButton("Add Media", func() {
// 			pages.SwitchToPage("mediaModal")
// 		}).
// 		AddButton("Quit", func() {
// 			app.Stop()
// 		})

// 	// Side bar
// 	sideBar := newPrimitive("Side bar")

// 	// Add items to grid
// 	grid.AddItem(menu, 1, 0, 1, 1, 0, 100, true).
// 		AddItem(mainTextView, 1, 1, 1, 1, 0, 100, false). // **Use mainTextView here**
// 		AddItem(sideBar, 1, 2, 1, 1, 0, 100, false)

// 	// Add pages and start app
// 	pages.AddPage("mainMenu", grid, true, true)

// 	if err := app.SetRoot(pages, true).SetFocus(menu).Run(); err != nil {
// 		panic(err)
// 	}
// }

// Old One to Modify
func App(db *sql.DB) {
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

	// Setting Views for each media input
	addBookMedia(db, pages)
	addMovieMedia(db, pages)
	addVidoGameMedia(db, pages)
	addUser(db, pages)
	mediaModal(pages)

	// Main menu prmitive to interact with TUI
	searchInput := tview.NewInputField().SetLabel("Search: ")

	// Start view of the app when user first runs it
	mainTextView := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetWordWrap(true).
		SetText("Search Results")

	// updateMainView := updateMain(db, app, searchRequest, mainTextView, searchInput)

	menu := tview.NewForm().
		AddFormItem(searchInput).
		AddButton("Search", func() {
			updateMain(db, app, searchInput.GetText(), mainTextView, searchInput)
		}).
		AddButton("Add Media", func() {
			pages.SwitchToPage("mediaModal")
		}).
		AddButton("Quit", func() {
			app.Stop()
		})

	// Setting the layout in grid
	sideBar := newPrimitive("Side bar")
	main := mainTextView

	// adding menu, main, and sidebar to grid to draw.
	grid.AddItem(menu, 1, 0, 1, 1, 0, 100, true).
		AddItem(main, 1, 1, 1, 1, 0, 100, false).
		AddItem(sideBar, 1, 2, 1, 1, 0, 100, false)

	pages.AddPage("mainMenu", grid, true, true)

	// Startup of the App
	if err := app.SetRoot(pages, true).SetFocus(menu).Run(); err != nil {
		panic(err)
	}
}

// Original
// func App(db *sql.DB) {
// 	// This lay out is found: https://github.com/rivo/tview/wiki/Grid
// 	// Will be adapting it for my project.

// 	app := tview.NewApplication()
// 	pages := tview.NewPages()

// 	// Setting own things
// 	newPrimitive := func(text string) tview.Primitive {
// 		return tview.NewTextView().
// 			SetTextAlign(tview.AlignCenter).
// 			SetText(text)
// 	}

// 	// Layout of TUI
// 	grid := tview.NewGrid().
// 		SetRows(3, 0, 3).
// 		SetColumns(30, 0, 30).
// 		SetBorders(true).
// 		AddItem(newPrimitive("Header"), 0, 0, 1, 3, 0, 0, false).
// 		AddItem(newPrimitive("Footer"), 2, 0, 1, 3, 0, 0, false)

// 	// Setting Views for each media input
// 	addBookMedia(db, pages)
// 	addMovieMedia(db, pages)
// 	addVidoGameMedia(db, pages)
// 	addUser(db, pages)
// 	mediaModal(pages)

// 	// Main menu prmitive to interact with TUI
// 	searchInput := tview.NewInputField().SetLabel("Search: ")

// 	// This is the main interaction view on the TUI
// 	menuPrimitive := func() (tview.Primitive, []interface{}) {
// 		var searchResult []interface{}
// 		menuPrim := tview.NewForm().
// 			// TODO create function to seach data base
// 			AddFormItem(searchInput).
// 			// Have User Add New Media to DB
// 			AddButton("Add Media", func() {
// 				// pages.SwitchToPage("addBook")
// 				pages.SwitchToPage("mediaModal")

// 			}).
// 			// Exit App
// 			AddButton("Quit", func() {
// 				app.Stop()
// 			}).
// 			// Searching through the DB
// 			AddButton("Search", func() {
// 				//TODO: this may need to be bigger or in it's own function
// 				// IDEA: have the mainPrimitive call from here
// 				// On press call and update the mainPrimitive
// 				search := searchInput.GetText()
// 				if search != "" {
// 					// Temp just for right now
// 					// This will need to get the input and search the database for items
// 					// utils.AddUserInfo(search, db)

// 					searchValue, err := utils.SearchTables(db, search)
// 					if err != nil {
// 						log.Printf("IN APP: %s", err)
// 					}
// 					// Setting the results to new var
// 					searchResult = searchValue

// 					// Loging the values for Debugging
// 					// for i, ptr := range searchResult {
// 					// 	log.Printf("ValuePtrs[%d]: %v", i, *ptr.(*interface{}))

// 					// 	prtValue, ok := (*(searchResult[i].(*interface{}))).(string)
// 					// 	if !ok {
// 					// 		log.Println("Failed to convert valuePtrs[1] to str")
// 					// 	} else {
// 					// 		log.Printf("Converted value: %s", prtValue)
// 					// 	}
// 					// }

// 					// log.Printf("In Submit Button: %s\n", searchResult...)
// 					searchInput.SetText("")
// 				}
// 			})

// 		return menuPrim, searchResult
// 	}

// 	// This should be showing the search results.
// 	// IDEA may need to rethink this more
// 	mainPrimitive := func(searchResult []interface{}) tview.Primitive {
// 		// Creating the text view
// 		mainPrim := tview.NewTextView().
// 			SetTextAlign(tview.AlignCenter).
// 			SetWordWrap(true)

// 		log.Println("IN MAINPRIMITIVE")

// 		results := mainPrimitiveResults(searchResult)

// 		mainPrim.SetText(results)

// 		return mainPrim
// 	}

// 	// Setting the layout in grid
// 	menu, searchResult := menuPrimitive()
// 	main := mainPrimitive(searchResult)
// 	sideBar := newPrimitive("Side bar")

// 	// adding menu, main, and sidebar to grid to draw.
// 	grid.AddItem(menu, 1, 0, 1, 1, 0, 100, true).
// 		AddItem(main, 1, 1, 1, 1, 0, 100, false).
// 		AddItem(sideBar, 1, 2, 1, 1, 0, 100, false)

// 	pages.AddPage("mainMenu", grid, true, true)

// 	// Startup of the App
// 	if err := app.SetRoot(pages, true).SetFocus(menu).Run(); err != nil {
// 		panic(err)
// 	}
// }

func updateMain(db *sql.DB, app *tview.Application, searchRequest string, mainTextView *tview.TextView, searchInput *tview.InputField) {
	go func() {
		searchResults, err := utils.SearchTables(db, searchRequest)
		if err != nil {
			log.Printf("Couldn't Find: %v", err)
			return
		}

		// Converting results to a strings
		resultsText := mainPrimitiveResults(searchResults)

		// Updating UI
		app.QueueUpdateDraw(func() {
			mainTextView.SetText(resultsText)
			searchInput.SetText("") // Clears input fields
		})

	}()
}

func mainPrimitiveResults(searchResult []interface{}) string {
	var results string
	for i, prt := range searchResult {
		log.Printf("ValuePrts[%d]: %v", i, *prt.(*interface{}))

		prtValue, ok := (*(searchResult[i].(*interface{}))).(string)
		if !ok {
			log.Printf("Failed to Convert ValuePtrs[i] to str")
		} else {
			log.Printf("Converted Value: %s", prtValue)
			results = "" + prtValue
			log.Printf("Results: %s", results)
		}
	}

	return results
}

/*
This is for mangaging adding media to db
*/
func mediaModal(pages *tview.Pages) {
	// Pop up for adding media to app
	addMediaModal := tview.NewModal().
		AddButtons([]string{"Add Books", "Add Movie", "Add Game", "Add User", "Back"}).
		SetText("Add Media").
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Add Books" {
				pages.SwitchToPage("addBook")
			} else if buttonLabel == "Add Movie" {
				pages.SwitchToPage("addMovie")
			} else if buttonLabel == "Add Game" {
				pages.SwitchToPage("addGame")
			} else if buttonLabel == "Add User" {
				pages.SwitchToPage("addUser")
			} else if buttonLabel == "Back" {
				pages.SwitchToPage("mainMenu")
			}
		})

	addMediaModal.SetBackgroundColor(tcell.ColorBlack)

	pages.AddPage("mediaModal", addMediaModal, true, false)
}

func addBookMedia(db *sql.DB, pages *tview.Pages) {
	// TODO: have the input fields return user input and input the data into the db

	titleInput := tview.NewInputField().SetLabel("Input Title: ")
	pageNumInput := tview.NewInputField().SetLabel("Input Page Number: ")
	authorInput := tview.NewInputField().SetLabel("Input Author Name: ")

	addMediaForm := tview.NewForm().
		// AddInputField("Input Title: ", "", 0, nil, nil).
		AddFormItem(titleInput).
		AddFormItem(pageNumInput).
		AddFormItem(authorInput).
		AddButton("Submit", func() {
			title := titleInput.GetText()
			pageNum := pageNumInput.GetText()
			author := authorInput.GetText()

			// Check if all fields are filled
			if title != "" && pageNum != "" && author != "" {

				// Convert pageNum from str to int
				pageNumber, err := strconv.Atoi(pageNum)
				if err != nil {
					// TODO add error to page
					fmt.Println(pageNumber, "Not a Valid Page Number")
					return
				}

				utils.AddBookInfo(title, pageNumber, author, db)

			}
		}).
		AddButton("Back", func() {
			// app.SetRoot(grid, true)
			pages.SwitchToPage("mediaModal")
		})

	// This has to be set out side of the tview.NewForm() for the input fields and buttons to show
	addMediaForm.SetBorder(true).SetTitle("Add Book").SetTitleAlign(tview.AlignCenter)

	// Setting the addMediaForm to be shown
	addBookFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(addMediaForm, 0, 1, true)

	pages.AddPage("addBook", addBookFlex, true, false)
}

//-----Do Note that the following function are patterened after addBookMedia()-----//

func addMovieMedia(db *sql.DB, pages *tview.Pages) {
	// TODO: have the input fields return user input and input the data into the db

	titleInput := tview.NewInputField().SetLabel("Input Title Name: ")

	addMovieForm := tview.NewForm().
		AddFormItem(titleInput).
		AddButton("Submit", func() {
			title := titleInput.GetText()
			//TODO have input fields input into db
			if title != "" {
				// fmt.Println("Item Submited")
				utils.AddMovieInfo(title, db)

			}
		}).
		AddButton("Back", func() {
			// app.SetRoot(grid, true)
			pages.SwitchToPage("mediaModal")
		})

	addMovieForm.SetBorder(true).SetTitle("Add Movie").SetTitleAlign(tview.AlignCenter)

	addBookFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(addMovieForm, 0, 1, true)

	pages.AddPage("addMovie", addBookFlex, true, false)
}

func addVidoGameMedia(db *sql.DB, pages *tview.Pages) {
	// TODO: have the input fields return user input and input the data into the db

	titleInput := tview.NewInputField().SetLabel("Input Video Game Title: ")

	addGameForm := tview.NewForm().
		AddFormItem(titleInput).
		AddButton("Submit", func() {
			title := titleInput.GetText()

			//TODO have input fields input into db
			if title != "" {
				// fmt.Println("Item Submited")
				utils.AddVideoGameInfo(title, db)
			}
		}).
		AddButton("Back", func() {
			// app.SetRoot(grid, true)
			pages.SwitchToPage("mediaModal")
		})

	addGameForm.SetBorder(true).SetTitle("Add Movie").SetTitleAlign(tview.AlignCenter)

	addGameFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(addGameForm, 0, 1, true)

	pages.AddPage("addGame", addGameFlex, true, false)
}

func addUser(db *sql.DB, pages *tview.Pages) {
	// TODO: have the input fields return user input and input the data into the db

	userNameInput := tview.NewInputField().SetLabel("Input User Name")

	addUserForm := tview.NewForm().
		AddFormItem(userNameInput).
		AddButton("Submit", func() {
			user := userNameInput.GetText()

			//TODO have input fields input into db
			if user != "" {
				// fmt.Println("Item Submited")
				utils.AddUserInfo(user, db)
			}
		}).
		AddButton("Back", func() {
			// app.SetRoot(grid, true)
			pages.SwitchToPage("mediaModal")
		})

	addUserForm.SetBorder(true).SetTitle("Add Movie").SetTitleAlign(tview.AlignCenter)

	addUserFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(addUserForm, 0, 1, true)

	pages.AddPage("addUser", addUserFlex, true, false)

}
