package util

import (
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

const DateStringFormat = "2006-01-02"

type DateString struct {
	time.Time
}

func (date *DateString) UnmarshalJSON(input []byte) error {
	strInput := string(input[:11])
	strInput = strings.Trim(strInput, `"`)
	newTime, err := time.Parse(DateStringFormat, strInput)
	if err != nil {
		return err
	}
	date.Time = newTime
	return nil
}

func (date *DateString) MarshalJSON() ([]byte, error) {
	if date.IsZero() {
		return []byte("null"), nil
	}
	return []byte("\"" + date.Format(DateStringFormat) + "\""), nil
}

// Scan permet à go-pg de bien décoder la valeur en base
func (date *DateString) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	newTime, err := time.Parse(DateStringFormat, string(value.([]byte)))
	if err != nil {
		return err
	}
	date.Time = newTime
	return nil
}

func Date(datestr string) DateString {
	date, err := time.Parse("2006-01-02", datestr)
	if err != nil {
		log.Fatal().Msgf("unable to parse date: %v", err)
	}
	return DateString{date}
}

func DateTime(datestr string) time.Time {
	date, err := time.Parse("2006-01-02T15:04:05", datestr)
	if err != nil {
		log.Fatal().Msgf("unable to parse date: %v", err)
	}
	return date
}

func FormatDate(date time.Time) string {
	return date.Format("02/01/2006")
}

func FormatDateTime(date time.Time) string {
	return date.Format("02/01/2006 à 15h04")
}
