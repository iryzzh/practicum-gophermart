package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type Time struct {
	time.Time
	Valid bool
}

func (s *Time) String() string {
	// 2020-12-10T15:12:01+03:00
	_, offset := s.Zone()
	return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d+%02d:00",
		s.Year(),
		s.Month(),
		s.Day(),
		s.Hour(),
		s.Minute(),
		s.Second(),
		offset/3600)
}

func (s *Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Time) Scan(value interface{}) error {
	if value == nil {
		s.Time, s.Valid = time.Time{}, false
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		s.Time, s.Valid = v, true
		return nil
	}

	return fmt.Errorf("can't convert %T to time.Time", value)
}
