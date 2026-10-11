package accuracyvectors

import (
	"testing"
)

func TestRows(t *testing.T) {
	rows, err := Rows()
	if err != nil {
		t.Fatalf("Rows() failed: %v", err)
	}
	if len(rows) == 0 {
		t.Fatalf("expected rows, got 0")
	}
	for i, row := range rows {
		if len(row) != 4 {
			t.Errorf("row %d: expected 4 columns, got %d: %v", i, len(row), row)
		}
	}
}

func TestRowsInvalid(t *testing.T) {
	old := MPFR
	defer func() { MPFR = old }()

	MPFR = "invalid,header\n1,2\n"
	if _, err := Rows(); err != ErrInvalidFixture {
		t.Errorf("expected ErrInvalidFixture, got %v", err)
	}

	MPFR = "operation,1,2,3\n\"unclosed quote"
	if _, err := Rows(); err == nil {
		t.Errorf("expected error on invalid csv, got nil")
	}

	MPFR = ""
	if _, err := Rows(); err != ErrInvalidFixture {
		t.Errorf("expected ErrInvalidFixture on empty string, got %v", err)
	}
}
