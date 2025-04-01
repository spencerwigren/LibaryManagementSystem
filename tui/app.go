package tui

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

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
		SetRows(3, 0, 3).      // Setting rows height three rows
		SetColumns(35, 0, 35). // Setting columns with three columns
		SetBorders(true).
		AddItem(newPrimitive("Header"), 0, 0, 1, 3, 0, 0, false).
		AddItem(newPrimitive(footerFormate), 2, 0, 1, 3, 0, 0, false)

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

	// Setting Output for SideBar
	sideBarTextView := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetWordWrap(true).
		SetText("Newest Content")

	// Creating SideBar
	sideBar := tview.NewForm().
		AddFormItem(sideBarTextView)

	// This will get the title to show even if the text is there
	// sideBar.SetBorder(true).SetTitle("Newest Content").SetTitleAlign(tview.AlignCenter)

	sideBar.SetTitle("Newest Content")

	// Setting the layout in grid
	// sideBar := newPrimitive("Side bar")
	main := mainTextView

	// Setting Views for each media input
	addBookMedia(db, pages, sideBarTextView)
	addMovieMedia(db, pages, sideBarTextView)
	addVidoGameMedia(db, pages, sideBarTextView)
	addUser(db, pages)
	mediaModal(pages, db, sideBarTextView)

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
				resultOutput := fmt.Sprintf("Search Results\nTitle: %s\nPage Number: %s\nAuthor: %s\nDateTime: %s", resultsText[1], resultsText[2], resultsText[3], resultsText[4])
				// Dispalying the new result/title to the mainTextView
				mainTextView.SetText(resultOutput)

			} else if len(resultsText) > 1 && (tableName == "movies" || tableName == "videoGames") {
				resultOutput := fmt.Sprintf("Search Results\nTitle: %s\nRating: %s\nRelease Year: %s\nDateTime: %s", resultsText[1], resultsText[2], resultsText[3], resultsText[4])
				mainTextView.SetText(resultOutput)

				// This may be redundent, or not???
			} else if len(resultsText) > 1 { // Movies and VideGames have same fields
				resultOutput := fmt.Sprintf("Search Results\nTitle: %s\nDateTime: %s", resultsText[1], resultsText[2])
				mainTextView.SetText(resultOutput)

			} else {
				log.Printf("Couln't Remove first index of: %v", resultsText)
			}

			searchInput.SetText("") // Clears input fields
		})

	}()
}

func mainPrimitiveResults(searchResult []interface{}) []string {
	var results []string

	for i, prt := range searchResult {
		log.Printf("ValuePrts[%d]: %v", i, *prt.(*interface{}))
		// searchPrt := searchResult[i]

		if strValue, ok := (*(searchResult[i].(*interface{}))).(string); ok {
			log.Printf("Converted Value: %s", strValue)
			// results += " " + strValue
			results = append(results, strValue)
			log.Printf("Results: %s", results)

		} else if intValue, ok := (*(searchResult[i].((*interface{})))).(int64); ok {
			log.Printf("Converted Value: %d", intValue)
			intVal := int(intValue)

			// converting and setting to results)
			results = append(results, strconv.Itoa(intVal))
			log.Printf("Results: %s", results)

		} else if dateValue, ok := (*(searchResult[i].((*interface{})))).(time.Time); ok {
			log.Printf("Converted Value : %v", dateValue)
			convertedTime := dateValue.Format(time.UnixDate)

			results = append(results, convertedTime)

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
func mediaModal(pages *tview.Pages, db *sql.DB, sideBarTextView *tview.TextView) {
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
				mostRecentEntries(db, sideBarTextView)
			}
		})

	addMediaModal.SetBackgroundColor(tcell.ColorBlack)

	pages.AddPage("mediaModal", addMediaModal, true, false)
}

func addBookMedia(db *sql.DB, pages *tview.Pages, sideBarTextView *tview.TextView) {
	// Setting tview input fields
	titleInput := tview.NewInputField().SetLabel("Input Title: ")
	pageNumInput := tview.NewInputField().SetLabel("Input Page Number: ")
	authorInput := tview.NewInputField().SetLabel("Input Author Name: ")

	addMediaForm := tview.NewForm().
		// AddInputField("Input Title: ", "", 0, nil, nil).
		AddFormItem(titleInput).
		AddFormItem(pageNumInput).
		AddFormItem(authorInput).
		AddButton("Submit", func() {
			// Getting user text
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

				// Setting info into db
				utils.AddBookInfo(title, pageNumber, author, db)

				// Clearning the text fields
				titleInput.SetText("")
				pageNumInput.SetText("")
				authorInput.SetText("")

				// TODO IDEA: make this a list of the 5 most recent entries in the db
				// May need to make this it's own function at some point to draw all of them
				// sideBarTextUpdate := fmt.Sprintf("Newest Content\nTitle: %s\n", title)
				// sideBarTextView.SetText(sideBarTextUpdate)

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

//-----Note that the following function are patterened after addBookMedia()-----//

func addMovieMedia(db *sql.DB, pages *tview.Pages, sideBarTextView *tview.TextView) {
	titleInput := tview.NewInputField().SetLabel("Input Title Name: ")
	ratingInput := tview.NewInputField().SetLabel("Input Rating: ")
	releaseYearInput := tview.NewInputField().SetLabel("Input Release Year: ")

	addMovieForm := tview.NewForm().
		AddFormItem(titleInput).
		AddFormItem(ratingInput).
		AddFormItem(releaseYearInput).
		AddButton("Submit", func() {
			title := titleInput.GetText()
			rating := ratingInput.GetText()
			year := releaseYearInput.GetText()

			if title != "" && rating != "" && year != "" {

				yearConverted, err := strconv.Atoi(year)
				if err != nil {
					log.Println("ERROR:", err)
				}
				utils.AddMovieInfo(db, title, rating, yearConverted)

				titleInput.SetText("")
				ratingInput.SetText("")
				releaseYearInput.SetText("")

				// sideBarTextUpdate := fmt.Sprintf("Newest Content\nTitle: %s\n", title)
				// sideBarTextView.SetText(sideBarTextUpdate)

			}
		}).
		AddButton("Back", func() {
			pages.SwitchToPage("mediaModal")
		})

	addMovieForm.SetBorder(true).SetTitle("Add Movie").SetTitleAlign(tview.AlignCenter)

	addBookFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(addMovieForm, 0, 1, true)

	pages.AddPage("addMovie", addBookFlex, true, false)
}

func addVidoGameMedia(db *sql.DB, pages *tview.Pages, sideBarTextView *tview.TextView) {
	titleInput := tview.NewInputField().SetLabel("Input Video Game Title: ")
	ratingInput := tview.NewInputField().SetLabel("Input Rating: ")
	releaseYearInput := tview.NewInputField().SetLabel("Input Release Year: ")

	addGameForm := tview.NewForm().
		AddFormItem(titleInput).
		AddFormItem(ratingInput).
		AddFormItem(releaseYearInput).
		AddButton("Submit", func() {
			title := titleInput.GetText()
			rating := ratingInput.GetText()
			year := releaseYearInput.GetText()

			if title != "" && rating != "" && year != "" {

				yearConv, err := strconv.Atoi(year)
				if err != nil {
					log.Println("ERROR:", err)
				}

				utils.AddVideoGameInfo(db, title, rating, yearConv)

				titleInput.SetText("")
				ratingInput.SetText("")
				releaseYearInput.SetText("")

				// sideBarTextUpdate := fmt.Sprintf("Newest Content\nTitle: %s\n", title)
				// sideBarTextView.SetText(sideBarTextUpdate)
			}
		}).
		AddButton("Back", func() {
			pages.SwitchToPage("mediaModal")
		})

	addGameForm.SetBorder(true).SetTitle("Add Movie").SetTitleAlign(tview.AlignCenter)

	addGameFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(addGameForm, 0, 1, true)

	pages.AddPage("addGame", addGameFlex, true, false)
}

func addUser(db *sql.DB, pages *tview.Pages) {
	userNameInput := tview.NewInputField().SetLabel("Input User Name")

	addUserForm := tview.NewForm().
		AddFormItem(userNameInput).
		AddButton("Submit", func() {
			user := userNameInput.GetText()

			if user != "" {
				utils.AddUserInfo(user, db)

				userNameInput.SetText("")
			}
		}).
		AddButton("Back", func() {
			pages.SwitchToPage("mediaModal")
		})

	addUserForm.SetBorder(true).SetTitle("Add Movie").SetTitleAlign(tview.AlignCenter)

	addUserFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(addUserForm, 0, 1, true)

	pages.AddPage("addUser", addUserFlex, true, false)

}

//---------------------------------------------------------------------------------//

func mostRecentEntries(db *sql.DB, sideBarTextView *tview.TextView) {
	rows := utils.QueryMostRecent(db)

	// checking if no rows
	if rows == nil {
		log.Println("NO Rows Returned")
		return
	}
	defer rows.Close()

	sideBarTitle := "Most Recent Entries"
	output := []string{sideBarTitle}
	var title string
	var time time.Time
	count := 0

	// Getting the data from the rows
	// And saving it to the output
	for rows.Next() {
		if err := rows.Scan(&title, &time); err != nil {
			log.Printf("ERROR %s", err)
			continue
		}
		count++
		output = append(output, fmt.Sprintf("[%d] Title: %s", count, title))
		log.Printf("[%d] Title: %s", count, title)
	}

	// creating one single strings for output
	// Setting it to the TUI view
	outputFinal := strings.Join(output, "\n")
	sideBarTextView.SetText(outputFinal)

	// sideBarTextView.SetText("This is a test")
}
