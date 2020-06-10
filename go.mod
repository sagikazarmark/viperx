module github.com/sagikazarmark/viperx

go 1.12

require (
	emperror.dev/errors v0.7.0
	github.com/banzaicloud/bank-vaults/pkg/sdk v0.3.1
	github.com/hashicorp/vault/api v1.0.4
	github.com/spf13/viper v1.7.0
)

replace k8s.io/client-go => k8s.io/client-go v0.17.2
