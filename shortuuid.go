package shortuuid

import (
	"fmt"
	"math"
	"math/big"
	"sort"
	"strings"

	crand "crypto/rand"
	"github.com/satori/go.uuid"
)

type ShortUUID struct {
	alphabet    []string
	alphaLength int
}

func (su *ShortUUID) GetAlphabet() string {
	return strings.Join(su.alphabet, "")
}

func (su *ShortUUID) SetAlphabet(alphabet string) {
	alphabetTable := removeDuplicates(splitStr(alphabet))
	sort.Strings(alphabetTable)
	su.alphabet = alphabetTable
	su.alphaLength = len(alphabetTable)
}

func (su *ShortUUID) GetLength() int {
	return int(math.Ceil(mathLog(math.Pow(2, 128), float64(su.alphaLength))))
}

func (su *ShortUUID) Encode(uuid string, padLength int) string {
	if padLength == 0 {
		padLength = su.GetLength()
	}
	return IntToString(uuidToInt(uuid), su.GetAlphabet(), padLength)
}

func (su *ShortUUID) Decode(str string) uuid.UUID {
	originBytes := StringToInt(str, su.GetAlphabet()).Bytes()
	diff := 16 - len(originBytes)

	var fullBytes []byte
	for diff > 0 {
		fullBytes = append(fullBytes, 0)
		diff--
	}

	if len(fullBytes) != 0 {
		fullBytes = append(fullBytes, originBytes...)
	} else {
		fullBytes = originBytes
	}

	u, err := uuid.FromBytes(fullBytes)
	if err != nil {
		panic(err)
	}
	return u
}

func (su *ShortUUID) GetUUID(name string, padLength int) string {
	u := new(uuid.UUID)
	name = strings.ToLower(name)

	if name == "" {
		uuidV4 := uuid.NewV4()
		u = &uuidV4
	} else if strings.HasPrefix(name, "http://") || strings.HasPrefix(name, "https://") {
		uuidV5 := uuid.NewV5(uuid.NamespaceURL, name)
		u = &uuidV5
	} else {
		uuidV5 := uuid.NewV5(uuid.NamespaceDNS, name)
		u = &uuidV5
	}

	return su.Encode(u.String(), padLength)
}

func (su *ShortUUID) GetUUIDWithNameSpace(name string) string {
	uuidStr := su.GetUUID(name, su.GetLength())
	return uuidStr

}

func (su *ShortUUID) UUID() string {
	return su.GetUUIDWithNameSpace("")
}

func (su *ShortUUID) Random(length int) string {
	if length == 0 {
		length = su.GetLength()
	}

	// read random file
	b := make([]byte, length)
	if _, err := crand.Read(b); err != nil {
		panic(err)
	}

	//hexadecimal
	hexStr := fmt.Sprintf("%x", b)

	i := newBigInt()
	i.SetString(hexStr, 16)
	return IntToString(i, su.GetAlphabet(), length)[:length]

}

func NewShortUUID() *ShortUUID {
	su := new(ShortUUID)
	su.SetAlphabet("23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")
	return su
}

func IntToString(number *big.Int, alphabet string, padding int) string {
	output := ""
	alphabetTable := splitStr(alphabet)
	alphabetLength := big.NewInt(int64(len(alphabet)))

	for number.Cmp(big.NewInt(0)) == 1 {
		quotient, digit := newBigInt().DivMod(number, alphabetLength, newBigInt())
		output += alphabetTable[digit.Int64()]
		number = quotient
	}

	if padding != 0 {
		remainder := int(math.Max(float64(padding-len(output)), 0))
		other := ""
		for remainder > 0 {
			other += alphabetTable[0]
			remainder--
		}
		output = output + other
	}

	return reverseStr(output)
}

func StringToInt(str string, alphabet string) *big.Int {
	number := newBigInt()
	alphabetTable := splitStr(alphabet)
	alphabetLength := big.NewInt(int64(len(alphabet)))

	for _, char := range splitStr(str) {
		number = number.Mul(number, alphabetLength)
		charIndex := big.NewInt(int64(indexOf(char, alphabetTable)))
		number.Add(number, charIndex)
	}

	return number
}

func indexOf(char string, chars []string) int {
	for i, _ := range chars {
		if chars[i] == char {
			return i
		}
	}

	return -1
}

func reverseStr(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func splitStr(s string) []string {
	var ss []string
	for _, r := range s {
		c := string(r)
		ss = append(ss, c)
	}
	return ss
}

func mathLog(x, base float64) float64 {
	if base == 0 {
		base = math.E
	}
	r := math.Log2(x) / math.Log2(base)
	return r
}

func uuidToInt(u string) *big.Int {
	i := newBigInt()
	i.SetString(strings.Replace(u, "-", "", 4), 16)
	return i
}

func newBigInt() *big.Int {
	return new(big.Int)
}

func removeDuplicates(chars []string) []string {
	encountered := map[string]bool{}
	var result []string

	for v := range chars {
		if encountered[chars[v]] == true {
		} else {
			encountered[chars[v]] = true
			result = append(result, chars[v])
		}
	}
	return result
}
