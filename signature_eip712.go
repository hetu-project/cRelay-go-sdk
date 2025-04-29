//go:build !libsecp256k1

package nostr

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// NostrTypedData represents the EIP-712 typed data structure for Nostr events
type NostrTypedData struct {
	Types       apitypes.Types           `json:"types"`
	PrimaryType string                   `json:"primaryType"`
	Domain      apitypes.TypedDataDomain `json:"domain"`
	Message     map[string]interface{}   `json:"message"`
}

// NewNostrTypedData creates a new EIP-712 typed data structure for a Nostr event
func NewNostrTypedData(evt *Event) *NostrTypedData {
	// Define the types for EIP-712
	types := apitypes.Types{
		"EIP712Domain": []apitypes.Type{
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
		},
		"NostrEvent": []apitypes.Type{
			{Name: "id", Type: "string"},
			{Name: "pubkey", Type: "string"},
			{Name: "created_at", Type: "uint64"},
			{Name: "kind", Type: "uint32"},
			{Name: "tags", Type: "string[]"},
			{Name: "content", Type: "string"},
		},
	}

	// Create the domain data
	domain := apitypes.TypedDataDomain{
		Name:    "Nostr",
		Version: "1",
		ChainId: math.NewHexOrDecimal256(1), // Mainnet
	}

	// Create the message data
	message := map[string]interface{}{
		"id":         evt.ID,
		"pubkey":     evt.PubKey,
		"created_at": evt.CreatedAt,
		"kind":       evt.Kind,
		"tags":       evt.Tags,
		"content":    evt.Content,
	}

	return &NostrTypedData{
		Types:       types,
		PrimaryType: "NostrEvent",
		Domain:      domain,
		Message:     message,
	}
}

// CheckSignature checks if the event signature is valid using EIP-712
func (evt Event) CheckSignature_eip712() (bool, error) {
	// Create typed data
	typedData := NewNostrTypedData(&evt)

	// Get the domain separator
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return false, fmt.Errorf("failed to hash domain: %w", err)
	}

	// Get the message hash
	messageHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return false, fmt.Errorf("failed to hash message: %w", err)
	}

	// Get the final hash
	hash := crypto.Keccak256(
		[]byte("\x19\x01"),
		domainSeparator,
		messageHash,
	)

	// Decode signature
	sig, err := hex.DecodeString(evt.Sig)
	if err != nil {
		return false, fmt.Errorf("signature '%s' is invalid hex: %w", evt.Sig, err)
	}

	// Recover public key
	pubKey, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return false, fmt.Errorf("failed to recover public key: %w", err)
	}

	// Get address from public key
	recoveredAddr := crypto.PubkeyToAddress(*pubKey).Hex()
	recoveredAddr = strings.TrimPrefix(recoveredAddr, "0x")

	// Compare addresses
	return (recoveredAddr == evt.PubKey), nil
}

// Sign signs an event using EIP-712
func (evt *Event) Sign_eip712(secretKey string) error {
	// Create typed data
	typedData := NewNostrTypedData(evt)

	// Get the domain separator
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return fmt.Errorf("failed to hash domain: %w", err)
	}

	// Get the message hash
	messageHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return fmt.Errorf("failed to hash message: %w", err)
	}

	// Get the final hash
	hash := crypto.Keccak256(
		[]byte("\x19\x01"),
		domainSeparator,
		messageHash,
	)

	// Convert private key
	s, err := crypto.HexToECDSA(secretKey)
	if err != nil {
		return fmt.Errorf("invalid secret key '%s': %w", secretKey, err)
	}

	// Sign the hash
	sig, err := crypto.Sign(hash, s)
	if err != nil {
		return fmt.Errorf("failed to sign: %w", err)
	}

	// Set the public key
	evt.PubKey = strings.TrimPrefix(crypto.PubkeyToAddress(s.PublicKey).Hex(), "0x")

	// Set the signature
	evt.Sig = hex.EncodeToString(sig)

	// Set the ID (hash of the event)
	evt.ID = hex.EncodeToString(hash)

	return nil
}

// HashStruct hashes a struct according to EIP-712
func (typedData *NostrTypedData) HashStruct(primaryType string, data map[string]interface{}) ([]byte, error) {
	encodedData, err := typedData.EncodeData(primaryType, data, 1)
	if err != nil {
		return nil, err
	}
	return crypto.Keccak256(encodedData), nil
}

// EncodeData encodes the data according to EIP-712
func (typedData *NostrTypedData) EncodeData(primaryType string, data map[string]interface{}, depth int) ([]byte, error) {
	if depth > 2 {
		return nil, fmt.Errorf("depth too high")
	}

	// Get the type definition
	types, ok := typedData.Types[primaryType]
	if !ok {
		return nil, fmt.Errorf("type %s not found", primaryType)
	}

	// Encode the type hash
	typeHash := crypto.Keccak256([]byte(typedData.EncodeType(primaryType)))

	// Encode the data
	var encodedData []byte
	encodedData = append(encodedData, typeHash...)

	for _, field := range types {
		value, ok := data[field.Name]
		if !ok {
			return nil, fmt.Errorf("field %s not found in data", field.Name)
		}

		encodedValue, err := typedData.EncodeValue(field.Type, value, depth+1)
		if err != nil {
			return nil, err
		}
		encodedData = append(encodedData, encodedValue...)
	}

	return encodedData, nil
}

// EncodeType encodes a type according to EIP-712
func (typedData *NostrTypedData) EncodeType(primaryType string) string {
	types := typedData.Types[primaryType]
	var result string
	result += primaryType + "("
	for i, field := range types {
		if i > 0 {
			result += ","
		}
		result += field.Type + " " + field.Name
	}
	result += ")"
	return result
}

// EncodeValue encodes a value according to EIP-712
func (typedData *NostrTypedData) EncodeValue(fieldType string, value interface{}, depth int) ([]byte, error) {
	switch fieldType {
	case "string":
		return crypto.Keccak256([]byte(value.(string))), nil
	case "uint32", "uint64":
		return crypto.Keccak256([]byte(fmt.Sprintf("%d", value))), nil
	case "string[]":
		// For tags, we'll encode them as a JSON string and hash it
		jsonStr, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}
		return crypto.Keccak256(jsonStr), nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", fieldType)
	}
}
