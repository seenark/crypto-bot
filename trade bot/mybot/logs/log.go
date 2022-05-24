package logs

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Log struct {
	fileName string
	file     *os.File
}

func NewLog(filename string) *Log { // "trades.log"
	workingDir, err := os.Getwd()
	log := Log{
		fileName: filename,
	}
	if err != nil {
		fmt.Println("Log", err)
	}
	tradeLogFile, err := os.OpenFile(fmt.Sprintf("%s/logs/%s", workingDir, log.fileName), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println("Cannot open file")
	}
	log.file = tradeLogFile
	return &log
}

func (l *Log) WriteText(msg string) {
	wrt := io.MultiWriter(os.Stdout, l.file)
	log.SetOutput(wrt)
	log.Println(msg)
}

func (l *Log) Close() {
	l.file.Close()
}
