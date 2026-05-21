package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"
)

// correlationIDCounter is used to generate sequential IDs within the same millisecond
var correlationIDCounter uint64

// GenerateCorrelationID generates a unique correlation ID for tracing API requests
// Format: <timestamp>-<counter>-<random>
// Example: 1683123456789-001-a1b2c3d4
// This ID will be the same for all API calls within a single Terraform operation (Create/Read/Update/Delete)
func GenerateCorrelationID() string {
	// Get current timestamp in milliseconds
	timestamp := time.Now().UnixMilli()
	
	// Increment counter atomically
	counter := atomic.AddUint64(&correlationIDCounter, 1)
	
	// Generate random bytes for uniqueness
	randomBytes := make([]byte, 4)
	if _, err := rand.Read(randomBytes); err != nil {
		// Fallback to timestamp-based random if crypto/rand fails
		randomBytes = []byte(fmt.Sprintf("%04d", counter%10000))
	}
	randomHex := hex.EncodeToString(randomBytes)
	
	return fmt.Sprintf("%d-%03d-%s", timestamp, counter%1000, randomHex)
}