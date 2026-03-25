#!/bin/bash
set -e

protoc --go_out=. --go_opt=paths=source_relative msg.proto
