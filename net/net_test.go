package net_test

import (
	"errors"
	"testing"

	netutil "github.com/mistifyio/util/net"
	"github.com/stretchr/testify/assert"
)

func TestSplitHostPort(t *testing.T) {
	tests := []struct {
		input string
		host  string
		port  string
		err   error
	}{
		{"localhost", "localhost", "", nil},
		{"localhost:1234", "localhost", "1234", nil},
		{"[localhost]", "localhost", "", nil},
		{"[localhost]:1234", "localhost", "1234", nil},
		{"2001:db8:85a3:8d3:1319:8a2e:370:7348", "2001:db8:85a3:8d3:1319:8a2e:370", "7348", nil},
		{"[2001:db8:85a3:8d3:1319:8a2e:370:7348]", "2001:db8:85a3:8d3:1319:8a2e:370:7348", "", nil},
		{"[2001:db8:85a3:8d3:1319:8a2e:370:7348]:443", "2001:db8:85a3:8d3:1319:8a2e:370:7348", "443", nil},
		{"2001:db8:85a3:8d3:1319:8a2e:370:7348:443", "2001:db8:85a3:8d3:1319:8a2e:370:7348", "443", nil},
		{":1234", "", "1234", nil},
		{"", "", "", nil},
		{":::", "::", "", nil},
		{"foo:1234:bar", "foo:1234", "bar", nil},
		{"[2001:db8:85a3:8d3:1319:8a2e:370:7348", "", "", errors.New("missing ']'")},
		{"[localhost", "", "", errors.New("missing ']'")},
		{"2001:db8:85a3:8d3:1319:8a2e:370:7348]", "", "", errors.New("missing '['")},
		{"localhost]", "", "", errors.New("missing '['")},
		{"[loca[lhost]:1234", "", "", errors.New("too many '['")},
		{"[loca]lhost]:1234", "", "", errors.New("too many ']'")},
		{"[localhost]:1234]", "", "", errors.New("too many ']'")},
		{"a[localhost]:1234", "", "", errors.New("nothing can come before '['")},
		{"[localhost]:1:234", "localhost", "", errors.New("poorly separated or formatted port")},
	}

	for _, tt := range tests {
		host, port, err := netutil.SplitHostPort(tt.input)
		assert.Equal(t, tt.host, host, tt.input)
		assert.Equal(t, tt.port, port, tt.input)
		assert.Equal(t, tt.err, err, tt.input)
	}
}

func TestLookupSRVPort(t *testing.T) {
	port, err := netutil.LookupSRVPort("_xmpp-server._tcp.google.com")
	assert.NoError(t, err)
	assert.EqualValues(t, 5269, port)

	port, err = netutil.LookupSRVPort("_xmpp-server._tcp.asduhaisudbfa.invalid")
	assert.Error(t, err)
	assert.Empty(t, port)
}

func TestHostWithPort(t *testing.T) {
	tests := []struct {
		input       string
		hostport    string
		expectedErr bool
	}{
		{":8080", ":8080", false},
		{"localhost:8080", "localhost:8080", false},
		{"_xmpp-server._tcp.google.com", "_xmpp-server._tcp.google.com:5269", false},
		{"localhost", "", true},
		{"[localhost", "", true},
	}

	for _, tt := range tests {
		hostport, err := netutil.HostWithPort(tt.input)
		if tt.expectedErr {
			assert.Error(t, err, tt.input)
		}
		assert.Equal(t, tt.hostport, hostport, tt.input)
	}
}
