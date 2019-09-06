package main

import (
	"os"
	"net/http"
	"sync"
	"strconv"
	"log"
	"fmt"
	"strings"
	"regexp"
	"github.com/PuerkitoBio/goquery"
)

type Result struct {
	rank int
	content string 
}

var resultsCh = make(chan Result, 10)

// parse web page through goquery
func parseUrl(url string, results []Result) []Result{
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	r := regexp.MustCompile(`[0-9]+`)
	doc.Find(".grid_view .item").Each(func(_ int, s *goquery.Selection) {
			rankStr := s.Find(".pic em").Text()
			multiTitle := s.Find(".info .hd a .title").Text()
			titleSplit := strings.Split(multiTitle,"/")
			title := titleSplit[0]
			star := s.Find(".info .bd .star .rating_num").Text()
			multiSpan := s.Find(".info .bd .star span").Text()
			personsSplit := strings.Split(multiSpan,star)
			personsMix := personsSplit[1]
			persons := r.FindString(personsMix)
			rank, _ := strconv.Atoi(rankStr)
			result := Result{rank,title+" "+star+" "+persons}
			results = append(results,result)
	})
	return results
}

// spider
func spider(i int, wg *sync.WaitGroup) {
	start := i * 25
	var url string = "https://movie.douban.com/top250?start="+strconv.Itoa(start)+"&filter="
	var parseResults []Result
	returnPR := parseUrl(url,parseResults)
	for _, result := range returnPR {
		resultsCh <- result
	}
	wg.Done()
}

// allocate tasks
func allocate() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go spider(i, &wg)
	}
	wg.Wait()
	close(resultsCh)
}

// write resultMap by channel
func writeResultMap(resultMap map[int]string, done chan bool) {
	for result := range resultsCh {
		resultMap[result.rank] = result.content
	}
	done <- true
}


func main() {
	resultMap := make(map[int]string)
	done := make(chan bool)

	file, err := os.Create("./douban_top250.txt")
	if err != nil {
		log.Println(err)
		return 
	}

	go allocate()
	go writeResultMap(resultMap,done)
	<- done

	for i := 0; i <= len(resultMap); i++ {
		if 0 == i {
			fmt.Fprintln(file,"rank  title  star  evaluations")
		} else {
			fmt.Fprintln(file,strconv.Itoa(i)+" "+resultMap[i])
		}
	}
	file.Close()
}
