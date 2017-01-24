// Package axpyspeed implements a few blas axpy functions such that their
// speeds can be compared.
package axpyspeed

/*
#cgo CFLAGS: -fopenmp -march=native -O3

void do_axpy(int n, double a, double* x, double* y)
{
	int i;
	for (i = 0; i < n; i++) {
		y[i] += a * x[i];
	}
}

void do_simd_axpy(int n, double a, double* x, double* y)
{
	int i;
#pragma omp simd
	for (i = 0; i < n; i++) {
		y[i] += a * x[i];
	}
}
*/
import "C"

import (
	"sync"

	"github.com/gonum/blas/blas64"
)

// Native is the most basic loop implementation.
func Native(a float64, x, y []float64) {
	n := len(x)
	for i := 0; i < n; i++ {
		y[i] += a * x[i]
	}
}

// NativeLen is a basic loop but calls len(x) in the for loop.
func NativeLen(a float64, x, y []float64) {
	for i := 0; i < len(x); i++ {
		y[i] += a * x[i]
	}
}

// NativeRange is the basic loop implemented with range.
func NativeRange(a float64, x, y []float64) {
	for i, v := range x {
		y[i] += a * v
	}
}

// NativeGoRoutine launches a concurrent go routine for each loop iteration.
func NativeGoRoutine(a float64, x, y []float64) {
	var w sync.WaitGroup
	for i, v := range x {
		w.Add(1)
		go func(i int, v float64) {
			y[i] += a * v
			w.Done()
		}(i, v)
	}
	w.Wait()
}

// Cgo calls the basic C implementation of axpy
func Cgo(a float64, x, y []float64) {
	n := len(x)
	C.do_axpy(
		C.int(n),
		C.double(a),
		(*C.double)(&x[0]),
		(*C.double)(&y[0]))
}

// CgoSIMD calls the basic C implementation of axpy with pragma omp simd
func CgoSIMD(a float64, x, y []float64) {
	n := len(x)
	C.do_simd_axpy(
		C.int(n),
		C.double(a),
		(*C.double)(&x[0]),
		(*C.double)(&y[0]))
}

// GonumBlas calls gonum's implementation of axpy
func GonumBlas(a float64, x, y []float64) {
	n := len(x)
	blas64.Implementation().Daxpy(n, a, x, 1, y, 1)
}
