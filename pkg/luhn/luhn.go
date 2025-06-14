package luhn

import (
	"strconv"
	"strings"
)

type LuhnService struct{}

func New() *LuhnService {
	return &LuhnService{}
}

func (*LuhnService) IsValid(number string) bool {
	number = strings.ReplaceAll(number, " ", "")
	if len(number) == 0 {
		return false
	}

	for _, r := range number {
		if r < '0' || r > '9' {
			return false
		}
	}

	digits := make([]int, len(number))
	for i, r := range number {
		digit, err := strconv.Atoi(string(r))
		if err != nil {
			return false
		}
		digits[i] = digit
	}

	sum := 0
	doubleNext := false

	for i := len(digits) - 1; i >= 0; i-- {
		digit := digits[i]

		if doubleNext {
			digit *= 2
			if digit > 9 {
				digit = digit/10 + digit%10
			}
		}

		sum += digit
		doubleNext = !doubleNext
	}

	return sum%10 == 0
}
