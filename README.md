# DSdriver
Distributed Systems driver program written in Golang

## Purpose
This package is build for development and testing of distributed system protocols in Golang.
The designer of a distributed system should not have to concern themself with messy network code during development.
This package solves the problem of making sure the protocol code is correct and tested.
Once the development is done, this package can be replaced with one that puts bits on the wire and the protocol should work as intended! (provided that networking code is correct and fault-tolerant :) )


