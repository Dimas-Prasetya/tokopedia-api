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

func getShopData(user string) (id string, gold bool) {

	username := strings.ToLower(user)
	url := "https://www.tokopedia.com/" + username
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("h1").Each(func(i int, s *goquery.Selection) {

		_, ok := s.Attr("class")
		gold = ok

	})

	doc.Find(`input[name="shop_id"]`).Each(func(i int, s *goquery.Selection) {

		shopid, _ := s.Attr("value")
		id = shopid

	})

	return

}

func getProdList(user, id, page string) (name, price, link, thumb []string) {

	ajax := "https://www.tokopedia.com/ajax/shop/shop.pl"
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
		prodprice := s.Find(".price").Text()
		prodlink, _ := s.Find("a").First().Attr("href")
		prodthumb, _ := s.Find(".product-image img").Attr("src")
		name = append(name, prodname)
		price = append(price, prodprice)
		link = append(link, prodlink)
		thumb = append(thumb, prodthumb)

	})

	return

}

func getProdDetail(url string) (desc string) {

	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".tab-content").Each(func(i int, s *goquery.Selection) {
		proddesc := s.Find("p[itemprop=description]").Text()
		desc = proddesc
	})

	return

}

func main() {

	user := "idealmuslimshop"
	page := "1"
	id, _ := getShopData(user)
	_, _, link, _ := getProdList(user, id, page)
	desc := getProdDetail(link[0])
	fmt.Println(desc)

}
