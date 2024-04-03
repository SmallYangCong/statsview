package viewer

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
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
	p, _ := process.NewProcess(int32(os.Getpid()))

	graph := NewBasicView(VGCCPUFraction)
	graph.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "CPUFraction"}),
		charts.WithYAxisOpts(opts.YAxis{Name: "Percent", AxisLabel: &opts.AxisLabel{Show: true, Formatter: "{value} %", Rotate: 35}}),
	)
	graph.AddSeries("GC CPUFraction", []opts.LineData{})
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
	percent, _ := p.Percent(time.Second)
	return percent
}
