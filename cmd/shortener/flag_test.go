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

	os.Args = []string{"cmd", "-a=" + aValue, "-b=" + bValue}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flag.CommandLine = flag.NewFlagSet(os.Args[1], flag.ExitOnError)

	var a string
	flag.StringVar(&a, "a", "", "")

	var b string
	flag.StringVar(&b, "b", "", "")

	flag.Parse()

	if a != aValue {
		t.Errorf("expected %v, got %v", aValue, a)
	}

	if b != bValue {
		t.Errorf("expected %v, got %v", bValue, b)
	}
}
