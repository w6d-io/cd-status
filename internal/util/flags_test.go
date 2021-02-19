package util_test

import (
	"errors"
	"flag"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap/zapcore"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/w6d-io/ci-status/internal/util"
)

var _ = Describe("Flags", func() {
	var (
		opts zap.Options
	)
	BeforeEach(func() {
		opts = zap.Options{
			Encoder: zapcore.NewConsoleEncoder(util.TextEncoderConfig()),
		}
	})
	Context("Check function", func() {
		It("JsonEncoderConfig", func() {
			Expect(util.JsonEncoderConfig()).ToNot(BeNil())
		})
		It("TextEncoderConfig", func() {
			Expect(util.TextEncoderConfig()).ToNot(BeNil())
		})
		It("BindFlags", func() {
			util.BindFlags(&opts, flag.CommandLine)
		})
	})
	Context("Check flags methods check", func() {
		configFlag := util.ConfigFlag{}
		When("config flag is used", func() {
			It("Flag is empty", func() {
				Expect(configFlag.Set("")).Should(Equal(errors.New("config cannot be empty")))
			})
			It("it is a directory", func() {
				Expect(configFlag.Set("/tmp")).Should(Equal(errors.New("file /tmp does not exist")))
			})
			It("File does not exist", func() {
				Expect(configFlag.Set("/tmp/no-file.yaml")).Should(Equal(errors.New("file /tmp/no-file.yaml does not exist")))
			})
			It("File has got errors", func() {
				Expect(configFlag.Set("../../test/config/file2.yaml").Error()).Should(ContainSubstring("instanciate config returns "))
			})
			It("File is correct", func() {
				Expect(configFlag.Set("../../test/config/file1.yaml")).To(BeNil())
			})
		})
		levelFlag := util.LevelFlag{}
		levelFlag.ZapOptions = &opts
		When("level flag is used", func() {
			It("Flag is empty", func() {
				Expect(levelFlag.Set("")).Should(Equal(errors.New("log-level cannot be empty")))
			})
			It("invalid string level", func() {
				Expect(levelFlag.Set("no-level")).Should(Equal(errors.New(`invalid log level "no-level"`)))
			})
			It("invalid integer level", func() {
				Expect(levelFlag.Set("-1")).Should(Equal(errors.New(`invalid log level "-1"`)))
			})
			It("valid integer level", func() {
				Expect(levelFlag.Set("1")).To(BeNil())
			})
			It("valid string level", func() {
				Expect(levelFlag.Set("debug")).To(BeNil())
			})
		})
		outputFlag := util.OutputFormatFlag{}
		outputFlag.ZapOptions = &opts
		When("output format flag is used", func() {
			It("Flag is empty", func() {
				Expect(outputFlag.Set("")).Should(Equal(errors.New("log-format cannot be empty")))
			})
			It("invalid format", func() {
				Expect(outputFlag.Set("wrong-format")).Should(Equal(errors.New(`invalid "wrong-format"`)))
			})
			It("valid json format", func() {
				Expect(outputFlag.Set("json")).To(BeNil())
			})
			It("valid text format", func() {
				Expect(outputFlag.Set("text")).To(BeNil())
			})
		})
	})
})
