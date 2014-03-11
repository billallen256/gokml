// The kml package is used for rendering Google Earth Keyhole Markup Language
// (KML) files.  Some of the terms used in this library are pulled from the
// KML specification (recommended reading/reference).
package kml

import (
	"fmt"
	"math"
	"strings"
	"sync"
)

type renderable interface {
	render() string
}

// KML represents the top-level KML document object.
type KML struct {
	folders []*Folder
	mutex   *sync.Mutex
}

// NewKML returns a pointer to a KML struct.
func NewKML() *KML {
	f := make([]*Folder, 0, 2)
	return &KML{f, new(sync.Mutex)}
}

// AddFolder adds a new Folder to the KML document.
func (k *KML) AddFolder(folder *Folder) {
	if folder != nil {
		k.mutex.Lock()
		k.folders = append(k.folders, folder)
		k.mutex.Unlock()
	}
}

// Renders the entire KML document.
func (k *KML) Render() string {
	ret := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
		"<kml xmlns=\"http://www.opengis.net/kml/2.2\">\n"

	for _, folder := range k.folders {
		ret += folder.render()
	}

	ret += "</kml>\n"

	return ret
}

// Folder represents a folder in the KML document.
type Folder struct {
	name        string
	description string
	features    []renderable
	mutex       *sync.Mutex
}

// Returns a pointer to a new Folder instance.
func NewFolder(name string, desc string) *Folder {
	f := make([]renderable, 0, 10)
	return &Folder{name, desc, f, new(sync.Mutex)}
}

// AddFeature adds a feature (Placemark, another Folder, etc.) to
// the Folder.
func (f *Folder) AddFeature(feature renderable) {
	if feature != nil {
		f.mutex.Lock()
		f.features = append(f.features, feature)
		f.mutex.Unlock()
	}
}

func (f *Folder) render() string {
	ret := "<Folder>\n" +
		fmt.Sprintf("<name>%s</name>\n", f.name) +
		fmt.Sprintf("<description>%s</description>\n", f.description)

	for _, feature := range f.features {
		ret += feature.render()
	}

	ret += "</Folder>\n"

	return ret
}

// Style represents a style used for a geometry object (point, line,
// polygon, etc.)
type Style struct {
	name      string
	alpha     uint8
	red       uint8
	green     uint8
	blue      uint8
	iconURL   string
	iconScale float64
	fill      int8
}

// NewStyle returns a new instance of a Style.  The alpha, red, green, and
// blue color properties are applied to point icon color as well as line and
// polygon color.  Name must be a single word (no spaces).
func NewStyle(name string, alpha uint8, red uint8, green uint8, blue uint8) *Style {
	return &Style{name, alpha, red, green, blue, "http://maps.google.com/mapfiles/kml/pushpin/ylw-pushpin.png", 1.1, 1}
}

// SetIconURL changes the icon that will be used for point placemarks.
// Built-in icon URL's can be seen in Google Earth when setting the
// placemark icon in the placemark properties dialog box.
func (s *Style) SetIconURL(url string) {
	url = strings.TrimSpace(url)

	if len(url) > 0 {
		s.iconURL = url
	}
}

// SetIconScale changes the icon scale from the default of 1.1.  Valid values
// are between 0.0 and 100.0.  Invalid values are ignored.
func (s *Style) SetIconScale(scale float64) {
	if scale >= 0.0 && scale <= 100.0 {
		s.iconScale = scale
	}
}

// SetPolygonFill specifies whether to fill in polygons.  The default is to
// not fill in the polygon.
func (s *Style) SetPolygonFill(fill bool) {
	if fill == true {
		s.fill = 1
	} else {
		s.fill = 0
	}
}

func (s *Style) render() string {
	colorStr := fmt.Sprintf("<color>%02x%02x%02x%02x</color>\n", s.alpha, s.blue, s.green, s.red) // yes, ABGR
	ret := fmt.Sprintf("<Style id=\"%s\">\n", s.name) +
		"<IconStyle>\n" +
		colorStr +
		fmt.Sprintf("<scale>%f</scale>\n", s.iconScale) +
		fmt.Sprintf("<Icon><href>%s</href></Icon>\n", s.iconURL) +
		"</IconStyle>\n" +
		"<LineStyle>\n" +
		colorStr +
		"<width>3</width>\n" +
		"</LineStyle>\n" +
		"<PolyStyle>\n" +
		colorStr +
		"<colorMode>normal</colorMode>\n" +
		fmt.Sprintf("<fill>%d</fill>\n", s.fill) +
		"<outline>1</outline>\n" +
		"</PolyStyle>\n" +
		"</Style>\n"

	return ret
}

// Point represents a point on the Earth
type Point struct {
	Lat float64 // latitude
	Lon float64 // longitude
	Alt float64 // altitude in meters
}

// NewPoint returns a pointer to a new Point instance.  Invalid points (those
// that with lat outside of +/-90.0 and lon outside of +/-180.0, NaN or Inf
// will return nil.
func NewPoint(lat float64, lon float64, alt float64) *Point {
	if math.IsNaN(lat) || math.IsInf(lat, 0) {
		return nil
	}

	if lat > 90.0 || lat < -90.0 {
		return nil
	}

	if math.IsNaN(lon) || math.IsInf(lon, 0) {
		return nil
	}

	if lon > 180.0 || lon < -180.0 {
		return nil
	}

	if math.IsNaN(alt) || math.IsInf(alt, 0) {
		alt = 0.0
	}

	return &Point{lat, lon, alt}
}

func (p *Point) render() string {
	ret := "<Point>\n" +
		"<extrude>0</extrude>\n" +
		"<altitudeMode>clampToGround</altitudeMode>\n" +
		fmt.Sprintf("<coordinates>%f,%f,%f</coordinates>\n", p.Lon, p.Lat, p.Alt) +
		"</Point>\n"

	return ret
}

// LineString represents a series of lines in a KML document.
type LineString struct {
	coordinates []*Point
	mutex       *sync.Mutex
}

// NewLineString returns a new instance of LineString.
func NewLineString() *LineString {
	ls := make([]*Point, 0, 10)
	return &LineString{ls, new(sync.Mutex)}
}

// Adds a Point to the LineString.  In order to render, the LineString
// needs at least two Points.  Points that are nil are ignored.
func (ls *LineString) AddPoint(point *Point) {
	if point != nil {
		ls.mutex.Lock()
		ls.coordinates = append(ls.coordinates, point)
		ls.mutex.Unlock()
	}
}

func (ls *LineString) render() string {
	if len(ls.coordinates) < 2 {
		return ""
	}

	ret := "<LineString>\n" +
		"<extrude>0</extrude>\n" +
		"<tessellate>1</tessellate>\n" +
		"<altitudeMode>clampToGround</altitudeMode>\n" +
		"<coordinates>\n"

	for _, coord := range ls.coordinates {
		ret += fmt.Sprintf("%f,%f,%f\n", coord.Lon, coord.Lat, coord.Alt)
	}

	ret += "</coordinates>\n" +
		"</LineString>\n"

	return ret
}

// Polygon represents a polygon in the KML document.  Must be added to a
// Placemark in order to render.
type Polygon struct {
	points []*Point
	mutex  *sync.Mutex
}

// NewPolygon returns a new instance of Polygon.
func NewPolygon() *Polygon {
	p := make([]*Point, 0, 4)
	return &Polygon{p, new(sync.Mutex)}
}

// AddPoint add a point (vertex) to the Polygon instance.  The Polygon will
// automatically close the ring if the last Point does not match the first
// Point.  For example, a box needs four Points, but if only three are added,
// then the fourth Point (which matches the first) will be added when the
// Polygon is rendered.
func (poly *Polygon) AddPoint(point *Point) {
	if point != nil {
		poly.mutex.Lock()
		poly.points = append(poly.points, point)
		poly.mutex.Unlock()
	}
}

func (poly *Polygon) render() string {
	if len(poly.points) == 0 {
		return ""
	}

	firstPoint := poly.points[0]
	lastPoint := poly.points[len(poly.points)-1]

	if *lastPoint != *firstPoint {
		poly.AddPoint(firstPoint) // close the polygon
	}

	ret := "<Polygon>\n" +
		"<extrude>1</extrude>\n" +
		"<altitudeMode>clampToGround</altitudeMode>\n" +
		"<outerBoundaryIs>\n" +
		"<LinearRing>\n" +
		"<coordinates>\n"

	for _, point := range poly.points {
		ret += fmt.Sprintf("%f,%f,%f\n", point.Lon, point.Lat, point.Alt)
	}

	ret += "</coordinates>\n" +
		"</LinearRing>\n" +
		"</Polygon>\n"

	return ret
}

// Placemark represents a placemark in the KML document.  All geometry
// objects (points, lines, polygons, etc.) must be within a Placemark
// instance.
type Placemark struct {
	name        string
	description string
	geometry    renderable
	style       string
}

// NewPlacemark returns a pointer to a new Placemark instance.  It takes a
// name, description, and a geometry object (Point, Polygon, etc.) as
// parameters.
func NewPlacemark(name string, desc string, geom renderable) *Placemark {
	return &Placemark{name, desc, geom, ""}
}

// SetStyle sets the style of the Placemark to the specified name.  The KML
// document must have a Style instance with a matching name (see NewStyle).
func (pm *Placemark) SetStyle(name string) {
	name = strings.TrimSpace(name)

	if len(name) > 0 {
		pm.style = name
	}
}

func (pm *Placemark) render() string {
	ret := "<Placemark>\n" +
		fmt.Sprintf("<name>%s</name>\n", pm.name) +
		fmt.Sprintf("<description>%s</description>\n", pm.description) +
		"<visibility>1</visibility>\n"

	if len(pm.style) > 0 {
		ret += fmt.Sprintf("<styleUrl>#%s</styleUrl>\n", pm.style)
	}

	ret += pm.geometry.render() +
		"</Placemark>\n"

	return ret
}
