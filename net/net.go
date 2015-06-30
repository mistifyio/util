// Package net provides additional network utility functions.
package net

import (
	"errors"
	"fmt"
	"net"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// hostportCache stores the resulting hostport value for a given input to
// prevent repeated SRV lookups.
var hostportCache = make(map[string]string)

const missingPortMsg = "missing port in address"

// SplitHostPort splits a network address of the form "host", "host:port", "[host]",
// "[host]:port", "[ipv6-host%zone]", or "[ipv6-host%zone]:port" into host or
// ipv6-host%zone and port. Port will be an empty string if not supplied.
func SplitHostPort(hostport string) (host string, port string, err error) {
	var rawport string

	if len(hostport) == 0 {
		return
	}

	// Limit literal brackets to max one open and one closed
	openPos := strings.Index(hostport, "[")
	if openPos != strings.LastIndex(hostport, "[") {
		err = errors.New("too many '['")
		return
	}
	closePos := strings.Index(hostport, "]")
	if closePos != strings.LastIndex(hostport, "]") {
		err = errors.New("too many ']'")
		return
	}

	// Break into host and port parts based on literal brackets
	if openPos > -1 {
		// Needs to open with the '['
		if openPos != 0 {
			err = errors.New("nothing can come before '['")
			return
		}
		// Must have a matching ']'
		if closePos == -1 {
			err = errors.New("missing ']'")
			return
		}
		host = hostport[1:closePos]
		rawport = hostport[closePos+1:]
	} else if closePos > -1 {
		// Did not have a matching '['
		err = errors.New("missing '['")
		return
	} else {
		// No literal brackets, split on the last :
		splitPos := strings.LastIndex(hostport, ":")
		if splitPos < 0 {
			host = hostport
		} else {
			host = hostport[0:splitPos]
			rawport = hostport[splitPos:]
		}
	}

	if rawport != "" {
		if strings.LastIndex(rawport, ":") != 0 {
			err = errors.New("poorly separated or formatted port")
			return
		}
		port = rawport[1:]
	}
	return
}

// LookupSRVPort determines the port for a service via an SRV lookup
func LookupSRVPort(name string) (uint16, error) {
	_, addrs, err := net.LookupSRV("", "", name)
	if err != nil {
		return 0, err
	}

	if len(addrs) == 0 {
		err := errors.New("no srv results")
		return 0, err
	}

	return addrs[0].Port, nil
}

// HostWithPort returns a host:port or [host]:port, performing the necessary
// port lookup if one is not provided. Results are cached.
func HostWithPort(input string) (string, error) {
	hostport, ok := hostportCache[input]
	if ok {
		return hostport, nil
	}

	host, port, err := SplitHostPort(input)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"input": input,
		}).Error("failed to split host and port")
		return "", err
	}

	if port == "" {
		srvPort, err := LookupSRVPort(host)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"host":  host,
			}).Error("srv lookup failed")
			return "", err
		}
		port = fmt.Sprintf("%d", srvPort)
	}

	hostportCache[input] = net.JoinHostPort(host, port)

	return hostportCache[input], nil
}

func isMissingPort(err error) bool {
	addrError, ok := err.(*net.AddrError)
	if !ok {
		return false
	}
	return addrError.Err == missingPortMsg
}
