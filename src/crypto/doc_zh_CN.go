// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package crypto collects common cryptographic constants.

// Package crypto collects common cryptographic constants.
package crypto

import (
    "hash"
    "io"
    "strconv"
)


const (
	MD4        Hash = 1 + iota // import golang.org/x/crypto/md4
	MD5                        // import crypto/md5
	SHA1                       // import crypto/sha1
	SHA224                     // import crypto/sha256
	SHA256                     // import crypto/sha256
	SHA384                     // import crypto/sha512
	SHA512                     // import crypto/sha512
	MD5SHA1                    // no implementation; MD5+SHA1 used for TLS RSA
	RIPEMD160                  // import golang.org/x/crypto/ripemd160
	SHA3_224                   // import golang.org/x/crypto/sha3
	SHA3_256                   // import golang.org/x/crypto/sha3
	SHA3_384                   // import golang.org/x/crypto/sha3
	SHA3_512                   // import golang.org/x/crypto/sha3
	SHA512_224                 // import crypto/sha512
	SHA512_256                 // import crypto/sha512

)


// Decrypter is an interface for an opaque private key that can be used for
// asymmetric decryption operations. An example would be an RSA key
// kept in a hardware module.
type Decrypter interface {
	// Public returns the public key corresponding to the opaque,
	// private key.
	Public() PublicKey

	// Decrypt decrypts msg. The opts argument should be appropriate for
	// the primitive used. See the documentation in each implementation for
	// details.
	Decrypt(rand io.Reader, msg []byte, opts DecrypterOpts) (plaintext []byte, err error)
}



type DecrypterOpts interface{}


// Hash identifies a cryptographic hash function that is implemented in another
// package.
type Hash uint


// PrivateKey represents a private key using an unspecified algorithm.
type PrivateKey interface{}


// PublicKey represents a public key using an unspecified algorithm.
type PublicKey interface{}


// Signer is an interface for an opaque private key that can be used for
// signing operations. For example, an RSA key kept in a hardware module.
type Signer interface {
	// Public returns the public key corresponding to the opaque,
	// private key.
	Public() PublicKey

	// Sign signs digest with the private key, possibly using entropy from
	// rand. For an RSA key, the resulting signature should be either a
	// PKCS#1 v1.5 or PSS signature (as indicated by opts). For an (EC)DSA
	// key, it should be a DER-serialised, ASN.1 signature structure.
	//
	// Hash implements the SignerOpts interface and, in most cases, one can
	// simply pass in the hash function used as opts. Sign may also attempt
	// to type assert opts to other types in order to obtain algorithm
	// specific values. See the documentation in each package for details.
	//
	// Note that when a signature of a hash of a larger message is needed,
	// the caller is responsible for hashing the larger message and passing
	// the hash (as digest) and the hash function (as opts) to Sign.
	Sign(rand io.Reader, digest []byte, opts SignerOpts) (signature []byte, err error)
}


// SignerOpts contains options for signing with a Signer.
type SignerOpts interface {
	// HashFunc returns an identifier for the hash function used to produce
	// the message passed to Signer.Sign, or else zero to indicate that no
	// hashing was done.
	HashFunc() Hash
}


// RegisterHash registers a function that returns a new instance of the given
// hash function. This is intended to be called from the init function in
// packages that implement hash functions.
func RegisterHash(h Hash, f func() hash.Hash)

// Available reports whether the given hash function is linked into the binary.
func (Hash) Available() bool

// HashFunc simply returns the value of h so that Hash implements SignerOpts.
func (Hash) HashFunc() Hash

// New returns a new hash.Hash calculating the given hash function. New panics
// if the hash function is not linked into the binary.
func (Hash) New() hash.Hash

// Size returns the length, in bytes, of a digest resulting from the given hash
// function. It doesn't require that the hash function in question be linked
// into the program.
func (Hash) Size() int

