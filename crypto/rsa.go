package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"

	"gopkg.in/jwt.v1/util"
)

type SigningRSA struct {
	Signing
}

func init() {
	util.Debug("SigningRSA init")
	signingMethodRS256 := &SigningRSA{Signing{"RS256", crypto.SHA256}}
	signingMethodRS384 := &SigningRSA{Signing{"RS384", crypto.SHA384}}
	signingMethodRS512 := &SigningRSA{Signing{"RS512", crypto.SHA512}}
	AddSigningFunc("RS256", signingMethodRS256)
	AddSigningFunc("RS384", signingMethodRS384)
	AddSigningFunc("RS512", signingMethodRS512)
}

func (s *SigningRSA) Verify(data, sign string, key interface{}) error { // Returns nil if signature is valid
	// Verify the key is the right type
	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return ErrorInvalidKey
	}
	fmt.Println(data, sign)

	st, e := base64.RawURLEncoding.DecodeString(sign)
	if e != nil {
		return e
	}

	if !s.Hash.Available() {
		return ErrorHashUnavailable
	}
	hashed := s.Hash.New()
	hashed.Write([]byte(data))

	return rsa.VerifyPKCS1v15(rsaKey, s.Hash, hashed.Sum(nil), st)
}

func (s *SigningRSA) Sign(data string, key interface{}) (string, error) { // Returns encoded signature or error
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return "", ErrorInvalidKey
	}
	hashed := s.Hash.New()
	hashed.Write([]byte(data))

	sigBytes, err := rsa.SignPKCS1v15(rand.Reader, rsaKey, s.Hash, hashed.Sum(nil))
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(sigBytes), nil
}

func (s *SigningRSA) Alg() string { // returns the alg identifier for this method (example: 'HS256')
	return s.Name
}

func (s *SigningRSA) HashType() crypto.Hash {
	return s.Hash
}
