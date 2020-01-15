package proxy

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"domain-proxy/src/base/config"
)

func TestGetHostAndPortFromKey(t *testing.T) {
	as := assert.New(t)

	config.MockConfig(config.MainConfig{
		Alias: map[string]string{"demo": "192.168.1.10:8888"},
	})

	cases := [][]string{
		{"127-0-0-1-8080", "127.0.0.1", "8080"},
		{"168-3-1-3333", "192.168.3.1", "3333"},
		{"100-1-9080", "192.168.100.1", "9080"},
		{"1-88", "192.168.1.1", "88"},
		{"demo", "192.168.1.10", "8888"},
	}

	for _, c := range cases {
		host, port := getHostAndPortFromKey(c[0])
		as.Equal(c[1], host)
		as.Equal(c[2], port)
	}
}

func BenchmarkGetHostAndPortFromKeyLong(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getHostAndPortFromKey("127-0-0-1-8080")
	}
}

func BenchmarkGetHostAndPortFromKeyShort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getHostAndPortFromKey("1-3333")
	}
}
