package main

import "testing"

func TestRun(t *testing.T) {
	if err := run("https://www.mercadolivre.com.br/"); err != nil {
		t.Fatalf("run() error = %v, want nil", err)
	}
}

func TestRunRejectsInvalidURL(t *testing.T) {
	if err := run("not a url"); err == nil {
		t.Fatal("run() error = nil, want an error")
	}
}
