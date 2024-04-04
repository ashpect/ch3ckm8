/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ch3ckm8",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

const colorRed = "\033[0;31m"
const colorNone = "\033[0m"
const colorYellow = "\033[1;33m"

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ch3ckm8.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	fmt.Printf("\033[2J")
	fmt.Println(colorYellow, "									    WELCOME TO									", colorNone)
	fmt.Println(colorRed, "						    ______  __    __   ____     ______  __  ___ .___  ___.   ___   ", colorNone)
	fmt.Println(colorRed, "						   /      ||  |  |  | |___ \\   /      ||  |/  / |   \\/   |  / _ \\  ", colorNone)
	fmt.Println(colorRed, "						  |  ,----'|  |__|  |   __) | |  ,----'|  '  /  |  \\  /  | | (_) | ", colorNone)
	fmt.Println(colorRed, "						  |  |     |   __   |  |__ <  |  |     |    <   |  |\\/|  |  > _ <  ", colorNone)
	fmt.Println(colorRed, "						  |  `----.|  |  |  |  ___) | |  `----.|  .  \\  |  |  |  | | (_) | ", colorNone)
	fmt.Println(colorRed, "						   \\______||__|  |__| |____/   \\______||__|\\__\\ |__|  |__|  \\___/  ", colorNone)

	fmt.Println()
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
