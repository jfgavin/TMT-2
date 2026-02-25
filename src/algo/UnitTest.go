package algo

import (
	"image/color"
	"math/rand"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func (syn BiexpSynapse) TestSynapse() {

	upperLimit := 100.0
	tickTotal := int(upperLimit / syn.Dt)

	output := make([]float64, tickTotal)
	var spikeTimes []float64

	for i := 0; i < tickTotal; i++ {

		t := float64(i) * syn.Dt

		input := rand.Intn(100) == 99
		if input {
			syn.Inject(1.0)
			spikeTimes = append(spikeTimes, t)
		}

		output[i] = syn.Advance()
	}

	// Create plot
	p := plot.New()
	p.Title.Text = "Biexponential Synapse"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Amplitude"

	// ---- Output Line ----
	pts := make(plotter.XYs, tickTotal)
	for i := range output {
		pts[i].X = float64(i) * syn.Dt
		pts[i].Y = output[i]
	}

	line, _ := plotter.NewLine(pts)
	p.Add(line)

	// ---- Spike Markers (Red) ----
	spikePts := make(plotter.XYs, len(spikeTimes))
	for i, t := range spikeTimes {
		spikePts[i].X = t
		index := int(t / syn.Dt)
		spikePts[i].Y = output[index]
	}

	scatter, _ := plotter.NewScatter(spikePts)
	scatter.GlyphStyle.Color = color.RGBA{R: 255, A: 255}
	scatter.GlyphStyle.Radius = vg.Points(3)
	p.Add(scatter)

	// ---- Threshold Line at 1.0 ----
	threshold := plotter.NewFunction(func(x float64) float64 {
		return 1.0
	})
	threshold.Color = color.RGBA{R: 200, G: 0, B: 0, A: 255}
	threshold.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}

	p.Add(threshold)

	p.Save(10*vg.Inch, 4*vg.Inch, "synapse_output.png")
}
