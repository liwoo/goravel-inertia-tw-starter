import React, { useEffect } from 'react';
import { ArrowLeft, Edit, Save, Command } from 'lucide-react';
import { cn } from '@/lib/utils';
import { DrawerProps } from '@/types/crud';
import { Button } from '@/components/ui/button';

interface CrudDrawerProps extends DrawerProps {
  type?: 'create' | 'edit' | 'view';
  showEditButton?: boolean;
  onEdit?: () => void;
  onSave?: () => void;
  canEdit?: boolean;
  canSave?: boolean;
  isSaving?: boolean;
  resourceName?: string;
}

export function CrudDrawer({ 
  isOpen, 
  onClose, 
  title, 
  size = 'lg', 
  children,
  className,
  overlayClassName,
  type,
  showEditButton = false,
  onEdit,
  onSave,
  canEdit = false,
  canSave = false,
  isSaving = false,
  resourceName = ''
}: CrudDrawerProps) {
  // Keyboard shortcuts
  useEffect(() => {
    if (!isOpen) return;

    const handleKeyDown = (e: KeyboardEvent) => {
      // Cmd/Ctrl + S to save
      if ((e.metaKey || e.ctrlKey) && e.key === 's' && canSave && onSave) {
        e.preventDefault();
        onSave();
      }
      
      // Cmd/Ctrl + E to edit (from view mode)
      if ((e.metaKey || e.ctrlKey) && e.key === 'e' && type === 'view' && canEdit && onEdit) {
        e.preventDefault();
        onEdit();
      }
      
      // Escape to close
      if (e.key === 'Escape') {
        e.preventDefault();
        onClose();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [isOpen, type, canSave, canEdit, onSave, onEdit, onClose]);
  const sizeClasses = {
    sm: 'w-[400px]',
    md: 'w-[600px]',
    lg: 'w-[800px]',
    xl: 'w-[1000px]',
    full: 'w-full',
  };

  if (!isOpen) return null;

  // Capitalize resource name for title
  const capitalizedResourceName = resourceName ? 
    resourceName.charAt(0).toUpperCase() + resourceName.slice(1, -1) : // Remove 's' at the end
    '';

  // Build title based on type
  const displayTitle = type && capitalizedResourceName ? 
    (type === 'create' ? `Create New ${capitalizedResourceName}` :
     type === 'edit' ? `Edit ${capitalizedResourceName}` :
     type === 'view' ? `${capitalizedResourceName} Details` : title) : title;

  return (
    <>
      {/* Overlay */}
      <div 
        className={cn(
          "fixed inset-0 bg-black/20 z-40 transition-opacity",
          overlayClassName
        )}
        onClick={onClose}
      />
      
      {/* Drawer */}
      <div 
        className={cn(
          "fixed right-0 top-0 h-full bg-white dark:bg-gray-950 shadow-xl z-50 flex flex-col",
          "transform transition-transform duration-300 ease-in-out",
          isOpen ? "translate-x-0" : "translate-x-full",
          sizeClasses[size],
          className
        )}
      >
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b border-gray-200 dark:border-gray-800">
          <div className="flex items-center gap-4">
            <Button
              variant="ghost"
              size="icon"
              onClick={onClose}
              className="h-8 w-8"
            >
              <ArrowLeft className="h-4 w-4" />
            </Button>
            <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              {displayTitle}
            </h2>
          </div>
          
          <div className="flex items-center gap-2">
            {/* Edit button for View mode */}
            {type === 'view' && canEdit && onEdit && (
              <Button 
                onClick={onEdit}
                size="sm"
                variant="outline"
                className="group"
              >
                <Edit className="h-4 w-4 mr-2" />
                Edit
                <kbd className="ml-2 pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium opacity-100 group-hover:bg-background">
                  <Command className="h-3 w-3" />E
                </kbd>
              </Button>
            )}
            
            {/* Save button for Create/Edit modes */}
            {(type === 'create' || type === 'edit') && canSave && onSave && (
              <Button 
                onClick={onSave}
                size="sm"
                disabled={isSaving}
                className="bg-primary hover:bg-primary/90 text-primary-foreground group"
              >
                <Save className="h-4 w-4 mr-2" />
                {isSaving ? 'Saving...' : 'Save'}
                {!isSaving && (
                  <kbd className="ml-2 pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border bg-primary-foreground/20 px-1.5 font-mono text-[10px] font-medium opacity-100 group-hover:bg-primary-foreground/30">
                    <Command className="h-3 w-3" />S
                  </kbd>
                )}
              </Button>
            )}
          </div>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-y-auto">
          <div className="p-6">
            {children}
          </div>
        </div>
      </div>
    </>
  );
}