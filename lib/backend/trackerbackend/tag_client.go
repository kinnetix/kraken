package trackerbackend

import (
	"errors"
	"fmt"
	"io"

	"code.uber.internal/infra/kraken/lib/backend/backenderrors"
	"code.uber.internal/infra/kraken/lib/serverset"
	"code.uber.internal/infra/kraken/tracker/tagclient"
)

// DockerTagClient is a read-only backend client which resolves tags to manifest
// digest lookups from the tracker.
type DockerTagClient struct {
	client tagclient.Client
}

// NewDockerTagClient creates a new DockerTagClient.
func NewDockerTagClient(config Config) (*DockerTagClient, error) {
	servers, err := serverset.NewRoundRobin(config.RoundRobin)
	if err != nil {
		return nil, fmt.Errorf("round robin: %s", err)
	}
	return &DockerTagClient{tagclient.New(servers)}, nil
}

// Download downloads the manifest digest that the given tag name maps to.
func (c *DockerTagClient) Download(name string, dst io.Writer) error {
	v, err := c.client.Get(name)
	if err != nil {
		if err == tagclient.ErrNotFound {
			return backenderrors.ErrBlobNotFound
		}
		return fmt.Errorf("get tag: %s", err)
	}
	if _, err := io.WriteString(dst, v); err != nil {
		return fmt.Errorf("write to dst: %s", err)
	}
	return nil
}

// Upload is not supported.
func (c *DockerTagClient) Upload(name string, src io.Reader) error {
	return errors.New("upload not supported")
}
