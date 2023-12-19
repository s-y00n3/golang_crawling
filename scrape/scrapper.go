package scrapper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Web struct {
	title    string
	location string
}

// 스크랩 함수
func Scrapper(query string) {
	var jobs []Web
	var baseURL string = "https://www.jobkorea.co.kr/Search/?stext=" + query
	c1 := make(chan []Web)
	TotalPage := getPages(baseURL)
	fmt.Println("TotalPage : ", TotalPage)

	for i := 0; i < TotalPage; i++ {
		go getCard(i, baseURL, c1)
	}

	for i := 0; i < TotalPage; i++ {
		extractJob := <-c1
		jobs = append(jobs, extractJob...)
	}
	writeJobs(jobs)
	fmt.Println("done")
}

// 크롤링할 페이지 개수 확인
func getPages(baseURL string) int {
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	checkErr(err)

	doc.Find(".tplPagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("li").Length()
	})

	return pages
}

// 특정 페이지 결과 데이터 가져오기
func getCard(page int, baseURL string, c1 chan []Web) {
	var jobs []Web
	c := make(chan Web)
	URL := baseURL + "&Page_No=2" + strconv.Itoa(page*10) // int 타입을 string 으로 변환
	fmt.Println(URL)

	res, err := http.Get(URL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".post-list-info")

	searchCards.Each(func(i int, s *goquery.Selection) {
		go extractJob(s, c)
	})

	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}

	c1 <- jobs
}

// 상세 내용 가져오기
func extractJob(s *goquery.Selection, c chan<- Web) {
	//id, _ := s.Attr("a")
	title := CleanString(s.Find("a").Text())
	location := CleanString("location...")

	c <- Web{
		title:    title,
		location: location}
}

// 결과 저장
func writeJobs(jobs []Web) {
	file, err := os.Create("web.csv")
	checkErr(err)

	w := csv.NewWriter(file)

	defer w.Flush()

	header := []string{"TITLE", "LOCATION"}

	wErr := w.Write(header)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{job.location}
		jobErr := w.Write(jobSlice)
		checkErr(jobErr)
	}
}

// 에러 check
func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// http 요청 결과 상태 코드 확인
func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalf("Status code err: %d %s", res.StatusCode, res.Status)
	}
}

// 문자열 앞뒤 공백, 띄어쓰기 제거
func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}
