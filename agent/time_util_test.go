package agent

import (
	"testing"
	"time"
)

func Test_tzIsNumeric(t *testing.T) {
	testCases := []struct {
		in  string
		out bool
	}{
		{
			"+0000",
			true,
		},
		{
			"-0500",
			true,
		},
		{
			"+0",
			false,
		},
		{
			"+00",
			false,
		},
		{
			"+000",
			false,
		},
		{
			"UTC",
			false,
		},
	}

	for idx, testCase := range testCases {
		matched, err := tzIsNumeric(testCase.in)

		if err != nil {
			t.Error(err)
		}

		if matched != testCase.out {
			t.Errorf("[Test case %d] expecting %v, got %v", idx, testCase.out, matched)
		}
	}
}

func Test_tzIsAbbr(t *testing.T) {
	testCases := []struct {
		in  string
		out bool
	}{
		{
			"+",
			false,
		},
		{
			"0",
			false,
		},
		{
			"",
			false,
		},
		{
			" ",
			false,
		},
		{
			"A",
			true,
		},
		{
			"AB",
			true,
		},
		{
			"ABC",
			true,
		},
	}

	for idx, testCase := range testCases {
		matched, err := tzIsAbbr(testCase.in)

		if err != nil {
			t.Error(err)
		}

		if matched != testCase.out {
			t.Errorf("[Test case %d] expecting %v, got %v", idx, testCase.out, matched)
		}
	}
}

func Test_parsePubDate(t *testing.T) {
	testCases := []struct {
		in    string // in
		year  int    // out
		month time.Month
		day   int
		hour  int
		min   int
		sec   int
	}{
		{
			"Tue, 19 Apr 2016 17:25:18 +0000",
			2016,
			4,
			19,
			17,
			25,
			18,
		},
		{
			"Tue, 19 Apr 2016 17:25:18 EDT",
			2016,
			4,
			19,
			17,
			25,
			18,
		},
	}

	for idx, testCase := range testCases {
		parsed, err := parsePubDate(testCase.in)
		if err != nil {
			t.Error(err)
		}

		year, month, day := parsed.Date()
		hour, min, sec := parsed.Clock()
		if year != testCase.year || month != testCase.month || day != testCase.day || hour != testCase.hour || min != testCase.min || sec != testCase.sec {
			t.Errorf("[Test case %d] expecting (year, month, day, hour, min, sec) as (%d, %d, %d, %d, %d, %d), got (%d, %d, %d, %d, %d, %d)", idx, year, month, day, hour, min, sec, testCase.year, testCase.month, testCase.day, testCase.hour, testCase.min, testCase.sec)
		}
	}
}

func Test_parsePubDate_ConvertToUtc(t *testing.T) {
	testCases := []struct {
		in    string // in
		year  int    // out
		month time.Month
		day   int
		hour  int
		min   int
		sec   int
	}{
		{
			"Tue, 19 Apr 2016 17:25:18 +0000",
			2016,
			4,
			19,
			17,
			25,
			18,
		},
		{
			"Tue, 19 Apr 2016 17:25:18 +0100",
			2016,
			4,
			19,
			16,
			25,
			18,
		},
		{
			"Tue, 19 Apr 2016 17:25:18 -0100",
			2016,
			4,
			19,
			18,
			25,
			18,
		},
		{
			"Tue, 19 Apr 2016 17:25:18 EDT",
			2016,
			4,
			19,
			21,
			25,
			18,
		},
	}

	for idx, testCase := range testCases {
		parsed, err := parsePubDate(testCase.in)
		if err != nil {
			t.Error(err)
		}

		parsed = parsed.In(time.UTC)

		year, month, day := parsed.Date()
		hour, min, sec := parsed.Clock()
		if year != testCase.year || month != testCase.month || day != testCase.day || hour != testCase.hour || min != testCase.min || sec != testCase.sec {
			t.Errorf("[Test case %d] expecting (year, month, day, hour, min, sec) as (%d, %d, %d, %d, %d, %d), got (%d, %d, %d, %d, %d, %d)", idx, year, month, day, hour, min, sec, testCase.year, testCase.month, testCase.day, testCase.hour, testCase.min, testCase.sec)
		}
	}
}
