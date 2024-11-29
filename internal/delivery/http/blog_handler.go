package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"blog-api/config"
	"blog-api/internal/entity"
	"blog-api/internal/usecase"
	"blog-api/pkg/jwt"
	"blog-api/pkg/middleware"

	"github.com/gorilla/mux"
)

type BlogHandler struct {
	BlogUsecase usecase.BlogUsecase
}

func NewBlogHandler(r *mux.Router, blogUsecase usecase.BlogUsecase, secretKey string) {
	handler := &BlogHandler{
		BlogUsecase: blogUsecase,
	}

	// User can read all blogs and post cmment
	r.HandleFunc("/blogs", handler.GetAllBlogs).Methods("GET")
	r.HandleFunc("/blogs/{id}", handler.GetBlogByID).Methods("GET")
	r.Handle("/comments/{blogID}", middleware.AuthMiddleware(secretKey, http.HandlerFunc(handler.CreateComment))).Methods("POST")

	// Author can create, update and delete blogs
	r.Handle("/blogs", middleware.AuthorMiddleware(secretKey)(http.HandlerFunc(handler.CreateBlog))).Methods("POST")
	r.Handle("/blogs/{id}", middleware.AuthorMiddleware(secretKey)(http.HandlerFunc(handler.UpdateBlog))).Methods("PUT")
	r.Handle("/blogs/{id}", middleware.AuthorMiddleware(secretKey)(http.HandlerFunc(handler.DeleteBlog))).Methods("DELETE")

}

func (h *BlogHandler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	// Parse the form data to handle both fields and files
	err := r.ParseMultipartForm(10 << 20) // Limit file size to 10MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Extract text fields (title, content, etc.)
	title := r.FormValue("title")
	content := r.FormValue("content")

	// Extract the file (thumbnail)
	thumbnailFile, _, err := r.FormFile("thumbnail")
	if err != nil {
		http.Error(w, "Thumbnail image is required", http.StatusBadRequest)
		return
	}
	defer thumbnailFile.Close()

	// Validate that the title and content are not empty
	if title == "" || content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	// Assuming userID is extracted from the JWT token or context
	userID, ok := r.Context().Value(jwt.UserIDKey).(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	// Ensure the "uploads" directory exists
	err = os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		log.Println("Error creating uploads directory:", err)
		http.Error(w, "Unble to create directory for thumbnail", http.StatusInternalServerError)
		return
	}

	// Generate a file name for the thumbnail and save it
	thumbnailPath := "uploads/" + "thumbnail.jpg" // Store the thumbnail in the "uploads" folder
	outFile, err := os.Create(thumbnailPath)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to save thumbnail", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	// Copy the uploaded thumbnail to the file
	_, err = io.Copy(outFile, thumbnailFile)
	if err != nil {
		http.Error(w, "Unable to save thumbnail", http.StatusInternalServerError)
		return
	}

	// Create the blog entity
	blog := &entity.Blog{
		Title:     title,
		Content:   content,
		UserID:    userID,
		Thumbnail: thumbnailPath,
	}

	// Log the user ID to ensure it is correct
	log.Printf("Creating blog with user ID: %d", userID)

	// Save the blog in the database
	if err := h.BlogUsecase.Create(blog); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusCreated)
}

func (h *BlogHandler) GetAllBlogs(w http.ResponseWriter, r *http.Request) {
	blogs, err := h.BlogUsecase.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(blogs)
}

func (h *BlogHandler) GetBlogByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	blog, err := h.BlogUsecase.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(blog)
}

func (h *BlogHandler) UpdateBlog(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	secretKey := config.LoadConfig().JWTSecret
	// Extract user ID from JWT token
	claims, err := jwt.ExtractClaims(r, secretKey)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := claims.UserID

	// Retrieve the existing blog from the database
	existingBlog, err := h.BlogUsecase.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the user is the author of the blog
	if existingBlog.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the form data to handle both fields and files
	err = r.ParseMultipartForm(10 << 20) // Limit file size to 10MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Extract text fields (title, content, etc.)
	title := r.FormValue("title")
	content := r.FormValue("content")

	// Validate that the title and content are not empty
	if title == "" || content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	// Extract the file (thumbnail) if provided
	thumbnailFile, _, err := r.FormFile("thumbnail")
	var thumbnailPath string
	if err == nil {
		defer thumbnailFile.Close()

		// Ensure the "uploads" directory exists
		err = os.MkdirAll("uploads", os.ModePerm)
		if err != nil {
			log.Println("Error creating uploads directory:", err)
			http.Error(w, "Unable to create directory for thumbnail", http.StatusInternalServerError)
			return
		}

		// Generate a file name for the thumbnail and save it
		thumbnailPath = "uploads/" + "thumbnail.jpg" // Store the thumbnail in the "uploads" folder
		outFile, err := os.Create(thumbnailPath)
		if err != nil {
			log.Println(err)
			http.Error(w, "Unable to save thumbnail", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		// Copy the uploaded thumbnail to the file
		_, err = io.Copy(outFile, thumbnailFile)
		if err != nil {
			http.Error(w, "Unable to save thumbnail", http.StatusInternalServerError)
			return
		}
	}

	// Update the blog entity
	existingBlog.Title = title
	existingBlog.Content = content
	if thumbnailPath != "" {
		existingBlog.Thumbnail = thumbnailPath
	}

	// Save the updated blog in the database
	if err := h.BlogUsecase.Update(existingBlog); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusNoContent)
}

func (h *BlogHandler) DeleteBlog(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Delete blog using usecase
	if err := h.BlogUsecase.Delete(id); err != nil {
		// If blog does not exist, respond with 404 Not Found
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Log success
	fmt.Println("Successfully deleted blog with ID:", id)

	w.WriteHeader(http.StatusNoContent)
}

func (h *BlogHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment entity.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	blogID, err := strconv.Atoi(mux.Vars(r)["blogID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check that the blog exists
	_, err = h.BlogUsecase.GetByID(blogID)
	if err != nil {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}

	comment.BlogID = blogID

	// Assuming user ID is extracted from context or JWT token
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	comment.UserID = userID

	if err := h.BlogUsecase.CreateComment(&comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
