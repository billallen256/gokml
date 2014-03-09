package kml

import (
	"fmt"
	"testing"
)

func TestKML(t *testing.T) {
	k := NewKML()
	f := NewFolder("Test Folder", "This is a test folder")
	k.AddFolder(f)

	s := NewStyle("TestStyle", 240, 0, 255, 0)
	s.SetIconURL("http://maps.google.com/mapfiles/kml/paddle/wht-circle.png")
	f.AddFeature(s)

	p, err := NewPoint(40.67, -73.9, 0.0)

	if err != nil {
		t.Errorf("%s", err)
	}

	pm := NewPlacemark("Manhattan", "The Big Apple", p)
	pm.SetStyle("TestStyle")
	f.AddFeature(pm)

	fmt.Printf("%s", k.Render())
}
