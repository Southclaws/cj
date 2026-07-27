package admin

import (
	"errors"
	"net"
	"strings"
)

func resolveAllowedHosts(cfg dashboardConfig, resolvedAddr string) ([]string, error) {
	if cfg.AllowedHosts != "" {
		parts := strings.Split(cfg.AllowedHosts, ",")
		hosts := make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.ToLower(strings.TrimSpace(p))
			if p != "" {
				hosts = append(hosts, p)
			}
		}
		if len(hosts) == 0 {
			return nil, errors.New("CJ_DASHBOARD_ALLOWED_HOSTS is set but contains no usable entries")
		}
		return hosts, nil
	}

	if cfg.ListenAddr != "" {
		return nil, errors.New("CJ_DASHBOARD_ALLOWED_HOSTS must be set when CJ_DASHBOARD_LISTEN_ADDR is configured explicitly")
	}

	hosts := []string{"localhost", "127.0.0.1"}
	if host, _, err := net.SplitHostPort(resolvedAddr); err == nil && host != "" && host != "127.0.0.1" {
		hosts = append(hosts, host)
	}
	return hosts, nil
}
