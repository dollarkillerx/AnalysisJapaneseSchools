package query

import (
	"bytes"
	"fmt"
	"github.com/dollarkillerx/analysis_japanese_schools/internal/pkg/models"
	"github.com/dollarkillerx/analysis_japanese_schools/utils"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
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

func TestT2(t *testing.T) {
	AsyncInternal("")
}

func AsyncInternal(url string) {
	school := models.LanguageSchool{}

	//code, rdata, err := urllib.Get(url).Byte()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//fmt.Println(code)
	//
	//os.WriteFile("niss_l1_l2.html", rdata, 00666)

	rdata, err := os.ReadFile("niss_l1_l2.html")
	if err != nil {
		log.Fatalln(err)
	}

	reader, err := goquery.NewDocumentFromReader(bytes.NewReader(rdata))
	if err != nil {
		log.Fatalln(err)
	}

	schoolName := reader.Find(".collegeTitle").Text()
	schoolName = strings.ReplaceAll(schoolName, reader.Find(".collegeTitle").Find("span").Text(), "")
	school.Name = schoolName

	bg1 := reader.Find(".floatBox,.clearfix").Find(".floL").Text()
	split := strings.Split(bg1, "\n")
	for _, v := range split {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		index := strings.Index(v, "番号：")
		if index != -1 {
			school.NissID = v[len("番号：")+index:]
		} else {
			parse, err := time.Parse("2006年1月更新", v)
			if err != nil {
				log.Println(err)
				continue
			}

			school.UpDataTime = parse
		}
	}

	reader.Find(".tableStyle04").Each(func(i int, selection *goquery.Selection) {
		selection.Find("table").Each(func(i int, selection *goquery.Selection) {
			switch i {
			case 0:
				selection.Find("tr").Each(func(i int, selection *goquery.Selection) {
					switch {
					case strings.Contains(selection.Text(), "電話番号"):
						selection.Find("td").Each(func(i int, selection *goquery.Selection) {
							//fmt.Println(fmt.Sprintf("%d - %s \n", i, selection.Text()))
							switch i {
							case 1:
								school.Tel = selection.Text()
							case 2:
								tmp := strings.TrimSpace(utils.WhitespaceOptimization(selection.Text()))
								if strings.Contains(tmp, "最寄駅か") {
									jr := strings.Split(tmp, "\n")
									if len(jr) == 2 {
										school.NearestStation = strings.TrimSpace(jr[1])
									}
									//log.Println(school.NearestStation)
								}
							}
						})
					case strings.Contains(selection.Text(), "URL"):
						selection.Find("td").Each(func(i int, selection *goquery.Selection) {
							if i == 1 {
								school.Website = strings.TrimSpace(selection.Text())
							}
						})
					case strings.Contains(selection.Text(), "E-Mail"):
						selection.Find("td").Each(func(i int, selection *goquery.Selection) {
							if i == 1 {
								school.Email = strings.TrimSpace(selection.Text())
							}
						})
					}
				})
			case 1:
				selection.Find("tr").Each(func(i int, selection *goquery.Selection) {
					fmt.Println(i)
					fmt.Println(selection.Text())
					fmt.Println("===============")

					switch {
					case strings.Contains(selection.Text(), "設置者名"):
						selection.Find("td").Each(func(i int, selection *goquery.Selection) {
							switch i {
							case 1:
								school.Founder = strings.TrimSpace(selection.Text())
							case 3:
								parse, err := time.Parse("2006年1月2日", strings.TrimSpace(selection.Text()))
								if err != nil {
									log.Println(err)
									return
								}
								school.TeachingStartTime = parse
							}
						})
					case strings.Contains(selection.Text(), "代表者"):
						selection.Find("td").Each(func(i int, selection *goquery.Selection) {
							//fmt.Printf("%d - %s \n", i, selection.Text())
							switch i {
							case 1:
								school.Representative = strings.TrimSpace(selection.Text())
							case 3:
								numbers := utils.ExtractNumbers(strings.TrimSpace(selection.Text()))
								if len(numbers) >= 1 {

								}
								if len(numbers) >= 2 {

								}
							}
						})
					}
				})
			}
		})

		os.Exit(0)
	})
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
