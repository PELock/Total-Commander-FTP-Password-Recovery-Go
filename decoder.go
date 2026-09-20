/******************************************************************************
 * Total Commander FTP Password Recovery
 *
 * Version        : v1.0.1
 * Language       : Go
 * Author         : Bartosz Wójcik
 * Web page       : https://www.pelock.com
 *
 *****************************************************************************/

package totalcommanderftppasswordrecovery

import (
	"fmt"
	"strings"
	"unicode"
)

// Decoder decrypts Total Commander FTP passwords stored in wcx_ftp.ini.
// Checksum validation present in Total Commander is not implemented here.
type Decoder struct {
	randomSeed uint32
}

// New creates a password decoder.
func New() *Decoder {
	return &Decoder{}
}

// HexStringToByteArray decodes an even-length hex string. Returns nil on invalid input.
func HexStringToByteArray(str string) []byte {
	lowered := strings.ToLower(str)
	n := len(lowered)
	if n == 0 || n&1 != 0 {
		return nil
	}
	result := make([]byte, n/2)
	for i := 0; i < n; i += 2 {
		hi := hexNibble(lowered[i])
		lo := hexNibble(lowered[i+1])
		if hi < 0 || lo < 0 {
			return nil
		}
		result[i/2] = byte((hi << 4) | lo)
	}
	return result
}

// SeedPrng sets the uint32 LCG seed.
func (d *Decoder) SeedPrng(seed int) {
	d.randomSeed = uint32(seed)
}

// NextRandMax advances seed = seed*0x8088405+1 and returns the high 32 bits of seed*nMax.
func (d *Decoder) NextRandMax(nMax int) int {
	d.randomSeed = d.randomSeed*0x8088405 + 1
	return int((uint64(d.randomSeed) * uint64(uint32(nMax))) >> 32)
}

// Rol8 rotates the low 8 bits of value left by counter&7.
func Rol8(value, counter int) int {
	b := value & 0xFF
	c := counter & 7
	if c == 0 {
		return b
	}
	return ((b << c) | (b >> (8 - c))) & 0xFF
}

// DecryptPassword decrypts hex ciphertext and returns UTF-8/raw plaintext bytes.
func (d *Decoder) DecryptPassword(passwordHex string) ([]byte, error) {
	var b strings.Builder
	for _, r := range passwordHex {
		if !unicode.IsSpace(r) {
			b.WriteRune(r)
		}
	}
	buf := HexStringToByteArray(strings.ToLower(b.String()))
	if buf == nil {
		return nil, fmt.Errorf("invalid hex password")
	}

	passwordLength := len(buf)
	if passwordLength <= 4 {
		return nil, fmt.Errorf("password too short")
	}
	passwordLength -= 4
	work := append([]byte(nil), buf...)

	d.SeedPrng(849521)
	for i := 0; i < passwordLength; i++ {
		work[i] = byte(Rol8(int(work[i]), d.NextRandMax(8)))
	}

	d.SeedPrng(12345)
	for i := 0; i < 256; i++ {
		x := d.NextRandMax(passwordLength)
		y := d.NextRandMax(passwordLength)
		work[x], work[y] = work[y], work[x]
	}

	d.SeedPrng(42340)
	for i := 0; i < passwordLength; i++ {
		work[i] ^= byte(d.NextRandMax(256))
	}

	d.SeedPrng(54321)
	for i := 0; i < passwordLength; i++ {
		work[i] = byte((int(work[i]) - d.NextRandMax(256)) & 0xFF)
	}

	out := make([]byte, passwordLength)
	copy(out, work[:passwordLength])
	return out, nil
}

// DecryptPasswordString decrypts hex ciphertext to a Latin-1 string (matches PHP chr() output).
func (d *Decoder) DecryptPasswordString(passwordHex string) (string, error) {
	raw, err := d.DecryptPassword(passwordHex)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func hexNibble(c byte) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'a' && c <= 'f':
		return int(c - 'a' + 10)
	default:
		return -1
	}
}
