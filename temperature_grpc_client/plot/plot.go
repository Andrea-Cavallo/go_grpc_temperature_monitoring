package plot

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"time"
)

// TemperatureData rappresenta una singola lettura di temperatura.
type TemperatureData struct {
	Timestamp time.Time
	Value     float64
}

// PlotTemperatureGraph genera un grafico delle temperature e lo salva come immagine PNG.
func PlotTemperatureGraph(data []TemperatureData, outputFile string) error {
	// Crea un nuovo grafico
	p := plot.New()
	p.Title.Text = "Temperature Over Time"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Temperature (°C)"

	// Crea un set di punti per il grafico
	pts := make(plotter.XYs, len(data))
	for i, d := range data {
		pts[i].X = float64(d.Timestamp.Unix())
		pts[i].Y = d.Value
	}

	// Aggiungi una linea al grafico
	line, err := plotter.NewLine(pts)
	if err != nil {
		return err
	}
	p.Add(line)

	// Aggiungi i punti al grafico
	points, err := plotter.NewScatter(pts)
	if err != nil {
		return err
	}
	p.Add(points)

	// Format the X axis as time
	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04"}

	// Salva il grafico come immagine PNG
	if err := p.Save(10*vg.Inch, 4*vg.Inch, outputFile); err != nil {
		return err
	}

	return nil
}
