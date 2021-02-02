package gen

import (
	"egoctl/cmd"
	"egoctl/internal/app/module/web"
	"github.com/spf13/cobra"
)

var CmdGenerate = &cobra.Command{
	Use:   "web [command]",
	Short: "Web Generator",
	Long:  ``,
}

var (
	flagConfig string
	flagSql    string
)

func init() {
	codeCmd := &cobra.Command{
		Use:   "start",
		Short: "front-end code or backend-code generator",
		Run: func(cmd *cobra.Command, args []string) {
			web.DefaultWebContainer.Run()
		},
	}
	CmdGenerate.PersistentFlags().StringVarP(&flagConfig, "start", "s", "./egoctl.toml", "")
	CmdGenerate.AddCommand(codeCmd)
	cmd.RootCommand.AddCommand(CmdGenerate)
}
