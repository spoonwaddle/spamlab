#!/bin/bash

export GOPATH=/vagrant/spam_classifier
VAGRANT_PROFILE=/home/vagrant/.profile

install_git ()
{
    add-apt-repository ppa:git-core/ppa
    apt-get update
    apt-get install -y git
}

install_golang_and_go_deps ()
{
    apt-get install -y golang
    echo "export GOPATH=$GOPATH" >> $VAGRANT_PROFILE
    go get github.com/garyburd/redigo/redis
    go get github.com/rafaeljusto/redigomock
}

install_redis ()
{
    apt-get install -y redis-server
    echo "export REDIS_URL=127.0.0.1:6379" >> $VAGRANT_PROFILE
}

install_git
install_golang_and_go_deps
install_redis
