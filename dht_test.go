package logging_test

import (
	"net"
	"testing"

	"github.com/anacrolix/dht/v2"
	"github.com/anacrolix/logging"
)

func getNewLogger(b *testing.B) *logging.Logger {
	nl := logging.GetLogger(b.Name())
	nl.Propagate = false
	nl.SetHandler(logging.Discard)
	nl.DefaultLevel = logging.Debug
	nl.FilterLevel = logging.Info
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
			nl.LazyLog(nl.DefaultLevel, func() logging.Msg {
				return logging.Fmsg("reply to %q", addr).AddValues(nl)
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
