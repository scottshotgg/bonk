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

package main

import (
	"flag"
	"fmt"
	"os"

	"git.6740.io/scottshotgg/bonk/pkg/agent"
	"git.6740.io/scottshotgg/bonk/pkg/controller"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/sample-controller/pkg/signals"
)

func main() {
	var (
		masterURL  string
		kubeconfig string
		namespace  string
		deployment string
		agentAddr  string
		isAgent    bool
	)

	// TODO: replace all this config shit with actual JSON config
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&namespace, "namespace", "default", "")
	flag.StringVar(&deployment, "deployment", "ingress-nginx-controller", "")
	flag.StringVar(&agentAddr, "agent-addr", "", "Address of the agent you want to send IPs to; 10.32.0.1:9876")

	flag.BoolVar(&isAgent, "is-agent", false, "")

	flag.Parse()

	// set up signals so we handle the shutdown signal gracefully
	ctx := signals.SetupSignalHandler()

	if isAgent {
		var err = agent.New().Start(ctx)
		if err != nil {
			fmt.Println("Error running controller:", err)
			os.Exit(9)
		}
	}

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		fmt.Println("Error building kubeconfig:", err)
		os.Exit(9)
	}

	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		fmt.Println("Error building kubernetes client:", err)
		os.Exit(9)
	}

	fmt.Println("namespace, deployment:", namespace, deployment)
	controller := controller.New(ctx, client, namespace, deployment, agentAddr)

	err = controller.Run(ctx, 2)
	if err != nil {
		fmt.Println("Error running controller:", err)
		os.Exit(9)
	}
}
