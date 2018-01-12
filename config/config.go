package config

import "github.com/astaxie/beego/config"

type DBConfig struct {
	DBUser string
	DBPass string
	DBName string
}

func (dc *DBConfig) Read() error {
	fullConfigIni, err := config.NewConfig("ini", "config.ini")
	if err != nil {
		return err
	}

	configIni, err := fullConfigIni.GetSection("default")

	if err != nil {
		return err
	}

	dc.DBUser = configIni["user"]
	dc.DBPass = configIni["pass"]
	dc.DBName = configIni["name"]
	return nil
}
