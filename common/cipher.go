package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"

	"github.com/pkg/errors"
)

// BlockCipher はブロック暗号を表現する
type BlockCipher interface {
	Encrypt(plain []byte) ([]byte, error)
	Decrypt(encrypted []byte) ([]byte, error)
}

// AesCbcPkcs7Cipher はAES/CBC/PKCS7 のブロック暗号を表現する
type AesCbcPkcs7Cipher struct {
	// 初期ベクトル
	initialVector []byte
	// ブロック暗号
	block cipher.Block
}

// NewAesCbcPkcs7Cipher は AES/CBC/PKCS#7 のブロック暗号を作成し、返却する
func NewAesCbcPkcs7Cipher(key, iv []byte) (*AesCbcPkcs7Cipher, error) {
	// 鍵長チェック
	keyLen := len(key)
	if (keyLen != 16) && (keyLen != 24) && (keyLen != 32) {
		return nil, errors.Errorf("illegal key length [%d]. key length for AES must be 128, 192, 256 bit", keyLen)
	}
	// 初期ベクトル長チェック
	if len(iv) != aes.BlockSize {
		return nil, errors.Errorf("illegal initial vector size [%d]byte. initial vector size must be [%d]byte", len(iv), aes.BlockSize)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create AES cipher block")
	}
	return &AesCbcPkcs7Cipher{
		initialVector: iv,
		block:         block,
	}, nil
}

// pad は RFC 5652 6.3. Content-encryption Process に記述された通りに
// b にパディングとしてのバイトを追加する (PKCS#7 Padding)
func (c *AesCbcPkcs7Cipher) pad(b []byte) []byte {
	padSize := aes.BlockSize - (len(b) % aes.BlockSize)
	pad := bytes.Repeat([]byte{byte(padSize)}, padSize)
	return append(b, pad...)
}

// unpad は PKCS#7 Padding に従って付与されたパディングを削除する
func (c *AesCbcPkcs7Cipher) unpad(b []byte) []byte {
	padSize := int(b[len(b)-1])
	return b[:len(b)-padSize]
}

// Encrypt は plain を AES/CBC/PKCS#7 で暗号化する。
func (c *AesCbcPkcs7Cipher) Encrypt(plain []byte) ([]byte, error) {
	encrypter := cipher.NewCBCEncrypter(c.block, c.initialVector)

	// PKCS#7 に沿ってパディングを付与
	padded := c.pad(plain)
	// 暗号化
	encrypted := make([]byte, len(padded))
	encrypter.CryptBlocks(encrypted, padded)
	return encrypted, nil
}

// Decrypt は encrypted を AES/CBC/PKCS#7 で復号化する
func (c *AesCbcPkcs7Cipher) Decrypt(encrypted []byte) ([]byte, error) {
	mode := cipher.NewCBCDecrypter(c.block, c.initialVector)

	plain := make([]byte, len(encrypted))
	mode.CryptBlocks(plain, encrypted)
	// パディングを除去
	return c.unpad(plain), nil
}
