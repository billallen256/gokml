package gokml

import (
	"fmt"
	"testing"
	"time"
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

	states := NewStyle("StateStyle", 240, 0, 0, 255)
	f.AddFeature(states)

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

	colorado := NewPolygon()
	colorado.AddPoint(NewPoint(41.071904, -101.868843, 0.0))
	colorado.AddPoint(NewPoint(36.926393, -101.868843, 0.0))
	colorado.AddPoint(NewPoint(36.926393, -109.279635, 0.0))
	colorado.AddPoint(NewPoint(41.071904, -109.279635, 0.0))
	pm = NewPlacemark("Colorado", "The Centennial State", colorado)
	pm.SetStyle("StateStyle")
	pm.SetTime(time.Now().Add(-10*time.Hour), time.Now())
	f.AddFeature(pm)

	fmt.Printf("%s", k.Render())
}
