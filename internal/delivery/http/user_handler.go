package http

import (
	"encoding/json"
	"net/http"
	"regexp"

	"blog-api/internal/entity"
	"blog-api/internal/usecase"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	UserUsecase usecase.UserUsecase
}

func NewUserHandler(r *mux.Router, userUsecase usecase.UserUsecase) {
	handler := &UserHandler{
		UserUsecase: userUsecase,
	}

	r.HandleFunc("/register", handler.Register).Methods("POST")
	r.HandleFunc("/login", handler.Login).Methods("POST")
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if user.Username == "" || user.Password == "" || user.Email == "" {
		http.Error(w, "Username, Password, and Email are required", http.StatusBadRequest)
		return
	}

	// Validate Email format
	if !isValidEmail(user.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Check if the user already exists
	existingUser, err := h.UserUsecase.GetByUsernameOrEmail(user.Username, user.Email)
	if err != nil {
		http.Error(w, "Failed to check existing user", http.StatusInternalServerError)
		return
	}
	if existingUser != nil {
		http.Error(w, "Username or Email already exists", http.StatusConflict)
		return
	}

	// Check if the user role is provided
	if user.Role == "" {
		// Check if this is the first attempt by looking for a specific cookie
		cookie, err := r.Cookie("first_attempt")
		if err != nil || cookie.Value != "true" {
			// Set a cookie to mark the first attempt
			http.SetCookie(w, &http.Cookie{
				Name:  "first_attempt",
				Value: "true",
				Path:  "/",
				// Optionally set an expiration time
				MaxAge: 30,
			})
			http.Error(w, "Please specify your role (either 'author' or 'user').", http.StatusBadRequest)
			return
		} else {
			// If it's a subsequent attempt, set the default role to 'user'
			user.Role = "user"
		}
	}

	// Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Proceed with registration by calling the use case
	if err := h.UserUsecase.Register(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Clear the first attempt cookie after successful registration
	http.SetCookie(w, &http.Cookie{
		Name:   "first_attempt",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Deletes the cookie
	})

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully with role: " + user.Role,
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loggedInUser, err := h.UserUsecase.Login(user.Username, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userDetails, err := h.UserUsecase.GetByUsernameOrEmail(user.Username, user.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Return both user details and the JWT toke
	response := map[string]interface{}{
		"token": loggedInUser,
		"user": map[string]interface{}{
			"user_id":  userDetails.ID,
			"username": userDetails.Username,
			"email":    userDetails.Email,
			"role":     userDetails.Role,
		},
	}

	// // Generate the JWT token
	// token, err := h.UserUsecase.GenerateJWTToken(loggedInUser.ID, loggedInUser.Role)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusUnauthorized)
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
