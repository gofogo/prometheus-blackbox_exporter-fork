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

// main generates docs/monitoring/metrics.md from prober.AllSpecs.
// No network calls, no dummy probes — all metadata is sourced from the
// MetricSpec vars defined in prober/*_metrics.go.
//
// Run via: go run ./internal/gen/docs/metrics
// Or via:  make generate-metrics-documentation
package main

import (
	"embed"
	"log"
	"sort"
	"strings"

	"github.com/prometheus/blackbox_exporter/config"
	"github.com/prometheus/blackbox_exporter/internal/gen/docs/render"
	"github.com/prometheus/blackbox_exporter/prober"
)

//go:embed templates
var templateFS embed.FS

// metricRow is the flat view used by the template.
type metricRow struct {
	Name   string
	Type   string
	Prober string
	Labels []string
	Help   string
}

// ColWidths holds computed column widths for the markdown table.
type ColWidths struct {
	Name   int
	Type   int
	Prober int
	Labels int
	Help   int
}

// TemplateData is passed to the template engine.
type TemplateData struct {
	Metrics   []metricRow
	ColWidths ColWidths
}

func buildRows() []metricRow {
	// Config / exporter-level metrics (config package can't use MetricSpec
	// due to import cycles, so we handle them inline here).
	rows := []metricRow{
		{
			Name:   config.ConfigReloadSuccessOpts.Namespace + "_" + config.ConfigReloadSuccessOpts.Name,
			Type:   "gauge",
			Prober: "-",
			Help:   config.ConfigReloadSuccessOpts.Help,
		},
		{
			Name:   config.ConfigReloadSuccessTimestampOpts.Namespace + "_" + config.ConfigReloadSuccessTimestampOpts.Name,
			Type:   "gauge",
			Prober: "-",
			Help:   config.ConfigReloadSuccessTimestampOpts.Help,
		},
	}

	// All prober metrics from the single source of truth.
	for _, s := range prober.AllSpecs {
		rows = append(rows, metricRow{
			Name:   s.Opts.Name,
			Type:   s.Type(),
			Prober: s.Prober,
			Labels: s.Labels,
			Help:   s.Opts.Help,
		})
	}
	return rows
}

func sortRows(rows []metricRow) {
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Prober != rows[j].Prober {
			return rows[i].Prober < rows[j].Prober
		}
		return rows[i].Name < rows[j].Name
	})
}

func computeColWidths(rows []metricRow) ColWidths {
	return ColWidths{
		Name:   render.MapColumn("Name", rows, func(r metricRow) string { return r.Name }),
		Type:   render.MapColumn("Metric Type", rows, func(r metricRow) string { return r.Type }),
		Prober: render.MapColumn("Prober", rows, func(r metricRow) string { return r.Prober }),
		Labels: render.MapColumn("Labels", rows, func(r metricRow) string { return strings.Join(r.Labels, ", ") }),
		Help:   render.MapColumn("Help", rows, func(r metricRow) string { return r.Help }),
	}
}

func main() {
	rows := buildRows()
	sortRows(rows)

	data := TemplateData{
		Metrics:   rows,
		ColWidths: computeColWidths(rows),
	}

	content, err := render.RenderTemplate(templateFS, "metrics.gotpl", data)
	if err != nil {
		log.Fatalf("render template: %v", err)
	}

	const outputPath = "docs/monitoring/metrics.md"
	if err := render.WriteToFile(outputPath, content); err != nil {
		log.Fatalf("write %s: %v", outputPath, err)
	}

	log.Printf("wrote %s (%d metrics)", outputPath, len(rows))
}
