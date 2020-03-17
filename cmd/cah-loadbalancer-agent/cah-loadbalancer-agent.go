package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"net/http"

	"k8s.io/klog"

	"github.com/cloudandheat/cah-loadbalancer/pkg/agent"
	"github.com/cloudandheat/cah-loadbalancer/pkg/config"
)

var (
	configPath string
)

func main() {
	flag.Parse()

	fileCfg, err := config.ReadAgentConfigFromFile(configPath)
	if err != nil {
		klog.Fatalf("Failed reading config: %s", err.Error())
	}

	config.FillAgentConfig(&fileCfg)

	sharedSecret, err := base64.StdEncoding.DecodeString(fileCfg.SharedSecret)
	if err != nil {
		klog.Fatalf("shared-secret failed to decode: %s", err.Error())
	}

	_ = sharedSecret

	http.Handle("/v1/apply", &agent.ApplyHandlerv1{
		MaxRequestSize: 1048576,
		KeepalivedGenerator: &agent.KeepalivedConfigGenerator{
			VRIDBase: fileCfg.Keepalived.VRIDBase,
			VRRPPassword: fileCfg.Keepalived.VRRPPassword,
			Interface: fileCfg.Keepalived.Interface,
			Priority: fileCfg.Keepalived.Priority,
		},
		KeepalivedOutputFile: fileCfg.Keepalived.OutputFile,
	})

	http.ListenAndServe(fmt.Sprintf("%s:%d", fileCfg.BindAddress, fileCfg.BindPort), nil)
}

func init() {
	flag.StringVar(&configPath, "config", "agent-config.toml", "Path to the agent config file.")
}
