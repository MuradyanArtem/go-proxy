package servercfg

import (
	"proxy/internal/interfaces/web/server"
	"time"
)

func Export(cfg Config, prefix string) *server.Config {
	var (
		addr = cfg.String(
			prefix+"addr",
			":8000",
			"host:port",
		)
		read = cfg.Duration(
			prefix+"read_timeout",
			time.Duration(3*time.Second),
			"read timeout",
		)
		header = cfg.Duration(
			prefix+"read_header_timeout",
			time.Duration(3*time.Second),
			"read header timeout",
		)
		write = cfg.Duration(
			prefix+"write_timeout",
			time.Duration(3*time.Second),
			"write timeout",
		)
		idle = cfg.Duration(
			prefix+"idle_timeout",
			time.Duration(3*time.Second),
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
