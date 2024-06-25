package utils

import (
	"encoding/json"
	"log"
)

type Logger struct {
	PrintDebugLogs bool
	PrintErrorLogs bool
}

func (l *Logger) DEBUG(v ...interface{}) {
	if l.PrintDebugLogs {
		bytes, err := json.Marshal(v)
		if err != nil {
			log.Println(err)
		}
		log.Println(string(bytes))
	}
}

func (l *Logger) ERROR(v ...interface{}) {
	if l.PrintErrorLogs {
		bytes, err := json.Marshal(v)
		if err != nil {
			log.Println(err)
		}
		log.Println("------ [ERROR] ------", string(bytes))
	}
}
