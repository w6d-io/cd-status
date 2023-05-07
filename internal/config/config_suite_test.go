package config_test

import (
    "testing"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "github.com/spf13/viper"
    zapraw "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func TestConfig(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Config Suite")
}

var _ = BeforeSuite(func() {
    encoder := zapcore.EncoderConfig{
        // Keys can be anything except the empty string.
        TimeKey:        "T",
        LevelKey:       "L",
        NameKey:        "N",
        CallerKey:      "C",
        MessageKey:     "M",
        StacktraceKey:  "S",
        LineEnding:     zapcore.DefaultLineEnding,
        EncodeLevel:    zapcore.CapitalLevelEncoder,
        EncodeTime:     zapcore.ISO8601TimeEncoder,
        EncodeDuration: zapcore.StringDurationEncoder,
        EncodeCaller:   zapcore.FullCallerEncoder,
    }
    opts := zap.Options{
        Encoder:         zapcore.NewConsoleEncoder(encoder),
        Development:     true,
        StacktraceLevel: zapcore.PanicLevel,
    }
    ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts), zap.RawZapOpts(zapraw.AddCaller(), zapraw.AddCallerSkip(-2))))

    viper.SetConfigType("yaml")
}, 60)
