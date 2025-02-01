package prompter

import (
	"fmt"
)

// GetTimezone - prompt the user to enter a timezone, starting with their region and then the city
//
// Params:
//   - question: the question to ask the user
func GetTimezone(question string) (string, error) {
	regions, cities := Timezones()
	region, err := Select(fmt.Sprintf("%s: Region", question), false, regions...)
	if err != nil {
		return "", err
	}

	city, err := Select(fmt.Sprintf("%s: City", question), false, cities[region]...)
	if err != nil {
		return "", err
	}

	return region + "/" + city, nil
}
