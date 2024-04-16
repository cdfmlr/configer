# config

Lightweight generic configuration package for normal human beings that is ready to go.






## Getting Started

Install the package to your go module:

```bash
go get github.com/cdfmlr/config
```

Define a struct to hold the configuration and load it from a file:

```go
package main

import (
	"fmt"
	"github.com/cdfmlr/config"
	"time"
)

type appConfig struct {
	Version         string
	BackendServices []struct {
		Addr    string
		Timeout time.Duration
		Labels  map[string]string
	}
	DB struct {
		URL  string
		Auth struct{ User, Password string }
	}
}

// AppConfig singleton
var AppConfig appConfig

func main() {
	err := config.NewConfiger(&AppConfig, config.TOML).ReadFromFile("./config.toml")

	if err != nil {
		panic("Failed to read config file: " + err.Error())
	}

	fmt.Printf("AppConfig: Version=%s, len(BackendServices)=%d, DB.Auth.User=%s\n", AppConfig.Version, len(AppConfig.BackendServices), AppConfig.DB.Auth.User)
}
```

Create a config file, e.g. `config.toml`:

```toml
version = "1.0.0"

[[BackendServices]]
Addr = "https://backend1"
Timeout = "5s"
Labels = { env = "dev", region = "us" }

[[BackendServices]]
Addr = "https://backend3"
[BackendServices.Labels]
env = "prod"
region = "eu"
disabled = "true"

[DB]
URL = "mysql://localhost:3306"

[DB.Auth]
User = "root"
Password = "pswd123"
```

Run the program to read the configuration:

```bash
go run main.go
```

Output:

```bash
AppConfig: Version=1.0.0, len(BackendServices)=2, DB.Auth.User=root
```

## Motivation & Philosophy

<details>
  <summary>prologue</summary>

I have relied on [viper](https://github.com/spf13/viper) for years, and I love it.

For complex use cases, viper is the best in the wild.
Especially for docker-like or kubernetes-like projects I have worked on,
I found Viper’s fangs charming to handle the multiple sources 
of configurations with priority rules.
We have few choices to make it a wrap for those Man vs. Wild, 
though viper is also ferocious and hard to tame.

However, for personal toy projects, for demo versions, for baby microservices,
I think viper is too much. I don't need the fangs, I don't need the venom,
I just need to read a config file.

Then I wrote a piece of `config.go` for one of my projects, which is essentially
the same as this package but supports only JSON. I copy-paste it to every
small project I start, it works well, and I am happy with it. A YAML version
is added later, and a TOML version born, too.
I think it’s time to turn it into a standalone package and put an end to the era
of copy-pasting.

So here it is.

</details>

<details>
  <summary>struct is a first-class citizen</summary>

### struct is a first-class citizen

Instead of offering the popular Key-Value configuration,
this package provides a way to load configuration data into a struct,
which is more type-safe and easier to use.

I prefer this:

```go
mqgoUrl := Config.Mqgo.URL
```

to this:

```go
mqgoUrl := Config.Get("Mygo.URL")  // encourage making typos.
```

I prefer this:

```go
type Config struct {
    Mqgo struct {
        URL string
    }
}

mqgoUrl := Config.Mqgo.URL
```

to this:

```go
type configKey int

const (
    configKeyMqgoURL configKey = iota
)

mqgoUrl := Config.GetString(configKeyMqgoURL)
```

</details>

<details>
  <summary>single source of truth</summary>

### single source of truth

Instead of allowing configuration from multiple sources
(multiple files, environment variables, command line flags, etc.),
that costs you a whole weekend to learn the priority rules,
that takes you the following weekdays to debug the unexpected behaviors,
this package only supports loading configuration from a single source
(a file or a `io.Reader`), which is simple and clear.

I prefer this:

```bash
$ cat > config.toml <<EOF
listen = ":8080"
EOF

$ ./myservice -c config.toml

$ curl http://localhost:8080
```

to this:

```bash
$ cat > /etc/myservice/config.toml <<EOF
listen = ":8080"
EOF

$ cat > ~/.myservice/config.yaml <<EOF
listen: ":8081"
EOF

$ cat > ./config.json <<EOF
{"listen": ":8082"}
EOF

$ export MYSERVICE_LISTEN=":8083"

$ ./myservice --config etcd=etcd://localhost:2379 --listen ":8084"

$ curl http://what.the.hell.is.the.port:8086?
```

</details>

<details>
  <summary>keep it simple, stupid</summary>

### keep it simple, stupid

Instead of providing a lot of features that you may never use and
a lot of dependencies that scare you every time you open the `go.mod` file,
which is definitely an overkill for your simple 10,000-line mirco-service project,
this package only depends on the standard library + `gopkg.in/yaml.v3` + `github.com/BurntSushi/toml`.

I prefer this (for small projects):

```go
module github.com/cdfmlr/config

go 1.21.5

require (
	github.com/BurntSushi/toml v1.3.2
	gopkg.in/yaml.v3 v3.0.1
)
```

to this:

```go
module github.com/spf13/viper

go 1.20

require (
	github.com/fsnotify/fsnotify v1.7.0
	github.com/hashicorp/hcl v1.0.0
	github.com/magiconair/properties v1.8.7
	github.com/mitchellh/mapstructure v1.5.0
	github.com/pelletier/go-toml/v2 v2.2.0
	github.com/sagikazarmark/crypt v0.19.0
	github.com/sagikazarmark/locafero v0.4.0
	github.com/sagikazarmark/slog-shim v0.1.0
	github.com/spf13/afero v1.11.0
	github.com/spf13/cast v1.6.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.9.0
	github.com/subosito/gotenv v1.6.0
	gopkg.in/ini.v1 v1.67.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	cloud.google.com/go v0.112.1 // indirect
	cloud.google.com/go/compute v1.24.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/firestore v1.15.0 // indirect
	cloud.google.com/go/longrunning v0.5.5 // indirect
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fatih/color v1.14.1 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.3 // indirect
	github.com/hashicorp/consul/api v1.28.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.5.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.17.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nats-io/nats.go v1.34.0 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	go.etcd.io/etcd/api/v3 v3.5.12 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.12 // indirect
	go.etcd.io/etcd/client/v2 v2.305.12 // indirect
	go.etcd.io/etcd/client/v3 v3.5.12 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.49.0 // indirect
	go.opentelemetry.io/otel v1.24.0 // indirect
	go.opentelemetry.io/otel/metric v1.24.0 // indirect
	go.opentelemetry.io/otel/trace v1.24.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/oauth2 v0.18.0 // indirect
	golang.org/x/sync v0.6.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	google.golang.org/api v0.171.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto v0.0.0-20240213162025-012b6fc9bca9 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240311132316-a219d84964c2 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240314234333-6e1732d8331c // indirect
	google.golang.org/grpc v1.62.1 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)
```

</details>

## Features

- [x] Read configuration from file or `io.Reader`;
- [x] Marshaling and unmarshalling of configuration data;
- [x] Support for JSON, YAML and TOML formats;
- [x] Binding configuration to a struct;
- [x] Write configuration back to file or `io.Writer`;

## TODO

- [ ] Export `encoding`, `decoder` and `encoder` interfaces for custom formats;
- [ ] Watch config file changes (I am not sure if it is necessary, due to a config changing, in many cases, indicates a restart of everything related to. This is far beyond the responsibility of a config loader package. Consider use `fsnotify` at a higher level to handle this properly if necessary.)

## License

MIT License, Copyright (c) 2023-present CDFMLR
