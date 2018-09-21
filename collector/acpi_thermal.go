// returns data points from MSAcpi_ThermalZoneTemperature
// https://msdn.microsoft.com/en-us/library/aa394317(v=vs.90).aspx
package collector

import (
	"strings"

	"github.com/StackExchange/wmi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

func init() {
	Factories["acpi_thermal"] = NewAcpiThermalCollector
}

// A MSAcpi_ThermalZoneCollector is a Prometheus collector for the WMI MSAcpi_ThermalZoneTemperature metric
type AcpiThermalCollector struct {
	CurrentTemperature *prometheus.Desc
}

func NewAcpiThermalCollector() (Collector, error) {
	const subsystem = "MSAcpi"
	return &AcpiThermalCollector{
		CurrentTemperature: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "current_temperature"),
			"Current motherboard temperature.",
			[]string{"current_temperature"},
			nil,
		),
	}, nil
}

// Collect sends the metric values for each metric
// to the provided prometheus Metric channel.
func (c *AcpiThermalCollector) Collect(ch chan<- prometheus.Metric) error {
	if desc, err := c.collect(ch); err != nil {
		log.Error("failed collecting acpi_thermal metrics:", desc, err)
		return err
	}
	return nil
}

type MSAcpi_ThermalZoneTemperaure struct {
	CurrentTemperature    uint32
}

func (c *AcpiThermalCollector) collect(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	var dst []MSAcpi_ThermalZoneTemperaure
	q := queryAll(&dst)
	if err := wmi.Query(q, &dst); err != nil {
		return nil, err
	}

	for _, data := range dst {

		ch <- prometheus.MustNewConstMetric(
			c.CurrentTemperature,
			prometheus.GaugeValue,
			float64(data.CurrentTemperature)/10-273.15,
			"current_temperature",
		)
	}

	return nil, nil
}
