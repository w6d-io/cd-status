module github.com/w6d-io/ci-status

go 1.15

replace (
	github.com/go-logr/logr => github.com/go-logr/logr v0.2.1
	github.com/go-logr/zapr => github.com/go-logr/zapr v0.3.0
	k8s.io/api => k8s.io/api v0.19.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.0
	k8s.io/client-go => k8s.io/client-go v0.19.0
	k8s.io/client-go/informers => k8s.io/client-go/informers v0.19.0
)

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/go-logr/zapr v0.3.0 // indirect
	github.com/google/go-containerregistry v0.2.1 // indirect
	github.com/tektoncd/pipeline v0.18.1
	go.uber.org/zap v1.16.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	knative.dev/pkg v0.0.0-20210107022335-51c72e24c179 // indirect
	knative.dev/test-infra v0.0.0-20200921012245-37f1a12adbd3 // indirect
	sigs.k8s.io/controller-runtime v0.6.1
)
