package plot

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"time"
)

const (
	imageWidth  = 12 * vg.Inch
	imageHeight = 6 * vg.Inch
)

// TemperatureData represents a single temperature reading.
type TemperatureData struct {
	Timestamp time.Time
	Value     float64
}

// createPlotPoints creates plot points from temperature data.
func createPlotPoints(data []TemperatureData) plotter.XYs {
	pts := make(plotter.XYs, len(data))
	for i, d := range data {
		pts[i].X = float64(d.Timestamp.Unix())
		pts[i].Y = d.Value
	}
	return pts
}

// createLine creates a line plot from plot points with custom styles.
func createLine(pts plotter.XYs) (*plotter.Line, error) {
	line, err := plotter.NewLine(pts)
	if err != nil {
		return nil, err
	}
	line.Color = plotutil.Color(1)
	line.Width = vg.Points(2)
	line.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	return line, nil
}

// createScatter creates scatter points with custom styles.
func createScatter(pts plotter.XYs) (*plotter.Scatter, error) {
	scatter, err := plotter.NewScatter(pts)
	if err != nil {
		return nil, err
	}
	scatter.GlyphStyle.Color = plotutil.Color(2)
	scatter.GlyphStyle.Radius = vg.Points(3)
	return scatter, nil
}

// PlotTemperatureGraph generates a temperature graph and saves it as a PNG image.
func PlotTemperatureGraph(data []TemperatureData, outputFile string) error {
	p := plot.New()
	p.Title.Text = "Temperature Over Time"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Temperature (°C)"

	// Add grid for better readability
	p.Add(plotter.NewGrid())

	// Create plot points
	pts := createPlotPoints(data)

	// Add line to the plot
	line, err := createLine(pts)
	if err != nil {
		return err
	}
	p.Add(line)

	// Add scatter points to the plot
	scatter, err := createScatter(pts)
	if err != nil {
		return err
	}
	p.Add(scatter)

	// Format the X axis as time
	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04"}

	// Save the plot as PNG image
	if err := p.Save(imageWidth, imageHeight, outputFile); err != nil {
		return err
	}
	return nil
}
