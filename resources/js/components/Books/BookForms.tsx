import React, { useState, useEffect } from 'react';
import { router } from '@inertiajs/react';
import { toast } from 'sonner';
import { Book, BookCreateData, BookUpdateData, BookFormErrors, BOOK_VALIDATION_RULES, BookStatus, BOOK_STATUS_CONFIG } from '@/types/book';
import { CrudFormProps, CrudEditFormProps, CrudDetailViewProps } from '@/types/crud';
import { FormField } from '@/components/Form/FormField';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { StatusBadge } from '@/components/Crud/StatusBadge';
import { 
  BookOpen, 
  User, 
  Hash, 
  DollarSign, 
  Calendar, 
  Tag, 
  AlertCircle,
  CheckCircle,
  Clock,
  Edit,
  Trash2,
  ExternalLink
} from 'lucide-react';

/**
 * Book Create Form Component
 */
export function BookCreateForm({ onSuccess, onError, onCancel, isLoading }: CrudFormProps) {
  const [formData, setFormData] = useState<BookCreateData>({
    title: '',
    author: '',
    isbn: '',
    description: '',
    price: 0,
    status: 'AVAILABLE',
    publishedAt: '',
    tags: [],
  });

  const [errors, setErrors] = useState<BookFormErrors>({});
  const [tagInput, setTagInput] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Validation function
  const validateForm = (): boolean => {
    const newErrors: BookFormErrors = {};

    // Title validation
    if (!formData.title.trim()) {
      newErrors.title = 'Title is required';
    } else if (formData.title.length > BOOK_VALIDATION_RULES.title.maxLength) {
      newErrors.title = `Title must be less than ${BOOK_VALIDATION_RULES.title.maxLength} characters`;
    }

    // Author validation
    if (!formData.author.trim()) {
      newErrors.author = 'Author is required';
    } else if (formData.author.length > BOOK_VALIDATION_RULES.author.maxLength) {
      newErrors.author = `Author must be less than ${BOOK_VALIDATION_RULES.author.maxLength} characters`;
    }

    // ISBN validation
    if (!formData.isbn.trim()) {
      newErrors.isbn = 'ISBN is required';
    } else if (!BOOK_VALIDATION_RULES.isbn.pattern.test(formData.isbn)) {
      newErrors.isbn = 'Invalid ISBN format (10-17 digits and dashes)';
    }

    // Price validation
    if (formData.price < BOOK_VALIDATION_RULES.price.min) {
      newErrors.price = 'Price must be greater than or equal to 0';
    }

    // Description validation
    if (formData.description && formData.description.length > BOOK_VALIDATION_RULES.description.maxLength) {
      newErrors.description = `Description must be less than ${BOOK_VALIDATION_RULES.description.maxLength} characters`;
    }

    // Tags validation
    if (formData.tags && formData.tags.length > BOOK_VALIDATION_RULES.tags.maxItems) {
      newErrors.tags = `Maximum ${BOOK_VALIDATION_RULES.tags.maxItems} tags allowed`;
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    setIsSubmitting(true);
    
    // Use fetch to handle the request manually since Inertia seems to have issues
    try {
      // Get CSRF token from meta tag
      const csrfToken = document.querySelector('meta[name="csrf-token"]')?.getAttribute('content');
      
      const response = await fetch('/api/books', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
          'X-Inertia': 'true',
          'X-Inertia-Version': '1.0.0',
          ...(csrfToken && { 'X-CSRF-TOKEN': csrfToken }),
        },
        body: JSON.stringify(formData),
      });

      const responseData = await response.json();
      console.log('Response status:', response.status);
      console.log('Response data:', responseData);

      if (response.ok) {
        onSuccess('Book created successfully');
      } else {
        // Handle validation errors
        setErrors(responseData.errors || {});
        onError?.(responseData);
      }
    } catch (error) {
      console.error('Error creating book:', error);
      const errorMessage = error instanceof Error ? error.message : 'Failed to create book. Please try again.';
      setErrors({ general: errorMessage });
      onError?.(error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const addTag = () => {
    if (tagInput.trim() && !formData.tags?.includes(tagInput.trim())) {
      const newTags = [...(formData.tags || []), tagInput.trim()];
      if (newTags.length <= BOOK_VALIDATION_RULES.tags.maxItems) {
        setFormData({ ...formData, tags: newTags });
        setTagInput('');
      }
    }
  };

  const removeTag = (index: number) => {
    const newTags = formData.tags?.filter((_, i) => i !== index) || [];
    setFormData({ ...formData, tags: newTags });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6 p-6">
      {errors.general && (
        <div className="bg-red-50 border border-red-200 rounded-md p-4">
          <div className="flex">
            <AlertCircle className="h-5 w-5 text-red-400" />
            <div className="ml-3">
              <h3 className="text-sm font-medium text-red-800">Error</h3>
              <div className="mt-2 text-sm text-red-700">{errors.general}</div>
            </div>
          </div>
        </div>
      )}

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <FormField label="Title" required error={errors.title}>
          <Input
            value={formData.title}
            onChange={(e) => setFormData({ ...formData, title: e.target.value })}
            placeholder="Enter book title"
            className="pl-10"
          />
          <BookOpen className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
        </FormField>

        <FormField label="Author" required error={errors.author}>
          <Input
            value={formData.author}
            onChange={(e) => setFormData({ ...formData, author: e.target.value })}
            placeholder="Enter author name"
            className="pl-10"
          />
          <User className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
        </FormField>

        <FormField label="ISBN" required error={errors.isbn}>
          <Input
            value={formData.isbn}
            onChange={(e) => setFormData({ ...formData, isbn: e.target.value })}
            placeholder="978-0123456789"
            className="pl-10"
          />
          <Hash className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
        </FormField>

        <FormField label="Price" required error={errors.price}>
          <Input
            type="number"
            step="0.01"
            min="0"
            value={formData.price}
            onChange={(e) => setFormData({ ...formData, price: parseFloat(e.target.value) || 0 })}
            placeholder="0.00"
            className="pl-10"
          />
          <DollarSign className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
        </FormField>

        <FormField label="Status">
          <select
            value={formData.status}
            onChange={(e) => setFormData({ ...formData, status: e.target.value as BookStatus })}
            className="w-full rounded-md border border-gray-300 px-3 py-2 focus:border-blue-500 focus:ring-blue-500"
          >
            {Object.entries(BOOK_STATUS_CONFIG).map(([value, config]) => (
              <option key={value} value={value}>
                {config.icon} {config.label}
              </option>
            ))}
          </select>
        </FormField>

        <FormField label="Published Date" error={errors.publishedAt}>
          <Input
            type="date"
            value={formData.publishedAt}
            onChange={(e) => setFormData({ ...formData, publishedAt: e.target.value })}
            className="pl-10"
          />
          <Calendar className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
        </FormField>
      </div>

      <FormField label="Description" error={errors.description}>
        <Textarea
          value={formData.description}
          onChange={(e) => setFormData({ ...formData, description: e.target.value })}
          placeholder="Enter book description"
          rows={3}
          className="resize-none"
        />
      </FormField>

      <FormField label="Tags" error={errors.tags}>
        <div className="space-y-3">
          <div className="flex space-x-2">
            <Input
              value={tagInput}
              onChange={(e) => setTagInput(e.target.value)}
              placeholder="Add a tag"
              onKeyPress={(e) => e.key === 'Enter' && (e.preventDefault(), addTag())}
              className="pl-10"
            />
            <Tag className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
            <Button type="button" onClick={addTag} variant="outline" size="sm">
              Add
            </Button>
          </div>
          {formData.tags && formData.tags.length > 0 && (
            <div className="flex flex-wrap gap-2">
              {formData.tags.map((tag, index) => (
                <Badge key={index} variant="outline" className="flex items-center gap-1">
                  {tag}
                  <button
                    type="button"
                    onClick={() => removeTag(index)}
                    className="ml-1 text-gray-500 hover:text-red-500"
                  >
                    ×
                  </button>
                </Badge>
              ))}
            </div>
          )}
        </div>
      </FormField>

      <div className="flex justify-end space-x-3 pt-6 border-t">
        <Button type="button" variant="outline" onClick={onCancel}>
          Cancel
        </Button>
        <Button type="submit" disabled={isLoading || isSubmitting}>
          {isLoading || isSubmitting ? 'Creating...' : 'Create Book'}
        </Button>
      </div>
    </form>
  );
}

/**
 * Book Edit Form Component
 */
export function BookEditForm({ item, onSuccess, onError, onCancel, isLoading }: CrudEditFormProps<Book>) {
  const [formData, setFormData] = useState<BookUpdateData>({
    title: item.title || '',
    author: item.author || '',
    isbn: item.isbn || '',
    description: item.description || '',
    price: item.price || 0,
    status: item.status || 'AVAILABLE',
    publishedAt: item.publishedAt ? item.publishedAt.split('T')[0] : '',
    tags: item.tags || [],
  });

  const [errors, setErrors] = useState<BookFormErrors>({});
  const [tagInput, setTagInput] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const validateForm = (): boolean => {
    const newErrors: BookFormErrors = {};

    if (formData.title && !formData.title.trim()) {
      newErrors.title = 'Title cannot be empty';
    } else if (formData.title && formData.title.length > BOOK_VALIDATION_RULES.title.maxLength) {
      newErrors.title = `Title must be less than ${BOOK_VALIDATION_RULES.title.maxLength} characters`;
    }

    if (formData.author && !formData.author.trim()) {
      newErrors.author = 'Author cannot be empty';
    } else if (formData.author && formData.author.length > BOOK_VALIDATION_RULES.author.maxLength) {
      newErrors.author = `Author must be less than ${BOOK_VALIDATION_RULES.author.maxLength} characters`;
    }

    if (formData.isbn && !BOOK_VALIDATION_RULES.isbn.pattern.test(formData.isbn)) {
      newErrors.isbn = 'Invalid ISBN format';
    }

    if (formData.price !== undefined && formData.price < BOOK_VALIDATION_RULES.price.min) {
      newErrors.price = 'Price must be greater than or equal to 0';
    }

    if (formData.description && formData.description.length > BOOK_VALIDATION_RULES.description.maxLength) {
      newErrors.description = `Description must be less than ${BOOK_VALIDATION_RULES.description.maxLength} characters`;
    }

    if (formData.tags && formData.tags.length > BOOK_VALIDATION_RULES.tags.maxItems) {
      newErrors.tags = `Maximum ${BOOK_VALIDATION_RULES.tags.maxItems} tags allowed`;
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    setIsSubmitting(true);
    
    try {
      // Get CSRF token from meta tag
      const csrfToken = document.querySelector('meta[name="csrf-token"]')?.getAttribute('content');
      
      const response = await fetch(`/api/books/${item.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
          'X-Inertia': 'true',
          'X-Inertia-Version': '1.0.0',
          ...(csrfToken && { 'X-CSRF-TOKEN': csrfToken }),
        },
        body: JSON.stringify(formData),
      });

      const responseData = await response.json();
      console.log('Update response status:', response.status);
      console.log('Update response data:', responseData);

      if (response.ok) {
        onSuccess('Book updated successfully');
      } else {
        // Handle validation errors
        setErrors(responseData.errors || {});
        onError?.(responseData);
      }
    } catch (error) {
      console.error('Error updating book:', error);
      const errorMessage = error instanceof Error ? error.message : 'Failed to update book. Please try again.';
      setErrors({ general: errorMessage });
      onError?.(error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const addTag = () => {
    if (tagInput.trim() && !formData.tags?.includes(tagInput.trim())) {
      const newTags = [...(formData.tags || []), tagInput.trim()];
      if (newTags.length <= BOOK_VALIDATION_RULES.tags.maxItems) {
        setFormData({ ...formData, tags: newTags });
        setTagInput('');
      }
    }
  };

  const removeTag = (index: number) => {
    const newTags = formData.tags?.filter((_, i) => i !== index) || [];
    setFormData({ ...formData, tags: newTags });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6 p-6">
      {errors.general && (
        <div className="bg-red-50 border border-red-200 rounded-md p-4">
          <div className="flex">
            <AlertCircle className="h-5 w-5 text-red-400" />
            <div className="ml-3">
              <div className="text-sm text-red-700">{errors.general}</div>
            </div>
          </div>
        </div>
      )}

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <FormField label="Title" error={errors.title}>
          <Input
            value={formData.title}
            onChange={(e) => setFormData({ ...formData, title: e.target.value })}
            placeholder="Enter book title"
          />
        </FormField>

        <FormField label="Author" error={errors.author}>
          <Input
            value={formData.author}
            onChange={(e) => setFormData({ ...formData, author: e.target.value })}
            placeholder="Enter author name"
          />
        </FormField>

        <FormField label="ISBN" error={errors.isbn}>
          <Input
            value={formData.isbn}
            onChange={(e) => setFormData({ ...formData, isbn: e.target.value })}
            placeholder="978-0123456789"
          />
        </FormField>

        <FormField label="Price" error={errors.price}>
          <Input
            type="number"
            step="0.01"
            min="0"
            value={formData.price}
            onChange={(e) => setFormData({ ...formData, price: parseFloat(e.target.value) || 0 })}
            placeholder="0.00"
          />
        </FormField>

        <FormField label="Status">
          <select
            value={formData.status}
            onChange={(e) => setFormData({ ...formData, status: e.target.value as BookStatus })}
            className="w-full rounded-md border border-gray-300 px-3 py-2 focus:border-blue-500 focus:ring-blue-500"
          >
            {Object.entries(BOOK_STATUS_CONFIG).map(([value, config]) => (
              <option key={value} value={value}>
                {config.icon} {config.label}
              </option>
            ))}
          </select>
        </FormField>

        <FormField label="Published Date" error={errors.publishedAt}>
          <Input
            type="date"
            value={formData.publishedAt}
            onChange={(e) => setFormData({ ...formData, publishedAt: e.target.value })}
          />
        </FormField>
      </div>

      <FormField label="Description" error={errors.description}>
        <Textarea
          value={formData.description}
          onChange={(e) => setFormData({ ...formData, description: e.target.value })}
          placeholder="Enter book description"
          rows={3}
        />
      </FormField>

      <FormField label="Tags" error={errors.tags}>
        <div className="space-y-3">
          <div className="flex space-x-2">
            <Input
              value={tagInput}
              onChange={(e) => setTagInput(e.target.value)}
              placeholder="Add a tag"
              onKeyPress={(e) => e.key === 'Enter' && (e.preventDefault(), addTag())}
            />
            <Button type="button" onClick={addTag} variant="outline" size="sm">
              Add
            </Button>
          </div>
          {formData.tags && formData.tags.length > 0 && (
            <div className="flex flex-wrap gap-2">
              {formData.tags.map((tag, index) => (
                <Badge key={index} variant="outline" className="flex items-center gap-1">
                  {tag}
                  <button
                    type="button"
                    onClick={() => removeTag(index)}
                    className="ml-1 text-gray-500 hover:text-red-500"
                  >
                    ×
                  </button>
                </Badge>
              ))}
            </div>
          )}
        </div>
      </FormField>

      <div className="flex justify-end space-x-3 pt-6 border-t">
        <Button type="button" variant="outline" onClick={onCancel}>
          Cancel
        </Button>
        <Button type="submit" disabled={isLoading || isSubmitting}>
          {isLoading || isSubmitting ? 'Updating...' : 'Update Book'}
        </Button>
      </div>
    </form>
  );
}

/**
 * Book Detail View Component
 */
export function BookDetailView({ 
  item: book, 
  onEdit, 
  onClose, 
  canEdit 
}: CrudDetailViewProps<Book>) {
  const [isPerformingAction, setIsPerformingAction] = useState(false);

  const handleBorrow = async () => {
    if (book.status !== 'AVAILABLE') return;
    
    setIsPerformingAction(true);
    try {
      await router.post(`/api/books/${book.id}/borrow`, {}, {
        onSuccess: () => {
          // The parent component should handle refresh
        },
        onError: (error) => {
          console.error('Error borrowing book:', error);
        },
        onFinish: () => {
          setIsPerformingAction(false);
        },
      });
    } catch (error) {
      console.error('Error borrowing book:', error);
      setIsPerformingAction(false);
    }
  };

  const handleReturn = async () => {
    if (book.status !== 'BORROWED') return;
    
    setIsPerformingAction(true);
    try {
      await router.post(`/api/books/${book.id}/return`, {}, {
        onSuccess: () => {
          // The parent component should handle refresh
        },
        onError: (error) => {
          console.error('Error returning book:', error);
        },
        onFinish: () => {
          setIsPerformingAction(false);
        },
      });
    } catch (error) {
      console.error('Error returning book:', error);
      setIsPerformingAction(false);
    }
  };

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex items-start space-x-4">
        <div className="flex-shrink-0">
          <div className="w-16 h-20 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center">
            <BookOpen className="w-8 h-8 text-white" />
          </div>
        </div>
        <div className="flex-grow">
          <h2 className="text-2xl font-bold text-gray-900 mb-1">{book.title}</h2>
          <p className="text-lg text-gray-600 mb-2">by {book.author}</p>
          <div className="flex items-center space-x-4">
            <StatusBadge 
              status={book.status}
              variant={
                book.status === 'AVAILABLE' ? 'success' :
                book.status === 'BORROWED' ? 'info' :
                'warning'
              }
            />
            <div className="text-xl font-semibold text-green-600">
              {new Intl.NumberFormat('en-US', {
                style: 'currency',
                currency: 'USD',
              }).format(book.price)}
            </div>
          </div>
        </div>
      </div>

      <Separator />

      {/* Book Information */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="space-y-4">
          <h3 className="text-lg font-semibold text-gray-900">Book Details</h3>
          
          <div className="space-y-3">
            <div className="flex items-center space-x-3">
              <Hash className="h-4 w-4 text-gray-400" />
              <span className="text-sm font-medium text-gray-500">ISBN:</span>
              <span className="text-sm text-gray-900 font-mono">{book.isbn}</span>
            </div>
            
            {book.publishedAt && (
              <div className="flex items-center space-x-3">
                <Calendar className="h-4 w-4 text-gray-400" />
                <span className="text-sm font-medium text-gray-500">Published:</span>
                <span className="text-sm text-gray-900">
                  {new Date(book.publishedAt).toLocaleDateString()}
                </span>
              </div>
            )}
            
            <div className="flex items-center space-x-3">
              <Clock className="h-4 w-4 text-gray-400" />
              <span className="text-sm font-medium text-gray-500">Added:</span>
              <span className="text-sm text-gray-900">
                {new Date(book.createdAt).toLocaleDateString()}
              </span>
            </div>
          </div>
        </div>

        <div className="space-y-4">
          <h3 className="text-lg font-semibold text-gray-900">Status Information</h3>
          
          <div className="space-y-3">
            <div className="flex items-center space-x-3">
              <div className="flex items-center space-x-2">
                {book.status === 'AVAILABLE' && <CheckCircle className="h-4 w-4 text-green-500" />}
                {book.status === 'BORROWED' && <Clock className="h-4 w-4 text-blue-500" />}
                {book.status === 'MAINTENANCE' && <AlertCircle className="h-4 w-4 text-orange-500" />}
                <span className="text-sm font-medium text-gray-900">
                  {BOOK_STATUS_CONFIG[book.status].label}
                </span>
              </div>
            </div>
            
            <p className="text-sm text-gray-600">
              {BOOK_STATUS_CONFIG[book.status].description}
            </p>

            {book.status === 'BORROWED' && book.dueDate && (
              <div className="bg-blue-50 rounded-lg p-3">
                <div className="text-sm font-medium text-blue-900">Due Date</div>
                <div className="text-sm text-blue-700">
                  {new Date(book.dueDate).toLocaleDateString()}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Description */}
      {book.description && (
        <>
          <Separator />
          <div>
            <h3 className="text-lg font-semibold text-gray-900 mb-3">Description</h3>
            <p className="text-sm text-gray-700 leading-relaxed">{book.description}</p>
          </div>
        </>
      )}

      {/* Tags */}
      {book.tags && book.tags.length > 0 && (
        <>
          <Separator />
          <div>
            <h3 className="text-lg font-semibold text-gray-900 mb-3">Tags</h3>
            <div className="flex flex-wrap gap-2">
              {book.tags.map((tag, index) => (
                <Badge key={index} variant="outline" className="flex items-center gap-1">
                  <Tag className="w-3 h-3" />
                  {tag}
                </Badge>
              ))}
            </div>
          </div>
        </>
      )}

      {/* Actions */}
      <div className="flex justify-between items-center pt-6 border-t">
        <div className="flex space-x-3">
          {book.status === 'AVAILABLE' && (
            <Button 
              onClick={handleBorrow} 
              disabled={isPerformingAction}
              className="bg-blue-600 hover:bg-blue-700"
            >
              {isPerformingAction ? 'Borrowing...' : 'Borrow Book'}
            </Button>
          )}
          
          {book.status === 'BORROWED' && (
            <Button 
              onClick={handleReturn} 
              disabled={isPerformingAction}
              className="bg-green-600 hover:bg-green-700"
            >
              {isPerformingAction ? 'Returning...' : 'Return Book'}
            </Button>
          )}
        </div>

        <div className="flex space-x-3">
          <Button type="button" variant="outline" onClick={onClose}>
            Close
          </Button>
          {canEdit && onEdit && (
            <Button type="button" onClick={onEdit}>
              <Edit className="w-4 h-4 mr-2" />
              Edit Book
            </Button>
          )}
        </div>
      </div>
    </div>
  );
}