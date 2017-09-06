package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getProdList() {

	ajax := "https://www.tokopedia.com/ajax/shop/shop.pl"
	user := "rahmataligos"
	id := "138023"
	page := "1"
	max := "80"
	res, err := http.Get(ajax + "?u=/" + user + "/page/" +
		page + "&a=reload_data&action=reload_data&qs=perpage%3D" +
		max + "&uri_path=%2Fibishop%2Fpage%2F" + page + "&s_id=" + id)
	if err != nil {
		log.Fatal(err)
	}
	byt, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	var dat map[string]interface{}

	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}

	var data = dat["content"].(map[string]interface{})
	prodsHTML := "<html><head></head><body>" + data["showcase"].(string) + "</body></html>"
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(prodsHTML))
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".product").Each(func(i int, s *goquery.Selection) {

		prodname := s.Find("div[class=meta-product] b").Text()
		prodlink, _ := s.Find("a").First().Attr("href")
		prodthumb, _ := s.Find(".product-image img").Attr("src")

		fmt.Printf("Product %d:\n * %s\n * %s\n * %s\n", i, prodname, prodlink, prodthumb)

	})

}

func main() {
	getProdList()
}
