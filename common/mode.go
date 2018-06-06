package common

import (
	"os"
)

const (
	// mode_env is the name environment variable that set the running mode of core.
	// See below constants for list of available modes.
	mode_env = "KYBER_ENV"

	// DEV_MODE is default running mode. It is used for development.
	DEV_MODE = "dev"
	// PRODUCTION_MODE is the running mode in staging environment.
	STAGING_MODE = "staging"
	// PRODUCTION_MODE is the running mode in production environment.
	PRODUCTION_MODE = "production"
	// MAINNET_MODE is the same as production mode, just for backward compatibility.
	MAINNET_MODE = "mainnet"
	// KOVAN_MODE is the running mode for testing kovan network.
	KOVAN_MODE = "kovan"
	// KOVAN_MODE is the running mode for testing ropsten network.
	ROPSTEN_MODE = "ropsten"
	// SIMULATION_MODE is running mode in simulation.
	SIMULATION_MODE = "simulation"
	// SIMULATION_MODE is running mode for analytic development.
	ANALYTIC_DEV_MODE = "analytic_dev"
)

var validModes = map[string]struct{}{
	DEV_MODE:          {},
	STAGING_MODE:      {},
	PRODUCTION_MODE:   {},
	MAINNET_MODE:      {},
	KOVAN_MODE:        {},
	ROPSTEN_MODE:      {},
	SIMULATION_MODE:   {},
	ANALYTIC_DEV_MODE: {},
}

// RunningMode returns the current running mode of application.
func RunningMode() string {
	mode, ok := os.LookupEnv(mode_env)
	if !ok {
		return DEV_MODE
	}
	_, valid := validModes[mode]
	if !valid {
		return DEV_MODE
	}
	return mode
}
