// +build relic

package crypto

// #cgo CFLAGS: -g -Wall -std=c99 -I./ -I./relic/include -I./relic/include/low
// #cgo LDFLAGS: -Lrelic/build/lib -l relic_s
// #include "bls_include.h"
import "C"

// TODO: remove -wall after reaching a stable version
// TODO: enable QUIET in relic

// Go wrappers to C types
type pointG1 C.ep_st
type pointG2 C.ep2_st
type scalar C.bn_st
type ctx struct {
	relicCtx *C.ctx_t
	precCtx  *C.prec_st
}

var signatureLengthBLS_BLS12381 = int(C._getSignatureLengthBLS_BLS12381())
var pubKeyLengthBLS_BLS12381 = int(C._getPubKeyLengthBLS_BLS12381())
var prKeyLengthBLS_BLS12381 = int(C._getPrKeyLengthBLS_BLS12381())

// init sets the context of BLS12381 curve
func (a *BLS_BLS12381Algo) init() error {
	// sanity checks of lengths
	if a.prKeyLength != PrKeyLenBLS_BLS12381 ||
		a.pubKeyLength != PubKeyLenBLS_BLS12381 ||
		a.signatureLength != SignatureLenBLS_BLS12381 {
		return cryptoError{"BLS Lengths in types.go are not matching bls_include.h"}
	}

	// Inits relic context and sets the B12_381 context
	c := C.relic_init_BLS12_381()
	if c == nil {
		return cryptoError{"Relic core init failed"}
	}
	a.context.relicCtx = c
	a.context.precCtx = C.init_precomputed_data_BLS12_381()
	return nil
}

// reinit the context of BLS12381 curve assuming there was a previous call to init()
// If the implementation evolves and relic has multiple contexts,
// reinit should be called at every a. operation.
func (a *BLS_BLS12381Algo) reinit() {
	C.core_set(a.context.relicCtx)
	C.precomputed_data_set(a.context.precCtx)
}

// Exponentiation in G1 (scalar point multiplication)
func (p *pointG1) _G1scalarPointMult(res *pointG1, expo *scalar) {
	C._G1scalarPointMult((*C.ep_st)(res), (*C.ep_st)(p), (*C.bn_st)(expo))
}

// Exponentiation of g1 in G1
// This function is for DEBUG/TESTs purpose only
func _G1scalarGenMult(res *pointG1, expo *scalar) {
	C._G1scalarGenMult((*C.ep_st)(res), (*C.bn_st)(expo))
}

// Exponentiation of g2 in G2
func _G2scalarGenMult(res *pointG2, expo *scalar) {
	C._G2scalarGenMult((*C.ep2_st)(res), (*C.bn_st)(expo))
}

// TEST/DEBUG

// seeds the internal relic random function
func seedRelic(seed []byte) {
	// TODO: define the length of seed
	C._seed_relic((*C.uchar)(&seed[0]), (C.int)(len(seed)))
}

// returns a random number on Z/Z.r
func randZr(x *scalar) error {
	C._bn_randZr((*C.bn_st)(x))
	if x == nil {
		return cryptoError{"the memory allocation of the random number has failed"}
	}
	return nil
}

// TEST/DEBUG/BENCH
// returns the hash to G1 point
func hashToG1(data []byte) *pointG1 {
	l := len(data)
	var h pointG1
	C.mapToG1((*C.ep_st)(&h), (*C.uchar)(&data[0]), (C.int)(l))
	return &h
}

// TEST/DEBUG/BENCH
// wraps a call to optimized SwU algorithm since cgo can't be used
// in go test files
func OpSwUUnitTest(output []byte, input []byte) {
	C.opswu_test((*C.uchar)(&output[0]),
		(*C.uchar)(&input[0]),
		SignatureLenBLS_BLS12381)
}

// sets a scalar to a small integer
func (x *scalar) setInt(a int) {
	C.bn_set_dig((*C.bn_st)(x), (C.uint64_t)(a))
}

// writePointG2 writes a G2 point in a slice of bytes
func writePointG2(dest []byte, a *pointG2) {
	C._ep2_write_bin_compact((*C.uchar)(&dest[0]),
		(*C.ep2_st)(a),
		(C.int)(pubKeyLengthBLS_BLS12381),
	)
}

// readVerifVector reads a G2 point from a slice of bytes
func readPointG2(a *pointG2, data []byte) {
	C._ep2_read_bin_compact((*C.ep2_st)(a),
		(*C.uchar)(&data[0]),
		(C.int)(len(data)),
	)
}

// computes a bls signature
func (a *BLS_BLS12381Algo) blsSign(sk *scalar, data []byte) Signature {
	s := make([]byte, a.signatureLength)

	C._blsSign((*C.uchar)(&s[0]),
		(*C.bn_st)(sk),
		(*C.uchar)(&data[0]),
		(C.int)(len(data)))
	return s
}

// Checks the validity of a bls signature
func (a *BLS_BLS12381Algo) blsVerify(pk *pointG2, s Signature, data []byte) bool {
	verif := C._blsVerify((*C.ep2_st)(pk),
		(*C.uchar)(&s[0]),
		(*C.uchar)(&data[0]),
		(C.int)(len(data)))

	const sigValid = 1 // same value as in include.h
	const sigErr = 0xFF

	if verif == sigErr {
		panic("Relic memory allocation failed")
	}
	return (verif == sigValid)
}
