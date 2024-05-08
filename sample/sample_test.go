package sample

import (
	"flag"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	featureOptions string
	appName        string
	incremental    bool
)

func TestMain(m *testing.M) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	flag.StringVar(&appName, "appName", "", "Your application name")
	flag.StringVar(&featureOptions, "featureOptions", "", "Comma separated list of feature options")
	flag.BoolVar(&incremental, "incremental", false, "Incremental mode for CI/CD pipeline. If set to true, only check for MRs")
	flag.Parse()

	log.Info().Msgf("appName: %s", appName)
	log.Info().Msgf("featureOptions: %s", featureOptions)
	log.Info().Msgf("incremental: %t", incremental)
	log.Info().Msgf("args: %+v", flag.Args())

	os.Exit(m.Run())
}
