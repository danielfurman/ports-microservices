package ports

import "errors"

type Port struct {
	ID          string
	Name        string
	City        string
	Country     string
	Alias       []string
	Regions     []string
	Coordinates []float64
	Province    string
	Timezone    string
	Unlocs      []string
	Code        string
}

func (p Port) Validate() error {
	if p.ID == "" {
		return errors.New("ID is required")
	}
	if p.Name == "" {
		return errors.New("name is required")
	}

	return nil
}
