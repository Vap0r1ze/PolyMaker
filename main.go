package main

import (
	"fmt"
	"image/jpeg"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/vap0r1ze/PolyMaker/shapes"

	"github.com/anthonynsimon/bild/effect"
	"github.com/fogleman/delaunay"
)

func main() {
	var index int
	rand.Seed(time.Now().UTC().UnixNano())

	if len(os.Args) < 7 {
		fmt.Println("not enough args")
		return
	}
	lineLow, err := strconv.Atoi(os.Args[1])
	lineHigh, err := strconv.Atoi(os.Args[2])
	linePointCount, err := strconv.Atoi(os.Args[3])
	randPointCount, err := strconv.Atoi(os.Args[4])
	edgePointCount, err := strconv.Atoi(os.Args[5])
	pointMargin, err := strconv.Atoi(os.Args[6])
	if err != nil {
		fmt.Println("could not parse arg")
		fmt.Println(err)
		return
	}

	var (
		faces           [][]delaunay.Point
		facePointsCount int
	)
	if len(os.Args) > 7 {
		var curVal string

		faceCount, err := strconv.Atoi(os.Args[8])
		if err != nil {
			fmt.Println("could not parse arg")
			fmt.Println(err)
			return
		}

		facePointCountsStr := os.Args[9]

		faces = make([][]delaunay.Point, faceCount)
		facePointCounts := make([]int, faceCount)

		index = 0
		for _, b := range facePointCountsStr + "\x00" {
			// if strIndex >= len(facePointCountsStr) {
			// 	break
			// }
			if b == ',' || b == 0 {
				val, err := strconv.Atoi(curVal)
				if err != nil {
					fmt.Println("could not parse arg")
					fmt.Println(err)
					return
				}
				facePointCounts[index] = val
				curVal = ""
				index++
			} else {
				curVal += string(b)
			}
		}

		for _, cnt := range facePointCounts {
			facePointsCount += cnt
		}

		facePointsStr := os.Args[7]

		index = 0
		pointIndex := 0
		curFace := make([]delaunay.Point, facePointCounts[index])
		curPoint := delaunay.Point{}
		curVal = ""
		for _, b := range facePointsStr + "\x00" {
			if b == ';' || b == 0 {
				val, err := strconv.Atoi(curVal)
				if err != nil {
					fmt.Println("could not parse arg")
					fmt.Println(err)
					return
				}
				curVal = ""
				curPoint.Y = float64(val)
				curFace[pointIndex] = curPoint
				curPoint = delaunay.Point{}
				pointIndex++
				continue
			}
			if b == '|' {
				val, err := strconv.Atoi(curVal)
				if err != nil {
					fmt.Println("could not parse arg")
					fmt.Println(err)
					return
				}
				curVal = ""
				curPoint.Y = float64(val)
				curFace[pointIndex] = curPoint
				faces[index] = curFace
				curPoint = delaunay.Point{}
				index++
				pointIndex = 0
				curFace = make([]delaunay.Point, facePointCounts[index])
				continue
			}
			if b == ',' {
				val, err := strconv.Atoi(curVal)
				if err != nil {
					fmt.Println("could not parse arg")
					fmt.Println(err)
					return
				}
				curVal = ""
				curPoint.X = float64(val)
				continue
			}
			curVal += string(b)
		}
		faces[index] = curFace
	}

	pointCount := linePointCount + randPointCount + edgePointCount

	points := make([]delaunay.Point, pointCount)

	content, err := os.Open("input.jpg")
	if err != nil {
		fmt.Println("file not found")
		fmt.Println(err)
		return
	}

	img, err := jpeg.Decode(content)
	if err != nil {
		fmt.Println("cannot decode image")
		fmt.Println(err)
		return
	}

	bounds := img.Bounds()

	emboss := effect.Emboss(img)

	linePix := make(map[delaunay.Point]bool)

	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			pix := emboss.At(x, y)
			r, g, b, _ := pix.RGBA()
			avg := int(r+g+b) / 3
			if avg < lineLow || avg > lineHigh {
				point := delaunay.Point{
					X: float64(x),
					Y: float64(y),
				}
				linePix[point] = true
			}
		}
	}

	linePixCount := len(linePix)
	linePointIndexes := make(map[int]bool)
	for len(linePointIndexes) < linePointCount {
		pointIndex := rand.Intn(linePixCount)
		linePointIndexes[pointIndex] = true
	}
	linePoints := make([]delaunay.Point, linePointCount)
	index = 0
	linePointIndex := 0
	for point := range linePix {
		if linePointIndexes[index] {
			linePoints[linePointIndex] = point
			linePointIndex++
		}
		index++
	}

	randPointInts := make(map[int]bool)
	for len(randPointInts) < randPointCount {
		pointInt := IntnMargin(bounds.Max.Y, pointMargin)*bounds.Max.X + IntnMargin(bounds.Max.X, pointMargin)
		randPointInts[pointInt] = true
	}
	randPoints := make([]delaunay.Point, randPointCount)
	index = 0
	for pointInt := range randPointInts {
		X := pointInt % bounds.Max.X
		Y := pointInt / bounds.Max.X
		randPoints[index] = delaunay.Point{
			X: float64(X),
			Y: float64(Y),
		}
		index++
	}

	edgePointInts := make(map[int]bool)
	for len(edgePointInts) < edgePointCount {
		var x int
		var y int

		if rand.Intn(2) == 0 {
			x = IntnMargin(bounds.Max.X, pointMargin)
			if rand.Intn(2) == 0 {
				y = bounds.Max.Y
			}
		} else {
			y = IntnMargin(bounds.Max.Y, pointMargin)
			if rand.Intn(2) == 0 {
				x = bounds.Max.X
			}
		}
		pointInt := y*bounds.Max.X + x
		edgePointInts[pointInt] = true
	}
	edgePoints := make([]delaunay.Point, edgePointCount)
	index = 0
	for pointInt := range edgePointInts {
		X := pointInt % bounds.Max.X
		Y := pointInt / bounds.Max.X
		edgePoints[index] = delaunay.Point{
			X: float64(X),
			Y: float64(Y),
		}
		index++
	}

	index = 0
	for _, point := range linePoints {
		points[index] = point
		index++
	}
	for _, point := range randPoints {
		points[index] = point
		index++
	}
	for _, point := range edgePoints {
		points[index] = point
		index++
	}

	notFacePoints := make(map[delaunay.Point]bool)
	for _, point := range points {
		p := shapes.Point{
			X: int(point.X),
			Y: int(point.Y),
		}
		var inFace bool
		for _, facePoints := range faces {
			t, err := delaunay.Triangulate(facePoints)
			if err != nil {
				fmt.Println("could not triangulate points")
				fmt.Println(err)
				return
			}

			ts := shapes.FromTriangulation(t)
			for _, tr := range ts {
				if tr.ContainsPoint(p) {
					inFace = true
					break
				}
			}
			if inFace {
				break
			}
		}
		if !inFace {
			notFacePoints[point] = true
		}
	}

	points = make([]delaunay.Point, len(notFacePoints)+facePointsCount+4)
	index = 0
	for p := range notFacePoints {
		points[index] = p
		index++
	}
	for _, face := range faces {
		for _, p := range face {
			points[index] = p
			index++
		}
	}
	points[index] = delaunay.Point{
		X: 0,
		Y: 0,
	}
	index++
	points[index] = delaunay.Point{
		X: float64(bounds.Max.X),
		Y: 0,
	}
	index++
	points[index] = delaunay.Point{
		X: 0,
		Y: float64(bounds.Max.Y),
	}
	index++
	points[index] = delaunay.Point{
		X: float64(bounds.Max.X),
		Y: float64(bounds.Max.Y),
	}
	index++

	t, err := delaunay.Triangulate(points)
	if err != nil {
		fmt.Println("could not triangulate points")
		fmt.Println(err)
		return
	}

	ts := shapes.FromTriangulation(t)

	tSums := make(map[int]int)
	tPixs := make(map[int]int)
	lastTrIndex := -1
	var lastTr shapes.Triangle
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			r, g, b, a := img.At(x, y).RGBA()
			point := shapes.Point{
				X: x,
				Y: y,
			}
			notLast := true
			if lastTrIndex > -1 {
				if lastTr.ContainsPoint(point) {
					notLast = false
					tSums[lastTrIndex*4] += int(r*uint32(math.Pow(2, 12))) / 0xFFFF * 0xFF
					tSums[lastTrIndex*4+1] += int(g*uint32(math.Pow(2, 8))) / 0xFFFF * 0xFF
					tSums[lastTrIndex*4+2] += int(b*uint32(math.Pow(2, 4))) / 0xFFFF * 0xFF
					tSums[lastTrIndex*4+3] += int(a) / 0xFFFF * 0xFF
					tPixs[lastTrIndex]++
				}
			}
			if notLast {
				for trIndex, tr := range ts {
					if tr.ContainsPoint(point) {
						tSums[trIndex*4] += int(r*uint32(math.Pow(2, 12))) / 0xFFFF * 0xFF
						tSums[trIndex*4+1] += int(g*uint32(math.Pow(2, 8))) / 0xFFFF * 0xFF
						tSums[trIndex*4+2] += int(b*uint32(math.Pow(2, 4))) / 0xFFFF * 0xFF
						tSums[trIndex*4+3] += int(a) / 0xFFFF * 0xFF
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
	_, err = f.WriteString("<svg viewBox=\"0 0 " + strconv.FormatInt(int64(bounds.Max.X), 10) + " " + strconv.FormatInt(int64(bounds.Max.Y), 10) + "\">")
	_, err = f.WriteString("<defs><filter id=\"median\"><feConvolveMatrix in=\"SourceGraphic\" order=\"3\" kernelMatrix=\"0 0 0 1 0 0 0 0 0\" result=\"1\" preserveAlpha=\"true\"/><feConvolveMatrix in=\"SourceGraphic\" order=\"3\" kernelMatrix=\"1 0 0 0 0 0 0 0 0\" result=\"2\" preserveAlpha=\"true\"/><feConvolveMatrix in=\"SourceGraphic\" order=\"3\" kernelMatrix=\"0 1 0 0 0 0 0 0 0\" result=\"3\" preserveAlpha=\"true\"/><feConvolveMatrix in=\"SourceGraphic\" order=\"3\" kernelMatrix=\"0 0 1 0 0 0 0 0 0\" result=\"4\" preserveAlpha=\"true\"/><feConvolveMatrix in=\"SourceGraphic\" order=\"3\" kernelMatrix=\"0 0 0 0 0 1 0 0 0\" result=\"5\" preserveAlpha=\"true\"/><feConvolveMatrix in=\"SourceGraphic\" order=\"3\" kernelMatrix=\"0 0 0 0 0 0 0 0 1\" result=\"6\" preserveAlpha=\"true\"/><feConvolveMatrix in=\"SourceGraphic\" order=\"3\" kernelMatrix=\"0 0 0 0 0 0 0 1 0\" result=\"7\" preserveAlpha=\"true\"/><feConvolveMatrix in=\"SourceGraphic\" order=\"3\" kernelMatrix=\"0 0 0 0 0 0 1 0 0\" result=\"8\" preserveAlpha=\"true\" /><feConvolveMatrix in=\"SourceGraphic\" order=\"3\" kernelMatrix=\"0 0 0 0 1 0 0 0 0\" result=\"9\" preserveAlpha=\"true\" /><feBlend in=\"1\" in2=\"2\" mode=\"lighten\" result=\"a1\"/><feBlend in=\"1\" in2=\"2\" mode=\"darken\" result=\"a2\"/><feBlend in=\"a2\" in2=\"3\" mode=\"lighten\" result=\"a3\"/><feBlend in=\"a2\" in2=\"3\" mode=\"darken\" result=\"a4\"/><feBlend in=\"a4\" in2=\"4\" mode=\"lighten\" result=\"a5\"/><feBlend in=\"a4\" in2=\"4\" mode=\"darken\" result=\"a6\"/><feBlend in=\"a6\" in2=\"5\" mode=\"lighten\" result=\"a7\"/><feBlend in=\"a6\" in2=\"5\" mode=\"darken\" result=\"a8\"/><feBlend in=\"a8\" in2=\"6\" mode=\"lighten\" result=\"a9\"/><feBlend in=\"a8\" in2=\"6\" mode=\"darken\" result=\"a10\"/><feBlend in=\"a10\" in2=\"7\" mode=\"lighten\" result=\"a11\"/><feBlend in=\"a10\" in2=\"7\" mode=\"darken\" result=\"a12\"/><feBlend in=\"a12\" in2=\"8\" mode=\"lighten\" result=\"a13\"/><feBlend in=\"a13\" in2=\"8\" mode=\"darken\" result=\"a14\"/><feBlend in=\"1\" in2=\"2\" mode=\"lighten\" result=\"a15\"/><feBlend in=\"1\" in2=\"2\" mode=\"darken\" result=\"a16\"/>    <feBlend in=\"a1\" in2=\"a3\" mode=\"lighten\" result=\"b1\"/><feBlend in=\"a1\" in2=\"a3\" mode=\"darken\" result=\"b2\"/><feBlend in=\"b2\" in2=\"a5\" mode=\"lighten\" result=\"b3\"/><feBlend in=\"b2\" in2=\"a5\" mode=\"darken\" result=\"b4\"/><feBlend in=\"b4\" in2=\"a7\" mode=\"lighten\" result=\"b5\"/><feBlend in=\"b4\" in2=\"a7\" mode=\"darken\" result=\"b6\"/><feBlend in=\"b6\" in2=\"a9\" mode=\"lighten\" result=\"b7\"/><feBlend in=\"b6\" in2=\"a9\" mode=\"darken\" result=\"b8\"/><feBlend in=\"b8\" in2=\"a11\" mode=\"lighten\" result=\"b9\"/><feBlend in=\"b8\" in2=\"a11\" mode=\"darken\" result=\"b10\"/><feBlend in=\"b10\" in2=\"a13\" mode=\"lighten\" result=\"b11\"/><feBlend in=\"b10\" in2=\"a13\" mode=\"darken\" result=\"b12\"/><feBlend in=\"b12\" in2=\"a15\" mode=\"lighten\" result=\"b13\"/><feBlend in=\"b12\" in2=\"a15\" mode=\"darken\" result=\"b14\"/><feBlend in=\"b1\" in2=\"b3\" mode=\"lighten\" result=\"c1\"/><feBlend in=\"b1\" in2=\"b3\" mode=\"darken\" result=\"c2\"/><feBlend in=\"c2\" in2=\"b5\" mode=\"lighten\" result=\"c3\"/><feBlend in=\"c2\" in2=\"b5\" mode=\"darken\" result=\"c4\"/><feBlend in=\"c4\" in2=\"b7\" mode=\"lighten\" result=\"c5\"/><feBlend in=\"c4\" in2=\"b7\" mode=\"darken\" result=\"c6\"/><feBlend in=\"c6\" in2=\"b9\" mode=\"lighten\" result=\"c7\"/><feBlend in=\"c6\" in2=\"b9\" mode=\"darken\" result=\"c8\"/><feBlend in=\"c8\" in2=\"b11\" mode=\"lighten\" result=\"c9\"/><feBlend in=\"c8\" in2=\"b11\" mode=\"darken\" result=\"c10\"/><feBlend in=\"c10\" in2=\"b13\" mode=\"lighten\" result=\"c11\"/><feBlend in=\"c10\" in2=\"b13\" mode=\"darken\" result=\"c12\"/><feBlend in=\"c1\" in2=\"c3\" mode=\"lighten\" result=\"d1\"/><feBlend in=\"d1\" in2=\"c3\" mode=\"darken\" result=\"d2\"/><feBlend in=\"d2\" in2=\"c5\" mode=\"lighten\" result=\"d3\"/><feBlend in=\"d2\" in2=\"c5\" mode=\"darken\" result=\"d4\"/><feBlend in=\"d4\" in2=\"c7\" mode=\"lighten\" result=\"d5\"/><feBlend in=\"d4\" in2=\"c7\" mode=\"darken\" result=\"d6\"/><feBlend in=\"d6\" in2=\"c9\" mode=\"lighten\" result=\"d7\"/><feBlend in=\"d6\" in2=\"c9\" mode=\"darken\" result=\"d8\"/><feBlend in=\"d8\" in2=\"c11\" mode=\"lighten\" result=\"d9\"/><feBlend in=\"d8\" in2=\"c11\" mode=\"darken\" result=\"d10\"/><feBlend in=\"d1\" in2=\"d3\" mode=\"darken\" result=\"e1\"/><feBlend in=\"e1\" in2=\"d5\" mode=\"darken\" result=\"e2\"/><feBlend in=\"e2\" in2=\"d7\" mode=\"darken\" result=\"e3\"/><feBlend in=\"e3\" in2=\"d9\" mode=\"darken\" result=\"e4\"/></filter></defs>")
	_, err = f.WriteString("<g filter=\"url(#median)\">")
	for tIndex, t := range ts {
		var svgElem string
		pixs := tPixs[tIndex]
		if pixs == 0 {
			continue
		}
		rSum := tSums[tIndex*4]
		gSum := tSums[tIndex*4+1]
		bSum := tSums[tIndex*4+2]
		aSum := tSums[tIndex*4+3]
		r := rSum / pixs / int(math.Pow(2, 12))
		g := gSum / pixs / int(math.Pow(2, 8))
		b := bSum / pixs / int(math.Pow(2, 4))
		a := aSum / pixs
		d := t.GetPathData()
		svgElem += "<path stroke=\"none\" fill=\"rgba("
		svgElem += strconv.FormatInt(int64(r), 10) + "," + strconv.FormatInt(int64(g), 10) + "," + strconv.FormatInt(int64(b), 10) + "," + strconv.FormatInt(int64(a), 10)
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

// IntnMargin - Generates a random integer within a certain margin
func IntnMargin(n, margin int) int {
	if margin >= n/2 {
		margin = n/2 - 1
	}

	r := rand.Intn(n)
	if r < margin {
		r = margin
	} else if r > (n - margin) {
		r = n - margin
	}

	return r
}
