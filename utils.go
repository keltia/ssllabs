// utils.go
//
// Copyright 2018 © by Ollivier Robert <roberto@keltia.net>

package ssllabs

import (
	"fmt"
	"log"
	"time"
)

// Private area

func toDuration(t int) time.Duration {
	s := fmt.Sprintf("%ds", t)
	d, _ := time.ParseDuration(s)
	return d
}

// debug displays only if fDebug is set
func (c *Client) debug(str string, a ...interface{}) {
	if c.level >= 2 {
		log.Printf(str, a...)
	}
}

// debug displays only if fVerbose is set
func (c *Client) verbose(str string, a ...interface{}) {
	if c.level >= 1 {
		log.Printf(str, a...)
	}
}
