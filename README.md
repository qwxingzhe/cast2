cast
====
[![GoDoc](https://godoc.org/github.com/qwxingzhe/cast2?status.svg)](https://godoc.org/github.com/qwxingzhe/cast2)
[![Build Status](https://github.com/qwxingzhe/cast2/actions/workflows/go.yml/badge.svg)](https://github.com/qwxingzhe/cast2/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/qwxingzhe/cast2)](https://goreportcard.com/report/github.com/qwxingzhe/cast2)

Easy and safe casting from one struct to another in Go

Don’t Panic! ... Cast

## What is Cast?

Cast is a library to convert between different go types in a consistent and easy way.

Cast provides simple functions to easily convert a struct to a another, an
struct into a map, etc. Cast does this intelligently when an obvious
conversion is possible. It doesn’t make any attempts to guess what you meant,
for example you can only convert a string to an int when it is a string
representation of an int such as “8”. Cast2 is used to supplement [Cast](https://github.com/spf13/cast)'s deficiency in aaa conversion

## Why use Cast?


## Usage


### Example ‘ToMap’:




### Example ‘ToMapByTag’:


