package dbcfg

import (
	"proxy/internal/infrastructure/storage"
)

func Export(cfg Config, prefix string) *storage.DBConfig {
	var (
		addr = cfg.String(
			prefix+"host",
			"127.0.0.1",
			"postgres host:port",
		)
		port = cfg.Int(
			prefix+"port",
			5432,
			"postgres user name",
		)
		user = cfg.String(
			prefix+"user",
			"",
			"postgres user name",
		)
		pass = cfg.String(
			prefix+"pass",
			"",
			"postgres user pass",
		)
		name = cfg.String(
			prefix+"name",
			"",
			"postgres DB name",
		)
		isSimple = cfg.Bool(
			prefix+"prefexSimpleProtocol",
			true,
			"postgres set protocol",
		)
		maxConn = cfg.Int(
			prefix+"maxConnections",
			10,
			"pool max connections",
		)
		timeout = cfg.Duration(
			prefix+"acquireTimeout",
			3,
			"postgres set timeout",
		)
	)

	return &storage.DBConfig{
		Host:                 *addr,
		Port:                 uint16(*port),
		User:                 *user,
		Password:             *pass,
		Database:             *name,
		PreferSimpleProtocol: *isSimple,
		MaxConnections:       *maxConn,
		AcquireTimeout:       *timeout,
	}
}
