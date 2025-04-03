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
		SetText("Most Recent Entries")

	main := mainTextView
	sideBar := sideBarTextView

	// Setting Views for each media input
	addBookFlex := addBookMedia(db, pages)
	addMovieFlex := addMovieMedia(db, pages)
	addGameFlex := addVidoGameMedia(db, pages)
	addUserFlex := addUser(db, pages)
	mediaModal(pages, db, sideBarTextView)

	pages.AddPage("addBook", addBookFlex, true, false)
	pages.AddPage("addGame", addGameFlex, true, false)
	pages.AddPage("addMovie", addMovieFlex, true, false)
	pages.AddPage("addUser", addUserFlex, true, false)

	// Setting Media Update Views
	updateBook := updateBookModal(db, pages)
	updateMovie := updateMovieModal(db, pages)
	updateGame := updateVideoGameModal(db, pages)

	pages.AddPage("addUpdateBook", updateBook, true, false)
	pages.AddPage("addUpdateMovie", updateMovie, true, false)
	pages.AddPage("addUpdateGame", updateGame, true, false)
	// Setting Formate for input errors
	// Default view
	pre := "Book"

	hpn := handlePageNumber(pages, pre)
	hy := handleYear(pages, pre)
	he := handleExisting(pages, pre)
	het := handleExistingTrue(pages, pre)

	pages.AddPage("errModalPageNum", hpn, true, false)
	pages.AddPage("errModalYear", hy, true, false)
	pages.AddPage("errModalExisting", he, true, false)
	pages.AddPage("errModalExistingTrue", het, true, false)

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

func addBookMedia(db *sql.DB, pages *tview.Pages) *tview.Flex {
	// Setting tview input fields
	titleInput := tview.NewInputField().SetLabel("Input Title: ")
	pageNumInput := tview.NewInputField().SetLabel("Input Page Number: ")
	authorInput := tview.NewInputField().SetLabel("Input Author Name: ")

	addMediaForm := tview.NewForm().
		AddFormItem(titleInput).
		AddFormItem(pageNumInput).
		AddFormItem(authorInput).
		AddButton("Submit", func() {
			// Getting user text
			title := titleInput.GetText()
			pageNum := pageNumInput.GetText()
			author := authorInput.GetText()

			// Trimming Any Beginging or End Spaces
			title = strings.TrimSpace(title)
			pageNum = strings.TrimSpace(pageNum)
			author = strings.TrimSpace(author)

			// Check if all fields are filled
			if title != "" && pageNum != "" && author != "" {

				// Convert pageNum from str to int
				pageNumber, err := strconv.Atoi(pageNum)
				if err != nil {
					hpn := handlePageNumber(pages, "Book")          // Setting correct page to return to
					pages.AddPage("errModalYear", hpn, true, false) // Creating that page
					pages.SwitchToPage("errModalPageNum")           // Switching to error page
					pageNumInput.SetText("")                        // Clearing wrong input

					return
					// checking out of bounds
				} else if pageNumber > 100100 || pageNumber <= 0 {
					hpn := handlePageNumber(pages, "Book")          // Setting correct page to return to
					pages.AddPage("errModalYear", hpn, true, false) // Creating that page
					pages.SwitchToPage("errModalPageNum")           // Switching to error page
					pageNumInput.SetText("")                        // Clearing wrong input

					return
				}
				// checking if entry exist
				if utils.CheckExisting(db, title) {
					he := handleExisting(pages, "Book")                // Setting correct page to return to
					pages.AddPage("errModalExisting", he, true, false) //Creating that page
					pages.SwitchToPage("errModalExisting")             // Swithcing to error page
					titleInput.SetText("")                             // Clearing wrong input

					return
				}

				// Setting info into db
				utils.AddBookInfo(title, pageNumber, author, db)

				// Clearning the text fields
				titleInput.SetText("")
				pageNumInput.SetText("")
				authorInput.SetText("")
			}
		}).
		AddButton("Update Entry", func() {
			pages.SwitchToPage("addUpdateBook")
		}).
		AddButton("Back", func() {
			titleInput.SetText("")
			pageNumInput.SetText("")
			authorInput.SetText("")

			pages.SwitchToPage("mediaModal")
		})

	// This has to be set out side of the tview.NewForm() for the input fields and buttons to show
	addMediaForm.SetBorder(true).SetTitle("Add Book").SetTitleAlign(tview.AlignCenter)

	// Setting the addMediaForm to be shown
	addBookFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(addMediaForm, 0, 1, true)

	return addBookFlex
}

//-----Note that the following function are patterened after addBookMedia()-----//

func addMovieMedia(db *sql.DB, pages *tview.Pages) *tview.Flex {
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

			title = strings.TrimSpace(title)
			rating = strings.ToUpper(strings.TrimSpace(rating))
			year = strings.TrimSpace(year)

			if title != "" && rating != "" && year != "" {

				yearConverted, err := strconv.Atoi(year)
				if err != nil {
					log.Println("ERROR:", err)
					hy := handleYear(pages, "Movie")
					pages.AddPage("errModalYear", hy, true, false)
					pages.SwitchToPage("errModalYear")
					releaseYearInput.SetText("")

					return

				} else if yearConverted < 1878 { // Said to be the first year a film was released https://historycooperative.org/first-movie-ever-made/
					hy := handleYear(pages, "Movie")
					pages.AddPage("errModalYear", hy, true, false)
					pages.SwitchToPage("errModalYear")
					releaseYearInput.SetText("")

					return
				}

				if utils.CheckExisting(db, title) {
					he := handleExisting(pages, "Movie")
					pages.AddPage("errModalExisting", he, true, false)
					pages.SwitchToPage("errModalExisting")
					titleInput.SetText("")

					return
				}

				// Checking for correct rating for movies
				ratingList := []string{"G", "PG", "PG-13", "PG13", "R", "NC-17", "NC17"}

				// Quicker look up
				lookup := make(map[string]bool)
				for _, v := range ratingList {
					lookup[v] = true
				}

				// checking if not in
				if !lookup[rating] {
					hr := handleRating(pages, "Movie", ratingList)
					pages.AddAndSwitchToPage("errModalRating", hr, true)
					ratingInput.SetText("")

					return
				}

				utils.AddMovieInfo(db, title, rating, yearConverted)

				titleInput.SetText("")
				ratingInput.SetText("")
				releaseYearInput.SetText("")
			}
		}).
		AddButton("Update Entry", func() {
			pages.SwitchToPage("addUpdateMovie")
		}).
		AddButton("Back", func() {
			titleInput.SetText("")
			ratingInput.SetText("")
			releaseYearInput.SetText("")

			pages.SwitchToPage("mediaModal")
		})

	addMovieForm.SetBorder(true).SetTitle("Add Movie").SetTitleAlign(tview.AlignCenter)

	addMovieFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(addMovieForm, 0, 1, true)

	return addMovieFlex
}

func addVidoGameMedia(db *sql.DB, pages *tview.Pages) *tview.Flex {
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

			title = strings.TrimSpace(title)
			rating = strings.ToUpper(strings.TrimSpace(rating))
			year = strings.TrimSpace(year)

			if title != "" && rating != "" && year != "" {

				yearConv, err := strconv.Atoi(year)
				if err != nil {
					log.Println("ERROR:", err)
					hpn := handleYear(pages, "Game")
					pages.AddPage("errModalYear", hpn, true, false)
					pages.SwitchToPage("errModalYear")
					releaseYearInput.SetText("")

					return
				} else if yearConv <= 1968 { // Said to be the first year a video game was released "Tennis for Two"
					hpn := handleYear(pages, "Game")
					pages.AddPage("errModalYear", hpn, true, false)
					pages.SwitchToPage("errModalYear")
					releaseYearInput.SetText("")

					return
				}

				if utils.CheckExisting(db, title) {
					he := handleExisting(pages, "Game")
					pages.AddPage("errModalExisting", he, true, false)
					pages.SwitchToPage("errModalExisting")
					titleInput.SetText("")

					return
				}

				// Checking for correct rating for video games
				ratingList := []string{"E", "E10+", "T", "M", "AO", "18+", "RP"}

				// Quicker look up
				lookup := make(map[string]bool)
				for _, v := range ratingList {
					lookup[v] = true
				}

				// checking if not in
				if !lookup[rating] {
					hr := handleRating(pages, "Game", ratingList)
					pages.AddAndSwitchToPage("errModalRating", hr, true)
					ratingInput.SetText("")

					return
				}

				utils.AddVideoGameInfo(db, title, rating, yearConv)

				titleInput.SetText("")
				ratingInput.SetText("")
				releaseYearInput.SetText("")
			}
		}).
		AddButton("Update Entry", func() {
			pages.SwitchToPage("addUpdateGame")
		}).
		AddButton("Back", func() {
			titleInput.SetText("")
			ratingInput.SetText("")
			releaseYearInput.SetText("")

			pages.SwitchToPage("mediaModal")
		})

	addGameForm.SetBorder(true).SetTitle("Add Video Game").SetTitleAlign(tview.AlignCenter)

	addGameFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(addGameForm, 0, 1, true)

	return addGameFlex
}

func addUser(db *sql.DB, pages *tview.Pages) *tview.Flex {
	userNameInput := tview.NewInputField().SetLabel("Input User Name")

	addUserForm := tview.NewForm().
		AddFormItem(userNameInput).
		AddButton("Submit", func() {
			user := userNameInput.GetText()

			user = strings.TrimSpace(user)

			if user != "" {
				utils.AddUserInfo(user, db)

				userNameInput.SetText("")
			}
		}).
		AddButton("Back", func() {
			userNameInput.SetText("")

			pages.SwitchToPage("mediaModal")
		})

	addUserForm.SetBorder(true).SetTitle("Add Movie").SetTitleAlign(tview.AlignCenter)

	addUserFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(addUserForm, 0, 1, true)

	return addUserFlex

}

//---------------------------------------------------------------------------------//

func mostRecentEntries(db *sql.DB, sideBarTextView *tview.TextView) {
	rows := utils.QueryMostRecent(db)
	log.Println("IN MOST RECENT ENTREIS")

	// checking if no rows
	if rows == nil {
		log.Println("NO Rows Returned")
		noEnteries := fmt.Sprintln("Most Recent Entries\n NO MEDIA")
		sideBarTextView.SetText(noEnteries)
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
}

func handlePageNumber(pages *tview.Pages, prePage string) *tview.Modal {
	errModal := tview.NewModal().
		SetText("Warning! Not a Valid Page Number").
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {

			pre := fmt.Sprintf("add%s", prePage)

			pages.SwitchToPage(pre)
		})

	return errModal
}

func handleYear(pages *tview.Pages, prePage string) *tview.Modal {
	errModal := tview.NewModal().
		SetText("Warning! Not a Valid Year ").
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pre := fmt.Sprintf("add%s", prePage)
			log.Println("ERROR:", pre)

			pages.SwitchToPage(pre)
		})

	return errModal
}

func handleExisting(pages *tview.Pages, prePage string) *tview.Modal {
	errModal := tview.NewModal().
		SetText("Warning! Entry already exist").
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pre := fmt.Sprintf("add%s", prePage)

			pages.SwitchToPage(pre)
		})

	return errModal
}

func handleRating(pages *tview.Pages, prePage string, ratingList []string) *tview.Modal {
	message := fmt.Sprintf("Warning! Entry is not a Rating\n Rating: %s", ratingList)

	errModal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pre := fmt.Sprintf("add%s", prePage)

			pages.SwitchToPage(pre)
		})

	return errModal
}

func handleExistingTrue(pages *tview.Pages, prePage string) *tview.Modal {
	errModal := tview.NewModal().
		SetText("Warning! Entry Doesn't exist").
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pre := fmt.Sprintf("add%s", prePage)

			pages.SwitchToPage(pre)
		})

	return errModal
}

func updateBookModal(db *sql.DB, pages *tview.Pages) *tview.Flex {
	oringalTitle := tview.NewInputField().SetLabel("Enter Title of Media to Update")
	titleUpdate := tview.NewInputField().SetLabel("Update Title:")
	pageNumUpdate := tview.NewInputField().SetLabel("Update Page Number:")
	authorUpdate := tview.NewInputField().SetLabel("Update Author:")

	updateForm := tview.NewForm().
		AddFormItem(oringalTitle).
		AddFormItem(titleUpdate).
		AddFormItem(pageNumUpdate).
		AddFormItem(authorUpdate).
		AddButton("Update", func() {
			orgTitle := strings.TrimSpace(oringalTitle.GetText())
			titleUp := strings.TrimSpace(titleUpdate.GetText())
			pageNumUp := strings.TrimSpace(pageNumUpdate.GetText())
			authUp := strings.TrimSpace(authorUpdate.GetText())

			log.Println("CHECKING EXISTING", utils.CheckExisting(db, orgTitle))

			if orgTitle != "" && utils.CheckExisting(db, orgTitle) {
				pageNum, err := strconv.Atoi(pageNumUp)
				if err != nil {
					hpn := handlePageNumber(pages, "UpdateBook")    // Setting correct page to return to
					pages.AddPage("errModalYear", hpn, true, false) // Creating that page
					pages.SwitchToPage("errModalPageNum")           // Switching to error page
					pageNumUpdate.SetText("")

					return
					// Clearing wrong input
				} else if pageNum > 100100 || pageNum <= 0 {
					hpn := handlePageNumber(pages, "UpdateBook")    // Setting correct page to return to
					pages.AddPage("errModalYear", hpn, true, false) // Creating that page
					pages.SwitchToPage("errModalPageNum")           // Switching to error page
					pageNumUpdate.SetText("")                       // Clearing wrong input

					return
				}

				utils.UpdateEntryBook(db, titleUp, pageNum, authUp, orgTitle)

				oringalTitle.SetText("")
				titleUpdate.SetText("")
				pageNumUpdate.SetText("")
				authorUpdate.SetText("")

			} else {
				he := handleExistingTrue(pages, "UpdateBook")          // Setting correct page to return to
				pages.AddPage("errModalExistingTrue", he, true, false) //Creating that page
				pages.SwitchToPage("errModalExistingTrue")             // Swithcing to error page
				oringalTitle.SetText("")                               // Clearing wrong input

				return
			}
		}).
		AddButton("Back to Book", func() {
			oringalTitle.SetText("")
			titleUpdate.SetText("")
			pageNumUpdate.SetText("")
			authorUpdate.SetText("")

			pages.SwitchToPage("addBook")
		})

	updateForm.SetBorder(true).SetTitle("Update Book Media").SetTitleAlign(tview.AlignCenter)

	updateBookFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(updateForm, 0, 1, true)

	return updateBookFlex
}

func updateMovieModal(db *sql.DB, pages *tview.Pages) *tview.Flex {
	oringalTitle := tview.NewInputField().SetLabel("Enter Title of Media to Update")
	titleUpdate := tview.NewInputField().SetLabel("Update Title:")
	ratingUpdate := tview.NewInputField().SetLabel("Update Rating:")
	releaseYearUpdate := tview.NewInputField().SetLabel("Update Release Year:")

	updateForm := tview.NewForm().
		AddFormItem(oringalTitle).
		AddFormItem(titleUpdate).
		AddFormItem(ratingUpdate).
		AddFormItem(releaseYearUpdate).
		AddButton("Update", func() {
			orgTitle := strings.TrimSpace(oringalTitle.GetText())
			titleUp := strings.TrimSpace(titleUpdate.GetText())
			ratingUp := strings.ToUpper(strings.TrimSpace(ratingUpdate.GetText()))
			yearUp := strings.TrimSpace(releaseYearUpdate.GetText())

			log.Println("CHECKING EXISTING", utils.CheckExisting(db, orgTitle))

			if orgTitle != "" && utils.CheckExisting(db, orgTitle) {
				yearConverted, err := strconv.Atoi(yearUp)
				if err != nil {
					log.Println("ERROR:", err)
					hy := handleYear(pages, "UpdateMovie")
					pages.AddPage("errModalYear", hy, true, false)
					pages.SwitchToPage("errModalYear")
					releaseYearUpdate.SetText("")

					return

				} else if yearConverted < 1878 { // Said to be the first year a film was released https://historycooperative.org/first-movie-ever-made/
					hy := handleYear(pages, "UpdateMovie")
					pages.AddPage("errModalYear", hy, true, false)
					pages.SwitchToPage("errModalYear")
					releaseYearUpdate.SetText("")

					return
				}

				// Checking for correct rating for movies
				ratingList := []string{"G", "PG", "PG-13", "PG13", "R", "NC-17", "NC17"}

				// Quicker look up
				lookup := make(map[string]bool)
				for _, v := range ratingList {
					lookup[v] = true
				}

				// checking if not in
				if !lookup[ratingUp] {
					hr := handleRating(pages, "UpdateMovie", ratingList)
					pages.AddAndSwitchToPage("errModalRating", hr, true)
					ratingUpdate.SetText("")

					return
				}

				utils.UpdateEntryMovie(db, titleUp, ratingUp, yearConverted, orgTitle)

				oringalTitle.SetText("")
				titleUpdate.SetText("")
				ratingUpdate.SetText("")
				releaseYearUpdate.SetText("")

			} else {
				he := handleExistingTrue(pages, "UpdateMovie")
				pages.AddPage("errModalExistingTrue", he, true, false)
				pages.SwitchToPage("errModalExistingTrue")
				oringalTitle.SetText("")

				return
			}
		}).
		AddButton("Back to Movie", func() {
			oringalTitle.SetText("")
			titleUpdate.SetText("")
			ratingUpdate.SetText("")
			releaseYearUpdate.SetText("")

			pages.SwitchToPage("addMovie")
		})

	updateForm.SetBorder(true).SetTitle("Update Movie Media").SetTitleAlign(tview.AlignCenter)
	updateMovieFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(updateForm, 0, 1, true)

	return updateMovieFlex
}

func updateVideoGameModal(db *sql.DB, pages *tview.Pages) *tview.Flex {
	oringalTitle := tview.NewInputField().SetLabel("Enter Title of Video Game to Update")
	titleUpdate := tview.NewInputField().SetLabel("Update Title:")
	ratingUpdate := tview.NewInputField().SetLabel("Update Rating:")
	releaseYearUpdate := tview.NewInputField().SetLabel("Update Release Year:")

	updateForm := tview.NewForm().
		AddFormItem(oringalTitle).
		AddFormItem(titleUpdate).
		AddFormItem(ratingUpdate).
		AddFormItem(releaseYearUpdate).
		AddButton("Update", func() {
			orgTitle := strings.TrimSpace(oringalTitle.GetText())
			titleUp := strings.TrimSpace(titleUpdate.GetText())
			ratingUp := strings.ToUpper(strings.TrimSpace(ratingUpdate.GetText()))
			yearUp := strings.TrimSpace(releaseYearUpdate.GetText())

			log.Println("CHECKING EXISTING", utils.CheckExisting(db, orgTitle))

			if orgTitle != "" && utils.CheckExisting(db, orgTitle) {
				yearConverted, err := strconv.Atoi(yearUp)
				if err != nil {
					log.Println("ERROR:", err)
					hy := handleYear(pages, "UpdateGame")
					pages.AddPage("errModalYear", hy, true, false)
					pages.SwitchToPage("errModalYear")
					releaseYearUpdate.SetText("")

					return

				} else if yearConverted < 1878 { // Said to be the first year a film was released https://historycooperative.org/first-movie-ever-made/
					hy := handleYear(pages, "UpdateGame")
					pages.AddPage("errModalYear", hy, true, false)
					pages.SwitchToPage("errModalYear")
					releaseYearUpdate.SetText("")

					return
				}

				// Checking for correct rating for movies
				ratingList := []string{"E", "E10+", "T", "M", "AO", "18+", "RP"}

				// Quicker look up
				lookup := make(map[string]bool)
				for _, v := range ratingList {
					lookup[v] = true
				}

				// checking if not in
				if !lookup[ratingUp] {
					hr := handleRating(pages, "UpdateGame", ratingList)
					pages.AddAndSwitchToPage("errModalRating", hr, true)
					ratingUpdate.SetText("")

					return
				}

				utils.UpdateEntryGames(db, titleUp, ratingUp, yearConverted, orgTitle)

				oringalTitle.SetText("")
				titleUpdate.SetText("")
				ratingUpdate.SetText("")
				releaseYearUpdate.SetText("")

			} else {
				he := handleExistingTrue(pages, "UpdateGame")
				pages.AddPage("errModalExistingTrue", he, true, false)
				pages.SwitchToPage("errModalExistingTrue")
				oringalTitle.SetText("")

				return
			}
		}).
		AddButton("Back to Video Game", func() {
			oringalTitle.SetText("")
			titleUpdate.SetText("")
			ratingUpdate.SetText("")
			releaseYearUpdate.SetText("")

			pages.SwitchToPage("addGame")
		})

	updateForm.SetBorder(true).SetTitle("Update Video Game Media").SetTitleAlign(tview.AlignCenter)
	updateMovieFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(updateForm, 0, 1, true)

	return updateMovieFlex

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
