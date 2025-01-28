package main

import (
	"testing"
)

const (
	inputFile  = "./testdata/test1.md"
	goldenFile = "./testdata/test1.md.html"
)

func TestPackItems(t *testing.T) {
	totalItems := PackItems(2000)
	var expectedTotalItems int32
	expectedTotalItems = 2000

	if totalItems != expectedTotalItems {
		t.Errorf("Expected %d but got %d", expectedTotalItems, totalItems)
	}
}

// // func TestParseContent(t *testing.T) {
// // 	input, err := os.ReadFile(inputFile)
// // 	if err != nil {
// // 		t.Fatal(err)
// // 	}
// // 	result, err := parseContent(input, "")
// // 	if err != nil {
// // 		t.Fatal(err)
// // 	}
// // 	expected, err := os.ReadFile(goldenFile)
// // 	if err != nil {
// // 		t.Fatal(err)
// // 	}
// // 	if !bytes.Equal(result, expected) {
// // 		t.Logf("golden: \n%s\n", expected)
// // 		t.Logf("result: \n%s\n", result)
// // 		t.Error("Result content does not match golden file")
// // 	}
// // }

// func TestRun(t *testing.T) {
// 	var mockStdOut bytes.Buffer
// 	if err := run(inputFile, "", &mockStdOut, true); err != nil {
// 		t.Fatal(err)
// 	}
// 	resultFile := strings.TrimSpace(mockStdOut.String())
// 	result, err := os.ReadFile(resultFile)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	expected, err := os.ReadFile(goldenFile)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if !bytes.Equal(expected, result) {
// 		t.Logf("Expected: \n%s\n", expected)
// 		t.Logf("Result: \n%s\n", result)
// 		t.Logf("Rsult does not match golden file")
// 	}
// 	os.Remove(resultFile)
// }
