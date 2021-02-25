/*
Copyright 2015 The Kubernetes Authors.
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

package e2e

import (
	"flag"
	"math/rand"
	"os"
	"testing"
	"time"

	// Never, ever remove the line with "/ginkgo". Without it,
	// the ginkgo test runner will not detect that this
	// directory contains a Ginkgo test suite.
	// See https://github.com/kubernetes/kubernetes/issues/74827
	// "github.com/onsi/ginkgo"

	"k8s.io/kubernetes/test/e2e/framework"
	"k8s.io/kubernetes/test/e2e/framework/config"

	// test sources
	_ "github.com/xenitab/sonobuoy-plugins/e2e/cluster-health"
	_ "github.com/xenitab/sonobuoy-plugins/e2e/fluxv2"
	_ "github.com/xenitab/sonobuoy-plugins/e2e/high-availability"
	_ "github.com/xenitab/sonobuoy-plugins/e2e/ingress"
	_ "github.com/xenitab/sonobuoy-plugins/e2e/pod-security"
)

func handleFlags() {
	config.CopyFlags(config.Flags, flag.CommandLine)
	framework.RegisterCommonFlags(flag.CommandLine)
	framework.RegisterClusterFlags(flag.CommandLine)
	flag.Parse()
}

func TestMain(m *testing.M) {
	var versionFlag bool
	flag.CommandLine.BoolVar(&versionFlag, "version", false, "Displays version information.")

	// Register test flags, then parse flags.
	handleFlags()

	framework.AfterReadingAllFlags(&framework.TestContext)

	rand.Seed(time.Now().UnixNano())
	os.Exit(m.Run())
}

func TestE2E(t *testing.T) {
	RunE2ETests(t)
}
