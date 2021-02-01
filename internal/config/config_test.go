package config_test

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/w6d-io/ci-status/internal/config"
)

var _ = Describe("Config", func() {
	Describe("Parse yaml config", func() {
		Context("Manage Issue", func() {
			It("File does not exist", func() {
				Expect(config.New("../../test/config/no-file.yaml")).ToNot(BeNil())
			})
			It("Bad syntax in file", func() {
				Expect(config.New("../../test/config/file2.yaml")).ToNot(BeNil())
			})
			It("kafka webhook entry", func() {
				Expect(config.New("../../test/config/file3.yaml")).
					Should(Equal(errors.New("missing topic")))
			})
			It("webhook bad entry", func() {
				Expect(config.New("../../test/config/file4.yaml").Error()).
					Should(ContainSubstring("first path segment in URL cannot contain colon"))
			})
		})
		Context("Validate config", func() {
			It("File exist", func() {
				Expect(config.New("../../test/config/file1.yaml")).To(BeNil())
			})
			It("is auth enabled ?", func() {
				Expect(config.IsAuth()).Should(Equal(true))
			})
			It("get the config struct", func() {
				Expect(config.GetConfig()).ToNot(BeNil())
			})
			It("get listen api address", func() {
				Expect(config.GetListen()).To(Equal(":8080"))
			})
			It("get auth", func() {
				config.SetAuth([]config.Auth{{"test", "test"}})
				Expect(config.GetAuth()).To(Equal([]config.Auth{{"test", "test"}}))
			})
		})
	})
})
