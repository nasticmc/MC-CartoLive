package meshcore

import (
	"crypto/sha256"
	"encoding/hex"
)

func PacketHash(raw []byte) string {
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])[:16]
}
