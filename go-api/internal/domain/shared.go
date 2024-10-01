package domain

import (
	"encoding/json"
	"time"
)

const DateFormat = "2006-01-02"

type Date struct {
	time.Time
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time.Format(DateFormat))
}

func (d *Date) UnmarshalJSON(data []byte) error {
	// Strip the quotes around the date string
	dateStr := string(data)
	dateStr = dateStr[1 : len(dateStr)-1]

	// Parse the date string according to the custom format
	t, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return err
	}

	// Set the parsed time value into the Date struct
	d.Time = t
	return nil
}
