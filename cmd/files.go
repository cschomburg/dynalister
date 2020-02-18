package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/xconstruct/dynalister/pkg/dynalist"
)

func init() {
	rootCmd.AddCommand(filesCmd)
}

var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "Prints a tree of files in dynalist",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		list, err := dynalist.New()
		if err != nil {
			log.Fatalln(err)
		}

		if err := list.PrintFiles(os.Stdout); err != nil {
			log.Fatalln(err)
		}
	},
}
