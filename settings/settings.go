package settings

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

var (
	logger  *log.Logger
	BaseDir string
)

func init() {
	BaseDir, _ = os.Getwd()
}

type wechat struct {
	CorpID         string
	AgentID        string
	Token          string
	EncodingAESKey string
	Secret         string
}

//Config for server
type Config struct {
	Wechat   *wechat `toml:"wechat"`
	Server   string
	LogLevel string
}

//GetConfig get config obj
func GetConfig(configPath string) (config Config, err error) {
	_, err = toml.DecodeFile(configPath, &config)
	return
}

//GetLogger get logger
func GetLogger(tomlconfig *Config) *log.Logger {
	if logger == nil {
		if tomlconfig == nil || tomlconfig.LogLevel == "" {
			logger = log.New(os.Stdout, "DEBUG", log.Ldate|log.Ltime)
		} else {
			logger = log.New(os.Stdout, tomlconfig.LogLevel, log.Ldate|log.Ltime)
		}
	}
	return logger

}
