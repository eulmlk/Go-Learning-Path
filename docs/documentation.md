# Library Management System User Guide

## Overview

Welcome to the Library Management System! This console-based application allows you to manage books and members in a library. You can add and remove books, borrow and return books, and view available and borrowed books. This guide will help you understand how to use the system effectively.

## Getting Started

1. **Run the Application**

   To start the application, navigate to the project directory in your terminal and run:

   ```sh
   go run main.go
   ```

2. **Main Menu**

   Once the application starts, you will see the main menu with the following options:

   ```
   **************************************************
               Library Management System
   **************************************************
   Please select an option:
   1. Add a new book
   2. Remove an existing book
   3. Add a new member
   4. Remove an existing member
   5. Borrow a book
   6. Return a book
   7. List all available books
   8. List all borrowed books by a member
   9. Exit
   Enter your choice: 
   ```

## Functionalities

### 1. Add a New Book

   To add a new book to the library, select option `1` from the main menu. You will be prompted to enter the book ID, title, and author.

   Example:
   ```
   **************************************************
                  Add a new book
   **************************************************
   Enter book ID: 1
   Enter book title: Go Programming
   Enter book author: John Doe
   Book added successfully.
   Press Enter to continue...
   ```

### 2. Remove an Existing Book

   To remove a book from the library, select option `2` from the main menu. You will be prompted to enter the book ID of the book you wish to remove.

   Example:
   ```
   **************************************************
               Remove an existing book
   **************************************************
   Enter book ID to remove: 1
   Book removed successfully.
   Press Enter to continue...
   ```

### 3. Add a new Member

   To add a new member to the library, select option `3` from the main menu. You will be prompted to enter the member ID and name.

   Example:
   ```
   **************************************************
                  Add a new member
   **************************************************
   Enter member ID: 1
   Enter member name: Abebe Kebede
   Member added successfully.
   Press Enter to continue...
   ```

### 4. Remove an Existing Member

   To remove a member from the library, select option `4` from the main menu. You will be prompted to enter the member ID of the member you wish to remove.

   Example:
   ```
   **************************************************
               Remove an existing member
   **************************************************
   Enter member ID to remove: 1
   Member removed successfully.
   Press Enter to continue...
   ```

### 5. Borrow a Book

   To borrow a book, select option `5` from the main menu. You will be prompted to enter the book ID and the member ID.

   Example:
   ```
   **************************************************
                     Borrow a book
   **************************************************
   Enter book ID to borrow: 1
   Enter member ID: 1
   Book borrowed successfully.
   Press Enter to continue...
   ```

   **Note:** The book must be available to be borrowed, and the member must be registered.

### 6. Return a Book

   To return a borrowed book, select option `6` from the main menu. You will be prompted to enter the book ID and the member ID.

   Example:
   ```
   **************************************************
                     Return a book
   **************************************************
   Enter book ID to return: 1
   Enter member ID: 1
   Book returned successfully.
   Press Enter to continue...
   ```
   **Note:** The book must be borrowed by the specified member to be returned.

### 7. List All Available Books

   To list all available books in the library, select option `7` from the main menu. The system will display all books that are currently available for borrowing.

   Example:
   ```
   ----------------------------------------------------------------
   |         ID         |       Title        |       Author       |
   ----------------------------------------------------------------
   |         1          |   Go Programming   |      John Doe      |
   |         2          |   C Programming    |      Dohn Joe      |
   |         3          |        OOP         |      Jim Jack      |
   |         4          |  Data Structures   |     Mike Jack      |
   ----------------------------------------------------------------
   Press Enter to continue...
   ```

### 8. List All Borrowed Books by a Member

   To list all books borrowed by a specific member, select option `8` from the main menu. You will be prompted to enter the member ID.

   Example:
   ```
   **************************************************
         List all borrowed books by a member
   **************************************************
   Enter member ID: 1
   Borrowed books:
   ----------------------------------------------------------------
   |         ID         |       Title        |       Author       |
   ----------------------------------------------------------------
   |         1          |   Go Programming   |      John Doe      |
   ----------------------------------------------------------------
   Press Enter to continue...
   ```

### 9. Exit

   To exit the application, select option `9` from the main menu. The application will terminate, and you will return to the terminal.