# Portfolio CMS Server

![GitHub last commit](https://img.shields.io/github/last-commit/StGrozdanov/portfolio-cms-server)

A custom CMS backend for my portfolio website built with Golang. You can find the [Frontend here](https://github.com/StGrozdanov/portfolio-cms-client)

## Project Overview

### Technologies Used

- :whale: Docker
- :cloud: AWS SDK 
- :lock: Bcrypt 
- 🍹Gin 
- :key: JSON Web Token 
- 🏬 PostgreSQL 
- 📅 Moment 
- 👨‍💻 Logrus
- 🔦 Prometheus
- 🪝 pre-commit and pre-push hooks

### Database

The server uses PostgreSQL database with a mixture of JSONB and standard columns.

## Features

- CI/CD pipeline consisting of custom linters, unit tests, integration tests and a single deployment environment
- S3 bucket files download, upload, delete and read
- Authentication with Bcrypt and JWT
- Logging with Logrus
- Tracking by IP, Country, Browser, Device type, referer
- Analytics by date, week, month, quarter, year
- Metrics by prometheus
- Healths endpoint
- Authenticated and non authenticated paths
- Response time measuring on each endpoint
- A full range of CRUD operations for each portfolio resource
