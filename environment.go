package files_sdk

type Environment int64

const (
	Production Environment = iota
	Staging
	Local
)

func NewEnvironment(env string) Environment {
	switch env {
	case "staging":
		return Staging
	case "local":
		return Local
	default:
		return Production
	}
}

func (e Environment) String() string {
	switch e {
	case Staging:
		return "staging"
	case Local:
		return "local"
	default:
		return "production"
	}
}

const (
	ProductionEndpoint = "https://{SUBDOMAIN}.files.com"
	localEndpoint      = "https://{SUBDOMAIN}.filesrails.test"
	stagingEndpoint    = "https://{SUBDOMAIN}.filesstaging.av"
)

func (e Environment) Endpoint() string {
	switch e {
	case Staging:
		return stagingEndpoint
	case Local:
		return localEndpoint
	default:
		return ProductionEndpoint
	}
}
