package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fancyPrint("Welcome to Student Grade Calculator")

	name := readString("Please enter your name: ", "Please enter a valid name!")
	fmt.Printf("Hello, %s! ", name)
	subjectCount := readSubjectCount()

	var subjectNames []string
	var grades []float64
	for i := 0; i < subjectCount; i++ {
		subjectPrompt := fmt.Sprintf("Enter name of subject %d: ", i+1)
		subjectName := readString(subjectPrompt, "Please enter a valid subject name!")
		subjectNames = append(subjectNames, subjectName)

		grade := readGrade(i + 1)
		grades = append(grades, grade)
	}

	printReport(name, subjectNames, grades)
}

const cellLen = 20
const fullLen = (cellLen * 3) + 4

func readString(prompt string, errorPrompt string) string {
	fmt.Print(prompt)
	str := readLine()

	_, err := strconv.ParseFloat(str, 64)
	for err == nil || str == "" {
		fmt.Println(errorPrompt)
		fmt.Print(prompt)
		str = readLine()
		_, err = strconv.ParseFloat(str, 64)
	}

	return str
}

func readSubjectCount() int {
	prompt := "Please enter how many subjects you took: "
	fmt.Print(prompt)
	str := readLine()

	count, err := strconv.Atoi(str)
	for err != nil || count <= 0 {
		fmt.Println("Please enter a valid number of subjects.")
		fmt.Print(prompt)
		str = readLine()
		count, err = strconv.Atoi(str)
	}

	return count
}

func readGrade(index int) float64 {
	prompt := fmt.Sprintf("Enter grade for subject %d (should be between 0 and 100): ", index)
	fmt.Print(prompt)
	str := readLine()

	grade, err := strconv.ParseFloat(str, 64)
	for err != nil || grade < 0 || grade > 100 {
		fmt.Println("Please enter a valid grade between 0 and 100.")
		fmt.Print(prompt)
		str = readLine()
		grade, err = strconv.ParseFloat(str, 64)
	}

	return grade
}

func fancyPrint(str string) {
	var leftPadding = (fullLen - len(str)) / 2
	printChar("*", fullLen, true)
	printChar(" ", leftPadding, false)
	fmt.Println(str)
	printChar("*", fullLen, true)
}

func readLine() string {
	var str string

	reader := bufio.NewReader(os.Stdin)
	str, _ = reader.ReadString('\n')

	return strings.TrimSpace(str)
}

func calculateAverage(grades []float64) float64 {
	var total float64
	for i := 0; i < len(grades); i++ {
		total += grades[i]
	}

	return total / float64(len(grades))
}

func printReport(studentName string, subjectNames []string, grades []float64) {
	average := calculateAverage(grades)

	fancyPrint("Student Grade Report")
	fmt.Println("\n    Student Name:", studentName, "\n")

	printTable(subjectNames, grades)

	fmt.Printf("\n    Average Grade: %.2f\n\n", average)
}

func printTable(subjectNames []string, grades []float64) {
	subjectCount := len(subjectNames)

	printChar("-", (cellLen)*3+4, true)
	printRow([]string{"ID", "Subject Name", "Grade"})
	printChar("-", (cellLen)*3+4, true)
	for i := 0; i < subjectCount; i++ {
		rowData := []string{
			fmt.Sprintf("%d", i+1),
			subjectNames[i],
			fmt.Sprintf("%.2f", grades[i]),
		}

		printRow(rowData)
	}
	printChar("-", (cellLen)*3+4, true)
}

func printRow(rowData []string) {
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

func printChar(ch string, count int, newLine bool) {
	for i := 0; i < count; i++ {
		fmt.Print(ch)
	}

	if newLine {
		fmt.Println()
	}
}
