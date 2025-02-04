package main

// getting demo writen
// the struct are going to act as the database for right now
// the terminal is going ot act as the UI for now

type book struct {
	title      string
	pageNumber int
	author     string
}

type movie struct {
	title string
}

type videoGame struct {
	title string
}

func main() {
	bookInfo := addBookInfo("testBook", 125, "Me")
	movieInfo := addMovieInfo("testMovie")
	videoGameInfo := addVideoGameInfo("testVideoGame")

	println(bookInfo.title)
	println(movieInfo.title)
	println(videoGameInfo.title)

}

func addBookInfo(title string, pageNumber int, author string) *book {
	b := book{title: title}
	b.pageNumber = pageNumber
	b.author = author

	return &b
}

func addMovieInfo(title string) *movie {
	m := movie{title: title}
	return &m
}

func addVideoGameInfo(title string) *videoGame {
	g := videoGame{title: title}
	return &g
}
