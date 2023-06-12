package models

import "time"

type Currency struct {
	Model
	Symbol         string          `gorm:"uniqueIndex;idx_symbol;not null" json:"symbol"`
	Name           string          `gorm:"not null" json:"name"`
	CurrencyPrices []CurrencyPrice `gorm:"foreignKey:CurrencyID"`
}

type CurrencyPrice struct {
	Model
	Date       time.Time `gorm:"not null" json:"date"`
	Open       float64   `gorm:"not null" json:"open"`
	High       float64   `gorm:"not null" json:"high"`
	Low        float64   `gorm:"not null" json:"low"`
	Close      float64   `gorm:"not null" json:"close"`
	Volume     float64   `gorm:"not null" json:"volume"`
	CurrencyID int
}
