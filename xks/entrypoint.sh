#!/bin/sh
e2e.test -report-dir=/tmp/results -ginkgo.skip="${E2E_SKIP}" -ginkgo.focus="${E2E_FOCUS}"
