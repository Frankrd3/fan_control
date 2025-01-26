package commands

import (
	"context"
	"fan_control/configuration"
	"fan_control/device"
	"github.com/bougou/go-ipmi"
	"github.com/rs/zerolog/log"
	"time"
)

type RunCmd struct {
	Config   string `short:"c" help:"Path to the config file" default:"./config.yaml"`
	Interval int    `short:"i" help:"How often to check and adjust the fan speed in seconds" default:"10"`
}

func (cmd *RunCmd) Run() error {
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

	// check and adjust if needed fan speed every cmd.Interval seconds
	for {
		temp, err := getAvgTemperature(client, config.DeviceConfig.CPU0ID, config.DeviceConfig.CPU1ID)
		if err != nil {
			return err
		}

		// Connect will create an authenticated session for you.
		if err := adjustFanSpeed(client, temp); err != nil {
			return err
		}

		time.Sleep(time.Duration(cmd.Interval) * time.Second)

	}
}

// Get the average CPU temperature of cpu0 and cpu1
func getAvgTemperature(client *device.Client, cpu0id uint8, cpu1id uint8) (float64, error) {
	cpu0, err := client.GetCPUTemperature(cpu0id)
	cpu1, err := client.GetCPUTemperature(cpu1id)
	if err != nil {
		return 0, nil
	}

	avg := (cpu0 + cpu1) / 2

	return avg, nil
}

func adjustFanSpeed(client *device.Client, avgTemp float64) error {

	switch temp := avgTemp; {
	case temp >= 80:
		log.Debug().Float64("avg_temp", avgTemp).Msg("Got temperature")

		if err := client.DisableFanControl(); err != nil {
			return err
		}

	case temp >= 70:
		log.Debug().Float64("avg_temp", avgTemp).Msg("Got temperature")

		if err := client.SetFanSpeed(0x32); err != nil {
			return err
		}

		// set 50%
	case temp >= 60:
		log.Debug().Float64("avg_temp", avgTemp).Msg("Got temperature")

		if err := client.SetFanSpeed(0x28); err != nil {
			return err
		}

	case temp >= 55:
		log.Debug().Float64("avg_temp", avgTemp).Msg("Got temperature")

		if err := client.SetFanSpeed(0x19); err != nil {
			return err
		}

	case temp >= 50:
		log.Debug().Float64("avg_temp", avgTemp).Msg("Got temperature")

		if err := client.SetFanSpeed(0x14); err != nil {
			return err
		}

	case temp >= 45:
		log.Debug().Float64("avg_temp", avgTemp).Msg("Got temperature")

		if err := client.SetFanSpeed(0x0f); err != nil {
			return err
		}

	case temp >= 40:
		log.Debug().Float64("avg_temp", avgTemp).Msg("Got temperature")

		if err := client.SetFanSpeed(0x0a); err != nil {
			return err
		}

	default:
		log.Debug().Float64("avg_temp", avgTemp).Msg("Got temperature")

		if err := client.SetFanSpeed(0x0a); err != nil {
			return err
		}

	}

	return nil
}
