package keyper

import "fmt"

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

	return u[0], nil
}
