package prober

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	// "os"
	"strings"
	"time"

	"github.com/prometheus/blackbox_exporter/config"
	"github.com/prometheus/client_golang/prometheus"
)

func extractVMCURL(txt string) string {
	if !strings.HasPrefix(txt, "v=BIMI1") {
		return ""
	}

	parts := strings.Split(txt, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "l=") {
			return strings.TrimPrefix(part, "l=")
		}
	}
	return ""
}

func ProbeVMC(ctx context.Context, target string, module config.Module, registry *prometheus.Registry, logger *slog.Logger) bool {
	bimiDomain := "default._bimi." + target
	resolver := net.DefaultResolver
	if module.VMC.ResolverAddr != "" {
		resolver = &net.Resolver{
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: 5 * time.Second,
				}
				return d.DialContext(ctx, network, module.VMC.ResolverAddr)
			},
		}
	}

	txtRecords, err := resolver.LookupTXT(ctx, bimiDomain)
	if err != nil {
		fmt.Printf("Error looking up TXT record for %s: %v\n", bimiDomain, err)
		return false
	}

	found := false
	for _, txt := range txtRecords {
		fmt.Println("Found TXT record:", txt)
		if strings.HasPrefix(txt, "v=BIMI1") {
			parts := strings.Split(txt, ";")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if strings.HasPrefix(part, "l=") {
					url := strings.TrimPrefix(part, "l=")
					fmt.Println("Found URL:", url)
					found = true
					break
				}
			}
		}
	}

	if !found {
		fmt.Println("No valid BIMI record found.")
	}

	return true
}
