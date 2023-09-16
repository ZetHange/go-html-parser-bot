package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type Mark struct {
	Title        string `json:"title"`
	Link         string `json:"link"`
	Image        string `json:"image"`
	Cost         string `json:"cost"`
	Availability bool   `json:"availability"`
}

var OldMarks []*Mark

func ScrapeMyCollection(availability string) ([]*Mark, error) {
	res, err := http.Get("https://mycollection48.ru/novinki/novinki-pochtovye-marki?limit=1000")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var marks []*Mark

	doc.Find(".category_products_array").Each(func(i int, s *goquery.Selection) {
		s.Find(".product-thumb").Each(func(i int, s *goquery.Selection) {
			link, _ := s.Find("a").Attr("href")
			image, _ := s.Find("img").Attr("src")
			findAvailability := s.Find(".product_stock").Text()
			findAvailability = strings.ReplaceAll(findAvailability, "\n", "")
			isAvailable := findAvailability != "нет в наличии"

			mark := Mark{
				Title:        s.Find("a").Text(),
				Link:         link,
				Image:        image,
				Cost:         s.Find(".product_normal_price").Text(),
				Availability: isAvailable,
			}
			if availability == "" {
				marks = append(marks, &mark)
			} else if availability == "true" && isAvailable {
				marks = append(marks, &mark)
			} else if availability == "false" && !isAvailable {
				marks = append(marks, &mark)
			}
		})
	})

	return marks, nil
}

func GetMarks(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		query := r.URL.Query()
		availability := query.Get("availability")
		//var availabilityBool bool
		//
		//if availability == "true" {
		//	availabilityBool = true
		//} else if availability == "false" {
		//	availabilityBool = false
		//}

		collection, err := ScrapeMyCollection(availability)
		if err != nil {
			log.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(collection)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func main() {
	for {
		collection, err := ScrapeMyCollection("")
		if err != nil {
			fmt.Println(err)
		}

		if OldMarks == nil {
			OldMarks = collection
		}

		if !reflect.DeepEqual(collection, OldMarks) {
			if len(collection) >= 5 {
				var Actions []*Action
				Actions = append(Actions, &Action{Url: "https://mycollection48.ru/novinki/novinki-pochtovye-marki",
					Action: "view",
					Label:  "Открыть сайт",
				})
				message := ""
				firstFive := collection[:5]

				for i, mark := range firstFive {
					formattedMark := fmt.Sprintf("%v. %s, цена: %s\n", i+1, mark.Title, mark.Cost)
					message = message + formattedMark
				}

				go SendNotification(&Alert{
					Topic:    "my_collection_48",
					Title:    "Обновление \"Моя коллекция\"",
					Message:  message,
					Markdown: true,
					Tags:     []string{"warning"},
					Priority: 4,
					Actions:  Actions,
				})
			}
		}

		OldMarks = collection

		time.Sleep(30 * time.Second)
	}
	//http.HandleFunc("/", GetMarks)
	//http.ListenAndServe(":8080", nil)

	//for _, mark := range collection {
	//	//fmt.Printf("", mark.Title)
	//	//fmt.Printf("%v %s, стоимость: %s\n", mark.Availability, mark.Title, mark.Cost)
	//}
}
