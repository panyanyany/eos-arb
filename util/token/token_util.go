package token

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseQuantity(qty string) (amt float64, symbolCode string, err error) {
	piece := strings.Split(qty, " ")
	if len(piece) != 2 {
		err = fmt.Errorf("invalid quantity format")
		return
	}
	amt, err = strconv.ParseFloat(piece[0], 64)
	if err != nil {
		err = fmt.Errorf("parsing amount: %w", err)
		return
	}
	return
}
