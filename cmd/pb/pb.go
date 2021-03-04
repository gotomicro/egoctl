package pb

import (
	"github.com/uber/prototool/pcmd"

	"github.com/gotomicro/egoctl/cmd"
)

func init() {
	cmd.RootCommand.AddCommand(pcmd.GetCommand())
}

