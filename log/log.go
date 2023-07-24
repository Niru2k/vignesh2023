package logs

import (
	//Inbuild packages
	"fmt"
	"io"
	"os"

	//Third-party packages
	"github.com/sirupsen/logrus"
	call "github.com/t-tomalak/logrus-easy-formatter"
)

func Logs() *logrus.Logger {
	file, err := os.OpenFile("logs.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create a log file")
	}

	log := &logrus.Logger{
		Out:   io.MultiWriter(file, os.Stdout),
		Level: logrus.DebugLevel,
		Formatter: &call.Formatter{
			TimestampFormat: "02-01-2006 15:04:05",
			LogFormat:       "[%lvl%]: %time% - %msg%\n",
		},
	}
	return log
}
