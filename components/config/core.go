package config

import (
	"fmt"

	"github.com/geomyidia/zylog/logger"
	cfg "github.com/spf13/viper"
)

// Metaconfigs ...
const (
	ConfigFile string = "configs/state.yml"
)

/////////////////////////////////////////////////////////////////////////////
///   In-Memory Configuration Mapper   //////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////
///
/// We use structs for converting from Viper's config to something that we
/// can hold in-memory in order to avoid filesystem hits/penalties for
/// hi-performance situations that require faster access to configuration data.
/// Note that Viper's support for refreshing config from disk showed up in
/// previous profiling efforts as a red flag, thus this change/feature.

// BadgerConfig ...
type BadgerConfig struct {
	Directory      string
	ValueDirectory string
}

// CockroachConfig ...
type CockroachConfig struct {
	Host string
	Port int
}

// DBConfig ...
type DBConfig struct {
	Type      string
	Badger    *BadgerConfig
	Cockroach *CockroachConfig
}

// GRPCDConfig ...
type GRPCDConfig struct {
	Host string
	Port int
}

// ClientConfig ...
type ClientConfig struct {
	Logging *logger.ZyLogOptions
}

// Config ...
type Config struct {
	Client  *ClientConfig
	DB      *DBConfig
	GRPCD   *GRPCDConfig
	Logging *logger.ZyLogOptions
}

// NewConfig is a constructor that creates the full coniguration data structure
// for use by our application(s) and client(s) as an in-memory copy of the
// config data (saving from having to make repeated and somewhat expensive
// calls to the viper library).
//
// Note that Viper does provide both the AllSettings() and Unmarshall()
// functions, but these require that you have a struct defined that will be
// used to dump the Viper config data into. We've already got that set up, so
// there's no real benefit to switching.
//
// Furthermore, in our case, we're utilizing structs from other libraries to
// be used when setting those up (see how we initialize the logging component
// in ./components/logging.go, Setup).
func NewConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Type: cfg.GetString("db.type"),
			Badger: &BadgerConfig{
				Directory:      cfg.GetString("db.badger.directory"),
				ValueDirectory: cfg.GetString("db.badger.value-directory"),
			},
			Cockroach: &CockroachConfig{
				Host: cfg.GetString("db.cockroach.host"),
				Port: cfg.GetInt("db.cockroach.port"),
			},
		},
		GRPCD: &GRPCDConfig{
			Host: cfg.GetString("grpc.host"),
			Port: cfg.GetInt("grpc.port"),
		},
		Logging: &logger.ZyLogOptions{
			Colored:      cfg.GetBool("logging.colored"),
			Level:        cfg.GetString("logging.level"),
			Output:       cfg.GetString("logging.output"),
			ReportCaller: cfg.GetBool("logging.report-caller"),
		},
		Client: &ClientConfig{
			Logging: &logger.ZyLogOptions{
				Colored:      cfg.GetBool("client.logging.colored"),
				Level:        cfg.GetString("client.logging.level"),
				Output:       cfg.GetString("client.logging.output"),
				ReportCaller: cfg.GetBool("client.logging.report-caller"),
			},
		},
	}
}

/////////////////////////////////////////////////////////////////////////////
///   Configuration Utility Methods   ///////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////

// GRPCConnectionString ...
func (c *Config) GRPCConnectionString() string {
	return fmt.Sprintf("%s:%d", c.GRPCD.Host, c.GRPCD.Port)
}

// BadgerConnectionString ...
func (c *Config) BadgerConnectionString() string {
	return fmt.Sprintf("%s:%s", c.DB.Badger.Directory, c.DB.Badger.ValueDirectory)
}

// CockroachConnectionString ...
func (c *Config) CockroachConnectionString() string {
	return fmt.Sprintf("%s:%d", c.DB.Cockroach.Host, c.DB.Cockroach.Port)
}
