# API Config

**API Config** is a Go library that provides a simple and flexible way to manage application configurations using YAML files. This library allows for structured configuration access and supports environment variable overrides for production deployments.

## Table of Contents

- [Getting Started](#getting-started)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration and Features](#configuration)

## Getting Started

To get started with **API Config**, follow the steps below to install and use the library in your Go application.

### Installation

1. **Get the repository:**

    ```bash
    go get https://gitlab.cept.gov.in/it-2.0-common/api-config
    ```

2. **Set up your configuration:**

   Create a `config.yaml` file in your project root. Here’s an example configuration file:

   Example yaml files can be found in [configs](https://gitlab.cept.gov.in/it-2.0-common/api-examples/-/tree/main/config?ref_type=heads) folder 

  ```yaml
    db:
      username: "postgres"
      password: "yourpassword"
      host: "localhost"
      port: "5432"
      database: "yourdatabase"
      schema: "public"
    
    server:
      bodylimit: 10485
      readbuffersize: 16384
      addr: 8080
      readtimeout: 10s
      writetimeout: 10s
      timeout: 40s
    
    log:
      level: "debug"
      format: "json"
      output: "stdout"
  ```

### Usage

1. **Basic Setup:**

   Import the necessary packages in your `main.go`:

    ```go
    package main

    import (
        "fmt"
        config "gitlab.cept.gov.in/it-2.0-common/api-config"
    )
    ```

2. **Accessing Configuration:**

   Use the following example to access various configuration parameters:

   By Default, the module expects configuration files:
   1. To be present in .(root) , ./configs directories of your project
   2. To offer env overrides files named config.{env}.{format} based on the env var APP_ENV (ex: config.test.yaml if env     var       APP_ENV=test). In production, APP_ENV  will be pushed as kubernetes secrets/vault secrets during deployment phase by DevOps.
   Based on env value application picks config.{env}.{format}. 
   For example, if APP_ENV is test, application checks for config.test.yaml file load values.


   For a Default folder structure 

  ```go
   cfg, _ := config.NewDefaultConfigFactory().Create()
  ```

  To change the folder structure or name of the config following code can be used 

   ```go
    func main() {
        cfg, _ := config.NewDefaultConfigFactory().Create(
        config.WithFileName("config"), // Base name of the config file
        config.WithFilePaths(".", "./config", "./configs"), // Lookup paths
        config.WithAppEnv("test"), 
        )
   ```

   Access default functions related to application details
   ```go
    fmt.Printf("App Name: %s\n", cfg.AppName())           // Retrieves application name
    fmt.Printf("Environment: %s\n", cfg.AppEnv())          // Retrieves application environment (development/production)
    fmt.Printf("App Version: %s\n", cfg.AppVersion())     // Retrieves application version
    fmt.Printf("Debug Mode: %t\n", cfg.AppDebug())        // Retrieves whether debug mode is enabled
   ```

   Access values from the config object using a tree structure

   ```go
    fmt.Printf("DB Host: %s\n", cfg.GetString("db.host"))                                  // Retrieves DB host
    fmt.Printf("Max Concurrency: %d\n", cfg.GetInt("cron.scheduler.concurrency.limit.max")) // Retrieves max concurrency limit
   ```

  Access a specific section of the config object
   ```go
    serverCfg, e := cfg.Of("server") // Returns a subroot of the config object with "server" as the root
	if e != nil {
		fmt.Println("Error:", e)
	} 
    fmt.Println("Body Limit:", serverCfg.GetInt("bodylimit")) // Retrieves body limit from server config
    fmt.Println("Read Buffer Size:", serverCfg.GetInt("readbuffersize")) // Retrieves read buffer size
   ```

   Get DB config subroot and Access DB configuration values

  ```go
      	db, e := cfg.Of("db") // Get DB config subroot
	      if e != nil {
		    fmt.Println("DB config not found")
	      }// Check if DB config section exists

        // Access DB configuration values
        fmt.Printf("DB Host: %s\n", db.GetString("host"))         // Retrieves DB host
        fmt.Printf("DB Username: %s\n", db.GetString("username")) // Retrieves DB username
        fmt.Println("DB Username Exists:", db.Exists("username"))  // Checks if username key exists in DB config
        fmt.Printf("DB Schema: %s\n", db.GetString("schema"))     // Retrieves DB schema
        fmt.Printf("DB Port: %d\n", db.GetInt("port"))             // Retrieves DB port
  ```

  Convert a subroot to struct. For example, log section of the config object to a struct:

  ```go
  var logConfig LogConfig
  config.ToStruct(cfg, "log", &logConfig) // Converts the "log" section to LogConfig struct
  fmt.Println("Log Level:", logConfig.Level) // Retrieves log level from log config
  fmt.Println("Log Format:", logConfig.Format) // Retrieves log format from log config
  ```

  Access environment variables

  ```go
  fmt.Println("App Env Variable:", cfg.GetEnvVar("APP_ENV")) // Retrieves APP_ENV variable
  ```



### Configuration

The **API Config** library supports configuration through a YAML file, and environment variables can override any values specified in the file. 

- **Hierarchy**: Configuration values are accessed in a hierarchical manner, e.g., `db.host` can be accessed as `cfg.GetString("db.host")`.
- **Environment Variable Overrides**: For example, to set the database host as environment variable, add `DB_HOST=localhost` in OS environment variable to override `db.host` in your YAML file.
https://drive.google.com/file/d/1PEW7gYIpGxfoIcqbBQsvMggGlv3elnd2/view?usp=sharing
- **Nested Configuration Access**: This line retrieves the maximum concurrency limit from a nested structure in the configuration. The value can be accessed using the hierarchical path cron.scheduler.concurrency.limit.max, allowing you to drill down through the layers of the configuration. This approach facilitates organizing related configuration options in a logical structure, making it easier to manage complex configurations.
- **Priority of values**: Environment variables > config<APP_ENV>.yaml > Default values






