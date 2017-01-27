# axpyspeed

Testing the speed of various `y[i] += a*x[i]` (axpy) implemenations in Go. To run the tests, from the project root directory run:

`go test -bench=.`

The `_n` variable in `axpyspeed_test.go` determines the length of the `x` and `y` vectors.

Disclaimer: these are simple tests which are _at best_ meant to give some intuition of simple vector operation speeds and overheads for various implementations.

The following observations were made using:
* Go: go1.7.4 linux/amd64
* OS: Ubuntu 16.04
* CPU: i5-4670K
* BLAS: OpenBLAS 0.2.18-1ubuntu amd64

General observations:
* The native loop using `range` appears to be slightly slower than the others: is it not generating the same machine code?
* The `Cgo` and `CgoSIMD` implementations run at the same speed: the SIMD directive has no effect beyond gcc's optimization.
  * `CgoSIMDBroken` is faster than `CgoBroken`: it appears that when x and y overlap gcc doesn't use SIMD.
* `BenchmarkNativeOut` has the same performance as `BenchmarkNative`: the optimizer doesn't seem to be optimizing out the function call.

Using `_n = 1`:
* The overhead for a go routine (with the sync) is about 400 ns
* The overhead for calling the C functions is about 150 ns
* The overhead for calling from the gonum package (has checks) is about 10 ns

Using `_n = 10000`:
* Spwaning the go routines is surprisingly fast. That test isn't meant to be fast at all.
* The C implementation runs in about 0.4 the time of the native implementations.
* The C implementation runs as fast as the OpenBLAS implementation.
* The gonum native implementation runs nearly as fast as the C implementation.
  * Digging around, the call goes to `github.com/gonum/internal/asm/daxpyinc_amd64.s` which is AMD64 assembly.

Using `_n = 1000000`:
* The C implementation is much closer to the native Go implementations (about 0.98 the time).
* The OpenBLAS implementation is faster than the C implementations (about 0.85 the time).
