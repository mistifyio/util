package net_test

import (
	"fmt"
	"os"
	"testing"
	"text/tabwriter"

	netutil "github.com/mistifyio/util/net"
	"github.com/stretchr/testify/assert"
)

func ExampleSplit() {
	examples := []string{
		"localhost",
		"localhost:1234",
		"[localhost]",
		"[localhost]:1234",
		"2001:db8:85a3:8d3:1319:8a2e:370:7348",
		"[2001:db8:85a3:8d3:1319:8a2e:370:7348]",
		"[2001:db8:85a3:8d3:1319:8a2e:370:7348]:443",
		"2001:db8:85a3:8d3:1319:8a2e:370:7348:443",
		":1234",
		"",
		":::",
		"foo:1234:bar",
		"[2001:db8:85a3:8d3:1319:8a2e:370:7348",
		"[localhost",
		"2001:db8:85a3:8d3:1319:8a2e:370:7348]",
		"localhost]",
		"[loca[lhost]:1234",
		"[loca]lhost]:1234",
		"[localhost]:1234]",
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "HOSTPORT\tHOST\tPORT\tERR")
	fmt.Fprintln(w, "========\t====\t====\t===")

	for _, hp := range examples {
		host, port, err := netutil.SplitHostPort(hp)

		fmt.Fprintf(w, "%s\t%s\t%s\t%v\n", hp, host, port, err)
	}
	w.Flush()

	// Output:
	// HOSTPORT					HOST					PORT	ERR
	// ========					====					====	===
	// localhost					localhost					<nil>
	// localhost:1234					localhost				1234	<nil>
	// [localhost]					localhost					<nil>
	// [localhost]:1234				localhost				1234	<nil>
	// 2001:db8:85a3:8d3:1319:8a2e:370:7348		2001:db8:85a3:8d3:1319:8a2e:370		7348	<nil>
	// [2001:db8:85a3:8d3:1319:8a2e:370:7348]		2001:db8:85a3:8d3:1319:8a2e:370:7348		<nil>
	// [2001:db8:85a3:8d3:1319:8a2e:370:7348]:443	2001:db8:85a3:8d3:1319:8a2e:370:7348	443	<nil>
	// 2001:db8:85a3:8d3:1319:8a2e:370:7348:443	2001:db8:85a3:8d3:1319:8a2e:370:7348	443	<nil>
	// :1234											1234	<nil>
	// 												<nil>
	// :::						::						<nil>
	// foo:1234:bar					foo:1234				bar	<nil>
	// [2001:db8:85a3:8d3:1319:8a2e:370:7348								missing ']'
	// [localhost											missing ']'
	// 2001:db8:85a3:8d3:1319:8a2e:370:7348]								missing '['
	// localhost]											missing '['
	// [loca[lhost]:1234										too many '['
	// [loca]lhost]:1234										too many ']'
	// [localhost]:1234]										too many ']'
}

func TestLookupSRVPort(t *testing.T) {
	port, err := netutil.LookupSRVPort("_xmpp-server._tcp.google.com")
	assert.NoError(t, err)
	assert.EqualValues(t, 5269, port)

	port, err = netutil.LookupSRVPort("_xmpp-server._tcp.asduhaisudbfa.invalid")
	assert.Error(t, err)
	assert.Empty(t, 0, port)
}

func TestHostWithPort(t *testing.T) {
	hostport, err := netutil.HostWithPort(":8080")
	assert.NoError(t, err)
	assert.EqualValues(t, ":8080", hostport)

	hostport, err = netutil.HostWithPort("localhost:8080")
	assert.NoError(t, err)
	assert.EqualValues(t, "localhost:8080", hostport)

	hostport, err = netutil.HostWithPort("_xmpp-server._tcp.google.com")
	assert.NoError(t, err)
	assert.EqualValues(t, "_xmpp-server._tcp.google.com:5269", hostport)

	hostport, err = netutil.HostWithPort("localhost")
	assert.Error(t, err)
}
