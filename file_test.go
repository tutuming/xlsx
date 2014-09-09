package xlsx

import (
	. "gopkg.in/check.v1"
)

type FileSuite struct {}

var _ = Suite(&FileSuite{})

// Test we can correctly open a XSLX file and return a xlsx.File
// struct.
func (l *FileSuite) TestOpenFile(c *C) {
	var xlsxFile *File
	var error error

	xlsxFile, error = OpenFile("testfile.xlsx")
	c.Assert(error, IsNil)
	c.Assert(xlsxFile, NotNil)
}

// Test we can create a File object from scratch
func (l *FileSuite) TestCreateFile(c *C) {
	var xlsxFile *File

	xlsxFile = NewFile()
	c.Assert(xlsxFile, NotNil)
}

// Test that when we open a real XLSX file we create xlsx.Sheet
// objects for the sheets inside the file and that these sheets are
// themselves correct.
func (l *FileSuite) TestCreateSheet(c *C) {
	var xlsxFile *File
	var err error
	var sheet *Sheet
	var row *Row
	xlsxFile, err = OpenFile("testfile.xlsx")
	c.Assert(err, IsNil)
	c.Assert(xlsxFile, NotNil)
	sheetLen := len(xlsxFile.Sheets)
	c.Assert(sheetLen, Equals, 3)
	sheet = xlsxFile.Sheets["Tabelle1"]
	rowLen := len(sheet.Rows)
	c.Assert(rowLen, Equals, 2)
	row = sheet.Rows[0]
	c.Assert(len(row.Cells), Equals, 2)
	cell := row.Cells[0]
	cellstring := cell.String()
	c.Assert(cellstring, Equals, "Foo")
}

// Test that we can add a sheet to a File
func (l *FileSuite) TestAddSheet(c *C) {
	var f *File
	f = NewFile()
	sheet := f.AddSheet("MySheet")
	c.Assert(sheet, NotNil)
	c.Assert(len(f.Sheets), Equals, 1)
	c.Assert(f.Sheets["MySheet"], Equals, sheet)
}

// Test that we can marshall a File to a collection of xml files
func (l *FileSuite) TestMarshalFile(c *C) {
	var f *File
	f = NewFile()
	sheet1 := f.AddSheet("MySheet")
	row1 := sheet1.AddRow()
	cell1 := row1.AddCell()
	cell1.Value = "A cell!"
	sheet2 := f.AddSheet("AnotherSheet")
	row2 := sheet2.AddRow()
	cell2 := row2.AddCell()
	cell2.Value = "A cell!"
	parts, err := f.MarshallParts()
	c.Assert(err, IsNil)
	c.Assert(len(parts), Equals, 4)
	expectedSheet := `<?xml version="1.0" encoding="UTF-8"?>
  <worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">
    <dimension ref="A1:A1"></dimension>
    <sheetData>
      <row r="1">
        <c r="A1" t="s">
          <v>0</v>
        </c>
      </row>
    </sheetData>
  </worksheet>`
	c.Assert(parts["worksheets/sheet1.xml"], Equals, expectedSheet)
	c.Assert(parts["worksheets/sheet2.xml"], Equals, expectedSheet)
	expectedXLSXSST := `<?xml version="1.0" encoding="UTF-8"?>
  <sst count="1" uniqueCount="1">
    <si>
      <t>A cell!</t>
    </si>
  </sst>`
	c.Assert(parts["xl/sharedStrings.xml"], Equals, expectedXLSXSST)
	expectedXLSXWorkbookRels := `<?xml version="1.0" encoding="UTF-8"?>
  <Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
    <Relationship Id="rId1" Target="worksheets/sheet1.xml" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/sharedStrings"></Relationship>
    <Relationship Id="rId2" Target="worksheets/sheet2.xml" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/sharedStrings"></Relationship>
  </Relationships>`
	c.Assert(parts["xl/_rels/workbook.xml.rels"], Equals, expectedXLSXWorkbookRels)

}