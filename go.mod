module github.com/scorix/grib-go

go 1.23.2

toolchain go1.23.3

require (
	github.com/aliyun/aliyun-oss-go-sdk v3.0.2+incompatible
	github.com/scorix/aliyun-oss-io v0.3.3
	github.com/scorix/go-eccodes v0.1.5
	github.com/scorix/walg v0.4.0
	github.com/stretchr/testify v1.9.0
	go.opentelemetry.io/otel/trace v1.32.0
	golang.org/x/exp v0.0.0-20240904232852-e7e105dedf7e
	golang.org/x/sync v0.9.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/otel v1.32.0 // indirect
	golang.org/x/time v0.7.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// replace github.com/scorix/walg => ../walg
