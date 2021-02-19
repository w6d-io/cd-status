package router_test

import (
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/w6d-io/ci-status/pkg/router"
	"net/http"
)

var _ = Describe("Log", func() {
	var (
		c *gin.Context
	)
	Describe("For gin gonic", func() {
		Context("log output", func() {
			outputFunc := router.LogMiddleware()
			It("gin output handlerFunc", func() {
				Expect(outputFunc).ToNot(BeNil())
			})
			Context("Correlation ID", func() {
				correlationID := router.CorrelationID()
				It("gin id header handlerFunc", func() {
					Expect(correlationID).ToNot(BeNil())
				})
				When("is not empty", func() {
					It("set the correlation id in the http header", func() {
						c2 := &gin.Context{Request: &http.Request{Header: http.Header{}}}
						cr := router.CorrelationID()
						cr(c2)
					})
				})
			})
			Context("Get client ip address", func() {
				c = &gin.Context{Request: &http.Request{Header: http.Header{}}}
				c.Request.RemoteAddr = "10.0.0.2,10.0.0.3"
				When("X-Real-IP is not set", func() {
					It("return the remote address", func() {
						Expect(router.GetClientIP(c)).Should(Equal("10.0.0.2"))
					})
				})
				When("X-Real-IP is set", func() {
					It("returns the ip address", func() {
						c.Request.Header.Set("X-Real-IP", "10.0.0.1")
						Expect(router.GetClientIP(c)).Should(Equal("10.0.0.1"))
					})
				})
			})
			Context("Gin handler function", func() {
				It("Json log", func() {
					jsonLog := router.LogMiddleware()
					c.Request.Method = "POST"
					jsonLog(c)
				})
			})
		})
	})
})
