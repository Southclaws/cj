package admin

func resolveListenAddr(cfg dashboardConfig) (addr, source string) {
	if cfg.ListenAddr != "" {
		return cfg.ListenAddr, "explicit"
	}
	return "127.0.0.1:0", "loopback"
}
