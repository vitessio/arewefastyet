/*
Copyright 2021 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package report

import (
	"github.com/vitessio/arewefastyet/go/storage/influxdb"
	"strconv"
	"strings"

	"github.com/vitessio/arewefastyet/go/tools/microbench"

	"github.com/jung-kurt/gofpdf"
	"github.com/vitessio/arewefastyet/go/storage/mysql"
	"github.com/vitessio/arewefastyet/go/tools/git"
	"github.com/vitessio/arewefastyet/go/tools/macrobench"
)

// pdfMargin is the margin that is used while creating the pdf. It is used as the left, right, top and bottom margin.
const pdfMargin = 15.0

// cellStyle is a struct that has the elements used for styling a cell in gofpdf
type cellStyle struct {
	textCol   [3]int // RGB colors of the text
	fillCol   [3]int // RGB colors of the background
	borderStr string // border style string
	alignStr  string // alignment style string
}

// cellStyles contains the different cell styles that we use in the pdf
var cellStyles = []cellStyle{
	{
		textCol:   [3]int{224, 224, 224},
		fillCol:   [3]int{64, 64, 64},
		borderStr: "1",
		alignStr:  "CM",
	}, {
		textCol:   [3]int{24, 24, 24},
		fillCol:   [3]int{255, 255, 255},
		borderStr: "",
		alignStr:  "C",
	}, {
		textCol:   [3]int{255, 255, 255},
		fillCol:   [3]int{0, 0, 0},
		borderStr: "1",
		alignStr:  "CM",
	},
}

// tableCell contains the string value along with the styling index to use
type tableCell struct {
	value string
	// styleIndex is the index of cellStyles to use for styling
	styleIndex int
	// linkUrl if non empty, will make the cell a clickable link to the given URL
	linkUrl string
}

// GenerateCompareReport is used to generate a comparison report between the 2 SHAs provided. It uses the client connection
// to read the results. It also takes as an argument the name of the report that will be generated
func GenerateCompareReport(client *mysql.Client, metricsClient *influxdb.Client, fromSHA, toSHA, reportFile string) error {
	// Compare macrobenchmark results for the 2 SHAs
	macrosMatrices, err := macrobench.CompareMacroBenchmarks(client, metricsClient, fromSHA, toSHA)
	if err != nil {
		return err
	}

	// Compare microbenchmark results for the 2 SHAs
	microsMatrix, err := microbench.CompareMicroBenchmarks(client, fromSHA, toSHA)
	if err != nil {
		return err
	}

	// Create a new pdf and set the margins
	pdf := gofpdf.New(gofpdf.OrientationPortrait, "mm", "A4", "")
	pdf.SetMargins(pdfMargin, pdfMargin, pdfMargin)
	// Add a page to start
	pdf.AddPage()
	pdf.SetAutoPageBreak(true, pdfMargin)
	// Set the font for the title
	pdf.SetFont("Arial", "B", 28)
	// Print the title
	pdf.WriteAligned(0, 10, "Comparison Results", "C")
	pdf.Ln(-1)
	pdf.Ln(2)

	if len(macrosMatrices) != 0 {
		// Print the subtitle
		writeSubtitle(pdf, "Macro-benchmarks")
		macroTable := [][]tableCell{
			{
				tableCell{value: "Metric", styleIndex: 2},
				tableCell{value: git.ShortenSHA(fromSHA), styleIndex: 2, linkUrl: "https://github.com/vitessio/vitess/tree/" + fromSHA + "/"},
				tableCell{value: git.ShortenSHA(toSHA), styleIndex: 2, linkUrl: "https://github.com/vitessio/vitess/tree/" + toSHA + "/"},
			},
		}
		// range over all the macrobenchmarks
		for key, value := range macrosMatrices {
			// the map stores the comparisonArrays
			macroCompArr := value.(macrobench.ComparisonArray)
			if len(macroCompArr) > 0 {
				macroComp := macroCompArr[0]
				macroTable = append(macroTable, []tableCell{{value: strings.ToUpper(key.String()), styleIndex: 2}})
				macroTable = append(macroTable, []tableCell{{value: "TPS", styleIndex: 0}, {value: convertFloatToString(macroComp.Reference.Result.TPS), styleIndex: 1}, {value: convertFloatToString(macroComp.Compare.Result.TPS), styleIndex: 1}})
				macroTable = append(macroTable, []tableCell{{value: "QPS Reads", styleIndex: 0}, {value: convertFloatToString(macroComp.Reference.Result.QPS.Reads), styleIndex: 1}, {value: convertFloatToString(macroComp.Compare.Result.QPS.Reads), styleIndex: 1}})
				macroTable = append(macroTable, []tableCell{{value: "QPS Writes", styleIndex: 0}, {value: convertFloatToString(macroComp.Reference.Result.QPS.Writes), styleIndex: 1}, {value: convertFloatToString(macroComp.Compare.Result.QPS.Writes), styleIndex: 1}})
				macroTable = append(macroTable, []tableCell{{value: "QPS Total", styleIndex: 0}, {value: convertFloatToString(macroComp.Reference.Result.QPS.Total), styleIndex: 1}, {value: convertFloatToString(macroComp.Compare.Result.QPS.Total), styleIndex: 1}})
				macroTable = append(macroTable, []tableCell{{value: "QPS Others", styleIndex: 0}, {value: convertFloatToString(macroComp.Reference.Result.QPS.Other), styleIndex: 1}, {value: convertFloatToString(macroComp.Compare.Result.QPS.Other), styleIndex: 1}})
				macroTable = append(macroTable, []tableCell{{value: "Latency", styleIndex: 0}, {value: convertFloatToString(macroComp.Reference.Result.Latency), styleIndex: 1}, {value: convertFloatToString(macroComp.Compare.Result.Latency), styleIndex: 1}})
				macroTable = append(macroTable, []tableCell{{value: "Errors", styleIndex: 0}, {value: convertFloatToString(macroComp.Reference.Result.Errors), styleIndex: 1}, {value: convertFloatToString(macroComp.Compare.Result.Errors), styleIndex: 1}})
				macroTable = append(macroTable, []tableCell{{value: "Reconnects", styleIndex: 0}, {value: convertFloatToString(macroComp.Reference.Result.Reconnects), styleIndex: 1}, {value: convertFloatToString(macroComp.Compare.Result.Reconnects), styleIndex: 1}})
				macroTable = append(macroTable, []tableCell{{value: "Time", styleIndex: 0}, {value: strconv.Itoa(macroComp.Reference.Result.Time), styleIndex: 1}, {value: strconv.Itoa(macroComp.Compare.Result.Time), styleIndex: 1}})
				macroTable = append(macroTable, []tableCell{{value: "Threads", styleIndex: 0}, {value: convertFloatToString(macroComp.Reference.Result.Threads), styleIndex: 1}, {value: convertFloatToString(macroComp.Compare.Result.Threads), styleIndex: 1}})
			}
		}
		// write the table to pdf
		writeTableToPdf(pdf, macroTable)
		pdf.AddPage()
	}

	if len(microsMatrix) != 0 {
		// Print the subtitle
		writeSubtitle(pdf, "Micro-benchmarks")
		microTable := [][]tableCell{
			{
				tableCell{value: "Metric", styleIndex: 2},
				tableCell{value: git.ShortenSHA(fromSHA), styleIndex: 2, linkUrl: "https://github.com/vitessio/vitess/tree/" + fromSHA + "/"},
				tableCell{value: git.ShortenSHA(toSHA), styleIndex: 2, linkUrl: "https://github.com/vitessio/vitess/tree/" + toSHA + "/"},
			},
		}
		// range over all the microbenchmarks
		for _, microComp := range microsMatrix {
			microTable = append(microTable, []tableCell{{value: microComp.PkgName + "." + microComp.Name, styleIndex: 2}})
			microTable = append(microTable, []tableCell{{value: "Ops", styleIndex: 0}, {value: strconv.Itoa(microComp.Current.Ops), styleIndex: 1}, {value: strconv.Itoa(microComp.Last.Ops), styleIndex: 1}})
			microTable = append(microTable, []tableCell{{value: "NSPerOp", styleIndex: 0}, {value: convertFloatToString(microComp.Current.NSPerOp), styleIndex: 1}, {value: convertFloatToString(microComp.Last.NSPerOp), styleIndex: 1}})
			microTable = append(microTable, []tableCell{{value: "MBPerSec", styleIndex: 0}, {value: convertFloatToString(microComp.Current.MBPerSec), styleIndex: 1}, {value: convertFloatToString(microComp.Last.MBPerSec), styleIndex: 1}})
			microTable = append(microTable, []tableCell{{value: "BytesPerOp", styleIndex: 0}, {value: convertFloatToString(microComp.Current.BytesPerOp), styleIndex: 1}, {value: convertFloatToString(microComp.Last.BytesPerOp), styleIndex: 1}})
			microTable = append(microTable, []tableCell{{value: "AllocsPerOp", styleIndex: 0}, {value: convertFloatToString(microComp.Current.AllocsPerOp), styleIndex: 1}, {value: convertFloatToString(microComp.Last.AllocsPerOp), styleIndex: 1}})
			microTable = append(microTable, []tableCell{{value: "Ratio of NSPerOp", styleIndex: 0}, {value: convertFloatToString(microComp.CurrLastDiff), styleIndex: 1}})
		}
		// write the table to pdf
		writeTableToPdf(pdf, microTable)
	}
	err = pdf.OutputFileAndClose(reportFile)
	return err
}

// convertFloatToString converts the float input into a string so that it can be printed
func convertFloatToString(in float64) string {
	return strconv.FormatFloat(in, 'f', -1, 64)
}

// writeSubtitle is used to write a subtitle in the pdf provided
func writeSubtitle(pdf *gofpdf.Fpdf, subtitle string) {
	pdf.SetFont("Arial", "B", 20)
	pdf.WriteAligned(0, 10, subtitle, "C")
	pdf.Ln(-1)
	pdf.Ln(2)
}

// writeTableToPdf is used to write a table to the pdf provided
func writeTableToPdf(pdf *gofpdf.Fpdf, table [][]tableCell) {
	pageWidth, pageHeight := pdf.GetPageSize()
	lineHt := 5.5
	cellGap := 2.0

	type cellType struct {
		str  string
		list [][]byte
		ht   float64
	}
	var (
		cellList []cellType
		cell     cellType
	)

	pdf.SetFont("Arial", "", 14)

	// Rows
	y := pdf.GetY()
	for rowJ := 0; rowJ < len(table); rowJ++ {
		cellList = nil
		colCount := len(table[rowJ])
		colWd := (pageWidth - pdfMargin - pdfMargin) / float64(colCount)
		maxHt := lineHt
		// Cell height calculation loop
		// required because a cell might span multiple lines but then the entire row must span that many lines
		for colJ := 0; colJ < colCount; colJ++ {
			cell.str = table[rowJ][colJ].value
			cell.list = pdf.SplitLines([]byte(cell.str), colWd-cellGap-cellGap)
			cell.ht = float64(len(cell.list)) * lineHt
			if cell.ht > maxHt {
				maxHt = cell.ht
			}
			cellList = append(cellList, cell)
		}
		// Cell render loop
		x := pdfMargin
		for colJ := 0; colJ < len(cellList); colJ++ {
			if y+maxHt+cellGap+cellGap > pageHeight-pdfMargin-pdfMargin {
				pdf.AddPage()
				y = pdf.GetY()
			}
			pdf.Rect(x, y, colWd, maxHt+cellGap+cellGap, "D")
			cell = cellList[colJ]
			cellY := y + cellGap + (maxHt-cell.ht)/2
			// Get the styling from the index
			styling := cellStyles[table[rowJ][colJ].styleIndex]
			// Use it to set the colours
			pdf.SetTextColor(styling.textCol[0], styling.textCol[1], styling.textCol[2])
			pdf.SetFillColor(styling.fillCol[0], styling.fillCol[1], styling.fillCol[2])
			for splitJ := 0; splitJ < len(cell.list); splitJ++ {
				pdf.SetXY(x+cellGap, cellY)
				pdf.CellFormat(colWd-cellGap-cellGap, lineHt, string(cell.list[splitJ]), styling.borderStr, 0,
					styling.alignStr, true, 0, table[rowJ][colJ].linkUrl)
				cellY += lineHt
			}
			x += colWd
		}
		y += maxHt + cellGap + cellGap
	}
	pdf.Ln(-1)
	pdf.Ln(4)
}
