package csv

import (
	"encoding/csv"
	"strings"
)

func Parse(input string) ([][]string, error) {
	return csv.NewReader(strings.NewReader(input)).ReadAll()
}
