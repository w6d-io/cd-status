package router_test

import (
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/w6d-io/ci-status/pkg/router"
)

var _ = Describe("Auth", func() {
	Context("Handler function", func() {
		It("return the gin handler function", func() {
			auth := router.Auth()
			auth(&gin.Context{})
			Expect(auth).ToNot(BeNil())
		})
	})
})
