package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Errors

var BookNotFound = errors.New("Book not found")

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "title1", Author: "auther1", Quantity: 3},
	{ID: "2", Title: "title2", Author: "auther2", Quantity: 4},
	{ID: "3", Title: "title3", Author: "auther3", Quantity: 9},
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func createBook(c *gin.Context) {
	var newbook book
	if err := c.BindJSON(&newbook); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	books = append(books, newbook)
	c.IndentedJSON(http.StatusCreated, newbook)
}

func bookById(c *gin.Context) {
	id := c.Param("id")

	book, err := getBookbyId(id)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		log.Println("BookbyID err:", err)
		return
	}

	c.IndentedJSON(http.StatusAccepted, book)
}

// Implement book checkout function

func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameters."})
		return
	}

	book, err := getBookbyId(id)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not found."})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available."})
		return
	}
	book.Quantity -= 1
	updateBook(book)
	c.IndentedJSON(http.StatusOK, book)
}

func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameters."})
		return
	}

	book, err := getBookbyId(id)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not found."})
		return
	}

	book.Quantity += 1
	updateBook(book)
	c.IndentedJSON(http.StatusOK, book)
}

func getBookbyId(id string) (*book, error) {
	for _, book := range books {
		if book.ID == id {
			return &book, nil
		}
	}
	return nil, BookNotFound
}

func updateBook(b *book) {
	for i, v := range books {
		if v.ID == b.ID {
			books[i] = *b
		}
	}
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.POST("/books", createBook)
	router.GET("/books/:id", bookById)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.Run(":8080")
}
