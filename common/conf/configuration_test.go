package conf

import (
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Encoding の String() が正しいことを確認する
func TestEncoding_String(t *testing.T) {
	testCases := []struct {
		enc      Encoding
		expected string
	}{
		{enc: UTF8, expected: "utf-8"},
		{enc: HEX, expected: "hex"},
	}

	for _, tc := range testCases {
		actual := tc.enc.String()
		if tc.enc.String() != tc.expected {
			t.Errorf("expected: [%s] but got [%s]", tc.expected, actual)
		}
	}
}

func createConfiguration(t *testing.T) *Configuration {
	// プロジェクトのルートディレクトリを指定して設定を取得
	path := filepath.Join("..", "..", "conf")
	c := NewConfiguration("authorizer", "unittest", []string{path})
	err := c.Initialize()
	if err != nil {
		t.Fatalf("reading configuration failed: %s", err)
		return nil
	}
	return c
}

// conf/unittest.yaml から設定値を読み出して assert する
func TestConfigurationLifecycleResource_Get(t *testing.T) {

	t.Run("設定ファイルから設定値が読み取れること", func(t *testing.T) {
		c := createConfiguration(t)

		// int
		actualInt := c.GetInt("unittest.i")
		if actualInt != 10 {
			t.Errorf("expected 10, but got %d", actualInt)
		}
		// string
		actualStr := c.GetString("unittest.s")
		if actualStr != "hoge" {
			t.Errorf("expected hoge but got [%s]", actualStr)
		}
		// slice
		actualStringSlice := c.GetStringSlice("unittest.ss")
		expected := []string{"hoge", "fuga"}
		if !reflect.DeepEqual(actualStringSlice, expected) {
			t.Errorf("expected %s, but got %s", expected, actualStringSlice)
		}
		// bool
		actualBool := c.GetBool("unittest.b")
		expectedBool := true
		if actualBool != expectedBool {
			t.Errorf("expected %t, but got %t", expectedBool, actualBool)
		}
		// byte
		actualByte, err := c.GetByte("unittest.utf8byte", UTF8)
		if err != nil {
			t.Errorf("err should be nil, but got %s", err)
		} else {
			expectedByte := []byte{0x61, 0x62, 0x63, 0x64, 0x65} // abcde
			if !reflect.DeepEqual(actualByte, expectedByte) {
				t.Errorf("expected %x, but got %x", expectedByte, actualByte)
			}
		}
		// duration
		actualDuration := c.GetDuration("unittest.dr")
		expectedDuration, err := time.ParseDuration("1h10m10s")
		if err != nil {
			t.Errorf("error should be nil, but got %s", err)
		}
		if actualDuration != expectedDuration {
			t.Errorf("expected %s, but got %s", expectedDuration, actualDuration)
		}
	})

	t.Run("環境変数で上書きできること", func(t *testing.T) {

		expected := 100
		os.Setenv("AUTHORIZER_PORT", strconv.Itoa(expected))

		c := createConfiguration(t)
		actual := c.GetInt("port")
		if actual != expected {
			t.Errorf("expected %d, but got %d", expected, actual)
		}
	})

}

func TestConfiguration_ReadConfig(t *testing.T) {
	c, err := NewConfigurationFromReader("properties", strings.NewReader("hoge=fuga"))
	if err != nil {
		t.Fatalf("err must be nil, but got %s", err)
	}
	actual := c.GetString("hoge")
	if "fuga" != actual {
		t.Errorf("expected fuga, but got %s", actual)
	}
}
