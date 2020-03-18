module github.com/nlnwa/veidemannctl

go 1.13

require (
	github.com/coreos/go-oidc v2.0.0+incompatible
	github.com/dgraph-io/badger/v2 v2.0.0
	github.com/ghodss/yaml v1.0.0
	github.com/golang/protobuf v1.4.0-rc.4
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/magiconair/properties v1.8.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/nlnwa/veidemann-api-go v1.0.0-beta12
	github.com/pkg/errors v0.8.1
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/afero v1.2.1 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.4.0
	golang.org/x/crypto v0.0.0-20190313024323-a1f597ede03a // indirect
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a
	golang.org/x/oauth2 v0.0.0-20190226205417-e64efc72b421
	google.golang.org/grpc v1.28.0
	google.golang.org/protobuf v1.20.1
	gopkg.in/square/go-jose.v2 v2.3.0 // indirect
	gopkg.in/yaml.v2 v2.2.3
)

replace github.com/nlnwa/veidemann-api-go => /home/johnh/prosjekter/veidemann/veidemann-api-go
