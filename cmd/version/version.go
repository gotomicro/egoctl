package version

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/gotomicro/egoctl/cmd"
	"github.com/gotomicro/egoctl/internal/config"
	"github.com/gotomicro/egoctl/internal/logger"
	"github.com/gotomicro/egoctl/internal/logger/colors"
)

const verboseVersionBanner string = `%s%s
 ______                          _     _ 
|  ____|                        | |   | |
| |__      __ _    ___     ___  | |_  | |
|  __|    / _  |  / _ \   / __| | __| | |
| |____  | (_| | | (_) | | (__  | |_  | |
|______|  \__, |  \___/   \___|  \__| |_| {{ .Version }}
          __/  |
         |____/
%s
%s%s
├── Version   	    : {{ .Version }}
├── buildGitVersion : {{ .BuildGitVersion }}
├── buildTag        : {{ .BuildTag }}
├── BuildStatus     : {{ .BuildStatus }}
├── BuildUser       : {{ .BuildUser }}
├── BuildHost       : {{ .BuildHost }}
└── BuildTime       : {{ .BuildTime }}
`

const shortVersionBanner = `
 ______                          _     _ 
|  ____|                        | |   | |
| |__      __ _    ___     ___  | |_  | |
|  __|    / _  |  / _ \   / __| | __| | |
| |____  | (_| | | (_) | | (__  | |_  | |
|______|  \__, |  \___/   \___|  \__| |_| {{ .Version }}
          __/  |
         |____/
`

var CmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Prints the current egoctl version",
	Long: `
Prints the current egoctl, ego and Go version alongside the platform information.
`,
	Run: versionCmd,
}
var outputFormat string

const version = config.Version

func init() {
	CmdVersion.PersistentFlags().StringVarP(&outputFormat, "o", "", "", "Set the output format. Either json or yaml.")
	cmd.RootCommand.AddCommand(CmdVersion)
}

func versionCmd(cmd *cobra.Command, args []string) {
	stdout := os.Stdout
	if outputFormat != "" {
		runtimeInfo := map[string]interface{}{
			"Version":         buildVersion,
			"BuildGitVersion": buildGitVersion,
			"BuildTag":        buildTag,
			"BuildStatus":     buildStatus,
			"BuildUser":       buildUser,
			"BuildHost":       buildHost,
			"BuildTime":       buildTime,
		}
		switch outputFormat {
		case "json":
			{
				b, err := json.MarshalIndent(runtimeInfo, "", "    ")
				if err != nil {
					logger.Log.Error(err.Error())
				}
				fmt.Println(string(b))
				return
			}
		case "yaml":
			{
				b, err := yaml.Marshal(&runtimeInfo)
				if err != nil {
					logger.Log.Error(err.Error())
				}
				fmt.Println(string(b))
				return
			}
		}
	}

	coloredBanner := fmt.Sprintf(verboseVersionBanner, "\x1b[35m", "\x1b[1m",
		"\x1b[0m", "\x1b[32m", "\x1b[1m")
	InitBanner(stdout, bytes.NewBufferString(coloredBanner))
	return
}

// ShowShortVersionBanner prints the short version banner.
func ShowShortVersionBanner() {
	output := colors.NewColorWriter(os.Stdout)
	InitBanner(output, bytes.NewBufferString(colors.MagentaBold(shortVersionBanner)))
}
