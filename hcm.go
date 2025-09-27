package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

// BuildHCM builds the Hash Challenge Metadata structure based on the WebAssembly analysis
func BuildHCM(challenge, publicSalt, nonce string, difficulty int, wasmResult string) map[string]interface{} {
	baseSeed := challenge + publicSalt + nonce

	// Generate the core structure based on WebAssembly patterns
	res := map[string]interface{}{}

	// Generate 5-6 random keys with different lengths (observed pattern)
	randomKeys := generateRandomKeys(baseSeed)
	hashValues := generateHashValues(publicSalt, challenge, nonce, wasmResult)

	// Add the random key-value pairs (first 5 are hash values)
	for i := 0; i < 5; i++ {
		if hashValues[i] != "" { // Only add non-empty values
			res[randomKeys[i]] = hashValues[i]
		}
	}

	// The 6th key typically holds difficulty info
	res[randomKeys[5]] = fmt.Sprintf("%d-chr", difficulty)

	// Special underscore keys as documented in README
	res["_"] = generateValidationHash(challenge, publicSalt, nonce, wasmResult, "primary")
	res["_2"] = generateValidationHash(challenge, publicSalt, nonce, wasmResult, "secondary")

	res["_s"] = wasmResult

	// Add __meta object as documented: 5x "64" and 1x "5"
	res["__meta"] = generateMetaObject(nonce, publicSalt)

	return res
}

// generateRandomKeys generates random-looking keys for HCM structure
func generateRandomKeys(baseSeed string) []string {
	keys := make([]string, 6)
	lengths := []int{10, 8, 8, 9, 10, 9}

	for i := 0; i < 6; i++ {
		var keyInput string
		switch i {
		case 0:
			keyInput = baseSeed + "key_alpha"
		case 1:
			keyInput = baseSeed + "key_beta"
		case 2:
			keyInput = baseSeed + "key_gamma"
		case 3:
			keyInput = baseSeed + "key_delta"
		case 4:
			keyInput = baseSeed + "key_epsilon"
		default:
			keyInput = baseSeed + fmt.Sprintf("key_%d", i)
		}

		hash := sha256.Sum256([]byte(keyInput))
		encoded := base64.StdEncoding.EncodeToString(hash[:8])
		cleaned := cleanForAlphanumeric(encoded, lengths[i])
		keys[i] = cleaned
	}
	return keys
}

// generateHashValues generates the hash values for HCM structure
func generateHashValues(publicSalt, challenge, nonce, wasmResult string) []string {
	values := make([]string, 5)

	// val1 - SHA256(wasmResult + challenge)
	hash1 := sha256.Sum256([]byte(wasmResult + challenge))
	values[0] = hex.EncodeToString(hash1[:])

	// val2 - SHA256(publicSalt + challenge)
	hash2 := sha256.Sum256([]byte(publicSalt + challenge))
	values[1] = hex.EncodeToString(hash2[:])

	// val3 - SHA256(nonce + challenge)
	hash3 := sha256.Sum256([]byte(nonce + challenge))
	values[2] = hex.EncodeToString(hash3[:])

	// val4 - SHA256(nonce + publicSalt)
	hash4 := sha256.Sum256([]byte(nonce + publicSalt))
	values[3] = hex.EncodeToString(hash4[:])

	// val5 - SHA256(wasmResult + challenge + publicSalt + nonce)
	hash5 := sha256.Sum256([]byte(wasmResult + challenge + publicSalt + nonce))
	values[4] = hex.EncodeToString(hash5[:])

	return values
}

// generateValidationHash generates validation hashes for special keys
func generateValidationHash(challenge, publicSalt, nonce, wasmResult, hashType string) string {
	return ""
}

// generateMetaObject generates __meta object with hex keys and integer values
func generateMetaObject(nonce string, salt string) map[string]int {
	metaMap := make(map[string]int)
	metaValues := []int{64, 64, 64, 64, 64, 5}

	for i := 0; i < 6; i++ {
		keyInput := fmt.Sprintf("%s:%s:%d", salt, nonce, i)
		hash := sha256.Sum256([]byte(keyInput))
		hexKey := hex.EncodeToString(hash[:])[:15]
		metaMap[hexKey] = metaValues[i]
	}
	return metaMap
}

// cleanForAlphanumeric cleans base64 string to alphanumeric format
func cleanForAlphanumeric(encoded string, targetLength int) string {
	cleaned := strings.ReplaceAll(encoded, "=", "")
	cleaned = strings.ReplaceAll(cleaned, "/", "")
	cleaned = strings.ReplaceAll(cleaned, "+", "")

	if len(cleaned) < targetLength {
		hash := sha256.Sum256([]byte(encoded))
		extra := hex.EncodeToString(hash[:])
		cleaned += extra
	}

	if len(cleaned) > targetLength {
		cleaned = cleaned[:targetLength]
	}

	return cleaned
}
