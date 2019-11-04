package main

import (
	"os"

	"github.com/bluest-eel/state/components/logging"
	"github.com/bluest-eel/state/tool"
)

func main() {
	// Create the tool object and assign components to it
	t := new(tool.Tool)
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
