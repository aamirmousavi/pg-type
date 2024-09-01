package pgtype

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

type Point struct {
	X float64
	Y float64
}

func (p *Point) Scan(val interface{}) (err error) {
	if bb, ok := val.([]uint8); ok {
		tmp := bb[1 : len(bb)-1]
		coors := strings.Split(string(tmp[:]), ",")
		if p.X, err = strconv.ParseFloat(coors[0], 64); err != nil {
			return err
		}
		if p.Y, err = strconv.ParseFloat(coors[1], 64); err != nil {
			return err
		}
	}
	return nil
}

func (p Point) Value() (driver.Value, error) {
	return fmt.Sprintf("(%f, %f)", p.X, p.Y), nil
}

func (p Point) String() string {
	return fmt.Sprintf("(%v, %v)", p.X, p.Y)
}

type NullPoint struct {
	Point
	Valid bool
}

func (np *NullPoint) Scan(val interface{}) (err error) {
	if val == nil {
		np.Valid = false
		return nil
	}
	np.Valid = true
	return np.Point.Scan(val)
}

func (np NullPoint) Value() (driver.Value, error) {
	if !np.Valid {
		return nil, nil
	}
	return np.Point.Value()
}

func (np NullPoint) String() string {
	if !np.Valid {
		return ""
	}
	return np.Point.String()
}
