/*
 * Copyright (C) 2020, MinIO, Inc.
 *
 * This code is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License, version 3,
 * as published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License, version 3,
 * along with this program.  If not, see <http://www.gnu.org/licenses/>
 *
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/minio/minio/cmd/logger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Version provides the version of pvgen
var Version = "DEVELOPMENT.GOGET"

var (
	masterURL    string
	kubeconfig   string
	inputPath    string
	checkVersion bool

	onlyOneSignalHandler = make(chan struct{})
	shutdownSignals      = []os.Signal{os.Interrupt, syscall.SIGTERM}
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "the address of the Kubernetes API server. Overrides any value in kubeconfig.")
	flag.StringVar(&inputPath, "input", "./pv-details.toml", "path to input data for PV.")
	flag.BoolVar(&checkVersion, "version", false, "print version")
}

func main() {
	flag.Parse()

	if checkVersion {
		fmt.Println(Version)
		return
	}

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	logger.FatalIf(err, "Error building kubeconfig")

	kubeClient, err := kubernetes.NewForConfig(cfg)
	logger.FatalIf(err, "Error building Kubernetes client")

	err = createPVs(kubeClient)
	logger.FatalIf(err, "Error creating PV %s")

	logger.Info("Completed PV creation")
}

// setupSignalHandler registered for SIGTERM and SIGINT. A stop channel is returned
// which is closed on one of these signals. If a second signal is caught, the program
// is terminated with exit code 1.
func setupSignalHandler() (stopCh <-chan struct{}) {
	// panics when called twice
	close(onlyOneSignalHandler)

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c
		// second signal. Exit directly.
		os.Exit(1)
	}()

	return stop
}
