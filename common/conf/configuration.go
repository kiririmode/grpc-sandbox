package conf

import (
	"encoding/hex"
	"io"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Configuration は、アプリケーションの持つ設定を表現する
type Configuration struct {
	// アプリケーション名
	AppName string
	// 環境名
	EnvironmentName string
	// 設定ファイルの探索パス
	paths []string
	// Viper のインスタンス
	viper *viper.Viper
}

// Encoding は設定ファイルの文字列をバイト化するときのエンコーディングを表現する
type Encoding int

const (
	// UTF8 は utf-8 encoding を表現する
	UTF8 Encoding = iota
	// HEX は hex encoding を表現する
	HEX
)

// String は Encoding の名称を返却する
func (e Encoding) String() string {
	switch e {
	case UTF8:
		return "utf-8"
	case HEX:
		return "hex"
	}
	return "unknown"
}

// NewConfiguration は、envname という環境名に紐付く新しい設定を返却する。
// appName はアプリケーション名で、設定を上書きするための環境変数の接頭語を定義する。
// searchPaths は、設定ファイルを探索するパスを表現する。
func NewConfiguration(appName, envname string, searchPaths []string) *Configuration {
	return &Configuration{
		AppName:         appName,
		EnvironmentName: envname,
		paths:           searchPaths,
		viper:           viper.New(),
	}
}

// NewConfigurationFromReader は、in から読み込んだ内容に紐付く新しい設定を返却する。
func NewConfigurationFromReader(format string, in io.Reader) (*Configuration, error) {
	v := viper.New()
	v.SetConfigType(format)
	if err := v.ReadConfig(in); err != nil {
		return nil, errors.Wrapf(err, "failed to init Configuration with %s", in)
	}
	return &Configuration{viper: v}, nil
}

// Name は初期化対象である "configuration" を返却する
func (c *Configuration) Name() string {
	return "configuration"
}

// Initialize は、設定を読み込み、利用できるようにする
func (c *Configuration) Initialize() error {
	if strings.TrimSpace(c.EnvironmentName) == "" {
		return errors.Errorf("environment name is missing")
	}

	// 設定ファイルを探索するパスを設定する
	c.viper.SetConfigName(c.EnvironmentName)
	for _, path := range c.paths {
		c.viper.AddConfigPath(path)
	}

	// 設定ファイルを読み込み
	err := c.viper.ReadInConfig()
	if err != nil {
		return errors.Wrapf(err, "failed to read config file: [%s] (suffix ommitted)", c.EnvironmentName)
	}

	// アプリケーション名を接頭語として付与した環境変数を設定することで設定を上書きできるようにする
	c.viper.SetEnvPrefix(c.AppName)
	c.viper.AutomaticEnv()
	// ネストした設定項目も環境変数で上書きできるようにする
	// ex.) database.host => AUTHORIZER_DATABASE_HOST
	c.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return nil
}

// Finalize は何も実行しない
func (c *Configuration) Finalize() error {
	return nil
}

// GetInt は、key に対応する設定値を int で返却する
func (c *Configuration) GetInt(key string) int {
	return c.viper.GetInt(key)
}

// GetString は、key に対応する設定値を string で返却する
func (c *Configuration) GetString(key string) string {
	return c.viper.GetString(key)
}

// GetStringSlice は、key に対応する設定値を string のスライスで返却する
func (c *Configuration) GetStringSlice(key string) []string {
	return c.viper.GetStringSlice(key)
}

// GetBool は、key に対応する設定値を bool で返却する
func (c *Configuration) GetBool(key string) bool {
	return c.viper.GetBool(key)
}

// GetByte は、key に対応する設定値(string) を enc で
// 表現されるエンコーディングでデコードし、その結果としての byte スライスを返却する。
func (c *Configuration) GetByte(key string, enc Encoding) (b []byte, err error) {
	v := c.viper.GetString(key)

	switch enc {
	case UTF8:
		return []byte(v), nil
	case HEX:
		b, err = hex.DecodeString(v)
	default:
		return nil, errors.Errorf("unsupported encoding [%s]", enc)
	}

	return b, err
}

// GetDuration は、key に対応する設定値を Duration として返却する
func (c *Configuration) GetDuration(key string) time.Duration {
	return c.viper.GetDuration(key)
}

// SetFormat は設定ファイルのフォーマットを指定する。
// 設定に使用しているライブラリである viper は自動的にフォーマットを検知してくれるので、
// 本メソッドは主としてテスト用である。
func (c *Configuration) SetFormat(formatType string) {
	c.viper.SetConfigType(formatType)
}

// ReadConfig は reader から設定を読み込む。
func (c *Configuration) ReadConfig(in io.Reader) error {
	return c.viper.ReadConfig(in)
}
