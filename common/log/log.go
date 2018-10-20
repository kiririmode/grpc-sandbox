package log

import (
	"io"
	"os"

	"github.com/kiririmode/grpc-sandbox/common/conf"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Log は 本アプリケーションの利用するロギング用クラスを表現する
type Log struct {
	config *conf.Configuration
	rl     *rotatelogs.RotateLogs
	Logger *logrus.Logger
}

// NewLog は、設定 c に基いた新しい Log オブジェクトを返却する
func NewLog(c *conf.Configuration) *Log {
	return &Log{
		config: c,
	}
}

// Name は初期化対象である "log" を返却する
func (l *Log) Name() string {
	return "log"
}

// Initialize はログの初期化を行う
func (l *Log) Initialize() error {

	// rotatelogs の初期化
	rl, err := l.initializeRotateLog()
	if err != nil {
		return errors.Wrap(err, "failed to initialize rotatelogs")
	}

	// logrus の初期化
	logger, err := l.initializeLogrus(rl)
	if err != nil {
		return errors.Wrap(err, "failed to initialize logrus")
	}

	l.rl, l.Logger = rl, logger
	return nil
}

// initializeRotateLog は rotatelog の初期化を行い、そのインスタンスを返却する
func (l *Log) initializeRotateLog() (*rotatelogs.RotateLogs, error) {
	// 古いログのパージ期間、ログローテーションの間隔を設定
	logOption := []rotatelogs.Option{
		rotatelogs.WithRotationCount(l.config.GetInt("log.rotation_counts")),
		rotatelogs.WithRotationTime(l.config.GetDuration("log.rotation_interval")),
	}
	// ローテーション用ログのオブジェクトを作成
	return rotatelogs.New(
		l.config.GetString("log.basename"),
		logOption...,
	)
}

// initializeLogrus は、出力先を writer とした新たな logrus の Logger を作成・返却する。
func (l *Log) initializeLogrus(writer io.Writer) (*logrus.Logger, error) {
	logger := logrus.New()

	// 出力先の設定
	if l.config.GetBool("log.output_stdout") {
		writer = io.MultiWriter(os.Stdout, writer)
	}
	logger.SetOutput(writer)

	// ログフォーマット
	var formatter logrus.Formatter
	f := l.config.GetString("log.format")
	switch f {
	case "json":
		formatter = &logrus.JSONFormatter{}
	case "text":
		formatter = &logrus.TextFormatter{FullTimestamp: true, QuoteEmptyFields: true}
	default:
		return nil, errors.Errorf("illegal log format [%s], specify \"text\" or \"json\" with \"log.format\" key", f)
	}
	logger.SetFormatter(formatter)

	// ログレベル
	v := l.config.GetString("log.level")
	level, err := logrus.ParseLevel(v)
	if err != nil {
		return nil, errors.Errorf("illegal log level [%s]", v)
	}
	logger.SetLevel(level)

	return logger, nil
}

func (l *Log) Finalize() error {
	err := l.rl.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close rotatelog")
	}
	return nil
}
