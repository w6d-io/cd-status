package play_test

import (
	"errors"
	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-status/pkg/handler/watch/play"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/util/framer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

var _ = Describe("Watch", func() {
	Describe("a payload has been received", func() {
		When("Payload is correct", func() {
			Context("Pipelinerun scan", func() {
				var (
					payload string
				)
				It("watch for pipelinerun", func() {
					payload = `
{
  "object": {
    "kind": "pipelinerun",
    "namespaced_name": {
      "name": "pipeline-run-1-1",
      "namespace": "default"
    }
  },
  "project_id": 1,
  "pipeline_id": 1,
  "repo_url": " https://github.com/w6d-io/nodejs-sample.git"
}
`
					f := func(_ logr.Logger, _ types.NamespacedName, _ int64, _ int64, _ string, _ string, _ string, _ string) error {
						return nil
					}
					play.AddWatcher("test", f)
					r := ioutil.NopCloser(strings.NewReader(payload))
					w := httptest.NewRecorder()
					c, _ := gin.CreateTestContext(w)
					c.Request = &http.Request{
						Body: framer.NewJSONFramedReader(r),
					}
					play.Play(c)
					Expect(c.Writer.Status()).To(Equal(200))
				})
				It("scan failed", func() {
					payload = `
{
  "object": {
    "kind": "pipelinerun",
    "namespaced_name": {
      "name": "pipeline-run-1-1",
      "namespace": "default"
    }
  },
  "project_id": 1,
  "pipeline_id": 1,
  "repo_url": " https://github.com/w6d-io/nodejs-sample.git"
}
`
					f := func(_ logr.Logger, _ types.NamespacedName, _ int64, _ int64, _ string, _ string, _ string, _ string) error {
						return errors.New("test")
					}
					play.AddWatcher("test", f)
					r := ioutil.NopCloser(strings.NewReader(payload))
					w := httptest.NewRecorder()
					c, _ := gin.CreateTestContext(w)
					c.Request = &http.Request{
						Body: framer.NewJSONFramedReader(r),
					}
					play.Play(c)
					Expect(c.Writer.Status()).To(Equal(200))
				})
			})
		})
		When("Payload is not correct", func() {
			Context("pipelinerun scan", func() {
				It("payload badly formatted", func() {
					payload := `
{
  "object": {
    "kind": "pipelinerun",
	"namespaced_name"
    "name": "pipeline-run-1-1"
  },
  "project_id": 1,
  "pipeline_id": 1,
  "repo_url": " https://github.com/w6d-io/nodejs-sample.git",
`
					f := func(_ logr.Logger, _ types.NamespacedName, _ int64, _ int64, _ string, _ string, _ string, _ string) error {
						return nil
					}
					play.AddWatcher("test", f)
					r := ioutil.NopCloser(strings.NewReader(payload))
					w := httptest.NewRecorder()
					c, _ := gin.CreateTestContext(w)
					c.Request = &http.Request{
						Body: framer.NewJSONFramedReader(r)}
					play.Play(c)
					Expect(c.Writer.Status()).To(Equal(400))
				})
				It("payload object kind is not supported", func() {
					payload := `
{
  "object": {
    "kind": "toto",
    "name": "pipeline-run-1-1"
  },
  "project_id": 1,
  "pipeline_id": 1,
  "repo_url": " https://github.com/w6d-io/nodejs-sample.git"
}
`
					f := func(_ logr.Logger, _ types.NamespacedName, _ int64, _ int64, _ string, _ string, _ string, _ string) error {
						return nil
					}
					play.AddWatcher("test", f)
					r := ioutil.NopCloser(strings.NewReader(payload))
					w := httptest.NewRecorder()
					c, _ := gin.CreateTestContext(w)
					c.Request = &http.Request{
						Body: framer.NewJSONFramedReader(r)}
					play.Play(c)
					Expect(c.Writer.Status()).To(Equal(402))
				})
				It("Scan return an error", func() {
					payload := `
{
  "object": {
    "kind": "pipelinerun",
    "name": "pipeline-run-1-1"
  },
  "project_id": 1,
  "pipeline_id": 1,
  "repo_url": " https://github.com/w6d-io/nodejs-sample.git"
}
`
					f := func(_ logr.Logger, _ types.NamespacedName, _ int64, _ int64, _ string, _ string, _ string, _ string) error {
						return errors.New("test")
					}
					play.AddWatcher("test", f)
					r := ioutil.NopCloser(strings.NewReader(payload))
					w := httptest.NewRecorder()
					c, _ := gin.CreateTestContext(w)
					c.Request = &http.Request{
						Body: framer.NewJSONFramedReader(r)}
					play.Play(c)
					Expect(c.Writer.Status()).To(Equal(200))
				})
				It("name or namespace missing handler return an error", func() {
					payload := `
{
  "object": {
    "kind": "pipelinerun",
	"namespaced_name": {
      "name": "",
      "namespace": ""
    }
  },
  "project_id": 1,
  "pipeline_id": 1,
  "repo_url": " https://github.com/w6d-io/nodejs-sample.git"
}
`
					f := func(_ logr.Logger, _ types.NamespacedName, _ int64, _ int64, _ string, _ string, _ string, _ string) error {
						return nil
					}
					play.AddWatcher("test", f)
					r := ioutil.NopCloser(strings.NewReader(payload))
					w := httptest.NewRecorder()
					c, _ := gin.CreateTestContext(w)
					c.Request = &http.Request{
						Body: framer.NewJSONFramedReader(r)}
					play.Play(c)
					Expect(c.Writer.Status()).To(Equal(401))
				})
			})
		})
	})
})
