package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dollarkillerx/analysis_japanese_schools/internal/pkg/models"
	"github.com/dollarkillerx/analysis_japanese_schools/utils"
	"github.com/dollarkillerx/urllib"
	"github.com/google/uuid"
)

type LanguageSchoolUpdate struct{}

func (l *LanguageSchoolUpdate) Update() ([]models.LanguageSchool, error) {
	l1, err := l.genL1()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	l2, err := l.genL2(l1)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return l2, nil
}

func (l *LanguageSchoolUpdate) genL1() ([]string, error) {
	url := "https://www.nisshinkyo.org/search/area.php?lng=3&area=%E6%9D%B1%E4%BA%AC%E9%83%BD#terms"
	code, rdata, err := urllib.Get(url).SetTimeout(time.Second * 10).Byte()
	if err != nil {
		return nil, err
	}

	if code != 200 {
		return nil, errors.New(string(rdata))
	}

	reader, err := goquery.NewDocumentFromReader(bytes.NewReader(rdata))
	if err != nil {
		return nil, err
	}

	var result []string

	reader.Find("#areajapan").Find("li").Each(func(i int, selection *goquery.Selection) {
		text := strings.TrimSpace(selection.Find("a").AttrOr("href", ""))
		if text == "" {
			return
		}

		result = append(result, text)
	})

	return result, nil
}

func (l *LanguageSchoolUpdate) genL2(keys []string) ([]models.LanguageSchool, error) {
	var tasks []models.LanguageSchool
	for _, v := range keys {
		item, err := l.genL2Item(v)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		if len(item) == 0 {
			continue
		}

		tasks = append(tasks, item...)
	}

	for idx, v := range tasks {
		code, rdata, err := urllib.Get(v.NissUrl).Byte()
		if err != nil {
			return nil, err
		}

		if code != 200 {
			return nil, errors.New(string(rdata))
		}

		reader, err := goquery.NewDocumentFromReader(bytes.NewReader(rdata))
		if err != nil {
			return nil, err
		}

		schoolName := reader.Find(".collegeTitle").Text()
		schoolName = strings.ReplaceAll(schoolName, reader.Find(".collegeTitle").Find("span").Text(), "")
		tasks[idx].Name = schoolName

		bg1 := reader.Find(".floatBox,.clearfix").Find(".floL").Text()
		split := strings.Split(bg1, "\n")

		for _, vc := range split {
			vc = strings.TrimSpace(vc)
			if vc == "" {
				continue
			}
			index := strings.Index(vc, "?????????")
			if index != -1 {
				tasks[idx].NissID = vc[len("?????????")+index:]
			} else {
				parse, err := time.Parse("2006???1?????????", vc)
				if err != nil {
					log.Println(err)
					continue
				}

				tasks[idx].UpDataTime = parse
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
							case strings.Contains(selection.Text(), "????????????"):
								selection.Find("td").Each(func(i int, selection *goquery.Selection) {
									switch i {
									case 1:
										tasks[idx].Tel = selection.Text()
									case 2:
										tmp := strings.TrimSpace(utils.WhitespaceOptimization(selection.Text()))
										if strings.Contains(tmp, "????????????") {
											jr := strings.Split(tmp, "\n")
											if len(jr) == 2 {
												tasks[idx].NearestStation = strings.TrimSpace(jr[1])
											}
											//log.Println(tasks[idx].NearestStation)
										}
									}
								})
							case strings.Contains(selection.Text(), "URL"):
								selection.Find("td").Each(func(i int, selection *goquery.Selection) {
									if i == 1 {
										tasks[idx].Website = strings.TrimSpace(selection.Text())
									}
								})
							case strings.Contains(selection.Text(), "E-Mail"):
								selection.Find("td").Each(func(i int, selection *goquery.Selection) {
									if i == 1 {
										tasks[idx].Email = strings.TrimSpace(selection.Text())
									}
								})
							}
						})
					case 1:
						selection.Find("tr").Each(func(i int, selection *goquery.Selection) {
							switch {
							case strings.Contains(selection.Text(), "????????????"):
								selection.Find("td").Each(func(i int, selection *goquery.Selection) {
									switch i {
									case 1:
										tasks[idx].Founder = strings.TrimSpace(selection.Text())
									case 3:
										parse, err := time.Parse("2006???1???2???", strings.TrimSpace(selection.Text()))
										if err != nil {
											log.Println(err)
											return
										}
										tasks[idx].TeachingStartTime = parse
									}
								})
							case strings.Contains(selection.Text(), "?????????"):
								selection.Find("td").Each(func(i int, selection *goquery.Selection) {
									//fmt.Printf("%d - %s \n", i, selection.Text())
									switch i {
									case 1:
										tasks[idx].Representative = strings.TrimSpace(selection.Text())
									case 3:
										numbers := utils.ExtractNumbers(strings.TrimSpace(selection.Text()))
										if len(numbers) >= 1 {
											tasks[idx].NumberOfTeachers = uint16(numbers[0])
										}
										if len(numbers) >= 2 {
											tasks[idx].NumberOfTeachersFull = uint16(numbers[1])
										}
									}
								})
							case strings.Contains(selection.Text(), "???????????????"):
								selection.Find("td").Each(func(i int, selection *goquery.Selection) {
									switch i {
									case 1:
										if strings.Contains(strings.TrimSpace(selection.Text()), "????????????") {
											tasks[idx].SchoolType = models.SchoolCorporation
										} else {
											tasks[idx].SchoolType = models.FinancialCorporation
										}
									case 3:
										tr := strings.TrimSpace(selection.Text())
										i2 := strings.Split(tr, "???")
										if len(i2) >= 1 {
											parse, err := time.Parse("2006???1???1???", strings.TrimSpace(i2[0]))
											if err != nil {
												log.Println(err)
												return
											}
											tasks[idx].CertificationPeriodStart = parse
										}
									}
								})
							case strings.Contains(selection.Text(), "?????????"):
								selection.Find("td").Each(func(i int, selection *goquery.Selection) {
									switch i {
									case 1:
										tasks[idx].SchoolMaster = strings.TrimSpace(selection.Text())
									case 3:
										tasks[idx].Quota = strings.TrimSpace(selection.Text())
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
					case strings.Contains(selection.Text(), "??????"):
						tasks[idx].ChinaPeople = uim
					case strings.Contains(selection.Text(), "??????"):
						tasks[idx].KoreaPeople = uim
					case strings.Contains(selection.Text(), "??????"):
						tasks[idx].TaiwanPeople = uim
					case strings.Contains(selection.Text(), "???????????????"): // ??????
						tasks[idx].VietnamPeople = uim
					case strings.Contains(selection.Text(), "???????????????"): // ?????????
						tasks[idx].NepalPeople = uim
					case strings.Contains(selection.Text(), "??????"): // ??????
						tasks[idx].ThailandPeople = uim
					case strings.Contains(selection.Text(), "???????????????"): // ??????
						tasks[idx].MyanmarPeople = uim
					case strings.Contains(selection.Text(), "???????????????"): // ??????
						tasks[idx].MongoliaPeople = uim
					case strings.Contains(selection.Text(), "?????????????????????"): // ???????????????
						tasks[idx].IndonesiaPeople = uim
					case strings.Contains(selection.Text(), "???????????????"): // ????????????
						tasks[idx].SriLankaPeople = uim
					case strings.Contains(selection.Text(), "?????????????????????"): // ??????
						tasks[idx].SwedenPeople = uim
					case strings.Contains(selection.Text(), "???????????????"): // ????????????
						tasks[idx].MalaysiaPeople = uim
					case strings.Contains(selection.Text(), "????????????"): // ??????
						tasks[idx].AmericaPeople = uim
					case strings.Contains(selection.Text(), "????????????"): // ??????
						tasks[idx].IndonesiaPeople = uim
					case strings.Contains(selection.Text(), "????????????"): // ??????
						tasks[idx].FrancePeople = uim
					case strings.Contains(selection.Text(), "?????????"): // ?????????
						tasks[idx].RussiaPeople = uim
					case strings.Contains(selection.Text(), "??????????????????"): // ?????????
						tasks[idx].PhilippinesPeople = uim
					case strings.Contains(selection.Text(), "???????????????????????????"): // ???????????????
						tasks[idx].SaudiArabiaPeople = uim
					case strings.Contains(selection.Text(), "????????????"): // ?????????
						tasks[idx].ItalyPeople = uim
					case strings.Contains(selection.Text(), "???????????????"): // ?????????
						tasks[idx].SpainPeople = uim
					case strings.Contains(selection.Text(), "???????????????"): // ??????
						tasks[idx].EnglandPeople = uim
					case strings.Contains(selection.Text(), "????????????"): // ?????????
						tasks[idx].CanadaPeople = uim
					case strings.Contains(selection.Text(), "??????????????????????????????"): // ????????????
						tasks[idx].BangladeshPeople = uim
					case strings.Contains(selection.Text(), "?????????????????????"): // ?????????
						tasks[idx].CambodiaPeople = uim
					case strings.Contains(selection.Text(), "????????????????????????"): // ?????????
						tasks[idx].SingaporePeople = uim
					case strings.Contains(selection.Text(), "?????????"): // ??????
						tasks[idx].SwitzerlandPeople = uim
					case strings.Contains(selection.Text(), "????????????"): // ??????
						tasks[idx].GermanyPeople = uim
					case strings.Contains(selection.Text(), "?????????????????????"): // ????????????
						//case strings.Contains(selection.Text(), "?????????????????????"): // ????????????
						tasks[idx].AustraliaPeople = uim
					case strings.Contains(selection.Text(), "?????????"): // ??????
						tasks[idx].OtherPeople = uim
					case strings.Contains(selection.Text(), "??????"): // ??????
						tasks[idx].TotalPeople = uim
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
					return
				}
				tasks[idx].CourseInfoJson = string(marshal)
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
									tasks[idx].EnterPostgraduate = uint16(nu)
								case 1:
									tasks[idx].EnterUniversity = uint16(nu)
								case 2:
									tasks[idx].EnterJuniorCollege = uint16(nu)
								case 3:
									tasks[idx].EnterCollegeTechnology = uint16(nu)
								case 4:
									tasks[idx].EnterVocationalSchool = uint16(nu)
								case 5:
									tasks[idx].EnterVariousSchools = uint16(nu)
								case 6:
									tasks[idx].EnterOtherSchools = uint16(nu)
								}
							}
						})
					case 2:
						tasks[idx].University = strings.TrimSpace(selection.Find("td").Text())
					}
				})
			}
		})

		tasks[idx].GoogleMap = reader.Find("iframe").AttrOr("src", "")
	}

	return tasks, nil
}

func (l *LanguageSchoolUpdate) genL2Item(key string) ([]models.LanguageSchool, error) {
	url := fmt.Sprintf("https://www.nisshinkyo.org/search/%s", key)
	code, rdata, err := urllib.Get(url).SetTimeout(time.Second * 10).Byte()
	if err != nil {
		return nil, err
	}

	if code != 200 {
		return nil, errors.New(string(rdata))
	}

	var ls []models.LanguageSchool

	city := l.parseKey(key)

	reader, err := goquery.NewDocumentFromReader(bytes.NewReader(rdata))
	if err != nil {
		return nil, err
	}

	reader.Find(".termsDetail").Find("tr").Each(func(i int, selection *goquery.Selection) {
		if i == 0 {
			return
		}
		schoolName := selection.Find("th").Text()
		schoolAddr := selection.Find("td").Text()
		thIn := selection.Find("a").AttrOr("href", "")
		ur := fmt.Sprintf("https://www.nisshinkyo.org/search/%s", strings.ReplaceAll(thIn, "lng=3", "lng=1"))

		it := models.LanguageSchool{
			BaseModel: models.BaseModel{
				ID: uuid.New().String(),
			},
			NameCh:  schoolName,
			Addr:    schoolAddr,
			NissUrl: ur,
			City:    city,
		}

		ls = append(ls, it)
	})

	return ls, nil
}

func (l *LanguageSchoolUpdate) parseKey(key string) (city string) {
	split := strings.Split(key, "=")

	if len(split) >= 3 {
		city = strings.ReplaceAll(split[2], "#terms", "")
	}
	return city
}
