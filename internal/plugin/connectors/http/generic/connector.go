package generic

import (
	"fmt"
	gohttp "net/http"

	"github.com/cyberark/secretless-broker/pkg/secretless/log"
	"github.com/cyberark/secretless-broker/pkg/secretless/plugin/connector"
)

// Connector injects an HTTP requests with authorization headers.
type Connector struct {
	logger log.Logger
	config *config
}

// Connect implements the http.Connector func signature.
func (c *Connector) Connect(
	r *gohttp.Request,
	credentialsByID connector.CredentialValuesByID,
) error {
	// Validate credential values match expected patterns
	if err := c.config.validate(credentialsByID); err != nil {
		return err
	}

	// Fulfill SSL requests
	if c.config.ForceSSL {
		r.URL.Scheme = "https"
	}

	// Add configured headers to request
	headers, err := c.config.renderedHeaders(credentialsByID)
	if err != nil {
		return fmt.Errorf("failed to render headers: %s", err)
	}
	for headerName, headerVal := range headers {
		r.Header.Set(headerName, headerVal)
	}

	return nil
}
