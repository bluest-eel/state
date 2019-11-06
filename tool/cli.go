package tool

import (
	"encoding/json"
	"fmt"
	"strings"

	cfglib "github.com/bluest-eel/common/config"
	utilib "github.com/bluest-eel/common/util"
	"github.com/bluest-eel/state/common"
	"github.com/bluest-eel/state/components/config"
	"github.com/bluest-eel/state/components/logging"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

/////////////////////////////////////////////////////////////////////////////
///   Constants, etc.   /////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////

const shortDescription = "The Bluest Eel state tool"
const longDescription = `

TBD
`

const template = `
%s{{if len .Authors}}
AUTHORS:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}
WEBSITE: 
   https://github.com/bluest-eel/state

SUPPORT: 
   https://github.com/bluest-eel/state/issues/new

`

var versionData = common.VersionData()

/////////////////////////////////////////////////////////////////////////////
///   CLI Setup   ///////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////

// SetupCLI ...
func (t *Tool) SetupCLI() {
	cli.AppHelpTemplate = fmt.Sprintf(template, cli.AppHelpTemplate)

	t.CLIApp = cli.NewApp()
	t.CLIApp.Authors = []cli.Author{
		cli.Author{
			Name:  "Bluest Eel Team",
			Email: "bluesteel@billo.systems",
		},
	}
	t.CLIApp.Compiled = *versionData.BuildTime
	t.CLIApp.Copyright = strings.Join([]string{
		"(c) 2019 Antonio Macías Ojeda",
		"(c) 2019 BilloSystems, Ltd. Co."}, "\n   ")
	t.CLIApp.Name = "tool"
	t.CLIApp.Usage = shortDescription + longDescription
	t.CLIApp.Version = versionData.Semantic
	t.SetFlags()
	t.SetCommands()
	t.SetActions()
	t.SetBefore()
}

// SetEnvVars looks for specific values in the envionment that are not part of
// configuration and pulls them in. Note that 99% of the time, you'll actually
// want to update the config and not mess with this. ParseEnv is only needed
// for pulling data out of the environment that can impact how configuration
// is read.
func (t *Tool) SetEnvVars() {
	cfgFile := cfglib.EnvConfigFile()
	if cfgFile != "" {
		log.Warn("Overwriting config file flag with ENV var ...")
		t.ConfigFile = cfgFile
	}
}

// SetBefore sets up the code that will be called before the Run client method.
func (t *Tool) SetBefore() {
	t.CLIApp.Before = func(cl *cli.Context) error {
		log.Debugf("Args: %v", cl.Args())
		log.Debugf("Set config file (before env)? '%s'", t.ConfigFile)
		t.SetEnvVars()
		log.Debugf("Set config file (after env)? '%s'", t.ConfigFile)
		log.Debug("Post-bootstrap configuration setup ...")
		// We now may need to make changes to the setup; let's redo it with any of
		// the new inputs that may affect things:
		t.Config = t.SetupConfiguration()
		// If any options have come in which override logging settings (e.g., in a
		// new config file), we'll need to redo this, too:
		log.Debugf("Updated config: %#v", t.Config)
		log.Debugf("Updated logging config: %#v", t.Config.Logging)
		log.Debugf("Updated client logging config: %#v", t.Config.Client.Logging)
		t.Logger = logging.LoadClient(t.Config)
		log.Debug("Post-reconfiguration setup ...")
		t.SetupDBConnection()
		return nil
	}
}

// Run ...
func (t *Tool) Run(args []string) {
	err := t.CLIApp.Run(args)
	if err != nil {
		log.Fatal(err)
	}
}

/////////////////////////////////////////////////////////////////////////////
///   CLI Flags   ///////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////

// SetFlags ...
func (t *Tool) SetFlags() {
	t.CLIApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Value:       config.ConfigFile,
			Usage:       "override the default configuration",
			Destination: &t.ConfigFile,
		},
	}
}

/////////////////////////////////////////////////////////////////////////////
///   CLI Actions   /////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////

// SetActions setup the action to take when the compiled binary will be called
// without any commands/subcommands. This is essentially the same as other
// languages/frameworks that define a default `main` function.
func (t *Tool) SetActions() {
	t.CLIApp.Action = func(cl *cli.Context) error {
		err := cli.ShowAppHelp(cl)
		if err != nil {
			log.Fatal(err)
		}
		log.Error("A command is required for this CLI tool; " +
			"see the help message above.")
		return nil
	}
}

/////////////////////////////////////////////////////////////////////////////
///   CLI Commands   ////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////

// SetCommands ...
func (t *Tool) SetCommands() {
	t.CLIApp.Commands = []cli.Command{
		{
			Name:  "run",
			Usage: "XXX",
			Action: func(cl *cli.Context) error {
				t.RunTool()
				return nil
			},
		},
		{
			Name:  "version",
			Usage: "get all version data as JSON (to format, pipe to `jq`)",
			Action: func(cl *cli.Context) error {
				version := versionToJSON(common.VersionData())
				fmt.Println(version)
				return nil
			},
		},
	}
}

/////////////////////////////////////////////////////////////////////////////
///   Utility Functions for Output   ////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////

func versionToJSON(structData *utilib.Version) string {
	jsonData, err := json.Marshal(structData)
	if err != nil {
		log.Error(err)
		log.Fatalf("Couldn't marshal: %v", structData)
	}
	return string(jsonData)
}
