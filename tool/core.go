package tool

import (
	"os"

	cfglib "github.com/bluest-eel/common/config"
	"github.com/bluest-eel/state/common"
	"github.com/bluest-eel/state/components"
	"github.com/bluest-eel/state/components/config"
	"github.com/bluest-eel/state/components/db"
	log "github.com/sirupsen/logrus"
)

// Tool ...
type Tool struct {
	components.Base
	components.BaseApp
	components.BaseDB
	components.BaseCLI
}

// NewTool ...
func NewTool() *Tool {
	tool := Tool{}
	tool.AppName = common.AppName
	tool.AppAbbv = common.AppAbbreviation
	tool.ProjectPath = common.CallerPaths().DotPath
	tool.RawArgs = os.Args
	return &tool
}

// BootstrapConfiguration ...
func (t *Tool) BootstrapConfiguration() *config.Config {
	cfglib.Setup(t.AppAbbv, t.ProjectPath, config.ConfigFile)
	return config.NewConfig()
}

// // Close the gRPC connection
// func (t *Tool) Close() {
// 	t.APIConn.Close()
// }

// SetupConfiguration ...
func (t *Tool) SetupConfiguration() *config.Config {
	if t.ConfigFile != "" {
		log.Debug("Updating configuration ...")
		log.Debugf("Using project path '%s' ...", t.ProjectPath)
		log.Debugf("Using config file '%s' ...", t.ConfigFile)
		cfglib.Setup(t.AppAbbv, t.ProjectPath, t.ConfigFile)
		return config.NewConfig()
	}
	log.Debug("No config file passed; using bootstrapped config")
	return t.Config
}

// SetupDBConnection ...
func (t *Tool) SetupDBConnection() {
	log.Debug("Setting up database connection ...")
	db, err := db.Open(t.Config)
	if err != nil {
		log.Fatal(err)
	}
	t.DB = &db
}

// SetupGRPCConnection ...
func (t *Tool) SetupGRPCConnection() {
	// connectionOpts := c.Config.GRPCConnectionString()
	// conn, err := grpc.Dial(connectionOpts, grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatalf("did not connect to gRPC server: %v", err)
	// }
	// c.KVConn = conn
	// c.KVClient = api.NewKVServiceClient(conn)
}
