package main

import (
	"fan_control/commands"
	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

const version = "1.0.0"

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	cli := commands.CLI{
		Globals: commands.Globals{
			Version: commands.VersionFlag(version),
		},
	}

	ctx := kong.Parse(&cli,
		kong.Name("fanctl"),
		kong.Description("A command line utility to control fan speed of the Dell R720 via IPMI either locally or remotely."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{Compact: true}),
		kong.Vars{
			"version": version,
		})

	switch cli.Globals.Logging.Level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		break
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		break
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
		break
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		break
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// log to console by default, otherwise log using json (zerolog default)
	if cli.Globals.Logging.Type == "console" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// set logging level based on cmd args
	switch cli.Globals.Logging.Level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msgf("Global logging level set to: %s.", strings.ToUpper(zerolog.GlobalLevel().String()))
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	err := ctx.Run(&cli.Globals)
	ctx.FatalIfErrorf(err)

}
