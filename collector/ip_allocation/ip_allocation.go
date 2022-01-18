package collector

import (
	"context"
	"log"

	"github.com/mcbenjemaa/gke-metadata-exporter/pkg/gke"
	"github.com/mcbenjemaa/gke-metadata-exporter/pkg/k8s"
	"github.com/mcbenjemaa/gke-metadata-exporter/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/client-go/kubernetes"
)

const (
	prefix = "gke_info"
)

// KubernetesCollector represents GKE metadata
type IpAllocationCollector struct {
	gke gke.GKEMetadataFetcher
	k8sClient *kubernetes.Clientset

	AllocatableServiceIps *prometheus.Desc
	AllocatablePodIps    *prometheus.Desc
	ServiceIpAllocation *prometheus.Desc
	PodsIpAllocation    *prometheus.Desc
}

// NewKubernetesCollector creates a new MetadataCollector
func NewGKEIpAllocationCollector(gke gke.GKEMetadataFetcher, clientset *kubernetes.Clientset) *IpAllocationCollector {
	labelKeys := []string{
		"cluster",
		"zone",
		"projectId",
	}
	return &IpAllocationCollector{
		gke: gke,
		k8sClient: clientset,

		ServiceIpAllocation: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "capacity", "service_ip"),
			"Capacity of IP addresses of Services",
			labelKeys,
			nil,
		),
		PodsIpAllocation: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "capacity", "pods_ip"),
			"Capacity of IP addresses of Pods",
			labelKeys,
			nil,
		),
		AllocatableServiceIps: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "allocatable", "service_ip"),
			"Allocatable IP addresses of Services",
			labelKeys,
			nil,
		),
		AllocatablePodIps: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "allocatable", "pods_ip"),
			"Allocatable IP addresses of Pods",
			labelKeys,
			nil,
		),
	}
}

// Collect implements Prometheus' Collector interface and is used to collect metrics
func (c *IpAllocationCollector) Collect(ch chan<- prometheus.Metric) {
	ctx := context.Background()

	clusters, err := c.gke.FetchMetadata(ctx)
	if err != nil {
		log.Printf("unable to get clusters info %e", err)
	}

	for _, cluster := range clusters {
		g := c.gke.(gke.GKEClient)
		log.Printf("[IpAllocationCollector] cluster: %s", cluster.Name)
		ch <- prometheus.MustNewConstMetric(
			c.ServiceIpAllocation,
			prometheus.GaugeValue,
			float64(util.IPRangeSize(cluster.ServicesIPv4CIDR)),
			[]string{
				cluster.Name,
				cluster.Zone,
				g.ProjectID,
			}...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.PodsIpAllocation,
			prometheus.GaugeValue,
			float64(util.IPRangeSize(cluster.ContainerIPv4CIDR)),
			[]string{
				cluster.Name,
				cluster.Zone,
				g.ProjectID,
			}...,
		)

		svcCount, err := k8s.CountServices(ctx, c.k8sClient)
		if err != nil {
			log.Printf("unable to get k8s services count, %v", err)
		}
		ch <- prometheus.MustNewConstMetric(
			c.AllocatableServiceIps,
			prometheus.GaugeValue,
			float64(util.AllocatableIps(cluster.ServicesIPv4CIDR, svcCount)),
			[]string{
				cluster.Name,
				cluster.Zone,
				g.ProjectID,
			}...,
		)

		podsCount, err := k8s.CountPods(ctx, c.k8sClient)
		if err != nil {
			log.Printf("unable to get k8s pods count, %v", err)
		}
		ch <- prometheus.MustNewConstMetric(
			c.AllocatablePodIps,
			prometheus.GaugeValue,
			float64(util.AllocatableIps(cluster.ContainerIPv4CIDR, podsCount)),
			[]string{
				cluster.Name,
				cluster.Zone,
				g.ProjectID,
			}...,
		)
	}
}

// Describe implements Prometheus' Collector interface and is used to describe metrics
func (c *IpAllocationCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.ServiceIpAllocation
	ch <- c.PodsIpAllocation
}
