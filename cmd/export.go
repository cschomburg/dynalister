package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/xconstruct/dynalister/pkg/dynalist"
)

func init() {
	rootCmd.AddCommand(exportCmd)
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Exports all documents into a directory tree",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		list, err := dynalist.New()
		if err != nil {
			log.Fatalln(err)
		}

		err = list.ExportAll(args[0])
		if err != nil {
			log.Fatalln(err)
		}
	},
}
