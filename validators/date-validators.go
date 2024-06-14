package validators

import (
	"errors"
	"fmt"
	"time"
)

func InterfaceToDate(i interface{}) *time.Time {
	dateStr := i.(string)
	layouts := []string{
		"2006-01-02",
		time.Layout,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	}

	for _, layout := range layouts {
		if validDatetime, err := time.Parse(layout, dateStr); err == nil {
			return &validDatetime
		}
	}

	return nil
}

func IsValidDate(i interface{}, _ map[string]interface{}) error {
	date := InterfaceToDate(i)
	if date == nil {
		return errors.New("invalid date provided")
	}
	return nil
}

func IsLessThanNow(i interface{}, _ map[string]interface{}) error {
	date := InterfaceToDate(i)
	now := time.Now()
	if now.After(*date) {
		layout := "2006-01-02 15:04:05"
		return errors.New(fmt.Sprintf("%s has been passed now", date.Format(layout)))
	}
	return nil
}

func IsMoreThanNow(i interface{}, _ map[string]interface{}) error {
	date := InterfaceToDate(i)
	now := time.Now()
	if now.Before(*date) {
		layout := "2006-01-02 15:04:05"
		return errors.New(fmt.Sprintf("%s has not yet passed", date.Format(layout)))
	}
	return nil
}

func IsBefore(i interface{}, attr map[string]interface{}) error {
	date := InterfaceToDate(i)
	comparableDate := InterfaceToDate(attr["maxTime"])

	if date.After(*comparableDate) {
		layout := "2006-01-02 15:04:05"
		return errors.New(fmt.Sprintf("%s is after %s", date.Format(layout), comparableDate.Format(layout)))
	}

	return nil
}
func IsAfter(i interface{}, attr map[string]interface{}) error {
	date := InterfaceToDate(i)
	comparableDate := InterfaceToDate(attr["maxTime"])

	if date.Before(*comparableDate) {
		layout := "2006-01-02 15:04:05"
		return errors.New(fmt.Sprintf("%s is after %s", date.Format(layout), comparableDate.Format(layout)))
	}

	return nil
}
func IsInBetweenTime(i interface{}, _ map[string]interface{}) error {
	date := InterfaceToDate(i)
	comparableMaxDate := InterfaceToDate(attr["maxTime"])
	comparableMinDate := InterfaceToDate(attr["minTime"])
	if !(date.Before(*comparableMaxDate) && date.After(*comparableMinDate)) {
		layout := "2006-01-02 15:04:05"
		return errors.New(fmt.Sprintf("%s is before %s or after %s", date.Format(layout), comparableMinDate.Format(layout), comparableMaxDate.Format(layout)))
	}
	return nil
}
