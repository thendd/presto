package log

import "log"

func Info(text string, args ...any) {
	log.Printf(BLUE+"[INFO]"+RESET+" "+text+"\n", args...)
}

func Warn(text string, args ...any) {
	log.Printf(YELLOW+"[WARN]"+RESET+" "+text+"\n", args...)
}

func Error(text string, args ...any) {
	log.Printf(RED+"[ERROR]"+RESET+" "+text+"\n", args...)
}

func Fatal(text string, args ...any) {
	log.Printf(RED+"[FATAL]"+RESET+" "+text+"\n", args...)
}
