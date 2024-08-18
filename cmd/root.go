package cmd

import (
	"context"
	_ "net/http/pprof"
	"os"
	"os/signal"

	"github.com/madlabx/pkgx/errors"
	"github.com/madlabx/pkgx/log"
	"github.com/madlabx/pkgx/utils"
	"github.com/madlabx/pkgx/viperx"
	"github.com/spf13/cobra"

	"github.com/madlabx/fs/api"
	"github.com/madlabx/fs/common/cfg"
	"github.com/madlabx/fs/module/transfer"
	"github.com/madlabx/fs/pkg/buildcontext"
)

func checkFatalError(err error, megs ...string) {
	if err != nil {
		log.Fatalf("Panic: %+v, Message:%v", errors.WithStack(err), megs)
	}
}

var rootCmd = &cobra.Command{
	Use: "fs",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		checkFatalError(initConfigAndLog(ctx))

		printBanner()
		checkFatalError(transfer.Launch(ctx))

		agw, err := api.New(ctx, &cfg.Get().AccessLog, nil)
		checkFatalError(err)

		//block in agw.Run
		checkFatalError(agw.Run(cfg.Get().Sys.Address, cfg.Get().Sys.Port))
	},
}

func printBanner() {
	log.Infof("Start Service:%v", buildcontext.Get())
	log.Infof("Using config file:%v", viperx.ConfigFileUsed())
	log.Infof("Config:%v", utils.ToString(cfg.Get()))
}

// Execute executes the commands.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func Init() {
	ctx, cancel := context.WithCancel(context.Background())
	rootCmd.SetContext(ctx)

	if _, err := viperx.BindAllFlags(rootCmd.Flags(), cfg.Config{}); err != nil {
		log.Panicf("Failed to BindAllFlags, err:%v", err)
	}

	//add sub commands
	rootCmd.AddCommand(configCmd)
	sigCh := make(chan os.Signal, 1)
	go cleanupHandler(sigCh, cancel)
	signal.Notify(sigCh)
}

func initConfigAndLog(ctx context.Context) error {
	configFile := viperx.GetString("sys.ConfigFile", "./conf/file_server.toml")
	err := cfg.Parse("FS", configFile)
	if err != nil {
		return err
	}

	log.SetLoggerOutput(log.StandardLogger(), ctx, cfg.Get().MainLog.LogFile)
	if err = log.SetLevelStr(cfg.Get().MainLog.Level); err != nil {
		return err
	}

	log.SetFormatter(&log.TextFormatter{
		QuoteEmptyFields: true,
		DisableSorting:   true})

	return nil
}
