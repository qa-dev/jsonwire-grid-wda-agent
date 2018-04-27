//Package logger Добавляет имя файла и строчку кода из которой была вызванна функция логирования
package logger

import (
	"runtime"
	"strings"

	log "github.com/Sirupsen/logrus"
)

const (
	defSkipDepth = 2
	maxSkipDepth = 5
)

//Fields ...
type Fields log.Fields

func getFields(skip int) log.Fields {
	var pc uintptr
	var line int
	var ok bool

	/*
	   Так как некоторые пакеты в core тоже юзают лог, то полезнее будет видеть не эти пакеты - а тех кто их юзает
	   из реального приложения
	   пример go.avito.ru/gl/core/http/httputil
	   В старом варианте в логе было
	   function="go.avito.ru/gl/core/http/httputil.ValidateRequestTo" lineno=150
	   а хотелось
	   function="avito.ru/saver.(*HandlerSend).ServeHTTP"  lineno=55
	*/

	for i := defSkipDepth; i < maxSkipDepth; i++ {
		pc, _, line, ok = runtime.Caller(i)
		if ok == false {
			return log.Fields{}
		}
		f := runtime.FuncForPC(pc).Name()

		if strings.Contains(f, "go.avito.ru/gl/core") {
			continue
		}
		return log.Fields{"lineno": line, "function": runtime.FuncForPC(pc).Name()}
	}

	return log.Fields{"lineno": line, "function": runtime.FuncForPC(pc).Name()}

}

//WithField добавляет строчку кода и функция в дефолтный вызов logrus.WithField
func WithField(key string, value interface{}) *log.Entry {
	f := getFields(defSkipDepth)
	f[key] = value
	return log.WithFields(f)
}

//WithFields добавляет строчку кода и функция в дефолтный вызов logrus.WithFields
func WithFields(fields Fields) *log.Entry {
	f := getFields(defSkipDepth)
	for k, v := range fields {
		f[k] = v
	}
	return log.WithFields(f)
}

//WithError добавляет строчку кода и функция в дефолтный вызов logrus.WithError
func WithError(err error) *log.Entry {
	f := getFields(defSkipDepth)
	f[log.ErrorKey] = err
	return log.WithFields(f)
}

//Debugf добавляет строчку кода и функция в дефолтный вызов logrus.Debugf
func Debugf(format string, args ...interface{}) {
	log.WithFields(getFields(defSkipDepth)).Debugf(format, args...)
}

//Infof добавляет строчку кода и функция в дефолтный вызов logrus.Infof
func Infof(format string, args ...interface{}) {
	log.WithFields(getFields(defSkipDepth)).Infof(format, args...)
}

//Printf добавляет строчку кода и функция в дефолтный вызов logrus.Printf
func Printf(format string, args ...interface{}) {
	log.WithFields(getFields(defSkipDepth)).Printf(format, args...)
}

//Warnf добавляет строчку кода и функция в дефолтный вызов logrus.Warnf
func Warnf(format string, args ...interface{}) {
	log.WithFields(getFields(defSkipDepth)).Warnf(format, args...)
}

//Warningf добавляет строчку кода и функция в дефолтный вызов logrus.Warningf
func Warningf(format string, args ...interface{}) {
	log.WithFields(getFields(defSkipDepth)).Warnf(format, args...)
}

//Errorf добавляет строчку кода и функция в дефолтный вызов logrus.Errorf
func Errorf(format string, args ...interface{}) {
	log.WithFields(getFields(defSkipDepth)).Errorf(format, args...)
}

//Fatalf добавляет строчку кода и функция в дефолтный вызов logrus.Fatalf
func Fatalf(format string, args ...interface{}) {
	log.WithFields(getFields(defSkipDepth)).Fatalf(format, args...)
}

//Panicf добавляет строчку кода и функция в дефолтный вызов logrus.Panicf
func Panicf(format string, args ...interface{}) {
	log.WithFields(getFields(defSkipDepth)).Panicf(format, args...)
}

//Debug добавляет строчку кода и функция в дефолтный вызов logrus.Debug
func Debug(args ...interface{}) {
	log.WithFields(getFields(defSkipDepth)).Debug(args...)
}

//Info добавляет строчку кода и функция в дефолтный вызов logrus.Info
func Info(args ...interface{}) {
	log.WithFields(getFields(defSkipDepth)).Info(args...)
}

//Print добавляет строчку кода и функция в дефолтный вызов logrus.Print
func Print(args ...interface{}) {
	log.WithFields(getFields(defSkipDepth)).Print(args...)
}

//Warn добавляет строчку кода и функция в дефолтный вызов logrus.Warn
func Warn(args ...interface{}) {
	log.WithFields(getFields(defSkipDepth)).Warn(args...)
}

//Warning добавляет строчку кода и функция в дефолтный вызов logrus.Warning
func Warning(args ...interface{}) {
	log.WithFields(getFields(defSkipDepth)).Warning(args...)
}

//Error добавляет строчку кода и функция в дефолтный вызов logrus.Error
func Error(args ...interface{}) {
	log.WithFields(getFields(defSkipDepth)).Error(args...)
}

//Fatal добавляет строчку кода и функция в дефолтный вызов logrus.Fatal
func Fatal(args ...interface{}) {
	log.WithFields(getFields(defSkipDepth)).Fatal(args...)
}

//Panic добавляет строчку кода и функция в дефолтный вызов logrus.Panic
func Panic(args ...interface{}) {
	log.WithFields(getFields(defSkipDepth)).Panic(args...)
}

//Debugln добавляет строчку кода и функция в дефолтный вызов logrus.Debugln
func Debugln(args ...interface{}) {
	log.WithFields(getFields(defSkipDepth)).Debugln(args...)
}

//Infoln добавляет строчку кода и функция в дефолтный вызов logrus.Infoln
func Infoln(args ...interface{}) {
	log.WithFields(getFields(defSkipDepth)).Infoln(args...)
}

//Println добавляет строчку кода и функция в дефолтный вызов logrus.Println
func Println(args ...interface{}) {
	log.WithFields(getFields(defSkipDepth)).Println(args...)
}

//Warnln добавляет строчку кода и функция в дефолтный вызов logrus.Warnln
func Warnln(args ...interface{}) {
	log.WithFields(getFields(defSkipDepth)).Warnln(args...)
}

//Warningln добавляет строчку кода и функция в дефолтный вызов logrus.Warningln
func Warningln(args ...interface{}) {
	log.WithFields(getFields(defSkipDepth)).Warningln(args...)
}

//Errorln добавляет строчку кода и функция в дефолтный вызов logrus.Errorln
func Errorln(args ...interface{}) {
	log.WithFields(getFields(defSkipDepth)).Errorln(args...)
}

//Fatalln добавляет строчку кода и функция в дефолтный вызов logrus.Fatalln
func Fatalln(args ...interface{}) {
	log.WithFields(getFields(defSkipDepth)).Fatalln(args...)
}

//Panicln добавляет строчку кода и функция в дефолтный вызов logrus.Panicln
func Panicln(args ...interface{}) {
	log.WithFields(getFields(defSkipDepth)).Panicln(args...)
}
