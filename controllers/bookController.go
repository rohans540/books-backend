package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rohans540/books-backend/database"
	"github.com/rohans540/books-backend/kafka"
	"github.com/rohans540/books-backend/models"
	"github.com/rohans540/books-backend/redis"
)

// GetBooks godoc
// @Summary Get all books with pagination
// @Description Retrieve paginated details of all books
// @Tags books
// @Produce json
// @Param limit query int false "Limit the number of books per page (default: 10)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {array} models.Book
// @Router /books [get]
func GetBooks(ctx *gin.Context) {

	cachedBooks, err := redis.RedisClient.Get(context.Background(), "books").Result()
	if err == nil && cachedBooks != "" {
		ctx.Header("Content-Type", "application/json")
		ctx.String(http.StatusOK, cachedBooks)
		return
	}

	var books []models.Book
	result := database.DB.Find(&books)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching books"})
		return
	}

	booksJSON, _ := json.Marshal(books)
	redis.RedisClient.Set(context.Background(), "books", booksJSON, 0)
	ctx.JSON(http.StatusOK, books)
}

// GetBookByID godoc
// @Summary Get book by ID
// @Description Retrieve details of a book by its ID
// @Tags books
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} models.Book
// @Failure 404 {object} map[string]string "Book not found"
// @Router /books/{id} [get]
func GetBookByID(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	id := ctx.Param("id")
	var book models.Book

	cachedBook, err := redis.RedisClient.Get(context.Background(), "book:"+id).Result()
	if err == nil {
		json.Unmarshal([]byte(cachedBook), &book)
		ctx.JSON(http.StatusOK, book)
		return
	}

	result := database.DB.First(&book, id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	data, _ := json.Marshal(book)
	redis.RedisClient.Set(context.Background(), "book:"+id, data, 0)
	ctx.JSON(http.StatusOK, book)
}

// CreateBook godoc
// @Summary Create a new book
// @Description Add a new book to the collection
// @Tags books
// @Accept json
// @Produce json
// @Param book body models.Book true "Book object"
// @Success 201 {object} models.Book
// @Failure 400 {object} map[string]string "Invalid request body"
// @Router /books [post]
func CreateBook(ctx *gin.Context) {
	var book models.Book
	if err := ctx.ShouldBindJSON(&book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	if book.Title == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Title cannot be empty"})
		return
	}
	if book.Author == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Author cannot be empty"})
		return
	}
	if book.Year <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Year must be a valid positive number"})
		return
	}

	result := database.DB.Create(&book)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	redis.RedisClient.Del(context.Background(), "books")
	kafka.PublishMessage("book_events", "New book added: "+book.Title)

	ctx.JSON(http.StatusCreated, book)
}

// UpdateBook godoc
// @Summary Update an existing book
// @Description Modify the details of an existing book
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body models.Book true "Updated book object"
// @Success 200 {object} models.Book
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 404 {object} map[string]string "Book not found"
// @Router /books/{id} [put]
func UpdateBook(ctx *gin.Context) {
	id := ctx.Param("id")
	var book models.Book
	result := database.DB.First(&book, id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	var updatedBook models.Book
	if err := ctx.ShouldBindJSON(&updatedBook); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	if updatedBook.Title == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Title cannot be empty"})
		return
	}
	if updatedBook.Author == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Author cannot be empty"})
		return
	}
	if updatedBook.Year <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Year must be a valid positive number"})
		return
	}
	book.Title = updatedBook.Title
	book.Author = updatedBook.Author
	book.Year = updatedBook.Year
	database.DB.Save(&book)

	deletedKeys := []string{"books", "book:" + id}
	redis.RedisClient.Del(context.Background(), deletedKeys...)

	val, _ := redis.RedisClient.Get(context.Background(), "books").Result()
	log.Println("Redis books cache after update:", val)

	kafka.PublishMessage("book_events", "Book updated: "+book.Title)

	ctx.JSON(http.StatusOK, book)
}

// DeleteBook godoc
// @Summary Delete a book
// @Description Remove a book from the collection
// @Tags books
// @Param id path int true "Book ID"
// @Success 200 {object} map[string]string "Book deleted successfully"
// @Failure 404 {object} map[string]string "Book not found"
// @Router /books/{id} [delete]
func DeleteBook(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	id := ctx.Param("id")
	var book models.Book

	result := database.DB.First(&book, id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	database.DB.Unscoped().Delete(&book)
	deletedKeys := []string{"books", "book:" + id}
	redis.RedisClient.Del(context.Background(), deletedKeys...)
	val, _ := redis.RedisClient.Get(context.Background(), "books").Result()
	log.Println("Redis books cache after delete:", val)

	kafka.PublishMessage("book_events", "Book deleted: "+strconv.Itoa(int(book.ID)))

	ctx.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}
