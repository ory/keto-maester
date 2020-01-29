/*

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
	"net/http"
	"net/url"
	"os"

	"github.com/ory/keto-maester/keto"

	ketov1alpha1 "github.com/ory/keto-maester/api/v1alpha1"
	"github.com/ory/keto-maester/controllers"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {

	apiv1.AddToScheme(scheme)
	ketov1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var (
		metricsAddr, ketoURL, endpoint string
		ketoPort                       int
		enableLeaderElection           bool
	)

	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&ketoURL, "keto-url", "", "The address of ORY Keto")
	flag.IntVar(&ketoPort, "keto-port", 4444, "Port ORY Keto is listening on")
	flag.StringVar(&endpoint, "endpoint", "/engines", "ORY Keto's engines endpoint")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(zap.Logger(true))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		LeaderElection:     enableLeaderElection,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if ketoURL == "" {
		setupLog.Error(fmt.Errorf("keto URL can't be empty"), "unable to create controller", "controller", "ORYAccessControlPolicy")
		os.Exit(1)
	}

	defaultKeto := ketov1alpha1.Keto{
		URL:      ketoURL,
		Port:     ketoPort,
		Endpoint: endpoint,
	}
	ketoClientMaker := getKetoClientMaker(defaultKeto)
	ketoClient, err := ketoClientMaker(defaultKeto)
	if err != nil {
		setupLog.Error(err, "making default keto client", "controller", "ORYAccessControlPolicy")
		os.Exit(1)

	}

	err = (&controllers.ORYAccessControlPolicyReconciler{
		Client:          mgr.GetClient(),
		Log:             ctrl.Log.WithName("controllers").WithName("ORYAccessControlPolicy"),
		KetoClient:      ketoClient,
		KetoClientMaker: ketoClientMaker,
	}).SetupWithManager(mgr)
	if err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ORYAccessControlPolicy")
		os.Exit(1)
	}
	err = (&controllers.ORYAccessControlPolicyRoleReconciler{
		Client:          mgr.GetClient(),
		Log:             ctrl.Log.WithName("controllers").WithName("ORYAccessControlPolicyRole"),
		KetoClient:      ketoClient,
		KetoClientMaker: ketoClientMaker,
	}).SetupWithManager(mgr)
	if err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ORYAccessControlPolicyRole")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func getKetoClientMaker(defaultSpec controllers.KetoConfiger) controllers.KetoClientMakerFunc {

	defaultKeto := defaultSpec.GetKeto()

	return controllers.KetoClientMakerFunc(func(spec controllers.KetoConfiger) (controllers.KetoClientInterface, error) {

		ketoConf := spec.GetKeto()

		if ketoConf.URL == "" {
			ketoConf.URL = defaultKeto.URL
		}
		if ketoConf.Port == 0 {
			ketoConf.Port = defaultKeto.Port
		}
		if ketoConf.Endpoint == "" {
			ketoConf.Endpoint = defaultKeto.Endpoint
		}

		address := fmt.Sprintf("%s:%d", ketoConf.URL, ketoConf.Port)
		u, err := url.Parse(address)
		if err != nil {
			return nil, fmt.Errorf("unable to parse ORY Keto's URL: %w", err)
		}

		client := &keto.Client{
			KetoURL:    *u.ResolveReference(&url.URL{Path: ketoConf.Endpoint}),
			HTTPClient: &http.Client{},
		}

		return client, nil
	})

}
