package controllers

import (
	"bufio"
	"errors"
	"fmt"
	"lib_mngmt/models"
	"lib_mngmt/services"
	"os"
	"strconv"
	"strings"
)

type LibraryController struct {
	service services.LibraryManager
}

func NewLibraryController(service services.LibraryManager) *LibraryController {
	return &LibraryController{service: service}
}

func (lc *LibraryController) Run() {
	reader := bufio.NewReader(os.Stdin)

	for {
		services.ClearScreen()
		services.FancyPrint("Library Management System")
		fmt.Println("Please select an option:")
		fmt.Println("  1. Add a new book")
		fmt.Println("  2. Remove an existing book")
		fmt.Println("  3. Add a new member")
		fmt.Println("  4. Remove an existing member")
		fmt.Println("  5. Borrow a book")
		fmt.Println("  6. Return a book")
		fmt.Println("  7. List all available books")
		fmt.Println("  8. List all borrowed books by a member")
		fmt.Println("  9. Exit")
		fmt.Print("Enter your choice: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			lc.addBook(reader)
		case "2":
			lc.removeBook(reader)
		case "3":
			lc.addMember(reader)
		case "4":
			lc.removeMember(reader)
		case "5":
			lc.borrowBook(reader)
		case "6":
			lc.returnBook(reader)
		case "7":
			lc.listAvailableBooks()
		case "8":
			lc.listBorrowedBooks(reader)
		case "9":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}

		services.Pause()
	}
}

func (lc *LibraryController) addBook(reader *bufio.Reader) {
	services.ClearScreen()
	services.FancyPrint("Add a new book")

	idStr, id, err := "", 0, errors.New("")
	for err != nil {
		fmt.Print("Enter book ID: ")
		idStr, err = reader.ReadString('\n')

		if err != nil {
			fmt.Println("Invalid ID. Please try again.")
			continue
		}

		id, err = strconv.Atoi(strings.TrimSpace(idStr))

		if err != nil {
			fmt.Println("Invalid ID. Please try again.")
		}
	}

	title := ""
	err = nil
	for title == "" || err != nil {
		fmt.Print("Enter book title: ")
		title, err = reader.ReadString('\n')
		title = strings.TrimSpace(title)

		if err != nil {
			fmt.Println("Invalid title. Please try again.")
		} else if title == "" {
			fmt.Println("Title cannot be empty. Please try again.")
		}
	}

	author := ""
	err = nil
	for author == "" || err != nil {
		fmt.Print("Enter book author: ")
		author, err = reader.ReadString('\n')
		author = strings.TrimSpace(author)

		if err != nil {
			fmt.Println("Invalid author. Please try again.")
		} else if author == "" {
			fmt.Println("Author cannot be empty. Please try again.")
		}
	}

	book := models.Book{
		ID:     id,
		Title:  title,
		Author: author,
		Status: "Available",
	}

	err = lc.service.AddBook(book)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Book added successfully.")
	}
}

func (lc *LibraryController) removeBook(reader *bufio.Reader) {
	services.ClearScreen()
	services.FancyPrint("Remove an existing book")

	idStr, id, err := "", 0, errors.New("")
	for err != nil {
		fmt.Print("Enter book ID to remove: ")

		idStr, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Invalid ID. Please try again.")
			continue
		}

		id, err = strconv.Atoi(strings.TrimSpace(idStr))
		if err != nil {
			fmt.Println("Invalid ID. Please try again.")
		}
	}

	err = lc.service.RemoveBook(id)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Book removed successfully.")
	}
}

func (lc *LibraryController) addMember(reader *bufio.Reader) {
	services.ClearScreen()
	services.FancyPrint("Add a new member")

	idStr, id, err := "", 0, errors.New("")
	for err != nil {
		fmt.Print("Enter member ID: ")

		idStr, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Invalid ID. Please try again.")
			continue
		}

		id, err = strconv.Atoi(strings.TrimSpace(idStr))
		if err != nil {
			fmt.Println("Invalid ID. Please try again.")
		}
	}

	name := ""
	err = nil
	for name == "" || err != nil {
		fmt.Print("Enter member name: ")
		name, err = reader.ReadString('\n')
		name = strings.TrimSpace(name)

		if err != nil {
			fmt.Println("Invalid name. Please try again.")
		} else if name == "" {
			fmt.Println("Name cannot be empty. Please try again.")
		}
	}

	member := models.Member{
		ID:   id,
		Name: name,
	}

	err = lc.service.AddMember(member)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Member added successfully.")
	}
}

func (lc *LibraryController) removeMember(reader *bufio.Reader) {
	services.ClearScreen()
	services.FancyPrint("Remove an existing member")

	idStr, id, err := "", 0, errors.New("")
	for err != nil {
		fmt.Print("Enter member ID to remove: ")

		idStr, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Invalid ID. Please try again.")
			continue
		}

		id, err = strconv.Atoi(strings.TrimSpace(idStr))
		if err != nil {
			fmt.Println("Invalid ID. Please try again.")
		}
	}

	err = lc.service.RemoveMember(id)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Member removed successfully.")
	}
}

func (lc *LibraryController) borrowBook(reader *bufio.Reader) {
	services.ClearScreen()
	services.FancyPrint("Borrow a book")

	bookID, memberID := 0, 0
	bookIDStr, memberIDStr := "", ""
	err := errors.New("")

	for err != nil {
		fmt.Print("Enter book ID to borrow: ")

		bookIDStr, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Invalid ID. Please try again.")
			continue
		}

		bookID, err = strconv.Atoi(strings.TrimSpace(bookIDStr))
		if err != nil {
			fmt.Println("Invalid ID. Please try again.")
		}
	}

	err = errors.New("")
	for err != nil {
		fmt.Print("Enter member ID: ")
		memberIDStr, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Invalid ID. Please try again.")
			continue
		}

		memberID, err = strconv.Atoi(strings.TrimSpace(memberIDStr))
		if err != nil {
			fmt.Println("Invalid ID. Please try again.")
		}
	}

	err = lc.service.BorrowBook(bookID, memberID)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Book borrowed successfully.")
	}
}

func (lc *LibraryController) returnBook(reader *bufio.Reader) {
	services.ClearScreen()
	services.FancyPrint("Return a book")

	bookID, memberID := 0, 0
	bookIDStr, memberIDStr := "", ""
	err := errors.New("")

	for err != nil {
		fmt.Print("Enter book ID to return: ")
		bookIDStr, err = reader.ReadString('\n')

		if err != nil {
			fmt.Println("Invalid ID. Please try again.")
			continue
		}

		bookID, err = strconv.Atoi(strings.TrimSpace(bookIDStr))
		if err != nil {
			fmt.Println("Invalid ID. Please try again.")
		}
	}

	err = errors.New("")
	for err != nil {
		fmt.Print("Enter member ID: ")
		memberIDStr, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Invalid ID. Please try again.")
			continue
		}

		memberID, err = strconv.Atoi(strings.TrimSpace(memberIDStr))
		if err != nil {
			fmt.Println("Invalid ID. Please try again.")
		}
	}

	err = lc.service.ReturnBook(bookID, memberID)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Book returned successfully.")
	}
}

func (lc *LibraryController) listAvailableBooks() {
	books := lc.service.ListAvailableBooks()
	services.ClearScreen()
	if len(books) == 0 {
		fmt.Println("No books available.")
		return
	}

	services.PrintLine(3)
	services.PrintRow([]string{"ID", "Title", "Author"})
	services.PrintLine(3)
	for _, book := range books {
		services.PrintRow([]string{strconv.Itoa(book.ID), book.Title, book.Author})
	}
	services.PrintLine(3)
}

func (lc *LibraryController) listBorrowedBooks(reader *bufio.Reader) {
	services.ClearScreen()
	services.FancyPrint("List all borrowed books by a member")

	memberIDStr, memberID, err := "", 0, errors.New("")
	for err != nil {
		fmt.Print("Enter member ID: ")
		memberIDStr, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Invalid ID. Please try again.")
			continue
		}

		memberID, err = strconv.Atoi(strings.TrimSpace(memberIDStr))
		if err != nil {
			fmt.Println("Invalid ID. Please try again.")
		}
	}

	var books []models.Book
	books, err = lc.service.ListBorrowedBooks(memberID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(books) == 0 {
		fmt.Println("No borrowed books.")
		return
	}

	fmt.Println("Borrowed books:")
	services.PrintLine(3)
	services.PrintRow([]string{"ID", "Title", "Author"})
	services.PrintLine(3)
	for _, book := range books {
		services.PrintRow([]string{strconv.Itoa(book.ID), book.Title, book.Author})
	}
	services.PrintLine(3)
}
