package keyper

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jsnfwlr/keyper-cli/internal/feedback"
)

/*
[

	{
	  "accountLocked": false,
	  "cn": "alice",
	  "displayName": "Alice Parker",
	  "dn": "cn=alice,ou=people,dc=keyper,dc=example,dc=org",
	  "givenName": "Alice",
	  "mail": "alice@dbsentry.com",
	  "memberOfs": [
	    "cn=demo_servers,ou=groups,dc=keyper,dc=example,dc=org"
	  ],
	  "sn": "Parker",
	  "sshPublicKeys": [
	    {
	      "dateExpire": "20201204",
	      "hostGroups": [
	        "cn=demo_servers,ou=groups,dc=keyper,dc=example,dc=org"
	      ],
	      "key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC1KtJpPn6W9W5WgPU8+eYuuSKKyHA+Z62mVLYp50Ch/MMTUSxcFF/V1H81CStU4OrPv/pUxpHtqSDeTCMbVtTmP0Bbc5V7rCYQVgfhTB7CzKAwnfJSfJGY/JoJLCrC4kt40PMwyXTHiPUkrs4tOHiv7GIT4aZI/wmVPrg8x6oBFRgfCl1TQVgeSQl2kAnjkUHEsq2CsnZR9mKIJ31CWzeHLotYHNg82jmgylCWUsl6Pd5eigObUtk0j6Vnjn7FUKwSmffhEPInU1K+IzYMdFe1QElTSO7X+IOjedQZ2Y8nt3U9N9WPyd7FK13Sn8Ij22CIMmTuvfNXv/H4ja9vF0Ob"
	    }
	  ],
	  "uid": "alice"
	}

]
*/
type UserResponse struct {
	AccountLocked bool           `json:"accountLocked,omitempty"`
	CN            string         `json:"cn,omitempty"`
	DisplayName   string         `json:"displayName,omitempty"`
	DN            string         `json:"dn,omitempty"`
	GivenName     string         `json:"givenName,omitempty"`
	Mail          string         `json:"mail,omitempty"`
	MemberOfs     []string       `json:"memberOfs,omitempty"`
	SN            string         `json:"sn,omitempty"`
	SSHPublicKeys []SSHPublicKey `json:"sshPublicKeys"`
	UID           string         `json:"uid,omitempty"`
}

type DateExpire string

type SSHPublicKey struct {
	CN          string     `json:"cn,omitempty"`
	DateExpire  DateExpire `json:"dateExpire,omitempty"`
	Date        time.Time  `json:"-"`
	HostGroups  []string   `json:"hostGroups,omitempty"`
	Key         string     `json:"key"`
	KeyId       int        `json:"keyid,omitempty"`
	Fingerprint string     `json:"fingerprint"`
	Name        string     `json:"name"`
	KeyType     int        `json:"keyType,omitempty"`
	Local       bool       `json:"-"`
}

func (d DateExpire) Parse() (time.Time, error) {
	year := d[0:4]
	month := d[4:6]
	day := d[6:8]
	hours := d[8:10]
	minutes := d[10:12]
	seconds := d[12:14]

	return time.Parse("2006-01-02T15:04:05", fmt.Sprintf("%s-%s-%sT%s:%s:%s", year, month, day, hours, minutes, seconds))
}

func (c *Client) GetSSHKeys(username, host string) ([]SSHPublicKey, error) {
	var response []UserResponse
	err := c.Do("GET", fmt.Sprintf("/api/users/%s", username), true, nil, &response)
	if err != nil {
		return nil, err
	}

	for i := range response[0].SSHPublicKeys {
		date, err := response[0].SSHPublicKeys[i].DateExpire.Parse()
		if err != nil {
			return nil, err
		}
		response[0].SSHPublicKeys[i].Date = date

		if strings.Contains(response[0].SSHPublicKeys[i].Name, host) {
			response[0].SSHPublicKeys[i].Local = true
		}
	}

	return response[0].SSHPublicKeys, nil
}

func (c Client) RevokeSSHKey(username string, key SSHPublicKey) error {
	response := []UserResponse{}

	key.HostGroups = []string{}
	key.DateExpire = ""

	payload := UserResponse{
		SSHPublicKeys: []SSHPublicKey{key},
	}

	err := c.Do("PUT", fmt.Sprintf("/api/users/%s", username), true, payload, &response)
	if err != nil {
		return err
	}

	return nil
}

func (c Client) AddSSHKey(username string, key SSHPublicKey) error {
	response := []UserResponse{}

	key.KeyId = 0

	payload := UserResponse{
		SSHPublicKeys: []SSHPublicKey{key},
	}

	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}

	feedback.Print(feedback.Debug, false, "%s", string(b))

	err = c.Do("PUT", fmt.Sprintf("/api/users/%s", username), true, payload, &response)
	if err != nil {
		return err
	}

	return nil
}
