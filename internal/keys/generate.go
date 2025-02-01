package keys

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

func Generate(ctx context.Context, filePath, keyType, comment, passphrase string, size int) (authorizedKey, fingerprint string, fault error) {
	var pKey any
	var sKey any
	var err error

	justPub := false

	switch keyType {
	case "ed25519":
		pKey, sKey, err = ed25519.GenerateKey(nil)
		if err != nil {
			return "", "", fmt.Errorf("could not generate ED25519 key: %w", err)
		}
	case "ed25519-sk":
		pKey, err = GetKeyFromSK(keyType, size)
		if err != nil {
			return "", "", fmt.Errorf("could not get ED25519-SK key: %w", err)
		}

		justPub = true

	case "rsa":
		var secretKey *rsa.PrivateKey
		secretKey, err = rsa.GenerateKey(rand.Reader, size)
		if err != nil {
			return "", "", fmt.Errorf("could not generate ESA key: %w", err)
		}

		pKey = secretKey.Public()
		sKey = secretKey
	case "rsa-sk":
		pKey, err = GetKeyFromSK(keyType, size)
		if err != nil {
			return "", "", fmt.Errorf("could not get RSA%d-SK key: %w", size, err)
		}

		justPub = true
	case "ecdsa":
		var secretKey *ecdsa.PrivateKey
		var ec elliptic.Curve
		switch size {
		case 256:
			ec = elliptic.P256()
		case 384:
			ec = elliptic.P384()
		case 521:
			ec = elliptic.P521()
		default:
			return "", "", fmt.Errorf("unsupported ECDSA key size: %d", size)
		}

		secretKey, err = ecdsa.GenerateKey(ec, rand.Reader)
		if err != nil {
			return "", "", fmt.Errorf("could not generate ECDSA key: %w", err)
		}

		pKey = secretKey.Public()
		sKey = secretKey

	case "ecdsa-sk":
		pKey, err = GetKeyFromSK(keyType, size)
		if err != nil {
			return "", "", fmt.Errorf("could not get ECDSA%d-SK key: %w", size, err)
		}

		justPub = true

	default:
		return "", "", fmt.Errorf("unsupported key type: %s", keyType)
	}
	var sshAuthKey string
	var sha256FingerPrint string

	if justPub {
		sshAuthKey, sha256FingerPrint, err = SaveAsSSHPubKey(ctx, pKey, filePath, comment)
		if err != nil {
			return "", "", fmt.Errorf("could not generate %s key: %w", keyType, err)
		}
	} else {
		sshAuthKey, sha256FingerPrint, err = SaveAsSSHKeyPair(ctx, pKey, sKey, filePath, comment, passphrase)
		if err != nil {
			return "", "", fmt.Errorf("could not generate %s key: %w", keyType, err)
		}
	}

	return sshAuthKey, sha256FingerPrint, nil
}

func SaveAsSSHPubKey(ctx context.Context, pkey any, filePath, comment string) (sshAuthKey, fingerprint string, fault error) {
	pubSSH, err := ssh.NewPublicKey(pkey)
	if err != nil {
		return "", "", fmt.Errorf("could not create public key: %w", err)
	}

	authorizedKey := ssh.MarshalAuthorizedKey(pubSSH)

	ak := fmt.Sprintf("%s %s\n", string(bytes.TrimSpace(authorizedKey)), comment)

	err = os.WriteFile(filePath, []byte(ak), 0o644)
	if err != nil {
		return "", "", fmt.Errorf("could not write to %s: %w", filePath, err)
	}

	fp := ssh.FingerprintSHA256(pubSSH)

	return ak, fp, nil
}

func SaveAsSSHKeyPair(ctx context.Context, pKey, sKey any, filePath, comment, passphrase string) (sshAuthKey, fingerprint string, fault error) {
	pubSSH, err := ssh.NewPublicKey(pKey)
	if err != nil {
		return "", "", fmt.Errorf("could not create public key: %w", err)
	}

	var pemKey *pem.Block

	if passphrase == "" {
		pemKey, err = ssh.MarshalPrivateKey(sKey, comment)
		if err != nil {
			return "", "", fmt.Errorf("could not marshal private key: %w", err)
		}

	} else {
		pemKey, err = ssh.MarshalPrivateKeyWithPassphrase(sKey, comment, []byte(passphrase))
		if err != nil {
			return "", "", fmt.Errorf("could not marshal private key: %w", err)
		}
	}

	privSSH := pem.EncodeToMemory(pemKey)
	authorizedKey := ssh.MarshalAuthorizedKey(pubSSH)

	err = os.WriteFile(filePath, privSSH, 0o600)
	if err != nil {
		return "", "", fmt.Errorf("could not write to %s: %w", filePath, err)
	}

	ak := fmt.Sprintf("%s %s\n", string(bytes.TrimSpace(authorizedKey)), comment)

	err = os.WriteFile(fmt.Sprintf("%s.pub", filePath), []byte(ak), 0o644)
	if err != nil {
		return "", "", fmt.Errorf("could not write to %s: %w", fmt.Sprintf("%s.pub", filePath), err)
	}

	fp := ssh.FingerprintSHA256(pubSSH)

	return ak, fp, nil
}
