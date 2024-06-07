package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
	"log"

	"github.com/nsf/termbox-go"
	"gopkg.in/yaml.v2"
)

const (
	usage = `
 countdowndsm <pathtoconfig>

 Usage
	countdowndsm /path/to/config.yml

 Flags
`
	tick         = time.Second
	inputDelayMS = 500 * time.Millisecond
)

var (
	timer          	*time.Timer
	ticker         	*time.Ticker
	queues         	chan termbox.Event
	w, h           	int
	inputStartTime 	time.Time
	actTime				 	string
	actName					string
)


type yamlData struct {
	Acts []map[string]interface{} `yaml:'acts'`
	Persons []string `yaml:'persons'`
}

func ReadConfig() yamlData {
	var config yamlData

	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		stderr(usage)
		flag.PrintDefaults()
		os.Exit(2)
	}

	// Open YAML file
	file, err := os.Open(args[0])
	if err != nil {
			log.Println(err.Error())
	}
	defer file.Close()

	// Decode YAML file to struct
	if file != nil {
			decoder := yaml.NewDecoder(file)
			if err := decoder.Decode(&config); err != nil {
					log.Println(err.Error())
			}
	}

	return config
}

func parseTime(date string) (time.Duration, error) {
	targetTime, err := time.Parse(time.Kitchen, strings.ToUpper(date))
	if err != nil {
		targetTime, err = time.Parse("15:04", date)
		if err != nil {
			return time.Duration(0), err
		}
	}

	now := time.Now()
	originTime := time.Date(0, time.January, 1, now.Hour(), now.Minute(), now.Second(), 0, time.UTC)

	// The time of day has already passed, so target tomorrow.
	if targetTime.Before(originTime) {
		targetTime = targetTime.AddDate(0, 0, 1)
	}

	duration := targetTime.Sub(originTime)

	return duration, err
}

func main() {

	config := ReadConfig()

	fmt.Println(config.Acts)

	for actIndex, act := range config.Acts {

		actTime, _ := act["time"].(string)
		actName, _ := act["name"].(string)

		fmt.Println(actTime)
		fmt.Println(actName)

		timeLeft, err := parseTime(actTime)

		if err != nil {
			timeLeft, err = time.ParseDuration(actTime)
			if err != nil {
				stderr("error: Time: invalid duration or time: acts[%v]\n", actIndex)
				os.Exit(2)
			}
		}

		// Clean terminal
		err = termbox.Init()
		if err != nil {
			panic(err)
		}

		queues = make(chan termbox.Event)
		go func() {
			for {
				queues <- termbox.PollEvent()
			}
		}()

		countdown(timeLeft, actName)
	}

}

func countdown(totalDuration time.Duration, actTitle string) {
	timeLeft := totalDuration
	title := actTitle
	var exitCode int
	w, h = termbox.Size()
	start(timeLeft)

	draw(title, timeLeft, w, h)

loop:
	for {
		select {
		case ev := <-queues:
			// Ctrl+C/Esc
			if ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC {
				exitCode = 1
				break loop
			}

		case <-ticker.C:
			timeLeft -= tick
			draw(title, timeLeft, w, h)
		case <-timer.C:
			break loop
		}
	}

	termbox.Close()
	if exitCode != 0 {
		os.Exit(exitCode)
	}
}

func start(d time.Duration) {
	timer = time.NewTimer(d)
	ticker = time.NewTicker(tick)
}

func format(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	if h < 1 {
		return fmt.Sprintf("%02d:%02d", m, s)
	}
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}


func draw(t string, d time.Duration, w int, h int) {
	clear()

	str := format(d)
	timerText := toText(str)
	titleStatus := toTextSmall(strings.ToLower(t))

	xTitle, yTitle, xTimer, yTimer := w/2-titleStatus.width()/2, h/2-timerText.height()/2-2-titleStatus.height(), w/2-timerText.width()/2, h/2-timerText.height()/2

	for _, symbolTitle := range titleStatus {
		echo_symbol(symbolTitle, xTitle, yTitle)
		xTitle += symbolTitle.width()
	}

	for _, symbolTimer := range timerText {
		echo_symbol(symbolTimer, xTimer, yTimer)
		xTimer += symbolTimer.width()
	}

	flush()
}
