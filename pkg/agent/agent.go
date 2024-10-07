package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
)

type (
	Agent struct{}

	BanIPReq struct {
		IP string `json:"ip"`
	}
)

func New() *Agent {
	return &Agent{}
}

func (a *Agent) Start(ctx context.Context) error {
	// TODO: import Gin here or make a swagger better yet and gen some shit
	// TODO: should we have an unban?
	http.HandleFunc("/ban", a.handler)

	// TODO: be extra careful with which interfaces we serve on ...
	// TODO: obligatory ... configurable
	var err = http.ListenAndServe("10.32.0.1:9876", nil)
	if err != nil {
		return err
	}

	return nil
}

// TODO: would be cool to have a pubsub for this
func (a *Agent) handler(rw http.ResponseWriter, r *http.Request) {
	var (
		req BanIPReq
		err = json.NewDecoder(r.Body).Decode(&req)
	)

	if err != nil {
		// TODO: write error or something
		fmt.Println("error decoding body:", err)
	}

	// sudo ipset add myset-ip 107.115.171.57
	var (
		cmd  = "ipset"
		args = []string{
			"add",
			"myset-ip",
			req.IP,
		}
	)

	// TODO: look into setip list instead of this but this is easy
	// iptables -A INPUT -s 154.213.184.15  -j DROP
	output, err := exec.
		CommandContext(r.Context(), cmd, args...).
		CombinedOutput()

	if err != nil {
		// TODO: write error or something
		fmt.Println("error running command:", err)
		fmt.Println("command output:", string(output))
	}
}
