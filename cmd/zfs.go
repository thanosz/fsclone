package cmd

import (
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/bitfield/script"
	"github.com/moby/sys/mountinfo"
	"github.com/spf13/cobra"
	"gitlab.com/variadico/lctime"
)

// zfsCmd represents the zfs command
var zfsCmd = &cobra.Command{
	Use:   "zfs",
	Short: "rsync source to dest and take snapshot on ZFS",
	Long: `
fsclone zfs - rsync source to dest and take a snapshot on ZFS

Example:
fsclone zfs --source / --dest /mnt/backup --pool zfspool --dataset root-bak
	
...will sync / to /mnt/backup which is a mounted zfspool/root-bak ZFS pool and
create a snapshot in the form of zfspool/root-bak@Tue-YYYYMMDD. Existing snapshot
for the same day will be firstly deleted.`,

	Run: func(cmd *cobra.Command, args []string) {

		logger.Info("ZFS module started")

		id := os.Getuid()
		if id != 0 {
			logger.Fatal("root is needed to destroy/create zfs snapshots :(")
			os.Exit(1)
		}
		_, err := os.Stat("/usr/sbin/zfs")
		if err != nil {
			logger.Fatal("/usr/sbin/zfs not found")
			os.Exit(1)
		}
		poolName, _ := cmd.Flags().GetString("pool")
		dataset, _ := cmd.Flags().GetString("dataset")
		fsPoolName := poolName + "/" + dataset
		err = exec.Command("/usr/sbin/zfs", "list", fsPoolName).Run()
		if err != nil {
			logger.Fatal("ZFS pool " + poolName + " not found")
			os.Exit(1)
		}
		dest, _ := cmd.Flags().GetString("dest")
		ismount, _ := mountinfo.Mounted(dest)
		if !ismount {
			logger.Error(dest + "is not a mount point")
			os.Exit(1)
		}
		doRsync(cmd.Flags())

		removed := script.Exec("/usr/sbin/zfs list -t snapshot").
			Match(fsPoolName).
			Match("@" + lctime.Strftime("%a", time.Now()) + "-").
			Column(1).ExecForEach("/usr/sbin/zfs destroy -v {{.}}")
		out, _ := removed.String()
		if len(out) > 0 {
			logger.Info("Snapshot removal: " + strings.ReplaceAll(out, "\n", ", "))
		}

		timestamp := lctime.Strftime("%a-%Y%m%d", time.Now())
		snapshotName := fsPoolName + "@" + timestamp
		logger.Info("Creating new snapshot " + snapshotName)
		script.Exec("/usr/sbin/zfs snapshot " + snapshotName).Wait()

		logger.Info("ZFS module finished successfully")
	},
}

func init() {
	rootCmd.AddCommand(zfsCmd)
	zfsCmd.Flags().String("pool", "", "the ZFS pool name")
	zfsCmd.Flags().String("dataset", "", "the ZFS filesytem in the pool")
	zfsCmd.MarkFlagRequired("pool")
	zfsCmd.MarkFlagRequired("dataset")
	markCommonFlagsRequired(zfsCmd)

}
