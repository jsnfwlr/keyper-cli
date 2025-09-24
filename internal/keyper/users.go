package keyper

import (
	"encoding/json"
	"fmt"

	"github.com/jsnfwlr/keyper-cli/internal/feedback"
)

func (c *Client) GetUsers() (users []UserResponse, fault error) {
	u := []UserResponse{}

	err := c.Do("GET", "/api/users", true, nil, &u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (c *Client) GetUser(username string) (user UserResponse, fault error) {
	u := []UserResponse{}

	err := c.Do("GET", fmt.Sprintf("/api/users/%s", username), true, nil, &u)
	if err != nil {
		return UserResponse{}, err
	}

	if !feedback.ExceedsLimit(feedback.Debug) {
		b, err := json.MarshalIndent(u, "", "  ")
		if err == nil {
			feedback.Printf(feedback.Debug, "user: %s", string(b))
		} else {
			feedback.Printf(feedback.Error, "user: %+v", err)
		}

	}

	return u[0], nil
}
