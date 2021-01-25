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
	github.com/go-openapi/spec v0.20.2 // indirect
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/google/go-containerregistry v0.2.1 // indirect
	github.com/google/uuid v1.1.2
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	github.com/tektoncd/cli v0.15.0
	github.com/tektoncd/pipeline v0.19.0
	github.com/ugorji/go v1.2.3 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c // indirect
	golang.org/x/tools v0.1.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	knative.dev/pkg v0.0.0-20210107022335-51c72e24c179 // indirect
	knative.dev/test-infra v0.0.0-20200921012245-37f1a12adbd3 // indirect
	sigs.k8s.io/controller-runtime v0.6.1
)
