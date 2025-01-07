module est-proxy

go 1.23.0

require (
	est_back_go v0.0.1
	est_proxy_go v0.0.1
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.3
	github.com/jackc/pgx/v5 v5.7.1
	github.com/labstack/echo/v4 v4.12.0
	github.com/labstack/gommon v0.4.2
	github.com/rabbitmq/amqp091-go v1.10.0
	go.elastic.co/apm/module/apmechov4/v2 v2.6.2
	go.elastic.co/apm/module/apmhttp/v2 v2.6.2
	go.elastic.co/apm/v2 v2.6.2
)

require (
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/elastic/go-sysinfo v1.7.1 // indirect
	github.com/elastic/go-windows v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/joeshaw/multierror v0.0.0-20140124173710-69b34d4ec901 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/procfs v0.0.0-20190425082905-87a4384529e0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	go.elastic.co/fastjson v1.1.0 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	howett.net/plist v0.0.0-20181124034731-591f970eefbb // indirect
)

replace (
	est_back_go => ../api/build/est-back-go
	est_proxy_go => ../api/build/est-proxy-go
)
