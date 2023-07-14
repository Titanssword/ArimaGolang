package arima

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strconv"
	"testing"

	"log"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func TestArima(t *testing.T) {
	// dataArray := []float64{2, 1, 2, 5, 2, 1, 2, 5, 2, 1, 2, 5, 2, 1, 2, 5}
	dataArray := []float64{2, 1, 2, 5, 2, 1, 2, 5, 2, 1, 2, 5, 2, 1, 2, 5}
	// trueForecast := []float64{2, 1, 2, 5}
	// Set ARIMA model parameters.
	p := 4
	d := 1
	q := 2
	P := 1
	D := 1
	Q := 0
	m := 0
	forecastSize := 10
	forecastResult := ForeCastARIMA(dataArray, forecastSize, NewConfig(p, d, q, P, D, Q, m))
	forecast := forecastResult.GetForecast()
	upper := forecastResult.GetForecastUpperConf()
	lower := forecastResult.GetForecastLowerConf()
	log.Println(forecast)
	log.Println(upper)
	log.Println(lower)

	// Create a new line chart.
	bar := charts.NewBar()

	// Set chart title and x-axis label.
	// line.SetGlobalOptions(
	// 	charts.WithTitleOpts(opts.Title{
	// 		Title: "ARIMA Model Forecast",
	// 	}),
	// 	charts.WithXAxisOpts(opts.XAxis{
	// 		Name: "Time",
	// 	}),
	// )

	// Add predicted values to the chart.
	// Put data into instance
	bar.SetXAxis([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}).
		AddSeries("Category A", generateBarItems(forecast)).
		AddSeries("Category B", generateBarItems(upper)).
		AddSeries("C", generateBarItems(lower))

	// Create a new file for the chart.
	file, err := os.Create("arima_forecast.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// Save the chart as an HTML file.
	bar.Render(file)
}

func TestArima2(t *testing.T) {
	// dataArray := []float64{2, 1, 2, 5, 2, 1, 2, 5, 2, 1, 2, 5, 2, 1, 2, 5}
	dataArray := []float64{2, 1, 2, 5, 2, 1, 2, 5, 2, 1, 2, 5, 2, 1, 2, 5}
	// trueForecast := []float64{2, 1, 2, 5}
	// Set ARIMA model parameters.
	p := 4
	d := 1
	q := 2
	P := 1
	D := 1
	Q := 0
	m := 0
	forecastSize := 10
	forecastResult := ForeCastARIMA(dataArray, forecastSize, NewConfig(p, d, q, P, D, Q, m))
	forecast := forecastResult.GetForecast()
	upper := forecastResult.GetForecastUpperConf()
	lower := forecastResult.GetForecastLowerConf()
	log.Println(forecast)
	log.Println(upper)
	log.Println(lower)

	// Create a new line chart.
	// create a new line instance
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Line example in Westeros theme",
			Subtitle: "Line chart rendered by the http server this time",
		}))

	// Put data into instance
	line.SetXAxis([]string{}).
		AddSeries("Category A", generateLineItems(append(dataArray, forecast...))).
		AddSeries("Category B", generateLineItems(append(dataArray, upper...))).
		AddSeries("Category C", generateLineItems(append(dataArray, lower...))).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))

	// Create a new file for the chart.
	file, err := os.Create("arima_forecast_line.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// Save the chart as an HTML file.
	line.Render(file)
}

// Helper function to generate line items for the chart.
func generateLineItems(data []float64) []opts.LineData {
	items := make([]opts.LineData, len(data))
	for i := range items {
		items[i] = opts.LineData{Value: data[i]}
	}
	return items
}

// generate random data for bar chart
func generateBarItems(data []float64) []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < 7; i++ {
		items = append(items, opts.BarData{Value: data[i]})
	}
	return items
}

var cscchris_val = []float64{
	2674.8060304978917, 3371.1788109723193, 2657.161969121835, 2814.5583226655367, 3290.855749923403, 3103.622791045206, 3403.2011487950185, 2841.438925235243, 2995.312700153925, 3256.4042898633224, 2609.8702933486843, 3214.6409110870877, 2952.1736018157644, 3468.7045537306344, 3260.9227206904898, 2645.5024256492215, 3137.857549381811, 3311.3526531674556, 2929.7762119375716, 2846.05991810631, 2606.47822546165, 3174.9770937667918, 3140.910443979614, 2590.6601484185085, 3123.4299821259915, 2714.4060964141136, 3133.9561758319487, 2951.3288157912752, 2860.3114228342765, 2757.4279640677833}
var cscchris_answer = []float64{
	3147.816496825682, 3418.2300802476093, 2856.905414401418, 3419.0312162705545, 3307.9803365878442, 3527.68377555284}

func dbl2str(value float64) string {
	return fmt.Sprintf("%.5f", value)
}

func commonTestCalculateRMSE(name string, trainingData []float64, trueForecastData []float64, forecastSize int, p, d, q, P, D, Q, m int) float64 {

	// Make forecast
	forecastResult := ForeCastARIMA(trainingData, forecastSize, NewConfig(p, d, q, P, D, Q, m))
	//Get forecast data and confidence intervals
	forecast := forecastResult.GetForecast()
	upper := forecastResult.GetForecastUpperConf()
	lower := forecastResult.GetForecastLowerConf()
	//Building output
	var sb bytes.Buffer
	sb.WriteString(name)
	sb.WriteString("  ****************************************************\n")
	sb.WriteString("Input Params { ")
	sb.WriteString("p: ")
	sb.WriteString(strconv.FormatInt(int64(p), 10))
	sb.WriteString(", d: ")
	sb.WriteString(strconv.FormatInt(int64(d), 10))
	sb.WriteString(", q: ")
	sb.WriteString(strconv.FormatInt(int64(q), 10))
	sb.WriteString(", P: ")
	sb.WriteString(strconv.FormatInt(int64(P), 10))
	sb.WriteString(", D: ")
	sb.WriteString(strconv.FormatInt(int64(D), 10))
	sb.WriteString(", Q: ")
	sb.WriteString(strconv.FormatInt(int64(Q), 10))
	sb.WriteString(", m: ")
	sb.WriteString(strconv.FormatInt(int64(m), 10))
	sb.WriteString(" }")
	sb.WriteString("\n\nFitted Model RMSE: ")
	sb.WriteString(fmt.Sprintf("%f", forecastResult.GetRMSE()))
	sb.WriteString("\n\n      TRUE DATA    |     LOWER BOUND          FORECAST       UPPER BOUND\n")

	for i := 0; i < len(forecast); i++ {
		sb.WriteString(dbl2str(trueForecastData[i]))
		sb.WriteString("    | ")
		sb.WriteString(dbl2str(lower[i]))
		sb.WriteString("   ")
		sb.WriteString(dbl2str(forecast[i]))
		sb.WriteString("   ")
		sb.WriteString(dbl2str(upper[i]))
		sb.WriteString("\n")
	}

	sb.WriteString("\n")

	//Compute RMSE against true forecast data
	temp := 0.0
	for i := 0; i < len(forecast); i++ {
		temp += math.Pow(forecast[i]-trueForecastData[i], 2)
	}
	rmse := math.Pow(temp/float64(len(forecast)), 0.5)
	sb.WriteString("RMSE = ")
	sb.WriteString(dbl2str(rmse))
	sb.WriteString("\n\n")
	log.Println(sb.String())
	return rmse
}
