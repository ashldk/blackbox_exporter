package prober

import (
	"context"
	"testing"
	"time"

	"github.com/miekg/dns"
	"github.com/prometheus/blackbox_exporter/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/promslog"
)

func vmcHandler(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	answers := []string{
		"example.com. 3600 IN TXT \"v=spf1 a:example.com a:example.com include:_spf.google.com include:servers.mcsv.net -all\"",
		"example.com. 3600 IN TXT \"v=BIMI1;l=https://example.com/bimi-icon-tiny-svg-12.svg;a=https://example.com/certificate.pem\"",
	}
	for _, rr := range answers {
		a, err := dns.NewRR(rr)
		if err != nil {
			panic(err)
		}
		m.Answer = append(m.Answer, a)
	}
	if err := w.WriteMsg(m); err != nil {
		panic(err)
	}
}
func TestVMC(t *testing.T) {
	server, addr := startDNSServer("udp", vmcHandler)
	defer server.Shutdown()

	registry := prometheus.NewRegistry()
	testCTX, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result := ProbeVMC(testCTX, "example.com", config.Module{VMC: config.VMCProbe{Resolver: addr}}, registry, promslog.New(&promslog.Config{}))
	if !result {
		t.Error("failed")
	}
}
