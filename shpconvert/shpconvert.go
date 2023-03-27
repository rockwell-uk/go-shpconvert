package shpconvert

import (
	"fmt"

	"github.com/rockwell-uk/shapefile/shp"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkb"
	"github.com/twpayne/go-geom/encoding/wkt"
)

func ShpToWKB(s shp.Shape) ([]byte, error) {
	var funcName string = "shpconvert.ShpToWKB"

	r, err := getGeom(s)
	if err != nil {
		return []byte{}, fmt.Errorf("%v: %v", funcName, err.Error())
	}

	switch w := r.(type) {
	case *geom.Point:
		return wkb.Marshal(w, wkb.NDR)

	case *geom.LineString:
		return wkb.Marshal(w, wkb.NDR)

	case *geom.Polygon:
		return wkb.Marshal(w, wkb.NDR)
	}

	return []byte{}, nil
}

func ShpToWKT(s shp.Shape) (string, error) {
	var funcName string = "shpconvert.ShpToWKT"

	r, err := getGeom(s)
	if err != nil {
		return "", fmt.Errorf("%v: %v", funcName, err.Error())
	}

	switch w := r.(type) {
	case *geom.Point:
		return wkt.Marshal(w, wkt.EncodeOptionWithMaxDecimalDigits(6))

	case *geom.LineString:
		return wkt.Marshal(w, wkt.EncodeOptionWithMaxDecimalDigits(6))

	case *geom.Polygon:
		return wkt.Marshal(w, wkt.EncodeOptionWithMaxDecimalDigits(6))
	}

	return "", nil
}

func getGeom(s shp.Shape) (interface{}, error) {
	var funcName string = "shpconvert.getGeom"

	var r interface{}

	switch t := s.(type) {
	case *shp.Point:
		pt, _ := s.(*shp.Point)
		r = geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{pt.X, pt.Y})

	case *shp.PointZ:
		pt, _ := s.(*shp.PointZ)
		r = geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{pt.X, pt.Y})

	case *shp.Polyline:
		line, _ := s.(*shp.Polyline)
		r = addLinearPoints(geom.NewLineString(geom.XY), line.Points)

	case *shp.PolylineZ:
		line, _ := s.(*shp.PolylineZ)
		r = addLinearPoints(geom.NewLineString(geom.XY), line.Points)

	case *shp.Polygon:
		polygon, _ := s.(*shp.Polygon)
		r = addPolygonPoints(geom.NewPolygon(geom.XY), polygon.Points, polygon.Parts, polygon.NumberOfParts)

	case *shp.PolygonZ:
		polygon, _ := s.(*shp.PolygonZ)
		r = addPolygonPoints(geom.NewPolygon(geom.XY), polygon.Points, polygon.Parts, polygon.NumberOfParts)

	default:
		return "", fmt.Errorf("%v: %v", funcName, fmt.Sprintf("geometry type %v not yet supported", t))
	}

	return r, nil
}

func addLinearPoints(r *geom.LineString, points []shp.Point) *geom.LineString {
	coords := []geom.Coord{}

	for _, pt := range points {
		coords = append(coords, geom.Coord{pt.X, pt.Y})
	}

	r.MustSetCoords(coords)

	return r
}

func addPolygonPoints(r *geom.Polygon, points []shp.Point, parts []int32, numParts int32) *geom.Polygon {
	var allcoords [][]geom.Coord
	var startIndex, endIndex int32

	for i := 0; i < int(numParts); i++ {
		var coords []geom.Coord

		startIndex = parts[int32(i)]

		if int32(i) == numParts-1 {
			endIndex = int32(len(points))
		} else {
			endIndex = parts[i+1]
		}

		for j := startIndex; j < endIndex; j++ {
			coords = append(coords, geom.Coord{points[j].X, points[j].Y})
		}

		allcoords = append(allcoords, coords)
	}

	r.MustSetCoords(allcoords)

	return r
}
