package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
)

func GenerateURL(author_id uuid.UUID, id uuid.UUID) string {
	h := sha256.New()
	ibytes, _ := id.MarshalBinary()
	abytes, _ := author_id.MarshalBinary()

	h.Write(abytes)
	dir := hex.EncodeToString(h.Sum(nil))

	h.Reset()

	h.Write(ibytes)
	file := hex.EncodeToString(h.Sum(nil))

	return fmt.Sprintf("%s/%s", dir, file)
}
