package kml

import (
	"fmt"
	"testing"
)

func TestKML(t *testing.T) {
	k := NewKML()
	f := NewFolder("Test Folder", "This is a test folder")
	k.AddFolder(f)

	places := NewStyle("PlaceStyle", 240, 0, 255, 0)
	places.SetIconURL("http://maps.google.com/mapfiles/kml/paddle/wht-circle.png")
	f.AddFeature(places)

	flights := NewStyle("FlightStyle", 240, 255, 0, 0)
	f.AddFeature(flights)

	manhattan := NewPoint(40.67, -73.9, 0.0)
	pm := NewPlacemark("Manhattan", "The Big Apple", manhattan)
	pm.SetStyle("PlaceStyle")
	f.AddFeature(pm)

	london := NewPoint(51.51, 0.1275, 0.0)
	pm = NewPlacemark("London", "The City", london)
	pm.SetStyle("PlaceStyle")
	f.AddFeature(pm)

	paris := NewPoint(48.85, 2.35, 0.0)
	pm = NewPlacemark("Paris", "The City of Light", paris)
	pm.SetStyle("PlaceStyle")
	f.AddFeature(pm)

	tokyo := NewPoint(35.69, 139.7, 0.0)
	pm = NewPlacemark("Tokyo", "東京", tokyo)
	pm.SetStyle("PlaceStyle")
	f.AddFeature(pm)

	flightPath := NewLineString()
	flightPath.AddPoint(manhattan)
	flightPath.AddPoint(london)
	flightPath.AddPoint(paris)
	flightPath.AddPoint(tokyo)
	pm = NewPlacemark("Flight Path", "", flightPath)
	pm.SetStyle("FlightStyle")
	f.AddFeature(pm)

	fmt.Printf("%s", k.Render())
}
