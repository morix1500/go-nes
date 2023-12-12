package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	nestestFile, err := os.Open("../nestest.log")
	if err != nil {
		fmt.Println("Error opening nestest.log:", err)
		return
	}
	defer nestestFile.Close()

	resFile, err := os.Open("../res.log")
	if err != nil {
		fmt.Println("Error opening res.log:", err)
		return
	}
	defer resFile.Close()

	nestestScanner := bufio.NewScanner(nestestFile)
	resScanner := bufio.NewScanner(resFile)

	lineNumber := 1
	for nestestScanner.Scan() && resScanner.Scan() {
		nestestLine := strings.TrimSpace(nestestScanner.Text())
		resLine := strings.TrimSpace(resScanner.Text())

		nestestValues := strings.Split(nestestLine, " ")
		resValues := strings.Split(resLine, " ")

		// Remove empty strings from nestestValues
		var nestestValuesFiltered []string
		for _, value := range nestestValues {
			if value != "" {
				nestestValuesFiltered = append(nestestValuesFiltered, value)
			}
		}
		// Remove empty strings from nestestValues
		var resValuesFiltered []string
		for _, value := range resValues {
			if value != "" {
				resValuesFiltered = append(resValuesFiltered, value)
			}
		}

		for i := 0; i < len(resValuesFiltered); i++ {
			if nestestValuesFiltered[i] != resValuesFiltered[i] {
				fmt.Printf("Difference at line %d:\n", lineNumber)
				fmt.Printf("Diff: %s != %s\n", nestestValuesFiltered[i], resValuesFiltered[i])
				fmt.Println("nestest.log:", nestestLine)
				fmt.Println("res.log:", resLine)
				fmt.Println()
				break
			}
		}

		lineNumber++
	}

	if nestestScanner.Err() != nil {
		fmt.Println("Error reading nestest.log:", nestestScanner.Err())
	}

	if resScanner.Err() != nil {
		fmt.Println("Error reading res.log:", resScanner.Err())
	}
}
