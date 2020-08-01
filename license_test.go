package main

import (
	"github.com/google/licensecheck"
	"io/ioutil"
	"testing"
)

func TestLicense(t *testing.T) {
	b, err := ioutil.ReadFile("./LICENSE.md")
	if err != nil {
		t.Error(err)
	}
	coverage, ok := licensecheck.Cover(b, licensecheck.Options{
		MinLength: 10,
		Threshold: 40,
		Slop:      8,
	})
	if !ok {
		t.Errorf("No license")
	}
	if coverage.Percent < 97 {
		t.Errorf("Low license percent")
	}
	if len(coverage.Match) < 1 {
		t.Errorf("License match is empty")
	}
	if coverage.Match[0].Name != "MIT" {
		t.Errorf("License match not MIT")
	}
	if coverage.Match[0].Percent != 100 {
		t.Errorf("License match percent not 100")
	}

}
