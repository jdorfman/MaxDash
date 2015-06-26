package main

import (
	ui 		"github.com/gizak/termui"
	volt	"github.com/bmconklin/maxcdn_volt"

	"log"
	"flag"
	"math"
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
	db := volt.Connect(config.DbAddr)

	err = ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	p := ui.NewPar("MaxCDN Dashboard - PRESS q TO QUIT")
	p.Height = 3
	p.Width = 50
	p.TextFgColor = ui.ColorWhite
	p.Border.Label = "Text Box"
	p.Border.FgColor = ui.ColorCyan

	strs := []string{"bootstrap.min.js", "bootstrap.min.css", "font-awesome.min.css", "bootswatch-2.1.1.min.css", "font-awesome.3.1.4-min.css", "bootstrap-mobile.min.css", "bootstrap-all.css", "bootswatch-2.4.3.min.css"}
	list := ui.NewList()
	list.Items = strs
	list.ItemFgColor = ui.ColorYellow
	list.Border.Label = "List"
	list.Height = 7
	list.Width = 25
	list.Y = 4

	g := ui.NewGauge()
	g.Percent = 50
	g.Width = 50
	g.Height = 3
	g.Y = 11
	g.Border.Label = "Gauge"
	g.BarColor = ui.ColorRed
	g.Border.FgColor = ui.ColorWhite
	g.Border.LabelFgColor = ui.ColorCyan

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

	sp := ui.NewSparklines(spark, spark1)
	sp.Width = 25
	sp.Height = 7
	sp.Border.Label = "Browsers"
	sp.Y = 4
	sp.X = 25

	sinps := (func() []float64 {
		n := 220
		ps := make([]float64, n)
		for i := range ps {
			ps[i] = 1 + math.Sin(float64(i)/5)
		}
		return ps
	})()

	lc := ui.NewLineChart()
	lc.Border.Label = "dot-mode Line Chart"
	lc.Data = sinps
	lc.Width = 50
	lc.Height = 11
	lc.X = 0
	lc.Y = 14
	lc.AxesColor = ui.ColorWhite
	lc.LineColor = ui.ColorRed | ui.AttrBold
	lc.Mode = "dot"

	bc := ui.NewBarChart()
	bcdata := []int{3, 2, 5, 3, 9, 5, 3, 2, 5, 8, 3, 2, 4, 5, 3, 2, 5, 7, 5, 3, 2, 6, 7, 4, 6, 3, 6, 7, 8, 3, 6, 4, 5, 3, 2, 4, 6, 4, 8, 5, 9, 4, 3, 6, 5, 3, 6}
	bclabels := []string{"S0", "S1", "S2", "S3", "S4", "S5"}
	bc.Border.Label = "Bar Chart"
	bc.Width = 26
	bc.Height = 10
	bc.X = 51
	bc.Y = 0
	bc.DataLabels = bclabels
	bc.BarColor = ui.ColorGreen
	bc.NumColor = ui.ColorBlack

	lc1 := ui.NewLineChart()
	lc1.Border.Label = "Cache Hit %"
	lc1.Data = sinps
	lc1.Width = 26
	lc1.Height = 11
	lc1.X = 51
	lc1.Y = 14
	lc1.AxesColor = ui.ColorWhite
	lc1.LineColor = ui.ColorYellow | ui.AttrBold

	p1 := ui.NewPar("Hey!\nI am a borderless block!")
	p1.HasBorder = false
	p1.Width = 26
	p1.Height = 2
	p1.TextFgColor = ui.ColorMagenta
	p1.X = 52
	p1.Y = 11

	data := &struct{
		Hits []int
		Bytes []int64
		Time []string
		CacheHits []int
		CachePerc []float64
	} {
		make([]int, 0),
		make([]int64, 0),
		make([]string, 0),
		make([]int, 0),
		make([]float64, 0),
	}
	draw := func(t int) {
		rawlogs := db.QueryUrls("SELECT * FROM urls WHERE ci = 1738 AND ti = '" + time.Now().UTC().Add(-10* time.Second).Format("2006-01-02 15:04:05") + "'")
		data.Hits = append(data.Hits, len(rawlogs))
		data.Time = append(data.Time, rawlogs[0].Ti.Format("2006-01-02 15:04:05"))
		var cacheHits int
		var bytes int64
		for _, l := range rawlogs {
			bytes += l.By_tr
			if l.Es == "HIT" {
				cacheHits++
			}
		}
		data.CachePerc = append(data.CachePerc,float64(cacheHits)/float64(len(rawlogs)))
		data.CacheHits = append(data.CacheHits,cacheHits)
		data.Bytes = append(data.Bytes, bytes)

		g.Percent = int(data.CachePerc[t] * float64(100))
		list.Items = data.Time
		sp.Lines[0].Data = spdata[:30+t%50]
		sp.Lines[1].Data = spdata[:35+t%50]
		lc.Data = sinps[t/2:]
		lc1.Data = data.CachePerc
		bc.Data = bcdata[t/2%10:]
		ui.Render(p, list, g, sp, lc, bc, lc1, p1)
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