package cfg

import (
	"github.com/madlabx/pkgx/errors"
	"github.com/madlabx/pkgx/httpx"
	"github.com/madlabx/pkgx/log"
	"github.com/madlabx/pkgx/utils"
	"github.com/madlabx/pkgx/viperx"
	"github.com/spf13/viper"
)

var instance Config

func Parse(envPrefix string, cfgFile string, opts ...viper.DecoderConfigOption) error {
	if cfgFile == "" {
		return errors.New("Empty config file")
	}
	return viperx.ParseConfig(&instance, envPrefix, cfgFile, opts...)
}

func Get() *Config {
	return &instance
}

type Config struct {
	Sys       SysConfig
	MainLog   LogConfig
	AccessLog httpx.LogConfig
}

func (c *Config) String() string {
	return utils.ToString(c)
}

type SysConfig struct {
	ConfigFile        string `vx_name:"conf" vx_short:"c"`
	Port              string `vx_name:"port" vx_short:"p" vx_default:"8080" vx_desc:"port to listen on"`
	Address           string `vx_name:"address" vx_short:"a" vx_default:"127.0.0.1" vx_desc:"address to listen on"`
	ProfilePort       int    `vx_default:"16060"`
	MaxUploadParallel int    `vx_default:"100"`
	Domain            string `vx_must:"true"`
	Root              string `vx_must:"true"`

	DiskReservePercent string `vx_default:"-"` //   -: 遵从磁盘预留空间设置； 整数：预留百分比

}

type LogConfig struct {
	LogFile log.FileConfig
	Level   string `vx_default:"info"`
}
