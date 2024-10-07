package types

import (
	"net"
)

type (
	/*
		TODO:
		Make types for this and override the JSON marshaller to get the *correct* types
		i.e: Request.Length should be an int instead of a string but if we don't use it then
		it won't matter anyways so w/e
	*/
	NginxLog struct {
		Remote  Remote  `json:"remote"`
		Request Request `json:"request"`
		Base    Base    `json:"base"`
		// Proxy    Proxy    `json:"proxy"`
		// Upstream Upstream `json:"upstream"`
	}

	Remote struct {
		Address net.IP `json:"addr"`
		User    string `json:"user"`
	}

	Request struct {
		ID         string `json:"id"`
		Method     string `json:"method"`
		URI        string `json:"uri"`
		Completion string `json:"completion"`
		Filename   string `json:"filename"`
		Length     string `json:"length"`
	}

	Base struct {
		TimeLocal     string `json:"time_local"`
		Status        string `json:"status"`
		Scheme        string `json:"scheme"`
		HTTPHost      string `json:"http_host"`
		HTTPReferer   string `json:"http_referer"`
		HTTPUserAgent string `json:"http_user_agent"`
		BodyBytesSent string `json:"body_bytes_sent"`
	}

	Proxy struct {
		UpstreamName            string `json:"upstream_name"`
		AlternativeUpstreamName string `json:"alternative_upstream_name"`
	}

	Upstream struct {
		Address        net.IP `json:"addr"`
		ResponseLength string `json:"response_length"`
		Status         string `json:"status"`
	}
)
