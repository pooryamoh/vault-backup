/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Take a backup and upload using rclone",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		dockerId, _ := cmd.Flags().GetString("dockerid")
		volume, _ := cmd.Flags().GetString("volume")
		dateFormat, _ := cmd.Flags().GetString("dateformat")
		rclonePath, _ := cmd.Flags().GetString("rclonepath")

		_ = backup(dockerId, volume, dateFormat, rclonePath)
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	backupCmd.PersistentFlags().StringP("dockerid", "d", "", "Docker ID of container")
	backupCmd.PersistentFlags().StringP("volume", "v", "", "Path of docker volume")
	backupCmd.PersistentFlags().StringP("dateformat", "f", "", "Date format of filename")
	backupCmd.PersistentFlags().StringP("rclonepath", "p", "", "Path of rclone remote \"remotename:path\"")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// backupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func backup(dockerId string, volume string, dateFormat string, rclonePath string) bool {

	command := fmt.Sprintf("docker stop %s", dockerId)
	_, err := exec.Command("bash", "-c", command).Output()

	if err != nil {
		log.Fatal(err)
	}

	dt := time.Now()
	filename := fmt.Sprint(dt.Format(dateFormat) + ".tar")
	command = fmt.Sprintf("tar -cf %s %s", filename, volume)
	_, err = exec.Command("bash", "-c", command).Output()

	if err != nil {
		log.Fatal(err)
	}

	command = fmt.Sprintf("rclone copy /tmp/%s %s", filename, rclonePath)

	_, err = exec.Command("bash", "-c", command).Output()

	if err != nil {
		log.Fatal(err)
	}

	command = fmt.Sprintf("docker start %s", dockerId)
	_, err = exec.Command("bash", "-c", command).Output()

	if err != nil {
		log.Fatal(err)
	}

	return true
}
