// +build unit

package main

import (
	"testing"
	"time"
)



func TestUnitValidUsername(t *testing.T) {
	tests := []struct {
		username string
		valid    bool
	}{
		{"validName", true},
		{"Invalid Name", false},
		{"123Invalid", false},
		{"validname123", false},
	}

	for _, test := range tests {
		if got := ValidUsername(test.username); got != test.valid {
			t.Errorf("ValidUsername(%q) = %v; want %v", test.username, got, test.valid)
		}
	}
}

func TestUnitValidDate(t *testing.T) {
	tests := []struct {
		date  string
		valid bool
	}{
		{"2000-01-01", true},
		{"3000-01-01", false}, // Future date
		{"invalid-date", false},
	}

	for _, test := range tests {
		if got := ValidDate(test.date); got != test.valid {
			t.Errorf("ValidDate(%q) = %v; want %v", test.date, got, test.valid)
		}
	}
}

func TestUnitDaysUntilNextBirthday(t *testing.T) {
	today := time.Now().UTC().Truncate(24 * time.Hour)
	tests := []struct {
		date string
		want int
	}{
		{today.AddDate(0, 0, 1).Format("2006-01-02"), 1},
		{today.AddDate(0, 0, -1).Format("2006-01-02"), 364},
		{today.Format("2006-01-02"), 0},
	}

	for _, test := range tests {
		if got := daysUntilNextBirthday(test.date); got != test.want {
			t.Errorf("daysUntilNextBirthday(%q) = %d; want %d", test.date, got, test.want)
		}
	}
}
