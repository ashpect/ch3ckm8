/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/ashpect/ch3ckm8/engine"
	"github.com/spf13/cobra"
)

// newgameCmd represents the newgame command
var newgameCmd = &cobra.Command{
	Use:   "newgame",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("newgame called")
		new()
	},
}

func init() {
	rootCmd.AddCommand(newgameCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newgameCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newgameCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func new() {
	fmt.Println("Starting new game")
	engine.Uci(engine.Input())
}
