package core

import (
	"strconv"
	"strings"
	"time"
)

type (
	JSONBool      bool
	JSONFloat     float64
	JSONTime      time.Time
	JSONTimestamp time.Time
	JSONUint16    uint16
	JSONInt       int
	JSONBytes     []byte
)

func (j *JSONBool) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	bValue, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	*j = JSONBool(bValue)
	return nil
}

func (j JSONBool) ToBool() bool {
	return bool(j)
}

func (j *JSONFloat) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	fValue, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*j = JSONFloat(fValue)
	return nil
}

func (j JSONFloat) ToFloat64() float64 {
	return float64(j)
}

func (j *JSONTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	tValue, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	*j = JSONTime(time.Unix(tValue, 0))
	return nil
}

func (j JSONTime) ToTime() time.Time {
	return time.Time(j)
}

func (j *JSONTimestamp) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	tValue, err := time.Parse("2006-01-02 15:04:05.999999-07", s)
	if err != nil {
		return err
	}
	*j = JSONTimestamp(tValue)
	return nil
}

func (j JSONTimestamp) ToTime() time.Time {
	return time.Time(j)
}

func (j *JSONUint16) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	tValue, err := strconv.ParseInt(s, 10, 16)
	if err != nil {
		return err
	}
	*j = JSONUint16(tValue)
	return nil
}

func (j JSONUint16) ToUint16() uint16 {
	return uint16(j)
}

func (j *JSONInt) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	tValue, err := strconv.ParseInt(s, 10, 16)
	if err != nil {
		return err
	}
	*j = JSONInt(tValue)
	return nil
}

func (j JSONInt) ToInt() int {
	return int(j)
}

func (j *JSONBytes) UnmarshalJSON(b []byte) error {
	*j = JSONBytes([]byte(strings.Trim(string(b), "\"")))
	return nil
}
