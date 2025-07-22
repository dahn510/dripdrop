package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/flac"
	"github.com/gopxl/beep/speaker"
)

//go:embed resource/rain_loop.flac
var rainSound []byte

type model struct {
	ctrl   *beep.Ctrl
	volume *effects.Volume
}

func NewModel(ctrl *beep.Ctrl, volume *effects.Volume) model {
	speaker.Play(volume)
	return model{ctrl: ctrl, volume: volume}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "+":
			speaker.Lock()
			m.volume.Volume += 0.1
			speaker.Unlock()
		case "-":
			speaker.Lock()
			m.volume.Volume -= 0.1
			speaker.Unlock()
		case " ":
			speaker.Lock()
			m.ctrl.Paused = !m.ctrl.Paused
			speaker.Unlock()
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	s := fmt.Sprintf("Drip drop... (space to play/pause, +/-: inc or dec volume, q to quit)\n\nVolume: %.0f ", m.volume.Volume*100)
	switch m.ctrl.Paused {
	case true:
		s += "Paused"
	default:
		s += "Playing"
	}
	return s
}

func main() {

	fRain := bytes.NewReader(rainSound)

	streamerR, format, err := flac.Decode(fRain)
	if err != nil {
		panic(err)
	}

	// load file into memory
	buffer := beep.NewBuffer(format)
	buffer.Append(streamerR)

	err = streamerR.Close() // this closes the source file as well
	if err != nil {
		panic(err)
	}

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		log.Fatal(err)
	}

	rainLoop := beep.Loop(-1, buffer.Streamer(0, buffer.Len()))

	ctrl := &beep.Ctrl{Streamer: rainLoop, Paused: false}

	volume := &effects.Volume{Streamer: ctrl, Base: 10, Volume: 0, Silent: false}

	if _, err := tea.NewProgram(NewModel(ctrl, volume)).Run(); err != nil {
		fmt.Printf("Bad thing happened: %s", err)
		os.Exit(1)
	}
}
