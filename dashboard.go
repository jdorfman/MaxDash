package main

import (
	ui "github.com/gizak/termui"

	"fmt"
	"log"
	"flag"
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

// Start Dashboard

	p := ui.NewPar("                    MaxCDN Dashboard - PRESS q TO QUIT ")
	p.Height = 3
	p.Width = 76
	p.TextFgColor = ui.ColorWhite
	p.Border.Label = ""
	p.Border.FgColor = ui.ColorYellow

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

// Hits Per Region Bar Chart

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

	draw := func(t int) {
		g2.Percent = t % 101
		list.Items = strs[t%9:]
		bc.Data = bcdata[t/2%10:]
		bc2.Data = bc2data[t/2%10:]
		ui.Render(p, list, bc, bc2, g2)
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