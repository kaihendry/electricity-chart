package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type ElecBill struct {
	Date     time.Time `json:"month"`
	Usage    float64   `json:"usage"`
	PriceKwh float64   `json:"price_kwh"`
}

func httpserver(w http.ResponseWriter, _ *http.Request) {

	var data []ElecBill
	data = append(data, ElecBill{Date: time.Date(2021, time.December, 1, 0, 0, 0, 0, time.UTC),
		Usage: 186.57, PriceKwh: 1.1526})
	data = append(data, ElecBill{Date: time.Date(2021, time.November, 1, 0, 0, 0, 0, time.UTC),
		Usage: 1218.95, PriceKwh: 0.3693})
	data = append(data, ElecBill{Date: time.Date(2021, time.October, 1, 0, 0, 0, 0, time.UTC),
		Usage: 1285.23, PriceKwh: 0.6343})

	// sort by Date
	sort.Slice(data, func(i, j int) bool {
		return data[i].Date.Before(data[j].Date)
	})
	log.Printf("%+v", data)

	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Electricity bill",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Month",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "kw/h price",
		}),
	)

	var times []string
	for _, b := range data {
		times = append(times, b.Date.Format("2006-01"))
	}
	bar.SetXAxis(times).AddSeries("bills", generateBarItems(data))

	bar.Render(w)
}

func generateBarItems(data []ElecBill) []opts.BarData {
	items := make([]opts.BarData, 0)
	for _, v := range data {
		items = append(items, opts.BarData{
			Value: v.PriceKwh,
		})
	}
	return items
}

func main() {
	// ignore favicon.ico
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, _ *http.Request) {})
	http.HandleFunc("/", httpserver)
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
}
