module xo

go 1.17

require (
	github.com/aws/aws-sdk-go-v2 v1.16.5
	github.com/aws/aws-sdk-go-v2/config v1.15.11
	github.com/aws/aws-sdk-go-v2/service/s3 v1.26.11
	github.com/aws/aws-sdk-go-v2/service/sns v1.17.7
	github.com/aws/aws-sdk-go-v2/service/sts v1.16.7
	github.com/google/uuid v1.3.0
	github.com/itchyny/gojq v0.12.8
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/lightstep/otel-launcher-go v1.5.2
	github.com/manifoldco/promptui v0.9.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/rs/zerolog v1.27.0
	github.com/spf13/cobra v1.5.0
	github.com/spf13/viper v1.12.0
	github.com/xeipuuv/gojsonschema v1.2.0
	github.com/zclconf/go-cty v1.10.0
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/trace v1.7.0
	gopkg.in/ini.v1 v1.66.6
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.2 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.12.6 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.12.6 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.12 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.6 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.13 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.13.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.11.9 // indirect
	github.com/aws/smithy-go v1.11.3 // indirect
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/chzyer/readline v1.5.0 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.3 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/itchyny/timefmt-go v0.1.3 // indirect
	github.com/lightstep/otel-launcher-go/pipelines v1.5.2 // indirect
	github.com/lufia/plan9stats v0.0.0-20220517141722-cf486979b281 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.2 // indirect
	github.com/power-devops/perfstat v0.0.0-20220216144756-c35f1ee13d7c // indirect
	github.com/sethvargo/go-envconfig v0.7.0 // indirect
	github.com/shirou/gopsutil/v3 v3.22.5 // indirect
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.4.0 // indirect
	github.com/tklauser/go-sysconf v0.3.10 // indirect
	github.com/tklauser/numcpus v0.5.0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	go.opentelemetry.io/contrib/instrumentation/host v0.32.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/runtime v0.32.0 // indirect
	go.opentelemetry.io/contrib/propagators/b3 v1.7.0 // indirect
	go.opentelemetry.io/contrib/propagators/ot v1.7.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.7.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.30.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.30.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.7.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.7.0 // indirect
	go.opentelemetry.io/otel/metric v0.30.0 // indirect
	go.opentelemetry.io/otel/sdk v1.7.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v0.30.0 // indirect
	go.opentelemetry.io/proto/otlp v0.18.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/net v0.0.0-20220622184535-263ec571b305 // indirect
	golang.org/x/sys v0.0.0-20220622161953-175b2fd9d664 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220623142657-077d458a5694 // indirect
	google.golang.org/grpc v1.47.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
