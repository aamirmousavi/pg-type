package pgtype

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

type Polygon []Point

// ((1,2),(2,2))
func (p *Polygon) Scan(val interface{}) (err error) {
	if bb, ok := val.([]uint8); ok {
		tmp := bb[2 : len(bb)-2]
		coors := strings.Split(string(tmp[:]), "),(")
		for _, coor := range coors {
			xy := strings.Split(coor, ",")
			x, err := strconv.ParseFloat(xy[0], 64)
			if err != nil {
				return err
			}
			y, err := strconv.ParseFloat(xy[1], 64)
			if err != nil {
				return err
			}
			*p = append(*p, Point{X: x, Y: y})
		}
	}
	return nil
}

func (p Polygon) Value() (driver.Value, error) {
	var str string
	for i, point := range p {
		if i == 0 {
			str += "("
		} else {
			str += ",("
		}
		str += fmt.Sprintf("%f, %f", point.X, point.Y)
	}
	str += ")"
	return str, nil
}

func (p Polygon) String() string {
	var str string
	for i, point := range p {
		if i == 0 {
			str += "("
		} else {
			str += ",("
		}
		str += fmt.Sprintf("%f, %f", point.X, point.Y)
	}
	str += ")"
	return str
}

type NullPolygon struct {
	Polygon
	Valid bool
}

func (np *NullPolygon) Scan(val interface{}) (err error) {
	if val == nil {
		np.Valid = false
		return nil
	}
	np.Valid = true
	return np.Polygon.Scan(val)
}

func (np NullPolygon) Value() (driver.Value, error) {
	if !np.Valid {
		return nil, nil
	}
	return np.Polygon.Value()
}

func (np NullPolygon) String() string {
	if !np.Valid {
		return ""
	}
	return np.Polygon.String()
}
