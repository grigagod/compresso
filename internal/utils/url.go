package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
)

func GenerateURL(author_id uuid.UUID, id uuid.UUID) (string, error) {
	h := sha256.New()
	ibytes, _ := id.MarshalBinary()
	abytes, _ := author_id.MarshalBinary()

	_, err := h.Write(abytes)
	if err != nil {
		return "", err
	}

	dir := hex.EncodeToString(h.Sum(nil))
	h.Reset()

	_, err = h.Write(ibytes)
	if err != nil {
		return "", err
	}

	file := hex.EncodeToString(h.Sum(nil))

	return fmt.Sprintf("%s/%s", dir, file), nil
}
