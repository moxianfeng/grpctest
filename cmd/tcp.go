package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/spf13/cobra"
)

var tcpOptions struct {
	listenPort int
}

var tcpCommand = &cobra.Command{
	Use: "tcp",
	RunE: func(cmd *cobra.Command, args []string) error {
		runFlag := true

		if tcpOptions.listenPort <= 0 {
			return fmt.Errorf("invalid port")
		}

		li, err := net.Listen("tcp", fmt.Sprintf(":%d", tcpOptions.listenPort))
		if err != nil {
			return err
		}
		defer li.Close()

		for runFlag {
			conn, err := li.Accept()
			if err != nil {
				log.Printf("accept got error, %v\n", err)
			}

			go func(conn net.Conn) {
				defer conn.Close()
				remoteAddr := conn.RemoteAddr().String()

				log.Printf("new connection from %s\n", remoteAddr)
				buf := make([]byte, 200)
				for {
					n, err := conn.Read(buf)
					if err != nil {
						if !errors.Is(err, io.EOF) {
							log.Printf("conn[%s] read got error, close it, %v\n", remoteAddr, err)
						}
						break
					}
					_n, err := conn.Write(buf[:n])
					if err != nil {
						log.Printf("conn[%s] write got error, close it, %v\n", remoteAddr, err)
						break
					}
					if _n != n {
						log.Printf("conn[%s] write not completion, close it, expect %d, got %d\n", remoteAddr, n, _n)
						break
					}
				}

				log.Printf("conn[%s] be closed\n", remoteAddr)
			}(conn)
		}
		return nil
	},
}

func init() {
	tcpCommand.PersistentFlags().IntVarP(&tcpOptions.listenPort, "port", "p", 9999, "listen port")
	rootCmd.AddCommand(tcpCommand)
}
