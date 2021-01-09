package gen

import (
	"egoctl/cmd"
	"egoctl/internal/app/module/gen"
	"github.com/spf13/cobra"
)

var CmdGenerate = &cobra.Command{
	Use:   "gen [command]",
	Short: "Source code generator",
	Long:  ``,
}

var (
	flagConfig string
	flagSql    string
)

func init() {
	codeCmd := &cobra.Command{
		Use:   "code",
		Short: "front-end code or backend-code generator",
		Run: func(cmd *cobra.Command, args []string) {
			gen.DefaultEgoctlPro.Run()
		},
	}
	CmdGenerate.PersistentFlags().StringVarP(&flagConfig, "config", "c", "./egoctl.toml", "")
	CmdGenerate.AddCommand(codeCmd)
	cmd.RootCommand.AddCommand(CmdGenerate)
}
