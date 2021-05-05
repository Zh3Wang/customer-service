package log

import (
	"customerservice/internal/pkg/setting"
	"fmt"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
	"time"
)

type MineFormatter struct{}

const TimeFormat = "2006-01-02 15:04:05"

var (
	Log        = logrus.New()
	Info       = Log.Info
	Debug      = Log.Debug
	Error      = Log.Error
	Fatal      = Log.Fatal
	Panic      = Log.Panic
	Warn       = Log.Warn
	WithFields = Log.WithFields
)

func InitLog() {
	//默认输出到终端
	Log.Out = os.Stdout

	//设置默认日志级别
	//Panic：记录日志，然后panic。
	//Fatal：致命错误，出现错误时程序无法正常运转。输出日志后，程序退出；
	//Error：错误日志，需要查看原因；
	//Warn：警告信息，提醒程序员注意；
	//Info：关键操作，核心流程的日志；
	//Debug：一般程序中输出的调试信息；
	//Trace：很细粒度的信息，一般用不到
	//日志级别从上向下依次增加，Trace最大，Panic最小。logrus有一个日志级别，高于这个级别的日志不会输出。
	var logLevel logrus.Level
	err := logLevel.UnmarshalText([]byte(setting.Config.LogConf.Level))
	if err != nil {
		Log.Panic("设置log级别失败：%v", err)
	}
	Log.SetLevel(logLevel)

	//Hook(Log, setting.Config.LogConf.LogPath, setting.Config.LogConf.LogSave)
}

//钩子方法，输出日志之前执行，这里用于确定输出的文件，按天分割日志
//为不同级别的日志设置不同的钩子方法
func Hook(log *logrus.Logger, logPath string, save int) {
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer(logPath, "debug", save),
		logrus.InfoLevel:  writer(logPath, "info", save),
		logrus.WarnLevel:  writer(logPath, "warn", save),
		logrus.ErrorLevel: writer(logPath, "error", save),
		logrus.FatalLevel: writer(logPath, "fatal", save),
		logrus.PanicLevel: writer(logPath, "panic", save),
	}, &MineFormatter{})
	log.AddHook(lfHook)
}

func writer(logPath string, level string, save int) *rotatelogs.RotateLogs {
	logFullPath := path.Join(logPath, level)
	var cstSh, _ = time.LoadLocation("Asia/Shanghai")
	fileSuffix := time.Now().In(cstSh).Format("2006-01-02") + ".log"
	writer, err := rotatelogs.New(
		logFullPath+"-"+fileSuffix,
		rotatelogs.WithLinkName(logFullPath), //生成软连接，指向最新的文件
		rotatelogs.WithRotationCount(save),   //文件最大保存份数
	)

	if err != nil {
		panic(err)
	}
	return writer
}

func (s *MineFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var data string
	for k, v := range entry.Data {
		if setting.Config.LogConf.HideKeys {
			_ = k
			data += fmt.Sprintf("[%v]", v)
		} else {
			data += fmt.Sprintf("[%s:%v]", k, v)
		}
	}
	msg := fmt.Sprintf("[%s] [%s] %s %s\n", time.Now().Local().Format(TimeFormat), strings.ToUpper(entry.Level.String()), data, entry.Message)
	return []byte(msg), nil
}
