# CRUDUI.md - Goravel Inertia React CRUD UI Specification

## Overview

This specification outlines the implementation of a comprehensive CRUD UI system using Goravel with Inertia.js, TypeScript, React, and shadcn/ui. The system provides a generic, reusable CRUD interface that automatically handles data tables, forms, validation, pagination, filtering, and actions.

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Core Components](#core-components)
- [Type System](#type-system)
- [Implementation Requirements](#implementation-requirements)
- [Usage Examples](#usage-examples)
- [File Structure](#file-structure)
- [Development Workflow](#development-workflow)

## Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Goravel       │    │   Inertia.js    │    │     React       │
│   Backend       │◄──►│   Bridge        │◄──►│   Frontend      │
│                 │    │                 │    │                 │
│ • CrudAppService│    │ • Route Binding │    │ • Generic CRUD  │
│ • Controllers   │    │ • Data Transfer │    │ • Type Safety   │
│ • DTOs          │    │ • State Sync    │    │ • Components    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Routes    │    │   Shared Props  │    │   UI Library    │
│   /posts        │    │   Page Data     │    │   shadcn/ui     │
│   /posts/create │    │   Form State    │    │   Data Tables   │
│   /posts/{id}   │    │   Validation    │    │   Forms/Dialogs │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Core Components

### 1. Backend Route Configuration

**File:** `routes/web.go`

```go
package routes

import (
    "goravel/app/http/controllers"
    "goravel/app/http/middleware"
    
    "github.com/goravel/framework/facades"
)

func Web() {
    // Initialize controllers
    pageController := controllers.NewPageController()
    
    // Authentication middleware group
    auth := facades.Route().Middleware(middleware.Auth())
    
    // CRUD UI Routes - Single route per entity that renders Inertia page
    auth.Get("/posts", pageController.PostsIndex).Name("posts.index")
    auth.Get("/users", pageController.UsersIndex).Name("users.index")
    
    // All CRUD operations are handled by API endpoints (in routes/api.go)
    // The frontend makes API calls to handle Create, Read, Update, Delete operations
}
```
```

### 2. Simple Inertia Page Controller

**File:** `app/http/controllers/page_controller.go`

```go
package controllers

import (
    "goravel/app/dtos"
    "goravel/app/services"
    "strconv"
    
    "github.com/goravel/framework/contracts/http"
    "github.com/goravel/framework/facades"
)

type PageController struct {
    postAppService contracts.IPostAppService
    userAppService contracts.IUserAppService
}

func NewPageController() *PageController {
    return &PageController{
        postAppService: services.NewPostAppService(),
        userAppService: services.NewUserAppService(),
    }
}

// PostsIndex renders the Posts CRUD page with initial data
func (c *PageController) PostsIndex(ctx http.Context) http.Response {
    // Get initial data for the page (first page only)
    filter := dtos.PostFilterDTO{
        PagedRequest: dtos.PagedRequest{
            Page:     1,
            PageSize: 10,
        },
    }
    
    posts, err := c.postAppService.GetListAsync(ctx.Context(), filter)
    if err != nil {
        facades.Log().Error("Failed to get posts: " + err.Error())
        posts = &dtos.PagedResponse[dtos.PostDTO]{
            Data:       []dtos.PostDTO{},
            TotalCount: 0,
            Page:       1,
            PageSize:   10,
            TotalPages: 0,
        }
    }
    
    // Get reference data for dropdowns/forms
    users, _ := c.userAppService.GetListAsync(ctx.Context(), dtos.UserFilterDTO{
        PagedRequest: dtos.PagedRequest{Page: 1, PageSize: 100}, // Get all users for dropdown
    })
    
    return ctx.Response().View().Make("Posts", map[string]interface{}{
        "initialData": posts,
        "referenceData": map[string]interface{}{
            "users": users.Data,
        },
        "meta": map[string]interface{}{
            "title":       "Posts Management",
            "description": "Manage your blog posts",
            "entityName":  "post",
            "entityLabel": "Post",
            "permissions": map[string]bool{
                "create": facades.Gate().WithContext(ctx).Allows("create-post", map[string]any{}),
                "update": facades.Gate().WithContext(ctx).Allows("update-post", map[string]any{}),
                "delete": facades.Gate().WithContext(ctx).Allows("delete-post", map[string]any{}),
                "export": facades.Gate().WithContext(ctx).Allows("export-post", map[string]any{}),
            },
        },
        "config": map[string]interface{}{
            "apiEndpoint":   "/api/posts",
            "exportFormats": []string{"excel", "csv", "pdf"},
            "pageSize":      10,
            "sortBy":        "created_at",
            "sortOrder":     "desc",
        },
    })
}
```

### 3. TypeScript Type Definitions

**File:** `resources/js/types/post.d.ts`

```typescript
// Base types
export interface BaseEntityDTO {
  id: number;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

export interface PagedRequest {
  page: number;
  page_size: number;
  sort?: string;
  filter?: string;
}

export interface PagedResponse<T> {
  data: T[];
  total_count: number;
  page: number;
  page_size: number;
  total_pages: number;
}

// Post-specific types
export interface PostDTO extends BaseEntityDTO {
  title: string;
  content: string;
  published: boolean;
  user_id: number;
  user?: UserDTO;
}

export interface CreatePostDTO {
  title: string;
  content: string;
  published: boolean;
  user_id: number;
}

export interface UpdatePostDTO {
  title?: string;
  content?: string;
  published?: boolean;
}

export interface PostFilterDTO extends PagedRequest {
  title?: string;
  published?: boolean;
  user_id?: number;
}

export interface UserDTO extends BaseEntityDTO {
  name: string;
  email: string;
}

// UI-specific types
export interface PostFormData {
  title: string;
  content: string;
  published: boolean;
  user_id: number;
}

export interface PostTableRow {
  id: number;
  title: string;
  content: string;
  published: boolean;
  user_name: string;
  created_at: string;
  actions?: string;
}

export interface CrudPageProps<TEntity, TFilter> {
  initialData: PagedResponse<TEntity>;
  meta: {
    title: string;
    description: string;
    entityName: string;
    entityLabel: string;
    permissions: {
      create: boolean;
      update: boolean;
      delete: boolean;
      export: boolean;
    };
  };
  config: {
    apiEndpoint: string;
    exportFormats: string[];
    pageSize: number;
    sortBy: string;
    sortOrder: 'asc' | 'desc';
  };
  referenceData?: Record<string, any[]>;
}
```

### 4. Generic CRUD Page Component

**File:** `resources/js/Components/Crud/CrudPage.tsx`

```typescript
import React, { useState, useEffect } from 'react';
import { Head, router } from '@inertiajs/react';
import { toast } from 'sonner';
import { z } from 'zod';

import { DataTable } from './DataTable';
import { CreateUpdateForm } from './CreateUpdateForm';
import { DetailView } from './DetailView';
import { ExportDialog } from './ExportDialog';
import { BulkActionDialog } from './BulkActionDialog';

import { Button } from '@/Components/ui/button';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/Components/ui/dialog';
import { Sheet, SheetContent, SheetHeader, SheetTitle } from '@/Components/ui/sheet';

interface CrudPageProps<
  TEntity,
  TCreateDTO,
  TUpdateDTO,
  TFilterDTO,
  TFormData,
  TTableRow
> {
  initialData: PagedResponse<TEntity>;
  meta: {
    title: string;
    description: string;
    entityName: string;
    entityLabel: string;
    permissions: Record<string, boolean>;
  };
  config: {
    apiEndpoint: string;
    exportFormats: string[];
    pageSize: number;
    sortBy: string;
    sortOrder: 'asc' | 'desc';
  };
  referenceData?: Record<string, any[]>;
  
  // Component definitions
  tableColumns: TableColumnDef<TTableRow>[];
  formFields: FormFieldDef<TFormData>[];
  detailFields: DetailFieldDef<TEntity>[];
  validationSchema: z.ZodSchema<TFormData>;
  
  // Mapping functions
  entityToTableRow: (entity: TEntity) => TTableRow;
  entityToFormData: (entity: TEntity) => TFormData;
  formDataToCreateDTO: (formData: TFormData) => TCreateDTO;
  formDataToUpdateDTO: (formData: TFormData) => TUpdateDTO;
}

export function CrudPage<
  TEntity,
  TCreateDTO,
  TUpdateDTO,
  TFilterDTO,
  TFormData,
  TTableRow
>({
  initialData,
  meta,
  config,
  referenceData = {},
  tableColumns,
  formFields,
  detailFields,
  validationSchema,
  entityToTableRow,
  entityToFormData,
  formDataToCreateDTO,
  formDataToUpdateDTO,
}: CrudPageProps<TEntity, TCreateDTO, TUpdateDTO, TFilterDTO, TFormData, TTableRow>) {
  
  // State management
  const [data, setData] = useState<PagedResponse<TEntity>>(initialData);
  const [loading, setLoading] = useState(false);
  const [selectedRows, setSelectedRows] = useState<number[]>([]);
  const [filters, setFilters] = useState<Partial<TFilterDTO>>({});
  
  // Dialog states
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [showUpdateForm, setShowUpdateForm] = useState(false);
  const [showDetailView, setShowDetailView] = useState(false);
  const [showExportDialog, setShowExportDialog] = useState(false);
  const [showBulkDialog, setShowBulkDialog] = useState(false);
  
  const [selectedEntity, setSelectedEntity] = useState<TEntity | null>(null);

  // API calls
  const fetchData = async (newFilters?: Partial<TFilterDTO>) => {
    setLoading(true);
    try {
      const params = new URLSearchParams({
        page: String(newFilters?.page || 1),
        page_size: String(newFilters?.page_size || config.pageSize),
        sort: newFilters?.sort || `${config.sortBy} ${config.sortOrder}`,
        ...Object.fromEntries(
          Object.entries(newFilters || {}).filter(([_, value]) => 
            value !== undefined && value !== null && value !== ''
          )
        ),
      });

      const response = await fetch(`${config.apiEndpoint}?${params}`);
      if (!response.ok) throw new Error('Failed to fetch data');
      
      const result = await response.json();
      setData(result);
    } catch (error) {
      toast.error('Failed to load data');
      console.error('Error fetching data:', error);
    } finally {
      setLoading(false);
    }
  };

  const createEntity = async (formData: TFormData) => {
    try {
      const createDTO = formDataToCreateDTO(formData);
      const response = await fetch(config.apiEndpoint, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(createDTO),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.message || 'Failed to create entity');
      }

      toast.success(`${meta.entityLabel} created successfully`);
      setShowCreateForm(false);
      fetchData(filters);
    } catch (error) {
      toast.error(error.message || `Failed to create ${meta.entityName}`);
      throw error;
    }
  };

  const updateEntity = async (id: number, formData: TFormData) => {
    try {
      const updateDTO = formDataToUpdateDTO(formData);
      const response = await fetch(`${config.apiEndpoint}/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(updateDTO),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.message || 'Failed to update entity');
      }

      toast.success(`${meta.entityLabel} updated successfully`);
      setShowUpdateForm(false);
      setSelectedEntity(null);
      fetchData(filters);
    } catch (error) {
      toast.error(error.message || `Failed to update ${meta.entityName}`);
      throw error;
    }
  };

  const deleteEntity = async (id: number) => {
    try {
      const response = await fetch(`${config.apiEndpoint}/${id}`, {
        method: 'DELETE',
      });

      if (!response.ok) throw new Error('Failed to delete entity');

      toast.success(`${meta.entityLabel} deleted successfully`);
      fetchData(filters);
    } catch (error) {
      toast.error(`Failed to delete ${meta.entityName}`);
    }
  };

  const bulkDelete = async (ids: number[]) => {
    try {
      const response = await fetch(`${config.apiEndpoint}/bulk-delete`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ ids }),
      });

      if (!response.ok) throw new Error('Failed to delete entities');

      toast.success(`${ids.length} ${meta.entityName}(s) deleted successfully`);
      setSelectedRows([]);
      setShowBulkDialog(false);
      fetchData(filters);
    } catch (error) {
      toast.error(`Failed to delete selected ${meta.entityName}s`);
    }
  };

  const exportData = async (format: string, filters?: Partial<TFilterDTO>) => {
    try {
      const params = new URLSearchParams({
        format,
        ...Object.fromEntries(
          Object.entries(filters || {}).filter(([_, value]) => 
            value !== undefined && value !== null && value !== ''
          )
        ),
      });

      const response = await fetch(`${config.apiEndpoint}/export?${params}`, {
        method: 'POST',
      });

      if (!response.ok) throw new Error('Failed to export data');

      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `${meta.entityName}s_export.${format}`;
      link.click();
      window.URL.revokeObjectURL(url);

      toast.success('