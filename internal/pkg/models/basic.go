package models

import "gorm.io/gorm"

type AdministrativeArea struct {
	gorm.Model
	Name   string `gorm:"type:varchar(300)" json:"name"`
	Level  int    `json:"level"`
	Father string `gorm:"type:varchar(300);index" json:"father"`
}

type SchoolType string

const (
	FinancialCorporation SchoolType = "financial_corporation" // 財團法人
	SchoolCorporation    SchoolType = "school_corporation"    // 學校法人
)
