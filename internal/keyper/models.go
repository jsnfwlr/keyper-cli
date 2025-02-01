package keyper

import (
	"fmt"
	"strings"
	"time"
)

type UserResponse struct {
	AccountLocked bool           `json:"accountLocked,omitempty"`
	CN            string         `json:"cn,omitempty"`
	DisplayName   string         `json:"displayName,omitempty"`
	DN            string         `json:"dn,omitempty"`
	GivenName     string         `json:"givenName,omitempty"`
	Mail          string         `json:"mail,omitempty"`
	Groups        []Group        `json:"memberOfs,omitempty"`
	SN            string         `json:"sn,omitempty"`
	SSHPublicKeys []SSHPublicKey `json:"sshPublicKeys"`
	UID           string         `json:"uid,omitempty"`
}

type DateExpire string

type SSHPublicKey struct {
	CN          string     `json:"cn,omitempty"`
	DateExpire  DateExpire `json:"dateExpire,omitempty"`
	Date        time.Time  `json:"-"`
	HostGroups  []Group    `json:"hostGroups,omitempty"`
	Key         string     `json:"key"`
	KeyId       int        `json:"keyid,omitempty"`
	Fingerprint string     `json:"fingerprint"`
	Name        string     `json:"name"`
	KeyType     int        `json:"keyType,omitempty"`
	Local       bool       `json:"-"`
}

type GroupEntry struct {
	CN string
	OU string
	DC []string
}

type Group string

func (g Group) Parse() GroupEntry {
	ge := GroupEntry{}
	parts := strings.Split(string(g), ",")
	for _, part := range parts {
		switch {
		case strings.HasPrefix(part, "cn="):
			ge.CN = strings.TrimPrefix(part, "cn=")
		case strings.HasPrefix(part, "ou="):
			ge.OU = strings.TrimPrefix(part, "ou=")
		case strings.HasPrefix(part, "dc="):
			ge.DC = append(ge.DC, strings.TrimPrefix(part, "dc="))
		}
	}

	return ge
}

func (ge GroupEntry) String() Group {
	return Group(fmt.Sprintf("cn=%s,ou=%s,dc=%s", ge.CN, ge.OU, strings.Join(ge.DC, ",dc=")))
}
