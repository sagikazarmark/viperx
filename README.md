# ViperX: [Viper](https://github.com/spf13/viper) extensions

[![CircleCI](https://circleci.com/gh/sagikazarmark/viperx.svg?style=svg)](https://circleci.com/gh/sagikazarmark/viperx)
[![Go Version](https://img.shields.io/badge/go%20version-%3E=1.12-orange.svg?style=flat-square)](https://github.com/sagikazarmark/viperx)
[![Go Report Card](https://goreportcard.com/badge/github.com/sagikazarmark/viperx?style=flat-square)](https://goreportcard.com/report/github.com/sagikazarmark/viperx)
[![GolangCI](https://golangci.com/badges/github.com/sagikazarmark/viperx.svg)](https://golangci.com/r/github.com/sagikazarmark/viperx)
[![GoDoc](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/sagikazarmark/viperx)


## Installation

```bash
$ go get github.com/sagikazarmark/viperx
```

## Usage

### Remote config provider registry

Package `remote` provides a remote provider registry.

```go
package main

import (
    "github.com/spf13/viper"

    vaultremote "github.com/sagikazarmark/viperx/remote"
)

func main() {
	vaultremote.RegisterConfigProvider("vault", &myVaultProvider{})
	
	_ = viper.AddRemoteProvider("vault", "endpoint", "path")
}
```

### Hashicorp Vault Remote config provider

```go
package main

import (
    "github.com/spf13/viper"

    "github.com/sagikazarmark/viperx/remote/vault"
)

func main() {
	_ = viper.AddRemoteProvider("vault", "endpoint", "path")
	
	_ = viper.ReadRemoteConfig()
}
```


## License

The MIT License (MIT). Please see [License File](LICENSE) for more information.
