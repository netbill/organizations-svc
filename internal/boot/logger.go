package boot

import (
	"github.com/netbill/logium"
	"github.com/sirupsen/logrus"
)

func (c *Config) NewLogger() *logium.Entry {
	log := logium.New()

	lvl, err := logrus.ParseLevel(c.Log.Level)
	if err != nil {
		lvl = logrus.InfoLevel
		log.WithField("bad_level", c.Log.Level).Warn("unknown log level, fallback to info")
	}

	log.SetLevel(lvl)

	switch {
	case c.Log.Format == "json":
		log.SetFormatter(&logrus.JSONFormatter{})
	default:
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	return log.WithField("service", "a;n;;;fn;d")
}
