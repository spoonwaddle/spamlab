package main

import (
	"github.com/rafaeljusto/redigomock"
	"testing"
)

var spamCount interface{} = "spam:#texts"
var hamCount interface{} = "ham:#texts"
var totalCount interface{} = "#texts"
var spamHelloCount interface{} = "spam:hello"
var hamHelloCount interface{} = "ham:hello"
var spamWordCount interface{} = "spam:#words"
var hamWordCount interface{} = "ham:#words"

var pLabelResp = []interface{}{
	[]uint8{uint8('1')}, // # texts labeled as SPAM
	[]uint8{uint8('4')}, // # texts total
}
var spamWordCountResp = []interface{}{
	[]uint8{uint8('8')}, // # words in label
	[]uint8{uint8('2')}, // # times 'hello' appears
}
var hamWordCountResp = []interface{}{
	[]uint8{uint8('9')}, // # words in label
	[]uint8{uint8('8')}, // # times 'hello' appears
}
var text = []string{"hello"}

func mockRedisProbDist(conn *redigomock.Conn) RedisWordProbabilityDist {
	return RedisWordProbabilityDist{
		conn:            conn,
		label2namespace: map[SpamLabel]string{SPAM: "spam:", HAM: "ham:"},
		wordCountSuffix: "#words",
		textCountSuffix: "#texts",
	}
}

func TestProbabilityOfLabel(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("MGET", spamCount, totalCount).Expect(pLabelResp)
	dist := mockRedisProbDist(conn)
	result, _ := dist.ProbabilityOfLabel(SPAM)
	expected := .25
	if expected != result {
		t.Error("Expected", expected, "got", result)
	}
}

func TestProbabilityOfTextGivenLabel(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("MGET", spamWordCount, spamHelloCount).Expect(spamWordCountResp)
	dist := mockRedisProbDist(conn)
	result, _ := dist.ProbabilityOfTextGivenLabel(text, SPAM)
	expected := .25
	if expected != result {
		t.Error("Expected", expected, "got", result)
	}
}

func TestProbabilityOfTextAndLabel(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("MGET", spamCount, totalCount).Expect(pLabelResp)
	conn.Command("MGET", spamWordCount, spamHelloCount).Expect(spamWordCountResp)
	dist := mockRedisProbDist(conn)
	result, _ := dist.ProbabilityOfTextAndLabel(text, SPAM)
	expected := .25 * .25
	if expected != result {
		t.Error("Expected", expected, "got", result)
	}
}

func TestMostLikelyLabel(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("MGET", spamCount, totalCount).Expect(pLabelResp)
	conn.Command("MGET", spamCount, totalCount).Expect(pLabelResp)
	conn.Command("MGET", spamWordCount, spamHelloCount).Expect(spamWordCountResp)
	conn.Command("MGET", hamWordCount, hamHelloCount).Expect(hamWordCountResp)
	dist := mockRedisProbDist(conn)
	result, _ := dist.MostLikelyLabel(text)
	expected := HAM
	if expected != result {
		t.Error("Expected", expected, "got", result)
	}
}
