package blowback

import (
	"context"
	"strings"

	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/plugin/pkg/nonwriter"

	"github.com/miekg/dns"
)

var log = clog.NewWithPlugin("blowback")

// Dump implement the plugin interface.
type Blowback struct {
	Next plugin.Handler
	// the url of the proxy server which we will initiate
	// a callback to the dns client requestor
	// todo(bonedaddy): require access control
	ProxyServerURL string
}

// ServeDNS implements the plugin.Handler interface.
func (b Blowback) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	log.Info("initializing new nonwriter")
	nw := nonwriter.New(w)
	rcode, err := plugin.NextOrFailure(b.Name(), b.Next, ctx, nw, r)
	if err != nil {
		if strings.Contains(err.Error(), "no next plugin") {
			return 0, nil
		}
		log.Error("plugin.NextOrFailure error", err.Error())
		return rcode, err
	}
	log.Info("plugin.NextOrFailure ok. rcode %v", rcode)
	r = nw.Msg
	log.Info("nw.Msg %+v", r)
	err = w.WriteMsg(r)
	if err != nil {
		log.Error("w.WriteMsg failed", err.Error())
		return 1, err
	}
	return 0, nil
	//state := request.Request{W: w, Req: r}
	//rep := replacer.New()
	//trw := dnstest.NewRecorder(w)
	//fmt.Fprintln(output, rep.Replace(ctx, state, trw, format))
	//return plugin.NextOrFailure(b.Name(), b.Next, ctx, w, r)
}

// Name implements the Handler interface.
func (b Blowback) Name() string { return "blowback" }

func New() *Blowback { return &Blowback{} }
