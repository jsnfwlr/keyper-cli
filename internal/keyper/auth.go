package keyper

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	AccessToken   string         `json:"access_token"`
	AccountLocked bool           `json:"accountLocked"`
	CN            string         `json:"cn"`
	DisplayName   string         `json:"displayName"`
	DN            string         `json:"dn"`
	GivenName     string         `json:"givenName"`
	Mail          string         `json:"mail"`
	MemberOfs     []string       `json:"memberOfs"`
	RefreshToken  string         `json:"refresh_token"`
	SN            string         `json:"sn"`
	SSHPublicKeys []SSHPublicKey `json:"sshPublicKeys"`
	UID           string         `json:"uid"`
}

func (c *Client) Auth() error {
	payload := &authRequest{
		Username: c.config.Username,
		Password: c.config.Password,
	}

	var response authResponse
	err := c.Do("POST", "/api/login", false, payload, &response)
	if err != nil {
		return err
	}

	c.token = response.AccessToken

	return nil
}
