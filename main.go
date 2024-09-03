package main

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	goutils "go.viam.com/utils"
	"golang.org/x/exp/maps"

	"github.com/nicksanford/imageclock/clockdrawer"
	"go.viam.com/rdk/components/camera"
	"go.viam.com/rdk/gostream"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/module"
	"go.viam.com/rdk/pointcloud"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/rimage/transform"
)

var imageTypes = map[string]bool{
	"jpeg": true,
	"png":  true,
}

var imageTypeOptions = maps.Keys(imageTypes)

var colors = map[string]color.NRGBA{
	"white": {R: 255, G: 255, B: 255, A: 255},
	"red":   {R: 255, A: 255},
	"green": {G: 255, A: 255},
	"blue":  {B: 255, A: 255},
}

var colorOptions = maps.Keys(colors)

func init() {
	slices.Sort(colorOptions)
	slices.Sort(imageTypeOptions)
}

var Model = resource.NewModel("ncs", "camera", "nickcam")
var (
	Reset = "\033[0m"
	Green = "\033[32m"
	Cyan  = "\033[36m"
)

type fake struct {
	mu sync.Mutex
	resource.Named
	resource.AlwaysRebuild
	resource.TriviallyCloseable
	clockDrawer *clockdrawer.ClockDrawer
	big         bool
	logger      logging.Logger
}

type Config struct {
	Big       bool   `json:"big,omitempty"`
	Color     string `json:"color,omitempty"`
	ImageType string `json:"image_type,omitempty"`
}

func (c *Config) Validate(path string) ([]string, error) {

	if _, ok := colors[c.Color]; !ok {
		return nil, fmt.Errorf("config color %s invalid, valid colors: %s", c.Color, strings.Join(colorOptions, ", "))
	}

	if _, ok := imageTypes[c.ImageType]; !ok {
		return nil, fmt.Errorf("config image_type %s invalid, valid image types: %s", c.ImageType, strings.Join(imageTypeOptions, ", "))
	}

	return nil, nil
}

type s struct {
	clockDrawer *clockdrawer.ClockDrawer
}

func (s *s) Next(ctx context.Context) (image.Image, func(), error) {
	return s.clockDrawer.Image("image time: " + time.Now().Format(time.RFC3339Nano)), nil, nil
}

func (s *s) Close(ctx context.Context) error {
	return nil
}

func (f *fake) newStream() gostream.MediaStream[image.Image] {
	return &s{f.clockDrawer}
}

func newCam(
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
	cd, err := clockdrawer.New(named.Name().String(), colors[c.Color], c.ImageType, c.Big)
	if err != nil {
		return nil, err
	}

	return &fake{
		Named:       named,
		logger:      logger,
		clockDrawer: &cd,
	}, nil
}

func (f *fake) Images(ctx context.Context) ([]camera.NamedImage, resource.ResponseMetadata, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	ts1 := time.Now()
	nowStr1 := ts1.Format(time.RFC3339Nano)
	img1 := f.clockDrawer.Image("images1 time: " + nowStr1)

	ts2 := time.Now()
	nowStr2 := ts1.Format(time.RFC3339Nano)
	img2 := f.clockDrawer.Image("images2 time: " + nowStr1)

	return []camera.NamedImage{
		{Image: img1, SourceName: nowStr1 + f.clockDrawer.Ext()},
		{Image: img2, SourceName: nowStr2 + f.clockDrawer.Ext()},
	}, resource.ResponseMetadata{CapturedAt: ts2}, nil
}

func (f *fake) NextPointCloud(ctx context.Context) (pointcloud.PointCloud, error) {
	return nil, nil
}

func (f *fake) Projector(ctx context.Context) (transform.Projector, error) {
	return nil, errors.New("Projector unimplemented")
}

func (f *fake) Properties(ctx context.Context) (camera.Properties, error) {
	return camera.Properties{SupportsPCD: true}, nil
}

func (f *fake) Stream(ctx context.Context, eh ...gostream.ErrorHandler) (gostream.MediaStream[image.Image], error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.newStream(), nil
}

func (f *fake) DoCommand(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	_, boom := extra["boom"]
	if boom {
		f.logger.Info(Cyan + "Boom" + Reset)
		os.Exit(1)
	}
	return nil, nil
}

func mainWithArgs(ctx context.Context, args []string, logger logging.Logger) (err error) {
	resource.RegisterComponent(
		camera.API,
		Model,
		resource.Registration[camera.Camera, *Config]{Constructor: newCam})

	module, err := module.NewModuleFromArgs(ctx, logger)
	if err != nil {
		return err
	}
	if err := module.AddModelFromRegistry(ctx, camera.API, Model); err != nil {
		return err
	}

	err = module.Start(ctx)
	defer module.Close(ctx)
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func main() {
	goutils.ContextualMain(mainWithArgs, module.NewLoggerFromArgs(Model.String()))
}
