package security

import (
	"crypto/ecdh"
	"crypto/rand"
	"fmt"
)

// KeyPair represents an X25519 ephemeral key pair
type KeyPair struct {
	PrivKey *ecdh.PrivateKey
	PubKey  *ecdh.PublicKey
}

// GenerateKeyPair generates a new X25519 key pair using crypto/rand
func GenerateKeyPair() (*KeyPair, error) {
	curve := ecdh.X25519()
	priv, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate X25519 private key: %w", err)
	}

	return &KeyPair{
		PrivKey: priv,
		PubKey:  priv.PublicKey(),
	}, nil
}

// PubKeyBytes returns the 32-byte raw public key
func (kp *KeyPair) PubKeyBytes() []byte {
	return kp.PubKey.Bytes()
}

// ComputeSharedSecret derives the raw ECDH shared secret given a local private key and remote public key bytes
func ComputeSharedSecret(localPriv *ecdh.PrivateKey, peerPubBytes []byte) ([]byte, error) {
	curve := ecdh.X25519()
	peerPub, err := curve.NewPublicKey(peerPubBytes)
	if err != nil {
		return nil, fmt.Errorf("invalid peer public key bytes: %w", err)
	}

	sharedSecret, err := localPriv.ECDH(peerPub)
	if err != nil {
		return nil, fmt.Errorf("ECDH key exchange failed: %w", err)
	}

	return sharedSecret, nil
}
