package passthroughcam

import (
	"context"
	"errors"

	_ "embed"

	"go.viam.com/rdk/components/camera"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

var Model = resource.NewModel("ncs", "camera", "nickcam-passthrough")

type Config struct {
	Camera string `json:"camera,omitempty"`
}

func (c *Config) Validate(path string) ([]string, error) {
	if c.Camera == "" {
		return nil, errors.New("no camera provided")
	}

	return []string{c.Camera}, nil
}

func New(
	ctx context.Context,
	deps resource.Dependencies,
	conf resource.Config,
	logger logging.Logger,
) (camera.Camera, error) {
	c, err := resource.NativeConfig[*Config](conf)
	if err != nil {
		return nil, err
	}

	return camera.FromDependencies(deps, c.Camera)
}
