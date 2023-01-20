package shpconvert

import (
	"errors"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/rockwell-uk/shapefile"
)

//nolint:dupl
func TestShpToWKT(t *testing.T) {

	tests := map[string]struct {
		shapeFile      string
		expectedValues map[int]string
	}{
		"HP_Woodland": {
			shapeFile: "./../testdata/HP_Woodland.shp",
			expectedValues: map[int]string{
				0: "POLYGON ((462567.97 1209327.8, 462571.3 1209324.9, 462592.75 1209349.57, 462629.2 1209318.7, 462593.99 1209277.01, 462577.75 1209290.42, 462553.25 1209310.65, 462550.4 1209313, 462565.06 1209330.33, 462567.97 1209327.8))",
				1: "POLYGON ((462209.3 1209502.3, 462183.7 1209509.2, 462173.3 1209555, 462138.89 1209557.29, 462139.49 1209563.43, 462191.9 1209575.1, 462209.3 1209502.3))",
			},
		},
		"HP_AdministrativeBoundary": {
			shapeFile: "./../testdata/HP_AdministrativeBoundary.shp",
			expectedValues: map[int]string{
				0: "LINESTRING (456548.4 1201820.1, 456551.5 1201821.5)",
				1: "LINESTRING (454331.4 1202522.2, 454332.4 1202517.4)",
				2: "LINESTRING (465969.5 1206905, 465000 1205924.33, 464898.5 1205824.5, 464546.5 1205470, 464340 1205260, 464262.5 1205182, 464082.68 1205000, 463267 1204175, 462249 1203142.5)",
				3: "LINESTRING (465985 1209231, 465002.94 1209861.07)",
			},
		},
		"SD_MotorwayJunction": {
			shapeFile: "./../testdata/SD_MotorwayJunction.shp",
			expectedValues: map[int]string{
				0:  "POINT (353384.54 482523.28)",
				1:  "POINT (359579.6 493000.63)",
				2:  "POINT (392648.96 411628.48)",
				3:  "POINT (398850.47 414805.24)",
				4:  "POINT (367627.02 406673.01)",
				5:  "POINT (368465.55 424560.52)",
				6:  "POINT (370324.18 405452.04)",
				7:  "POINT (370569.21 425378.2)",
				8:  "POINT (371522.35 428998.69)",
				9:  "POINT (374220.06 404833.61)",
				10: "POINT (374259.53 429967.66)",
				11: "POINT (375080.39 403866.13)",
				12: "POINT (375562.17 401832.96)",
				13: "POINT (376083.08 403429.4)",
				14: "POINT (376336.47 403092.78)",
				15: "POINT (377978.34 403421.07)",
				16: "POINT (378390.48 431717.16)",
				17: "POINT (380199.35 432015.86)",
				18: "POINT (380418.88 414982.68)",
				19: "POINT (380911.6 404720.54)",
				20: "POINT (382112.24 410621.11)",
				21: "POINT (382120.34 409143.5)",
				22: "POINT (382601.53 432675.89)",
				23: "POINT (382785.8 406010.1)",
				24: "POINT (383544.57 433683.99)",
				25: "POINT (384388.6 404836.78)",
				26: "POINT (384648.01 437139.17)",
				27: "POINT (385696.74 438522.33)",
				28: "POINT (386273.02 404298.02)",
				29: "POINT (386353.44 408835.22)",
				30: "POINT (387457.63 439786.67)",
				31: "POINT (389413.37 403420.61)",
				32: "POINT (389611.25 410024.28)",
				33: "POINT (390539.95 402619.46)",
				34: "POINT (390815.01 406672.99)",
				35: "POINT (337007 400085)",
				36: "POINT (356311.79 401697.79)",
				37: "POINT (340036.48 402340.51)",
				38: "POINT (350437.42 404290.49)",
				39: "POINT (353821.68 404408.42)",
				40: "POINT (345321.88 404764.57)",
				41: "POINT (347806.77 405178.32)",
				42: "POINT (363982.64 408776.18)",
				43: "POINT (354147.79 410528.87)",
				44: "POINT (358815.54 419468.46)",
				45: "POINT (355366.97 422240.71)",
				46: "POINT (362980.49 424532.37)",
				47: "POINT (356309.08 424612.16)",
				48: "POINT (356602.06 424750.91)",
				49: "POINT (355740.62 424818.28)",
				50: "POINT (358422 424834)",
				51: "POINT (358424.31 424842.35)",
				52: "POINT (357186.14 427205.83)",
				53: "POINT (358190.11 429984.61)",
				54: "POINT (356623.08 432803.57)",
				55: "POINT (335261.16 433511.29)",
				56: "POINT (352918.01 434082.85)",
				57: "POINT (341491.07 434887.17)",
				58: "POINT (353980.69 434943.94)",
				59: "POINT (348861.02 454507.76)",
				60: "POINT (349628.98 464206.95)",
				61: "POINT (351053.07 470478.84)",
				62: "POINT (350897 471969)",
				63: "POINT (374632.7 400515.33)",
			},
		},
	}

	for name, test := range tests {
		// Open the shapefile for reading
		sr, err := os.Open(test.shapeFile)
		if err != nil {
			t.Fatal(err)
		}
		defer sr.Close()

		// Open the dbf for reading
		dr, err := os.Open(strings.Replace(test.shapeFile, ".shp", ".dbf", 1))
		if err != nil {
			t.Fatal(err)
		}
		defer dr.Close()

		r, err := shapefile.NewReader(sr, shapefile.WithDBF(dr))
		if err != nil {
			t.Fatal(err)
		}

		i := 0
		for {
			rec, err := r.Next()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				t.Fatal(err)
			}

			actual, err := ShpToWKT(rec.Shape)
			if err != nil {
				t.Fatal(err)
			}

			expected := test.expectedValues[i]

			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("%v [%v]: expected %v\ngot %v", name, i, expected, actual)
			}

			i++
		}
	}
}
