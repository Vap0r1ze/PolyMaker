package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/vap0r1ze/PolyMaker/shapes"

	"github.com/fogleman/delaunay"
)

func main() {
	var (
		bounds     image.Rectangle
		content    *os.File
		curPoint   delaunay.Point
		curVal     string
		err        error
		img        image.Image
		index      int
		inputFile  string
		pointCount int
		points     []delaunay.Point
		pointsStr  string
		t          *delaunay.Triangulation
	)
	rand.Seed(time.Now().UTC().UnixNano())

	if len(os.Args) < 4 {
		fmt.Println("not enough args")
		return
	}

	inputFile = os.Args[1]
	content, err = os.Open(inputFile)
	if err != nil {
		fmt.Println("error opening file:")
		fmt.Println(err)
		return
	}

	pointsStr = os.Args[2]

	pointCount, err = strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("error parsing point count as integer:")
		fmt.Println(err)
		return
	}

	points = make([]delaunay.Point, pointCount)

	for _, b := range pointsStr + "\x00" {
		if b == ';' || b == 0 {
			val, err := strconv.Atoi(curVal)
			if err != nil {
				fmt.Println("error parsing point list:")
				fmt.Println(err)
				return
			}
			curVal = ""
			curPoint.Y = float64(val)
			points[index] = curPoint
			curPoint = delaunay.Point{}
			index++
			continue
		}
		if b == ',' {
			val, err := strconv.Atoi(curVal)
			if err != nil {
				fmt.Println("error parsing point list:")
				fmt.Println(err)
				return
			}
			curVal = ""
			curPoint.X = float64(val)
			continue
		}
		curVal += string(b)
	}

	img, err = jpeg.Decode(content)
	if err != nil {
		fmt.Println("cannot decode image")
		fmt.Println(err)
		return
	}

	t, err = delaunay.Triangulate(points)
	if err != nil {
		fmt.Println("could not triangulate points")
		fmt.Println(err)
		return
	}

	ts := shapes.FromTriangulation(t)

	bounds = img.Bounds()
	tSums := make([]int, len(ts)*3)
	tPixs := make([]int, len(ts))
	lastTrIndex := -1
	var lastTr shapes.Triangle
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			pix := img.At(x, y)
			ycbcrPix := pix.(color.YCbCr)
			r, g, b := color.YCbCrToRGB(ycbcrPix.Y, ycbcrPix.Cb, ycbcrPix.Cr)
			point := shapes.Point{
				X: x,
				Y: y,
			}
			notLast := true
			if lastTrIndex > -1 {
				if lastTr.ContainsPoint(point) {
					notLast = false
					tSums[lastTrIndex*3] += int(r)
					tSums[lastTrIndex*3+1] += int(g)
					tSums[lastTrIndex*3+2] += int(b)
					tPixs[lastTrIndex]++
				}
			}
			if notLast {
				for trIndex, tr := range ts {
					if tr.ContainsPoint(point) {
						tSums[trIndex*3] += int(r)
						tSums[trIndex*3+1] += int(g)
						tSums[trIndex*3+2] += int(b)
						tPixs[trIndex]++
						lastTr = tr
						lastTrIndex = trIndex
						break
					}
				}
			}
		}
	}

	f, err := os.Create("output.svg")
	if err != nil {
		fmt.Println("cannot create file")
		fmt.Println(err)
		return
	}
	_, err = f.WriteString("<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 " + strconv.FormatInt(int64(bounds.Max.X), 10) + " " + strconv.FormatInt(int64(bounds.Max.Y), 10) + "\">")
	_, err = f.WriteString("<defs><filter id=\"median\"><feConvolveMatrix in=\"SourceGraphic\" order=\"3\" kernelMatrix=\"0 0 0 1 0 0 0 0 0\" result=\"1\" preserveAlpha=\"true\"/><feConvolveMatrix in=\"SourceGraphic\" order=\"3\" kernelMatrix=\"1 0 0 0 0 0 0 0 0\" result=\"2\" preserveAlpha=\"true\"/><feConvolveMatrix in=\"SourceGraphic\" order=\"3\" kernelMatrix=\"0 1 0 0 0 0 0 0 0\" result=\"3\" preserveAlpha=\"true\"/><feConvolveMatrix in=\"SourceGraphic\" order=\"3\" kernelMatrix=\"0 0 1 0 0 0 0 0 0\" result=\"4\" preserveAlpha=\"true\"/><feConvolveMatrix in=\"SourceGraphic\" order=\"3\" kernelMatrix=\"0 0 0 0 0 1 0 0 0\" result=\"5\" preserveAlpha=\"true\"/><feConvolveMatrix in=\"SourceGraphic\" order=\"3\" kernelMatrix=\"0 0 0 0 0 0 0 0 1\" result=\"6\" preserveAlpha=\"true\"/><feConvolveMatrix in=\"SourceGraphic\" order=\"3\" kernelMatrix=\"0 0 0 0 0 0 0 1 0\" result=\"7\" preserveAlpha=\"true\"/><feConvolveMatrix in=\"SourceGraphic\" order=\"3\" kernelMatrix=\"0 0 0 0 0 0 1 0 0\" result=\"8\" preserveAlpha=\"true\" /><feConvolveMatrix in=\"SourceGraphic\" order=\"3\" kernelMatrix=\"0 0 0 0 1 0 0 0 0\" result=\"9\" preserveAlpha=\"true\" /><feBlend in=\"1\" in2=\"2\" mode=\"lighten\" result=\"a1\"/><feBlend in=\"1\" in2=\"2\" mode=\"darken\" result=\"a2\"/><feBlend in=\"a2\" in2=\"3\" mode=\"lighten\" result=\"a3\"/><feBlend in=\"a2\" in2=\"3\" mode=\"darken\" result=\"a4\"/><feBlend in=\"a4\" in2=\"4\" mode=\"lighten\" result=\"a5\"/><feBlend in=\"a4\" in2=\"4\" mode=\"darken\" result=\"a6\"/><feBlend in=\"a6\" in2=\"5\" mode=\"lighten\" result=\"a7\"/><feBlend in=\"a6\" in2=\"5\" mode=\"darken\" result=\"a8\"/><feBlend in=\"a8\" in2=\"6\" mode=\"lighten\" result=\"a9\"/><feBlend in=\"a8\" in2=\"6\" mode=\"darken\" result=\"a10\"/><feBlend in=\"a10\" in2=\"7\" mode=\"lighten\" result=\"a11\"/><feBlend in=\"a10\" in2=\"7\" mode=\"darken\" result=\"a12\"/><feBlend in=\"a12\" in2=\"8\" mode=\"lighten\" result=\"a13\"/><feBlend in=\"a13\" in2=\"8\" mode=\"darken\" result=\"a14\"/><feBlend in=\"1\" in2=\"2\" mode=\"lighten\" result=\"a15\"/><feBlend in=\"1\" in2=\"2\" mode=\"darken\" result=\"a16\"/>    <feBlend in=\"a1\" in2=\"a3\" mode=\"lighten\" result=\"b1\"/><feBlend in=\"a1\" in2=\"a3\" mode=\"darken\" result=\"b2\"/><feBlend in=\"b2\" in2=\"a5\" mode=\"lighten\" result=\"b3\"/><feBlend in=\"b2\" in2=\"a5\" mode=\"darken\" result=\"b4\"/><feBlend in=\"b4\" in2=\"a7\" mode=\"lighten\" result=\"b5\"/><feBlend in=\"b4\" in2=\"a7\" mode=\"darken\" result=\"b6\"/><feBlend in=\"b6\" in2=\"a9\" mode=\"lighten\" result=\"b7\"/><feBlend in=\"b6\" in2=\"a9\" mode=\"darken\" result=\"b8\"/><feBlend in=\"b8\" in2=\"a11\" mode=\"lighten\" result=\"b9\"/><feBlend in=\"b8\" in2=\"a11\" mode=\"darken\" result=\"b10\"/><feBlend in=\"b10\" in2=\"a13\" mode=\"lighten\" result=\"b11\"/><feBlend in=\"b10\" in2=\"a13\" mode=\"darken\" result=\"b12\"/><feBlend in=\"b12\" in2=\"a15\" mode=\"lighten\" result=\"b13\"/><feBlend in=\"b12\" in2=\"a15\" mode=\"darken\" result=\"b14\"/><feBlend in=\"b1\" in2=\"b3\" mode=\"lighten\" result=\"c1\"/><feBlend in=\"b1\" in2=\"b3\" mode=\"darken\" result=\"c2\"/><feBlend in=\"c2\" in2=\"b5\" mode=\"lighten\" result=\"c3\"/><feBlend in=\"c2\" in2=\"b5\" mode=\"darken\" result=\"c4\"/><feBlend in=\"c4\" in2=\"b7\" mode=\"lighten\" result=\"c5\"/><feBlend in=\"c4\" in2=\"b7\" mode=\"darken\" result=\"c6\"/><feBlend in=\"c6\" in2=\"b9\" mode=\"lighten\" result=\"c7\"/><feBlend in=\"c6\" in2=\"b9\" mode=\"darken\" result=\"c8\"/><feBlend in=\"c8\" in2=\"b11\" mode=\"lighten\" result=\"c9\"/><feBlend in=\"c8\" in2=\"b11\" mode=\"darken\" result=\"c10\"/><feBlend in=\"c10\" in2=\"b13\" mode=\"lighten\" result=\"c11\"/><feBlend in=\"c10\" in2=\"b13\" mode=\"darken\" result=\"c12\"/><feBlend in=\"c1\" in2=\"c3\" mode=\"lighten\" result=\"d1\"/><feBlend in=\"d1\" in2=\"c3\" mode=\"darken\" result=\"d2\"/><feBlend in=\"d2\" in2=\"c5\" mode=\"lighten\" result=\"d3\"/><feBlend in=\"d2\" in2=\"c5\" mode=\"darken\" result=\"d4\"/><feBlend in=\"d4\" in2=\"c7\" mode=\"lighten\" result=\"d5\"/><feBlend in=\"d4\" in2=\"c7\" mode=\"darken\" result=\"d6\"/><feBlend in=\"d6\" in2=\"c9\" mode=\"lighten\" result=\"d7\"/><feBlend in=\"d6\" in2=\"c9\" mode=\"darken\" result=\"d8\"/><feBlend in=\"d8\" in2=\"c11\" mode=\"lighten\" result=\"d9\"/><feBlend in=\"d8\" in2=\"c11\" mode=\"darken\" result=\"d10\"/><feBlend in=\"d1\" in2=\"d3\" mode=\"darken\" result=\"e1\"/><feBlend in=\"e1\" in2=\"d5\" mode=\"darken\" result=\"e2\"/><feBlend in=\"e2\" in2=\"d7\" mode=\"darken\" result=\"e3\"/><feBlend in=\"e3\" in2=\"d9\" mode=\"darken\" result=\"e4\"/></filter></defs>")
	_, err = f.WriteString("<g filter=\"url(#median)\">")
	for tIndex, t := range ts {
		var svgElem string
		pixs := tPixs[tIndex]
		if pixs == 0 {
			continue
		}
		rSum := tSums[tIndex*3]
		gSum := tSums[tIndex*3+1]
		bSum := tSums[tIndex*3+2]
		r := rSum / pixs
		g := gSum / pixs
		b := bSum / pixs
		d := t.GetPathData()
		svgElem += "<path stroke=\"none\" fill=\"rgb("
		svgElem += strconv.FormatInt(int64(r), 10) + "," + strconv.FormatInt(int64(g), 10) + "," + strconv.FormatInt(int64(b), 10)
		svgElem += ")\" d=\""
		svgElem += d
		svgElem += "\"/>"
		_, err = f.WriteString(svgElem)
	}
	_, err = f.WriteString("</g></svg>")

	if err != nil {
		fmt.Println("cannot write to file")
		fmt.Println(err)
		return
	}

	err = f.Close()
	if err != nil {
		fmt.Println("cannot close file")
		fmt.Println(err)
		return
	}

	fmt.Println("done")
	return
}
