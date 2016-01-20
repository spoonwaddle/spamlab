#!/bin/bash

set -ex

testClassifier () {
    if !  ( redis-cli "INFO" | grep -q "loading:0" && redis-cli "INFO" | grep -q "tcp_port:6379"  ) ; then
        # make sure redis is up and running
        echo "redis server not running on localhost:6379 ! exiting..."
        exit 1
    fi
    redis-cli "FLUSHALL"  # clean redis out for tests
}

testClassifier
