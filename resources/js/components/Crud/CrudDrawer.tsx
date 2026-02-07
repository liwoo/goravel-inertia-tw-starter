import React, { useEffect } from 'react';
import { ArrowLeft, Edit, Save, Command, X } from 'lucide-react';
import { cn } from '@/lib/utils';
import { DrawerProps } from '@/types/crud';
import { Button } from '@/components/ui/button';
import { useTranslation } from 'react-i18next';
import {
  Sheet,
  SheetContent,
  SheetTitle,
} from '@/components/ui/sheet';
import {
  Dialog,
  DialogContent,
  DialogTitle,
} from '@/components/ui/dialog';
import { ScrollArea } from '@/components/ui/scroll-area';

interface CrudDrawerProps extends DrawerProps {
  type?: 'create' | 'edit' | 'view';
  showEditButton?: boolean;
  onEdit?: () => void;
  onSave?: () => void;
  canEdit?: boolean;
  canSave?: boolean;
  isSaving?: boolean;
  resourceName?: string;
  displayName?: string; // Optional display name override for cleaner titles
  fullscreen?: boolean; // Use fullscreen dialog instead of sidebar drawer
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
  resourceName = '',
  displayName,
  fullscreen = false
}: CrudDrawerProps) {
  const { t } = useTranslation(['common', 'crud']);

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
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [isOpen, type, canSave, canEdit, onSave, onEdit]);

  // Map size to width classes
  const sizeClasses = {
    sm: 'w-full sm:max-w-[400px]',
    md: 'w-full sm:max-w-[600px]',
    lg: 'w-full sm:max-w-[800px]',
    xl: 'w-full sm:max-w-[1000px]',
    full: 'w-full',
  };

  // Use displayName if provided, otherwise format resourceName
  const formattedName = displayName || (resourceName ?
    resourceName.charAt(0).toUpperCase() + resourceName.slice(1, -1).replace(/_/g, ' ') : // Remove 's' and replace underscores
    '');

  // Build title based on type
  const displayTitle = type && formattedName ?
    (type === 'create' ? t('crud:drawer.createNew', {name: formattedName}) :
     type === 'edit' ? t('crud:drawer.edit', {name: formattedName}) :
     type === 'view' ? t('crud:drawer.details', {name: formattedName}) : title) : title;

  // Fullscreen mode uses Dialog
  if (fullscreen) {
    return (
      <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
        <DialogContent
          showCloseButton={false}
          className={cn(
            "fixed inset-4 w-auto max-w-none h-auto translate-x-0 translate-y-0 top-0 left-0 p-0 flex flex-col gap-0",
            className
          )}
        >
          {/* Header */}
          <div className="flex items-center justify-between px-4 sm:px-6 py-4 border-b border-border shrink-0">
            <div className="flex items-center gap-4">
              <Button
                variant="ghost"
                size="icon"
                onClick={onClose}
                className="h-8 w-8"
              >
                <X className="h-4 w-4" />
              </Button>
              <DialogTitle className="text-lg font-semibold">
                {displayTitle}
              </DialogTitle>
            </div>

            <div className="flex items-center gap-2">
              {/* Save button for Create/Edit modes */}
              {(type === 'create' || type === 'edit') && canSave && onSave && (
                <Button
                  onClick={onSave}
                  size="sm"
                  disabled={isSaving}
                  className="bg-primary hover:bg-primary/90 text-primary-foreground group"
                >
                  <Save className="h-4 w-4 mr-2" />
                  {isSaving ? t('common:actions.saving') : t('common:actions.save')}
                  {!isSaving && (
                    <kbd className="ml-2 pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border bg-primary-foreground/20 px-1.5 font-mono text-[10px] font-medium opacity-100 group-hover:bg-primary-foreground/30">
                      <Command className="h-3 w-3" />S
                    </kbd>
                  )}
                </Button>
              )}
            </div>
          </div>

          {/* Scrollable Content */}
          <ScrollArea className="flex-1">
            <div className="p-4 sm:p-6">
              {children}
            </div>
          </ScrollArea>
        </DialogContent>
      </Dialog>
    );
  }

  // Default sidebar mode uses Sheet
  return (
    <Sheet open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <SheetContent
        side="right"
        className={cn(
          sizeClasses[size],
          "p-0 flex flex-col gap-0 [&>button]:hidden", // Hide default close button
          className
        )}
      >
        {/* Custom Header with actions */}
        <div className="flex items-center justify-between px-4 sm:px-6 py-4 border-b border-border shrink-0">
          <div className="flex items-center gap-4">
            <Button
              variant="ghost"
              size="icon"
              onClick={onClose}
              className="h-8 w-8"
            >
              <ArrowLeft className="h-4 w-4" />
            </Button>
            <SheetTitle className="text-lg font-semibold">
              {displayTitle}
            </SheetTitle>
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
                {t('common:actions.edit')}
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
                {isSaving ? t('common:actions.saving') : t('common:actions.save')}
                {!isSaving && (
                  <kbd className="ml-2 pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border bg-primary-foreground/20 px-1.5 font-mono text-[10px] font-medium opacity-100 group-hover:bg-primary-foreground/30">
                    <Command className="h-3 w-3" />S
                  </kbd>
                )}
              </Button>
            )}
          </div>
        </div>

        {/* Scrollable Content */}
        <ScrollArea className="flex-1">
          <div className="p-4 sm:p-6">
            {children}
          </div>
        </ScrollArea>
      </SheetContent>
    </Sheet>
  );
}
