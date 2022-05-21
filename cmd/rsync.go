package cmd

import (
	"os"
	"time"

	"github.com/spf13/pflag"
	"github.com/zloylos/grsync"
	"gitlab.com/variadico/lctime"
)

func doRsync(flags *pflag.FlagSet) bool {

	sourceDir, _ := flags.GetString("source")
	destDir, _ := flags.GetString("dest")
	extraExcludes, _ := flags.GetStringSlice("exclude")
	//verbose, _ := flags.GetBool("verbose")

	if len(sourceDir) == 0 || len(destDir) == 0 {
		logger.Fatal("Empty source or dest")
		os.Exit(1)
	}
	excludes := []string{"/dev/*", "/proc/*", "/sys/*", "/tmp/*", "/run/*", "/mnt/*", "/media/*", "/lost+found"}
	if len(extraExcludes) > 0 {
		excludes = append(excludes, extraExcludes...)
	}
	logger.Info("Cloning ", sourceDir, " to ", destDir)

	task := grsync.NewTask(
		sourceDir,
		destDir,
		grsync.RsyncOptions{
			Exclude:   excludes,
			Progress:  true,
			Archive:   true,
			Update:    true,
			Delete:    true,
			ACLs:      true,
			XAttrs:    true,
			HardLinks: true,
			Sparse:    true,
		},
	)

	// go func() {
	// 	for {
	// 		state := task.State()
	// 		fmt.Printf(
	// 			"progress: %.2f / rem. %d / tot. %d / sp. %s \n",
	// 			state.Progress,
	// 			state.Remain,
	// 			state.Total,
	// 			state.Speed,
	// 		)
	// 		<-time.After(time.Second)
	// 	}
	// }()

	if err := task.Run(); err != nil {
		logger.Warn(task.Log())
		logger.Warn("rsync returned no-zero code: ", err)
	}

	err := os.WriteFile(destDir+"/last_backup", []byte(lctime.Strftime("%a %b %d %T %Z %Y", time.Now())+"\n"), 0644)
	if err != nil {
		logger.Error("Could not create last_backup file on dest ")
	}

	return true
}
