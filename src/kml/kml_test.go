package kml

import (
	"fmt"
	"testing"
)

func TestKML(t *testing.T) {
	k := NewKML()
	f := NewFolder("Test Folder", "This is a test folder")
	k.AddFolder(f)
	p, err := NewPoint(40.67, -73.9, 0.0)

	if err != nil {
		t.Errorf("%s", err)
	}

	pm := NewPlacemark("Manhattan", "The Big Apple", p)
	f.AddFeature(pm)
	fmt.Printf("%s", k.Render())
}
