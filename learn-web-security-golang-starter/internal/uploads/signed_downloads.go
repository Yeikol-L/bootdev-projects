package uploads

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

const signedDownloadTTL = 5 * time.Minute

func CreateSignedDownloadPath(signingKey [32]byte, fileID int64, now time.Time) string {
	expires := now.Unix() + int64(signedDownloadTTL.Seconds())
	signature := signDownload(signingKey, fileID, expires)
	return fmt.Sprintf("/files/%d/signed-download?expires=%d&signature=%s", fileID, expires, signature)
}

func VerifySignedDownload(key [32]byte, fileID int64, expiresValue, signature string, now time.Time) bool {
	expiresInt, err := strconv.ParseInt(expiresValue, 10, 64)
	if err != nil {
		return false
	}
	expectedSignature := signDownload(key, fileID, expiresInt)
	return subtle.ConstantTimeCompare([]byte(expectedSignature), []byte(signature)) == 1 && expiresInt > now.Unix()
}

func signDownload(key [32]byte, fileID, expires int64) string {
	h := hmac.New(sha256.New, key[:])
	fmt.Fprintf(h, "GET\n/files/%d/signed-download\n%d", fileID, expires)
	result := h.Sum(nil)
	return hex.EncodeToString(result)
}
