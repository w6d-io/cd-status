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
				Expect(config.New("testdata/no-file.yaml")).ToNot(BeNil())
			})
			It("Bad syntax in file", func() {
				Expect(config.New("testdata/file2.yaml")).ToNot(BeNil())
			})
			It("kafka webhook entry", func() {
				Expect(config.New("testdata/file3.yaml")).
					Should(Equal(errors.New("missing topic")))
			})
			It("webhook bad entry", func() {
				Expect(config.New("testdata/file4.yaml").Error()).
					Should(ContainSubstring("first path segment in URL cannot contain colon"))
			})
			It("webhook does not support", func() {
				Expect(config.New("testdata/file5.yaml").Error()).
					Should(ContainSubstring("scheme mongodb not supported"))
			})
		})
		Context("Validate config", func() {
			It("File exist", func() {
				Expect(config.New("testdata/file1.yaml")).To(BeNil())
			})
			It("is auth enabled ?", func() {
				Expect(config.IsAuth()).Should(Equal(true))
			})
			It("get the config struct", func() {
				Expect(config.GetConfig()).ToNot(BeNil())
			})
			It("get listen api address", func() {
				config.SetListen(":8888")
				Expect(config.GetListen()).To(Equal(":8888"))
			})
			It("get auth", func() {
				config.SetAuth([]config.Auth{{"test", "test"}})
				Expect(config.GetAuth()).To(Equal([]config.Auth{{"test", "test"}}))
			})
			It("get webhooks", func() {
				Expect(len(config.GetWebhooks())).To(Equal(2))
			})
		})
	})
	Describe("check tools", func() {
		It("check is in array", func() {
			haystack := []string{"test1", "test2"}
			Expect(config.IsInArray("test3", haystack)).To(Equal(false))
		})
	})
})
