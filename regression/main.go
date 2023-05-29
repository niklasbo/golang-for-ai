package main

import (
	"fmt"
	"log"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

var (
	// x contains 'day in year' values
	x = []float64{1, 30, 60, 90, 120, 150, 180, 210, 240, 270, 300, 330, 365}
	// y contains the corresponding 'sunshine minutes on day' values
	y = []float64{280, 300, 330, 360, 380, 400, 420, 400, 380, 360, 330, 300, 280}
	// degree holds the regression degree
	degree = 2
)

func main() {
	vandermondeMatrix := Vandermonde(x, degree)
	yValuesColumn := mat.NewDense(len(y), 1, y)
	factorsCol := mat.NewDense(degree+1, 1, nil)

	var qr mat.QR
	qr.Factorize(vandermondeMatrix)

	err := qr.SolveTo(factorsCol, false, yValuesColumn)
	if err != nil {
		log.Fatalf("could not solve QR: %+v", err)
	}
	a, b, c := factorsCol.At(0, 0), factorsCol.At(1, 0), factorsCol.At(2, 0)
	fmt.Printf("y = %.3f + %.3f * x + %.3f * x^2", a, b, c)

	p := plot.New()
	p.Title.Text = "Lichtzeit im Jahr"
	p.X.Label.Text = "Tag im Jahr"
	p.Y.Label.Text = "Licht in Minuten"
	if err := plotutil.AddLinePoints(p,
		"gemessene Daten", getMeasuredPointsForPlotting(),
		"Regression", calculatePointsForPlotting(a, b, c),
	); err != nil {
		log.Fatal(err)
	}
	if err := p.Save(20*vg.Centimeter, 16*vg.Centimeter, "reg.png"); err != nil {
		log.Fatal(err)
	}
}

func getMeasuredPointsForPlotting() plotter.XYs {
	pointsToPlot := plotter.XYs{}
	for i, value := range x {
		nextPoint := plotter.XY{
			X: value,
			Y: y[i],
		}
		pointsToPlot = append(pointsToPlot, nextPoint)
	}
	return pointsToPlot
}

func calculatePointsForPlotting(a float64, b float64, c float64) plotter.XYs {
	pointsToPlot := plotter.XYs{}
	for i := 1; i < 366; i++ {
		nextPoint := plotter.XY{
			X: float64(i),
			Y: a + b*float64(i) + c*float64(i)*float64(i),
		}
		pointsToPlot = append(pointsToPlot, nextPoint)
	}
	return pointsToPlot
}

// Vandermonde creates a new Vandermonde matrix with the given degree
func Vandermonde(a []float64, degree int) *mat.Dense {
	d := degree + 1
	vandermondeMatrix := mat.NewDense(len(a), d, nil)
	for i := range a {
		for j, polynom := 0, 1.; j < d; j, polynom = j+1, polynom*a[i] {
			vandermondeMatrix.Set(i, j, polynom)
		}
	}
	return vandermondeMatrix
}
