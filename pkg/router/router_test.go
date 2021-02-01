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
	BeforeEach(func() {
		gin.SetMode(gin.ReleaseMode)
	})
	It("initialize engine", func() {
		Expect(router.New()).To(BeNil())
	})
	Context("update router engine", func() {
		It("add a post handler", func() {
			router.AddPOST("/test/unit", func(c *gin.Context) {})
		})
		When("Auth is not empty", func() {
			It("initialize engine", func() {
				auth := []config.Auth{{Username: "test", Password: "test"}}
				config.SetAuth(auth)
				Expect(router.New()).To(BeNil())
			})
		})
	})
	Context("The engine", func() {
		It("runs", func() {
			//Expect(router.Run().Error()).To(Equal(""))
		}, 0)
	})
})
