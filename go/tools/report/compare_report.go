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
	"strconv"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/vitessio/arewefastyet/go/mysql"
	"github.com/vitessio/arewefastyet/go/tools/git"
	"github.com/vitessio/arewefastyet/go/tools/macrobench"
	"github.com/vitessio/arewefastyet/go/tools/microbench"
)

const pdfMargin = 15.0

func GenerateCompareReport(client *mysql.Client, fromSHA, toSHA, reportFile string) error {
	macrosMatrices, err := macrobench.CompareMacroBenchmarks(client, fromSHA, toSHA)
	if err != nil {
		return err
	}
	microsMatrix, err := microbench.CompareMicroBenchmarks(client, fromSHA, toSHA)
	if err != nil {
		return err
	}
	pdf := gofpdf.New(gofpdf.OrientationPortrait, "mm", "A4", "")
	pdf.SetMargins(pdfMargin, pdfMargin, pdfMargin)
	pdf.AddPage()
	pdf.SetAutoPageBreak(true, pdfMargin)
	pdf.SetFont("Arial", "B", 28)
	pdf.WriteAligned(0, 10, "Comparison Results", "C")
	pdf.Ln(-1)
	pdf.Ln(2)

	if len(macrosMatrices) != 0 {
		writeSubtitle(pdf, "Macro-benchmarks")
		macroTable := [][]string{
			{"Metric", git.ShortenSHA(fromSHA), git.ShortenSHA(toSHA)},
		}
		for key, value := range macrosMatrices {
			macroCompArr := value.(macrobench.ComparisonArray)
			if len(macroCompArr) > 0 {
				macroComp := macroCompArr[0]
				macroTable = append(macroTable, []string{strings.ToUpper(key.String())})
				macroTable = append(macroTable, []string{"TPS", convertFloatToString(macroComp.Reference.Result.TPS), convertFloatToString(macroComp.Compare.Result.TPS)})
				macroTable = append(macroTable, []string{"QPS Reads", convertFloatToString(macroComp.Reference.Result.QPS.Reads), convertFloatToString(macroComp.Compare.Result.QPS.Reads)})
				macroTable = append(macroTable, []string{"QPS Writes", convertFloatToString(macroComp.Reference.Result.QPS.Writes), convertFloatToString(macroComp.Compare.Result.QPS.Writes)})
				macroTable = append(macroTable, []string{"QPS Total", convertFloatToString(macroComp.Reference.Result.QPS.Total), convertFloatToString(macroComp.Compare.Result.QPS.Total)})
				macroTable = append(macroTable, []string{"QPS Others", convertFloatToString(macroComp.Reference.Result.QPS.Other), convertFloatToString(macroComp.Compare.Result.QPS.Other)})
				macroTable = append(macroTable, []string{"Latency", convertFloatToString(macroComp.Reference.Result.Latency), convertFloatToString(macroComp.Compare.Result.Latency)})
				macroTable = append(macroTable, []string{"Errors", convertFloatToString(macroComp.Reference.Result.Errors), convertFloatToString(macroComp.Compare.Result.Errors)})
				macroTable = append(macroTable, []string{"Reconnects", convertFloatToString(macroComp.Reference.Result.Reconnects), convertFloatToString(macroComp.Compare.Result.Reconnects)})
				macroTable = append(macroTable, []string{"Time", strconv.Itoa(macroComp.Reference.Result.Time), strconv.Itoa(macroComp.Compare.Result.Time)})
				macroTable = append(macroTable, []string{"Threads", convertFloatToString(macroComp.Reference.Result.Threads), convertFloatToString(macroComp.Compare.Result.Threads)})
			}
		}
		writeTableToPdf(pdf, macroTable)
	}

	if len(microsMatrix) != 0 {
		writeSubtitle(pdf, "Micro-benchmarks")
		microTable := [][]string{
			{"Metric", git.ShortenSHA(fromSHA), git.ShortenSHA(toSHA)},
		}
		for _, microComp := range microsMatrix {
			microTable = append(microTable, []string{microComp.PkgName + "." + microComp.Name})
			microTable = append(microTable, []string{"Ops", strconv.Itoa(microComp.Current.Ops), strconv.Itoa(microComp.Last.Ops)})
			microTable = append(microTable, []string{"NSPerOp", convertFloatToString(microComp.Current.NSPerOp), convertFloatToString(microComp.Last.NSPerOp)})
			microTable = append(microTable, []string{"MBPerSec", convertFloatToString(microComp.Current.MBPerSec), convertFloatToString(microComp.Last.MBPerSec)})
			microTable = append(microTable, []string{"BytesPerOp", convertFloatToString(microComp.Current.BytesPerOp), convertFloatToString(microComp.Last.BytesPerOp)})
			microTable = append(microTable, []string{"AllocsPerOp", convertFloatToString(microComp.Current.AllocsPerOp), convertFloatToString(microComp.Last.AllocsPerOp)})
			microTable = append(microTable, []string{"Ratio of NSPerOp", convertFloatToString(microComp.CurrLastDiff)})
		}
		writeTableToPdf(pdf, microTable)
	}
	err = pdf.OutputFileAndClose(reportFile)
	return err
}

func convertFloatToString(in float64) string {
	return strconv.FormatFloat(in, 'f', -1, 64)
}

func writeSubtitle(pdf *gofpdf.Fpdf, subtitle string) {
	pdf.SetFont("Arial", "B", 20)
	pdf.WriteAligned(0, 10, subtitle, "C")
	pdf.Ln(-1)
	pdf.Ln(2)
}

func writeTableToPdf(pdf *gofpdf.Fpdf, table [][]string) {
	pageWidth, pageHeight := pdf.GetPageSize()
	colCount := len(table[0])
	colWd := (pageWidth - pdfMargin - pdfMargin) / float64(colCount)
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

	header := table[0]

	pdf.SetFont("Arial", "", 14)

	// Headers
	pdf.SetTextColor(224, 224, 224)
	pdf.SetFillColor(64, 64, 64)
	for colJ := 0; colJ < colCount; colJ++ {
		pdf.CellFormat(colWd, 10, header[colJ], "1", 0, "CM", true, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetTextColor(24, 24, 24)
	pdf.SetFillColor(255, 255, 255)

	// Rows
	y := pdf.GetY()
	for rowJ := 1; rowJ < len(table); rowJ++ {
		cellList = nil
		colCount = len(table[rowJ])
		colWd = 180.0 / float64(colCount)
		maxHt := lineHt
		// Cell height calculation loop
		for colJ := 0; colJ < colCount; colJ++ {
			cell.str = table[rowJ][colJ]
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
			for splitJ := 0; splitJ < len(cell.list); splitJ++ {
				pdf.SetXY(x+cellGap, cellY)
				pdf.CellFormat(colWd-cellGap-cellGap, lineHt, string(cell.list[splitJ]), "", 0,
					"C", false, 0, "")
				cellY += lineHt
			}
			x += colWd
		}
		y += maxHt + cellGap + cellGap
	}
	pdf.Ln(-1)
	pdf.Ln(4)
}
