module github.com/gotomicro/egoctl

go 1.16

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/davecgh/go-spew v1.1.1
	github.com/flosch/pongo2 v0.0.0-20200529170236-5abacdfa4915
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.4.1
	github.com/gorilla/websocket v1.4.2
	github.com/gotomicro/ego v0.3.10
	github.com/gotomicro/gotoant v0.0.0-20210105085109-df5f1354ac30
	github.com/pelletier/go-toml v1.8.1
	github.com/smartwalle/pongo2render v1.0.1
	github.com/spf13/cobra v1.1.3
	github.com/syndtr/goleveldb v1.0.0
	github.com/uber/prototool v1.10.0
	go.uber.org/zap v1.16.0
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/uber/prototool v1.10.0 => github.com/gotomicro/prototool v1.10.1-0.20210304081706-a1439f175b8c
