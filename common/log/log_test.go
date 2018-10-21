package log

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kiririmode/grpc-sandbox/common/conf"
	"github.com/sirupsen/logrus"
)

func TestLog_Name(t *testing.T) {
	path := filepath.Join("..", "..", "conf")
	config := conf.NewConfiguration("testapp", "unittest", []string{path})
	log := NewLog(config)
	if log.Name() != "log" {
		t.Errorf("log name must be \"log\"")
	}
}

func TestLog_initializeLogrus(t *testing.T) {

	t.Run("正常系", func(t *testing.T) {
		testCases := []struct {
			config                  string
			expectedIsJSONFormatter bool
			expectedLevel           logrus.Level
		}{
			{
				// json
				config: `log.format=json
			log.output_stdout=false
			log.level=info
			`,
				expectedIsJSONFormatter: true,
				expectedLevel:           logrus.InfoLevel,
			},
			{
				// text
				config: `log.format=text
			log.output_stdout=true
			log.level=debug
			`,
				expectedIsJSONFormatter: false,
				expectedLevel:           logrus.DebugLevel,
			},
		}

		for _, tc := range testCases {
			c, err := conf.NewConfigurationFromReader("properties", strings.NewReader(tc.config))
			if err != nil {
				t.Errorf("err must be nil, but got %s", err)
			}

			log := NewLog(c)
			logger, err := log.initializeLogrus(ioutil.Discard)
			if err != nil {
				t.Errorf("err must be nil, but got %s", err)
			}

			// Formatter が想定通りであること
			_, ok := logger.Formatter.(*logrus.JSONFormatter)
			if ok != tc.expectedIsJSONFormatter {
				t.Errorf("formatter should be json?: expected %t, but got %t", tc.expectedIsJSONFormatter, ok)
			}
			// LogLevel が想定通りであること
			if tc.expectedLevel != logger.GetLevel() {
				t.Errorf("expected log level is %s, but got %s", tc.expectedLevel, logger.GetLevel())
			}
		}
	})

	t.Run("ログフォーマットが対応していない形式", func(t *testing.T) {
		config := "log.format=toml"
		c, err := conf.NewConfigurationFromReader("properties", strings.NewReader(config))
		if err != nil {
			t.Fatalf("err should be nil, but got %s", err)
		}

		log := NewLog(c)
		if _, err := log.initializeLogrus(ioutil.Discard); err == nil {
			t.Error("toml is not supported, but no error occured")
		}
	})

	t.Run("ログレベルが対応していないレベル", func(t *testing.T) {
		config := `log.format=toml
		log.level=hoge`
		c, err := conf.NewConfigurationFromReader("properties", strings.NewReader(config))
		if err != nil {
			t.Fatalf("err should be nil, but got %s", err)
		}

		log := NewLog(c)
		if _, err := log.initializeLogrus(ioutil.Discard); err == nil {
			t.Error("log level hoge is not supported, but no error occured")
		}
	})
}
