package services

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

const fullLen = 50

func FancyPrint(str string) {
	var leftPadding = (fullLen - len(str)) / 2
	printChar("*", fullLen, true)
	printChar(" ", leftPadding, false)
	fmt.Println(str)
	printChar("*", fullLen, true)
}

func printChar(ch string, count int, newLine bool) {
	for i := 0; i < count; i++ {
		fmt.Print(ch)
	}

	if newLine {
		fmt.Println()
	}
}

const cellLen = 30

func PrintLine(colCount int) {
	tableWidth := (cellLen * colCount) + (colCount + 1)
	printChar("-", tableWidth, true)
}

func PrintRow(rowData []string) {
	fmt.Print("|")

	for i := 0; i < len(rowData); i++ {
		leftPadding := (cellLen - len(rowData[i])) / 2
		rightPadding := cellLen - len(rowData[i]) - leftPadding

		printChar(" ", leftPadding, false)
		fmt.Print(rowData[i])
		printChar(" ", rightPadding, false)
		fmt.Print("|")
	}

	fmt.Println()
}

func ClearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func Pause() {
	fmt.Println("Press Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
