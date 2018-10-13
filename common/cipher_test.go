package common

import (
	"crypto/aes"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"
)

func TestAesCbcPkcs7CipherImplementsBlockCipher(t *testing.T) {
	var _ BlockCipher = &AesCbcPkcs7Cipher{}
}

func TestAesCbcPkcs7Encrypt(t *testing.T) {
	key, _ := hex.DecodeString("1234567890123456789012345678901234567890123456789012345678901234")
	iv, _ := hex.DecodeString("1234567890ABCDEF1234567890ABCDEF")
	sut, err := NewAesCbcPkcs7Cipher(key, iv)
	if err != nil {
		t.Errorf("error must be nil, but [%s]", err)
	}

	// expected は以下のコマンドにより生成
	// echo -n plain | openssl aes-256-cbc -e -base64 -iv 1234567890ABCDEF1234567890ABCDEF -K 1234567890123456789012345678901234567890123456789012345678901234
	testCases := []struct {
		plain    string
		expected string
	}{
		{plain: "a", expected: "YPID0ng/IBlB2BS1fyya+w=="},
		{plain: "aa", expected: "LJ0dXyUo/1Nl37UlXeQG4w=="},
		{plain: "aaa", expected: "4MAA3s1H1h/odEihUjQFVQ=="},
		{plain: "aaaa", expected: "nvGPv/rJR+REVIg/OtYVxA=="},
		{plain: "aaaaa", expected: "KKRUUCl9PtFeihcOPNApsQ=="},
		{plain: "aaaaaa", expected: "KoJurpbPYgJaza9gT1kUBQ=="},
		{plain: "aaaaaaa", expected: "G47q0vrzjecWzpsI2yDdWQ=="},
		{plain: "aaaaaaaa", expected: "IFMMDMct88+/ZUAHcRAfug=="},
		{plain: "aaaaaaaaa", expected: "37ZNMfsH1+foemy0cEZLug=="},
		{plain: "aaaaaaaaaa", expected: "QkasonIVbK9QhGx1K5ZqwA=="},
		{plain: "aaaaaaaaaaa", expected: "7u3ohko2CyRuuqoUQGbs+Q=="},
		{plain: "aaaaaaaaaaaa", expected: "NrUopgHNi+NYmJq7wk217A=="},
		{plain: "aaaaaaaaaaaaa", expected: "QnoQppDKNWXeKhjf6odl9g=="},
		{plain: "aaaaaaaaaaaaaa", expected: "QUVsN+uL9cuZsHKlmbaw9g=="},
		{plain: "aaaaaaaaaaaaaaa", expected: "ovCBJo5FHvnO6iHJKyribQ=="},
		{plain: "aaaaaaaaaaaaaaaa", expected: "lnWGbKBVfS+/qdgIosh+VdGtkqIQmQDyBJLkEqnU0+k="},
	}

	for _, tc := range testCases {
		encrypted, err := sut.Encrypt([]byte(tc.plain))
		if err != nil {
			t.Errorf("error must be nil, but [%s]", err)
		}
		actual := base64.StdEncoding.EncodeToString(encrypted)
		if tc.expected != actual {
			t.Errorf("expected [%s], but got [%s]", tc.expected, actual)
		}
	}
}

func TestAesCbcPkcs7Decrypt(t *testing.T) {
	key, _ := hex.DecodeString("1234567890123456789012345678901234567890123456789012345678901234")
	iv, _ := hex.DecodeString("1234567890ABCDEF1234567890ABCDEF")
	sut, err := NewAesCbcPkcs7Cipher(key, iv)
	if err != nil {
		t.Errorf("error must be nil, but [%s]", err)
	}

	// expected は以下のコマンドにより生成
	// echo -n plain | openssl aes-256-cbc -e -base64 -iv 1234567890ABCDEF1234567890ABCDEF -K 1234567890123456789012345678901234567890123456789012345678901234
	testCases := []struct {
		plain    string
		expected string
	}{
		{plain: "a", expected: "YPID0ng/IBlB2BS1fyya+w=="},
		{plain: "aa", expected: "LJ0dXyUo/1Nl37UlXeQG4w=="},
		{plain: "aaa", expected: "4MAA3s1H1h/odEihUjQFVQ=="},
		{plain: "aaaa", expected: "nvGPv/rJR+REVIg/OtYVxA=="},
		{plain: "aaaaa", expected: "KKRUUCl9PtFeihcOPNApsQ=="},
		{plain: "aaaaaa", expected: "KoJurpbPYgJaza9gT1kUBQ=="},
		{plain: "aaaaaaa", expected: "G47q0vrzjecWzpsI2yDdWQ=="},
		{plain: "aaaaaaaa", expected: "IFMMDMct88+/ZUAHcRAfug=="},
		{plain: "aaaaaaaaa", expected: "37ZNMfsH1+foemy0cEZLug=="},
		{plain: "aaaaaaaaaa", expected: "QkasonIVbK9QhGx1K5ZqwA=="},
		{plain: "aaaaaaaaaaa", expected: "7u3ohko2CyRuuqoUQGbs+Q=="},
		{plain: "aaaaaaaaaaaa", expected: "NrUopgHNi+NYmJq7wk217A=="},
		{plain: "aaaaaaaaaaaaa", expected: "QnoQppDKNWXeKhjf6odl9g=="},
		{plain: "aaaaaaaaaaaaaa", expected: "QUVsN+uL9cuZsHKlmbaw9g=="},
		{plain: "aaaaaaaaaaaaaaa", expected: "ovCBJo5FHvnO6iHJKyribQ=="},
		{plain: "aaaaaaaaaaaaaaaa", expected: "lnWGbKBVfS+/qdgIosh+VdGtkqIQmQDyBJLkEqnU0+k="},
	}

	for _, tc := range testCases {
		decoded, err := base64.StdEncoding.DecodeString(tc.expected)
		if err != nil {
			t.Errorf("err must be nil, but [%s]", err)
		}
		plain, err := sut.Decrypt(decoded)
		if err != nil {
			t.Errorf("err must be nil, but [%s]", err)
		}
		plainText := string(plain)

		if plainText != tc.plain {
			t.Errorf("expected [%s], but got [%s]", tc.plain, plainText)
		}
	}
}

func TestIllegalLength(t *testing.T) {
	t.Run("illegal key length", func(t *testing.T) {
		iv, _ := hex.DecodeString("1234567890ABCDEF1234567890ABCDEF")
		for l := 1; l < 32; l++ {
			key, err := hex.DecodeString(strings.Repeat("00", l))
			if err != nil {
				t.Errorf("err must be nil, but [%s]", err)
			}

			_, err = NewAesCbcPkcs7Cipher(key, iv)
			switch l {
			case 16, 24, 32:
				if err != nil {
					t.Errorf("must succeed with key length [%d]byte, but got %s", l, err)
				}
			default:
				if err == nil {
					t.Errorf("must fail with key length [%d]byte", l)
				}
			}
		}
	})
	t.Run("illegal iv length", func(t *testing.T) {
		key, _ := hex.DecodeString("1234567890123456789012345678901234567890123456789012345678901234")
		for i := 1; i < aes.BlockSize; i++ {
			iv, err := hex.DecodeString(strings.Repeat("00", i))
			if err != nil {
				t.Errorf("err must be nil, but [%s]", err)
			}
			_, err = NewAesCbcPkcs7Cipher(key, iv)
			if err == nil {
				t.Errorf("must fail with invalid initial vector length [%d] byte", len(iv))
			}
		}
	})
}
