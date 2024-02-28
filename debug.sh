#!/bin/sh
dlv debug --headless -l 127.0.0.1:4009 cmd/zwooc/main.go -- watch dev