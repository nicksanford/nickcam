package passthroughcam

import (
	"context"
	"errors"
	"fmt"
	"image"

	_ "embed"

	"go.viam.com/rdk/components/camera"
	"go.viam.com/rdk/components/camera/rtppassthrough"
	"go.viam.com/rdk/gostream"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/pointcloud"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/rimage/transform"
	"go.viam.com/rdk/spatialmath"
)

var Model = resource.NewModel("ncs", "camera", "nickcam-passthrough")

type Config struct {
	Camera string `json:"camera,omitempty"`
}

func (c *Config) Validate(path string) ([]string, []string, error) {
	if c.Camera == "" {
		return nil, nil, errors.New("no camera provided")
	}

	return []string{c.Camera}, nil, nil
}

type fake struct {
	resource.Named
	resource.AlwaysRebuild
	resource.TriviallyCloseable
	cam                  camera.Camera
	rtpPassthroughSource rtppassthrough.Source
	logger               logging.Logger
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

	named := conf.ResourceName().AsNamed()
	cam, err := camera.FromDependencies(deps, c.Camera)
	if err != nil {
		return nil, err
	}

	source, ok := cam.(rtppassthrough.Source)
	if !ok {
		return nil, fmt.Errorf("%s doesn't implement rtppassthrough.Source", named.Name().String())
	}

	return &fake{
		Named:                named,
		logger:               logger,
		cam:                  cam,
		rtpPassthroughSource: source,
	}, nil
}

func (f *fake) SubscribeRTP(ctx context.Context, bufferSize int, packetsCB rtppassthrough.PacketCallback) (rtppassthrough.Subscription, error) {
	f.logger.Info("SubscribeRTP START")
	defer f.logger.Info("SubscribeRTP STOP")
	return f.rtpPassthroughSource.SubscribeRTP(ctx, bufferSize, packetsCB)
}

func (f *fake) Unsubscribe(ctx context.Context, id rtppassthrough.SubscriptionID) error {
	f.logger.Info("Unsubscribe START")
	defer f.logger.Info("Unsubscribe STOP")
	return f.rtpPassthroughSource.Unsubscribe(ctx, id)
}

func (f *fake) Image(ctx context.Context, mimeType string, extra map[string]interface{}) ([]byte, camera.ImageMetadata, error) {
	f.logger.Debug("GetImage (NEXT) START")
	defer f.logger.Debug("GetImage (NEXT) END")
	return f.cam.Image(ctx, mimeType, extra)
}

func (f *fake) Images(ctx context.Context, filterSourceNames []string, extra map[string]interface{}) ([]camera.NamedImage, resource.ResponseMetadata, error) {
	return f.cam.Images(ctx, filterSourceNames, extra)
}

func (f *fake) NextPointCloud(ctx context.Context, extra map[string]interface{}) (pointcloud.PointCloud, error) {
	return f.cam.NextPointCloud(ctx, extra)
}

func (f *fake) Geometries(context.Context, map[string]interface{}) ([]spatialmath.Geometry, error) {
	return nil, nil
}

func (f *fake) Projector(ctx context.Context) (transform.Projector, error) {
	return nil, errors.New("unimplemented")
}

func (f *fake) Properties(ctx context.Context) (camera.Properties, error) {
	return f.cam.Properties(ctx)
}

func (f *fake) Stream(ctx context.Context, eh ...gostream.ErrorHandler) (gostream.MediaStream[image.Image], error) {
	f.logger.Debug("Stream START")
	defer f.logger.Debug("Stream END")
	return nil, errors.New("unimplemented")
}

func (f *fake) DoCommand(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
	return f.cam.DoCommand(ctx, extra)
}
