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
		"default._bimi.vmc.example.com. 3600 IN TXT \"v=spf1 a:example.com a:example.com include:_spf.google.com include:servers.mcsv.net -all\"",
		// in this next line, you won't be able to leave it hard coded; the URL for the cert will only be known to you when the test HTTP server is running in the test harness.
		"default._bimi.vmc.example.com. 3600 IN TXT \"v=BIMI1;l=https://example.com/bimi-icon-tiny-svg-12.svg;a=https://example.com/certificate.pem\"",
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
	// You will need a cert, let Go make it for you:
	// https://pkg.go.dev/crypto/x509@go1.24.5#CreateCertificate

	// You will need to serve the cert, so that ProbeVMC can fetch it.
	// Start a test HTTP server:  https://pkg.go.dev/net/http/httptest@go1.24.5#NewServer
	// that will return the right stuff
	// and you will find out the URL of the test http server
	// (the url is: https://pkg.go.dev/net/http/httptest@go1.24.5#Server   field URL)
	// so you can make vmcHandler return it

	server, addr := startDNSServer("udp", vmcHandler)
	defer server.Shutdown()

	registry := prometheus.NewRegistry()
	testCTX, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res := ProbeVMC(testCTX, "vmc.example.com", config.Module{VMC: config.VMCProbe{ResolverAddr: addr.String()}}, registry, promslog.New(&promslog.Config{}))
	if !res {
		t.Fatal("Unexpected failure")
	}
}

func TestVMCProbeFailure(t *testing.T) {

	server, addr := startDNSServer("udp", recursiveDNSHandler)
	defer server.Shutdown()

	registry := prometheus.NewRegistry()
	testCTX, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res := ProbeVMC(testCTX, "vmc.example.com", config.Module{VMC: config.VMCProbe{ResolverAddr: addr.String()}}, registry, promslog.New(&promslog.Config{}))
	if res {
		t.Fatal("Unexpected success")
	}
}
