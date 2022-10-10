package query

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dollarkillerx/analysis_japanese_schools/internal/pkg/models"
	"github.com/dollarkillerx/analysis_japanese_schools/utils"
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
		switch i {
		case 0:
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
										school.NumberOfTeachers = uint16(numbers[0])
									}
									if len(numbers) >= 2 {
										school.NumberOfTeachersFull = uint16(numbers[1])
									}
								}
							})
						case strings.Contains(selection.Text(), "設置者種別"):
							selection.Find("td").Each(func(i int, selection *goquery.Selection) {
								switch i {
								case 1:
									if strings.Contains(strings.TrimSpace(selection.Text()), "学校法人") {
										school.SchoolType = models.SchoolCorporation
									} else {
										school.SchoolType = models.FinancialCorporation
									}
								case 3:
									tr := strings.TrimSpace(selection.Text())
									i2 := strings.Split(tr, "～")
									if len(i2) >= 1 {
										parse, err := time.Parse("2006年1月1日", strings.TrimSpace(i2[0]))
										if err != nil {
											log.Println(err)
											return
										}
										school.CertificationPeriodStart = parse
									}
								}
							})
						case strings.Contains(selection.Text(), "校長名"):
							selection.Find("td").Each(func(i int, selection *goquery.Selection) {
								switch i {
								case 1:
									school.SchoolMaster = strings.TrimSpace(selection.Text())
								case 3:
									school.Quota = strings.TrimSpace(selection.Text())
								}
							})
						}
					})
				}
			})
		case 1:
			selection.Find("td").Each(func(i int, selection *goquery.Selection) {
				numbers := utils.ExtractNumbers(utils.WhitespaceOptimization(selection.Text()))
				if len(numbers) <= 0 {
					return
				}

				uim := uint16(numbers[0])

				switch {
				case strings.Contains(selection.Text(), "中国"):
					school.ChinaPeople = uim
				case strings.Contains(selection.Text(), "韓国"):
					school.KoreaPeople = uim
				case strings.Contains(selection.Text(), "台湾"):
					school.TaiwanPeople = uim
				case strings.Contains(selection.Text(), "ﾍﾞﾄﾅﾑ"): // 越南
					school.VietnamPeople = uim
				case strings.Contains(selection.Text(), "ﾈﾊﾟｰﾙ"): // 尼泊尔
					school.NepalPeople = uim
				case strings.Contains(selection.Text(), "ﾀｲ"): // 泰国
					school.ThailandPeople = uim
				case strings.Contains(selection.Text(), "ﾐｬﾝﾏｰ"): // 缅甸
					school.MyanmarPeople = uim
				case strings.Contains(selection.Text(), "ﾓﾝｺﾞﾙ"): // 蒙古
					school.MongoliaPeople = uim
				case strings.Contains(selection.Text(), "ｲﾝﾄﾞﾈｼｱ"): // 印度尼西亚
					school.IndonesiaPeople = uim
				case strings.Contains(selection.Text(), "ｽﾘﾗﾝｶ"): // 斯里兰卡
					school.SriLankaPeople = uim
				case strings.Contains(selection.Text(), "ｽｳｪｰﾃﾞﾝ"): // 瑞典
					school.SwedenPeople = uim
				case strings.Contains(selection.Text(), "ﾏﾚｰｼｱ"): // 马来西亚
					school.MalaysiaPeople = uim
				case strings.Contains(selection.Text(), "ｱﾒﾘｶ"): // 美国
					school.AmericaPeople = uim
				case strings.Contains(selection.Text(), "ｲﾝﾄﾞ"): // 印度
					school.IndonesiaPeople = uim
				case strings.Contains(selection.Text(), "ﾌﾗﾝｽ"): // 法国
					school.FrancePeople = uim
				case strings.Contains(selection.Text(), "ﾛｼｱ"): // 俄罗斯
					school.RussiaPeople = uim
				case strings.Contains(selection.Text(), "ﾌｨﾘﾋﾟﾝ"): // 菲律宾
					school.PhilippinesPeople = uim
				case strings.Contains(selection.Text(), "ｻｳｼﾞｱﾗﾋﾞｱ"): // 沙特阿拉伯
					school.SaudiArabiaPeople = uim
				case strings.Contains(selection.Text(), "ｲﾀﾘｱ"): // 意大利
					school.ItalyPeople = uim
				case strings.Contains(selection.Text(), "ｽﾍﾟｲﾝ"): // 西班牙
					school.SpainPeople = uim
				case strings.Contains(selection.Text(), "ｲｷﾞﾘｽ"): // 英国
					school.EnglandPeople = uim
				case strings.Contains(selection.Text(), "ｶﾅﾀﾞ"): // 加拿大
					school.CanadaPeople = uim
				case strings.Contains(selection.Text(), "ﾊﾞﾝｸﾞﾗﾃﾞｼｭ"): // 孟加拉国
					school.BangladeshPeople = uim
				case strings.Contains(selection.Text(), "ｶﾝﾎﾞｼﾞｱ"): // 柬埔寨
					school.CambodiaPeople = uim
				case strings.Contains(selection.Text(), "ｼﾝｶﾞﾎﾟｰﾙ"): // 新加坡
					school.SingaporePeople = uim
				case strings.Contains(selection.Text(), "ｽｲｽ"): // 瑞士
					school.SwitzerlandPeople = uim
				case strings.Contains(selection.Text(), "ﾄﾞｲﾂ"): // 德国
					school.GermanyPeople = uim
				case strings.Contains(selection.Text(), "ｵｰｽﾄﾗﾘｱ"): // 澳大利亚
					//case strings.Contains(selection.Text(), "ｵｰｽﾄﾗﾘｱ"): // 澳大利亚
					school.AustraliaPeople = uim
				case strings.Contains(selection.Text(), "その他"): // 其他
					school.OtherPeople = uim
				case strings.Contains(selection.Text(), "合計"): // 合計
					school.TotalPeople = uim
				}
			})
		case 2:
			var list [][]string
			selection.Find("tr").Each(func(i int, selection *goquery.Selection) {
				var itemList []string
				if i > 1 {
					selection.Find("td").Each(func(i int, selection *goquery.Selection) {
						text := strings.TrimSpace(selection.Text())
						if text == "" {
							return
						}
						itemList = append(itemList, text)
					})
				}

				if itemList != nil {
					tc := itemList[len(itemList)-1]
					numbers := utils.ExtractNumbers(tc)
					if len(numbers) > 0 {
						list = append(list, itemList)
					}
				}
			})

			marshal, err := json.Marshal(list)
			if err != nil {
				log.Println(err)
			}
			school.CourseInfoJson = string(marshal)
		case 3:
		case 4:
			selection.Find("tr").Each(func(i int, selection *goquery.Selection) {
				switch i {
				case 1:
					selection.Find("td").Each(func(i int, selection *goquery.Selection) {
						numbers := utils.ExtractNumbers(selection.Text())
						if len(numbers) > 0 {
							nu := numbers[0]
							switch i {
							case 0:
								school.EnterPostgraduate = uint16(nu)
							case 1:
								school.EnterUniversity = uint16(nu)
							case 2:
								school.EnterJuniorCollege = uint16(nu)
							case 3:
								school.EnterCollegeTechnology = uint16(nu)
							case 4:
								school.EnterVocationalSchool = uint16(nu)
							case 5:
								school.EnterVariousSchools = uint16(nu)
							case 6:
								school.EnterOtherSchools = uint16(nu)
							}
						}
					})
				case 2:
					school.University = strings.TrimSpace(selection.Find("td").Text())
				}
			})
		}
	})

	school.GoogleMap = reader.Find("iframe").AttrOr("src", "")

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

func TestPbx(t *testing.T) {
	px := "area.php?lng=3&area=岩手#terms"
	split := strings.Split(px, "=")
	fmt.Println(split)
	city := strings.ReplaceAll(split[2], "#terms", "")
	fmt.Println(city)
}
