package cmd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xconstruct/dynalister/pkg/dynalist"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an item to your inbox",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		list, err := dynalist.New()
		if err != nil {
			log.Fatalln(err)
		}

		note := strings.Join(args, " ")
		if err := list.InboxAdd(note); err != nil {
			log.Fatalln(err)
		}
	},
}
