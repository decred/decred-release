package main

import "testing"

func TestRunning(t *testing.T) {
	c := &ctx{}
	r, err := c.isRunning("dcrd")
	if err != nil {
		t.Fatalf("%v", err)
	}

	t.Logf("dcrd running: %v\n", r)
}
