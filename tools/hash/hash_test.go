package hash

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"testing"

	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
	"golang.org/x/exp/slices"
)

func TestStringToHashes(t *testing.T) {
	type expectedData struct {
		hashType HashType
		err      error
	}
	type structTestData struct {
		input        string
		expectedData expectedData
	}
	testData := []structTestData{
		{
			input: "sha3-256",
			expectedData: expectedData{
				hashType: SHA3_256,
			},
		},
		{
			input: "sha3-512",
			expectedData: expectedData{
				hashType: SHA3_512,
			},
		},
		{
			input: "sha256",
			expectedData: expectedData{
				hashType: SHA256,
			},
		},
		{
			input: "sha512",
			expectedData: expectedData{
				hashType: SHA512,
			},
		},
		{
			input: "keccak-256",
			expectedData: expectedData{
				hashType: Keccak256,
			},
		},
		{
			input: "keccak-512",
			expectedData: expectedData{
				hashType: Keccak512,
			},
		},
		{
			input: "md5",
			expectedData: expectedData{
				hashType: MD5,
			},
		},
		{
			input: "ripemd160",
			expectedData: expectedData{
				hashType: Ripemd160,
			},
		},

		{
			input: "",
			expectedData: expectedData{
				hashType: -1,
				err:      NotFound,
			},
		},
		{
			input: "abcd",
			expectedData: expectedData{
				hashType: -1,
				err:      NotFound,
			},
		},
	}

	for _, td := range testData {
		h, err := StringToHashes(td.input)
		if h != td.expectedData.hashType {
			t.Errorf("Invalid hash code received, expected %v, received %v", td.expectedData.hashType, h)
		}
		if err != td.expectedData.err {
			t.Errorf("Wrong error received.\nInput string hash: %v;\nExpected hash code: %v;\n Expected err: %v;\n Received hash code: %v;\n Received err: %v.",
				td.input, int(td.expectedData.hashType), td.expectedData.err,
				h, err)
		}
	}
	return
}

func TestString(t *testing.T) {
	type structTestData struct {
		input    HashType
		expected string
	}
	testData := []structTestData{
		{
			input:    SHA3_256,
			expected: "sha3-256",
		},
		{
			input:    SHA3_512,
			expected: "sha3-512",
		},
		{
			input:    SHA256,
			expected: "sha256",
		},
		{
			input:    SHA512,
			expected: "sha512",
		},
		{
			input:    Keccak256,
			expected: "keccak-256",
		},
		{
			input:    Keccak512,
			expected: "keccak-512",
		},
		{
			input:    MD5,
			expected: "md5",
		},
		{
			input:    Ripemd160,
			expected: "ripemd160",
		},
		{
			input:    -1,
			expected: "",
		},
	}

	for _, td := range testData {
		hString := td.input.String()
		if hString != td.expected {
			t.Errorf("Incorrect conversion of hashType to string. Input hashType: %v;\n Expected %v';\n Received %v.", int(td.input), td.expected, hString)
		}
	}
	return
}

type ExpectedDataNewHash struct {
	NewHash func() hash.Hash
	err     error
}

type StructTestDataNewHash struct {
	input        HashType
	expectedData ExpectedDataNewHash
}

func TestNewHash(t *testing.T) {
	testWord := []byte("abcdefghiklmnopqrstvxyz")

	testData := []StructTestDataNewHash{
		{
			input: SHA3_256,
			expectedData: ExpectedDataNewHash{
				NewHash: sha3.New256,
			},
		},
		{
			input: SHA3_512,
			expectedData: ExpectedDataNewHash{
				NewHash: sha3.New512,
			},
		},
		{
			input: SHA256,
			expectedData: ExpectedDataNewHash{
				NewHash: sha256.New,
			},
		},
		{
			input: SHA512,
			expectedData: ExpectedDataNewHash{
				NewHash: sha512.New,
			},
		},
		{
			input: Keccak256,
			expectedData: ExpectedDataNewHash{
				NewHash: sha3.NewLegacyKeccak256,
			},
		},
		{
			input: Keccak512,
			expectedData: ExpectedDataNewHash{
				NewHash: sha3.NewLegacyKeccak512,
			},
		},
		{
			input: MD5,
			expectedData: ExpectedDataNewHash{
				NewHash: md5.New,
			},
		},
		{
			input: Ripemd160,
			expectedData: ExpectedDataNewHash{
				NewHash: ripemd160.New,
			},
		},
		{
			input: -1,
			expectedData: ExpectedDataNewHash{
				err: NotFound,
			},
		},
	}

	for _, td := range testData {
		h, err := td.input.NewHash()
		if err != td.expectedData.err {
			t.Errorf("Wrong error received.\nInput string hash: %v;\nExpected err: %v;\nReceived err: %v.",
				td.input, td.expectedData.err, err)
		} else if err == nil {
			if h == nil && td.expectedData.NewHash != nil {
				t.Errorf("The function must be returned by newHash. Hash name: %v", td.input)
			} else if h != nil && td.expectedData.NewHash == nil {
				t.Errorf("The function must be not returned by newHash. Hash name: %v", td.input)
			} else {
				hasher := td.expectedData.NewHash()
				_, err := hasher.Write(testWord)
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				eHashWord := hasher.Sum(nil)

				hasher = h()
				_, err = hasher.Write(testWord)
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				rHashWords := hasher.Sum(nil)

				if ok := slices.Equal(eHashWord, rHashWords); !ok {
					t.Errorf("Functions generates different responses. Expected: %v, received: %v.", hex.EncodeToString(eHashWord), hex.EncodeToString(rHashWords))
				}
			}
		}
	}
	return
}
