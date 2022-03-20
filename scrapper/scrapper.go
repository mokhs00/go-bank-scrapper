package scrapper

import (
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Job struct {
	id       string
	location string
	title    string
	summary  string
	company  string
}

var limit = 50
var baseURL = "https://kr.indeed.com/jobs"

// Scrape Indeed By term
func Scrape(term string) {
	var jobs []Job

	baseURL += "?q=" + term + "&limit=" + strconv.Itoa(limit)

	//totalPages := getPages(baseURL)
	totalPages := 2

	jobChannel := make(chan []Job)

	for i := 0; i < totalPages; i++ {
		go getJobsFromPage(baseURL, i, jobChannel)
	}

	for i := 0; i < totalPages; i++ {
		extractedJobs := <-jobChannel
		jobs = append(jobs, extractedJobs...)

	}

	writeJobs(jobs)
	fmt.Println("Done :", len(jobs))
}

func writeJobs(jobs []Job) {
	file, err := os.Create("jobs.csv")
	checkErr(err)
	writer := csv.NewWriter(file)

	defer file.Close()
	defer writer.Flush()

	headers := []string{
		"id",
		"location",
		"title",
		"summary",
		"company",
	}

	writer.Write(headers)

	for _, job := range jobs {
		records := []string{job.id, job.location, job.title, job.summary, job.company}
		writer.Write(records)
	}

}

func getJobsFromPage(baseURL string, page int, jobChannel chan<- []Job) {
	var jobs []Job

	pageUrl := baseURL + "&start=" + strconv.Itoa(page*limit)

	res, err := http.Get(pageUrl)

	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	reader, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	c := make(chan Job)

	find := reader.Find(".mosaic-provider-jobcards").Find("a")
	find.Each(func(i int, item *goquery.Selection) {
		go extractJob(item, c)
	})

	for i := 0; i < find.Length(); i++ {
		extractedJob := <-c

		if extractedJob.title != "" {
			jobs = append(jobs, extractedJob)
		}
	}

	jobChannel <- jobs
}

func extractJob(item *goquery.Selection, c chan<- Job) {
	id, exists := item.Attr("data-jk")

	if !exists {
		c <- Job{}
	}

	title := CleanString(item.Find(".jobTitle>span").Text())
	company := CleanString(item.Find(".companyName").Text())
	location := CleanString(item.Find(".companyLocation").Text())
	summary := CleanString(item.Find(".job-snippet").Text())

	c <- Job{
		id:       id,
		title:    title,
		company:  company,
		location: location,
		summary:  summary}

}

func getPages(baseURL string) int {
	pages := 0
	res, err := http.Get(baseURL)
	fmt.Println(baseURL)

	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	reader, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	reader.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	return pages
}

// CleanString clean string
func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")

}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(response *http.Response) {
	if response.StatusCode != 200 {
		log.Fatalln("request failed with status : ", response.StatusCode)
	}
}
