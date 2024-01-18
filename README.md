# goadmin

An admin panel developed with gin, gorm, jwt, casbin, mysql and vuejs2 with vue-element-admin

## Features

- `Gin` an API framework similar to martini but with better performance.
- `MySQL` uses the MySql database
- `Jwt` uses JWT lightweight authentication and provides active user Token refresh function
- `Casbin` Casbin is a powerful and efficient open source access control framework.
- `Gorm` using Gorm 2.0 version.
- `Validator` using validator v10.
- `Lumberjack` sets log file size, save quantity, save time and compression, etc.
- `Viper` for configuration solution
- `GoFunk` toolkit containing a large number of Slice operation methods

## middleware

- `AuthMiddleware` handles login, logout, and stateless token verification
- `RateLimitMiddleware` limits the number of user requests
- `OperationLogMiddleware` records all user operations
- `CORSMiddleware` solve cross-domain request problems
- `CasbinMiddleware` uses Casbin to control user access

## Todo

- Add image server
- Add promtail-loki-grafana log monitoring system
- Add swagger documentation

