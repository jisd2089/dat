package logs

import (
	"drcs/settings"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
)

type NewHookFunc func(se settings.Settings) (logrus.Hook, error)
type NewFormatterFunc func(se settings.Settings) (logrus.Formatter, error)
type NewAppenderFunc func(se settings.Settings) (io.Writer, error)

var (
	newHookFuncRegistry      = make(map[string]NewHookFunc)
	newFormatterFuncRegistry = make(map[string]NewFormatterFunc)
	newAppenderFuncRegistry  = make(map[string]NewAppenderFunc)
)

var (
	hooks      = make(map[string]logrus.Hook)
	formatters = make(map[string]logrus.Formatter)
	appenders  = make(map[string]io.Writer)
	loggers    = make(map[string]*logrus.Logger)
)

func Initialize(configPath string) error {
	if configPath != "" {
		se, err := settings.CreateSettingsFromYAML(configPath)
		if err != nil {
			return err
		}

		initHooks(se)
		initFormatters(se)
		initAppenders(se)
		initRootLogger(se)
		initLoggers(se)
	}
	return nil
}

func initHooks(se settings.Settings) {
	hooksMap, err := se.GetMap("Hooks")
	if err == settings.Nil {
		return
	}
	for name := range hooksMap {
		hookName := name.(string)
		hookType, err := se.GetString(fmt.Sprintf("Hooks.%s.Type", hookName))
		if err != nil {
			logrus.Panic(err)
		}

		options, err := se.GetSettings(fmt.Sprintf("Hooks.%s.Options", hookName))
		if err != nil && err != settings.Nil {
			logrus.Panic(err)
		}

		newFunc := newHookFuncRegistry[hookType]
		if newFunc == nil {
			logrus.Panicf("Unknown hook type %s", hookType)
		}

		hook, err := newFunc(options)
		if err != nil {
			logrus.Panic(err)
		}
		hooks[hookName] = hook
	}
}

func initFormatters(se settings.Settings) {
	// 默认添加json、text
	formatters["json"] = &logrus.JSONFormatter{}
	formatters["text"] = &logrus.TextFormatter{}

	formatterMap, err := se.GetMap("Formatters")
	if err == settings.Nil {
		return
	}
	for name := range formatterMap {
		formatterName := name.(string)
		formatterType, err := se.GetString(fmt.Sprintf("Formatters.%s.Type", formatterName))
		if err != nil {
			logrus.Panic(err)
		}

		options, err := se.GetSettings(fmt.Sprintf("Formatters.%s.Options", formatterName))
		if err != nil && err != settings.Nil {
			logrus.Panic(err)
		}

		newFunc := newFormatterFuncRegistry[formatterType]
		if newFunc == nil {
			logrus.Panicf("Unknown formatter type %s", formatterType)
		}

		formatter, err := newFunc(options)
		if err != nil {
			logrus.Panic(err)
		}
		formatters[formatterName] = formatter
	}
}

func initAppenders(se settings.Settings) {
	appenderMap, err := se.GetMap("Appenders")
	if err == settings.Nil {
		return
	}
	for name := range appenderMap {
		appenderName := name.(string)
		appenderType, err := se.GetString(fmt.Sprintf("Appenders.%s.Type", appenderName))
		if err != nil {
			logrus.Panic(err)
		}

		options, err := se.GetSettings(fmt.Sprintf("Appenders.%s.Options", appenderName))
		if err != nil && err != settings.Nil {
			logrus.Panic(err)
		}

		newFunc := newAppenderFuncRegistry[appenderType]
		if newFunc == nil {
			logrus.Panicf("Unknown appender type %s", appenderType)
		}

		appender, err := newFunc(options)
		if err != nil {
			logrus.Panic(err)
		}
		appenders[appenderName] = appender
	}
}

func initRootLogger(se settings.Settings) {
	options, err := se.GetSettings("Root")
	if err == settings.Nil {
		return
	}
	overideLoggerFromSettings(logrus.StandardLogger(), options)
}

func initLoggers(se settings.Settings) {
	loggerMap, err := se.GetMap("Loggers")
	if err == settings.Nil {
		return
	}
	for name := range loggerMap {
		loggerName := name.(string)
		options, err := se.GetSettings(fmt.Sprintf("Loggers.%s", loggerName))
		if err != nil {
			logrus.Panic(err)
		}
		logger := createLoggerInheritRoot()
		overideLoggerFromSettings(logger, options)
		loggers[loggerName] = logger
	}
}

func createLoggerInheritRoot() *logrus.Logger {
	logger := logrus.New()
	logger.Out = logrus.StandardLogger().Out
	logger.Hooks = logrus.StandardLogger().Hooks
	logger.Formatter = logrus.StandardLogger().Formatter
	logger.Level = logrus.StandardLogger().Level
	return logger
}

func overideLoggerFromSettings(logger *logrus.Logger, se settings.Settings) {
	appenderName, err := se.GetString("Appender")
	if err != nil && err != settings.Nil {
		logrus.Panic(err)
	}
	if err != settings.Nil {
		appender := appenders[appenderName]
		if appender == nil {
			logrus.Panicf("Unknown appender name %s", appenderName)
		}
		logger.Out = appender
	}

	hookNames, err := se.GetSlice("Hooks")
	if err != nil && err != settings.Nil {
		logrus.Panic(err)
	}
	if err != settings.Nil {
		levelHooks := logrus.LevelHooks{}
		for _, item := range hookNames {
			hookName := item.(string)
			hook := hooks[hookName]
			if hook == nil {
				logrus.Panicf("Unknown hook name %s", hookName)
			}
			levelHooks.Add(hook)
		}
		logger.Hooks = levelHooks
	}

	formatterName, err := se.GetString("Formatter")
	if err != nil && err != settings.Nil {
		logrus.Panic(err)
	}
	if err != settings.Nil {
		formatter := formatters[formatterName]
		if formatter == nil {
			logrus.Panicf("Unknown formatter name %s", formatterName)
		}
		logger.Formatter = formatter
	}

	levelName, err := se.GetString("Level")
	if err != nil && err != settings.Nil {
		logrus.Panic(err)
	}
	if err != settings.Nil {
		level := parseLevelAndCheck(levelName)
		logger.Level = level
	}
}

func parseLevelAndCheck(levelName string) logrus.Level {
	level, err := logrus.ParseLevel(levelName)
	if err != nil {
		logrus.Panic(err)
	}
	// if level == logrus.PanicLevel || level == logrus.FatalLevel {
	// logrus.Panicf("Unsupported level %s", levelName)
	// }
	return level
}

func RegisterHook(id string, newFunc NewHookFunc) {
	if newHookFuncRegistry[id] != nil {
		logrus.Panicf("duplicated hook id %s", id)
	}
	newHookFuncRegistry[id] = newFunc
}

func RegisgerFormatter(id string, newFunc NewFormatterFunc) {
	if newFormatterFuncRegistry[id] != nil {
		logrus.Panicf("duplicated formatter id %s", id)
	}
	newFormatterFuncRegistry[id] = newFunc
}

func RegisgerAppenders(id string, newFunc NewAppenderFunc) {
	if newAppenderFuncRegistry[id] != nil {
		logrus.Panicf("duplicated appender id %s", id)
	}
	newAppenderFuncRegistry[id] = newFunc
}

func GetLogger(loggerName string) *logrus.Logger {
	logger := loggers[loggerName]
	if logger != nil {
		return logger
	}
	return logrus.StandardLogger()
}
