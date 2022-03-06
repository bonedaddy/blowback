package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/miekg/dns"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "bloback"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "server.host",
			Usage: "blowback dns server host",
		},
		&cli.StringFlag{
			Name:  "server.port",
			Usage: "blowback dns server port",
		},
	}
	app.Action = func(c *cli.Context) error {
		listener, err := net.Listen("tcp", "127.0.0.1:1081")
		if err != nil {
			return err
		}
		listenerParts := strings.Split(listener.Addr().String(), ":")
		if len(listenerParts) < 2 {
			return errors.New("too few listener parts")
		}
		hostIp := listenerParts[0]
		port := listenerParts[1]
		fmt.Printf("host ip %s port %s\n", hostIp, port)
		hostIp = strings.Replace(hostIp, ".", "-", -1) // -1 indicates replace all
		go func() {
			for {
				conn, err := listener.Accept()
				if err != nil {
					fmt.Println("failed to accept connection ", err)
					continue
				}
				fmt.Println("received connection ", conn)
				func() {
					defer conn.Close()
					parts := strings.Split(conn.LocalAddr().String(), ":")
					if len(parts) < 2 {
						fmt.Println("parts is less than 2")
						return
					}
					port := parts[1]
					fmt.Println("got port ", port)
				}()
			}
		}()
		m := new(dns.Msg)
		m.SetQuestion(dns.Fqdn(hostIp+"-"+port+".gamma.nightstalker.dev"), dns.TypeA)
		client := new(dns.Client)
		in, rtt, err := client.Exchange(m, c.String("server.host")+":"+c.String("server.port"))
		fmt.Printf(
			"in: %+v\nrtt: %+v\nerr: %s\n", in, rtt, err,
		)
		time.Sleep(time.Second * 30)
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
