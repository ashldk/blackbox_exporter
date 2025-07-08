package prober

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/prometheus/blackbox_exporter/config"
	"github.com/prometheus/client_golang/prometheus"
)

func ProbeVMC(ctx context.Context, target string, module config.Module, registry *prometheus.Registry, logger *slog.Logger) bool {
	resolver := &net.Resolver{
		Dial: func(c context.Context, network string, address string) (net.Conn, error) {
			return net.Dial(network, module.VMC.Resolver.String())
		},
	}
	txt, err := resolver.LookupTXT(ctx, target)
	if err != nil {
		logger.Error(err.Error())
	}
	fmt.Print(txt)
	// for _, record := range txtRecords {
	// 	if strings.HasPrefix(record, "a=") {
	// 		return strings.TrimPrefix(record, "a=")
	// 	}
	// }

	logger.Error("not impl")
	return false
}
