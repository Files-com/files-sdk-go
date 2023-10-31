module github.com/Files-com/files-sdk-go/v3

go 1.20

require (
	github.com/appscode/go-querystring v0.0.0-20170504095604-0126cfb3f1dc
	github.com/bradfitz/iter v0.0.0-20191230175014-e8f45d346db8
	github.com/chilts/sid v0.0.0-20190607042430-660e94789ec9
	github.com/dnaeon/go-vcr v1.2.0
	github.com/fatih/structs v1.1.0
	github.com/gin-gonic/gin v1.9.1
	github.com/hashicorp/go-retryablehttp v0.7.4
	github.com/itchyny/timefmt-go v0.1.5
	github.com/lpar/date v1.0.0
	github.com/panjf2000/ants/v2 v2.8.2
	github.com/sabhiram/go-gitignore v0.0.0-20210923224102-525f6e181f06
	github.com/samber/lo v1.38.1
	github.com/snabb/httpreaderat v1.0.1
	github.com/stretchr/testify v1.8.4
	github.com/tunabay/go-infounit v1.1.3
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d
	golang.org/x/text v0.13.0
	moul.io/http2curl/v2 v2.3.0
)

require (
	github.com/bytedance/sonic v1.10.2 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d // indirect
	github.com/chenzhuoyu/iasm v0.9.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.15.5 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.1.0 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	golang.org/x/arch v0.5.0 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sync v0.4.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2023-44487
replace golang.org/x/net => golang.org/x/net v0.17.0
