package usecase

import (
	"blog-api/internal/entity"
	repoMysql "blog-api/internal/repository/mysql"
	"fmt"
)

type BlogUsecase interface {
	Create(blog *entity.Blog) error
	GetAll() ([]*entity.Blog, error)
	GetByID(id int) (*entity.Blog, error)
	Update(blog *entity.Blog) error
	Delete(id int) error
	CreateComment(comment *entity.Comment) error
}

type blogUsecase struct {
	blogRepo *repoMysql.BlogRepository
}

// CreateComment implements BlogUsecase.
func (u *blogUsecase) CreateComment(comment *entity.Comment) error {
	return u.blogRepo.CreateComment(comment)
}

func NewBlogUsecase(blogRepo *repoMysql.BlogRepository) BlogUsecase {
	return &blogUsecase{blogRepo: blogRepo}
}

func (u *blogUsecase) Create(blog *entity.Blog) error {
	return u.blogRepo.Create(blog)
}

func (u *blogUsecase) GetAll() ([]*entity.Blog, error) {
	return u.blogRepo.GetAll()
}

func (u *blogUsecase) GetByID(id int) (*entity.Blog, error) {
	return u.blogRepo.GetByID(id)
}

func (u *blogUsecase) Update(blog *entity.Blog) error {
	return u.blogRepo.Update(blog)
}

func (u *blogUsecase) Delete(id int) error {
	// Check if the blog exists before attempting to delete it
	_, err := u.blogRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("blog with ID %d not found", id)
	}

	// Proceed with the deletion
	return u.blogRepo.Delete(id)
}
