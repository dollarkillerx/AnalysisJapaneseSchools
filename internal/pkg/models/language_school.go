package models

import (
	"gorm.io/gorm"
	"time"
)

type LanguageSchool struct {
	gorm.Model
	// basic info
	Name   string `gorm:"type:varchar(300);index" json:"name"`
	NameCh string `gorm:"type:varchar(300);index" json:"name_ch"`

	UpDataTime time.Time `json:"up_data_time"`

	NissID         string `gorm:"type:varchar(100);index" json:"niss_id"` // 日本語教育振興協会 id
	NissUrl        string `gorm:"type:text" json:"niss_url"`
	Addr           string `gorm:"type:text" json:"addr"`
	Tel            string `gorm:"type:varchar(300)" json:"tel"`
	NearestStation string `gorm:"type:text" json:"nearest_station"` // 最近车站
	Website        string `gorm:"type:varchar(300)" json:"website"`
	Email          string `gorm:"type:varchar(300)" json:"email"`
	Founder        string `gorm:"type:varchar(300)" json:"founder"` // 設置者名

	// addr info
	Province string `gorm:"type:varchar(100);index" json:"province"` // 省
	City     string `gorm:"type:varchar(100);index" json:"city"`     // 市
	Area     string `gorm:"type:varchar(100);index" json:"area"`     // 區

	// school info
	SchoolType               SchoolType `gorm:"type:varchar(100)" json:"school_type"`                // 學校類型: 財團法人/學校法人
	SchoolMaster             string     `gorm:"type:varchar(200)" json:"school_master"`              // 校長
	TeachingStartTime        time.Time  `gorm:"type:varchar(200)" json:"teaching_start_time"`        // 教學開始時間
	CertificationPeriodStart time.Time  `gorm:"type:varchar(200)" json:"certification_period_start"` // 认定期间
	Representative           string     `gorm:"type:varchar(200)" json:"representative"`             // 代表者名
	NumberOfTeachers         uint16     `json:"number_of_teachers"`                                  // 教员人数
	NumberOfTeachersFull     uint16     `json:"number_of_teachers_full"`                             // 教员人数 專職
	Quota                    string     `gorm:"type:varchar(200)" json:"quota"`                      // 名額

	// 留學人數
	PeopleInfoStatisticsTime string `gorm:"type:varchar(200)" json:"people_info_statistics_time"` // 留學人數統計時間

	CourseInfoJson string `gorm:"type:text" json:"course_info_json"` // 课程信息  {"認定コース": "", "目的": ...}
	University     string `gorm:"type:text" json:"university"`       // 升學大學

	// other
	GoogleMap string `gorm:"type:text" json:"google_map"`

	ChinaPeople       uint16 `json:"china_people"`        // 中国
	KoreaPeople       uint16 `json:"korea_people"`        // 韓国
	TaiwanPeople      uint16 `json:"taiwan_people"`       // 台湾
	VietnamPeople     uint16 `json:"vietnam_people"`      // 越南
	NepalPeople       uint16 `json:"nepal_people"`        // 尼泊尔
	ThailandPeople    uint16 `json:"thailand_people"`     // 泰国
	MyanmarPeople     uint16 `json:"myanmar_people"`      //  缅甸
	MongoliaPeople    uint16 `json:"mongolia_people"`     // 蒙古
	IndonesiaPeople   uint16 `json:"indonesia_people"`    // 印度尼西亚
	SriLankaPeople    uint16 `json:"sri_lanka_people"`    // 斯里兰卡
	SwedenPeople      uint16 `json:"sweden_people"`       // 瑞典
	MalaysiaPeople    uint16 `json:"malaysia_people"`     // 马来西亚
	AmericaPeople     uint16 `json:"america_people"`      // 美国
	IndiaPeople       uint16 `json:"india_people"`        // 印度
	FrancePeople      uint16 `json:"france_people"`       // 法国
	RussiaPeople      uint16 `json:"russia_people"`       // 俄罗斯
	PhilippinesPeople uint16 `json:"philippines_people"`  // 菲律宾
	SaudiArabiaPeople uint16 `json:"saudi_arabia_people"` // 沙特阿拉伯
	ItalyPeople       uint16 `json:"italy_people"`        // 意大利
	SpainPeople       uint16 `json:"spain_people"`        // 西班牙
	EnglandPeople     uint16 `json:"england_people"`      // 英国
	CanadaPeople      uint16 `json:"canada_people"`       // 加拿大
	BangladeshPeople  uint16 `json:"bangladesh_people"`   // 孟加拉国
	CambodiaPeople    uint16 `json:"cambodia_people"`     // 柬埔寨
	SingaporePeople   uint16 `json:"singapore_people"`    // 新加坡
	SwitzerlandPeople uint16 `json:"switzerland_people"`  // 瑞士
	GermanyPeople     uint16 `json:"germany_people"`      // 德国
	AustraliaPeople   uint16 `json:"australia_people"`    // 澳大利亚
	OtherPeople       uint16 `json:"other_people"`        // その他
	TotalPeople       uint16 `json:"total_people"`        // 合計

	EnterPostgraduate      uint16 `json:"enter_postgraduate"`       // 进学研究生
	EnterUniversity        uint16 `json:"enter_university"`         // 进学大学
	EnterJuniorCollege     uint16 `json:"enter_junior_college"`     // 进学短期大学
	EnterCollegeTechnology uint16 `json:"enter_college_technology"` // 进学高等専門学校
	EnterVocationalSchool  uint16 `json:"enter_vocational_school"`  // 进学専門学校
	EnterVariousSchools    uint16 `json:"enter_various_schools"`    // 各种学校
	EnterOtherSchools      uint16 `json:"enter_other_schools"`      // 其他学校

}
