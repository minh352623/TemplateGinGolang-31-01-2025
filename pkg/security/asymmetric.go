package security

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"time"
)

type SecurityService struct {
	PrivKey string
	AESKey  string
	PubKey  string
}

// Generate RSA key pair
func (a *SecurityService) GenerateRSAKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	return privKey, &privKey.PublicKey, nil
}

// Generate ECDSA key pair
func (a *SecurityService) GenerateECDSAKeys() (*ecdsa.PrivateKey, error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return privKey, nil
}

// GenerateAESKey creates a random 32-byte AES key
func (a *SecurityService) GenerateAESKey() ([]byte, error) {
	key := make([]byte, 32) // 256-bit key
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// Generate Ed25519 key pair
func (a *SecurityService) GenerateEd25519Keys() (ed25519.PrivateKey, ed25519.PublicKey, error) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return privKey, pubKey, nil
}

// Encrypt a message using RSA-OAEP + AES-GCM
func (a *SecurityService) EncryptMessage(msg []byte, recipientPub *rsa.PublicKey) (string, error) {
	aesKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, aesKey); err != nil {
		return "", err
	}

	hash := sha256.New()
	encryptedKey, err := rsa.EncryptOAEP(hash, rand.Reader, recipientPub, aesKey, nil)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nil, nonce, msg, nil)
	result := append(encryptedKey, nonce...)
	result = append(result, ciphertext...)

	return base64.StdEncoding.EncodeToString(result), nil
}

// Decrypt a message
func (a *SecurityService) DecryptMessage(encryptedMsg string, recipientPriv *rsa.PrivateKey) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(encryptedMsg)
	if err != nil {
		return nil, err
	}

	keySize := recipientPriv.PublicKey.Size()
	encryptedKey := decoded[:keySize]
	nonce := decoded[keySize : keySize+12]
	ciphertext := decoded[keySize+12:]

	hash := sha256.New()
	aesKey, err := rsa.DecryptOAEP(hash, rand.Reader, recipientPriv, encryptedKey, nil)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// Sign a message
func (a *SecurityService) SignMessage(msg []byte, senderPriv *rsa.PrivateKey) (string, error) {
	hash := sha256.Sum256(msg)
	signature, err := rsa.SignPKCS1v15(rand.Reader, senderPriv, crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

// Verify a signature
func (a *SecurityService) VerifySignature(msg []byte, signature string, senderPub *rsa.PublicKey) error {
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}
	hash := sha256.Sum256(msg)
	return rsa.VerifyPKCS1v15(senderPub, crypto.SHA256, hash[:], sig)
}

// Sign message with ECDSA
func (a *SecurityService) SignMessageECDSA(msg []byte, privKey *ecdsa.PrivateKey) (string, string, error) {
	hash := sha256.Sum256(msg)
	r, s, err := ecdsa.Sign(rand.Reader, privKey, hash[:])
	if err != nil {
		return "", "", err
	}

	// Encode r and s as base64 strings
	rBase64 := base64.StdEncoding.EncodeToString(r.Bytes())
	sBase64 := base64.StdEncoding.EncodeToString(s.Bytes())

	return rBase64, sBase64, nil
}

// Verify ECDSA signature
func (a *SecurityService) VerifySignatureECDSA(msg []byte, rBase64, sBase64 string, pubKey *ecdsa.PublicKey) bool {
	hash := sha256.Sum256(msg)

	// Decode base64 r, s
	rBytes, err := base64.StdEncoding.DecodeString(rBase64)
	if err != nil {
		fmt.Println("Invalid r value")
		return false
	}
	sBytes, err := base64.StdEncoding.DecodeString(sBase64)
	if err != nil {
		fmt.Println("Invalid s value")
		return false
	}

	r := new(big.Int).SetBytes(rBytes)
	s := new(big.Int).SetBytes(sBytes)

	return ecdsa.Verify(pubKey, hash[:], r, s)
}

// Encrypt and sign a message
func (a *SecurityService) EncryptAndSign(msg []byte, senderPriv *rsa.PrivateKey, recipientPub *rsa.PublicKey) (string, error) {
	signature, err := a.SignMessage(msg, senderPriv)
	if err != nil {
		return "", err
	}

	msgWithSignature := append(msg, []byte(signature)...)
	return a.EncryptMessage(msgWithSignature, recipientPub)
}

// Decrypt and verify a message
func (a *SecurityService) DecryptAndVerify(encryptedMsg string, recipientPriv *rsa.PrivateKey, senderPub *rsa.PublicKey) ([]byte, error) {
	decryptedMsg, err := a.DecryptMessage(encryptedMsg, recipientPriv)
	if err != nil {
		return nil, err
	}

	msgLen := len(decryptedMsg) - 344
	msg := decryptedMsg[:msgLen]
	signature := string(decryptedMsg[msgLen:])

	err = a.VerifySignature(msg, signature, senderPub)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// Encrypt a message using AES-GCM
func (a *SecurityService) EncryptAESGCM(msg []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nil, nonce, msg, nil)
	return append(nonce, ciphertext...), nil
}

// Decrypt a message using AES-GCM
func (a *SecurityService) DecryptAESGCM(encryptedMsg []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(encryptedMsg) < nonceSize {
		return nil, fmt.Errorf("invalid ciphertext")
	}

	nonce, ciphertext := encryptedMsg[:nonceSize], encryptedMsg[nonceSize:]
	return aesGCM.Open(nil, nonce, ciphertext, nil)
}

// Encrypt and sign a message using ECDSA
func (a *SecurityService) EncryptAndSignECDSA(msg []byte, senderPriv *ecdsa.PrivateKey, aesKey []byte) (string, error) {
	rBase64, sBase64, err := a.SignMessageECDSA(msg, senderPriv)
	if err != nil {
		return "", err
	}

	signature := rBase64 + ":" + sBase64
	msgWithSignature := append(msg, []byte("\n"+signature)...)

	encryptedMsg, err := a.EncryptAESGCM(msgWithSignature, aesKey)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encryptedMsg), nil
}

// Decrypt and verify a message using ECDSA
func (a *SecurityService) DecryptAndVerifyECDSA(encryptedMsg string, recipientAESKey []byte, senderPub *ecdsa.PublicKey) ([]byte, error) {
	decodedMsg, err := base64.StdEncoding.DecodeString(encryptedMsg)
	if err != nil {
		return nil, err
	}

	decryptedMsg, err := a.DecryptAESGCM(decodedMsg, recipientAESKey)
	if err != nil {
		return nil, err
	}

	parts := bytes.SplitN(decryptedMsg, []byte("\n"), 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid decrypted format")
	}

	msg := parts[0]
	signatureParts := bytes.SplitN(parts[1], []byte(":"), 2)
	if len(signatureParts) != 2 {
		return nil, fmt.Errorf("invalid signature format")
	}

	rBase64 := string(signatureParts[0])
	sBase64 := string(signatureParts[1])

	if !a.VerifySignatureECDSA(msg, rBase64, sBase64, senderPub) {
		return nil, fmt.Errorf("signature verification failed")
	}

	return msg, nil
}

// Sign message with Ed25519
func (a *SecurityService) SignMessageEd25519(msg []byte, privKey ed25519.PrivateKey) string {
	signature := ed25519.Sign(privKey, msg)
	return base64.StdEncoding.EncodeToString(signature)
}

// Verify Ed25519 signature
func (a *SecurityService) VerifySignatureEd25519(msg []byte, signatureBase64 string, pubKey ed25519.PublicKey) bool {
	signature, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		fmt.Println("Invalid signature format")
		return false
	}
	return ed25519.Verify(pubKey, msg, signature)
}

// Encrypt and sign a message using Ed25519
func (a *SecurityService) EncryptAndSignEd25519(msg []byte, senderPriv interface{}, recipientAESKey interface{}) (string, error) {
	// If no AES key provided, use the one from struct
	if recipientAESKey == nil && a.AESKey != "" {
		recipientAESKey = a.AESKey
	}

	// If no private key provided, use the one from struct
	if senderPriv == nil && a.PrivKey != "" {
		senderPriv = a.PrivKey
	}

	// Handle AES key - can be string or []byte
	var aesKey []byte
	switch k := recipientAESKey.(type) {
	case string:
		decoded, err := base64.StdEncoding.DecodeString(k)
		if err != nil {
			return "", fmt.Errorf("invalid AES key format: %w", err)
		}
		aesKey = decoded
	case []byte:
		aesKey = k
	default:
		return "", fmt.Errorf("AES key must be string or []byte")
	}

	// Handle private key - can be string or ed25519.PrivateKey
	var privKey ed25519.PrivateKey
	switch p := senderPriv.(type) {
	case string:
		decoded, err := base64.StdEncoding.DecodeString(p)
		if err != nil {
			return "", fmt.Errorf("invalid private key format: %w", err)
		}
		privKey = decoded
	case ed25519.PrivateKey:
		privKey = p
	case []byte:
		privKey = ed25519.PrivateKey(p)
	default:
		return "", fmt.Errorf("private key must be string, []byte, or ed25519.PrivateKey")
	}

	// Rest of the encryption logic...
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	msgWithTimestamp := append([]byte(timestamp+"|"), msg...)

	signature := a.SignMessageEd25519(msgWithTimestamp, privKey)
	msgWithSignature := append(msgWithTimestamp, []byte("|"+signature)...)

	encryptedMsg, err := a.EncryptAESGCM(msgWithSignature, aesKey)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encryptedMsg), nil
}

// Decrypt and verify a message using Ed25519
func (a *SecurityService) DecryptAndVerifyEd25519(encryptedMsg string, recipientAESKey interface{}, senderPub interface{}) ([]byte, error) {
	// If no AES key provided, use the one from struct
	if recipientAESKey == nil && a.AESKey != "" {
		recipientAESKey = a.AESKey
	}

	// If no public key provided, use the one from struct
	if senderPub == nil && a.PubKey != "" {
		senderPub = a.PubKey
	}

	// Handle AES key - can be string or []byte
	var aesKey []byte
	switch k := recipientAESKey.(type) {
	case string:
		decoded, err := base64.StdEncoding.DecodeString(k)
		if err != nil {
			return nil, fmt.Errorf("invalid AES key format: %w", err)
		}
		aesKey = decoded
	case []byte:
		aesKey = k
	default:
		return nil, fmt.Errorf("AES key must be string or []byte")
	}

	// Handle public key - can be string or ed25519.PublicKey
	var pubKey ed25519.PublicKey
	switch p := senderPub.(type) {
	case string:
		decoded, err := base64.StdEncoding.DecodeString(p)
		if err != nil {
			return nil, fmt.Errorf("invalid public key format: %w", err)
		}
		pubKey = decoded
	case ed25519.PublicKey:
		pubKey = p
	case []byte:
		pubKey = ed25519.PublicKey(p)
	default:
		return nil, fmt.Errorf("public key must be string, []byte, or ed25519.PublicKey")
	}

	// Rest of the decryption logic...
	decodedMsg, err := base64.StdEncoding.DecodeString(encryptedMsg)
	if err != nil {
		return nil, err
	}

	decryptedMsg, err := a.DecryptAESGCM(decodedMsg, aesKey)
	if err != nil {
		return nil, err
	}

	parts := bytes.SplitN(decryptedMsg, []byte("|"), 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid decrypted format")
	}

	timestamp, msg, signature := string(parts[0]), parts[1], string(parts[2])

	timeSent, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil || time.Now().Unix()-timeSent > 30000 {
		return nil, fmt.Errorf("message expired or invalid timestamp")
	}

	msgWithTimestamp := append([]byte(timestamp+"|"), msg...)
	if !a.VerifySignatureEd25519(msgWithTimestamp, signature, pubKey) {
		return nil, fmt.Errorf("signature verification failed")
	}

	return msg, nil
}
