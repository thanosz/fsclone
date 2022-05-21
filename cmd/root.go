package cmd

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var logger *log.Logger

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fsclone",
	Short: "An rsync wrapper with snapshot support for ZFS, BTRFS and STRATIS",
	Long: `
fsclone - An rsync wrapper with snapshot support for ZFS, BTRFS and STRATIS

Uses rsync to create exact clone of a directory. It was mainly
designed to clone the / partition to another directory in order to be able to 
boot from it (https://wiki.archlinux.org/title/Rsync#As_a_backup_utility). 
Furthermore, it supports taking snapshots in ZFS, BTRFS and STRATIS

rsync options: progress, archive, update, delete, xattrs, hard-links, sparse, 
	        exclude: /dev/* /proc/* /sys/* /tmp/* /run/* /mnt/* /media/* /lost+found


`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Plain Rsync module started")
		_, err := os.Stat("/usr/bin/rsync")
		if err != nil {
			logger.Fatal("/usr/bin/rsync not found")
			os.Exit(1)
		}
		doRsync(cmd.Flags())
		logger.Info("Plain Rsync module finished successfully")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func markCommonFlagsRequired(cmd *cobra.Command) {
	cmd.MarkPersistentFlagRequired("source")
	cmd.MarkPersistentFlagRequired("dest")
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fsclone.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	logname := "/tmp/fsclone.log"
	logfile, _ := os.OpenFile(logname, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	os.Chmod(logname, 0666)

	mw := io.MultiWriter(os.Stdout, logfile)
	logger = log.StandardLogger()
	logger.SetOutput(mw)
	logger.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	//rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose")
	rootCmd.PersistentFlags().String("source", "", "the source directory")
	rootCmd.PersistentFlags().String("dest", "", "the destination directory")
	rootCmd.PersistentFlags().StringSlice("exclude", []string{""}, "extra directories to exclude")
	// rootCmd.MarkPersistentFlagRequired("source")
	// rootCmd.MarkPersistentFlagRequired("dest")

	markCommonFlagsRequired(rootCmd)
}
