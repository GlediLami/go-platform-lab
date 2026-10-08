// Command manager runs all lab controllers against the cluster in your
// current kubeconfig: go run ./solutions/12-manager/cmd
package main

import (
	"os"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	"github.com/GlediLami/go-gardener-lab/pkg/cloud"
	"github.com/GlediLami/go-gardener-lab/pkg/labscheme"
	manager "github.com/GlediLami/go-gardener-lab/solutions/12-manager"
)

func main() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
	log := ctrl.Log.WithName("setup")

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 labscheme.Scheme,
		Metrics:                metricsserver.Options{BindAddress: "0"}, // off for the lab
		HealthProbeBindAddress: "0",
	})
	if err != nil {
		log.Error(err, "Could not create manager")
		os.Exit(1)
	}

	if err := manager.Setup(mgr, cloud.NewFake()); err != nil {
		log.Error(err, "Could not set up controllers")
		os.Exit(1)
	}

	log.Info("Starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		log.Error(err, "Manager stopped with an error")
		os.Exit(1)
	}
}
