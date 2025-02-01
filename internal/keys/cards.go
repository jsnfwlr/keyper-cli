package keys

import (
	"crypto"
	"fmt"
	"strings"

	"github.com/go-piv/piv-go/v2/piv"
	"github.com/jsnfwlr/keyper-cli/internal/feedback"
	"github.com/jsnfwlr/keyper-cli/internal/prompter"
)

func GetKeyFromSK(algo string, size int) (publicKey *crypto.PublicKey, fault error) {
	cards, err := piv.Cards()
	if err != nil {
		return nil, fmt.Errorf("could not get PIV cards: %w", err)
	}

	feedback.Print(feedback.Extra, false, "Found %d PIV cards", len(cards))

	// Find a YubiKey and open the reader.
	var info string
	var yk *piv.YubiKey
	for _, card := range cards {
		if strings.Contains(strings.ToLower(card), "yubikey") {
			yk, err = piv.Open(card)
			if err != nil {
				return nil, fmt.Errorf("could not open YubiKey: %w", err)
			}
			info = card
			break
		}
	}

	// info, err := yk.KeyInfo()
	// if err != nil {
	// 	return nil, fmt.Errorf("could not get key info: %w", err)
	// }

	var pivAlgo piv.Algorithm

	switch algo {
	case "rsa-sk":
		switch size {
		case 1024:
			pivAlgo = piv.AlgorithmRSA1024
		case 2048:
			pivAlgo = piv.AlgorithmRSA2048
		case 3072:
			pivAlgo = piv.AlgorithmRSA3072
		case 4096:
			pivAlgo = piv.AlgorithmRSA4096
		}
	case "ecdsa-sk":
		switch size {
		case 256:
			pivAlgo = piv.AlgorithmEC256
		case 384:
			pivAlgo = piv.AlgorithmEC384
		}
	case "ed25519-sk":
		pivAlgo = piv.AlgorithmEd25519
	}

	key := piv.Key{
		Algorithm:   pivAlgo,
		PINPolicy:   piv.PINPolicyAlways,
		TouchPolicy: piv.TouchPolicyAlways,
	}

	prompt := prompter.New()
	pin, err := prompt.Password(fmt.Sprintf("Enter the management PIN for %s", info), false)
	if err != nil {
		return nil, fmt.Errorf("could not get management PIN: %w", err)
	}

	m, err := yk.Metadata(pin)
	if err != nil {
		return nil, fmt.Errorf("could not get metadata: %w", err)
	}

	if m.ManagementKey == nil {
		return nil, fmt.Errorf("management key is nil")
	}

	mKey := *m.ManagementKey

	pubKey, err := yk.GenerateKey(mKey, piv.SlotAuthentication, key)
	if err != nil {
		return nil, fmt.Errorf("could not generate key: %w", err)
	}

	return &pubKey, nil
}
