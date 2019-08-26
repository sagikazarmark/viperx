module github.com/sagikazarmark/viperx

go 1.12

require (
	emperror.dev/errors v0.4.2
	github.com/banzaicloud/bank-vaults/pkg/sdk v0.1.2
	github.com/hashicorp/vault/api v1.0.4
	github.com/spf13/viper v1.4.0
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20181213150558-05914d821849
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20181127025237-2b1284ed4c93
	k8s.io/client-go => k8s.io/client-go v2.0.0-alpha.0.0.20181213151034-8d9ed539ba31+incompatible
)
