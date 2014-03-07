package kml

import (
	"fmt"
)

type Renderable interface {
	Render() string
}

type KML struct {
	Name    string
	Folders []*Folder
}

func NewKML(name string) *KML {
	f := make([]*Folder, 0, 2)
	return &KML{name, f}
}

func (k *KML) AddFolder(folder *Folder) {
	k.Folders = append(k.Folders, folder)
}

func (k *KML) Render() string {
	ret := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
		"<kml xmlns=\"http://www.opengis.net/kml/2.2\">\n"

	for _, folder := range k.Folders {
		ret += folder.Render()
	}

	ret += "</kml>\n"

	return ret
}

type Folder struct {
	Name        string
	Description string
	Objects     []Renderable
}

func NewFolder(name string, desc string) *Folder {
	o := make([]Renderable, 0, 10)
	return &Folder{name, desc, o}
}

func (f *Folder) AddObject(object Renderable) {
	f.Objects = append(f.Objects, object)
}

func (f *Folder) Render() string {
	ret := "<Folder>\n" +
		fmt.Sprintf("<name>%s</name>\n", f.Name) +
		fmt.Sprintf("<description>%s</description>\n", f.Description)

	for _, object := range f.Objects {
		ret += object.Render()
	}

	return ret
}

type Point struct {
	Lat float64
	Lon float64
	Alt float64
}

func NewPoint(lat float64, lon float64, alt float64) *Point {
	return &Point{lat, lon, alt}
}

func (p *Point) Render() string {
	ret := "<Point>\n" +
		"<extrude>0</extrude>\n" +
		"<altitudeMode>clampToGround</altitudeMode>\n" +
		fmt.Sprintf("<coordinates>%f,%f,%f</coordinates>\n", p.Lon, p.Lat, p.Alt) +
		"</Point>\n"

	return ret
}

type Placemark struct {
	Name        string
	Description string
	Geometry    Renderable
}

func NewPlacemark(name string, desc string, geom Renderable) *Placemark {
	return &Placemark{name, desc, geom}
}

func (pm *Placemark) Render() string {
	ret := "<Placemark>\n" +
		fmt.Sprintf("<name>%s</name>\n", pm.Name) +
		fmt.Sprintf("<description>%s</description>\n", pm.Description) +
		"<visibility>1</visiblity>\n" +
		pm.Geometry.Render() +
		"</Placemark>\n"

	return ret
}
