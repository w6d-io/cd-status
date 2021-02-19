package play_test

import (
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
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

var _ = Describe("Watch", func() {
	Describe("a payload has been received", func() {
		Context("Payload is correct", func() {
			Context("Pipelinerun scan", func() {
				var (
					payload string
				)
				BeforeEach(func() {
					payload = `
{
  "object": {
    "kind": "pipelinerun",
    namespaced_name: {
      "name": "pipeline-run-1-1",
      "namespace": "default"
    }
  },
  "project_id": 1,
  "pipeline_id": 1,
  "repo_url": " https://github.com/w6d-io/nodejs-sample.git"
}
`
				})
				It("add watcher", func() {
					f := func(_ logr.Logger, _ types.NamespacedName, _ int64, _ int64) error { return nil }
					play.AddWatcher("test", f)
				})
				It("watch for pipelinerun", func() {
					r := ioutil.NopCloser(strings.NewReader(payload))
					w := httptest.NewRecorder()

					c, _ := gin.CreateTestContext(w)
					c.Request = &http.Request{
						Body: framer.NewJSONFramedReader(r),
					}
					play.Play(c)
				})
			})
		})
		Context("Payload is not correct", func() {
			Context("pipelinerun scan", func() {
				It("payload badly formatted", func() {
					payload := `
{
  "object": {
    "kind": "pipelinerun",
    "name": "pipeline-run-1-1"
  },
  "project_id": 1,
  "pipeline_id": 1,
  "repo_url": " https://github.com/w6d-io/nodejs-sample.git",
`
					r := ioutil.NopCloser(strings.NewReader(payload))
					w := httptest.NewRecorder()
					c, _ := gin.CreateTestContext(w)
					c.Request = &http.Request{
						Body: framer.NewJSONFramedReader(r)}
					play.Play(c)
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
					r := ioutil.NopCloser(strings.NewReader(payload))
					w := httptest.NewRecorder()
					c, _ := gin.CreateTestContext(w)
					c.Request = &http.Request{
						Body: framer.NewJSONFramedReader(r)}
					play.Play(c)
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
					r := ioutil.NopCloser(strings.NewReader(payload))
					w := httptest.NewRecorder()
					c, _ := gin.CreateTestContext(w)
					c.Request = &http.Request{
						Body: framer.NewJSONFramedReader(r)}
					play.Play(c)
				})
			})
		})
	})
})
