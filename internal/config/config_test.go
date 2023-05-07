package config_test

import (
    "os"
    "syscall"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "github.com/spf13/viper"

    "github.com/w6d-io/ci-status/internal/config"
)

var _ = Describe("Config", func() {
    Describe("Parse yaml config", func() {
        Context("Manage Issue", func() {
            BeforeEach(func() {
                config.SkipValidation = true
                viper.Reset()
            })
            It("File does not exist", func() {
                config.CfgFile = "testdata/no-file.yaml"
                config.Init()
            })
            It("Bad syntax in file", func() {
                config.CfgFile = "testdata/file2.yaml"
                config.Init()
            })
            It("kafka webhook entry", func() {
                config.CfgFile = "testdata/file3.yaml"
                config.Init()

            })
            It("webhook bad entry", func() {
                config.CfgFile = "testdata/file4.yaml"
                config.Init()
            })
        })
        Context("Validate config", func() {
            BeforeEach(func() {
                config.SkipValidation = true
                viper.Reset()
            })
            It("File exist", func() {
                config.CfgFile = "testdata/file1.yaml"
                config.Init()
            })
        })
        Context("Validate onChange", func() {
            BeforeEach(func() {
                config.SkipValidation = true
                viper.Reset()
            })
            It("produces a change", func() {
                f, err := os.CreateTemp("testdata", "test*.yaml")
                Expect(err).ToNot(HaveOccurred())
                defer func() {
                    _ = syscall.Unlink(f.Name())
                }()
                data := `---
listen: ":8080"
`
                Expect(os.WriteFile(f.Name(), []byte(data), 0644)).To(Succeed())
                config.CfgFile = f.Name()
                config.Init()
                data += `
timeout: 1m
`
                Expect(os.WriteFile(f.Name(), []byte(data), 0644)).To(Succeed())
            })
        })
    })
})
