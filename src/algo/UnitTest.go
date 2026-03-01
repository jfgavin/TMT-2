package algo

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func (syn BiexpSynapse) PlotSynapse() {
	// Create plot
	p := plot.New()
	p.Title.Text = "Biexponential Synapse"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Amplitude"

	// ---- Output Line ----
	pts := make(plotter.XYs, len(syn.Output))
	for i := range syn.Output {
		pts[i].X = float64(i) * syn.Dt
		pts[i].Y = syn.Output[i]
	}

	line, _ := plotter.NewLine(pts)
	p.Add(line)

	p.Save(10*vg.Inch, 4*vg.Inch, "synapse_output.png")
}
