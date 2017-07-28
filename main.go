package main

import (
	_ "github.com/gonum/floats"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	df "github.com/kniren/gota/dataframe"
	_ "github.com/sajari/regression"

	"fmt"
	_ "io/ioutil"
	"log"
	_ "math"
	"os"
)

func createDF(dataFile string) df.DataFrame {
	// Open the CSV dataset file
	f, err := os.Open(dataFile)
	if err != nil {
		fmt.Println(err.Error())
	}
	// Close the file
	defer f.Close()

	// Create a dataframe from the csv string
	// The types of colomns will be inferred
	dataDf := df.ReadCSV(f)
	return dataDf
}

func creatHistogram(dataFrame df.DataFrame) {
	// Create a Histogram for each of the features
	for _, colName := range dataFrame.Names() {

		//Extract the columns as a slice of floats
		// Float() converts []series.Series to []float values
		floatCol := dataFrame.Col(colName).Float()

		// Create a plotter.Values value and fill it with the
		// values from the respective column of the dataframe
		plotVals := make(plotter.Values, len(floatCol))
		summaryVals := make([]float64, len(floatCol))
		for i, floatVal := range floatCol {
			plotVals[i] = floatVal
			summaryVals[i] = floatVal
		}

		// Make a plot and set its title
		p, err := plot.New()
		if err != nil {
			log.Fatal(err)
		}
		p.Title.Text = fmt.Sprintf("Hisogram of a %s", colName)

		// Create a histogram of our values drawn
		// from the standard normal
		h, err := plotter.NewHist(plotVals, 16)
		if err != nil {
			log.Fatal(err)
		}

		// Normaliz histogram
		h.Normalize(1)

		// Add histogram to the plot
		p.Add(h)

		// Save the plot to a png file
		if err := p.Save(4*vg.Inch, 4*vg.Inch, colName+"_hist.png"); err != nil {
			log.Fatal(err)
		}
	}
}

func createScat(dataFrame df.DataFrame, yColumn string) {
	// Extract the target column
	yVals := dataFrame.Col(yColumn).Float()

	// Create a scattered plot for each feature in the dataset
	for _, colName := range dataFrame.Names() {

		// Extract the columns as a slice of floats
		floatCol := dataFrame.Col(colName).Float()

		// pts will hold the values for plotting
		pts := make(plotter.XYs, len(floatCol))

		// Fill pts with data
		// Loop on floatCol and map each value to
		// the corresponding Y value
		for i, floatVal := range floatCol {
			pts[i].X = floatVal
			pts[i].Y = yVals[i]
		}

		// Create the plot
		p, err := plot.New()
		if err != nil {
			log.Fatal(err)
		}

		p.X.Label.Text = colName
		p.Y.Label.Text = yColumn
		//p.Add(plotter.NewGrid())

		s, err := plotter.NewScatter(pts)
		if err != nil {
			log.Fatal(err)
		}
		s.GlyphStyle.Radius = vg.Points(3)

		// Save the plot in PNG file
		p.Add(s)
		if err := p.Save(4*vg.Inch, 4*vg.Inch, colName+"_scatter.png"); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	newDataFrame := createDF("diabetes.csv")
	fmt.Printf("Num of Rows: %d\nNum of Columsn: %d\nColumn Names: %s\n\n", newDataFrame.Nrow(), newDataFrame.Ncol(), newDataFrame.Names())
	fmt.Println(newDataFrame.Select([]string{"bmi", "ltg", "y"}).Subset([]int{0, 1, 2, 3, 4}))
	creatHistogram(newDataFrame.Select([]string{"bmi", "ltg", "y"}))
	createScat(newDataFrame.Select([]string{"bmi", "ltg"}), "y")
}
