package agent

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const tzNumeric = "[+-]{1}[0-9]{4}"

func tzIsNumeric(candidate string) (bool, error) {
	matched, err := regexp.MatchString(tzNumeric, candidate)
	if err != nil { // handle any error extrinsic to this function
		return false, errors.New(fmt.Sprintf("[ERR] %v", err))
	}

	return matched, nil
}

const tzAbbr = "[A-Z]{1,}"

func tzIsAbbr(candidate string) (bool, error) {
	matched, err := regexp.MatchString(tzAbbr, candidate)
	if err != nil { // handle any error extrinsic to this function
		return false, errors.New(fmt.Sprintf("[ERR] %v", err))
	}

	return matched, nil
}

var Tz = map[string]string{
	"EDT": "America/New_York",
}

// <pubDate> tag in RSS XML files contains the date the article was published
func parsePubDate(date string) (time.Time, error) {
	// locate the last space
	lastSpace := strings.LastIndex(date, " ")
	if lastSpace == -1 { // space not exist at all
		return time.Time{}, errors.New(fmt.Sprintf("[ERR] Date '%s' has wrong format", date))
	}
	if lastSpace == len(date)-1 { // last character is a space
		return time.Time{}, errors.New(fmt.Sprintf("[ERR] Date '%s' has wrong format", date))
	}

	tz := date[lastSpace+1 : len(date)] // extract time zone

	isNumeric, err := tzIsNumeric(tz)
	if err != nil { // error already formatted, return error as is
		return time.Time{}, err
	}
	if isNumeric { // Parse takes into consideration the time diff
		return time.Parse(time.RFC1123Z, date)
	}

	isAbbr, err := tzIsAbbr(tz)
	if err != nil { // error already formatted, return error as is
		return time.Time{}, err
	}
	if !isAbbr {
		return time.Time{}, errors.New(fmt.Sprintf("[ERR] Time zone '%s' has wrong format", tz))
	}

	tzVal, ok := Tz[tz]
	if !ok { // TODO log this
		return time.Time{}, errors.New(fmt.Sprintf("[ERR] Key '%s' not exists in dictionary", tz))
	}

	loc, _ := time.LoadLocation(tzVal)
	return time.ParseInLocation(time.RFC1123, date, loc) // parse with respect to the location
}
