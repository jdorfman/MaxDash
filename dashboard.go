package main

import (
	ui "github.com/gizak/termui"

	"fmt"
	"log"
	"flag"
//	"math"
	"time"
)

var (
	confFile = flag.String("config", "./config.json", "path to config file")
)

func main() {
	flag.Parse()
	config, err := getConfig(*confFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Config Loaded: ", config)

	err = ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	p := ui.NewPar("       MaxCDN Dashboard - PRESS q TO QUIT ")
	p.Height = 3
	p.Width = 53
	p.TextFgColor = ui.ColorWhite
	p.Border.Label = ""
	p.Border.FgColor = ui.ColorYellow

	strs := []string{"bootstrap.min.js", "bootstrap.min.css", "font-awesome.min.css", "bootswatch-2.1.1.min.css", "font-awesome.3.1.4-min.css", "bootstrap-mobile.min.css", "bootstrap-all.css", "bootswatch-2.4.3.min.css"}
	list := ui.NewList()
	list.Items = strs
	list.ItemFgColor = ui.ColorYellow
	list.Border.Label = " Popular Files "
	list.Height = 10
	list.Width = 25
	list.Y = 4

	spark := ui.Sparkline{}
	spark.Height = 1
	spark.Title = "Chrome:"
	spdata := []int{4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6, 4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6, 4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6, 4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6}
	spark.Data = spdata
	spark.LineColor = ui.ColorCyan
	spark.TitleColor = ui.ColorWhite

	spark1 := ui.Sparkline{}
	spark1.Height = 1
	spark1.Title = "Safari:"
	spark1.Data = spdata
	spark1.TitleColor = ui.ColorWhite
	spark1.LineColor = ui.ColorRed

	// spark2 := ui.Sparkline{}
	// spark2.Height = 1
	// spark2.Title = "IE:"
	// spark2.Data = spdata
	// spark2.TitleColor = ui.ColorWhite
	// spark2.LineColor = ui.ColorMagenta

	// spark3 := ui.Sparkline{}
	// spark3.Height = 1
	// spark3.Title = "Firefox:"
	// spark3.Data = spdata
	// spark3.TitleColor = ui.ColorWhite
	// spark3.LineColor = ui.ColorYellow

	sp := ui.NewSparklines(spark, spark1)
	sp.Width = 25
	sp.Height = 10
	sp.Border.Label = " Browsers "
	sp.Y = 14
	sp.X = 0

	// sinps := (func() []float64 {
	// 	n := 220
	// 	ps := make([]float64, n)
	// 	for i := range ps {
	// 		ps[i] = 1 + math.Sin(float64(i)/5)
	// 	}
	// 	return ps
	// })()

	g := ui.NewGauge()
	g.Percent = 50
	g.Width = 30
	g.Height = 3
	g.Y = 20
	g.Border.Label = "Gauge"
	g.BarColor = ui.ColorRed
	g.Border.FgColor = ui.ColorWhite
	g.Border.LabelFgColor = ui.ColorCyan

	bc := ui.NewBarChart()
	bcdata := []int{3, 2, 5, 3, 9, 5, 3, 2, 5, 8, 3, 2, 4, 5, 3, 2, 5, 7, 5, 3, 2, 6, 7, 4, 6, 3, 6, 7, 8, 3, 6, 4, 5, 3, 2, 4, 6, 4, 8, 5, 9, 4, 3, 6, 5, 3, 6}
	bclabels := []string{"200", "206", "304", "403", "400", "500", "499"}
	bc.Border.Label = " Status Codes "
	bc.Width = 26
	bc.Height = 10
	bc.X = 26
	bc.Y = 4
	bc.DataLabels = bclabels
	bc.BarColor = ui.ColorGreen
	bc.NumColor = ui.ColorBlack

	bc2 := ui.NewBarChart()
	bc2data := []int{3, 2, 5, 3, 9, 5, 3, 2, 5, 8, 3, 2, 4, 5, 3, 2, 5, 7, 5, 3, 2, 6, 7, 4, 6, 3, 6, 7, 8, 3, 6, 4, 5, 3, 2, 4, 6, 4, 8, 5, 9, 4, 3, 6, 5, 3, 6}
	bc2labels := []string{"NA", "EU", "ASIA", "OC", "SA"}
	bc2.Border.Label = " Hits Per Region "
	bc2.Width = 23
	bc2.Height = 10
	bc2.X = 53
	bc2.Y = 4
	bc2.DataLabels = bc2labels
	bc2.BarColor = ui.ColorRed
	bc2.NumColor = ui.ColorBlack

	// lc1 := ui.NewLineChart()
	// lc1.Border.Label = " Cache Hit % "
	// lc1.Data = sinps
	// lc1.Width = 26
	// lc1.Height = 11
	// lc1.X = 51
	// lc1.Y = 14
	// lc1.AxesColor = ui.ColorWhite
	// lc1.LineColor = ui.ColorYellow | ui.AttrBold

	draw := func(t int) {
		g.Percent = t % 101
		list.Items = strs[t%9:]
		sp.Lines[0].Data = spdata[:30+t%50]
		sp.Lines[1].Data = spdata[:35+t%50]
		//lc.Data = sinps[t/2:]
		//lc1.Data = sinps[2*t:]
		bc.Data = bcdata[t/2%10:]
		bc2.Data = bc2data[t/2%10:]
		ui.Render(p, list, sp, bc, bc2)
	}

	evt := ui.EventCh()

	i := 0
	for {
		select {
		case e := <-evt:
			if e.Type == ui.EventKey && e.Ch == 'q' {
				return
			}
		default:
			draw(i)
			i++
			if i == 102 {
				return
			}
			time.Sleep(time.Second / 2)
		}
	}
}