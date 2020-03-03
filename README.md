# Go Dotenv

![Go Unit Tests](https://github.com/navaz-alani/dotenv/workflows/Go%20Unit%20Tests/badge.svg)
[![GoDoc](https://godoc.org/github.com/navaz-alani/dotenv?status.svg)](https://godoc.org/github.com/navaz-alani/dotenv)
[![CodeFactor](https://www.codefactor.io/repository/github/navaz-alani/dotenv/badge)](https://www.codefactor.io/repository/github/navaz-alani/dotenv)

This package provides dotenv functionality (similar to that in NodeJS) for Golang projects.

It has been decoupled from the operating system in order to easily provide more flexibility
and additional features which can be used to gain greater control of the application's 
runtime parameters. This includes asserting that a set of required parameters have been 
initialized and chaining together environment variable source files using a special 
`__GO_LOAD` key and more. 

Please find the API documentation by clicking on the `godoc` tag above.
