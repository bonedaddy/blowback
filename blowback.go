package blowback

import (
	"context"
	"errors"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

var log = clog.NewWithPlugin("blowback")
var regIpDash = regexp.MustCompile(`^(\d{1,3}-\d{1,3}-\d{1,3}-\d{1,3})(-\d+)?\.`)

// Dump implement the plugin interface.
type Blowback struct {
	Next plugin.Handler
	// the url of the proxy server which we will initiate
	// a callback to the dns client requestor
	// todo(bonedaddy): require access control
	ProxyServerURL string
}

// taken from https://github.com/wenerme/coredns-ipin
func (b Blowback) Resolve(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (*dns.Msg, int, error) {
	state := request.Request{W: w, Req: r}

	a := new(dns.Msg)
	a.SetReply(r)
	a.Compress = true
	a.Authoritative = true

	matches := regIpDash.FindStringSubmatch(state.QName())
	if len(matches) > 1 {
		ip := matches[1]
		ip = strings.Replace(ip, "-", ".", -1)

		var rr dns.RR
		rr = new(dns.A)
		rr.(*dns.A).Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeA, Class: state.QClass()}
		rr.(*dns.A).A = net.ParseIP(ip).To4()

		a.Answer = []dns.RR{rr}

		if len(matches[2]) > 0 {
			srv := new(dns.SRV)
			srv.Hdr = dns.RR_Header{Name: "_port." + state.QName(), Rrtype: dns.TypeSRV, Class: state.QClass()}
			if state.QName() == "." {
				srv.Hdr.Name = "_port." + state.QName()
			}
			port, _ := strconv.Atoi(matches[2][1:])
			srv.Port = uint16(port)
			srv.Target = "."

			a.Extra = []dns.RR{srv}
		}
	} else {
		// return empty
	}

	return a, 0, nil
}

// ServeDNS implements the plugin.Handler interface.
func (b Blowback) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	a, i, err := b.Resolve(ctx, w, r)
	if err != nil {
		return i, err
	}
	if len(a.Answer) == 0 {
		return 1, plugin.Error("blowback", errors.New("failed to find any answers"))
	}
	if len(a.Extra) == 0 {
		return 1, plugin.Error("blowback", errors.New("failed to find srv record"))
	}
	aRecord, ok := a.Answer[0].(*dns.A)
	if !ok {
		log.Error("failed to parse answer into dns.A record")
		return 1, plugin.Error("blowback", errors.New("failed to parse answer into A record"))
	}
	srvRecord, ok := a.Extra[0].(*dns.SRV)
	if !ok {
		log.Error("failed to parse extra into SRV record")
		return 1, plugin.Error("blowback", errors.New("failed to parse into SRV record"))
	}
	log.Info("spawning background blowback task", "host ", aRecord.A.String(), " port ", srvRecord.Port)
	go func(msg *dns.Msg) {
		log.Info("connecting to proxy server")
	}(a)
	return 0, w.WriteMsg(a)
}

// Name implements the Handler interface.
func (b Blowback) Name() string { return "blowback" }

func New() *Blowback { return &Blowback{} }
