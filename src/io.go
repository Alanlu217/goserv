package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/rs/zerolog/log"
)

type Io interface {
	SetPin(pin int8, value bool)
	TogglePin(pin int8)
}

type LogIo struct{}

func (LogIo) SetPin(pin int8, value bool) {
	log.Info().Msgf("set pin %d to %t", pin, value)
}

func (LogIo) TogglePin(pin int8) {
	log.Info().Msgf("toggling pin %d", pin)
}

type RPiIo struct{}

func (RPiIo) SetPin(pin int8, value bool) {
	valint := int8(0)
	if value {
		valint = 1
	}

	cmd := exec.Command("scripts/ioset.sh", fmt.Sprintf("%d=%d", pin, valint))
	go cmd.Run()
}

func (r RPiIo) TogglePin(pin int8) {
	r.SetPin(pin, true)
	time.Sleep(200 * time.Millisecond)
	r.SetPin(pin, false)
}
