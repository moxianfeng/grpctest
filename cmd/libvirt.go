package cmd

import (
	"github.com/spf13/cobra"
	// "libvirt.org/go/libvirt"
)

var libvirtCmd = &cobra.Command{
	Use: "libvirt",
	RunE: func(cmd *cobra.Command, args []string) error {
		// conn, err := libvirt.NewConnect("qemu:///system")
		// if err != nil {
		// 	return err
		// }
		// defer conn.Close()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(libvirtCmd)
}
