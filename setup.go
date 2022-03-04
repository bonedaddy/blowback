package blowback

import (
	"fmt"
	"strings"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/v2fly/v2ray-core/v4/common/net"
)

func init() { plugin.Register("blowback", setup) }

func setup(c *caddy.Controller) error {
	blowback, err := parse(c)
	if err != nil {
		return plugin.Error("blowback", err)
	}

	// finally add the plugin to coredns
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		blowback.Next = next
		return blowback
	})

	return nil
}

func parse(c *caddy.Controller) (*Blowback, error) {
	blowbackPlugin := New()
	for c.Next() {
		fmt.Println("a")
		args := c.RemainingArgs()
		switch len(args) {
		case 0:
			return nil, c.ArgErr()
		case 1:
			return nil, c.ArgErr()
		case 2:
			if strings.EqualFold("proxy_server", args[0]) {
				if err := isValidHost(args[1]); err != nil {
					return nil, fmt.Errorf("proxy_server argument is invalid host: %s", err)
				}
				blowbackPlugin.ProxyServerURL = args[1]
			}
		}
	}
	return blowbackPlugin, nil
}

// ensures a valid host is being used.
func isValidHost(host string) error {
	if strings.HasPrefix(host, "http://") {
		host = strings.TrimPrefix(host, "http://")
	} else {
		host = strings.TrimPrefix(host, "https://")
	}
	_, err := net.ParseDestination(host)
	return err
}
