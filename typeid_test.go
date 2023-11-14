package typeid_test

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.jetpack.io/typeid"
	"gopkg.in/yaml.v2"
)

func ExampleWithPrefix() {
	tid := typeid.Must(typeid.WithPrefix("prefix"))
	fmt.Printf("New typeid: %s\n", tid)
}

func ExampleWithPrefix_emptyPrefix() {
	tid := typeid.Must(typeid.WithPrefix(""))
	fmt.Printf("New typeid without prefix: %s\n", tid)
}

func ExampleFromString() {
	tid := typeid.Must(typeid.FromString("prefix_00041061050r3gg28a1c60t3gf"))
	fmt.Printf("Prefix: %s\nSuffix: %s\n", tid.Type(), tid.Suffix())
	// Output:
	// Prefix: prefix
	// Suffix: 00041061050r3gg28a1c60t3gf
}

func TestNilIsEmpty(t *testing.T) {
	var emptyID typeid.AnyID
	nilID, err := typeid.FromString("00000000000000000000000000")
	assert.NoError(t, err)
	assert.Equal(t, nilID, emptyID)
	assert.Equal(t, nilID.String(), emptyID.String())
	assert.Equal(t, nilID.UUID(), emptyID.UUID())
	assert.Equal(t, nilID.UUIDBytes(), emptyID.UUIDBytes())
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
			_, err := typeid.WithPrefix(td.input)
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
			_, err := typeid.From("prefix", td.input)
			if err == nil {
				t.Errorf("Expected error for invalid suffix: %s", td.input)
			}
		})
	}
}

//go:embed testdata/invalid.yml
var invalidYML []byte

type InvalidExample struct {
	Name string `yaml:"name"`
	Tid  string `yaml:"typeid"`
}

func TestInvalidTestdata(t *testing.T) {
	assert.Greater(t, len(invalidYML), 0)
	var testdata []InvalidExample
	err := yaml.Unmarshal(invalidYML, &testdata)
	if err != nil {
		t.Errorf("Failed to unmarshal testdata: %s", err)
	}
	assert.Greater(t, len(testdata), 0)

	for _, td := range testdata {
		t.Run(td.Name, func(t *testing.T) {
			_, err := typeid.FromString(td.Tid)
			if err == nil {
				t.Errorf("Expected error for invalid typeid: %s", td.Tid)
			}
		})
	}
}

func TestEncodeDecode(t *testing.T) {
	// Generate a bunch of random typeids, encode and decode from a string
	// and make sure the result is the same as the original.
	for i := 0; i < 1000; i++ {
		tid := typeid.Must(typeid.WithPrefix("prefix"))
		decoded, err := typeid.FromString(tid.String())
		if err != nil {
			t.Error(err)
		}
		if tid != decoded {
			t.Errorf("Expected %s, got %s", tid, decoded)
		}
	}

	// Repeat with the empty prefix:
	for i := 0; i < 1000; i++ {
		tid := typeid.Must(typeid.WithPrefix(""))
		decoded, err := typeid.FromString(tid.String())
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
	for _, data := range testdata {
		t.Run(data.name, func(t *testing.T) {
			// Values should be equal if we start by parsing the typeid
			tid := typeid.Must(typeid.FromString(data.tid))
			if data.uuid != tid.UUID() {
				t.Errorf("Expected %s, got %s", data.uuid, tid.UUID())
			}

			// Values should be equal if we start by parsing the uuid
			tid = typeid.Must(typeid.FromUUID[typeid.AnyID]("", data.uuid))
			if data.tid != tid.String() {
				t.Errorf("Expected %s, got %s", data.tid, tid.String())
			}
		})
	}
}

//go:embed testdata/valid.yml
var validYML []byte

type ValidExample struct {
	Name   string `yaml:"name"`
	Tid    string `yaml:"typeid"`
	Prefix string `yaml:"prefix"`
	UUID   string `yaml:"uuid"`
}

func TestValidTestdata(t *testing.T) {
	assert.Greater(t, len(validYML), 0)
	var testdata []ValidExample
	err := yaml.Unmarshal(validYML, &testdata)
	if err != nil {
		t.Errorf("Failed to unmarshal testdata: %s", err)
	}
	assert.Greater(t, len(testdata), 0)

	for _, td := range testdata {
		t.Run(td.Name, func(t *testing.T) {
			tid := typeid.Must(typeid.FromString(td.Tid))
			if td.UUID != tid.UUID() {
				t.Errorf("Expected %s, got %s", td.UUID, tid.UUID())
			}
			if td.Prefix != tid.Type() {
				t.Errorf("Expected %s, got %s", td.Prefix, tid.Type())
			}
		})
	}
}
