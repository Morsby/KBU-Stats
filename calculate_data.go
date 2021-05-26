package kbu

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// calculateRound calculates the given round from
// the url (of form https://kbu.logbog.net/Ajax_getYYYYv[12].asp) and returns a more readable format.
func calculateRound(str string) (Round, error) {
	var year int
	var season Season
	var err error
	if str == latestUrl {
		year = 2020
		season = SeasonFall
	} else {
		yStr := str[31:35]
		year, err = strconv.Atoi(yStr)
		if err != nil {
			return Round{}, err
		}

		switch str[36] {
		case '1':
			season = SeasonSpring
		case '2':
			season = SeasonFall
		default:
			err = errors.New("cannot calculate season")
		}
	}
	if err != nil {
		return Round{}, err
	}
	return Round{Year: year, Season: season}, err
}

// calculateNumberUni calculates the pick number af university of the doctor
// or retuns an error if it's an invalid string (should be of form "n UNI")
func calculateNumberUni(s string) (number int, uni string, err error) {
	// number not selected
	if s == "" {
		return -1, "", nil
	}

	s = strings.ReplaceAll(s, " ", " ") // fix for weird spacing in raw data

	uniNumberRegexp := regexp.MustCompile(`^([0-9]*)\s*([A-za-z]*)$`)
	vals := uniNumberRegexp.FindStringSubmatch(s)
	if len(vals) != 3 {
		return 0, "", errors.New("cannot calculate number and uni; invalid input format")
	}
	if vals[1] != "" {
		number, err = strconv.Atoi(vals[1])
		if err != nil {
			return 0, "", err
		}
	}

	uni = vals[2]
	return number, uni, nil
}

func calculateDate(str string, yearOverride int) (date time.Time, err error) {
	if str == "" {
		return time.Time{}, nil
	}

	loc, err := time.LoadLocation("Europe/Copenhagen")
	if err != nil {
		return time.Time{}, err
	}

	var day int
	var month time.Month
	var year int
	var hour int
	var min int

	// Format: 10. okt 13:20; needs year || 1. feb 15
	// $1 = date
	// $2 = month
	// $3 = year, optional
	// $4 = time (needs year), optional
	str = strings.ReplaceAll(str, ".", "")
	dateRegexp := regexp.MustCompile(`^([0-9]+)\s+([a-z]+)\s+([0-9]{2})?\s*([0-9]{2}:[0-9]{2})?$`)
	times := dateRegexp.FindStringSubmatch(str)

	// day
	day, err = strconv.Atoi(times[1])
	if err != nil {
		return time.Time{}, err
	}

	// month
	switch times[2] {
	case "jan":
		month = time.January
	case "feb":
		month = time.February
	case "mar":
		month = time.March
	case "apr":
		month = time.April
	case "maj":
		month = time.May
	case "jun":
		month = time.June
	case "jul":
		month = time.July
	case "aug":
		month = time.August
	case "sep":
		month = time.September
	case "okt":
		month = time.October
	case "nov":
		month = time.November
	case "dec":
		month = time.December
	default:
		err = errors.New("cannot find month")
	}

	if err != nil {
		return time.Time{}, err
	}

	// year
	if times[3] != "" {
		year, err = strconv.Atoi(times[3])
		year += 2000 // hard fix for years
		if err != nil {
			return time.Time{}, err
		}
	}

	if yearOverride != 0 {
		year = yearOverride
	}

	// time of day
	if times[4] != "" {
		hour, err = strconv.Atoi(times[4][:2])
		if err != nil {
			return time.Time{}, err
		}

		min, err = strconv.Atoi(times[4][3:])
		if err != nil {
			return time.Time{}, err
		}
	}

	return time.Date(year, month, day, hour, min, 0, 0, loc), nil
}