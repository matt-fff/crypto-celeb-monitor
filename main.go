package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gocolly/colly"
)

func getPreviousCelebs(fName string) map[string][]string {

	previousRows := map[string][]string{}

	file, err := os.Open(fName)
	if err != nil {
		log.Fatalf("Cannot open file %q: %s\n", fName, err)
		return previousRows
	}
	defer file.Close()

	// Create a new reader.
	r := csv.NewReader(bufio.NewReader(file))
	for {
		record, err := r.Read()
		// Stop at EOF.
		if err == io.EOF {
			break
		}

		// I know, we shouldn't have magic indices.
		// I know, we're including the header here. Doesn't really matter.
		previousRows[record[0]] = record
	}

	return previousRows
}

func main() {

	nameRegex, _ := regexp.Compile("[a-zA-Z ]+[a-zA-Z]")
	priceRegex, _ := regexp.Compile("(Price: |Owner:.*)")
	txRegex, _ := regexp.Compile("[0-9]+")
	ownerRegex, _ := regexp.Compile("Owner: ")

	fName := "cryptocelebs.csv"

	previousCelebs := getPreviousCelebs(fName)

	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{"Name", "Price (ETH)", "Transactions", "Owner"})

	// Instantiate default collector
	c := colly.NewCollector()

	newCelebs := []string{}

	c.OnHTML(".item", func(e *colly.HTMLElement) {
		name := nameRegex.FindString(e.ChildText(".name"))
		price := priceRegex.ReplaceAllString(e.ChildText("#left"), "")
		transactions := txRegex.FindString(e.ChildText("#right"))
		owner := ownerRegex.ReplaceAllString(e.ChildText(".element-owner"), "")

		writer.Write([]string{
			name,
			price,
			transactions,
			owner,
		})

		if _, exists := previousCelebs[name]; !exists {
			newCelebs = append(newCelebs, name)
		}
	})

	c.Visit("https://cryptocelebrities.co/marketplace/?sort=newest")

	log.Printf("Scraping finished, check file %q for results\n", fName)

	if len(newCelebs) > 0 {
		reportNewCelebs(newCelebs)
	}
}

func reportNewCelebs(newCelebs []string) {
	msg := fmt.Sprintf("New entries detected: %v", newCelebs)
	log.Println(msg)

	alertToken := os.Getenv("alertToken")
	alertUser := os.Getenv("alertUser")

	// Bail if the API details aren't set
	if alertToken == "" || alertUser == "" {
		log.Println("Missing API details. Aborting push notification.")
		return
	}

	url := "https://api.pushover.net/1/messages.json"

	var jsonStr = []byte(
		fmt.Sprintf(`{
			"token":"%s",
			"user":"%s",
			"message":"%s",
			"url":"https://cryptocelebrities.co/marketplace/?sort=newest"
		}`, alertToken, alertUser, msg))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
