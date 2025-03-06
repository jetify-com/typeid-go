//nolint:all
package typeid_test

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"

	"github.com/gofrs/uuid/v5"
	"go.jetify.com/typeid"
	"go.jetify.com/typeid/base32"
)

func BenchmarkNew(b *testing.B) {
	b.Run("id=untyped", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			typeid.WithPrefix("prefix")
		}
	})
	b.Run("id=typed", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			typeid.New[TestID]()
		}
	})
	b.Run("id=uuid", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			uuid.NewV7()
		}
	})
	// Add benchmark for different prefix lengths
	b.Run("prefix=short", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			typeid.WithPrefix("s")
		}
	})
	b.Run("prefix=medium", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			typeid.WithPrefix("medium")
		}
	})
	b.Run("prefix=long", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			typeid.WithPrefix("thisislongprefix")
		}
	})
}

func BenchmarkString(b *testing.B) {
	b.Run("id=untyped", func(b *testing.B) {
		b.Run(benchUntypedString(1))
		b.Run(benchUntypedString(8))
		b.Run(benchUntypedString(64))
		b.Run(benchUntypedString(4096))
	})
	b.Run("id=typed", func(b *testing.B) {
		b.Run(benchTypedString(1))
		b.Run(benchTypedString(8))
		b.Run(benchTypedString(64))
		b.Run(benchTypedString(4096))
	})
	b.Run("id=uuid", func(b *testing.B) {
		b.Run(benchUUIDString(1))
		b.Run(benchUUIDString(8))
		b.Run(benchUUIDString(64))
		b.Run(benchUUIDString(4096))
	})
}

func benchUntypedString(n int) (string, func(*testing.B)) {
	ids := make([]typeid.AnyID, n)
	for i := range ids {
		ids[i] = typeid.Must(typeid.WithPrefix("prefix"))
	}
	return fmt.Sprintf("n=%d", n), func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, id := range ids {
				_ = id.String()
			}
		}
		b.ReportMetric(float64(n*b.N)/b.Elapsed().Seconds(), "id/s")
	}
}

func benchTypedString(n int) (string, func(*testing.B)) {
	ids := make([]TestID, n)
	for i := range ids {
		ids[i] = typeid.Must(typeid.New[TestID]())
	}
	return fmt.Sprintf("n=%d", n), func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, id := range ids {
				_ = id.String()
			}
		}
		b.ReportMetric(float64(n*b.N)/b.Elapsed().Seconds(), "id/s")
	}
}

func benchUUIDString(n int) (string, func(*testing.B)) {
	uuids := make([]uuid.UUID, n)
	for i := range uuids {
		uuids[i] = uuid.Must(uuid.NewV7())
	}
	return fmt.Sprintf("n=%d", n), func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, id := range uuids {
				_ = id.String()
			}
		}
		b.ReportMetric(float64(n*b.N)/b.Elapsed().Seconds(), "id/s")
	}
}

func BenchmarkFrom(b *testing.B) {
	b.Run("id=untyped", func(b *testing.B) {
		b.Run(benchUntypedFrom(1))
		b.Run(benchUntypedFrom(8))
		b.Run(benchUntypedFrom(64))
		b.Run(benchUntypedFrom(4096))
	})
	b.Run("id=typed", func(b *testing.B) {
		b.Run(benchTypedFrom(1))
		b.Run(benchTypedFrom(8))
		b.Run(benchTypedFrom(64))
		b.Run(benchTypedFrom(4096))
	})
}

func benchUntypedFrom(n int) (string, func(*testing.B)) {
	ids := make([]struct{ prefix, suffix string }, n)
	for i := range ids {
		id := typeid.Must(typeid.WithPrefix("prefix"))
		ids[i].prefix, ids[i].suffix = id.Prefix(), id.Suffix()
	}
	return fmt.Sprintf("n=%d", n), func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, id := range ids {
				typeid.From(id.prefix, id.suffix)
			}
		}
		b.ReportMetric(float64(n*b.N)/b.Elapsed().Seconds(), "id/s")
	}
}

func benchTypedFrom(n int) (string, func(*testing.B)) {
	suffixes := make([]string, n)
	for i := range suffixes {
		suffixes[i] = typeid.Must(typeid.New[TestID]()).Suffix()
	}
	return fmt.Sprintf("n=%d", n), func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, suffix := range suffixes {
				typeid.FromSuffix[TestID](suffix)
			}
		}
		b.ReportMetric(float64(n*b.N)/b.Elapsed().Seconds(), "id/s")
	}
}

func BenchmarkFromString(b *testing.B) {
	b.Run("id=untyped", func(b *testing.B) {
		b.Run(benchUntypedFromString(1))
		b.Run(benchUntypedFromString(8))
		b.Run(benchUntypedFromString(64))
		b.Run(benchUntypedFromString(4096))
	})
	b.Run("id=typed", func(b *testing.B) {
		b.Run(benchTypedFromString(1))
		b.Run(benchTypedFromString(8))
		b.Run(benchTypedFromString(64))
		b.Run(benchTypedFromString(4096))
	})
	b.Run("id=uuid", func(b *testing.B) {
		b.Run(benchUUIDFromString(1))
		b.Run(benchUUIDFromString(8))
		b.Run(benchUUIDFromString(64))
		b.Run(benchUUIDFromString(4096))
	})
}

func benchUntypedFromString(n int) (string, func(*testing.B)) {
	ids := make([]string, n)
	for i := range ids {
		ids[i] = typeid.Must(typeid.WithPrefix("prefix")).String()
	}
	return fmt.Sprintf("n=%d", n), func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, id := range ids {
				typeid.FromString(id)
			}
		}
		b.ReportMetric(float64(n*b.N)/b.Elapsed().Seconds(), "id/s")
	}
}

func benchTypedFromString(n int) (string, func(*testing.B)) {
	ids := make([]string, n)
	for i := range ids {
		ids[i] = typeid.Must(typeid.New[TestID]()).String()
	}
	return fmt.Sprintf("n=%d", n), func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, id := range ids {
				typeid.Parse[TestID](id)
			}
		}
		b.ReportMetric(float64(n*b.N)/b.Elapsed().Seconds(), "id/s")
	}
}

func benchUUIDFromString(n int) (string, func(*testing.B)) {
	uuids := make([]string, n)
	for i := range uuids {
		uuids[i] = uuid.Must(uuid.NewV7()).String()
	}
	return fmt.Sprintf("n=%d", n), func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, id := range uuids {
				uuid.FromString(id)
			}
		}
		b.ReportMetric(float64(n*b.N)/b.Elapsed().Seconds(), "id/s")
	}
}

func BenchmarkFromBytes(b *testing.B) {
	b.Run("id=untyped", func(b *testing.B) {
		b.Run(benchUntypedFromBytes(1))
		b.Run(benchUntypedFromBytes(8))
		b.Run(benchUntypedFromBytes(64))
		b.Run(benchUntypedFromBytes(4096))
	})
	b.Run("id=typed", func(b *testing.B) {
		b.Run(benchTypedFromBytes(1))
		b.Run(benchTypedFromBytes(8))
		b.Run(benchTypedFromBytes(64))
		b.Run(benchTypedFromBytes(4096))
	})
	b.Run("id=uuid", func(b *testing.B) {
		b.Run(benchUUIDFromBytes(1))
		b.Run(benchUUIDFromBytes(8))
		b.Run(benchUUIDFromBytes(64))
		b.Run(benchUUIDFromBytes(4096))
	})
}

func benchUntypedFromBytes(n int) (string, func(*testing.B)) {
	ids := make([][]byte, n)
	for i := range ids {
		ids[i] = typeid.Must(typeid.WithPrefix("prefix")).UUIDBytes()
	}
	return fmt.Sprintf("n=%d", n), func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, id := range ids {
				typeid.FromUUIDBytesWithPrefix("prefix", id)
			}
		}
		b.ReportMetric(float64(n*b.N)/b.Elapsed().Seconds(), "id/s")
	}
}

func benchTypedFromBytes(n int) (string, func(*testing.B)) {
	ids := make([][]byte, n)
	for i := range ids {
		ids[i] = typeid.Must(typeid.New[TestID]()).UUIDBytes()
	}
	return fmt.Sprintf("n=%d", n), func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, id := range ids {
				typeid.FromUUIDBytesWithPrefix("prefix", id)
			}
		}
		b.ReportMetric(float64(n*b.N)/b.Elapsed().Seconds(), "id/s")
	}
}

func benchUUIDFromBytes(n int) (string, func(*testing.B)) {
	uuids := make([][]byte, n)
	for i := range uuids {
		uuids[i] = uuid.Must(uuid.NewV7()).Bytes()
	}
	return fmt.Sprintf("n=%d", n), func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, id := range uuids {
				uuid.FromBytes(id)
			}
		}
		b.ReportMetric(float64(n*b.N)/b.Elapsed().Seconds(), "id/s")
	}
}

func BenchmarkSuffix(b *testing.B) {
	b.Run("id=untyped", func(b *testing.B) {
		id := typeid.Must(typeid.WithPrefix("prefix"))

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			id.Suffix()
		}
	})
	b.Run("id=typed", func(b *testing.B) {
		id := typeid.Must(typeid.New[TestID]())

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			id.Suffix()
		}
	})
}

func BenchmarkUUIDBytes(b *testing.B) {
	b.Run("id=untyped", func(b *testing.B) {
		id := typeid.Must(typeid.WithPrefix("prefix"))

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			id.UUIDBytes()
		}
	})
	b.Run("id=typed", func(b *testing.B) {
		id := typeid.Must(typeid.New[TestID]())

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			id.UUIDBytes()
		}
	})
	b.Run("id=uuid", func(b *testing.B) {
		id := uuid.Must(uuid.NewV7())

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			id.Bytes()
		}
	})
}

func BenchmarkNewWithPrefix(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = typeid.Must(typeid.WithPrefix("prefix"))
	}
}

func BenchmarkEncodeDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tid := typeid.Must(typeid.WithPrefix("prefix"))
		_ = typeid.Must(typeid.FromString(tid.String()))
	}
}

// Benchmark Base32 operations directly
func BenchmarkBase32(b *testing.B) {
	b.Run("encode", func(b *testing.B) {
		uid := uuid.Must(uuid.NewV7())
		var bytes [16]byte
		copy(bytes[:], uid.Bytes())
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = base32.Encode(bytes)
		}
	})
	b.Run("decode", func(b *testing.B) {
		uid := uuid.Must(uuid.NewV7())
		var bytes [16]byte
		copy(bytes[:], uid.Bytes())
		encoded := base32.Encode(bytes)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = base32.Decode(encoded)
		}
	})
}

// Benchmark memory usage with different batch sizes
func BenchmarkMemoryUsage(b *testing.B) {
	benchSizes := []int{100, 1000, 10000}
	for _, size := range benchSizes {
		size := size // capture range variable
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			// Pre-allocate a slice to avoid measuring slice growth
			ids := make([]typeid.AnyID, size)
			b.Cleanup(func() {
				// Clear the slice to help GC
				for i := range ids {
					ids[i] = typeid.AnyID{}
				}
				ids = nil
			})

			b.ReportAllocs()
			b.SetBytes(int64(size * 16)) // Each UUID is 16 bytes
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				for j := range ids {
					ids[j] = typeid.Must(typeid.WithPrefix("prefix"))
				}
				// Force memory pressure to ensure GC behavior is measured
				runtime.GC()
			}
		})
	}
}

// Benchmark validation
func BenchmarkValidation(b *testing.B) {
	validIDs := make([]string, 100)
	invalidIDs := make([]string, 100)

	for i := range validIDs {
		validIDs[i] = typeid.Must(typeid.WithPrefix("prefix")).String()
		if i < len(invalidIDs) {
			// Create definitely invalid IDs by:
			// 1. Using invalid prefix characters
			// 2. Wrong separator
			// 3. Invalid base32 characters
			switch i % 3 {
			case 0:
				// Invalid prefix (contains number)
				invalidIDs[i] = "prefix1_01h2xcejqtf2nbrexx3vqjhp41"
			case 1:
				// Wrong separator (using . instead of _)
				invalidIDs[i] = "prefix.01h2xcejqtf2nbrexx3vqjhp41"
			case 2:
				// Invalid base32 character in suffix (using 'u' which isn't in the alphabet)
				invalidIDs[i] = "prefix_u1h2xcejqtf2nbrexx3vqjhp41"
			}
		}
	}

	b.Run("valid", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			// Test all valid IDs in each iteration
			for _, id := range validIDs {
				_, err := typeid.FromString(id)
				if err != nil {
					b.Fatalf("Expected valid ID to pass validation: %v", err)
				}
			}
		}
	})

	b.Run("invalid", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			// Test all invalid IDs in each iteration
			for _, id := range invalidIDs {
				_, err := typeid.FromString(id)
				if err == nil {
					b.Fatalf("Expected invalid ID to fail validation for ID: %s", id)
				}
			}
		}
	})
}

// Benchmark parallel ID generation
func BenchmarkParallelGeneration(b *testing.B) {
	benchCases := []struct {
		name      string
		procs     int
		batchSize int
	}{
		{"procs=4_batch=100", 4, 100},
		{"procs=8_batch=100", runtime.GOMAXPROCS(0) * 2, 100},
		{"procs=4_batch=1000", 4, 1000},
		{"procs=8_batch=1000", runtime.GOMAXPROCS(0) * 2, 1000},
	}

	for _, bc := range benchCases {
		bc := bc // capture range variable
		b.Run(bc.name, func(b *testing.B) {
			// Pre-allocate a slice of prefixes for each processor
			prefixes := make([]string, bc.procs)
			for i := range prefixes {
				// Use valid prefixes with only [a-z_]
				prefixes[i] = fmt.Sprintf("prefix_%c", 'a'+i)
			}

			b.SetParallelism(bc.procs)
			b.ReportAllocs()
			b.ResetTimer()

			b.RunParallel(func(pb *testing.PB) {
				// Pre-allocate slice with capacity
				ids := make([]typeid.AnyID, 0, bc.batchSize)
				// Keep track of which prefix we're using
				prefixIdx := 0

				for pb.Next() {
					// Clear the slice but keep capacity
					ids = ids[:0]

					// Use all prefixes in each iteration
					for j := 0; j < bc.batchSize; j++ {
						// Cycle through prefixes
						prefix := prefixes[prefixIdx%len(prefixes)]
						prefixIdx++
						ids = append(ids, typeid.Must(typeid.WithPrefix(prefix)))
					}
				}
			})
		})
	}
}

// Benchmark mixed operations to simulate real-world usage patterns
func BenchmarkMixedOperations(b *testing.B) {
	// Pre-generate all operations and test data
	type op struct {
		kind string
		id   typeid.AnyID // For operations that need an existing ID
	}

	// Create a slice of operations in the ratio we want to test
	numOps := 100 // Number of operations to pre-generate
	ops := make([]op, numOps)
	ids := make([]typeid.AnyID, numOps/2) // Pre-generate some IDs for operations that need them

	for i := range ids {
		ids[i] = typeid.Must(typeid.WithPrefix("prefix"))
	}

	// Initialize random number generator
	src := rand.NewSource(1234)
	rnd := rand.New(src)

	// Generate operations in desired ratios
	for i := range ops {
		r := rnd.Float32()
		switch {
		case r < 0.1: // 10% new IDs
			ops[i] = op{kind: "create"}
		case r < 0.4: // 30% toString
			ops[i] = op{kind: "toString", id: ids[rnd.Intn(len(ids))]}
		case r < 0.7: // 30% parse
			ops[i] = op{kind: "parse", id: ids[rnd.Intn(len(ids))]}
		default: // 30% validate
			ops[i] = op{kind: "validate", id: ids[rnd.Intn(len(ids))]}
		}
	}

	b.ReportAllocs()
	b.ResetTimer()

	// Ensure we go through all operations at least once
	for i := 0; i < b.N; i++ {
		// Do all operations once per iteration
		for _, op := range ops {
			switch op.kind {
			case "create":
				_ = typeid.Must(typeid.WithPrefix("prefix"))
			case "toString":
				_ = op.id.String()
			case "parse":
				_, _ = typeid.FromString(op.id.String())
			case "validate":
				_, _ = typeid.FromString(op.id.String())
			}
		}
	}
}

// TODO: define these in a shared file if we're gonna use in several tests.

type TestPrefix struct{}

func (TestPrefix) Prefix() string { return "prefix" }

type TestID struct {
	typeid.TypeID[TestPrefix]
}
