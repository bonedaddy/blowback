package blowback

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/nonwriter"
	"github.com/coredns/coredns/plugin/pkg/replacer"

	"github.com/miekg/dns"
)

// Dump implement the plugin interface.
type Blowback struct {
	Next plugin.Handler
	// the url of the proxy server which we will initiate
	// a callback to the dns client requestor
	// todo(bonedaddy): require access control
	ProxyServerURL string
}

const format = `{remote} ` + replacer.EmptyValue + ` {>id} {type} {class} {name} {proto} {port}`

var output io.Writer = os.Stdout

// ServeDNS implements the plugin.Handler interface.
func (b Blowback) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {

	nw := nonwriter.New(w)
	rcode, err := plugin.NextOrFailure(b.Name(), b.Next, ctx, nw, r)
	if err != nil {
		return rcode, err
	}
	r = nw.Msg
	fmt.Printf("%+v\n", r)
	err = w.WriteMsg(r)
	if err != nil {
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
