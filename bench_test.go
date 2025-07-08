//nolint:all
package typeid_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"go.jetify.com/typeid"
)

// Benchmark suite for TypeID operations.
//
// These benchmarks measure performance across realistic usage patterns:
// - Different prefix lengths (empty, short, medium, long)
// - Various operations (creation, string conversion, parsing)
// - Comparison with raw UUID operations
//
// All benchmarks use sink variables to prevent compiler optimizations
// and cycle through varied test data for representative results.

// Sink variables to prevent compiler optimizations
var (
	sinkTypeID typeid.TypeID
	sinkString string
	sinkBytes  []byte
	sinkError  error
	sinkUUID   uuid.UUID
)

// Test data patterns for varied input benchmarks
var prefixPatterns = []string{
	"",             // empty prefix (untyped)
	"u",            // single character
	"usr",          // short prefix
	"user",         // common prefix
	"account",      // medium prefix
	"customer",     // medium prefix
	"organization", // longer prefix
}

// Pre-generated UUIDs for consistent benchmarking
var testUUIDs = func() []uuid.UUID {
	uuids := make([]uuid.UUID, 100)
	for i := range uuids {
		uuids[i] = uuid.Must(uuid.NewV7())
	}
	return uuids
}()

// Pre-generated TypeIDs with various prefixes
var testTypeIDs = func() []typeid.TypeID {
	var ids []typeid.TypeID
	for _, prefix := range prefixPatterns {
		for i := 0; i < 10; i++ {
			ids = append(ids, typeid.MustGenerate(prefix))
		}
	}
	return ids
}()

// Pre-generated TypeID strings for parsing benchmarks
var testTypeIDStrings = func() []string {
	strings := make([]string, len(testTypeIDs))
	for i, id := range testTypeIDs {
		strings[i] = id.String()
	}
	return strings
}()

// BenchmarkNew measures TypeID creation performance
func BenchmarkNew(b *testing.B) {
	b.Run("prefix=empty", func(b *testing.B) {
		b.ReportAllocs()
		var tid typeid.TypeID
		var err error

		for b.Loop() {
			tid, err = typeid.Generate("")
		}

		sinkTypeID = tid
		sinkError = err
	})

	b.Run("prefix=short", func(b *testing.B) {
		b.ReportAllocs()
		var tid typeid.TypeID
		var err error

		for b.Loop() {
			// Cycle through short prefixes
			prefix := prefixPatterns[1+(b.N%3)]
			tid, err = typeid.Generate(prefix)
		}

		sinkTypeID = tid
		sinkError = err
	})

	b.Run("prefix=medium", func(b *testing.B) {
		b.ReportAllocs()
		var tid typeid.TypeID
		var err error

		for b.Loop() {
			// Cycle through medium prefixes
			prefix := prefixPatterns[4+(b.N%3)]
			tid, err = typeid.Generate(prefix)
		}

		sinkTypeID = tid
		sinkError = err
	})

	// Compare with raw UUID generation
	b.Run("uuid=v7", func(b *testing.B) {
		b.ReportAllocs()
		var uid uuid.UUID
		var err error

		for b.Loop() {
			uid, err = uuid.NewV7()
		}

		sinkUUID = uid
		sinkError = err
	})
}

// BenchmarkString measures string conversion performance
func BenchmarkString(b *testing.B) {
	b.Run("typeid", func(b *testing.B) {
		b.ReportAllocs()
		var s string

		for b.Loop() {
			tid := testTypeIDs[b.N%len(testTypeIDs)]
			s = tid.String()
		}

		sinkString = s
	})

	b.Run("uuid", func(b *testing.B) {
		b.ReportAllocs()
		var s string

		for b.Loop() {
			uid := testUUIDs[b.N%len(testUUIDs)]
			s = uid.String()
		}

		sinkString = s
	})
}

// BenchmarkParse measures parsing performance
func BenchmarkParse(b *testing.B) {
	b.Run("typeid", func(b *testing.B) {
		b.ReportAllocs()
		var tid typeid.TypeID
		var err error

		for b.Loop() {
			s := testTypeIDStrings[b.N%len(testTypeIDStrings)]
			tid, err = typeid.Parse(s)
		}

		sinkTypeID = tid
		sinkError = err
	})

	// Pre-generated UUID strings
	uuidStrings := make([]string, len(testUUIDs))
	for i, u := range testUUIDs {
		uuidStrings[i] = u.String()
	}

	b.Run("uuid", func(b *testing.B) {
		b.ReportAllocs()
		var uid uuid.UUID
		var err error

		for b.Loop() {
			s := uuidStrings[b.N%len(uuidStrings)]
			uid, err = uuid.FromString(s)
		}

		sinkUUID = uid
		sinkError = err
	})
}

// BenchmarkPrefix measures prefix extraction performance
func BenchmarkPrefix(b *testing.B) {
	// Group TypeIDs by prefix length for targeted benchmarking
	var emptyPrefix, shortPrefix, mediumPrefix []typeid.TypeID
	for _, tid := range testTypeIDs {
		switch len(tid.Prefix()) {
		case 0:
			emptyPrefix = append(emptyPrefix, tid)
		case 1, 2, 3, 4:
			shortPrefix = append(shortPrefix, tid)
		default:
			mediumPrefix = append(mediumPrefix, tid)
		}
	}

	b.Run("empty", func(b *testing.B) {
		b.ReportAllocs()
		var s string

		for b.Loop() {
			tid := emptyPrefix[b.N%len(emptyPrefix)]
			s = tid.Prefix()
		}

		sinkString = s
	})

	b.Run("short", func(b *testing.B) {
		b.ReportAllocs()
		var s string

		for b.Loop() {
			tid := shortPrefix[b.N%len(shortPrefix)]
			s = tid.Prefix()
		}

		sinkString = s
	})

	b.Run("medium", func(b *testing.B) {
		b.ReportAllocs()
		var s string

		for b.Loop() {
			tid := mediumPrefix[b.N%len(mediumPrefix)]
			s = tid.Prefix()
		}

		sinkString = s
	})
}

// BenchmarkSuffix measures suffix extraction performance
func BenchmarkSuffix(b *testing.B) {
	b.ReportAllocs()
	var s string

	for b.Loop() {
		tid := testTypeIDs[b.N%len(testTypeIDs)]
		s = tid.Suffix()
	}

	sinkString = s
}

// BenchmarkUUIDBytes measures UUID bytes extraction performance
func BenchmarkUUIDBytes(b *testing.B) {
	b.Run("typeid", func(b *testing.B) {
		b.ReportAllocs()
		var bytes []byte

		for b.Loop() {
			tid := testTypeIDs[b.N%len(testTypeIDs)]
			bytes = tid.Bytes()
		}

		sinkBytes = bytes
	})

	b.Run("uuid", func(b *testing.B) {
		b.ReportAllocs()
		var bytes []byte

		for b.Loop() {
			uid := testUUIDs[b.N%len(testUUIDs)]
			bytes = uid.Bytes()
		}

		sinkBytes = bytes
	})
}

// BenchmarkFromUUID measures TypeID creation from UUID
func BenchmarkFromUUID(b *testing.B) {
	uuidStrings := make([]string, len(testUUIDs))
	for i, u := range testUUIDs {
		uuidStrings[i] = u.String()
	}

	b.ReportAllocs()
	var tid typeid.TypeID
	var err error

	for b.Loop() {
		prefix := prefixPatterns[b.N%len(prefixPatterns)]
		uuidStr := uuidStrings[b.N%len(uuidStrings)]

		tid, err = typeid.FromUUID(prefix, uuidStr)
	}

	sinkTypeID = tid
	sinkError = err
}

// BenchmarkFromBytes measures TypeID creation from UUID bytes with zero allocation
func BenchmarkFromBytes(b *testing.B) {
	uuidBytes := make([][]byte, len(testUUIDs))
	for i, u := range testUUIDs {
		uuidBytes[i] = u.Bytes()
	}

	b.ReportAllocs()
	var tid typeid.TypeID
	var err error

	for b.Loop() {
		prefix := prefixPatterns[b.N%len(prefixPatterns)]
		bytes := uuidBytes[b.N%len(uuidBytes)]

		tid, err = typeid.FromBytes(prefix, bytes)
	}

	sinkTypeID = tid
	sinkError = err
}

// BenchmarkValidation measures validation performance
func BenchmarkValidation(b *testing.B) {
	// Pre-generate valid and invalid TypeID strings
	validStrings := testTypeIDStrings

	invalidStrings := []string{
		"prefix1_01h2xcejqtf2nbrexx3vqjhp41", // invalid prefix (number)
		"prefix.01h2xcejqtf2nbrexx3vqjhp41",  // wrong separator
		"prefix_u1h2xcejqtf2nbrexx3vqjhp41",  // invalid base32 char
		"prefix_01h2xcejqtf2nbrexx3vqjhp4",   // too short
		"prefix_01h2xcejqtf2nbrexx3vqjhp411", // too long
		"_01h2xcejqtf2nbrexx3vqjhp41",        // prefix starts with _
		"prefix__01h2xcejqtf2nbrexx3vqjhp41", // double underscore
	}

	b.Run("valid", func(b *testing.B) {
		b.ReportAllocs()
		var tid typeid.TypeID
		var err error

		for b.Loop() {
			s := validStrings[b.N%len(validStrings)]
			tid, err = typeid.Parse(s)
			if err != nil {
				b.Fatalf("expected valid ID to parse: %v", err)
			}
		}

		sinkTypeID = tid
		sinkError = err
	})

	b.Run("invalid", func(b *testing.B) {
		b.ReportAllocs()
		var tid typeid.TypeID
		var err error

		for b.Loop() {
			s := invalidStrings[b.N%len(invalidStrings)]
			tid, err = typeid.Parse(s)
			if err == nil {
				b.Fatalf("expected invalid ID to fail: %s", s)
			}
		}

		sinkTypeID = tid
		sinkError = err
	})
}

// BenchmarkRoundTrip measures full encode/decode cycle
func BenchmarkRoundTrip(b *testing.B) {
	b.ReportAllocs()
	var tid typeid.TypeID
	var err error

	for b.Loop() {
		// Create new TypeID
		prefix := prefixPatterns[b.N%len(prefixPatterns)]

		tid, err = typeid.Generate(prefix)
		if err != nil {
			b.Fatal(err)
		}

		// Convert to string and back
		s := tid.String()
		tid, err = typeid.Parse(s)
		if err != nil {
			b.Fatal(err)
		}
	}

	sinkTypeID = tid
}

// BenchmarkConcurrent measures concurrent TypeID generation
func BenchmarkConcurrent(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var tid typeid.TypeID
		var err error
		iterCount := 0

		for pb.Next() {
			prefix := prefixPatterns[iterCount%len(prefixPatterns)]
			iterCount++

			tid, err = typeid.Generate(prefix)
			if err != nil {
				b.Fatal(err)
			}
		}

		// Store results to sinks to prevent elimination
		sinkTypeID = tid
		sinkError = err
	})
}

// BenchmarkTypicalUsage simulates realistic usage patterns
func BenchmarkTypicalUsage(b *testing.B) {
	// Client pattern: few creates, many string conversions
	b.Run("client", func(b *testing.B) {
		// Pre-create TypeIDs
		userID := typeid.MustGenerate("user")
		sessionID := typeid.MustGenerate("session")

		b.ReportAllocs()
		b.ResetTimer()

		var s string
		for b.Loop() {
			// Simulate typical client usage
			s = userID.String()
			s = sessionID.String()
			s = userID.String() // Used multiple times
			s = userID.Prefix()
		}

		sinkString = s
	})

	// Server pattern: many creates, few string conversions per ID
	b.Run("server", func(b *testing.B) {
		b.ReportAllocs()

		var s string
		var tid typeid.TypeID
		for b.Loop() {
			// Create new request ID
			tid = typeid.MustGenerate("req")
			s = tid.String()

			// Create response ID
			tid = typeid.MustGenerate("resp")
			s = tid.String()
		}

		sinkString = s
		sinkTypeID = tid
	})

	// Mixed pattern: balance of creates and uses
	b.Run("mixed", func(b *testing.B) {
		b.ReportAllocs()

		var s string
		var tid typeid.TypeID
		var err error

		// Pre-create some IDs
		cachedIDs := testTypeIDs[:10]

		for b.Loop() {
			switch b.N % 4 {
			case 0: // Create new
				tid, err = typeid.Generate("event")
			case 1, 2: // Use cached
				tid = cachedIDs[b.N%len(cachedIDs)]
				s = tid.String()
			case 3: // Parse
				s = testTypeIDStrings[b.N%len(testTypeIDStrings)]
				tid, err = typeid.Parse(s)
			}
		}

		sinkString = s
		sinkTypeID = tid
		sinkError = err
	})
}

// TODO: Move test types to shared location
type TestPrefix struct{}

func (TestPrefix) Prefix() string { return "test" }

type TestID struct {
	typeid.TypeID
}
