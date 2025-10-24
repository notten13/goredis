package main

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func stringToReader(s string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(s))
}

func TestEmptyInput(t *testing.T) {
	_, err := parseRESP(stringToReader(""))
	if err == nil || err.Error() != "unexpected end of input" {
		t.Fail()
	}
}

func TestUnsupportedType(t *testing.T) {
	// ignoring anything that is not an array or bulk string for now
	res, err := parseRESP(stringToReader("+OK\r\n"))
	if err != nil || len(res) != 0 {
		t.Fail()
	}
}

func TestBulkStringPing(t *testing.T) {
	res, _ := parseRESP(stringToReader("$4\r\nPING\r\n"))
	if !reflect.DeepEqual(res, []string{"PING"}) {
		t.Fail()
	}
}

func TestBulkStringGet(t *testing.T) {
	res, _ := parseRESP(stringToReader("$3\r\nGET\r\n"))
	if !reflect.DeepEqual(res, []string{"GET"}) {
		t.Fail()
	}
}

func TestEmptyBulkString(t *testing.T) {
	res, _ := parseRESP(stringToReader("$0\r\n\r\n"))
	if !reflect.DeepEqual(res, []string{""}) {
		t.Fail()
	}
}

func TestInvalidBulkStringLengthLetter(t *testing.T) {
	_, err := parseRESP(stringToReader("$x\r\nPING\r\n"))
	if err == nil || err.Error() != "invalid bulk string length" {
		t.Fail()
	}
}

func TestInvalidBulkStringLengthNegative(t *testing.T) {
	_, err := parseRESP(stringToReader("$-1\r\nPING\r\n"))
	if err == nil || err.Error() != "invalid bulk string length" {
		t.Fail()
	}
}

func TestWrongBulkStringLength(t *testing.T) {
	_, err := parseRESP(stringToReader("$5\r\nPING\r\n"))
	if err == nil || err.Error() != "bulk string length mismatch" {
		t.Fail()
	}
}

func TestMissingBulkStringData(t *testing.T) {
	_, err := parseRESP(stringToReader("$4\r\n"))
	if err == nil || err.Error() != "unexpected end of input" {
		t.Fail()
	}
}

func TestInvalidArrayLengthLetter(t *testing.T) {
	_, err := parseRESP(stringToReader("*x\r\n$4\r\nPING\r\n"))
	if err == nil || err.Error() != "invalid array length" {
		t.Fail()
	}
}

func TestInvalidArrayLengthNegative(t *testing.T) {
	_, err := parseRESP(stringToReader("*-1\r\n$4\r\nPING\r\n"))
	if err == nil || err.Error() != "invalid array length" {
		t.Fail()
	}
}

func TestWrongArrayLengthTooShort(t *testing.T) {
	_, err := parseRESP(stringToReader("*2\r\n$4\r\nPING\r\n"))
	if err == nil || err.Error() != "unexpected end of input" {
		t.Fail()
	}
}

func TestWrongArrayLengthTooLong(t *testing.T) {
	res, err := parseRESP(stringToReader("*1\r\n$4\r\nPING\r\n$3\r\nGET\r\n"))
	// should not error, just ignore extra data
	if err != nil || !reflect.DeepEqual(res, []string{"PING"}) {
		t.Fail()
	}
}

func TestEmptyArray(t *testing.T) {
	res, _ := parseRESP(stringToReader("*0\r\n"))
	if !reflect.DeepEqual(res, []string{}) {
		t.Fail()
	}
}

func TestArrayOfOneBulkString(t *testing.T) {
	res, _ := parseRESP(stringToReader("*1\r\n$4\r\nPING\r\n"))
	if !reflect.DeepEqual(res, []string{"PING"}) {
		t.Fail()
	}
}

func TestArrayOfTwoBulkStrings(t *testing.T) {
	res, _ := parseRESP(stringToReader("*2\r\n$3\r\nGET\r\n$4\r\nkey1\r\n"))
	if !reflect.DeepEqual(res, []string{"GET", "key1"}) {
		t.Fail()
	}
}

func TestArrayOfThreeBulkStrings(t *testing.T) {
	res, _ := parseRESP(stringToReader("*3\r\n$3\r\nSET\r\n$4\r\nkey1\r\n$4\r\nval1\r\n"))
	if !reflect.DeepEqual(res, []string{"SET", "key1", "val1"}) {
		t.Fail()
	}
}

func TestArrayWithEmptyBulkString(t *testing.T) {
	res, _ := parseRESP(stringToReader("*3\r\n$3\r\nSET\r\n$4\r\nkey1\r\n$0\r\n\r\n"))
	if !reflect.DeepEqual(res, []string{"SET", "key1", ""}) {
		t.Fail()
	}
}
