package main

import (
	"fmt"
	"os/exec"

	"github.com/rs/zerolog/log"
)

type Io interface {
	SetPin(pin int8, value bool)
}

type LogIo struct{}

func (LogIo) SetPin(pin int8, value bool) {
	log.Info().Msg(fmt.Sprintf("set pin %d to %t", pin, value))
}

type RPiIo struct{}

func (RPiIo) SetPin(pin int8, value bool) {
	valint := int8(0)
	if value {
		valint = 1
	}

	cmd := exec.Command("scripts/ioset.sh", fmt.Sprintf("GPIO%d=%d", pin, valint))
	go cmd.Run()
}
