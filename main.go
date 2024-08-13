package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	mainMenu()
}

func clearScreen() {
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

func pause() {
	fmt.Println("Press Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func mainMenu() {
	clearScreen()
	fancyPrint("Welcome to String Processor")
	fmt.Println("Please select an option:")
	fmt.Println("  1. Count Word Frequency")
	fmt.Println("  2. Check Palindrome")
	fmt.Println("  3. Exit")

	fmt.Print("Enter your choice: ")
	choice := readLine()

	switch choice {
	case "1":
		countWordMenu()
	case "2":
		checkPalindromeMenu()
	case "3":
		fmt.Println("Exiting program...")
		os.Exit(0)
	default:
		fmt.Println("Invalid choice. Please try again.")
		pause()
		mainMenu()
	}
}

func countWordMenu() {
	clearScreen()
	fancyPrint("Welcome to Word Frequency Counter")
	fmt.Print("Please enter a string to count word frequency: ")

	str := readLine()
	for word, count := range CountWords(str) {
		fmt.Printf("Counts of %s: %d\n", word, count)
	}

	fmt.Print("Do you want to try again? (Y/N): ")
	choice := readLine()
	if strings.ToLower(choice) == "y" {
		countWordMenu()
	} else {
		mainMenu()
	}
}

func CountWords(str string) map[string]int {
	str = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r == ' ' {
			return r
		}
		return -1
	}, str)
	words := strings.Fields(str)

	counts := make(map[string]int)
	for _, word := range words {
		counts[strings.ToLower(word)]++
	}

	return counts
}

func checkPalindromeMenu() {
	clearScreen()
	fancyPrint("Welcome to Palindrome Checker")
	fmt.Print("Please enter a string to check palindrome: ")

	str := readLine()
	if IsPalindrome(str) {
		fmt.Println("The string is a palindrome.")
	} else {
		fmt.Println("The string is not a palindrome.")
	}

	fmt.Print("Do you want to try again? (Y/N): ")
	choice := readLine()
	if strings.ToLower(choice) == "y" {
		checkPalindromeMenu()
	} else {
		mainMenu()
	}
}

func IsPalindrome(str string) bool {
	str = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' {
			return r
		}
		return -1
	}, str)

	str = strings.ToLower(str)
	left, right := 0, len(str)-1

	for left < right {
		if str[left] != str[right] {
			return false
		}

		left++
		right--
	}

	return true
}

const cellLen = 20
const fullLen = (cellLen * 3) + 4

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

func printChar(ch string, count int, newLine bool) {
	for i := 0; i < count; i++ {
		fmt.Print(ch)
	}

	if newLine {
		fmt.Println()
	}
}
