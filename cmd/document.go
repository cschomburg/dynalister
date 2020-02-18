package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/xconstruct/dynalister/pkg/dynalist"
)

func init() {
	rootCmd.AddCommand(documentCmd)
}

var documentCmd = &cobra.Command{
	Use:   "document",
	Short: "Prints a document",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		list, err := dynalist.New()
		if err != nil {
			log.Fatalln(err)
		}

		err = list.PrintDocument(os.Stdout, args[0])
		if err != nil {
			log.Fatal(err)
		}
	},
}
