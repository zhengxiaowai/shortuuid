package shortuuid

import (
	"sort"
	"strings"
	"testing"

	"github.com/satori/go.uuid"
)

func TestShortUUID_UUID(t *testing.T) {
	su := NewShortUUID()

	u1 := su.UUID()
	u2 := su.GetUUIDWithNameSpace("http://www.example.com/")
	u3 := su.GetUUIDWithNameSpace("HTTP://www.example.com/")
	u4 := su.GetUUIDWithNameSpace("example.com/")
	cases := []string{u1, u2, u3, u4}

	for i, c := range cases {
		if !(20 < len(c) && len(c) < 24) {
			t.Errorf("test generation fail, index= %d; uuid = %s", i, c)
		}
	}

}

func TestShortUUID_Encode(t *testing.T) {
	su := NewShortUUID()

	caseTable := map[string]string{
		"3b1f8b40-222c-4a6e-b77e-779d5a94e21c": "CXc85b4rqinB7s5J52TRYb",
	}

	for i, o := range caseTable {
		r := su.Encode(i, 0)
		if r != o {
			t.Errorf("test encode fail, %s != %s", i, o)
		}
	}
}

func TestShortUUID_Decode(t *testing.T) {
	su := NewShortUUID()

	caseTable := map[string]string{
		"CXc85b4rqinB7s5J52TRYb": "3b1f8b40-222c-4a6e-b77e-779d5a94e21c",
	}

	for i, o := range caseTable {
		r := su.Decode(i)
		if r.String() != o {
			t.Errorf("test encode fail, %s != %s", r.String(), o)
		}
	}
}

func TestShortUUID_Random(t *testing.T) {
	su := NewShortUUID()
	for i := 1; i <= 1000; i++ {
		u := su.Random(0)
		if len(u) != 22 {
			t.Errorf("test random fail, length %d != %d", len(u), 22)
		}
	}

	for i := 1; i <= 100; i++ {
		u := su.Random(i)
		if len(u) != i {
			t.Errorf("test random fail, length %d != %d", len(u), i)
		}
	}

}

func TestShortUUID_SetAlphabet(t *testing.T) {
	alphabet := "01"
	su1 := NewShortUUID()
	su1.SetAlphabet(alphabet)

	if alphabet != su1.GetAlphabet() {
		t.Errorf("test SetAlphabet fail, %s != %s", alphabet, su1.GetAlphabet())
	}

	su1.SetAlphabet("01010101010101")

	if alphabet != su1.GetAlphabet() {
		t.Errorf("test SetAlphabet fail, %s != %s", alphabet, su1.GetAlphabet())
	}

	su1UUID := su1.UUID()
	uuidCharArray := removeDuplicates(splitStr(su1UUID))
	sort.Strings(uuidCharArray)
	if strings.Join(uuidCharArray, "") != "01" {
		t.Errorf("test SetAlphabet fail, %s != %s", su1UUID, "01")
	}

	if !(116 < len(su1.UUID()) && len(su1.UUID()) < 140) {
		t.Errorf("test SetAlphabetfail, 116 < (%d) < 140", len(su1.UUID()))
	}

	su2 := NewShortUUID()
	su2UUID := su2.UUID()
	if !(20 < len(su2.UUID()) && len(su2UUID) < 24) {
		t.Errorf("test SetAlphabet fail, 20 < (%d) < 24", len(su2UUID))
	}

	for i := 1; i <= 100; i++ {
		u := uuid.NewV4()

		uuidStr := u.String()
		encodeUUID := su1.Encode(uuidStr, 0)
		decodeUUID := su1.Decode(encodeUUID).String()
		if uuidStr != decodeUUID {
			t.Errorf("test SetAlphabet fail, %s %s %s", uuidStr, encodeUUID, decodeUUID)
		}
	}

}

func TestShortUUID_Padding(t *testing.T) {
	su := NewShortUUID()

	RandomUID := uuid.NewV4()

	zeroBytes := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	smallestUID, err := uuid.FromBytes(zeroBytes)
	if err != nil {
		t.Fatalf("uuid.FromBytes call err, error = %v", err)
	}

	encodedRandom := su.Encode(RandomUID.String(), 0)
	encodedSmall := su.Encode(smallestUID.String(), 0)
	if len(encodedRandom) != len(encodedSmall) {
		t.Errorf("%d != %d", len(encodedRandom), len(encodedSmall))
	}
}

func TestShortUUID_Padding_Decode(t *testing.T) {
	su := NewShortUUID()

	RandomUID := uuid.NewV4()

	zeroBytes := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	smallestUID, err := uuid.FromBytes(zeroBytes)
	if err != nil {
		t.Fatalf("uuid.FromBytes call err, error = %v", err)
	}

	encodedRandom := su.Encode(RandomUID.String(), 0)
	encodedSmall := su.Encode(smallestUID.String(), 0)

	if su.Decode(encodedSmall).String() != smallestUID.String() {
		t.Errorf("%s != %s", su.Decode(encodedSmall).String(), smallestUID.String())
	}

	if su.Decode(encodedRandom).String() != RandomUID.String() {
		t.Errorf("%s != %s", su.Decode(encodedRandom).String(), RandomUID.String())
	}
}

func TestShortUUID_Padding_Consistency(t *testing.T) {
	su := NewShortUUID()
	total := 1000
	uidLengths := make(map[int]int)

	for i := 0; i < total; i++ {
		RandomUID := uuid.NewV4()

		encodedRandom := su.Encode(RandomUID.String(), 0)
		uidLengths[len(encodedRandom)] += 1
		decodeRandom := su.Decode(encodedRandom).String()

		if RandomUID.String() != decodeRandom {
			t.Errorf("%s != %s", RandomUID.String(), decodeRandom)
		}
	}

	if len(uidLengths) != 1 {
		t.Errorf("%d != %d", len(uidLengths), 1)
	}

	for _, v := range uidLengths {
		if v != total {
			t.Errorf("%d != %d", v, total)
		}
	}
}
