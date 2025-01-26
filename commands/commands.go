package commands

import (
	"context"
	"fan_control/configuration"
	"fan_control/device"
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/bougou/go-ipmi"
	"github.com/rs/zerolog/log"
)

type Globals struct {
	Version VersionFlag `name:"version" help:"Print version information and quit"`
	Logging struct {
		Level string `enum:"debug,info,warn,error" default:"info" help:"Set the level the logger should use"`
		Type  string `enum:"json,console" default:"console" help:"Set the type of the logger to use"`
	} `embed:"" prefix:"log."`
}
type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Println(vars["version"])
	app.Exit(0)
	return nil
}

type CLI struct {
	Globals

	Get struct {
		Temp GetCmd `cmd:"" help:"Get CPU temperature information from the device sensors"`
	} `cmd:"" help:"Get device sensor information from IPMI"`

	Run RunCmd `cmd:"" help:"Run the fan control"`
}

func connect(config configuration.Config) (*device.Client, *ipmi.Client, error) {

	// create device client
	deviceClient, err := ipmi.NewClient(config.DeviceConfig.Host, config.DeviceConfig.Port, config.Auth.Username, config.Auth.Password)
	if err != nil {
		return nil, nil, err
	}

	// set the ipmi interface type
	client := device.NewClient(deviceClient)
	client.IPMIClient.Interface = ipmi.InterfaceLanplus

	log.Info().Str("host", config.DeviceConfig.Host).Int("port", config.DeviceConfig.Port).Msg("Fetching CPU temperature from IPMI")

	// try and connect to device
	if err := client.IPMIClient.Connect(context.Background()); err != nil {
		return nil, nil, err
	}

	return client, deviceClient, nil
}
