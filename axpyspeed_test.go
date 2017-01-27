package axpyspeed

import (
	"testing"

	"github.com/gonum/blas/blas64"
	"github.com/gonum/blas/cgo"
	"github.com/gonum/blas/native"
)

const (
	_n = 10000 // length of x and y
	_a = 10    // value of a
)

// Make the x and y arrays of length _n, ranging 0.._n
func makeData() ([]float64, []float64) {
	x := make([]float64, _n)
	y := make([]float64, _n)
	for i := 0; i < _n; i++ {
		x[i] = float64(i)
		y[i] = float64(i)
	}
	return x, y
}

// Test that the axpy calcuation went as expected
func assertAXPY(t *testing.T, x, y []float64) {
	if len(x) != _n {
		t.Fatalf("len(x) is not n")
	}
	if len(y) != _n {
		t.Fatalf("len(y) is not n")
	}
	for i, v := range x {
		if v != float64(i) {
			t.Fatalf("x is not unchanged")
		}
		if y[i] != _a*v+float64(i) {
			t.Fatalf("y is not a*x+y", y[i], _a*v)
		}
	}
}

func TestNative(t *testing.T) {
	x, y := makeData()
	Native(_a, x, y)
	assertAXPY(t, x, y)
}

func TestNativeLen(t *testing.T) {
	x, y := makeData()
	NativeLen(_a, x, y)
	assertAXPY(t, x, y)
}

func TestNativeRange(t *testing.T) {
	x, y := makeData()
	NativeRange(_a, x, y)
	assertAXPY(t, x, y)
}

func TestNativeGoRoutine(t *testing.T) {
	x, y := makeData()
	NativeGoRoutine(_a, x, y)
	assertAXPY(t, x, y)
}

func TestCgo(t *testing.T) {
	x, y := makeData()
	Cgo(_a, x, y)
	assertAXPY(t, x, y)
}

func TestCgoBroken(t *testing.T) {
	// passes because x==y in the test so it is invariant to the failure
	x, y := makeData()
	CgoBroken(_a, x, y)
	assertAXPY(t, x, y)
}

func TestCgoSIMDBroken(t *testing.T) {
	// passes because x==y in the test so it is invariant to the failure
	x, y := makeData()
	CgoSIMDBroken(_a, x, y)
	assertAXPY(t, x, y)
}

func TestCgoSIMD(t *testing.T) {
	x, y := makeData()
	CgoSIMD(_a, x, y)
	assertAXPY(t, x, y)
}

func TestGonumBLAS(t *testing.T) {
	x, y := makeData()
	GonumBLAS(_a, x, y)
	assertAXPY(t, x, y)
}

func BenchmarkNative(b *testing.B) {
	x, y := makeData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Native(_a, x, y)
	}
}

// Check that the benchmark isn't optimizing out the call
var EnsureOutput float64

func BenchmarkNativeOut(b *testing.B) {
	x, y := makeData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Native(_a, x, y)
	}
	b.StopTimer()
	EnsureOutput = y[len(y)-1]
}

func BenchmarkNativeLen(b *testing.B) {
	x, y := makeData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NativeLen(_a, x, y)
	}
}

func BenchmarkNativeRange(b *testing.B) {
	x, y := makeData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NativeRange(_a, x, y)
	}
}

func BenchmarkNativeGoRoutine(b *testing.B) {
	x, y := makeData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NativeGoRoutine(_a, x, y)
	}
}

func BenchmarkCgo(b *testing.B) {
	x, y := makeData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Cgo(_a, x, y)
	}
}

func BenchmarkCgoSIMD(b *testing.B) {
	x, y := makeData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CgoSIMD(_a, x, y)
	}
}

func BenchmarkCgoBroken(b *testing.B) {
	x, y := makeData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CgoBroken(_a, x, y)
	}
}

func BenchmarkCgoSIMDBroken(b *testing.B) {
	x, y := makeData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CgoSIMDBroken(_a, x, y)
	}
}

func BenchmarkGonumBLAS(b *testing.B) {
	blas64.Use(native.Implementation{})
	x, y := makeData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GonumBLAS(_a, x, y)
	}
}

func BenchmarkGonumBLASCGO(b *testing.B) {
	blas64.Use(cgo.Implementation{})
	x, y := makeData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GonumBLAS(_a, x, y)
	}
}
