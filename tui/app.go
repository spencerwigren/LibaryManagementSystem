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

	footerFormate := fmt.Sprintln("Commands\nSee all Data: $All\nSearch the name of your input")

	// Layout of TUI
	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0, 30).
		SetBorders(true).
		AddItem(newPrimitive("Header"), 0, 0, 1, 3, 0, 0, false).
		AddItem(newPrimitive(footerFormate), 2, 0, 1, 3, 0, 0, false)

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
			if searchInput.GetText() == "$All" {
				// rowsEnteries := utils.QueryAllEntry(db)
				// log.Printf("Log Entery: %s", rowsEnteries)
				updateMainAll(db, app, mainTextView, searchInput)

			} else {
				updateMain(db, app, searchInput.GetText(), mainTextView, searchInput)
			}
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

func updateMain(db *sql.DB, app *tview.Application, searchRequest string, mainTextView *tview.TextView, searchInput *tview.InputField) {
	// TODO: the display of the search result needs to be better formatted
	go func() {
		searchResults, tableName, err := utils.SearchTables(db, searchRequest)
		if err != nil {
			log.Printf("Couldn't Find: %v", err)
			return
		}

		log.Printf("Table Name in App: %s", tableName)

		// Converting results to a strings
		resultsText := mainPrimitiveResults(searchResults)

		// Updating UI
		app.QueueUpdateDraw(func() {
			if len(resultsText) > 1 && tableName == "books" {
				// result := strings.Join(resultsWords[1:], " ")
				// TODO: FIX this so it can take in all the
				// Also it would be good to put this in a dictonary is possible to output instead of a list.
				resultOutput := fmt.Sprintf("Search Results\nTitle: %s\nPageNumber: %s\nAuthor: %s\n", resultsText[1], resultsText[2], resultsText[3])

				// Dispalying the new result/title to the mainTextView
				mainTextView.SetText(resultOutput)
			} else if len(resultsText) > 1 { // Movies and VideGames have same fields
				resultOutput := fmt.Sprintf("Search Results\nTitle: %s\n", resultsText[1])
				mainTextView.SetText(resultOutput)
			} else {
				log.Printf("Couln't Remove first index of: %v", resultsText)
			}

			searchInput.SetText("") // Clears input fields
		})

	}()
}

func mainPrimitiveResults(searchResult []interface{}) []string {
	// var results string
	var results []string

	for i, prt := range searchResult {
		log.Printf("ValuePrts[%d]: %v", i, *prt.(*interface{}))

		if strValue, ok := (*(searchResult[i].(*interface{}))).(string); ok {
			log.Printf("Converted Value: %s", strValue)
			// results += " " + strValue
			results = append(results, strValue)
			log.Printf("Results: %s", results)

		} else if intValue, ok := (*(searchResult[i].((*interface{})))).(int64); ok {
			log.Printf("Converte Value: %d", intValue)
			intVal := int(intValue)
			// converting and setting to results
			// results += " " + strconv.Itoa(intVal)
			results = append(results, strconv.Itoa(intVal))
			log.Printf("Results: %s", results)

		} else {
			log.Printf("Failed to Convert to string")
		}

	}

	return results
}

/*
For Now this will have to do on the search all
TODO: come back and fix this, so that it works
*/
func updateMainAll(db *sql.DB, app *tview.Application, mainTextView *tview.TextView, searchInput *tview.InputField) {

	go func() {
		rowsEnteries := utils.QueryAllEntry(db)
		log.Printf("Row Enteries: %s", rowsEnteries...)

	}()
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

				bookType := "book"

				utils.AddBookInfo(title, pageNumber, author, bookType, db)

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
