// Copyright 2021 The Okteto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubeconfig

import (
	"log"
	"os"

	"k8s.io/client-go/rest"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// InCluster returns true if Okteto is running on a Kubernetes cluster
func InCluster() bool {
	_, err := rest.InClusterConfig()
	return err == nil
}

//Create creates a new  kubeconfig file
func Create() *clientcmdapi.Config {
	return clientcmdapi.NewConfig()
}

//Get retrieves a kubeconfig file
func Get(kubeconfigPath string) *clientcmdapi.Config {
	_, err := os.Stat(kubeconfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		log.Fatalf("error accessing your KUBECONFIG file '%s': %v", kubeconfigPath, err)
	}

	cfg, err := clientcmd.LoadFromFile(kubeconfigPath)
	if err != nil {
		log.Fatalf("error accessing your KUBECONFIG file '%s': %v", kubeconfigPath, err)
	}

	return cfg
}

//Write stores a kubeconfig file
func Write(cfg *clientcmdapi.Config, kubeconfigPath string) error {
	return clientcmd.WriteToFile(*cfg, kubeconfigPath)
}

// CurrentContext returns the name of the current context
func CurrentContext(kubeconfigPath string) string {
	cfg := Get(kubeconfigPath)
	if cfg == nil {
		return ""
	}
	return cfg.CurrentContext
}

// CurrentNamespace returns the name of the namespace in use by a given context
func CurrentNamespace(kubeconfigPath string) string {
	cfg := Get(kubeconfigPath)
	if cfg == nil {
		return ""
	}
	if currentContext, ok := cfg.Contexts[cfg.CurrentContext]; ok {
		return currentContext.Namespace
	}
	return "default"
}