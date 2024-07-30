package services

import (
	"errors"
	"lib_mngmt/models"
)

type LibraryManager interface {
	AddBook(book models.Book) error
	RemoveBook(bookID int) error
	AddMember(member models.Member) error
	RemoveMember(memberID int) error
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) ([]models.Book, error)
}

type Library struct {
	books   map[int]models.Book
	members map[int]models.Member
}

func NewLibrary() *Library {
	return &Library{
		books:   make(map[int]models.Book),
		members: make(map[int]models.Member),
	}
}

func (l *Library) AddBook(book models.Book) error {
	_, bookExists := l.books[book.ID]
	if bookExists {
		return errors.New("a book with the given id already exists")
	}

	l.books[book.ID] = book
	return nil
}

func (l *Library) RemoveBook(bookID int) error {
	_, bookExists := l.books[bookID]
	if !bookExists {
		return errors.New("book not found")
	}

	delete(l.books, bookID)
	return nil
}

func (l *Library) AddMember(member models.Member) error {
	_, memberExists := l.members[member.ID]
	if memberExists {
		return errors.New("a member with the given id already exists")
	}

	l.members[member.ID] = member
	return nil
}

func (l *Library) RemoveMember(memberID int) error {
	_, memberExists := l.members[memberID]
	if !memberExists {
		return errors.New("member not found")
	}

	delete(l.members, memberID)
	return nil
}

func (l *Library) BorrowBook(bookID int, memberID int) error {
	book, bookExists := l.books[bookID]
	if !bookExists {
		return errors.New("book not found")
	}

	if book.Status == "Borrowed" {
		return errors.New("book already borrowed")
	}

	member, memberExists := l.members[memberID]
	if !memberExists {
		return errors.New("member not found")
	}

	book.Status = "Borrowed"
	l.books[bookID] = book

	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.members[memberID] = member

	return nil
}

func (l *Library) ReturnBook(bookID int, memberID int) error {
	member, memberExists := l.members[memberID]
	if !memberExists {
		return errors.New("member not found")
	}

	book, bookExists := l.books[bookID]
	if !bookExists {
		return errors.New("book not found")
	}

	if book.Status == "Available" {
		return errors.New("book is not borrowed")
	}

	book.Status = "Available"
	l.books[bookID] = book

	for i, b := range member.BorrowedBooks {
		if b.ID == bookID {
			member.BorrowedBooks = append(member.BorrowedBooks[:i], member.BorrowedBooks[i+1:]...)
			break
		}
	}

	l.members[memberID] = member

	return nil
}

func (l *Library) ListAvailableBooks() []models.Book {
	availableBooks := []models.Book{}
	for _, book := range l.books {
		if book.Status == "Available" {
			availableBooks = append(availableBooks, book)
		}
	}
	return availableBooks
}

func (l *Library) ListBorrowedBooks(memberID int) ([]models.Book, error) {
	member, memberExists := l.members[memberID]
	if !memberExists {
		return nil, errors.New("member not found")
	}

	return member.BorrowedBooks, nil
}
