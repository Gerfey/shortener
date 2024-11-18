package main

import (
	"flag"
	"os"
	"testing"
)

func TestFlagParsing(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	var aValue = "localhost:8081"
	var bValue = "http://localhost:8082"
	var dValue = "host=localhost port=5432 user=shortener password=shortener dbname=shortener sslmode=disable"

	os.Args = []string{"cmd", "-a=" + aValue, "-b=" + bValue, "-d=" + dValue}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flag.CommandLine = flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	flag.CommandLine = flag.NewFlagSet(os.Args[2], flag.ExitOnError)

	var a string
	flag.StringVar(&a, "a", "", "")

	var b string
	flag.StringVar(&b, "b", "", "")

	var d string
	flag.StringVar(&d, "d", "", "")

	flag.Parse()

	if a != aValue {
		t.Errorf("expected %v, got %v", aValue, a)
	}

	if b != bValue {
		t.Errorf("expected %v, got %v", bValue, b)
	}

	if d != dValue {
		t.Errorf("expected %v, got %v", dValue, d)
	}
}
