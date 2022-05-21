package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// stratisCmd represents the stratis command
var stratisCmd = &cobra.Command{
	Use:   "stratis",
	Short: "todo",
	Long:  `todo`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("stratis not implemeted")
	},
}

func init() {
	rootCmd.AddCommand(stratisCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stratisCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stratisCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
