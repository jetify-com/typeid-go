package base32

import (
	"testing"
)

// Benchmark suite for base32 encoding/decoding operations.
//
// These benchmarks measure performance of different approaches:
// - Allocating vs zero-allocation functions
// - Fresh allocation vs slice reuse patterns
//
// All benchmarks cycle through realistic UUID patterns to ensure representative performance.
// Results are stored in package-level sinks to prevent compiler optimizations.

// Test data patterns for varied input benchmarks
var testPatterns = []struct {
	name string
	data [16]byte
}{
	// UUIDv7 example - time-ordered UUID with timestamp prefix
	{"uuidv7", [16]byte{0x01, 0x8D, 0x5C, 0x9F, 0x12, 0x34, 0x70, 0x00, 0x80, 0x00, 0x1A, 0x2B, 0x3C, 0x4D, 0x5E, 0x6F}},

	// UUIDv4 example - random UUID
	{"uuidv4", [16]byte{0xF4, 0x7A, 0xC1, 0x0B, 0x58, 0xCC, 0x4F, 0x72, 0xAA, 0x24, 0x89, 0x15, 0xD3, 0x78, 0x42, 0x9E}},

	// Alternating pattern - tests non-uniform distribution
	{"alternating", [16]byte{0, 255, 0, 255, 0, 255, 0, 255, 0, 255, 0, 255, 0, 255, 0, 255}},

	// All zeros - edge case testing
	{"zeros", [16]byte{}},
}

// Pre-encoded data for decode benchmarks - encodings of all test patterns
var (
	benchmarkEncodedStrings = make([]string, len(testPatterns))
	benchmarkEncodedBytes   = make([][]byte, len(testPatterns))
)

func init() {
	for i, pattern := range testPatterns {
		benchmarkEncodedStrings[i] = EncodeToString(pattern.data)
		benchmarkEncodedBytes[i] = []byte(benchmarkEncodedStrings[i])
	}
}

// Sink variables to prevent compiler optimizations
var (
	sink      string
	sinkBytes []byte
	sinkInt   int
	sinkError error
)

func BenchmarkEncodeToString(b *testing.B) {
	b.ReportAllocs()
	var r string
	for b.Loop() {
		// Cycle through all test patterns for realistic performance measurement
		data := testPatterns[b.N%len(testPatterns)].data
		r = EncodeToString(data)
	}
	// always store the result to a package level variable
	// so the compiler cannot eliminate the Benchmark itself
	sink = r
}

func BenchmarkEncode(b *testing.B) {
	dst := make([]byte, 26)
	b.ReportAllocs()
	// No need for b.ResetTimer() with b.Loop() - setup is automatically excluded
	for b.Loop() {
		data := testPatterns[b.N%len(testPatterns)].data
		Encode(dst, data)
	}
	// Store dst to prevent elimination - Encode modifies dst in place
	sinkBytes = dst
}

// BenchmarkAppendEncode benchmarks AppendEncode starting with nil slice
func BenchmarkAppendEncode(b *testing.B) {
	b.ReportAllocs()
	var r []byte
	for b.Loop() {
		// always record the result to prevent compiler eliminating the function call
		data := testPatterns[b.N%len(testPatterns)].data
		r = AppendEncode(nil, data)
	}
	// always store the result to a package level variable
	// so the compiler cannot eliminate the Benchmark itself
	sinkBytes = r
}

// BenchmarkAppendEncodeReuse benchmarks AppendEncode with slice reuse
func BenchmarkAppendEncodeReuse(b *testing.B) {
	result := make([]byte, 0, 26)
	b.ReportAllocs()
	// No need for b.ResetTimer() with b.Loop() - setup is automatically excluded
	for b.Loop() {
		data := testPatterns[b.N%len(testPatterns)].data
		result = AppendEncode(result[:0], data)
	}
	// Store final result to prevent elimination
	sinkBytes = result
}

func BenchmarkDecodeString(b *testing.B) {
	b.ReportAllocs()
	var r []byte
	var err error
	for b.Loop() {
		// Cycle through all encoded patterns
		encoded := benchmarkEncodedStrings[b.N%len(benchmarkEncodedStrings)]
		r, err = DecodeString(encoded)
	}
	// always store the result to package level variables
	// so the compiler cannot eliminate the Benchmark itself
	sinkBytes = r
	sinkError = err
}

func BenchmarkDecode(b *testing.B) {
	dst := make([]byte, 16)
	b.ReportAllocs()
	// No need for b.ResetTimer() with b.Loop() - setup is automatically excluded
	var n int
	var err error
	for b.Loop() {
		// Cycle through all encoded patterns
		encoded := benchmarkEncodedBytes[b.N%len(benchmarkEncodedBytes)]
		n, err = Decode(dst, encoded)
	}
	// always store the results to package level variables
	// so the compiler cannot eliminate the Benchmark itself
	sinkInt = n
	sinkError = err
	sinkBytes = dst
}

// BenchmarkAppendDecode benchmarks AppendDecode starting with nil slice
func BenchmarkAppendDecode(b *testing.B) {
	b.ReportAllocs()
	var r []byte
	var err error
	for b.Loop() {
		// Cycle through all encoded patterns
		encoded := benchmarkEncodedBytes[b.N%len(benchmarkEncodedBytes)]
		r, err = AppendDecode(nil, encoded)
	}
	// always store the result to package level variables
	// so the compiler cannot eliminate the Benchmark itself
	sinkBytes = r
	sinkError = err
}

// BenchmarkAppendDecodeReuse benchmarks AppendDecode with slice reuse
func BenchmarkAppendDecodeReuse(b *testing.B) {
	result := make([]byte, 0, 16)
	b.ReportAllocs()
	// No need for b.ResetTimer() with b.Loop() - setup is automatically excluded
	var err error
	for b.Loop() {
		// Cycle through all encoded patterns
		encoded := benchmarkEncodedBytes[b.N%len(benchmarkEncodedBytes)]
		result, err = AppendDecode(result[:0], encoded)
	}
	// Store final results to prevent elimination
	sinkBytes = result
	sinkError = err
}

// BenchmarkInputPatterns tests encoding performance with different input patterns
// This helps ensure our performance is representative across various data types
func BenchmarkInputPatterns(b *testing.B) {
	for _, pattern := range testPatterns {
		b.Run(pattern.name, func(b *testing.B) {
			var r string
			for b.Loop() {
				r = EncodeToString(pattern.data)
			}
			sink = r
		})
	}
}
