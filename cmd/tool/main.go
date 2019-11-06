package main

import (
	"os"

	"github.com/bluest-eel/state/components/logging"
	"github.com/bluest-eel/state/tool"
)

// XXX Note that any of this which ends up being useful will be moved into the
//     bluest-eel/cli repo; the `tool` code here is just an experiment to get
//     familiar with the s2 library.
func main() {
	// Create the tool object and assign components to it
	t := tool.NewTool()
	// Bootstrap configuration and logging with defaults; this is to assist with
	// any debugging, e.g., logging output, etc.
	t.Config = t.BootstrapConfiguration()
	t.Logger = logging.LoadClient(t.Config)
	// Now that configuration has been boostrapped, let's pull in anything new
	// from parsed args, opts, etc.:
	t.SetupCLI()
	// defer t.Close()
	t.Run(os.Args)
}
