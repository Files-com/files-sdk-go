package files_sdk

// Environment selects the Files.com API environment.
type Environment int64

const (
	Production Environment = iota
	Staging
	Development
)

// NewEnvironment recognizes "staging" and "development". All other values select production.
func NewEnvironment(env string) Environment {
	switch env {
	case "staging":
		return Staging
	case "development":
		return Development
	default:
		return Production
	}
}

// String returns the environment name. Unknown values return "production".
func (e Environment) String() string {
	switch e {
	case Staging:
		return "staging"
	case Development:
		return "development"
	default:
		return "production"
	}
}

const (
	ProductionEndpoint  = "https://{{SUBDOMAIN}}.files.com"
	developmentEndpoint = "https://{{SUBDOMAIN}}.filesrails.test"
	stagingEndpoint     = "https://{{SUBDOMAIN}}.filesstaging.av"
)

// Endpoint returns the base URL template for the environment. Config.Endpoint
// replaces {{SUBDOMAIN}} with the configured site subdomain.
func (e Environment) Endpoint() string {
	switch e {
	case Staging:
		return stagingEndpoint
	case Development:
		return developmentEndpoint
	default:
		return ProductionEndpoint
	}
}
