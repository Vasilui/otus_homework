package main

import "fmt"

type Bar struct {
	percent int64  // percent - progress percentage
	cur     int64  // cur - current progress
	total   int64  // total - value for progress
	rate    string // rate - the actual progress bar to be printed
	graph   string // graph - the fill value for progress bar
}

func (bar *Bar) NewOption(start, total int64) {
	bar.cur = start
	bar.total = total
	bar.graph = placeholder
	bar.percent = bar.getPercent()
	for i := 0; i < int(bar.percent); i += 2 {
		bar.rate += bar.graph // initial progress position
	}
}

func (bar *Bar) getPercent() int64 {
	return int64((float32(bar.cur) / float32(bar.total)) * 100)
}

func (bar *Bar) Play(cur int64) {
	bar.cur = cur
	last := bar.percent
	bar.percent = bar.getPercent()
	if bar.percent != last && bar.percent%2 == 0 {
		bar.rate += bar.graph
	}
	fmt.Printf("\r[%-50s]%3d%% %8d/%d", bar.rate, bar.percent, bar.cur, bar.total)
}

func (bar *Bar) Finish() {
	fmt.Println()
}
