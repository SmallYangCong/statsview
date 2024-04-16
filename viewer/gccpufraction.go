package viewer

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/process"
)

const (
	// VGCCPUFraction is the name of GCCPUFractionViewer
	VGCCPUFraction = "gccpufraction"
)

// GCCPUFractionViewer collects the GC-CPU fraction metric via `runtime.ReadMemStats()`
type GCCPUFractionViewer struct {
	smgr  *StatsMgr
	graph *charts.Line
	p     *process.Process
}

// NewGCCPUFractionViewer returns the GCCPUFractionViewer instance
// Series: Fraction
func NewGCCPUFractionViewer() Viewer {
	return NewGCCPUFractionViewerWithNumCPU()
}

// NewGCCPUFractionViewer returns the GCCPUFractionViewer instance
// Series: Fraction
func NewGCCPUFractionViewerWithNumCPU() Viewer {
	p, _ := process.NewProcess(int32(os.Getpid()))

	graph := NewBasicView(VGCCPUFraction)
	graph.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "CPUFraction"}),
		charts.WithYAxisOpts(opts.YAxis{Name: "Percent", AxisLabel: &opts.AxisLabel{Show: true, Formatter: "{value} %", Rotate: 35}}),
	)
	graph.AddSeries("GC CPUFraction", []opts.LineData{})
	graph.AddSeries("Server CPUFraction", []opts.LineData{})
	graph.AddSeries("App CPUFraction", []opts.LineData{})

	return &GCCPUFractionViewer{graph: graph, p: p}
}

func (vr *GCCPUFractionViewer) SetStatsMgr(smgr *StatsMgr) {
	vr.smgr = smgr
}

func (vr *GCCPUFractionViewer) Name() string {
	return VGCCPUFraction
}

func (vr *GCCPUFractionViewer) View() *charts.Line {
	return vr.graph
}

func (vr *GCCPUFractionViewer) Serve(w http.ResponseWriter, _ *http.Request) {
	vr.smgr.Tick()

	metrics := Metrics{
		Values: []float64{
			FixedPrecision(memstats.Stats.GCCPUFraction, 6),
			FixedPrecision(vr.getServerCPUFraction(), 6),
			FixedPrecision(vr.getAppCPUFraction(), 6),
		},
		Time: memstats.T,
	}

	bs, _ := json.Marshal(metrics)
	w.Write(bs)
}

func (vr *GCCPUFractionViewer) getAppCPUFraction() float64 {
	p := vr.p
	if p == nil {
		return 0.0
	}
	percent, _ := p.Percent(0)
	return percent / 100
}

func (vr *GCCPUFractionViewer) getServerCPUFraction() float64 {
	totalUsage, _ := cpu.PercentWithContext(context.Background(), time.Second, false)
	if len(totalUsage) == 0 {
		return 0.0
	}

	return totalUsage[0] / 100
}
