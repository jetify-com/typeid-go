//nolint:all
package typeid_test

import (
	"fmt"
	"testing"

	"github.com/gofrs/uuid/v5"
	"go.jetify.com/typeid"
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

// TODO: define these in a shared file if we're gonna use in several tests.

type TestPrefix struct{}

func (TestPrefix) Prefix() string { return "prefix" }

type TestID struct {
	typeid.TypeID[TestPrefix]
}
