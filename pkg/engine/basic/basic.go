package basic

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"git.6740.io/scottshotgg/bonk/pkg/types"
)

type (
	Basic struct {
		paths      map[string]struct{}
		userAgents map[string]struct{}
	}
)

func New(ps, uas []string) *Basic {
	var paths = map[string]struct{}{}
	for _, path := range ps {
		paths[path] = struct{}{}
	}

	var useragents = map[string]struct{}{}
	for _, useragent := range uas {
		useragents[useragent] = struct{}{}
	}

	return &Basic{
		paths:      paths,
		userAgents: useragents,
	}
}

// TODO: probably needs to be configurable
func validLog(line []byte) bool {
	return !(len(line) == 0 || line[0] != '{')
}

// TODO: this needs to be its own pkg
// TODO: tbh MAYBE we should just allow a regex on the json string
func (b *Basic) Run(line []byte) (net.IP, bool, error) {
	// If the line is empty or it is not valid JSON then just ignore it
	if !validLog(line) {
		return nil, false, nil
	}

	var (
		l   types.NginxLog
		err = json.Unmarshal(line, &l)
	)

	if err != nil {
		fmt.Println("err unmarshaling:", err)
		return nil, false, err
	}

	// TODO: don't parse this everytime; save this in the pkg on creation
	_, cidr, err := net.ParseCIDR("10.32.0.0/16")
	if err != nil {
		return nil, false, err
	}

	// If it is from the 10.32.0.0/16 subnet then ignore it
	// TODO: technically there is a security hole here but probably safe
	// Ultimately, it should probably be a flag to provide the CIDR exemption range/s
	// TODO: probably want to omit localhost as well
	if cidr.Contains(l.Remote.Address) {
		return nil, false, nil
	}

	// I have seen weird shit like this
	if isInvalidRequest(l) {
		return l.Remote.Address, true, nil
	}

	// If they aren't exclusively requesting a subdomain of 6740.io (nothing is on
	// 6740.io so they shouldn't be hitting root) then flat out ban 'em
	// TODO: this needs to be configurable to provide the exempt hostnames/IPs and ports
	// TODO: this also contains the port; either fix that in Nginx logs or parse it here
	if !strings.Contains(l.Base.HTTPHost, ".6740.io") {
		return l.Remote.Address, true, nil
	}

	// TODO: Hmmm ... this should probably be a list of regex shit
	// OR just a list of paths maybe
	for path := range b.paths {
		if strings.Contains(l.Request.URI, path) {
			return l.Remote.Address, true, nil
		}
	}

	for userAgent := range b.userAgents {
		if strings.Contains(l.Base.HTTPUserAgent, userAgent) {
			return l.Remote.Address, true, nil
		}
	}

	return nil, false, nil
}

// weird shit here ...
// TODO: probably needs to be configurable
func isInvalidRequest(l types.NginxLog) bool {
	if l.Request.Method == "" {
		return true
	}

	if l.Request.URI == "" {
		return true
	}

	if l.Base.HTTPHost == "" {
		return true
	}

	return false
}
