package cmd

import (
	"os"
	"strings"
	"time"

	"github.com/bitfield/script"
	"github.com/moby/sys/mountinfo"
	"github.com/spf13/cobra"
	"gitlab.com/variadico/lctime"
)

// btrfsCmd represents the btrfs command
var btrfsCmd = &cobra.Command{
	Use:   "btrfs",
	Short: "rsync source to dest and take snapshot on BTRFS",
	Long: `
fsclone btrfs - rsync source to dest and take a snapshot on BTRFS

Example:
	fsclone btrfs --source / --dest /mnt/backup --fslabel FS_BAK
	
BTRFS mode assumes a special way that the btrfs volumes/subvolumes are created. 
vol id 5 contains two folders: @rootfs and .snapshots. The @rootfs is where the 
cloned filesystem resides and mounted under /mnt/backup and the snapshots are 
stored in .snapshots in the form of .snapshots/@Mon-YYYYMMDD. This kind of setup 
was inspired by SuSE's layout of their distros

To avoid inputing the /dev entry, the filesystem is expected to be labeled. In the
example, the --fslabel option will mount the partition with label FS_BAK in a temporary 
folder, execute the sync and then delete/create the snapshot

You can use the --create flag in order to create the schema required.
`,
	Run: func(cmd *cobra.Command, args []string) {

		logger.Info("BTRFS module started")
		id := os.Getuid()
		if id != 0 {
			logger.Fatal("root is needed to destroy/create btrfs snapshots :(")
			os.Exit(1)
		}
		_, err := os.Stat("/usr/bin/btrfs")
		if err != nil {
			logger.Fatal("/usr/bin/btrfs not found")
			os.Exit(1)
		}

		dest, _ := cmd.Flags().GetString("dest")
		ismount, _ := mountinfo.Mounted(dest)
		if !ismount {
			logger.Error("%s is not a mount point\n", dest)
			os.Exit(1)
		}

		doRsync(cmd.Flags())

		tmpMount, err := os.MkdirTemp("", "fsclone-")
		logger.Info("Creating tmp folder " + tmpMount)
		if err != nil && !os.IsExist(err) {
			logger.Fatal(err)
		}
		fslabel, _ := cmd.Flags().GetString("fslabel")
		logger.Info("Mounting subvolid=5 of " + fslabel + " on " + tmpMount)
		script.Exec("/usr/bin/mount -o subvolid=5 LABEL=" + fslabel + " " + tmpMount).Wait()

		removed := script.Exec("btrfs subvolume list " + tmpMount).
			Match("@" + lctime.Strftime("%a", time.Now()) + "-").
			Column(9).ExecForEach("/usr/bin/btrfs subvolume delete " + tmpMount + "/{{.}}")
		out, _ := removed.String()
		if len(out) > 0 {
			logger.Info(strings.ReplaceAll(out, "\n", ", "))
		}

		snapshot := tmpMount + "/.snapshots/@" + lctime.Strftime("%a-%Y%m%d", time.Now())
		logger.Info("Creating new snapshot " + snapshot)
		script.Exec("/usr/bin/btrfs subvolume snapshot " + tmpMount + "/@rootfs " + snapshot).Wait()

		logger.Info("Unmounting and removing " + tmpMount)
		script.Exec("unmount " + tmpMount).Wait()
		os.Remove(tmpMount)

		logger.Info("BTRFS module finished successfully")
	},
}

func init() {
	rootCmd.AddCommand(btrfsCmd)
	btrfsCmd.Flags().String("fslabel", "", "the BTRFS filesystem label")
	btrfsCmd.MarkFlagRequired("fslabel")
	btrfsCmd.Flags().String("create", "", "TODO: create the BTRFS layout")
	markCommonFlagsRequired(rootCmd)

}
