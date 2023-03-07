module github.com/w6d-io/ci-status

go 1.15

replace (
	github.com/go-logr/logr => github.com/go-logr/logr v0.4.0
	github.com/go-logr/zapr => github.com/go-logr/zapr v0.4.0
	k8s.io/api => k8s.io/api v0.19.7
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.7
	k8s.io/client-go => k8s.io/client-go v0.19.7
	k8s.io/client-go/informers => k8s.io/client-go/informers v0.19.7
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.8.3
)

require (
	github.com/avast/retry-go v3.0.0+incompatible
	github.com/gin-gonic/gin v1.7.7
	github.com/go-logr/logr v0.4.0
	github.com/go-logr/zapr v0.4.0 // indirect
	github.com/google/uuid v1.2.0
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/onsi/ginkgo v1.15.0
	github.com/onsi/gomega v1.10.5
	github.com/tektoncd/cli v0.16.0
	github.com/tektoncd/pipeline v0.21.0
	github.com/w6d-io/hook v0.1.2
	go.uber.org/zap v1.16.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	knative.dev/pkg v0.0.0-20210219034502-29092fe5b88d
	sigs.k8s.io/controller-runtime v0.8.3
)
