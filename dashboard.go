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
	company = flag.Int("company", 4348, "company ID to look at") // company id should be pulled from config.json
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

// Theme Chooser, future state: via config.json

	//Colors
	ui.UseTheme("helloworld")

	//Use system default colors
	//ui.UseTheme("default")

// Start Dashboard

	p := ui.NewPar("                    MaxCDN Dashboard - PRESS q TO QUIT ")
	p.Height = 3
	p.Width = 76
	p.TextFgColor = ui.ColorCyan
	p.Border.FgColor = ui.ColorWhite

// Account Summary

	par2 := ui.NewPar("Account: jdorfman\nZones: 5\n- - - - - - - - - - - -\nTransfered: 784.99 TB\n- - - - - - - - - - - -\nHits: 35,090,636,610\nMisses: 89,799,462\nPercentage: 99.74%")
	par2.Height = 10
	par2.Width = 25
	par2.Y = 4
	par2.Border.Label = " 24 Hour Summary "
	par2.Border.FgColor = ui.ColorWhite
	par2.Border.LabelFgColor = ui.ColorCyan

// Popular Files List Chart

	list := ui.NewList()
	list.ItemFgColor = ui.ColorCyan
	list.Border.Label = " Popular Files "
	list.Height = 10
	list.Width = 76
	list.Y = 17
	list.Border.LabelFgColor = ui.ColorCyan

// Status Codes Bar Chart

	bc := ui.NewBarChart()
	bclabels := []string{"200", "206", "304", "403", "400", "500", "499"} // Should be configurable
	bc.Border.Label = " Status Codes "
	bc.Width = 26
	bc.Height = 10
	bc.X = 26
	bc.Y = 4
	bc.DataLabels = bclabels
	bc.BarColor = ui.ColorCyan
	bc.NumColor = ui.ColorBlack
	bc.Border.LabelFgColor = ui.ColorCyan

// Hits Per Region Bar Chart

	bc2 := ui.NewBarChart()
	bc2labels := []string{"NA", "EU", "ASIA", "OC", "SA"}
	bc2.Border.Label = " Hits Per Region "
	bc2.Width = 23
	bc2.Height = 10
	bc2.X = 53
	bc2.Y = 4
	bc2.DataLabels = bc2labels
	bc2.BarColor = ui.ColorCyan
	bc2.NumColor = ui.ColorBlack
	bc2.Border.LabelFgColor = ui.ColorCyan

// Cache Hit Percentage Gauge

	g2 := ui.NewGauge()
	g2.Percent = 50
	g2.Width = 76
	g2.Height = 3
	g2.Y = 14
	g2.Border.Label = " Cache Hit Percentage "
	g2.BarColor = ui.ColorCyan
	g2.Border.FgColor = ui.ColorWhite
	g2.Border.LabelFgColor = ui.ColorCyan
	g2.PercentColor = ui.ColorBlack

// Draw Funtion

	data := &struct{
		Hits []int
		Bytes []int64
		Time []string
		CacheHits []int
		CachePerc []float64
		Status []int
		Cont []int
		PopUrl []string
	} {
		make([]int, 0),
		make([]int64, 0),
		make([]string, 0),
		make([]int, 0),
		make([]float64, 0),
		make([]int, 0),
		make([]int, 0),
		make([]string, 0),
	}
	draw := func(t int) {
		//tm := time.Now().UTC().Add(-10*time.Second)
		tm, _ := time.Parse("2006-01-02 15:04:05", "2015-06-26 19:22:15")
		urls := db.QueryUrls("SELECT * FROM urls WHERE company_id = " + strconv.Itoa(*company) + " ORDER BY hits DESC LIMIT 10")
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

	
		g2.Percent = int(data.CachePerc[t] * float64(100))
		list.Items = data.PopUrl
		bc.Data = data.Status
		bc2.Data = data.Cont
		ui.Render(p, par2, list, g2, bc, bc2)

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