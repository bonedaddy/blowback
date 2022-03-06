package blowback

import (
	"net"
	"time"

	"github.com/miekg/dns"
)

func DnsClientTest() {
	// https://pkg.go.dev/github.com/miekg/dns?utm_source=godoc
	m1 := new(dns.Msg)
	m1.Id = dns.Id()
	m1.RecursionDesired = true
	m1.Question = make([]dns.Question, 1)
	m1.SetEdns0(4096, true)
	m1.Question[0] = dns.Question{"miek.nl.", dns.TypeMX, dns.ClassINET}
	c := new(dns.Client)
	laddr := net.UDPAddr{
		IP:   net.ParseIP("[::1]"),
		Port: 12345,
		Zone: "",
	}
	c.Dialer = &net.Dialer{
		Timeout:   200 * time.Millisecond,
		LocalAddr: &laddr,
	}
	in, rtt, err := c.Exchange(m1, "8.8.8.8:53")
	_ = in
	_ = rtt
	_ = err
}
