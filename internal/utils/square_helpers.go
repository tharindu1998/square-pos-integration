package utils

import (
	"strconv"
	"github.com/square/square-go-sdk/v2"
)

func SafeString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func SafeInt64(i *int64) int64 {
	if i != nil {
		return *i
	}
	return 0
}

func SafeCurrency(c *square.Currency) string {
	if c != nil {
		return string(*c) // Convert enum to string value like "USD"
	}
	return ""
}

func ParseQuantity(q string) int {
	i, _ := strconv.Atoi(q) // You may want error handling here
	return i
}

