package main

import (
	"fmt"

	config "gitlab.cept.gov.in/it-2.0-common/api-config"
)

type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Output string `yaml:"output"`
}

func main() {
	//var cfg, _ = config.NewDefaultConfigFactory().Create() // this will create a config object with default values
	var cfg, _ = config.NewDefaultConfigFactory().Create(
		config.WithFileName("config"),                      // config files base name
		config.WithFilePaths(".", "./config", "./configs"), // config files lookup paths. As per example config file can be in root of the project or in config or configs directory
		config.WithAppEnv("test"),                          // environment to load config for
	)

	//Default Functions relating app details can be can be accessed as follows:
	fmt.Printf("name: %s\n ", cfg.AppName())           // myapp
	fmt.Printf("env: %s\n", cfg.AppEnv())              // development
	fmt.Printf("App Version: %s\n ", cfg.AppVersion()) // 1.0.0
	fmt.Printf("debug: %t\n", cfg.AppDebug())          // false

	// Others can be accessed as follows:
	//Accessing values from the config object
	// You can access values using tree structure of config yaml file
	// For example: cron.scheduler.concurrency.limit.max gives you the value of max key under concurrency under scheduler under cron
	fmt.Printf("db.host: %s\n", cfg.GetString("db.host"))                                            // localhost
	fmt.Printf("cron.scheduler.limit.max: %d\n", cfg.GetInt("cron.scheduler.concurrency.limit.max")) // 3

	// You can also access values using the Of method which returns a subroot of the config object
	// This is useful when you want to access a specific section of the config file
	// For example: cfg.Of("server") will return a subroot of the config object with server as the root
	serverCfg, e := cfg.Of("server")
	if e != nil {
		fmt.Println("Error:", e)
	}
	fmt.Println("Body Limit:", serverCfg.GetInt("bodylimit"))
	fmt.Println("Read Buffer Size:", serverCfg.GetInt("readbuffersize"))

	db, e := cfg.Of("db")
	if e != nil {
		fmt.Println("DB config not found")
	}

	if db != nil {
		fmt.Printf("DB Host: %s\n", db.GetString("host"))            // localhost
		fmt.Printf("DB Username: %s\n", db.GetString("username"))    //postgres
		fmt.Println("DB Username Exists:", db.Exists("username"))    //Checks whether key Username under db subroot exists in the config file
		fmt.Printf("DB schema exists: %s\n", db.GetString("schema")) //public
		fmt.Printf("DB Port: %d\n", db.GetInt("port"))               //5432
	} else {
		fmt.Println("DB config not found")
	}

	// You can also convert the config object to a struct
	// This is useful when you want to convert the config object to a struct
	// For example: config.ToStruct(cfg, "log", &LogConfig) will convert the log section of the config object to a struct
	var logConfig LogConfig
	e = config.ToStruct(cfg, "afdsf", &logConfig)
	if e != nil {
		fmt.Println("Error:", e)
	}
	fmt.Println("log level:", logConfig.Level)
	fmt.Println("log format:", logConfig.Format)

	// You can also access environment variables
	fmt.Println(" App Env Variable:", cfg.GetEnvVar("APP_ENV")) // test

	//Environment variables can override the values in the config file
	// For example: If you set DB_HOST=localhost in the environment variables, cfg.GetString("db.host") will return localhost instead of the value in the config file
	// This is useful for getting values set in kubernets secrets for  production environment
	// During deployment Devops team pulls values from kuberntes secrets and sets them as environment variables
	// All values in config can be overridden in above fashion by replacing . with _ and converting to uppercase
	// For example: db.host can be accessed as DB_HOST
	// Priority of values is as follows: Environment variables > config<APP_ENV>.yaml values > Config.yaml  values > Default values

}
