package query

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/dollarkillerx/urllib"
)

func TestQuery1(t *testing.T) {
	url := "https://www.nisshinkyo.org/search/area.php?lng=3&area=%E6%9D%B1%E4%BA%AC%E9%83%BD#terms"
	code, rdata, err := urllib.Get(url).Byte()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(code)

	err = os.WriteFile("niss.html", rdata, 00666)
	if err != nil {
		log.Fatalln(err)
	}
}

func TestQuery3(t *testing.T) {
	DownloadInitList("area.php?lng=3&area=東京#terms")
}

func DownloadInitList(addUrl string) {
	//url := fmt.Sprintf("https://www.nisshinkyo.org/search/%s", addUrl)
	//code, rdata, err := urllib.Get(url).Byte()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//fmt.Println(code)
	//err = os.WriteFile("niss_l1.html", rdata, 00666)
	//if err != nil {
	//	log.Fatalln(err)
	//}

	rdata, err := os.ReadFile("niss_l1.html")
	if err != nil {
		log.Fatalln(err)
	}

	reader, err := goquery.NewDocumentFromReader(bytes.NewReader(rdata))
	if err != nil {
		log.Fatalln(err)
	}

	reader.Find(".termsDetail").Find("tr").Each(func(i int, selection *goquery.Selection) {
		if i == 0 {
			return
		}
		th := selection.Find("th").Text()
		td := selection.Find("td").Text()
		thIn := selection.Find("a").AttrOr("href", "")
		url := fmt.Sprintf("https://www.nisshinkyo.org/search/%s", strings.ReplaceAll(thIn, "lng=3", "lng=1"))
		fmt.Println(th, "   =   ", td, "   href:   ", url)
		AsyncInternal(url)
		os.Exit(0)
	})
}

func AsyncInternal(url string) {
	code, rdata, err := urllib.Get(url).Byte()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(code)

	os.WriteFile("niss_l1_l2.html", rdata, 00666)
	reader, err := goquery.NewDocumentFromReader(bytes.NewReader(rdata))
	if err != nil {
		log.Fatalln(err)
	}

	reader.Find()

}

func TestQuery2(t *testing.T) {
	rdata, err := os.ReadFile("niss.html")
	if err != nil {
		log.Fatalln(err)
	}

	reader, err := goquery.NewDocumentFromReader(bytes.NewReader(rdata))
	if err != nil {
		log.Fatalln(err)
	}

	reader.Find("#areajapan").Find("li").Each(func(i int, selection *goquery.Selection) {
		text := selection.Find("a").AttrOr("href", "")
		fmt.Println(text)
	})
}
