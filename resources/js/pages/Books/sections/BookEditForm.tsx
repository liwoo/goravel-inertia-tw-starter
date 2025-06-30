import React, { useState, forwardRef, useImperativeHandle } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Badge } from '@/components/ui/badge';
import { CrudEditFormProps } from '@/types/crud';
import { Book, BookUpdateData, BookStatus } from '@/types/book';
import { Calendar, DollarSign, Hash, Plus, Tag, X } from 'lucide-react';

interface BookEditFormProps extends CrudEditFormProps<Book> {
  setIsSaving?: (saving: boolean) => void;
}

export const BookEditForm = forwardRef<any, BookEditFormProps>(({ 
  item: book,
  onSuccess,
  onError,
  onCancel, 
  isLoading = false,
  setIsSaving
}, ref) => {
  const [formData, setFormData] = useState<BookUpdateData>({
    title: book.title,
    author: book.author,
    isbn: book.isbn,
    description: book.description || '',
    price: book.price,
    status: book.status,
    publishedAt: book.publishedAt ? new Date(book.publishedAt).toISOString().split('T')[0] : '',
    tags: book.tags || [],
  });

  const [errors, setErrors] = useState<Record<string, string>>({});
  const [tagInput, setTagInput] = useState('');

  const handleSubmit = async () => {
    
    // Basic validation
    const newErrors: Record<string, string> = {};
    if (!formData.title?.trim()) {
      newErrors.title = 'Title is required';
    }
    if (!formData.author?.trim()) {
      newErrors.author = 'Author is required';
    }
    if (!formData.isbn?.trim()) {
      newErrors.isbn = 'ISBN is required';
    }
    if (formData.price !== undefined && formData.price < 0) {
      newErrors.price = 'Price must be a positive number';
    }
    
    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }
    
    setErrors({});
    setIsSaving?.(true);
    
    try {
      const response = await fetch(`/api/books/${book.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify(formData),
      });

      if (response.ok) {
        onSuccess('Book updated successfully');
      } else {
        const errorData = await response.json().catch(() => ({}));
        onError?.(errorData);
      }
    } catch (error) {
      onError?.(error);
    } finally {
      setIsSaving?.(false);
    }
  };

  // Expose handleSubmit to parent component
  useImperativeHandle(ref, () => ({
    handleSubmit
  }));

  const handleAddTag = () => {
    if (tagInput.trim() && formData.tags && !formData.tags.includes(tagInput.trim())) {
      setFormData({
        ...formData,
        tags: [...formData.tags, tagInput.trim()]
      });
      setTagInput('');
    }
  };

  const handleRemoveTag = (tagToRemove: string) => {
    setFormData({
      ...formData,
      tags: formData.tags?.filter(tag => tag !== tagToRemove) || []
    });
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      handleAddTag();
    }
  };

  return (
    <form onSubmit={(e) => e.preventDefault()} className="space-y-6">
      {/* Book Information */}
      <div className="space-y-4">
        <h3 className="text-lg font-semibold">Book Information</h3>
        
        <div className="space-y-2">
          <Label htmlFor="title">Title *</Label>
          <Input
            id="title"
            value={formData.title}
            onChange={(e) => setFormData({ ...formData, title: e.target.value })}
            placeholder="Enter book title"
            className={errors.title ? 'border-destructive' : ''}
          />
          {errors.title && (
            <p className="text-sm text-destructive">{errors.title}</p>
          )}
        </div>

        <div className="space-y-2">
          <Label htmlFor="author">Author *</Label>
          <Input
            id="author"
            value={formData.author}
            onChange={(e) => setFormData({ ...formData, author: e.target.value })}
            placeholder="Enter author name"
            className={errors.author ? 'border-destructive' : ''}
          />
          {errors.author && (
            <p className="text-sm text-destructive">{errors.author}</p>
          )}
        </div>

        <div className="space-y-2">
          <Label htmlFor="isbn">ISBN *</Label>
          <div className="relative">
            <Input
              id="isbn"
              value={formData.isbn}
              onChange={(e) => setFormData({ ...formData, isbn: e.target.value })}
              placeholder="Enter ISBN number"
              className={`pl-10 ${errors.isbn ? 'border-destructive' : ''}`}
            />
            <Hash className="absolute left-3 top-2.5 h-5 w-5 text-gray-400" />
          </div>
          {errors.isbn && (
            <p className="text-sm text-destructive">{errors.isbn}</p>
          )}
        </div>

        <div className="space-y-2">
          <Label htmlFor="description">Description</Label>
          <Textarea
            id="description"
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            placeholder="Enter book description"
            rows={4}
          />
        </div>
      </div>

      {/* Publication Details */}
      <div className="space-y-4">
        <h3 className="text-lg font-semibold">Publication Details</h3>
        
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label htmlFor="price">Price</Label>
            <div className="relative">
              <Input
                id="price"
                type="number"
                step="0.01"
                value={formData.price}
                onChange={(e) => setFormData({ ...formData, price: parseFloat(e.target.value) || 0 })}
                placeholder="0.00"
                className={`pl-10 ${errors.price ? 'border-destructive' : ''}`}
              />
              <DollarSign className="absolute left-3 top-2.5 h-5 w-5 text-gray-400" />
            </div>
            {errors.price && (
              <p className="text-sm text-destructive">{errors.price}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="publishedAt">Published Date</Label>
            <div className="relative">
              <Input
                id="publishedAt"
                type="date"
                value={formData.publishedAt}
                onChange={(e) => setFormData({ ...formData, publishedAt: e.target.value })}
                className="pl-10"
              />
              <Calendar className="absolute left-3 top-2.5 h-5 w-5 text-gray-400" />
            </div>
          </div>
        </div>

        <div className="space-y-2">
          <Label htmlFor="status">Status</Label>
          <Select
            value={formData.status}
            onValueChange={(value) => setFormData({ ...formData, status: value as BookStatus })}
          >
            <SelectTrigger>
              <SelectValue placeholder="Select status" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="AVAILABLE">Available</SelectItem>
              <SelectItem value="CHECKED_OUT">Checked Out</SelectItem>
              <SelectItem value="RESERVED">Reserved</SelectItem>
              <SelectItem value="LOST">Lost</SelectItem>
              <SelectItem value="DAMAGED">Damaged</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div className="space-y-2">
          <Label htmlFor="tags">Tags</Label>
          <div className="flex gap-2">
            <div className="relative flex-1">
              <Input
                id="tags"
                value={tagInput}
                onChange={(e) => setTagInput(e.target.value)}
                onKeyPress={handleKeyPress}
                placeholder="Add tags"
                className="pl-10"
              />
              <Tag className="absolute left-3 top-2.5 h-5 w-5 text-gray-400" />
            </div>
            <Button
              type="button"
              variant="outline"
              onClick={handleAddTag}
            >
              <Plus className="h-4 w-4" />
            </Button>
          </div>
          {formData.tags && formData.tags.length > 0 && (
            <div className="flex flex-wrap gap-2 mt-2">
              {formData.tags.map((tag, index) => (
                <Badge key={index} variant="secondary" className="gap-1">
                  {tag}
                  <button
                    type="button"
                    onClick={() => handleRemoveTag(tag)}
                    className="ml-1 hover:text-destructive"
                  >
                    <X className="h-3 w-3" />
                  </button>
                </Badge>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Metadata */}
      <div className="space-y-4">
        <h3 className="text-lg font-semibold">Metadata</h3>
        <div className="grid grid-cols-2 gap-4 text-sm">
          <div>
            <p className="text-gray-500">Book ID</p>
            <p className="font-medium">#{book.id}</p>
          </div>
          <div>
            <p className="text-gray-500">Created</p>
            <p className="font-medium">
              {new Date(book.createdAt).toLocaleDateString()}
            </p>
          </div>
          <div>
            <p className="text-gray-500">Last Updated</p>
            <p className="font-medium">
              {new Date(book.updatedAt).toLocaleDateString()}
            </p>
          </div>
        </div>
      </div>

    </form>
  );
});

BookEditForm.displayName = 'BookEditForm';