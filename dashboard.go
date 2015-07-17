package main

import (
	ui 		"github.com/gizak/termui"
	volt	"github.com/bmconklin/maxcdn_volt"

	"log"
	"flag"
	"time"
	"strconv"
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


// Start Dashboard

	p := ui.NewPar("                    MaxCDN Dashboard - PRESS q TO QUIT ")
	p.Height = 3
	p.Width = 76
	p.TextFgColor = ui.ColorWhite
	p.Border.Label = ""
	p.Border.FgColor = ui.ColorYellow

// Line Chart

	lc0 := ui.NewLineChart()
	lc0.Border.Label = " Hits Per Second "
	lc0.Data = []float64{}
	lc0.Width = 25
	lc0.Height = 10
	lc0.X = 0
	lc0.Y = 4
	lc0.AxesColor = ui.ColorWhite
	lc0.LineColor = ui.ColorGreen | ui.AttrBold

// Popular Files List Chart

	strs := []string{"bootstrap.min.js", "bootstrap.min.css", "font-awesome.min.css", "bootswatch-2.1.1.min.css", "font-awesome.3.1.4-min.css", "bootstrap-mobile.min.css", "bootstrap-all.css", "bootswatch-2.4.3.min.css"}
	list := ui.NewList()
	list.Items = strs
	list.ItemFgColor = ui.ColorYellow
	list.Border.Label = " Popular Files "
	list.Height = 10
	list.Width = 76
	list.Y = 17

// Status Codes Bar Chart

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

// Hits Per Region Bar Chart

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

// Cache Hit Percentage Gauge

	g2 := ui.NewGauge()
	g2.Percent = 50
	g2.Width = 76
	g2.Height = 3
	g2.Y = 14
	g2.Border.Label = " Cache Hit Percentage "
	g2.BarColor = ui.ColorRed
	g2.Border.FgColor = ui.ColorWhite
	g2.Border.LabelFgColor = ui.ColorWhite

// Draw Funtion

	data := &struct{
		Hits int
		Bytes int64
		Time string
		CacheHits []int
		CachePerc []float64
		Status []int
		Cont []int
		PopUrl []string
	} {
		0,
		int64(0),
		"",
		make([]int, 0),
		make([]float64, 0),
		make([]int, 0),
		make([]int, 0),
		make([]string, 0),
	}
	draw := func(t int) {
		tm := time.Now().UTC().Add(-10*time.Second)
		resp, err := db.QueryAll("SELECT count(*) as hits, sc, hn, ui FROM rawlogs WHERE ci = " + strconv.Itoa(*company) + " AND ti = '" + tm.Format("2006-01-02 15:04:05") + "' GROUP BY hn, sc, ui ORDER BY hits DESC LIMIT 10")
		if err == nil {
			data.PopUrl = make([]string, 0)
			url := struct{
				Hits 	int
				Sc 		string
				Hn 		string
				Ui 		string
			} {}
			for resp.Table(0).HasNext() {
				resp.Table(0).Next(&url)
				data.PopUrl = append(data.PopUrl, url.Sc + "://" + url.Hn + url.Ui)
			}
		} else {
			log.Println(err)
		}

		rawlogs := db.QueryRawLogs("SELECT * FROM rawlogs WHERE ci = " + strconv.Itoa(*company) + " AND ti = '" + tm.Format("2006-01-02 15:04:05") + "'")
		data.Hits = len(rawlogs)
		data.Time = tm.Format("2006-01-02 15:04:05")
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
		}
		lc0.Data = append(lc0.Data, float64(bytes))
		if len(lc0.Data) > 26 {
			lc0.Data = lc0.Data[1:]
		}
		data.CachePerc = append(data.CachePerc,float64(cacheHits)/float64(len(rawlogs)))
		data.CacheHits = append(data.CacheHits,cacheHits)

		data.Status = make([]int, 0)
		data.Status = append(data.Status, h200)
		data.Status = append(data.Status, h206)
		data.Status = append(data.Status, h304)
		data.Status = append(data.Status, h400)
		data.Status = append(data.Status, h403)
		data.Status = append(data.Status, h500)
		data.Status = append(data.Status, h499)

		data.Cont = make([]int, 0)
		data.Cont = append(data.Cont, na)
		data.Cont = append(data.Cont, eu)
		data.Cont = append(data.Cont, as)
		data.Cont = append(data.Cont, oc)
		data.Cont = append(data.Cont, sa)


		g2.Percent = int(data.CachePerc[t] * float64(100))
		list.Items = data.PopUrl
		bc.Data = data.Status
		bc2.Data = data.Cont
		ui.Render(p, list, g2, bc, bc2, lc0)

	}

	evt := ui.EventCh()

	tick := time.Tick(1 * time.Second)
	i := 0
	for {
		select {
		case e := <-evt:
			if e.Type == ui.EventKey && e.Ch == 'q' {
				return
			}
		case <- tick:
			draw(i)
			i++
		}
	}
}
