package main

import (
	"fmt"
	"os/exec"

	"github.com/rs/zerolog/log"
)

type Io interface {
	Handle(dash DashRequest)
}

type LogIo struct{}

func (LogIo) Handle(dash DashRequest) {
	log.Info().Msg(fmt.Sprint(dash))
}

type RPiIo struct{}

func (RPiIo) setPin(pin int8, value bool) {
	valint := int8(0)
	if value {
		valint = 1
	}

	cmd := exec.Command("scripts/ioset.sh", fmt.Sprintf("GPIO%d=%d", pin, valint))
	go cmd.Run()
}

func (r RPiIo) Handle(dash DashRequest) {
	switch dash.Id {
	case "b1":
		log.Info().Msg("Turning LED ON")
		r.setPin(18, true)

	case "b2":
		log.Info().Msg("Turning LED OFF")
		r.setPin(18, false)
	}
}
