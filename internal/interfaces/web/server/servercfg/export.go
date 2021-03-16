package servercfg

import "proxy/internal/interfaces/web/server"

func Export(cfg Config, prefix string) *server.Config {
	var (
		addr = cfg.String(
			prefix+"addr",
			"127.0.0.1:8000",
			"host:port",
		)
		read = cfg.Duration(
			prefix+"read_timeout",
			3,
			"read timeout",
		)
		header = cfg.Duration(
			prefix+"read_header_timeout",
			3,
			"read header timeout",
		)
		write = cfg.Duration(
			prefix+"write_timeout",
			3,
			"write timeout",
		)
		idle = cfg.Duration(
			prefix+"idle_timeout",
			3,
			"idle timeout",
		)
	)

	return &server.Config{
		Addr:              *addr,
		ReadTimeout:       *read,
		ReadHeaderTimeout: *header,
		WriteTimeout:      *write,
		IdleTimeout:       *idle,
	}
}
