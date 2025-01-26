package commands

import (
	"context"
	"fan_control/configuration"
	"github.com/bougou/go-ipmi"
	"github.com/rs/zerolog/log"
)

type GetCmd struct {
	Config string `short:"c" help:"Path to the config file" default:"./config.yaml"`
}

func (cmd *GetCmd) Run() error {
	// get and read config file
	config, err := configuration.ReadConfigFile(cmd.Config)
	if err != nil {
		return err
	}

	// connect to ipmi and create client
	client, deviceClient, err := connect(config)
	if err != nil {
		return err
	}

	defer func(client *ipmi.Client) {
		err := client.Close(context.Background())
		log.Debug().Str("host", client.Host).Int("port", client.Port).Msgf("Disconnected from device")
		if err != nil {
			log.Error().Err(err)
		}
	}(deviceClient)

	if err := client.PrintDeviceInfo(); err != nil {
		return err
	}

	// get cpu temperature for cpu 0 and cpu 1
	cpu0, err := client.GetCPUTemperature(config.DeviceConfig.CPU0ID)
	cpu1, err := client.GetCPUTemperature(config.DeviceConfig.CPU1ID)
	if err != nil {
		return err
	}

	log.Info().Float64("cpu0", cpu0).Float64("cpu1", cpu1).Msg("Received CPU temperature (°C)")

	return nil
}
