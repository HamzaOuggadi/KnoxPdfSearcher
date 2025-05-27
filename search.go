package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ledongthuc/pdf"
)

func SearchPdfForName(pdfPath, name string) (bool, string) {
	f, r, err := pdf.Open(pdfPath)
	if err != nil {
		fmt.Println("Failed to open PDF:", err)
		return false, ""
	}

	defer f.Close()

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		fmt.Println("Failed to read PDF: ", err)
		return false, ""
	}

	buf.ReadFrom(b)

	lines := strings.Split((buf.String()), "\n")

	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), strings.ToLower(name)) {
			return true, line
		}
	}

	return false, ""
}
