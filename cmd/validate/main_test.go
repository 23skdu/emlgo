package main

import (
	"testing"
)

func TestValidateCommands(t *testing.T) {
	verbose = false
	failedOnly = false
	typeFilter = ""

	validateIntTypes()
	validateUintTypes()
	validateFloatTypes()
	validateComplexTypes()
	printSummary()
}
