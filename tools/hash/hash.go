package hash

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"hash"

	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

type HashType int

const (
	SHA3_256 HashType = iota
	SHA3_512
	SHA256
	SHA512
	Keccak256
	Keccak512
	MD5
	Ripemd160
)

func StringToHashes(data string) (HashType, error) {
	switch data {
	case "sha3-256":
		return SHA3_256, nil
	case "sha3-512":
		return SHA3_512, nil
	case "sha256":
		return SHA256, nil
	case "sha512":
		return SHA512, nil
	case "keccak-256":
		return Keccak256, nil
	case "keccak-512":
		return Keccak512, nil
	case "md5":
		return MD5, nil
	case "ripemd160":
		return Ripemd160, nil
	}
	return -1, NotFound
}

func (h HashType) String() string {
	switch h {
	case SHA3_256:
		return "sha3-256"
	case SHA3_512:
		return "sha3-512"
	case SHA256:
		return "sha256"
	case SHA512:
		return "sha512"
	case Keccak256:
		return "keccak-256"
	case Keccak512:
		return "keccak-512"
	case MD5:
		return "md5"
	case Ripemd160:
		return "ripemd160"
	}

	return ""
}

func (h HashType) NewHash() (func() hash.Hash, error) {
	switch h {
	case SHA3_256:
		return sha3.New256, nil
	case SHA3_512:
		return sha3.New512, nil
	case SHA256:
		return sha256.New, nil
	case SHA512:
		return sha512.New, nil
	case Keccak256:
		return sha3.NewLegacyKeccak256, nil
	case Keccak512:
		return sha3.NewLegacyKeccak512, nil
	case MD5:
		return md5.New, nil
	case Ripemd160:
		return ripemd160.New, nil
	}
	return nil, NotFound
}

var (
	FileNotFound = errors.New("File not found")
	NotFound     = errors.New("Not found")
)
