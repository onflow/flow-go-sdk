/*
 * Flow Go SDK
 *
 * Copyright 2019-2020 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package hash

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func checkBytes(t *testing.T, input, expected, result []byte) {
	if !bytes.Equal(expected, result) {
		t.Errorf("hash mismatch: expect: %x have: %x, input is %x", expected, result, input)
	} else {
		t.Logf("hash test ok: expect: %x, input: %x", expected, input)
	}
}

// Sanity checks of SHA3_256
func TestSha3_256(t *testing.T) {
	input := []byte("test")
	expected, _ := hex.DecodeString("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab80")

	alg := NewSHA3_256()
	hash := alg.ComputeHash(input)
	checkBytes(t, input, expected, hash)

	alg.Reset()
	_, _ = alg.Write([]byte("te"))
	_, _ = alg.Write([]byte("s"))
	_, _ = alg.Write([]byte("t"))
	hash = alg.SumHash()
	checkBytes(t, input, expected, hash)
}

// Sanity checks of SHA3_384
func TestSha3_384(t *testing.T) {
	input := []byte("test")
	expected, _ := hex.DecodeString("e516dabb23b6e30026863543282780a3ae0dccf05551cf0295178d7ff0f1b41eecb9db3ff219007c4e097260d58621bd")

	alg := NewSHA3_384()
	hash := alg.ComputeHash(input)
	checkBytes(t, input, expected, hash)

	alg.Reset()
	_, _ = alg.Write([]byte("te"))
	_, _ = alg.Write([]byte("s"))
	_, _ = alg.Write([]byte("t"))
	hash = alg.SumHash()
	checkBytes(t, input, expected, hash)
}

// Sanity checks of SHA2_256
func TestSha2_256(t *testing.T) {
	input := []byte("test")
	expected, _ := hex.DecodeString("9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08")

	alg := NewSHA2_256()
	hash := alg.ComputeHash(input)
	checkBytes(t, input, expected, hash)

	alg.Reset()
	_, _ = alg.Write([]byte("te"))
	_, _ = alg.Write([]byte("s"))
	_, _ = alg.Write([]byte("t"))
	hash = alg.SumHash()
	checkBytes(t, input, expected, hash)
}

// Sanity checks of SHA2_256
func TestSha2_384(t *testing.T) {
	input := []byte("test")
	expected, _ := hex.DecodeString("768412320f7b0aa5812fce428dc4706b3cae50e02a64caa16a782249bfe8efc4b7ef1ccb126255d196047dfedf17a0a9")

	alg := NewSHA2_384()
	hash := alg.ComputeHash(input)
	checkBytes(t, input, expected, hash)

	alg.Reset()
	_, _ = alg.Write([]byte("te"))
	_, _ = alg.Write([]byte("s"))
	_, _ = alg.Write([]byte("t"))
	hash = alg.SumHash()
	checkBytes(t, input, expected, hash)
}

// SHA3_256 bench
func BenchmarkSha3_256(b *testing.B) {
	a := []byte("Bench me!")
	alg := NewSHA3_256()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		alg.ComputeHash(a)
	}
	return
}

// SHA3_384 bench
func BenchmarkSha3_384(b *testing.B) {
	a := []byte("Bench me!")
	alg := NewSHA3_384()
	for i := 0; i < b.N; i++ {
		alg.ComputeHash(a)
	}
	return
}

// SHA2_256 bench
func BenchmarkSha2_256(b *testing.B) {
	a := []byte("Bench me!")
	alg := NewSHA2_256()
	for i := 0; i < b.N; i++ {
		alg.ComputeHash(a)
	}
	return
}

// SHA2_256 bench
func BenchmarkSha2_384(b *testing.B) {
	a := []byte("Bench me!")
	alg := NewSHA2_384()
	for i := 0; i < b.N; i++ {
		alg.ComputeHash(a)
	}
	return
}
