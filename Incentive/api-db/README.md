# DB Library

**API DB** is a Go package that provides an abstraction layer for interacting with postgreSQL databases. It simplifies database connection management, transaction handling, and configuration using `pgxpool`.


## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [Examples](#examples)
- [Notes](#notes)

## Features

- **Connection Pooling**: Efficient resource management with `pgxpool`.
- **Transactional Operations**: Easy-to-use functions for read and write transactions.
- **Customizable Isolation Levels**: Support for different isolation levels for transactions.
- **Configurable Connection Settings**: Flexible options for database connection setup.
- **Graceful Cleanup**: Ensures database connections are closed properly.


## Installation

Add the `db` library to your Go project:

```bash
go get gitlab.cept.gov.in/it-2.0-common/n-api-db
```

Import the required package in your Go files:

```go
import (
dblib "gitlab.cept.gov.in/it-2.0-common/n-api-db"
)
```

## Configuration

The library uses a configuration struct (`DBConfig`) to set up database connection parameters. Below are the fields available in the configuration:

### DBConfig Fields

- `DBUsername`: Database username.
- `DBPassword`: Database password.
- `DBHost`: Host address of the database server.
- `DBPort`: Port of the database server.
- `DBDatabase`: Name of the database.
- `Schema`: Database schema.
- `MaxConns`: Maximum number of connections in the pool.
- `MinConns`: Minimum number of connections in the pool.
- `MaxConnLifetime`: Maximum connection lifetime in minutes.
- `MaxConnIdleTime`: Maximum idle time for connections in minutes.
- `HealthCheckPeriod`: Period for health checks in minutes.
- `AppName`: Name of the application for database logging purposes.

### Setting Up Configuration

```go
 dbConfig := dblib.DBConfig{ 
        DBUsername:        c.DBUsername(),
        DBPassword:        c.DBPassword(),
        DBHost:            c.DBHost(),
        DBPort:            c.DBPort(),
        DBDatabase:        c.DBDatabase(),
        Schema:            c.DBSchema(),
        MaxConns:          int32(c.MaxConns()),
        MinConns:          int32(c.MinConns()),
        MaxConnLifetime:   time.Duration(c.MaxConnLifetime()),   // In minutes
        MaxConnIdleTime:   time.Duration(c.MaxConnIdleTime()),   // In minutes
        HealthCheckPeriod: time.Duration(c.HealthCheckPeriod()), // In minutes
        AppName:           c.AppName(),
    }
```
## Usage

### Initialize Database Connection

Use the `DefaultDbFactory` to prepare the configuration and establish a connection:

```go
 preparedConfig := dblib.NewDefaultDbFactory().NewPreparedDBConfig(dbConfig)

    // Step 3: Establish the database connection
    dbConn, err := dblib.NewDefaultDbFactory().CreateConnection(preparedConfig)

    if err != nil {
        log.Warn(nil,"error in db connection %s", err)
        return nil, err
    }
    log.Info(nil,"Successfully connected to the database %s", c.DBConnection())
    defer dbConn.Close()

    if dbConn.Ping() == nil {
		log.info("Connection Established")
	} else {
		log.warn("Failed to establish database connection")
	}
```
## Examples

For complete example usage and implementation, please refer to the [API DB Example](https://gitlab.cept.gov.in/it-2.0-common/n-api-db/-/tree/main/example?ref_type=heads).

## Notes

- Ensure your PostgreSQL server is running and accessible with the provided configuration.
- Use appropriate isolation levels and access modes for transactions based on your application requirements.
- To use the functions from `util.go` or `helper.go`, prefix them with dblib. (assuming the import alias is used as `dblib`) to access those functions.
- If you encounter any issues or have questions about using the library, feel free to reach out via the following channels:

- **Issue Tracker**: [GitLab Issue Tracker](https://gitlab.cept.gov.in/it-2.0-common/n-api-db/issues)
- **Email**: [maintainer.email@example.com](mailto:soumyabrata.t@indiapost.gov.in)
 




