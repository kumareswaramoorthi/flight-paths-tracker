package logging

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kumareswaramoorthi/flight-paths-tracker/api/constants"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

const (
	SPRINT    = "sprint"
	SPRINT_F  = "sprintf"
	SPRINT_LN = "sprintln"
)

const (
	PANIC   = "panic"
	FATAL   = "fatal"
	ERROR   = "error"
	WARN    = "warn"
	WARNING = "warning"
	INFO    = "info"
	PRINT   = "print"
	DEBUG   = "debug"
	TRACE   = "trace"
)

type ApiLoggerEntry interface {
	WithError(err error) *apiLoggerEntry
	WithField(key string, value interface{}) *apiLoggerEntry
	WithFields(fields ApiLoggerFields) *apiLoggerEntry
	WithContext(ctx context.Context) *apiLoggerEntry
	SetFormatter(format string)
	SetLevel(lvl string) error

	Tracef(format string, value ...interface{})
	Debugf(format string, value ...interface{})
	Infof(format string, value ...interface{})
	Printf(format string, value ...interface{})
	Warnf(format string, value ...interface{})
	Warningf(format string, value ...interface{})
	Errorf(format string, value ...interface{})
	Fatalf(format string, value ...interface{})
	Panicf(format string, value ...interface{})

	Trace(value ...interface{})
	Debug(value ...interface{})
	Info(value ...interface{})
	Print(value ...interface{})
	Warn(value ...interface{})
	Warning(value ...interface{})
	Error(value ...interface{})
	Fatal(value ...interface{})
	Panic(value ...interface{})

	Traceln(value ...interface{})
	Debugln(value ...interface{})
	Infoln(value ...interface{})
	Println(value ...interface{})
	Warnln(value ...interface{})
	Warningln(value ...interface{})
	Errorln(value ...interface{})
	Fatalln(value ...interface{})
	Panicln(value ...interface{})
}

type ApiLoggerFields map[string]interface{}

type ApiLoggerHook logrus.Hook

type apiLoggerEntry struct {
	context  context.Context
	stdEntry *logrus.Entry
	Data     ApiLoggerFields
	err      string
}

func (l apiLoggerEntry) WithError(err error) *apiLoggerEntry {
	return l.WithField("error", err)
}

func (l apiLoggerEntry) WithField(key string, value interface{}) *apiLoggerEntry {
	return l.WithFields(ApiLoggerFields{key: value})
}

func (l apiLoggerEntry) AddHook(hk ApiLoggerHook) {
	l.stdEntry.Logger.AddHook(hk)
}

func (l apiLoggerEntry) WithFields(fields ApiLoggerFields) *apiLoggerEntry {
	logrusFields := make(logrus.Fields, len(fields))
	for k, v := range fields {
		logrusFields[k] = v
	}
	newEntry := l.stdEntry.WithFields(logrusFields)
	//span data
	data := make(ApiLoggerFields, len(l.Data)+len(fields))
	for k, v := range l.Data {
		data[k] = v
	}
	fieldErr := l.err
	for k, v := range fields {
		isErrField := false
		if t := reflect.TypeOf(v); t != nil {
			switch t.Kind() {
			case reflect.Func:
				isErrField = true
			case reflect.Ptr:
				isErrField = t.Elem().Kind() == reflect.Func
			}
		}
		if isErrField {
			tmp := fmt.Sprintf("can not add field %q", k)
			if fieldErr != "" {
				fieldErr = l.err + ", " + tmp
			} else {
				fieldErr = tmp
			}
		} else {
			data[k] = v
		}
	}
	return &apiLoggerEntry{context: l.context, stdEntry: newEntry, Data: data, err: fieldErr}
}

func (l apiLoggerEntry) WithContext(ctx context.Context) *apiLoggerEntry {
	return &apiLoggerEntry{
		context:  ctx,
		stdEntry: l.stdEntry.WithContext(ctx),
		Data:     l.Data,
		err:      l.err,
	}
}

//Happens only for logrus logs and not span logs
func (l apiLoggerEntry) SetFormatter(format string) {
	var formatter logrus.Formatter
	switch format {
	case constants.JSON:
		formatter = &logrus.JSONFormatter{}
		l.stdEntry.Logger.SetFormatter(formatter)
	}
}

func (l apiLoggerEntry) SetLevel(lvl string) error {
	setLevel := l.stdEntry.Logger.SetLevel
	lvlString := strings.ToLower(lvl)
	switch lvlString {
	case PANIC:
		setLevel(logrus.PanicLevel)
		return nil
	case FATAL:
		setLevel(logrus.FatalLevel)
		return nil
	case ERROR:
		setLevel(logrus.ErrorLevel)
		return nil
	case WARN, WARNING:
		setLevel(logrus.WarnLevel)
		return nil
	case INFO, PRINT:
		setLevel(logrus.InfoLevel)
		return nil
	case DEBUG:
		setLevel(logrus.DebugLevel)
		return nil
	case TRACE:
		setLevel(logrus.TraceLevel)
		return nil
	}
	return fmt.Errorf("not a valid optimusLogger Level: %q", lvl)
}

func (l apiLoggerEntry) Tracef(format string, value ...interface{}) {
	logToSpan(l, logrus.TraceLevel, SPRINT_F, format, value...)
	l.stdEntry.Tracef(format, value...)
}

func (l apiLoggerEntry) Debugf(format string, value ...interface{}) {
	logToSpan(l, logrus.DebugLevel, SPRINT_F, format, value...)
	l.stdEntry.Debugf(format, value...)
}

func (l apiLoggerEntry) Infof(format string, value ...interface{}) {
	logToSpan(l, logrus.InfoLevel, SPRINT_F, format, value...)
	l.stdEntry.Infof(format, value...)
}

func (l apiLoggerEntry) Printf(format string, value ...interface{}) {
	l.Infof(format, value...)
}

func (l apiLoggerEntry) Warnf(format string, value ...interface{}) {
	logToSpan(l, logrus.WarnLevel, SPRINT_F, format, value...)
	l.stdEntry.Warnf(format, value...)
}

func (l apiLoggerEntry) Warningf(format string, value ...interface{}) {
	l.Warnf(format, value...)
}

func (l apiLoggerEntry) Errorf(format string, value ...interface{}) {
	logToSpan(l, logrus.ErrorLevel, SPRINT_F, format, value...)
	l.stdEntry.Errorf(format, value...)
}

func (l apiLoggerEntry) Fatalf(format string, value ...interface{}) {
	//Always logToSpan first and then stdEntry for fatal
	logToSpan(l, logrus.FatalLevel, SPRINT_F, format, value...)
	l.stdEntry.Fatalf(format, value...)
}

func (l apiLoggerEntry) Panicf(format string, value ...interface{}) {
	logToSpan(l, logrus.PanicLevel, SPRINT_F, format, value...)
	l.stdEntry.Panicf(format, value...)
}

func (l apiLoggerEntry) Trace(value ...interface{}) {
	logToSpan(l, logrus.TraceLevel, SPRINT, "", value...)
	l.stdEntry.Trace(value...)
}

func (l apiLoggerEntry) Debug(value ...interface{}) {
	logToSpan(l, logrus.DebugLevel, SPRINT, "", value...)
	l.stdEntry.Debug(value...)
}

func (l apiLoggerEntry) Info(value ...interface{}) {
	logToSpan(l, logrus.InfoLevel, SPRINT, "", value...)
	l.stdEntry.Info(value...)
}

func (l apiLoggerEntry) Print(value ...interface{}) {
	l.Info(value...)
}

func (l apiLoggerEntry) Warn(value ...interface{}) {
	logToSpan(l, logrus.WarnLevel, SPRINT, "", value...)
	l.stdEntry.Warn(value...)
}

func (l apiLoggerEntry) Warning(value ...interface{}) {
	l.Warn(value...)
}

func (l apiLoggerEntry) Error(value ...interface{}) {
	logToSpan(l, logrus.ErrorLevel, SPRINT, "", value...)
	l.stdEntry.Error(value...)
}

func (l apiLoggerEntry) Fatal(value ...interface{}) {
	//Always logToSpan first and then stdEntry for fatal
	logToSpan(l, logrus.FatalLevel, SPRINT, "", value...)
	l.stdEntry.Fatal(value...)
}

func (l apiLoggerEntry) Panic(value ...interface{}) {
	logToSpan(l, logrus.PanicLevel, SPRINT, "", value...)
	l.stdEntry.Panic(value...)
}

func (l apiLoggerEntry) Traceln(value ...interface{}) {
	logToSpan(l, logrus.TraceLevel, SPRINT_LN, "", value...)
	l.stdEntry.Traceln(value...)
}

func (l apiLoggerEntry) Debugln(value ...interface{}) {
	logToSpan(l, logrus.DebugLevel, SPRINT_LN, "", value...)
	l.stdEntry.Debugln(value...)
}

func (l apiLoggerEntry) Infoln(value ...interface{}) {
	logToSpan(l, logrus.InfoLevel, SPRINT_LN, "", value...)
	l.stdEntry.Infoln(value...)
}

func (l apiLoggerEntry) Println(value ...interface{}) {
	l.Infoln(value...)
}

func (l apiLoggerEntry) Warnln(value ...interface{}) {
	logToSpan(l, logrus.WarnLevel, SPRINT_LN, "", value...)
	l.stdEntry.Warnln(value...)
}

func (l apiLoggerEntry) Warningln(value ...interface{}) {
	l.Warnln(value...)
}

func (l apiLoggerEntry) Errorln(value ...interface{}) {
	logToSpan(l, logrus.ErrorLevel, SPRINT_LN, "", value...)
	l.stdEntry.Errorln(value...)
}

func (l apiLoggerEntry) Fatalln(value ...interface{}) {
	//Always logToSpan first and then stdEntry for fatal
	logToSpan(l, logrus.FatalLevel, SPRINT_LN, "", value...)
	l.stdEntry.Fatalln(value...)
}

func (l apiLoggerEntry) Panicln(value ...interface{}) {
	logToSpan(l, logrus.PanicLevel, SPRINT_LN, "", value...)
	l.stdEntry.Panicln(value...)
}

func logToSpan(l apiLoggerEntry, level logrus.Level, printType string, format string, value ...interface{}) {
	var formattedValue string
	if l.stdEntry.Logger.IsLevelEnabled(level) {
		span := getSpan(l)

		switch printType {
		case SPRINT:
			formattedValue = fmt.Sprint(value...)
		case SPRINT_F:
			formattedValue = fmt.Sprintf(format, value...)
		case SPRINT_LN:
			//The following is same as logrus
			formattedValue = fmt.Sprintln(value...)
			formattedValue = formattedValue[:len(formattedValue)-1]
		}

		span.Annotate(getSpanAttributes(l, level), formattedValue)
	}
}

func getSpanAttributes(l apiLoggerEntry, level logrus.Level) []trace.Attribute {
	var logAttributes []trace.Attribute
	logAttributes = append(logAttributes, trace.StringAttribute("level", fmt.Sprint(level)))
	for key, data := range l.Data {
		logAttributes = append(logAttributes, trace.StringAttribute(key, fmt.Sprint(data)))
	}
	return logAttributes
}

func getSpan(l apiLoggerEntry) *trace.Span {
	ctx := l.context
	if ginContext, ok := ctx.(*gin.Context); ok {
		ctx = ginContext.Request.Context()
	}
	return trace.FromContext(ctx)
}

func newLoggerEntry(ctx context.Context, stdEntry *logrus.Entry) *apiLoggerEntry {
	return &apiLoggerEntry{context: ctx, stdEntry: stdEntry, err: ""}
}

func NewLoggerEntry() *apiLoggerEntry {
	ctx := context.TODO()
	return newLoggerEntry(ctx, logrus.StandardLogger().WithContext(ctx))
}

func GetLogger(c context.Context) *apiLoggerEntry {
	if c == nil {
		return NewLoggerEntry()
	}
	if ginContext, ok := c.(*gin.Context); ok {
		logger, ok := ginContext.Get(constants.LOGGER_KEY)
		if ok {
			return logger.(*apiLoggerEntry)
		}
	}
	return NewLoggerEntry()
}

func LoggingMiddleware(apiLogger *apiLoggerEntry) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(constants.LOGGER_KEY, apiLogger)
		c.Next()
	}
}
