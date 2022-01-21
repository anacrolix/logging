package log_test

import (
	"net"
	"testing"

	"github.com/anacrolix/dht/v2"
	"github.com/anacrolix/log"
	"github.com/anacrolix/torrent"
)

// Mirrors usage seen for a particularly expensive logging callsite in anacrolix/dht.
func BenchmarkDhtServerReplyLogger(b *testing.B) {
	nl := log.GetLogger(b.Name())
	nl.Propagate = false
	nl.SetHandler(log.Discard)
	// Wrap the NewLogger for old-style Logger use.
	l := log.Logger{log.RootLoggerImpl{log.GetLogger(b.Name())}}
	l = l.FilterLevel(log.Info).WithValues(&torrent.Client{}).WithContextText("some dht prefix").WithDefaultLevel(log.Debug)
	addr := dht.NewAddr(&net.UDPAddr{
		IP:   net.IPv6loopback,
		Port: 42069,
		Zone: "sup",
	})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		log.Fmsg("reply to %q", addr).Log(l)
	}
}

func getNewLogger(b *testing.B) *log.NewLogger {
	nl := log.GetLogger(b.Name())
	nl.Propagate = false
	nl.SetHandler(log.Discard)
	nl.DefaultLevel = log.Debug
	nl.FilterLevel = log.Info
	return nl
}

// Mirrors usage seen for a particularly expensive logging callsite in anacrolix/dht.
func BenchmarkDhtServerReplyNewLogger(b *testing.B) {
	nl := getNewLogger(b)
	addr := dht.NewAddr(&net.UDPAddr{
		IP:   net.IPv6loopback,
		Port: 42069,
		Zone: "sup",
	})
	b.Run("LazyLog", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			nl.LazyLog(nl.DefaultLevel, func() log.Msg {
				return log.Fmsg("reply to %q", addr).AddValues(nl)
			})
		}
	})
	b.Run("LogLevel", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if ll := nl.LogLevel(nl.DefaultLevel); ll.Ok() {
				ll.Value().Logf("reply to %q", addr)
			}
		}
	})
}
