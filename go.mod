module github.com/sagikazarmark/viperx

go 1.12

require (
	emperror.dev/errors v0.6.0
	github.com/banzaicloud/bank-vaults/pkg/sdk v0.2.1
	github.com/hashicorp/vault/api v1.0.4
	github.com/spf13/viper v1.4.0
)

replace k8s.io/client-go => k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90
