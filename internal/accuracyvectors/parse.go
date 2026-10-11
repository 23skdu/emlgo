package accuracyvectors

import (
	"encoding/csv"
	"errors"
	"strings"
)

var ErrInvalidFixture = errors.New("invalid MPFR vector fixture header")

// Rows returns the MPFR CSV records after the provenance line and header.
func Rows() ([][]string, error) {
	reader := csv.NewReader(strings.NewReader(MPFR))
	reader.Comment = '#'
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) == 0 || len(records[0]) != 4 || records[0][0] != "operation" {
		return nil, ErrInvalidFixture
	}
	return records[1:], nil
}
