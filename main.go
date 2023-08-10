package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/speaker"
)

type model struct {
	streamer *beep.Ctrl
}

func NewModel(streamer *beep.Ctrl) model {
	speaker.Play(streamer)
	return model{streamer: streamer}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			speaker.Lock()
			m.streamer.Paused = !m.streamer.Paused
			speaker.Unlock()
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "Drip drop... (space to play/pause, q to quit)\n\n"
	switch m.streamer.Paused {
	case true:
		s += "Paused"
	default:
		s += "Playing"
	}
	return s
}

func main() {
	fRain, err := os.Open("resource/rain_loop.flac")
	if err != nil {
		log.Fatal(err)
	}

	streamerR, format, err := flac.Decode(fRain)
	if err != nil {
		log.Fatal(err)
	}

	// load file into memory
	buffer := beep.NewBuffer(format)
	buffer.Append(streamerR)

	streamerR.Close() // this closes the source file as well

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		log.Fatal(err)
	}

	rainLoop := beep.Loop(-1, buffer.Streamer(0, buffer.Len()))

	ctrl := &beep.Ctrl{Streamer: rainLoop, Paused: false}

	if _, err := tea.NewProgram(NewModel(ctrl)).Run(); err != nil {
		fmt.Printf("Bad thing happened: %s", err)
		os.Exit(1)
	}
}
