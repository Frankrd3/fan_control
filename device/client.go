package device

import (
	"context"
	"github.com/bougou/go-ipmi"
	"github.com/rs/zerolog/log"
)

type Client struct {
	IPMIClient *ipmi.Client
	Context    context.Context
}

func NewClient(client *ipmi.Client) *Client {
	c := &Client{
		IPMIClient: client,
		Context:    context.Background(),
	}

	return c
}

func (c *Client) GetCPUTemperature(cpuId uint8) (float64, error) {
	res, err := c.IPMIClient.GetSDRBySensorID(c.Context, cpuId)
	if res == nil || err != nil {
		return 0.0, err
	}

	return res.Full.SensorValue, nil
}

func (c *Client) GetDeviceInfo() (string, error) {
	res, err := c.IPMIClient.GetDeviceGUID(c.Context)
	if res == nil || err != nil {
		return "", err
	}

	return FormatGUIDAsString(res.GUID), nil
}

func (c *Client) PrintDeviceInfo() error {
	id, err := c.GetDeviceInfo()
	if err != nil {
		return err
	}

	log.Debug().Str("guid", id).Str("host", c.IPMIClient.Host).Str("username", c.IPMIClient.Username).Int("port", c.IPMIClient.Port).Str("interface", string(c.IPMIClient.Interface)).Msg("Established a connection with the device")

	return nil
}

func (c *Client) SetFanSpeed(data uint8) error {
	log.Info().Uint8("speed", data).Msg("Setting new fan speed")
	_, err := c.IPMIClient.RawCommand(c.Context, ipmi.NetFnOEMSupermicroRequest, 0x30, []byte{0x02, 0xff, data}, "")
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DisableFanControl() error {
	log.Info().Msg("Disabling automatic fan control")
	_, err := c.IPMIClient.RawCommand(c.Context, ipmi.NetFnOEMSupermicroRequest, 0x30, []byte{0x01, 0x01}, "")
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) EnableFanControl() error {
	log.Info().Msg("Enabling manual fan control")
	_, err := c.IPMIClient.RawCommand(c.Context, ipmi.NetFnOEMSupermicroRequest, 0x30, []byte{0x01, 0x00}, "")
	if err != nil {
		return err
	}

	return nil
}
