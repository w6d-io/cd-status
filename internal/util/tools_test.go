package util_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/w6d-io/ci-status/internal/util"
)

var _ = Describe("Tools", func() {
	Context("test tools", func() {
		It("RemoveIndex", func() {
			arr := []string{"A", "B", "C", "D"}
			Expect(util.RemoveIndex(arr, 2)).To(Equal([]string{"A", "B", "D"}))
		})
	})
})
