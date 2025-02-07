package models

import (
	"strings"
	"time"
)


type CivilTime time.Time

func (c *CivilTime) UnmarshalJSON(b []byte) error {
    value := strings.Trim(string(b), `"`) //get rid of "
    if value == "" || value == "null" {
        return nil
    }

    t, err := time.Parse("2006-01-02", value) //parse time
    if err != nil {
        return err
    }
    *c = CivilTime(t) //set result using the pointer
    return nil
}

func (c CivilTime) MarshalJSON() ([]byte, error) {
    return []byte(`"` + time.Time(c).Format("2006-01-02") + `"`), nil
}