package main

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

type extractedJob struct {
	id       string
	title    string
	location string
	salary   string

	summary string
}

var baseURL string = "https://kr.indeed.com/jobs?q=python&limit=50"

func main() {

	c := make(chan []extractedJob)

	pages := getPages()

	var jobs []extractedJob

	for i := 0; i < pages; i++ {
		go getPage(i, c)

		// extractedJobs :=

		//jobs = append(jobs, extractedJobs...)
	}

	for i := 0; i < pages; i++ {
		jobs = append(jobs, <-c...)
	}

	writeJobs(jobs)

	fmt.Println("Done: ", len(jobs))

}

func getPage(page int, mainC chan<- []extractedJob) {
	var jobs []extractedJob

	c := make(chan extractedJob)

	pageUrl := baseURL + "$start=" + strconv.Itoa(page*50)

	fmt.Println("Requesting", pageUrl)

	res, err := http.Get(pageUrl)

	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".tapItem")

	searchCards.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, c)

		//jobs = append(jobs, job)
	})
	for i := 0; i < searchCards.Length(); i++ {
		job := <-c

		jobs = append(jobs, job)
	}

	mainC <- jobs
}

func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	id, _ := card.Attr("data-jk")

	title := cleanString(card.Find(".heading4>h2").Text())

	location := cleanString(card.Find(".companyLocation").Text())

	salary := cleanString(card.Find(".salary-snippet>span").Text())

	summary := cleanString(card.Find(".job-snippet").Text())

	c <- extractedJob{id: id, title: title, location: location, summary: summary, salary: salary}
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func writeJobs(jobs []extractedJob) {

	file, err := os.Create("jobs.csv")

	checkErr(err)

	w := csv.NewWriter(file)

	defer w.Flush()

	headers := []string{"ID", "Title", "Location", "Salary", "Summary"}

	wErr := w.Write(headers)

	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{"https://kr.indeed.com/viewjob?jk=" + job.id, job.title, job.location, job.salary, job.summary}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}

}

func getPages() int {
	pages := 0
	res, err := http.Get(baseURL)

	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	checkErr(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	fmt.Println(doc)

	return pages

}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed", res.StatusCode)
	}
}
