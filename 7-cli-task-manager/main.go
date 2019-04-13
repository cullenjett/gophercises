package main

import (
	"fmt"
	"os"
	"task/cmd"
	"task/db"
)

func main() {
	dbPath := "tasks.db"
	must(db.Init(dbPath))
	must(cmd.RootCmd.Execute())
}

func must(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
