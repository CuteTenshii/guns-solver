package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

// BuildHCM builds the Hash Challenge Metadata structure based on the WebAssembly analysis
func BuildHCM(challenge, publicSalt, nonce string, difficulty int) map[string]interface{} {
	// From the WebAssembly analysis, the key fields are:
	// 1. Random alphanumeric keys with hash values
	// 2. Special keys: _, _2, _s
	// 3. Meta object with hex keys and values

	baseSeed := challenge + publicSalt + nonce

	// Generate the core structure based on WebAssembly patterns
	res := map[string]interface{}{}

	// Generate 5-6 random keys with different lengths (observed pattern)
	randomKeys := generateRandomKeys(baseSeed, 6)
	hashValues := generateHashValues(publicSalt, challenge, 5)

	// Add the random key-value pairs (first 5 are hash values)
	for i := 0; i < 5; i++ {
		res[randomKeys[i]] = hashValues[i]
	}

	// The 6th key typically holds difficulty info
	res[randomKeys[5]] = fmt.Sprintf("%d-chr", difficulty)

	// Special underscore keys as documented in README
	res["_"] = generateValidationHash(challenge, publicSalt, nonce, "primary")
	res["_2"] = generateValidationHash(challenge, publicSalt, nonce, "secondary")

	// Generate the 5-char hex WASM result like in the payload
	wasmResult := generateWasmResultHex(challenge, publicSalt, nonce)
	res["_s"] = wasmResult

	// Add __meta object as documented: 5x "64" and 1x "5"
	res["__meta"] = generateMetaObject(baseSeed, nonce, publicSalt)

	return res
}

// generateWasmResultHex creates a 5-character hex result for HCM
// TODO: Need to reverse-engineer the exact algorithm from WebAssembly
func generateWasmResultHex(challenge, publicSalt, nonce string) string {
	// The real algorithm might involve different inputs or processing
	hash := sha256.Sum256([]byte(challenge + publicSalt + nonce + "wasm_result"))
	return hex.EncodeToString(hash[:])[:5]
}

// Generate random-looking keys based on hash derivatives
// Based on WebAssembly analysis: keys appear to be base64-like alphanumeric strings
func generateRandomKeys(baseSeed string, count int) []string {
	keys := make([]string, count)
	lengths := []int{10, 8, 8, 9, 10, 9} // Observed pattern lengths

	for i := 0; i < count; i++ {
		// Use different hash inputs to create more realistic variation
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
		// Convert to base64 and clean it up to match the pattern
		encoded := base64.StdEncoding.EncodeToString(hash[:8])
		cleaned := cleanForAlphanumeric(encoded, lengths[i])
		keys[i] = cleaned
	}
	return keys
}

// Generate hash values for the random keys
// These should be 64-character SHA-256 hashes, not MD5
func generateHashValues(publicSalt, challenge string, count int) []string {
	values := make([]string, count)
	for i := 0; i < count; i++ {
		var data string
		if i == 0 {
			// CONFIRMED: First one is public salt + challenge (user is SURE)
			data = publicSalt + challenge
		} else {
			// Others are variations - need to reverse engineer these
			data = fmt.Sprintf("%s%s%d", challenge, publicSalt, i)
		}
		hash := sha256.Sum256([]byte(data))
		values[i] = hex.EncodeToString(hash[:]) // Full 64-char SHA-256
	}
	return values
}

// Generate validation hash for special keys (_ and _2)
// TODO: Need to reverse-engineer the exact algorithm used by WebAssembly
func generateValidationHash(challenge, publicSalt, nonce string, hashType string) string {
	// This needs to be the actual algorithm that generates the expected values
	// For now using a placeholder - the real algorithm might involve:
	// - Different input ordering
	// - Additional secret/salt data
	// - Multiple hash rounds
	// - Different encoding (hex vs base64 vs binary)

	var input string
	switch hashType {
	case "primary":
		// Try different combinations until we find the right one
		input = challenge + publicSalt + nonce + "primary"
	case "secondary":
		input = challenge + publicSalt + nonce + "secondary"
	default:
		input = challenge + publicSalt + nonce + hashType
	}

	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

// Generate meta object with hex keys and integer values
// Based on WebAssembly analysis: __meta contains 6 hex keys (5×64, 1×5)
// Keys appear to be derived from hash operations in the PoW solver
func generateMetaObject(baseSeed string, nonce string, salt string) map[string]int {
	metaMap := make(map[string]int)
	metaValues := []int{64, 64, 64, 64, 64, 5} // Pattern from analysis

	// Generate hex keys using similar inputs to what WebAssembly might use
	for i := 0; i < 6; i++ {
		var keyInput string
		if i < 5 {
			// Standard meta keys (value 64) - use nonce + index
			keyInput = fmt.Sprintf("%s:%s:%d", salt, nonce, i)
		} else {
			// Special key (value 5) - might be from different source
			keyInput = fmt.Sprintf("%s:%s:special", baseSeed, salt)
		}

		hash := sha256.Sum256([]byte(keyInput))
		hexKey := hex.EncodeToString(hash[:])[:15] // 15 hex chars
		metaMap[hexKey] = metaValues[i]
	}
	return metaMap
}

// Clean base64-like string to alphanumeric matching observed patterns
func cleanForAlphanumeric(encoded string, targetLength int) string {
	// Remove base64 padding and problematic chars
	cleaned := strings.ReplaceAll(encoded, "=", "")
	cleaned = strings.ReplaceAll(cleaned, "/", "")
	cleaned = strings.ReplaceAll(cleaned, "+", "")

	// Ensure we have enough chars, pad with hash if needed
	if len(cleaned) < targetLength {
		// Extend with more characters if needed
		hash := sha256.Sum256([]byte(encoded))
		extra := hex.EncodeToString(hash[:])
		cleaned += extra
	}

	if len(cleaned) > targetLength {
		cleaned = cleaned[:targetLength]
	}

	return cleaned
}
