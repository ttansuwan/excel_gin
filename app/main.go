package main

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type BindFile struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type InvalidRow struct {
	rowNo  int
	email  string
	reason string
}

func validateProductMeta(product []string) (int, error) {
	// Slice the array
	metaData := product[9:]
	prevCol := metaData[0]

	for i, currCol := range metaData[1:] {
		// Check if prev is empty and curr is not
		if (prevCol == "") && (currCol != "") {
			return i, errors.New("missing information for current column")
		}
		// Shift to the next column
		prevCol = currCol
	}
	return -1, nil
}

func validateHeader(header []string) {
	validHeader := []string{"Email", "First name", "Last name", "Position", "Brand", "Company", "Local Supermarket", "MOM&POP", "Distributor", "Division", "Category", "Segment", "Manufacturer", "Brand", "Campaign"}

	// Check length and content
	if !reflect.DeepEqual(header, validHeader) {
		panic("Invalid header")
	}
}

func main() {
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.Static("/", "./public")

	router.POST("/upload", func(c *gin.Context) {
		var bindFile BindFile

		// Bind file
		if err := c.ShouldBind(&bindFile); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("err: %s", err.Error()))
			return
		}

		// Open excel sheet
		file := bindFile.File
		readFile, err := file.Open()
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("err: %s", err.Error()))
			return
		}

		f, err := excelize.OpenReader(readFile)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("err: %s", err.Error()))
			return
		}
		defer f.Close()

		// Get workbook first sheet
		firstSheet := f.WorkBook.Sheets.Sheet[0].Name
		rows, err := f.GetRows(firstSheet)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer func() {
			c.String(http.StatusForbidden, "err: Invalid file")
		}()

		var invalidRows []InvalidRow
		for i, row := range rows {
			if i == 0 {
				validateHeader(row)
			}
			colNo, err := validateProductMeta(row)

			if err != nil {
				invalidRows = append(invalidRows, InvalidRow{rowNo: i, email: row[0], reason: "Invalid data"})
			}
			fmt.Print(i, row)
			fmt.Println()
		}

		c.String(http.StatusOK, fmt.Sprintf("File %s uploaded.", file.Filename))
	})
	router.Run(":8080")
}
