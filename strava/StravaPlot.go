package strava

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func PlotHistogram(data []StravaListActivitiesResponse) {
	p := plot.New()
	p.Title.Text = "histogram plot"
	pts := make(plotter.XYs, len(data))
	for k, v := range data {
		pts[k].X = float64(len(data) - k)
		pts[k].Y = float64(v.DistanceKm)
	}
	p.Title.Text = "Running distance trend"
	p.X.Label.Text = "Session"
	p.Y.Label.Text = "Distance (km)"

	err := plotutil.AddLinePoints(p,
		"Distance", pts)
	if err != nil {
		panic(err)
	}
	if err := p.Save(6*vg.Inch, 4*vg.Inch, "web/images/plotXY.png"); err != nil {
		panic(err)
	}
}
