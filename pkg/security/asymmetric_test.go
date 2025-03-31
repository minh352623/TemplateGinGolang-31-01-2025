package security

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"
)

var securityService = &SecurityService{}

func TestEncryptionDecryption(t *testing.T) {
	privKey, pubKey, err := securityService.GenerateRSAKeys()
	if err != nil {
		t.Fatalf("Failed to generate keys: %v", err)
	}

	originalMessage := []byte("Hello, secure world!")
	encryptedMsg, err := securityService.EncryptMessage(originalMessage, pubKey)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	decryptedMsg, err := securityService.DecryptMessage(encryptedMsg, privKey)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if !bytes.Equal(originalMessage, decryptedMsg) {
		t.Errorf("Decrypted message does not match original message")
	}
}

func TestSigningVerification(t *testing.T) {
	privKey, pubKey, err := securityService.GenerateRSAKeys()
	if err != nil {
		t.Fatalf("Failed to generate keys: %v", err)
	}

	message := []byte("This is a signed message")
	signature, err := securityService.SignMessage(message, privKey)
	if err != nil {
		t.Fatalf("Signing failed: %v", err)
	}

	err = securityService.VerifySignature(message, signature, pubKey)
	if err != nil {
		t.Errorf("Signature verification failed: %v", err)
	}
}

func TestEncryptAndSignDecryptAndVerify(t *testing.T) {
	senderPriv, senderPub, err := securityService.GenerateRSAKeys()
	if err != nil {
		t.Fatalf("Failed to generate sender keys: %v", err)
	}

	recipientPriv, recipientPub, err := securityService.GenerateRSAKeys()
	if err != nil {
		t.Fatalf("Failed to generate recipient keys: %v", err)
	}

	message := []byte("Confidential data with signature")
	encryptedSignedMsg, err := securityService.EncryptAndSign(message, senderPriv, recipientPub)
	if err != nil {
		t.Fatalf("Encryption and signing failed: %v", err)
	}

	decryptedMsg, err := securityService.DecryptAndVerify(encryptedSignedMsg, recipientPriv, senderPub)
	if err != nil {
		t.Fatalf("Decryption and verification failed: %v", err)
	}

	if !bytes.Equal(message, decryptedMsg) {
		t.Errorf("Decrypted message does not match original message")
	}
}

func TestDecryptWithWrongKey(t *testing.T) {
	_, pubKey1, _ := securityService.GenerateRSAKeys()
	privKey2, _, _ := securityService.GenerateRSAKeys()

	message := []byte("Test message")
	encryptedMsg, _ := securityService.EncryptMessage(message, pubKey1)

	_, err := securityService.DecryptMessage(encryptedMsg, privKey2)
	if err == nil {
		t.Errorf("Decryption with wrong key should fail")
	}
}

func TestVerifySignatureWithWrongKey(t *testing.T) {
	privKey1, _, _ := securityService.GenerateRSAKeys()
	_, pubKey2, _ := securityService.GenerateRSAKeys()

	message := []byte("Signed message")
	signature, _ := securityService.SignMessage(message, privKey1)

	err := securityService.VerifySignature(message, signature, pubKey2)
	if err == nil {
		t.Errorf("Signature verification with wrong key should fail")
	}
}

func TestEncryptAESGCM(t *testing.T) {
	key, err := securityService.GenerateAESKey()
	if err != nil {
		t.Fatalf("Failed to generate AES key: %v", err)
	}
	message := []byte("Hello, AES-GCM encryption!")

	encryptedMsg, err := securityService.EncryptAESGCM(message, key)
	if err != nil {
		t.Fatalf("AES-GCM Encryption failed: %v", err)
	}

	decryptedMsg, err := securityService.DecryptAESGCM(encryptedMsg, key)
	if err != nil {
		t.Fatalf("AES-GCM Decryption failed: %v", err)
	}

	if !bytes.Equal(message, decryptedMsg) {
		t.Errorf("Decrypted message does not match original message")
	}
}

func TestECDSASigningVerification(t *testing.T) {
	privKey, err := securityService.GenerateECDSAKeys()
	if err != nil {
		t.Fatalf("Failed to generate ECDSA keys: %v", err)
	}
	pubKey := &privKey.PublicKey

	message := []byte("ECDSA signed message")
	r, s, err := securityService.SignMessageECDSA(message, privKey)
	if err != nil {
		t.Fatalf("ECDSA signing failed: %v", err)
	}

	valid := securityService.VerifySignatureECDSA(message, r, s, pubKey)
	if !valid {
		t.Errorf("ECDSA signature verification failed")
	}
}

func TestEncryptAndSignECDSA(t *testing.T) {
	privKey, err := securityService.GenerateECDSAKeys()
	if err != nil {
		t.Fatalf("Failed to generate ECDSA keys: %v", err)
	}
	pubKey := &privKey.PublicKey

	key, err := securityService.GenerateAESKey()
	if err != nil {
		t.Fatalf("Failed to generate ECDSA keys: %v", err)
	}
	message := []byte("Confidential data with ECDSA signature")

	encryptedSignedMsg, err := securityService.EncryptAndSignECDSA(message, privKey, key)
	if err != nil {
		t.Fatalf("ECDSA Encryption and signing failed: %v", err)
	}

	decryptedMsg, err := securityService.DecryptAndVerifyECDSA(encryptedSignedMsg, key, pubKey)
	if err != nil {
		t.Fatalf("ECDSA Decryption and verification failed: %v", err)
	}

	if !bytes.Equal(message, decryptedMsg) {
		t.Errorf("Decrypted message does not match original message")
	}
}

func TestEncryptAndSignEd25519(t *testing.T) {

	aesKey, err := securityService.GenerateAESKey()
	if err != nil {
		t.Fatalf("Failed to generate AES key: %v", err)
	}

	privKey, pubKey, err := securityService.GenerateEd25519Keys()
	fmt.Println("privKey", base64.StdEncoding.EncodeToString(privKey))
	fmt.Println("pubKey", base64.StdEncoding.EncodeToString(pubKey))
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 keys: %v", err)
	}

	message := []byte("This is a secret message")

	encryptedMsg, err := securityService.EncryptAndSignEd25519(message, privKey, aesKey)
	if err != nil {
		t.Fatalf("Failed to encrypt and sign: %v", err)
	}

	decryptedMsg, err := securityService.DecryptAndVerifyEd25519(encryptedMsg, aesKey, pubKey)
	if err != nil {
		t.Fatalf("Failed to decrypt and verify: %v", err)
	}

	if !bytes.Equal(message, decryptedMsg) {
		t.Errorf("Decrypted message doesn't match original: got %s, want %s", decryptedMsg, message)
	}
}

func TestEncryptAndSignEd25519Hardcoded(t *testing.T) {
	// 🔐 Hardcoded AES Key (Base64 Decode)
	aesKeyBase64 := "6180tOC6nefupjVb5zKDLRXUn9BQ+kgZqam4MEe3+pU="
	aesKey, err := base64.StdEncoding.DecodeString(aesKeyBase64)
	if err != nil {
		t.Fatalf("Failed to decode AES key: %v", err)
	}

	// 🔑 Hardcoded Ed25519 Private & Public Key (Base64 Decode)
	privKeyBase64 := "YMI25ylK0hRntZc5e0bXXKXmYTmBbrVFF4W4eE+1Le69xVcOsIV0S+U361orE7BTrOCmlfmo2NUYWjcUETU8Eg=="
	pubKeyBase64 := "jxtAWzs4ptkDCP4zFqsDYApqjrkurWMPc4ETnBDB1xY="

	privKeyBytes, _ := base64.StdEncoding.DecodeString(privKeyBase64)
	pubKeyBytes, _ := base64.StdEncoding.DecodeString(pubKeyBase64)

	privKey := ed25519.PrivateKey(privKeyBytes)
	pubKey := ed25519.PublicKey(pubKeyBytes)

	// 📝 JSON Payload
	payload := map[string]interface{}{
		"userID":          "01c06926-a397-433b-821f-a1afb1f90320",
		"currency":        "usdt",
		"Amount":          1.5,
		"RateUsd":         1,
		"transactionType": "deposit",
		"ProviderKey":     "123456abc",
		"Platform":        "web",
		"WebhookUrl":      "http://localhost:8001/trade/webhook",
		"RateCurrency": map[string]interface{}{
			"TON": map[string]interface{}{
				"USD": 2.3,
			},
			"USDT": map[string]interface{}{
				"USD": 1,
			},
			"USD": map[string]interface{}{
				"USD": 1,
			},
		},
	}
	// payload := map[string]interface{}{
	// 	"UserID":          "8a092ead-b85f-4412-81fb-9cc3bf173566",
	// 	"Platform":        "tbk",
	// 	"FromDate":        "2024-09-25T04:12:20.053Z",
	// 	"ToDate":          "2025-03-24T04:12:20.053Z",
	// 	"TransactionType": "all",
	// 	"Page":            1,
	// 	"Limit":           10,
	// }
	// fee request
	// payload := map[string]interface{}{
	// 	"userID":          "01c06926-a397-433b-821f-a1afb1f90320",
	// 	"currency":        "ton",
	// 	"Amount":          1.5,
	// 	"RateUsd":         2.1,
	// 	"ProviderKey":     "123456abc",
	// 	"TransactionType": "deposit",
	// 	"Platform":        "web",
	// 	"RateCurrency": map[string]interface{}{
	// 		"TON": map[string]interface{}{
	// 			"USD": 2.3,
	// 		},
	// 		"USD": map[string]interface{}{
	// 			"USD": 1,
	// 		},
	// 	},
	// }
	// payload := map[string]interface{}{
	// 	"providerKey": "123456abc",
	// 	"userID":      "01c06926-a397-433b-821f-a1afb1f90320",
	// }
	// charge fee request
	// payload := map[string]interface{}{
	// 	"userID":      "01c06926-a397-433b-821f-a1afb1f90320",
	// 	"providerKey": "123456abc",
	// 	"Platform":    "web",
	// 	"webhookUrl":  "http://localhost:8080/trade/webhook",
	// 	"updateWallet": []map[string]interface{}{
	// 		{
	// 			"currency": "ton",
	// 			"amount":   -1.3,
	// 			"rateUsd":  1.2,
	// 		},
	// 		// {
	// 		// 	"currency": "usdt",
	// 		// 	"amount":   -1.2,
	// 		// 	"rateUsd":  1.2,
	// 		// },
	// 	},
	// }

	// user info
	// payload := map[string]interface{}{
	// 	"userID":      "01c06926-a397-433b-821f-a1afb1f90320",
	// 	"providerKey": "123456abc",
	// }

	// get transaction by user id and Platform
	// payload := map[string]interface{}{
	// 	"userID":          "01c06926-a397-433b-821f-a1afb1f90320",
	// 	"Platform":        "web",
	// 	"fromDate":        "2024-01-01",
	// 	"toDate":          "2025-12-01",
	// 	"transactionType": "deposit",
	// 	"page":            1,
	// 	"limit":           10,
	// }
	// 📦 Chuyển payload thành JSON []byte
	message, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// 🔒 Encrypt & Sign
	encryptedMsg, err := securityService.EncryptAndSignEd25519(message, privKey, aesKey)
	fmt.Println("🔐 Encrypted Message:", encryptedMsg)
	if err != nil {
		t.Fatalf("Failed to encrypt and sign: %v", err)
	}

	// 🔓 Decrypt & Verify
	decryptedMsg, err := securityService.DecryptAndVerifyEd25519(encryptedMsg, aesKey, pubKey)
	if err != nil {
		t.Fatalf("Failed to decrypt and verify: %v", err)
	}

	// ✅ Kiểm tra nếu message gốc == message sau khi decrypt
	if !bytes.Equal(message, decryptedMsg) {
		t.Errorf("Decrypted message doesn't match original: got %s, want %s", decryptedMsg, message)
	}

	// 🎯 In kết quả để kiểm tra
	fmt.Println("🔐 Original Message:", string(message))
	fmt.Println("🔓 Decrypted Message:", string(decryptedMsg))
}

func TestEncryptAndSignEd25519_TamperedData(t *testing.T) {

	aesKey, err := securityService.GenerateAESKey()
	if err != nil {
		t.Fatalf("Failed to generate AES key: %v", err)
	}

	privKey, pubKey, err := securityService.GenerateEd25519Keys()
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 keys: %v", err)
	}
	// convert privKey to base64
	privKeyBase64 := base64.StdEncoding.EncodeToString(privKey)
	pubKeyBase64 := base64.StdEncoding.EncodeToString(pubKey)
	fmt.Println("privKey", privKeyBase64)
	fmt.Println("pubKey", pubKeyBase64)
	message := []byte("Original secure messageaaa")

	encryptedSignedMsg, err := securityService.EncryptAndSignEd25519(message, privKey, aesKey)
	if err != nil {
		t.Fatalf("Encryption and signing failed: %v", err)
	}

	decodedMsg, err := base64.StdEncoding.DecodeString(encryptedSignedMsg)
	if err != nil {
		t.Fatalf("Base64 decode failed: %v", err)
	}

	decodedMsg[len(decodedMsg)-1] ^= 0xFF

	tamperedMsg := base64.StdEncoding.EncodeToString(decodedMsg)

	_, err = securityService.DecryptAndVerifyEd25519(tamperedMsg, aesKey, pubKey)
	if err == nil {
		t.Errorf("Tampered data should not be successfully decrypted and verified")
	} else {
		t.Logf("Expected failure: %v", err)
	}
}

func TestDecryptWithWrongKeyEd25519(t *testing.T) {

	privKey, pubKey, _ := securityService.GenerateEd25519Keys()
	wrongAESKey := make([]byte, 32)
	rand.Read(wrongAESKey)
	message := []byte("Test message")

	aesKey := make([]byte, 32)
	rand.Read(aesKey)
	encryptedMsg, _ := securityService.EncryptAndSignEd25519(message, privKey, aesKey)

	_, err := securityService.DecryptAndVerifyEd25519(encryptedMsg, wrongAESKey, pubKey)
	if err == nil {
		t.Errorf("Decryption with wrong key should fail")
	}
}

func BenchmarkEncryptMessageRSA(b *testing.B) {
	privKey, pubKey, _ := securityService.GenerateRSAKeys()
	message := []byte("Benchmark encryption")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encryptedMsg, _ := securityService.EncryptMessage(message, pubKey)
		_ = encryptedMsg
		decryptedMsg, _ := securityService.DecryptMessage(encryptedMsg, privKey)
		_ = decryptedMsg
	}
}

func BenchmarkEncryptAndSignRSA(b *testing.B) {
	senderPriv, senderPub, _ := securityService.GenerateRSAKeys()
	recipientPriv, recipientPub, _ := securityService.GenerateRSAKeys()
	message := []byte("Benchmark encryption and signing")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encryptedSignedMsg, _ := securityService.EncryptAndSign(message, senderPriv, recipientPub)
		decryptedMsg, _ := securityService.DecryptAndVerify(encryptedSignedMsg, recipientPriv, senderPub)
		_ = decryptedMsg
	}
}

func BenchmarkEncryptAESGCM(b *testing.B) {
	securityService := &SecurityService{}
	key, err := securityService.GenerateAESKey()
	if err != nil {
		b.Fatalf("Failed to generate AES key: %v", err)
	}
	message := []byte("Benchmark AES-GCM encryption")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encryptedMsg, _ := securityService.EncryptAESGCM(message, key)
		decryptedMsg, _ := securityService.DecryptAESGCM(encryptedMsg, key)
		_ = decryptedMsg
	}
}

func BenchmarkEncryptAndSignECDSA(b *testing.B) {
	privKey, _ := securityService.GenerateECDSAKeys()
	pubKey := &privKey.PublicKey
	key, err := securityService.GenerateAESKey()
	if err != nil {
		b.Fatalf("Failed to generate AES key: %v", err)
	}
	message := []byte("Benchmark encryption and signing with ECDSA")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encryptedSignedMsg, _ := securityService.EncryptAndSignECDSA(message, privKey, key)
		decryptedMsg, _ := securityService.DecryptAndVerifyECDSA(encryptedSignedMsg, key, pubKey)
		_ = decryptedMsg
	}
}

func BenchmarkEncryptAndSignEd25519(b *testing.B) {

	privKey, pubKey, err := securityService.GenerateEd25519Keys()
	if err != nil {
		b.Fatalf("Failed to generate Ed25519 keys: %v", err)
	}

	key, err := securityService.GenerateAESKey()
	if err != nil {
		b.Fatalf("Failed to generate AES key: %v", err)
	}

	message := []byte("Benchmark encryption and signing with Ed25519")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encryptedSignedMsg, _ := securityService.EncryptAndSignEd25519(message, privKey, key)
		decryptedMsg, _ := securityService.DecryptAndVerifyEd25519(encryptedSignedMsg, key, pubKey)
		_ = decryptedMsg
	}
}
