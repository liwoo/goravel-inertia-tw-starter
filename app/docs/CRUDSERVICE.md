# Implementing CrudAppService Pattern in Goravel

A comprehensive guide to implementing ABP.io-style CrudAppService pattern in Goravel with automated validation, permissions, DTOs, and Swagger documentation.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Implementation Guide](#implementation-guide)
    - [1. Models and Relationships](#1-models-and-relationships)
    - [2. DTOs (Request/Response Objects)](#2-dtos-requestresponse-objects)
    - [3. CrudAppService Interface](#3-crudappservice-interface)
    - [4. Validation Layer](#4-validation-layer)
    - [5. Authorization with Gates](#5-authorization-with-gates)
    - [6. Resource Controllers](#6-resource-controllers)
    - [7. Swagger Documentation](#7-swagger-documentation)
- [Code Generation](#code-generation)
- [Usage Examples](#usage-examples)
- [Best Practices](#best-practices)
- [Conclusion](#conclusion)

## Overview

This implementation brings the powerful CrudAppService pattern from ABP.io to Goravel, providing:

- **Standardized CRUD operations** with pagination, sorting, and filtering
- **Automatic validation** using Goravel's validation system
- **Dynamic permissions** using Goravel Gates
- **Type-safe DTOs** for requests and responses
- **Auto-generated Swagger documentation**
- **Artisan commands** for scaffolding

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Controller    │────│   AppService    │────│   Repository    │
│  (HTTP Layer)   │    │ (Business Logic)│    │  (Data Access)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   DTOs/Forms    │    │   Validation    │    │     Models      │
│   Swagger       │    │     Gates       │    │  Relationships  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Prerequisites

- Goravel >= 1.14
- Go >= 1.21
- Swaggo for API documentation

```bash
# Install Swaggo
go install github.com/swaggo/swag/cmd/swag@latest
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

## Implementation Guide

### 1. Models and Relationships

First, let's create our example models with proper relationships:

**User Model (`app/models/user.go`)**:
```go
package models

import (
    "github.com/goravel/framework/database/orm"
    "time"
)

type User struct {
    orm.Model
    Name     string `gorm:"not null" json:"name"`
    Email    string `gorm:"unique;not null" json:"email"`
    Posts    []Post `gorm:"foreignKey:UserID" json:"posts,omitempty"`
    orm.SoftDeletes
}

func (u *User) TableName() string {
    return "users"
}
```

**Post Model (`app/models/post.go`)**:
```go
package models

import (
    "github.com/goravel/framework/database/orm"
    "time"
)

type Post struct {
    orm.Model
    Title       string `gorm:"not null" json:"title"`
    Content     string `gorm:"type:text" json:"content"`
    Published   bool   `gorm:"default:false" json:"published"`
    UserID      uint   `gorm:"not null" json:"user_id"`
    User        User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
    orm.SoftDeletes
}

func (p *Post) TableName() string {
    return "posts"
}
```

**Migration (`database/migrations/create_posts_table.go`)**:
```go
package migrations

import (
    "github.com/goravel/framework/contracts/database/schema"
    "github.com/goravel/framework/facades"
)

type CreatePostsTable struct{}

func (receiver *CreatePostsTable) Up() error {
    return facades.Schema().Create("posts", func(table schema.Blueprint) {
        table.ID()
        table.String("title")
        table.Text("content")
        table.Boolean("published").Default(false)
        table.UnsignedBigInteger("user_id")
        table.Timestamps()
        table.SoftDeletes()
        
        table.Foreign("user_id").References("id").On("users")
        table.Index("user_id")
    })
}

func (receiver *CreatePostsTable) Down() error {
    return facades.Schema().DropIfExists("posts")
}
```

### 2. DTOs (Request/Response Objects)

**Base DTO Types (`app/dtos/base.go`)**:
```go
package dtos

import (
    "time"
)

// PagedRequest provides pagination and sorting
type PagedRequest struct {
    Page     int    `form:"page" json:"page" default:"1" minimum:"1"`
    PageSize int    `form:"page_size" json:"page_size" default:"10" minimum:"1" maximum:"100"`
    Sort     string `form:"sort" json:"sort" example:"id desc"`
    Filter   string `form:"filter" json:"filter" example:"search term"`
}

// PagedResponse provides paginated results
type PagedResponse[T any] struct {
    Data       []T `json:"data"`
    TotalCount int `json:"total_count"`
    Page       int `json:"page"`
    PageSize   int `json:"page_size"`
    TotalPages int `json:"total_pages"`
}

// BaseEntityDTO provides common fields
type BaseEntityDTO struct {
    ID        uint       `json:"id"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
```

**Post DTOs (`app/dtos/post_dto.go`)**:
```go
package dtos

// PostDTO for reading operations
type PostDTO struct {
    BaseEntityDTO
    Title     string  `json:"title"`
    Content   string  `json:"content"`
    Published bool    `json:"published"`
    UserID    uint    `json:"user_id"`
    User      *UserDTO `json:"user,omitempty"`
}

// CreatePostDTO for creation
type CreatePostDTO struct {
    Title     string `json:"title" validate:"required,min:3,max:255" example:"My First Post"`
    Content   string `json:"content" validate:"required,min:10" example:"This is the content of my post"`
    Published bool   `json:"published" example:"false"`
    UserID    uint   `json:"user_id" validate:"required,min:1" example:"1"`
}

// UpdatePostDTO for updates
type UpdatePostDTO struct {
    Title     *string `json:"title,omitempty" validate:"omitempty,min:3,max:255"`
    Content   *string `json:"content,omitempty" validate:"omitempty,min:10"`
    Published *bool   `json:"published,omitempty"`
}

// PostFilterDTO for filtering
type PostFilterDTO struct {
    PagedRequest
    Title     *string `form:"title" json:"title,omitempty"`
    Published *bool   `form:"published" json:"published,omitempty"`
    UserID    *uint   `form:"user_id" json:"user_id,omitempty"`
}

// UserDTO for reading operations
type UserDTO struct {
    BaseEntityDTO
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

### 3. CrudAppService Interface

**Base CrudAppService Interface (`app/contracts/crud_app_service.go`)**:
```go
package contracts

import (
    "context"
    "goravel/app/dtos"
)

// CrudAppService defines the contract for CRUD operations
type CrudAppService[TEntity any, TEntityDTO any, TKey any, TGetListInput any, TCreateInput any, TUpdateInput any] interface {
    GetAsync(ctx context.Context, id TKey) (*TEntityDTO, error)
    GetListAsync(ctx context.Context, input TGetListInput) (*dtos.PagedResponse[TEntityDTO], error)
    CreateAsync(ctx context.Context, input TCreateInput) (*TEntityDTO, error)
    UpdateAsync(ctx context.Context, id TKey, input TUpdateInput) (*TEntityDTO, error)
    DeleteAsync(ctx context.Context, id TKey) error
}

// IPostAppService specific interface
type IPostAppService interface {
    CrudAppService[models.Post, dtos.PostDTO, uint, dtos.PostFilterDTO, dtos.CreatePostDTO, dtos.UpdatePostDTO]
}
```

**Base CrudAppService Implementation (`app/services/base_crud_app_service.go`)**:
```go
package services

import (
    "context"
    "errors"
    "fmt"
    "goravel/app/dtos"
    "math"
    "reflect"
    "strings"
    
    "github.com/goravel/framework/contracts/database/orm"
    "github.com/goravel/framework/facades"
)

// BaseCrudAppService provides base implementation
type BaseCrudAppService[TEntity any, TEntityDTO any, TKey any, TGetListInput any, TCreateInput any, TUpdateInput any] struct {
    Repository orm.Query
    EntityName string
}

func NewBaseCrudAppService[TEntity any, TEntityDTO any, TKey any, TGetListInput any, TCreateInput any, TUpdateInput any](
    entityName string,
) *BaseCrudAppService[TEntity, TEntityDTO, TKey, TGetListInput, TCreateInput, TUpdateInput] {
    return &BaseCrudAppService[TEntity, TEntityDTO, TKey, TGetListInput, TCreateInput, TUpdateInput]{
        Repository: facades.Orm().Query(),
        EntityName: entityName,
    }
}

func (s *BaseCrudAppService[TEntity, TEntityDTO, TKey, TGetListInput, TCreateInput, TUpdateInput]) GetAsync(ctx context.Context, id TKey) (*TEntityDTO, error) {
    var entity TEntity
    err := s.Repository.Find(&entity, id)
    if err != nil {
        return nil, fmt.Errorf("entity not found: %w", err)
    }
    
    dto, err := s.MapToEntityDTO(entity)
    if err != nil {
        return nil, fmt.Errorf("mapping error: %w", err)
    }
    
    return dto, nil
}

func (s *BaseCrudAppService[TEntity, TEntityDTO, TKey, TGetListInput, TCreateInput, TUpdateInput]) GetListAsync(ctx context.Context, input TGetListInput) (*dtos.PagedResponse[TEntityDTO], error) {
    query := s.Repository
    
    // Apply filtering if input implements filtering interface
    if filterable, ok := any(input).(FilterableInput); ok {
        query = s.ApplyFilters(query, filterable)
    }
    
    // Apply sorting if input implements sorting interface
    if sortable, ok := any(input).(SortableInput); ok {
        query = s.ApplySorting(query, sortable)
    }
    
    // Get pagination info
    var page, pageSize int = 1, 10
    if pageable, ok := any(input).(PageableInput); ok {
        page = pageable.GetPage()
        pageSize = pageable.GetPageSize()
    }
    
    // Get total count
    var totalCount int64
    countQuery := query
    countQuery.Count(&totalCount)
    
    // Apply pagination
    offset := (page - 1) * pageSize
    query = query.Offset(offset).Limit(pageSize)
    
    // Execute query
    var entities []TEntity
    err := query.Get(&entities)
    if err != nil {
        return nil, fmt.Errorf("query execution error: %w", err)
    }
    
    // Map to DTOs
    var dtos []TEntityDTO
    for _, entity := range entities {
        dto, err := s.MapToEntityDTO(entity)
        if err != nil {
            return nil, fmt.Errorf("mapping error: %w", err)
        }
        dtos = append(dtos, *dto)
    }
    
    totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
    
    return &dtos.PagedResponse[TEntityDTO]{
        Data:       dtos,
        TotalCount: int(totalCount),
        Page:       page,
        PageSize:   pageSize,
        TotalPages: totalPages,
    }, nil
}

func (s *BaseCrudAppService[TEntity, TEntityDTO, TKey, TGetListInput, TCreateInput, TUpdateInput]) CreateAsync(ctx context.Context, input TCreateInput) (*TEntityDTO, error) {
    entity, err := s.MapToEntity(input)
    if err != nil {
        return nil, fmt.Errorf("mapping error: %w", err)
    }
    
    err = s.Repository.Create(entity)
    if err != nil {
        return nil, fmt.Errorf("creation error: %w", err)
    }
    
    dto, err := s.MapToEntityDTO(*entity)
    if err != nil {
        return nil, fmt.Errorf("mapping error: %w", err)
    }
    
    return dto, nil
}

func (s *BaseCrudAppService[TEntity, TEntityDTO, TKey, TGetListInput, TCreateInput, TUpdateInput]) UpdateAsync(ctx context.Context, id TKey, input TUpdateInput) (*TEntityDTO, error) {
    var entity TEntity
    err := s.Repository.Find(&entity, id)
    if err != nil {
        return nil, fmt.Errorf("entity not found: %w", err)
    }
    
    // Apply updates using reflection for partial updates
    s.ApplyUpdates(&entity, input)
    
    err = s.Repository.Save(&entity)
    if err != nil {
        return nil, fmt.Errorf("update error: %w", err)
    }
    
    dto, err := s.MapToEntityDTO(entity)
    if err != nil {
        return nil, fmt.Errorf("mapping error: %w", err)
    }
    
    return dto, nil
}

func (s *BaseCrudAppService[TEntity, TEntityDTO, TKey, TGetListInput, TCreateInput, TUpdateInput]) DeleteAsync(ctx context.Context, id TKey) error {
    var entity TEntity
    err := s.Repository.Find(&entity, id)
    if err != nil {
        return fmt.Errorf("entity not found: %w", err)
    }
    
    err = s.Repository.Delete(&entity)
    if err != nil {
        return fmt.Errorf("deletion error: %w", err)
    }
    
    return nil
}

// Helper interfaces for type assertions
type FilterableInput interface {
    GetFilters() map[string]interface{}
}

type SortableInput interface {
    GetSort() string
}

type PageableInput interface {
    GetPage() int
    GetPageSize() int
}

// Mapping methods (to be overridden by concrete implementations)
func (s *BaseCrudAppService[TEntity, TEntityDTO, TKey, TGetListInput, TCreateInput, TUpdateInput]) MapToEntity(input TCreateInput) (*TEntity, error) {
    // This should be overridden in concrete implementations
    // For now, return a generic mapping error
    return nil, errors.New("MapToEntity not implemented")
}

func (s *BaseCrudAppService[TEntity, TEntityDTO, TKey, TGetListInput, TCreateInput, TUpdateInput]) MapToEntityDTO(entity TEntity) (*TEntityDTO, error) {
    // This should be overridden in concrete implementations
    return nil, errors.New("MapToEntityDTO not implemented")
}

func (s *BaseCrudAppService[TEntity, TEntityDTO, TKey, TGetListInput, TCreateInput, TUpdateInput]) ApplyFilters(query orm.Query, input FilterableInput) orm.Query {
    filters := input.GetFilters()
    for field, value := range filters {
        if value != nil {
            query = query.Where(field+" = ?", value)
        }
    }
    return query
}

func (s *BaseCrudAppService[TEntity, TEntityDTO, TKey, TGetListInput, TCreateInput, TUpdateInput]) ApplySorting(query orm.Query, input SortableInput) orm.Query {
    sort := input.GetSort()
    if sort != "" {
        query = query.Order(sort)
    }
    return query
}

func (s *BaseCrudAppService[TEntity, TEntityDTO, TKey, TGetListInput, TCreateInput, TUpdateInput]) ApplyUpdates(entity *TEntity, input TUpdateInput) {
    entityValue := reflect.ValueOf(entity).Elem()
    inputValue := reflect.ValueOf(input)
    
    for i := 0; i < inputValue.NumField(); i++ {
        field := inputValue.Field(i)
        if field.Kind() == reflect.Ptr && !field.IsNil() {
            fieldName := inputValue.Type().Field(i).Name
            entityField := entityValue.FieldByName(fieldName)
            if entityField.IsValid() && entityField.CanSet() {
                entityField.Set(field.Elem())
            }
        }
    }
}
```

### 4. Validation Layer

**Form Request for Posts (`app/http/requests/post_request.go`)**:
```go
package requests

import (
    "goravel/app/dtos"
    
    "github.com/goravel/framework/contracts/http"
    "github.com/goravel/framework/contracts/validation"
    "github.com/goravel/framework/facades"
)

type CreatePostRequest struct {
    dtos.CreatePostDTO
}

func (r *CreatePostRequest) Authorize(ctx http.Context) error {
    // Check if user can create posts
    if !facades.Gate().WithContext(ctx).Allows("create-post", map[string]any{}) {
        return errors.New("unauthorized to create posts")
    }
    return nil
}

func (r *CreatePostRequest) Rules(ctx http.Context) map[string]string {
    return map[string]string{
        "title":     "required|min_len:3|max_len:255",
        "content":   "required|min_len:10",
        "published": "bool",
        "user_id":   "required|numeric|min:1",
    }
}

func (r *CreatePostRequest) Messages(ctx http.Context) map[string]string {
    return map[string]string{
        "title.required":    "Post title is required",
        "title.min_len":     "Post title must be at least 3 characters",
        "title.max_len":     "Post title cannot exceed 255 characters",
        "content.required":  "Post content is required",
        "content.min_len":   "Post content must be at least 10 characters",
        "user_id.required":  "User ID is required",
        "user_id.numeric":   "User ID must be a number",
        "user_id.min":       "User ID must be greater than 0",
    }
}

func (r *CreatePostRequest) Attributes(ctx http.Context) map[string]string {
    return map[string]string{
        "user_id": "User",
    }
}

type UpdatePostRequest struct {
    dtos.UpdatePostDTO
}

func (r *UpdatePostRequest) Authorize(ctx http.Context) error {
    postID := ctx.Request().Route("id")
    // Check if user can update this specific post
    if !facades.Gate().WithContext(ctx).Allows("update-post", map[string]any{
        "post_id": postID,
    }) {
        return errors.New("unauthorized to update this post")
    }
    return nil
}

func (r *UpdatePostRequest) Rules(ctx http.Context) map[string]string {
    return map[string]string{
        "title":   "min_len:3|max_len:255",
        "content": "min_len:10",
    }
}
```

### 5. Authorization with Gates

**Auth Service Provider (`app/providers/auth_service_provider.go`)**:
```go
package providers

import (
    "context"
    "goravel/app/models"
    
    "github.com/goravel/framework/auth/access"
    "github.com/goravel/framework/contracts/auth/access"
    "github.com/goravel/framework/contracts/foundation"
    "github.com/goravel/framework/facades"
)

type AuthServiceProvider struct{}

func (receiver *AuthServiceProvider) Register(app foundation.Application) {
    // Register services if needed
}

func (receiver *AuthServiceProvider) Boot(app foundation.Application) {
    // Define Gates for Post operations
    facades.Gate().Define("create-post", func(ctx context.Context, arguments map[string]any) contractsaccess.Response {
        user := ctx.Value("user")
        if user == nil {
            return access.NewDenyResponse("User not authenticated")
        }
        
        // Example: Any authenticated user can create posts
        return access.NewAllowResponse()
    })
    
    facades.Gate().Define("update-post", func(ctx context.Context, arguments map[string]any) contractsaccess.Response {
        user := ctx.Value("user")
        if user == nil {
            return access.NewDenyResponse("User not authenticated")
        }
        
        postID, exists := arguments["post_id"]
        if !exists {
            return access.NewDenyResponse("Post ID not provided")
        }
        
        // Get the post to check ownership
        var post models.Post
        err := facades.Orm().Query().Find(&post, postID)
        if err != nil {
            return access.NewDenyResponse("Post not found")
        }
        
        currentUser := user.(models.User)
        if post.UserID == currentUser.ID {
            return access.NewAllowResponse()
        }
        
        // Check if user is admin (example additional logic)
        if currentUser.Email == "admin@example.com" {
            return access.NewAllowResponse()
        }
        
        return access.NewDenyResponse("You can only update your own posts")
    })
    
    facades.Gate().Define("delete-post", func(ctx context.Context, arguments map[string]any) contractsaccess.Response {
        // Similar logic to update-post
        user := ctx.Value("user")
        if user == nil {
            return access.NewDenyResponse("User not authenticated")
        }
        
        postID, exists := arguments["post_id"]
        if !exists {
            return access.NewDenyResponse("Post ID not provided")
        }
        
        var post models.Post
        err := facades.Orm().Query().Find(&post, postID)
        if err != nil {
            return access.NewDenyResponse("Post not found")
        }
        
        currentUser := user.(models.User)
        if post.UserID == currentUser.ID {
            return access.NewAllowResponse()
        }
        
        return access.NewDenyResponse("You can only delete your own posts")
    })
    
    facades.Gate().Define("view-post", func(ctx context.Context, arguments map[string]any) contractsaccess.Response {
        // Example: Anyone can view published posts, owners can view their unpublished posts
        post, exists := arguments["post"]
        if !exists {
            return access.NewDenyResponse("Post not provided")
        }
        
        postModel := post.(models.Post)
        if postModel.Published {
            return access.NewAllowResponse()
        }
        
        user := ctx.Value("user")
        if user == nil {
            return access.NewDenyResponse("Must be authenticated to view unpublished posts")
        }
        
        currentUser := user.(models.User)
        if postModel.UserID == currentUser.ID {
            return access.NewAllowResponse()
        }
        
        return access.NewDenyResponse("Can only view your own unpublished posts")
    })
}
```

### 6. Resource Controllers

**Post Controller (`app/http/controllers/post_controller.go`)**:
```go
package controllers

import (
    "goravel/app/contracts"
    "goravel/app/http/requests"
    "strconv"
    
    "github.com/goravel/framework/contracts/http"
    "github.com/goravel/framework/facades"
)

type PostController struct {
    postAppService contracts.IPostAppService
}

func NewPostController(postAppService contracts.IPostAppService) *PostController {
    return &PostController{
        postAppService: postAppService,
    }
}

// GetList godoc
// @Summary      List posts
// @Description  Get paginated list of posts with filtering and sorting
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        page       query   int     false  "Page number" default(1)
// @Param        page_size  query   int     false  "Items per page" default(10)
// @Param        sort       query   string  false  "Sort order" default("id desc")
// @Param        filter     query   string  false  "Search filter"
// @Param        title      query   string  false  "Filter by title"
// @Param        published  query   bool    false  "Filter by published status"
// @Param        user_id    query   int     false  "Filter by user ID"
// @Success      200  {object}  dtos.PagedResponse[dtos.PostDTO]
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/posts [get]
func (r *PostController) Index(ctx http.Context) http.Response {
    var filter dtos.PostFilterDTO
    errors, err := ctx.Request().ValidateRequest(&filter)
    if err != nil {
        return ctx.Response().Json(http.StatusBadRequest, map[string]interface{}{
            "error": "Validation failed",
            "details": errors,
        })
    }
    
    result, err := r.postAppService.GetListAsync(ctx.Context(), filter)
    if err != nil {
        facades.Log().Error("Failed to get posts: " + err.Error())
        return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
            "error": "Failed to retrieve posts",
        })
    }
    
    return ctx.Response().Json(http.StatusOK, result)
}

// Show godoc
// @Summary      Get post by ID
// @Description  Get single post by ID
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Post ID"
// @Success      200  {object}  dtos.PostDTO
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/posts/{id} [get]
func (r *PostController) Show(ctx http.Context) http.Response {
    idStr := ctx.Request().Route("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        return ctx.Response().Json(http.StatusBadRequest, map[string]string{
            "error": "Invalid post ID",
        })
    }
    
    // Check view permission
    if !facades.Gate().WithContext(ctx).Allows("view-post", map[string]any{
        "post_id": uint(id),
    }) {
        return ctx.Response().Json(http.StatusForbidden, map[string]string{
            "error": "Unauthorized to view this post",
        })
    }
    
    result, err := r.postAppService.GetAsync(ctx.Context(), uint(id))
    if err != nil {
        return ctx.Response().Json(http.StatusNotFound, map[string]string{
            "error": "Post not found",
        })
    }
    
    return ctx.Response().Json(http.StatusOK, result)
}

// Store godoc
// @Summary      Create new post
// @Description  Create a new post
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        request  body      dtos.CreatePostDTO  true  "Post data"
// @Success      201      {object}  dtos.PostDTO
// @Failure      400      {object}  map[string]interface{}
// @Failure      403      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/posts [post]
// @Security     BearerAuth
func (r *PostController) Store(ctx http.Context) http.Response {
    var request requests.CreatePostRequest
    errors, err := ctx.Request().ValidateRequest(&request)
    if err != nil {
        return ctx.Response().Json(http.StatusBadRequest, map[string]interface{}{
            "error": "Validation failed",
            "details": errors,
        })
    }
    
    result, err := r.postAppService.CreateAsync(ctx.Context(), request.CreatePostDTO)
    if err != nil {
        facades.Log().Error("Failed to create post: " + err.Error())
        return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
            "error": "Failed to create post",
        })
    }
    
    return ctx.Response().Json(http.StatusCreated, result)
}

// Update godoc
// @Summary      Update post
// @Description  Update an existing post
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        id       path      int                 true  "Post ID"
// @Param        request  body      dtos.UpdatePostDTO  true  "Post update data"
// @Success      200      {object}  dtos.PostDTO
// @Failure      400      {object}  map[string]interface{}
// @Failure      403      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/posts/{id} [put]
// @Security     BearerAuth
func (r *PostController) Update(ctx http.Context) http.Response {
    idStr := ctx.Request().Route("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        return ctx.Response().Json(http.StatusBadRequest, map[string]string{
            "error": "Invalid post ID",
        })
    }
    
    var request requests.UpdatePostRequest
    errors, err := ctx.Request().ValidateRequest(&request)
    if err != nil {
        return ctx.Response().Json(http.StatusBadRequest, map[string]interface{}{
            "error": "Validation failed",
            "details": errors,
        })
    }
    
    result, err := r.postAppService.UpdateAsync(ctx.Context(), uint(id), request.UpdatePostDTO)
    if err != nil {
        facades.Log().Error("Failed to update post: " + err.Error())
        return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
            "error": "Failed to update post",
        })
    }
    
    return ctx.Response().Json(http.StatusOK, result)
}

// Destroy godoc
// @Summary      Delete post
// @Description  Delete a post by ID
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Post ID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/posts/{id} [delete]
// @Security     BearerAuth
func (r *PostController) Destroy(ctx http.Context) http.Response {
    idStr := ctx.Request().Route("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        return ctx.Response().Json(http.StatusBadRequest, map[string]string{
            "error": "Invalid post ID",
        })
    }
    
    // Check delete permission
    if !facades.Gate().WithContext(ctx).Allows("delete-post", map[string]any{
        "post_id": uint(id),
    }) {
        return ctx.Response().Json(http.StatusForbidden, map[string]string{
            "error": "Unauthorized to delete this post",
        })
    }
    
    err = r.postAppService.DeleteAsync(ctx.Context(), uint(id))
    if err != nil {
        facades.Log().Error("Failed to delete post: " + err.Error())
        return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
            "error": "Failed to delete post",
        })
    }
    
    return ctx.Response().Json(http.StatusNoContent, nil)
}
```

**Post App Service Implementation (`app/services/post_app_service.go`)**:
```go
package services

import (
    "context"
    "goravel/app/contracts"
    "goravel/app/dtos"
    "goravel/app/models"
    
    "github.com/goravel/framework/facades"
)

type PostAppService struct {
    *BaseCrudAppService[models.Post, dtos.PostDTO, uint, dtos.PostFilterDTO, dtos.CreatePostDTO, dtos.UpdatePostDTO]
}

func NewPostAppService() contracts.IPostAppService {
    base := NewBaseCrudAppService[models.Post, dtos.PostDTO, uint, dtos.PostFilterDTO, dtos.CreatePostDTO, dtos.UpdatePostDTO]("Post")
    return &PostAppService{
        BaseCrudAppService: base,
    }
}

// Override GetListAsync to implement custom filtering
func (s *PostAppService) GetListAsync(ctx context.Context, input dtos.PostFilterDTO) (*dtos.PagedResponse[dtos.PostDTO], error) {
    query := facades.Orm().Query()
    
    // Apply custom filters
    if input.Title != nil && *input.Title != "" {
        query = query.Where("title LIKE ?", "%"+*input.Title+"%")
    }
    
    if input.Published != nil {
        query = query.Where("published = ?", *input.Published)
    }
    
    if input.UserID != nil {
        query = query.Where("user_id = ?", *input.UserID)
    }
    
    if input.Filter != "" {
        query = query.Where("title LIKE ? OR content LIKE ?", "%"+input.Filter+"%", "%"+input.Filter+"%")
    }
    
    // Apply sorting
    if input.Sort != "" {
        query = query.Order(input.Sort)
    } else {
        query = query.Order("created_at desc")
    }
    
    // Get total count
    var totalCount int64
    countQuery := query
    countQuery.Count(&totalCount)
    
    // Apply pagination
    offset := (input.Page - 1) * input.PageSize
    query = query.Offset(offset).Limit(input.PageSize)
    
    // Load relationships
    query = query.With("User")
    
    // Execute query
    var posts []models.Post
    err := query.Get(&posts)
    if err != nil {
        return nil, err
    }
    
    // Map to DTOs
    var postDTOs []dtos.PostDTO
    for _, post := range posts {
        dto := s.mapToPostDTO(post)
        postDTOs = append(postDTOs, dto)
    }
    
    totalPages := int(math.Ceil(float64(totalCount) / float64(input.PageSize)))
    
    return &dtos.PagedResponse[dtos.PostDTO]{
        Data:       postDTOs,
        TotalCount: int(totalCount),
        Page:       input.Page,
        PageSize:   input.PageSize,
        TotalPages: totalPages,
    }, nil
}

// MapToEntity implementation
func (s *PostAppService) MapToEntity(input dtos.CreatePostDTO) (*models.Post, error) {
    return &models.Post{
        Title:     input.Title,
        Content:   input.Content,
        Published: input.Published,
        UserID:    input.UserID,
    }, nil
}

// MapToEntityDTO implementation
func (s *PostAppService) MapToEntityDTO(entity models.Post) (*dtos.PostDTO, error) {
    return &s.mapToPostDTO(entity), nil
}

func (s *PostAppService) mapToPostDTO(post models.Post) dtos.PostDTO {
    dto := dtos.PostDTO{
        BaseEntityDTO: dtos.BaseEntityDTO{
            ID:        post.ID,
            CreatedAt: post.CreatedAt,
            UpdatedAt: post.UpdatedAt,
            DeletedAt: post.DeletedAt.Time,
        },
        Title:     post.Title,
        Content:   post.Content,
        Published: post.Published,
        UserID:    post.UserID,
    }
    
    // Map user if loaded
    if post.User.ID != 0 {
        dto.User = &dtos.UserDTO{
            BaseEntityDTO: dtos.BaseEntityDTO{
                ID:        post.User.ID,
                CreatedAt: post.User.CreatedAt,
                UpdatedAt: post.User.UpdatedAt,
                DeletedAt: post.User.DeletedAt.Time,
            },
            Name:  post.User.Name,
            Email: post.User.Email,
        }
    }
    
    return dto
}
```

### 7. Swagger Documentation

**Main Application Setup (`main.go`)**:
```go
package main

import (
    "goravel/app/providers"
    "goravel/bootstrap"
    _ "goravel/docs" // Swagger docs
    
    "github.com/goravel/framework/facades"
    "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Goravel CrudAppService API
// @version         1.0
// @description     A comprehensive CRUD API built with Goravel using the CrudAppService pattern
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.goravel.dev/support
// @contact.email  support@goravel.dev

// @license.name  MIT
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3000
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
    // Bootstrap the application
    app := bootstrap.NewApplication()
    
    // Boot service providers
    app.Boot()
    
    // Register Swagger routes
    facades.Route().Get("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
    // Start the application
    app.Run()
}
```

**Routes Setup (`routes/api.go`)**:
```go
package routes

import (
    "goravel/app/http/controllers"
    "goravel/app/http/middleware"
    "goravel/app/services"
    
    "github.com/goravel/framework/facades"
)

func Api() {
    // Initialize services
    postAppService := services.NewPostAppService()
    
    // Initialize controllers
    postController := controllers.NewPostController(postAppService)
    
    // API routes with authentication middleware
    api := facades.Route().Prefix("api").Middleware(middleware.Auth())
    {
        // Post routes
        posts := api.Prefix("posts")
        {
            posts.Get("/", postController.Index)
            posts.Get("/{id}", postController.Show)
            posts.Post("/", postController.Store)
            posts.Put("/{id}", postController.Update)
            posts.Delete("/{id}", postController.Destroy)
        }
    }
}
```

## Code Generation

Create custom Artisan commands to generate the boilerplate code:

**Make CrudAppService Command (`app/console/commands/make_crud_app_service.go`)**:
```go
package commands

import (
    "fmt"
    "os"
    "strings"
    "text/template"
    
    "github.com/goravel/framework/contracts/console"
    "github.com/goravel/framework/contracts/console/command"
)

type MakeCrudAppService struct{}

func (receiver *MakeCrudAppService) Signature() string {
    return "make:crud-app-service"
}

func (receiver *MakeCrudAppService) Description() string {
    return "Create a new CRUD app service with all necessary files"
}

func (receiver *MakeCrudAppService) Extend() command.Extend {
    return command.Extend{
        Category: "make",
        Flags: []command.Flag{
            &command.StringFlag{
                Name:    "name",
                Aliases: []string{"n"},
                Usage:   "Name of the entity (e.g., Post, User)",
                Value:   "",
            },
            &command.BoolFlag{
                Name:    "with-controller",
                Aliases: []string{"c"},
                Usage:   "Generate controller along with app service",
                Value:   true,
            },
            &command.BoolFlag{
                Name:    "with-requests",
                Aliases: []string{"r"},
                Usage:   "Generate form request classes",
                Value:   true,
            },
        },
    }
}

func (receiver *MakeCrudAppService) Handle(ctx console.Context) error {
    name := ctx.Option("name")
    if name == "" {
        return fmt.Errorf("entity name is required")
    }
    
    withController := ctx.OptionBool("with-controller")
    withRequests := ctx.OptionBool("with-requests")
    
    // Generate files
    if err := receiver.generateDTO(name); err != nil {
        return err
    }
    
    if err := receiver.generateAppService(name); err != nil {
        return err
    }
    
    if withController {
        if err := receiver.generateController(name); err != nil {
            return err
        }
    }
    
    if withRequests {
        if err := receiver.generateRequests(name); err != nil {
            return err
        }
    }
    
    ctx.Info(fmt.Sprintf("CRUD App Service for %s created successfully!", name))
    
    // Show next steps
    ctx.Info("Next steps:")
    ctx.Info("1. Update your model to implement necessary relationships")
    ctx.Info("2. Add routes in routes/api.go")
    ctx.Info("3. Run 'swag init' to generate Swagger documentation")
    ctx.Info("4. Register gates in app/providers/auth_service_provider.go")
    
    return nil
}

func (receiver *MakeCrudAppService) generateDTO(name string) error {
    // Template for DTO generation
    tmpl := `package dtos

// {{.Name}}DTO for reading operations
type {{.Name}}DTO struct {
    BaseEntityDTO
    // Add your fields here
    Name string ` + "`json:\"name\"`" + `
}

// Create{{.Name}}DTO for creation
type Create{{.Name}}DTO struct {
    Name string ` + "`json:\"name\" validate:\"required,min:3,max:255\" example:\"Example Name\"`" + `
}

// Update{{.Name}}DTO for updates
type Update{{.Name}}DTO struct {
    Name *string ` + "`json:\"name,omitempty\" validate:\"omitempty,min:3,max:255\"`" + `
}

// {{.Name}}FilterDTO for filtering
type {{.Name}}FilterDTO struct {
    PagedRequest
    Name *string ` + "`form:\"name\" json:\"name,omitempty\"`" + `
}
`
    
    return receiver.writeTemplate(fmt.Sprintf("app/dtos/%s_dto.go", strings.ToLower(name)), tmpl, map[string]string{
        "Name": name,
    })
}

func (receiver *MakeCrudAppService) generateAppService(name string) error {
    // Template for App Service generation
    tmpl := `package services

import (
    "context"
    "goravel/app/contracts"
    "goravel/app/dtos"
    "goravel/app/models"
)

type {{.Name}}AppService struct {
    *BaseCrudAppService[models.{{.Name}}, dtos.{{.Name}}DTO, uint, dtos.{{.Name}}FilterDTO, dtos.Create{{.Name}}DTO, dtos.Update{{.Name}}DTO]
}

func New{{.Name}}AppService() contracts.I{{.Name}}AppService {
    base := NewBaseCrudAppService[models.{{.Name}}, dtos.{{.Name}}DTO, uint, dtos.{{.Name}}FilterDTO, dtos.Create{{.Name}}DTO, dtos.Update{{.Name}}DTO]("{{.Name}}")
    return &{{.Name}}AppService{
        BaseCrudAppService: base,
    }
}

// MapToEntity implementation
func (s *{{.Name}}AppService) MapToEntity(input dtos.Create{{.Name}}DTO) (*models.{{.Name}}, error) {
    return &models.{{.Name}}{
        Name: input.Name,
        // Add other field mappings here
    }, nil
}

// MapToEntityDTO implementation
func (s *{{.Name}}AppService) MapToEntityDTO(entity models.{{.Name}}) (*dtos.{{.Name}}DTO, error) {
    return &dtos.{{.Name}}DTO{
        BaseEntityDTO: dtos.BaseEntityDTO{
            ID:        entity.ID,
            CreatedAt: entity.CreatedAt,
            UpdatedAt: entity.UpdatedAt,
            DeletedAt: entity.DeletedAt.Time,
        },
        Name: entity.Name,
        // Add other field mappings here
    }, nil
}
`
    
    return receiver.writeTemplate(fmt.Sprintf("app/services/%s_app_service.go", strings.ToLower(name)), tmpl, map[string]string{
        "Name": name,
    })
}

func (receiver *MakeCrudAppService) generateController(name string) error {
    // Template for Controller generation
    tmpl := `package controllers

import (
    "goravel/app/contracts"
    "goravel/app/http/requests"
    "strconv"
    
    "github.com/goravel/framework/contracts/http"
    "github.com/goravel/framework/facades"
)

type {{.Name}}Controller struct {
    {{.LowerName}}AppService contracts.I{{.Name}}AppService
}

func New{{.Name}}Controller({{.LowerName}}AppService contracts.I{{.Name}}AppService) *{{.Name}}Controller {
    return &{{.Name}}Controller{
        {{.LowerName}}AppService: {{.LowerName}}AppService,
    }
}

// Index godoc
// @Summary      List {{.LowerName}}s
// @Description  Get paginated list of {{.LowerName}}s with filtering and sorting
// @Tags         {{.LowerName}}s
// @Accept       json
// @Produce      json
// @Param        page       query   int     false  "Page number" default(1)
// @Param        page_size  query   int     false  "Items per page" default(10)
// @Param        sort       query   string  false  "Sort order" default("id desc")
// @Param        filter     query   string  false  "Search filter"
// @Success      200  {object}  dtos.PagedResponse[dtos.{{.Name}}DTO]
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/{{.LowerName}}s [get]
func (r *{{.Name}}Controller) Index(ctx http.Context) http.Response {
    // Implementation here
    return ctx.Response().Json(http.StatusOK, map[string]string{"message": "Index method"})
}

// Show godoc
// @Summary      Get {{.LowerName}} by ID
// @Description  Get single {{.LowerName}} by ID
// @Tags         {{.LowerName}}s
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "{{.Name}} ID"
// @Success      200  {object}  dtos.{{.Name}}DTO
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/{{.LowerName}}s/{id} [get]
func (r *{{.Name}}Controller) Show(ctx http.Context) http.Response {
    // Implementation here
    return ctx.Response().Json(http.StatusOK, map[string]string{"message": "Show method"})
}

// Store godoc
// @Summary      Create new {{.LowerName}}
// @Description  Create a new {{.LowerName}}
// @Tags         {{.LowerName}}s
// @Accept       json
// @Produce      json
// @Param        request  body      dtos.Create{{.Name}}DTO  true  "{{.Name}} data"
// @Success      201      {object}  dtos.{{.Name}}DTO
// @Failure      400      {object}  map[string]interface{}
// @Failure      403      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/{{.LowerName}}s [post]
// @Security     BearerAuth
func (r *{{.Name}}Controller) Store(ctx http.Context) http.Response {
    // Implementation here
    return ctx.Response().Json(http.StatusCreated, map[string]string{"message": "Store method"})
}

// Update godoc
// @Summary      Update {{.LowerName}}
// @Description  Update an existing {{.LowerName}}
// @Tags         {{.LowerName}}s
// @Accept       json
// @Produce      json
// @Param        id       path      int                 true  "{{.Name}} ID"
// @Param        request  body      dtos.Update{{.Name}}DTO  true  "{{.Name}} update data"
// @Success      200      {object}  dtos.{{.Name}}DTO
// @Failure      400      {object}  map[string]interface{}
// @Failure      403      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/{{.LowerName}}s/{id} [put]
// @Security     BearerAuth
func (r *{{.Name}}Controller) Update(ctx http.Context) http.Response {
    // Implementation here
    return ctx.Response().Json(http.StatusOK, map[string]string{"message": "Update method"})
}

// Destroy godoc
// @Summary      Delete {{.LowerName}}
// @Description  Delete a {{.LowerName}} by ID
// @Tags         {{.LowerName}}s
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "{{.Name}} ID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/{{.LowerName}}s/{id} [delete]
// @Security     BearerAuth
func (r *{{.Name}}Controller) Destroy(ctx http.Context) http.Response {
    // Implementation here
    return ctx.Response().Json(http.StatusNoContent, nil)
}
`
    
    return receiver.writeTemplate(fmt.Sprintf("app/http/controllers/%s_controller.go", strings.ToLower(name)), tmpl, map[string]string{
        "Name":      name,
        "LowerName": strings.ToLower(name),
    })
}

func (receiver *MakeCrudAppService) generateRequests(name string) error {
    // Template for Request classes
    tmpl := `package requests

import (
    "goravel/app/dtos"
    
    "github.com/goravel/framework/contracts/http"
    "github.com/goravel/framework/facades"
)

type Create{{.Name}}Request struct {
    dtos.Create{{.Name}}DTO
}

func (r *Create{{.Name}}Request) Authorize(ctx http.Context) error {
    if !facades.Gate().WithContext(ctx).Allows("create-{{.LowerName}}", map[string]any{}) {
        return errors.New("unauthorized to create {{.LowerName}}s")
    }
    return nil
}

func (r *Create{{.Name}}Request) Rules(ctx http.Context) map[string]string {
    return map[string]string{
        "name": "required|min_len:3|max_len:255",
        // Add more validation rules here
    }
}

type Update{{.Name}}Request struct {
    dtos.Update{{.Name}}DTO
}

func (r *Update{{.Name}}Request) Authorize(ctx http.Context) error {
    {{.LowerName}}ID := ctx.Request().Route("id")
    if !facades.Gate().WithContext(ctx).Allows("update-{{.LowerName}}", map[string]any{
        "{{.LowerName}}_id": {{.LowerName}}ID,
    }) {
        return errors.New("unauthorized to update this {{.LowerName}}")
    }
    return nil
}

func (r *Update{{.Name}}Request) Rules(ctx http.Context) map[string]string {
    return map[string]string{
        "name": "min_len:3|max_len:255",
        // Add more validation rules here
    }
}
`
    
    return receiver.writeTemplate(fmt.Sprintf("app/http/requests/%s_request.go", strings.ToLower(name)), tmpl, map[string]string{
        "Name":      name,
        "LowerName": strings.ToLower(name),
    })
}

func (receiver *MakeCrudAppService) writeTemplate(filename, tmplStr string, data map[string]string) error {
    tmpl, err := template.New("file").Parse(tmplStr)
    if err != nil {
        return err
    }
    
    // Create directory if it doesn't exist
    dir := filepath.Dir(filename)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }
    
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    return tmpl.Execute(file, data)
}
```

## Usage Examples

### Generate a complete CRUD setup:

```bash
# Generate a complete CRUD setup for a Product entity
go run . artisan make:crud-app-service --name=Product --with-controller --with-requests

# This generates:
# - app/dtos/product_dto.go
# - app/services/product_app_service.go
# - app/http/controllers/product_controller.go
# - app/http/requests/product_request.go
```

### API Usage Examples:

```bash
# Get paginated posts
GET /api/posts?page=1&page_size=10&sort=created_at desc&filter=golang

# Get specific post
GET /api/posts/123

# Create new post
POST /api/posts
{
  "title": "My New Post",
  "content": "This is the content of my post",
  "published": true,
  "user_id": 1
}

# Update post
PUT /api/posts/123
{
  "title": "Updated Title",
  "published": false
}

# Delete post
DELETE /api/posts/123
```

### Generate Swagger Documentation:

```bash
# Generate Swagger documentation
swag init

# Start the application and visit:
# http://localhost:3000/swagger/index.html
```

## Best Practices

### 1. **Separation of Concerns**
- Keep business logic in AppServices
- Use Controllers only for HTTP concerns
- Implement validation in Form Requests
- Handle authorization in Gates

### 2. **Type Safety**
- Use generics for reusable components
- Define clear DTO contracts
- Implement proper error handling

### 3. **Security**
- Always validate input data
- Implement proper authorization
- Use HTTPS in production
- Validate file uploads

### 4. **Performance**
- Implement proper pagination
- Use database indexes
- Consider caching for read operations
- Use lazy loading for relationships

### 5. **Documentation**
- Keep Swagger annotations up to date
- Document all API endpoints
- Provide examples in DTOs
- Include error response formats

### 6. **Testing**
- Write unit tests for AppServices
- Create integration tests for controllers
- Mock external dependencies
- Test authorization logic

## Conclusion

This implementation provides a robust, scalable foundation for building CRUD APIs in Goravel that mirrors the powerful CrudAppService pattern from ABP.io. Key benefits include:

- **Consistency**: Standardized patterns across all entities
- **Maintainability**: Clear separation of concerns and reusable components
- **Security**: Built-in validation and authorization
- **Documentation**: Auto-generated, always up-to-date API docs
- **Developer Experience**: Code generation and type safety

The pattern scales well from simple CRUD operations to complex business scenarios while maintaining clean, testable code. The automated tooling reduces boilerplate and ensures consistency across your application.

### Next Steps:

1. **Extend the pattern**: Add support for bulk operations, audit logging, and advanced filtering
2. **Performance optimization**: Implement caching, query optimization, and pagination strategies
3. **Advanced features**: Add support for file uploads, real-time notifications, and background jobs
4. **Testing suite**: Create comprehensive test templates and utilities
5. **CLI enhancements**: Add more sophisticated code generation options
```