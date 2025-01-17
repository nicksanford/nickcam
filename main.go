package main

import (
	"context"

	"github.com/nicksanford/nickcam/nickcam"
	goutils "go.viam.com/utils"

	"go.viam.com/rdk/components/camera"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/module"
	"go.viam.com/rdk/resource"
)

func mainWithArgs(ctx context.Context, args []string, logger logging.Logger) (err error) {
	resource.RegisterComponent(
		camera.API,
		nickcam.Model,
		resource.Registration[camera.Camera, *nickcam.Config]{Constructor: nickcam.New})

	module, err := module.NewModuleFromArgs(ctx, logger)
	if err != nil {
		return err
	}
	if err := module.AddModelFromRegistry(ctx, camera.API, nickcam.Model); err != nil {
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
	goutils.ContextualMain(mainWithArgs, module.NewLoggerFromArgs("nickcam"))
}
