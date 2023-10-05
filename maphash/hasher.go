// From https://github.com/dolthub/maphash
// Copyright 2022 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package maphash

import (
	"unsafe"
)

// Hasher hashes values of type K.
// Uses runtime AES-based hashing.
type Hasher[K comparable] struct {
	//hash func(key K) uintptr //hashfn
	hash hashfn
	seed uintptr
	typ int8
}

// NewHasher creates a new Hasher[K] with a random seed.
func NewHasher[K comparable](typ int8) Hasher[K] {
	h := Hasher[K]{}
	h.seed = newHashSeed()
	h.typ = typ
	switch(typ) {
	case 0:
		h.hash = getDefaultHasher[K]()
	case 1:
		h.hash = getExperimentalHasher[K]()
	}

	return h
}

// NewSeed returns a copy of |h| with a new hash seed.
func NewSeed[K comparable](h Hasher[K]) Hasher[K] {
	return Hasher[K]{
		hash: h.hash,
		seed: newHashSeed(),
	}
}

// Hash hashes |key|.
func (h Hasher[K]) Hash(key K) uint64 {
	return uint64(h.Hash2(key))
}

// Hash2 hashes |key| as more flexible uintptr.
func (h Hasher[K]) Hash2(key K) uintptr {
	// promise to the compiler that pointer
	// |p| does not escape the stack.

	p := noescape(unsafe.Pointer(&key))
	return h.hash(p, h.seed)
}