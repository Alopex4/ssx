package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"ssx/internal/lg"
	"ssx/ssx"
	"ssx/ssx/version"
)

var (
	logVerbose   bool
	printVersion bool
	// ssx Instance 核心结构体
	ssxInst *ssx.SSX
)

func NewRoot() *cobra.Command {
	opt := &ssx.CmdOption{}
	root := &cobra.Command{
		Use:   "ssx",
		Short: "🦅 ssx is a retentive ssh client",
		Example: `# If more than one flag of -i, -s ,-t specified, priority is ENTRY_ID > ADDRESS > TAG_NAME
ssx [-i ENTRY_ID] [-s [USER@]HOST[:PORT]] [-k IDENTITY_FILE] [-t TAG_NAME]`,
		SilenceUsage:       true,
		SilenceErrors:      true,
		DisableAutoGenTag:  true,
		DisableSuggestions: true,
		Args:               cobra.ArbitraryArgs, // accept arbitrary args for supporting quick login
		// 程序运行前执行，确认是否需要追加log(verbose)，是否需要打印版本信息 (version)
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			lg.SetVerbose(logVerbose)
			if !printVersion {
				// 使用工厂函数创建ssx, 将其赋值给SSX对象
				s, err := ssx.NewSSX(opt)
				if err != nil {
					return err
				}
				ssxInst = s
			}
			return nil
		},
		// 程序真正执行，返回值值为Error
		RunE: func(cmd *cobra.Command, args []string) error {
			if printVersion {
				fmt.Fprintln(os.Stdout, version.Detail())
				return nil
			}
			if len(args) > 0 {
				// just use first word as search key
				opt.Keyword = args[0]
			}
			// 内部应用启动入口
			return ssxInst.Main(cmd.Context())
		},
	}
	root.Flags().StringVarP(&opt.DBFile, "file", "f", "", "filepath to store auth data")
	root.Flags().Uint64VarP(&opt.EntryID, "id", "i", 0, "entry id")
	root.Flags().StringVarP(&opt.Addr, "server", "s", "", "target server address\nsupport formats: [user@]host[:port]")
	root.Flags().StringVarP(&opt.Tag, "tag", "t", "", "search entry by tag")
	root.Flags().StringVarP(&opt.IdentityFile, "keyfile", "k", "", "identity_file path")

	root.PersistentFlags().BoolVarP(&printVersion, "version", "v", false, "print ssx version")
	root.PersistentFlags().BoolVar(&logVerbose, "verbose", false, "output detail logs")

	root.AddCommand(newListCmd())
	root.AddCommand(newDeleteCmd())
	root.AddCommand(newTagCmd())

	root.CompletionOptions.HiddenDefaultCmd = true
	root.SetHelpCommand(&cobra.Command{Hidden: true})
	return root
}
