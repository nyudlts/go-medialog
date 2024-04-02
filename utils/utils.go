package utils

import (
	"fmt"
	"time"
)

type Pagination struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Sort   string `json:"sort"`
}

func Add(a int, b int) int { return a + b }

func Subtract(a int, b int) int { return a - b }

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%d-%d", year, month, day)
}
