/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"

	"git.6740.io/scottshotgg/bonk/pkg/agent"
	"git.6740.io/scottshotgg/bonk/pkg/engine"
	engine_basic "git.6740.io/scottshotgg/bonk/pkg/engine/basic"
	core_v1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	controllerAgentName = "nginx-ip-ban-controller"

	// TODO: this needs to be configurable
	ingressDeploymentName = "ingress-nginx-controller"
)

type (
	// Controller is the controller implementation for Foo resources
	Controller struct {
		// kubeclientset is a standard kubernetes clientset
		client kubernetes.Interface

		namespace  string
		deployment string

		agentAddr string
	}
)

// New returns a new ban hammer controller
func New(ctx context.Context, c kubernetes.Interface, ns, d, a string) *Controller {
	return &Controller{
		client:     c,
		namespace:  ns,
		deployment: d,
		agentAddr:  a,
	}
}

func (c *Controller) Run(ctx context.Context, workers int) error {
	fmt.Println("Starting Nginx IP ban controller")

	var pods, err = c.findNginxPods(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, pod := range pods {
		wg.Add(1)

		go func() {
			defer wg.Done()

			err = c.watchPodLogs(ctx, pod)
			if err != nil {
				// TODO: need a channel here for the errors
				fmt.Println(err)
			}
		}()
	}

	wg.Wait()

	// TODO: read from error channel here
	return nil
}

// var (
// 	// TODO: I don't think this should be our responsibility
// 	// but it might be cool to have some custom resources that maybe the
// 	// router would query to keep them in sync
// 	bannedIPs = map[string]struct{}{}
// )

func (c *Controller) watchPodLogs(ctx context.Context, pod core_v1.Pod) error {
	// TODO: namespace needs to be configurable
	var req = c.client.
		CoreV1().
		Pods(c.namespace).
		GetLogs(pod.Name, &core_v1.PodLogOptions{
			Follow: true,
		})

	stream, err := req.Stream(ctx)
	if err != nil {
		return err
	}

	defer func() {
		var err = stream.Close()
		if err != nil {
			fmt.Println("error closing stream:", err)
		}
	}()

	var (
		// TODO: this needs to be configurable
		paths = []string{
			"cgi-bin",
		}

		userAgents = []string{
			"zgrab",
		}

		// TODO: need to make this configurable
		engine engine.Engine = engine_basic.New(paths, userAgents)
		br                   = bufio.NewReader(stream)
	)

	for {
		line, _, err := br.ReadLine()
		switch {
		case err == nil:
			break

		default:
			return err
		}

		ip, shouldBan, err := engine.Run(line)
		if err != nil {
			fmt.Println("err running rules engine:", err)
			continue
		}

		// fmt.Println("Banning IP:", ip)
		// fmt.Println(string(line))

		if shouldBan {
			fmt.Println("bonk:", ip)
			// // TODO: take this out
			// var _, banned = bannedIPs[l.Remote.Address.String()]
			// if banned {
			// 	// fmt.Println("Already banned:", l.Remote.Address)
			// 	continue
			// }

			// fmt.Println("BAN HAMMER TIME:", l.Remote.Address)
			// bannedIPs[l.Remote.Address.String()] = struct{}{}
			// // TODO: yeet this over to the router agent

			err = c.banIP(ip)
			if err != nil {
				return err
			}
		}
	}
}

func (c *Controller) banIP(ip net.IP) error {
	var b, err = json.Marshal(agent.BanIPReq{
		IP: ip.String(),
	})

	if err != nil {
		return err
	}

	res, err := http.Post(fmt.Sprintf("http://%s/ban", c.agentAddr), "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	// Need to close this here since we are in a loop; otherwise it is effectively a memory leak
	// We don't even really need the response though so w/e
	defer res.Body.Close()

	return nil
}

func (c *Controller) findNginxPods(ctx context.Context) ([]core_v1.Pod, error) {
	// TODO: namespace needs to be configurable
	var podList, err = c.client.
		CoreV1().
		Pods(c.namespace).
		List(ctx, v1.ListOptions{})

	if err != nil {
		return nil, err
	}

	/*
		- get all pods with the correct ingress-class
			OR
		- get all pods associate with a deployment
			OR
		- get all pods from an endpoint: `k get endpoints ingress-nginx-controller -o yaml`

		For now we will just get all pods that statically contain the right name in their
		pod name. Kinda janky but w/e
	*/

	var pods []core_v1.Pod
	for _, pod := range podList.Items {
		fmt.Println("pod:", pod.Name)
		if !strings.Contains(pod.Name, c.deployment) {
			continue
		}

		pods = append(pods, pod)
	}

	return pods, nil
}
