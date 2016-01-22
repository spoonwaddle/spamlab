#!/bin/bash

set -e

clear

SOURCE_DIR=/vagrant/spam_classifier
TRAINING_DIR=/tmp/training_data
REDIS_PORT=5656
TEST_REDIS=127.0.0.1:$REDIS_PORT

prepare_redis () {
    echo "Creating test redis instance at $TEST_REDIS"
    redis-server --daemonize yes --port $REDIS_PORT
    if !  ( redis-cli -p $REDIS_PORT "INFO" | grep -q "loading:0" && redis-cli -p $REDIS_PORT "INFO" | grep -q "tcp_port:$REDIS_PORT"  ) ; then
        # make sure redis is up and running
        echo "Could not start test Redis instance on localhost:$REDIS_PORT ! exiting..."
        exit 1
    fi
    echo "Clearing..."
    redis-cli "FLUSHALL"  # clean redis out for tests
}

generate_training_data () {
    rm -rf $TRAINING_DIR && mkdir $TRAINING_DIR 
    echo "Generating training data in $TRAINING_DIR..."
    
    for i in `seq 1 10`; do
        echo "cat dog cat dog cat dog" > $TRAINING_DIR/doc$i.spam 
        echo "fish bird fish bird fish bird" > $TRAINING_DIR/doc$i.ham
    done
}

build_classifier () {
    echo "Building spam_classifier from source..."
    cd $SOURCE_DIR
    go build
}

train_classifier () {
    echo "Training classifier on data in $TRAINING_DIR..."
    $SOURCE_DIR/spam_classifier train -redis=$TEST_REDIS -spam=$TRAINING_DIR/*.spam -ham=$TRAINING_DIR/*.ham
}

test_classifier () {
    echo "Testing spam sample..."
    SPAM_RESULT=`echo "cat dog dog dog cat" | $SOURCE_DIR/spam_classifier classify -redis=$TEST_REDIS`
    echo "Result: $SPAM_RESULT -- Expected: SPAM"
    if [ "$SPAM_RESULT" = "SPAM" ] ; then
        echo "GREAT SUCCESS!!!"
    else
        echo "Something has gone horribly wrong!"
    fi
    
    echo "Testing ham sample..."
    HAM_RESULT=`echo "fish fish fish bird" | $SOURCE_DIR/spam_classifier classify -redis=$TEST_REDIS`
    echo "Result: $HAM_RESULT -- Expected: HAM"
    if [ "$HAM_RESULT" = "HAM" ] ;
    then
        echo "GREAT SUCCESS!!!"
    else
        echo "Something has gone horribly wrong!"
    fi
}

cleanup () {
    echo "Deleting training data..."
    rm -rf $TRAINING_DIR
    echo "Stopping test redis instance..."
    redis-cli -p $REDIS_PORT shutdown
    
}

prepare_redis
generate_training_data
build_classifier
train_classifier
test_classifier
cleanup
