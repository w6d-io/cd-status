package watch_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/w6d-io/ci-status/pkg/watch"
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
    "name": "pipeline-run-1-1"
  },
  "project_id": 1,
  "pipeline_id": 1,
  "repo_url": " https://github.com/w6d-io/nodejs-sample.git"
}
`
				})
				It("add watcher", func() {
					watch.AddWatcher("test", func(p watch.Payload) error { return nil })
				})
				It("watch for pipelinerun", func() {
					r := ioutil.NopCloser(strings.NewReader(payload))
					w := httptest.NewRecorder()

					c, _ := gin.CreateTestContext(w)
					c.Request = &http.Request{
						Body: framer.NewJSONFramedReader(r),
					}
					watch.Play(c)
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
					watch.Play(c)
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
					watch.Play(c)
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
					watch.Play(c)
				})
			})
		})
	})
})
