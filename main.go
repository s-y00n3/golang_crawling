package main

import (
	"fmt"
	"strings"

	"github.com/labstack/echo"

	scrapper "crawling/scrape"
)

var fileName = "web.csv"

// home.html로 연결
func Handler(c echo.Context) error {
	return c.File("home.html")
}

// 검색 쿼리에 대해 스크래핑 요청
func HandleFunc(c echo.Context) error {
	query := strings.ToLower(scrapper.CleanString(c.FormValue("query")))
	fmt.Println(query)
	scrapper.Scrapper(query)

	return c.Attachment(fileName, query+".csv")
}

func main() {
	e := echo.New()
	e.GET("/", Handler)
	e.POST("/scrape", HandleFunc)
	e.Logger.Fatal(e.Start(":1323"))
}
