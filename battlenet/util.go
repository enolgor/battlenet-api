package battlenet

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"time"
)

type oAuthStateGen struct {
	clientID []byte
	aesgcm   cipher.AEAD
	nonce    []byte
}

func newOAuthStateGen(clientID, clientSecret string) (*oAuthStateGen, error) {
	client, err := hex.DecodeString(clientID)
	if err != nil {
		return nil, err
	}
	key, err := base64.RawStdEncoding.DecodeString(clientSecret)
	if err != nil {
		return nil, fmt.Errorf("Wrong client secret format: %s", err.Error())
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return &oAuthStateGen{client, aesgcm, nonce}, nil
}

func (oasg *oAuthStateGen) generate() string {
	data := make([]byte, len(oasg.clientID)+4)
	copy(data[:len(oasg.clientID)], oasg.clientID)
	binary.BigEndian.PutUint32(data[len(oasg.clientID):], uint32(time.Now().Nanosecond()))
	return hex.EncodeToString(oasg.aesgcm.Seal(nil, oasg.nonce, data, nil))
}

func (oasg *oAuthStateGen) parse(state string) (bool, error) {
	decoded, err := hex.DecodeString(state)
	if err != nil {
		return false, err
	}
	data, err := oasg.aesgcm.Open(nil, oasg.nonce, decoded, nil)
	if err != nil {
		return false, err
	}
	if bytes.Compare(oasg.clientID, data[:len(oasg.clientID)]) == 0 {
		return true, nil
	}
	return false, nil
}
