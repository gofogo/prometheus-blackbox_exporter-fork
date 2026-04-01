// Copyright 2025 The Prometheus Authors
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

package main

import (
	"os"
	"testing"
)

// TestGenerateMatchesGoldenFile verifies that the generated metrics table
// matches the committed metrics.md. If this test fails, a metric was
// added, removed, or its name/type/labels/help changed — run
// `go run ./internal/gen/docs/metrics_v3` to regenerate the file.
func TestGenerateMatchesGoldenFile(t *testing.T) {
	golden, err := os.ReadFile("metrics.md")
	if err != nil {
		t.Fatalf("read golden file: %v", err)
	}

	got := generate()

	if got != string(golden) {
		t.Errorf("generated output differs from metrics.md\n" +
			"Run `go run ./internal/gen/docs/metrics_v3` to regenerate.")
	}
}
