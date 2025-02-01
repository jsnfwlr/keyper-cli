package keyper

func (c *Client) GetUsers() (users []UserResponse, fault error) {
	u := []UserResponse{}

	err := c.Do("GET", "/api/users", true, nil, &u)
	if err != nil {
		return nil, err
	}

	return u, nil
}
