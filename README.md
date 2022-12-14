# AnalysisJapaneseSchools
Analysis of Japanese schools  (日本語学校・専門学校・大学の統計分析)

(國内)日本語學習(至少到N5)  -> (日本)日本語學校 -> EJU留考 -> 日本大學/專門學校


### 數據采集：

#### 日本語学校：　

日本語教育振興協会：

https://www.nisshinkyo.org/index.php

``` 
type LanguageSchool struct {
	gorm.Model
	// basic info
	NissID  string `gorm:"type:varchar(100);index" json:"niss_id"` // 日本語教育振興協会 id
	Addr    string `gorm:"type:text" json:"addr"`
	Tel     string `gorm:"type:varchar(300)" json:"tel"`
	Website string `gorm:"type:varchar(300)" json:"website"`
	Email   string `gorm:"type:varchar(300)" json:"email"`

	// addr info
	Province string `gorm:"type:varchar(100);index" json:"province"` // 省
	City     string `gorm:"type:varchar(100);index" json:"city"`     // 市
	Area     string `gorm:"type:varchar(100);index" json:"area"`     // 區

	// school info
	SchoolType           SchoolType `gorm:"type:varchar(100)" json:"school_type"`          // 學校類型: 財團法人/學校法人
	SchoolMaster         string     `gorm:"type:varchar(200)" json:"school_master"`        // 校長
	TeachingStartTime    string     `gorm:"type:varchar(200)" json:"teaching_start_time"`  // 教學開始時間
	CertificationPeriod  string     `gorm:"type:varchar(200)" json:"certification_period"` // 认定期间
	NumberOfTeachers     uint16     `json:"number_of_teachers"`                            // 教员人数
	NumberOfTeachersFull uint16     `json:"number_of_teachers_full"`                       // 教员人数 專職
	Quota                uint16     `json:"quota"`                                         // 名額

	// 留學人數
	PeopleInfoStatisticsTime string `gorm:"type:varchar(200)" json:"people_info_statistics_time"` // 留學人數統計時間
	TotalPeople              uint16 `json:"total_people"`                                         // 縂人數
	ChinesePeople            uint16 `json:"chinese_people"`                                       // 中國留學生人數
	OtherPeopleJson          string `gorm:"type:text" json:"other_people_json"`                   // people json  {"us": 100, "ru": 20, "hk": 8 ...}

	CourseInfoJson         string `gorm:"type:text" json:"course_info_json"`                  // 课程信息  {"認定コース": "", "目的": ...}
	JLPTInfoJson           string `gorm:"type:text" json:"jlpt_info_json"`                    // JLPT课程信息 {"n1": {"total": 30, "ok": 20} ...}
	JLPTInfoStatisticsTime string `gorm:"type:varchar(200)" json:"jlpt_info_statistics_time"` // JLPT课程信息統計時間
	StudyInfoJson          string `gorm:"type:text" json:"study_info_json"`                   // 進學統計時間
	University             string `gorm:"type:text" json:"university"`                        // 升學大學

	// other
	GoogleMap string `gorm:"type:text" json:"google_map"`
}
```

#### 専門学校＆大学

スタディサプリ：

https://shingakunet.com/searchList/ksl_senkaku/

マナビジョン：

https://manabi.benesse.ne.jp/

大学受験パスナビ

https://passnavi.evidus.com/




