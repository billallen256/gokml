package kml

import (
	"fmt"
	"testing"
)

func TestKML(t *testing.T) {
	k := NewKML("Test")
	f := NewFolder("Test Folder", "This is a test folder")
	k.AddFolder(f)
	p := NewPoint(40.67, 73.9, 0.0)
	pm := NewPlacemark("Manhattan", "The Big Apple", p)
	f.AddObject(pm)
	fmt.Printf("%s", k.Render())
}
