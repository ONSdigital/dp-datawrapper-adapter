package mapper

import (
	"context"
	"testing"

	"github.com/ONSdigital/dp-datawrapper-adapter/config"

	. "github.com/smartystreets/goconvey/convey"
)

// TODO: remove example test case
func TestUnitMapper(t *testing.T) {
	ctx := context.Background()

	Convey("test mapper adds emphasis to hello world string when set in config", t, func() {
		cfg := config.Config{
			BindAddr:                   "1234",
			GracefulShutdownTimeout:    0,
			HealthCheckInterval:        0,
			HealthCheckCriticalTimeout: 0,
		}

		hm := HelloModel{
			Greeting: "Hello",
			Who:      "World",
		}

		hw := HelloWorld(ctx, hm, cfg)
		So(hw.HelloWho, ShouldEqual, "Hello World")
	})
}
