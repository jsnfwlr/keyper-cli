package keys

import (
	"fmt"

	enumFlag "github.com/thediveo/enumflag/v2"
)

type Definitions []Definition

type Definition struct {
	Algorithm TypeId
	short     string
	long      string
	suffix    string
}

type TypeId enumFlag.Flag

const (
	RSA TypeId = iota
	ED25519SK
	ED25519
	ECDSASK
	ECDSA
)

var definitions = Definitions{
	{Algorithm: RSA, short: "RSA", suffix: "rsa", long: "Rivest–Shamir–Adleman"},
	{Algorithm: ED25519SK, short: "Ed25519-SK", suffix: "ed25519-sk", long: "Elliptic Curve Digital Signature Algorithm with Curve25519 and Secure Key"},
	{Algorithm: ED25519, short: "Ed25519", suffix: "ed25519", long: "Elliptic Curve Digital Signature Algorithm with Curve25519"},
	{Algorithm: ECDSASK, short: "ECDSA-SK", suffix: "ecdsa-sk", long: "Elliptic Curve Digital Signature Algorithm with Secure Key"},
	{Algorithm: ECDSA, short: "ECDSA", suffix: "ecdsa", long: "Elliptic Curve Digital Signature Algorithm"},
}

func Types() map[TypeId][]string {
	f := make(map[TypeId][]string)
	for _, def := range definitions {
		f[def.Algorithm] = []string{def.suffix}
	}
	return f
}

func (a TypeId) Flag() string {
	for _, d := range definitions {
		if d.Algorithm == a {
			return d.suffix
		}
	}
	return ""
}

func (a TypeId) Short() string {
	for _, d := range definitions {
		if d.Algorithm == a {
			return d.short
		}
	}
	return ""
}

func (a TypeId) Long() string {
	for _, d := range definitions {
		if d.Algorithm == a {
			return fmt.Sprintf("%s (%s)", d.long, d.short)
		}
	}
	return ""
}

func (d Definition) Flag() string {
	return d.suffix
}

func (d Definition) Suffix() string {
	return d.suffix
}

func (d Definition) Short() string {
	return d.short
}

func (d Definition) Long() string {
	return fmt.Sprintf("%s (%s)", d.long, d.short)
}

func GetDefinitions(flag string) Definition {
	for _, set := range definitions {
		if set.suffix == flag {
			return set
		}
	}
	return Definition{}
}
