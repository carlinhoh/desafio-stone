package main

import (
	"time"
)

type Invoice struct {
	Id		*int64		`json:"id"`
	CreatedAt	time.Time	`json:"createdAt"`
	ReferenceMonth	int		`json:"referenceMonth"`
	ReferenceYear	int		`json:"referenceYear"`
	Document	string		`json:"document"`
	Description	*string		`json:"description,omitempty"`
	Amount		*float64	`json:"amount,omitempty"`
	IsActive	int       	`json:"-"`
	DeactiveAt	*time.Time	`json:"-"`
}