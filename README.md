# Waft

<p align="center">
  <img height="150" src="./logo.png"  alt="Waft" title="Waft">
</p>

[![Build Status](https://app.travis-ci.com/x-debug/waft.svg?branch=master)](https://app.travis-ci.com/x-debug/waft)
[![Go Report Card](https://goreportcard.com/badge/github.com/x-debug/waft)](https://goreportcard.com/report/github.com/x-debug/waft)
[![codecov](https://codecov.io/gh/x-debug/waft/branch/master/graph/badge.svg?token=IHVP92FLDV)](https://codecov.io/gh/x-debug/waft)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> An API Gateway written in Go

This is a lightweight API Gateway and Management Platform, It's **Developing** currently.

## Document
**Coming soon...**

## QuickStart
```
make build
./build/test_backend -port 1111 //open backend on port 1111 for loadbalancer 
./build/test_backend -port 2222 //open backend on port 2222 for loadbalancer
./build/waft start //start proxy server
./build/waft stop //stop proxy server
./build/waft restart //also, restart server easily
```

open web browser http://localhost or your domain config in waft.yml

## About Daemon
Waft not support run on background, but you can use nohup or supervisord.

## Key Features
* [ ] Hot reload configuration
* [x] Multi-Service
* [x] Load Balance
* [ ] Service Discovery
* [x] Plugin
* [ ] IpWhiteList
* [ ] JWT Authorization
* [ ] API Retry After Failure
* [ ] Backend Server Health Check
* [ ] Management of API
* [ ] Circuit Breaker
* [ ] RateLimit