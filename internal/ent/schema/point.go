package schema

//there may be some problem because it is made by chatgpt
//the point.go is used to convey the point in mysql
import (
	"database/sql/driver"
	"fmt"
	"strings"
)

// Point represents a geographical point with latitude and longitude.
type Point struct {
	Lat float64 // Latitude
	Lng float64 // Longitude
}

// Value implements the driver.Valuer interface for database serialization.
func (p Point) Value() (driver.Value, error) {
	return fmt.Sprintf("POINT(%f %f)", p.Lng, p.Lat), nil
}

// Scan implements the sql.Scanner interface for database deserialization.
func (p *Point) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to convert value to string: %v", value)
	}

	// Parse the string to extract latitude and longitude.
	str = strings.TrimPrefix(str, "POINT(")
	str = strings.TrimSuffix(str, ")")
	coords := strings.Split(str, " ")
	if len(coords) != 2 {
		return fmt.Errorf("invalid POINT format: %s", value)
	}

	var lng, lat float64
	_, err := fmt.Sscanf(coords[0], "%f", &lng)
	if err != nil {
		return fmt.Errorf("failed to parse longitude: %v", err)
	}
	_, err = fmt.Sscanf(coords[1], "%f", &lat)
	if err != nil {
		return fmt.Errorf("failed to parse latitude: %v", err)
	}

	p.Lat = lat
	p.Lng = lng
	return nil
}
