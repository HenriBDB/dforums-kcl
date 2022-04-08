package security

/*
	This code is largely inspired by:
	https://github.com/umahmood/hashcash/blob/master/hashcash.go
	Distributed under MIT Liscence as of 17 March 2022 (allowing reuse without credit required in any other project, commercial or private)
*/
import (
	"crypto/rand"
	"crypto/sha256"
	"dforum-app/configuration"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	maxIterations  int    = 1 << 30 // Max iterations to find a solution
	bytesToRead    int    = 8       // Bytes to read for random token
	bitsPerHexChar int    = 4       // Each hex character takes 4 bits
	zero           rune   = 48      // ASCII code for number zero
	hashcashLength int    = 4       // Number of items in a hashcash header
	version        string = "DF1"   // Version of hashcash algorithm used
	shaLength      int    = 64      // Size, in number of hex characters, of hashes compared for PoW
	stdDifficulty  int    = 24
)

// hashcash instance
type hashcash struct {
	// difficulty number of "partial pre-image" (zero) difficulty in the hashed code.
	difficulty int
	// salt to use, encoded in base-64 format.
	salt string
	// counter (up to 2^30), encoded in base-64 format.
	counter int
}

// compute a new hashcash header. If no solution can be found 'ErrSolutionFail'
// error is returned.
func (h *hashcash) compute(preImage string) (string, error) {
	// hex char: 0    0    0    0    0
	// binary  : 0000 0000 0000 0000 0000 = 4 bits per char = 20 bits total
	var (
		collisionSize = h.difficulty / bitsPerHexChar
		header        = h.createHeader()
		hash          = sha256Hash(header)
	)
	if len(preImage) != len(hash) {
		return "", ErrInvalidInput
	}
	for !acceptableHeader(hash, preImage, collisionSize) {
		h.counter++
		header = h.createHeader()
		hash = sha256Hash(header)
		if h.counter >= maxIterations {
			return "", ErrSolutionFail
		}
	}
	return header, nil
}

// Verify that a hashcash header is valid. If the header is not in a valid
// format, ErrInvalidHeader error is returned.
func verifyProofOfWork(header string, dataBytes []byte) (bool, error) {
	vals := strings.Split(header, ":")
	if len(vals) != hashcashLength {
		return false, ErrInvalidHeader
	}
	// vals: [version difficulty salt counter]
	difficulty, err := strconv.Atoi(vals[1])
	// TODO User parameterised difficulty
	if err != nil || difficulty < 0 || shaLength < difficulty {
		return false, ErrInvalidDifficulty
	}
	var (
		hash          = sha256Hash(header)
		collisionSize = difficulty / bitsPerHexChar
	)
	// test 1 - zero count
	if !acceptableHeader(hash, sha256HashFromBytes(dataBytes), collisionSize) {
		return false, ErrNoCollision
	}
	return true, nil
}

// New creates a new Hashcash instance
func newProofOfWork(dataBytes []byte) (string, error) {
	if dataBytes == nil {
		return "", ErrInvalidInput
	}
	salt, err := randomBytes(bytesToRead)
	if err != nil {
		return "", err
	}
	diff := configuration.GetMinNodeDifficulty()
	if diff < 16 || diff > 28 {
		diff = stdDifficulty
	}
	hc := hashcash{
		difficulty: diff,
		salt:       base64EncodeBytes(salt),
		counter:    1,
	}
	return hc.compute(sha256HashFromBytes(dataBytes))
}

// acceptableHeader determines if the string 'hash' is prefixed with 'n',
// 'char' characters.
func acceptableHeader(hash string, preImage string, n int) bool {
	preImageRunes := []rune(preImage)
	for i, val := range hash[:n] {
		if val != preImageRunes[i] {
			return false
		}
	}
	return true
}

// createHeader creates a new hashcash header
func (h *hashcash) createHeader() string {
	return fmt.Sprintf("%s:%d:%s:%s",
		version,
		h.difficulty,
		h.salt,
		base64EncodeInt(h.counter))
}

// randomBytes reads n cryptographically secure pseudo-random numbers.
func randomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// base64EncodeBytes
func base64EncodeBytes(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// base64EncodeInt
func base64EncodeInt(n int) string {
	return base64EncodeBytes([]byte(strconv.Itoa(n)))
}

func sha256HashFromBytes(dataBytes []byte) string {
	hash := sha256.New()
	_, err := hash.Write(dataBytes)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// sha256Hash
func sha256Hash(s string) string {
	hash := sha256.New()
	_, err := io.WriteString(hash, s)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}

var (
	// ErrSolutionFail error cannot compute a solution
	ErrSolutionFail = errors.New("exceeded 2^20 iterations failed to find solution")

	// ErrInvalidInput error empty or bad input object
	ErrInvalidInput = errors.New("invalid data object provided")

	// ErrInvalidHeader error invalid hashcash header format
	ErrInvalidHeader = errors.New("invalid hashcash header format")

	// ErrNoCollision error n 5 most significant hex digits (n most significant
	// bits are not 0.
	ErrNoCollision = errors.New("no collision most significant bits are not zero")

	// ErrInvalidDifficulty error avoid PoW wth too low difficulty settings
	ErrInvalidDifficulty = errors.New("difficulty too low or out of bounds")
)
