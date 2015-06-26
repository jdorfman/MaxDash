package main

import (
	ui 		"github.com/gizak/termui"
	volt	"github.com/bmconklin/maxcdn_volt"

	"log"
	"flag"
	"time"
	"strconv"
	"strings"
)

var (
	confFile = flag.String("config", "./config.json", "path to config file")
	company = flag.Int("company", 4348, "company ID to look at")
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
	//bcdata := []int{3, 2, 5, 3, 9, 5, 3, 2, 5, 8, 3, 2, 4, 5, 3, 2, 5, 7, 5, 3, 2, 6, 7, 4, 6, 3, 6, 7, 8, 3, 6, 4, 5, 3, 2, 4, 6, 4, 8, 5, 9, 4, 3, 6, 5, 3, 6}
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
	//bc2data := []int{3, 2, 5, 3, 9, 5, 3, 2, 5, 8, 3, 2, 4, 5, 3, 2, 5, 7, 5, 3, 2, 6, 7, 4, 6, 3, 6, 7, 8, 3, 6, 4, 5, 3, 2, 4, 6, 4, 8, 5, 9, 4, 3, 6, 5, 3, 6}
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

	data := &struct{
		Hits []int
		Bytes []int64
		Time []string
		CacheHits []int
		CachePerc []float64
		Status []int
		Cont []int
		Chrome []int
		Safari []int
		PopUrl []string
	} {
		make([]int, 0),
		make([]int64, 0),
		make([]string, 0),
		make([]int, 0),
		make([]float64, 0),
		make([]int, 0),
		make([]int, 0),
		make([]int, 0),
		make([]int, 0),
		make([]string, 0),
	}
	draw := func(t int) {
		tm := time.Now().UTC().Add(-3 * time.Hour)
		urls := db.QueryUrls("SELECT * FROM urls WHERE company_id = " + strconv.Itoa(*company) + " AND time_window = '" + tm.Truncate(24*time.Hour).Format("2006-01-02 15:04:05") + "' ORDER BY hits DESC LIMIT 10")
		for _, u := range urls {
			data.PopUrl = append(data.PopUrl, u.Url)
		}
		rawlogs := db.QueryRawLogs("SELECT * FROM rawlogs WHERE ci = " + strconv.Itoa(*company) + " AND ti = '" + tm.Format("2006-01-02 15:04:05") + "'")
		data.Hits = append(data.Hits, len(rawlogs))
		data.Time = append(data.Time, tm.Format("2006-01-02 15:04:05"))
		var cacheHits int
		var bytes int64
		var h200 int
		var h206 int
		var h304 int 
		var h403 int 
		var h400 int 
		var h500 int
		var h499 int

		var na int
		var eu int
		var as int
		var oc int
		var sa int

		var chrome int
		var safari int

		for _, l := range rawlogs {
			bytes += l.By_tr
			if l.Es == "HIT" {
				cacheHits++
			}
			if l.Ss == 200 {
				h200++
			} else if l.Ss == 206 {
				h206++
			} else if l.Ss == 304 {
				h304++
			} else if l.Ss == 403 {
				h403++
			} else if l.Ss == 400{
				h400++
			} else if l.Ss == 500 {
				h500++
			} else if l.Ss == 499 {
				h499++
			}

			if l.Co == "NA" {
				na++
			} else if l.Co == "SA" {
				sa++
			} else if l.Co == "EU" {
				eu++
			} else if l.Co == "AS" {
				as++
			} else if l.Co == "OC" {
				oc++
			}

			if strings.Contains(l.Ua, "Chrome") {
				chrome++
			} else if strings.Contains(l.Ua, "Safari") {
				safari++
			}
		}
		data.CachePerc = append(data.CachePerc,float64(cacheHits)/float64(len(rawlogs)))
		data.CacheHits = append(data.CacheHits,cacheHits)
		data.Bytes = append(data.Bytes, bytes)
		data.Status = append(data.Status, h200)
		data.Status = append(data.Status, h206)
		data.Status = append(data.Status, h304)
		data.Status = append(data.Status, h400)
		data.Status = append(data.Status, h403)
		data.Status = append(data.Status, h500)
		data.Status = append(data.Status, h499)

		data.Cont = append(data.Cont, na)
		data.Cont = append(data.Cont, eu)
		data.Cont = append(data.Cont, as)
		data.Cont = append(data.Cont, oc)
		data.Cont = append(data.Cont, sa)

		data.Chrome = append(data.Chrome, chrome)
		data.Safari = append(data.Safari, safari)

		g.Percent = int(data.CachePerc[t] * float64(100))
		list.Items = data.PopUrl
		sp.Lines[0].Data = data.Chrome
		sp.Lines[1].Data = data.Safari
		bc.Data = data.Status
		bc2.Data = data.Cont
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