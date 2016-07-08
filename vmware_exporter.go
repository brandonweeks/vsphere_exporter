// Copyright 2016 Square Inc.
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
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
	"github.com/serenize/snaker"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

var (
	listenAddress     = flag.String("web.listen-address", ":9103", "Address on which to expose metrics and web interface.")
	metricsPath       = flag.String("web.telemetry-path", "/metrics", "Path under which to expose Prometheus metrics.")
	vsphereHostname   = flag.String("vsphere.hostname", "localhost", "")
	vsphereUsername   = flag.String("vsphere.username", "administrator@vsphere.local", "")
	vspherePassword   = flag.String("vsphere.password", "vmware", "")
	vsphereDatacenter = flag.String("vsphere.datacenter", "Datacenter", "")
)

type Exporter struct {
	ctx                context.Context
	client             govmomi.Client
	performanceManager mo.PerformanceManager
}

func NewExporter() *Exporter {
	ctx, _ := context.WithCancel(context.Background())

	u, err := url.Parse(fmt.Sprintf("https://%s:%s@%s/sdk", *vsphereUsername, *vspherePassword, *vsphereHostname))
	if err != nil {
		log.Fatal(err)
	}

	client, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		log.Fatal(err)
	}

	var performanceManager mo.PerformanceManager
	err = client.RetrieveOne(ctx, *client.ServiceContent.PerfManager, nil, &performanceManager)
	if err != nil {
		log.Fatal(err)
	}

	return &Exporter{
		ctx:                ctx,
		client:             *client,
		performanceManager: performanceManager,
	}
}

var countersInfoMap = make(map[int]*prometheus.Desc)

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, perfCounterInfo := range e.performanceManager.PerfCounter {
		groupInfo := perfCounterInfo.GroupInfo.GetElementDescription()
		nameInfo := perfCounterInfo.NameInfo.GetElementDescription()
		metricName := fmt.Sprintf("vsphere_%s_%s", snaker.CamelToSnake(groupInfo.Key), strings.Join(strings.Split(snaker.CamelToSnake(nameInfo.Key), "."), "_"))
		labels := []string{"host", "instance", "entity"}
		desc := prometheus.NewDesc(metricName, nameInfo.Summary, labels, nil)
		countersInfoMap[int(perfCounterInfo.Key)] = desc
		ch <- desc
	}
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	finder := find.NewFinder(e.client.Client, true)
	datacenter, err := finder.Datacenter(e.ctx, *vsphereDatacenter)
	if err != nil {
		log.Fatal(err)
	}

	finder.SetDatacenter(datacenter)
	hosts, err := finder.HostSystemList(e.ctx, "*")
	if err != nil {
		log.Fatal(err)
	}

	for _, host := range hosts {
		hostName := host.Name()
		querySpec := types.PerfQuerySpec{
			Entity:     host.Reference(),
			MaxSample:  1,
			IntervalId: 20,
		}
		query := types.QueryPerf{
			This:      *e.client.ServiceContent.PerfManager,
			QuerySpec: []types.PerfQuerySpec{querySpec},
		}

		response, err := methods.QueryPerf(e.ctx, e.client, &query)
		if err != nil {
			log.Fatal(err)
		}

		for _, base := range response.Returnval {
			metric := base.(*types.PerfEntityMetric)
			for _, baseSeries := range metric.Value {
				series := baseSeries.(*types.PerfMetricIntSeries)
				desc := countersInfoMap[int(series.Id.CounterId)]
				ch <- prometheus.MustNewConstMetric(desc,
					prometheus.GaugeValue, float64(series.Value[0]), hostName, series.Id.Instance, metric.Entity.Type)
			}
		}
	}
}

func main() {
	flag.Parse()

	exporter := NewExporter()
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, prometheus.Handler())

	log.Infof("Starting Server: %s", *listenAddress)
	http.ListenAndServe(*listenAddress, nil)
}
