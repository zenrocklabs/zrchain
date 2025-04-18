//go:build integration
// +build integration

package integration

import (
	"fmt"
	"testing"
)

func TestSomethingIntegration(t *testing.T) {
	t.Log("Running dummy integration test...")
	if 1+1 != 2 {
		fmt.Println("error")
		t.Fatal("Error here")
	}
}
