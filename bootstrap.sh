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

download_spam_dataset ()
{
    DATADIR=/home/vagrant/training_data
    # TODO FIX THESE PERMISSIONS
    mkdir -p $DATADIR/ham
    mkdir $DATADIR/spam
    for i in {1..6}; do
        wget -q http://www.aueb.gr/users/ion/data/enron-spam/preprocessed/enron$i.tar.gz
        tar -zxvf enron$i.tar.gz
        mv enron$i/spam/* $DATADIR/spam
        mv enron$i/ham/* $DATADIR/ham
        rm -rf enron$i.tar.gz enron$i
    find $DATADIR -type d -exec chmod 755 {} \;
    find $DATADIR -type f -exec chmod 644 {} \;
    chown -R vagrant:vagrant $DATADIR
    done
}

install_git
install_golang_and_go_deps
install_redis
download_spam_dataset
