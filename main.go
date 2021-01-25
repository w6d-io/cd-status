/*
Copyright 2020 WILDCARD

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
Created on 22/01/2021
*/

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	resourcev1alpha1 "github.com/tektoncd/pipeline/pkg/apis/resource/v1alpha1"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/w6d-io/ci-status/internal/util"
	"github.com/w6d-io/ci-status/pkg/router"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")

	// Version microservice version
	Version = ""

	// Revision git commit
	Revision = ""

	// GoVersion ...
	GoVersion = ""

	// Built Date built
	Built = ""

	// OsArch ...
	OsArch = ""
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = tkn.AddToScheme(scheme)
	_ = resourcev1alpha1.AddToScheme(scheme)
}

//
func main() {
	var metricsAddr string
	var enableLeaderElection bool

	setupLog.Info("managed flag")
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Encoder: zapcore.NewConsoleEncoder(util.TextEncoderConfig()),
	}
	util.BindFlags(&opts, flag.CommandLine)
	flag.Parse()
	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	if !seen["config"] {
		fmt.Print("config file is missing\n")
		setupLog.Error(errors.New("flag error"), "config file is missing")
		os.Exit(1)
	}

	opts.Development = os.Getenv("RELEASE") != "prod"
	opts.StacktraceLevel = zapcore.PanicLevel
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	setupLog.Info("starting ci-status", "Version", Version, "Built",
		Built, "Revision", Revision, "Arch", OsArch, "GoVersion", GoVersion)
	if err := router.New(); err != nil {
		setupLog.Info("router init")
	}
}
