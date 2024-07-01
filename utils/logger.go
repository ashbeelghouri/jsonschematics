package utils

import (
	"encoding/json"
	"log"
)

type Logger struct {
	PrintDebugLogs bool
	PrintErrorLogs bool
	PreText        string
}

func (l *Logger) DEBUG(v ...interface{}) {
	if l.PrintDebugLogs {
		bytes, err := json.Marshal(v)
		if err != nil {
			log.Println(err)
		}
		log.Println(l.PreText, string(bytes))
	}
}

func (l *Logger) ERROR(v ...interface{}) {
	if l.PrintErrorLogs {
		bytes, err := json.Marshal(v)
		if err != nil {
			log.Println(err)
		}
		log.Println(l.PreText, "------ [ERROR] ------", string(bytes))
	}
}
