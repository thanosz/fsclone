# fsclone
An rsync wrapper with snapshot support for ZFS and BTRFS

fsclone uses rsync to create exact clone of a directory. It was mainly
designed to clone the / partition to another directory for backup
(https://wiki.archlinux.org/title/Rsync#As_a_backup_utility). 

The rsync options used are 
    progress, archive, update, delete, xattrs, hard-links, sparse
	  exclude: /dev/* /proc/* /sys/* /tmp/* /run/* /mnt/* /media/* /lost+found

Furthermore, it supports taking snapshots in ZFS, BTRFS. The idea is to keep a
snapshots of the backup filesystem for upto a week. The snapshots are taken after
the cloning has been made. To take snapshots, the standard zfs and btrfs commands
are used.

For example, this is the snapshots of my system
```
NAME                            USED  AVAIL     REFER  MOUNTPOINT
zfspool                         151G   748G       26K  /mnt/zfspool
zfspool/rhel-bak                109G   748G     98.3G  /mnt/zfspool/rhel-bak
zfspool/rhel-bak@Wed-20220518  2.26G      -     96.2G  -
zfspool/rhel-bak@Thu-20220519   958M      -     96.1G  -
zfspool/rhel-bak@Fri-20220520   908M      -     96.2G  -
zfspool/rhel-bak@Sat-20220521   742M      -     98.2G  -
zfspool/rhel-bak@Sun-20220522   768M      -     98.1G  -
zfspool/rhel-bak@Mon-20220523  1.13G      -     98.3G  -
zfspool/rhel-bak@Tue-20220524     0B      -     98.3G  -
```

