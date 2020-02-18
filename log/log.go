// Package log provides a simple levelled logger outputing to stdout and a file.
package log

import (
	"fmt"
	"io"
	stdLog "log"
	"os"
	"path"
	"runtime"

	"github.com/kyoukaya/rhine/utils"

	"github.com/logrusorgru/aurora"
	"github.com/mattn/go-colorable"
)

// Logger is the levelled logger interface implemented by Log.
type Logger interface {
	Flush()
	Printf(string, ...interface{})
	Println(...interface{})
	Verbosef(string, ...interface{})
	Verboseln(...interface{})
	Warnf(string, ...interface{})
	Warnln(...interface{})
}

// Log provides a simple levelled logger outputing to stdout and a file.
type Log struct {
	fileLogger   *stdLog.Logger
	stdOutLogger *stdLog.Logger
	verbose      bool
}

// New sets up and returns a new instance of Logger.
func New(stdOut, verbose bool, filePath string, flags int) Logger {
	// Support for colored stdout output on windows.
	logger := &Log{verbose: verbose}
	var output io.Writer
	if runtime.GOOS == "windows" {
		output = colorable.NewColorableStdout()
	} else {
		output = os.Stdout
	}
	if stdOut {
		logger.stdOutLogger = stdLog.New(output, "", flags)
	}
	if filePath == "/dev/null" {
		return logger
	}
	if filePath == "" {
		filePath = utils.BinDir + "/logs/proxy.log"
	} else if !path.IsAbs(filePath) {
		filePath = utils.BinDir + "/" + filePath
	}
	dir := path.Dir(filePath)
	err := os.MkdirAll(dir, 0755)
	utils.Check(err)
	f, err := os.Create(filePath)
	utils.Check(err)
	logger.fileLogger = stdLog.New(f, "", flags)
	return logger
}

// Flush all buffers associated with the standard logger, if any.
func (log *Log) Flush() {}

func (log *Log) output(calldepth int, color func(interface{}) aurora.Value, prefix, str string) {
	if log == nil {
		return
	}
	calldepth++
	if log.fileLogger != nil {
		utils.Check(log.fileLogger.Output(calldepth, prefix+str))
	}
	if log.stdOutLogger != nil {
		if color != nil {
			prefix = color(prefix).String()
		}
		// Don't print long strings in stdout, truncate them to 400 chars.
		if len(str) > 403 {
			str = str[0:400] + "..."
		}
		utils.Check(log.stdOutLogger.Output(calldepth, prefix+str))
	}
}

// Verbose calls output to print to the standard logger, only if program is launched with
// the verbose flag.
func (log *Log) Verboseln(v ...interface{}) {
	if log.verbose {
		log.output(2, aurora.Blue, "INFO ", fmt.Sprintln(v...))
	}
}

// Verbosef calls output to print to the standard logger, only if program is launched with
// the verbose flag.
func (log *Log) Verbosef(format string, v ...interface{}) {
	if log.verbose {
		log.output(2, aurora.Blue, "VERB ", fmt.Sprintf(format, v...))
	}
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func (log *Log) Printf(format string, v ...interface{}) {
	log.output(2, aurora.Green, "INFO ", fmt.Sprintf(format, v...))
}

// Info calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Info.
func (log *Log) Println(v ...interface{}) {
	log.output(2, aurora.Green, "INFO ", fmt.Sprintln(v...))
}

// Warnf calls Output to print a warning to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func (log *Log) Warnf(format string, v ...interface{}) {
	log.output(2, aurora.Red, "WARN ", fmt.Sprintf(format, v...))
}

// Warnln calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Infoln.
func (log *Log) Warnln(v ...interface{}) {
	log.output(2, aurora.Red, "WARN ", fmt.Sprintln(v...))
}
