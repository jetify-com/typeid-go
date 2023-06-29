package typeid

import (
	"fmt"
	"testing"
)

func ExampleNew() {
	tid := Must(New("prefix"))
	fmt.Printf("New typeid: %s\n", tid)
}

func ExampleNew_withoutPrefix() {
	tid := Must(New(""))
	fmt.Printf("New typeid without prefix: %s\n", tid)
}

func ExampleFromString() {
	tid := Must(FromString("prefix_00041061050r3gg28a1c60t3gf"))
	fmt.Printf("Prefix: %s\nSuffix: %s\n", tid.Type(), tid.suffix)
	// Output:
	// Prefix: prefix
	// Suffix: 00041061050r3gg28a1c60t3gf
}

func TestInvalidPrefix(t *testing.T) {
	testdata := []struct {
		name  string
		input string
	}{
		{"caps", "PREFIX"}, // Would be valid in lowercase
		{"numeric", "12323"},
		{"symbols", "pre.fix"},
		{"spaces", "  "},
		{"long", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz"},
	}

	for _, td := range testdata {
		t.Run(td.name, func(t *testing.T) {
			_, err := New(td.input)
			if err == nil {
				t.Errorf("Expected error for invalid prefix: %s", td.input)
			}
		})
	}
}

func TestInvalidSuffix(t *testing.T) {
	testdata := []struct {
		name  string
		input string
	}{
		{"spaces", "  "},
		{"short", "01234"},
		{"long", "012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789"},
		{"caps", "00041061050R3GG28A1C60T3GF"}, // Would be valid in lowercase
		{"hyphens", "00041061050-3gg28a1-60t3gf"},
		{"crockford_ambiguous", "ooo41o61o5or3gg28a1c6ot3gi"}, // Would be valid if we followed Crocksford's substitution rules
		{"symbols", "00041061050.3gg28a1_60t3gf"},
		{"wrong_alphabet", "ooooooiiiiiiuuuuuuulllllll"},
	}

	for _, td := range testdata {
		t.Run(td.name, func(t *testing.T) {
			_, err := From("prefix", td.input)
			if err == nil {
				t.Errorf("Expected error for invalid suffix: %s", td.input)
			}
		})
	}
}

func TestEncodeDecode(t *testing.T) {
	// Generate a bunch of random typeids, encode and decode from a string
	// and make sure the result is the same as the original.
	for i := 0; i < 1000; i++ {
		tid := Must(New("prefix"))
		decoded, err := FromString(tid.String())
		if err != nil {
			t.Error(err)
		}
		if tid != decoded {
			t.Errorf("Expected %s, got %s", tid, decoded)
		}
	}

	// Repeat with the empty prefix:
	for i := 0; i < 1000; i++ {
		tid := Must(New(""))
		decoded, err := FromString(tid.String())
		if err != nil {
			t.Error(err)
		}
		if tid != decoded {
			t.Errorf("Expected %s, got %s", tid, decoded)
		}
	}
}

func TestSpecialValues(t *testing.T) {
	testdata := []struct {
		name string
		tid  string
		uuid string
	}{
		{"nil", "00000000000000000000000000", "00000000-0000-0000-0000-000000000000"},
		{"one", "00000000000000000000000001", "00000000-0000-0000-0000-000000000001"},
		{"ten", "0000000000000000000000000a", "00000000-0000-0000-0000-00000000000a"},
		{"sixteen", "0000000000000000000000000g", "00000000-0000-0000-0000-000000000010"},
		{"thirty-two", "00000000000000000000000010", "00000000-0000-0000-0000-000000000020"},
	}
	for _, td := range testdata {
		t.Run(td.name, func(t *testing.T) {
			// Values should be equal if we start by parsing the typeid
			tid := Must(FromString(td.tid))
			if td.uuid != tid.UUID() {
				t.Errorf("Expected %s, got %s", td.uuid, tid.UUID())
			}

			// Values should be equal if we start by parsing the uuid
			tid = Must(FromUUID("", td.uuid))
			if td.tid != tid.String() {
				t.Errorf("Expected %s, got %s", td.tid, tid.String())
			}
		})
	}
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Must(New("prefix"))
	}
}

func BenchmarkEncodeDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tid := Must(New("prefix"))
		_ = Must(FromString(tid.String()))
	}
}
