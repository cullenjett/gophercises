package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a task to your task list",
	Run: func(cmd *cobra.Command, args []string) {
		tasks := strings.Join(args, " ")
		fmt.Printf("Added \"%s\" to your task list\n", tasks)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
