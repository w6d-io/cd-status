package router_test

import (
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/w6d-io/ci-status/internal/config"
	"github.com/w6d-io/ci-status/pkg/router"
	//
	//"github.com/w6d-io/ci-status/pkg/router"
)

var _ = Describe("Router", func() {
	Context("update router engine", func() {
		It("add a post handler", func() {
			router.AddPOST("/test/unit", func(c *gin.Context) {})
		})
	})
	Context("The engine", func() {
		It("run failed to listen", func() {
			config.SetListen("bad-address:ee")
			Expect(router.Run().Error()).To(ContainSubstring("listen tcp:"))
		})
		It("failed on stop", func() {
			Expect(router.Stop()).ToNot(HaveOccurred())
		})
	})
})
