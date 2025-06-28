import React, { Fragment } from 'react';
import { X } from 'lucide-react';
import { cn } from '@/lib/utils';
import { DrawerProps } from '@/types/crud';
import { Button } from '@/components/ui/button';
import { Sheet, SheetContent, SheetHeader, SheetTitle } from '@/components/ui/sheet';

export function CrudDrawer({ 
  isOpen, 
  onClose, 
  title, 
  size = 'md', 
  children,
  className,
  overlayClassName 
}: DrawerProps) {
  const sizeClasses = {
    sm: 'max-w-md',
    md: 'max-w-lg',
    lg: 'max-w-2xl',
    xl: 'max-w-4xl',
    full: 'max-w-full',
  };

  return (
    <Sheet open={isOpen} onOpenChange={onClose}>
      <SheetContent 
        side="right"
        className={cn(
          'flex flex-col h-full',
          sizeClasses[size],
          className
        )}
      >
        <SheetHeader className="border-b border-gray-200 pb-4">
          <div className="flex items-center justify-between">
            <SheetTitle className="text-lg font-semibold text-gray-900">
              {title}
            </SheetTitle>
            <Button
              variant="ghost"
              size="sm"
              onClick={onClose}
              className="h-8 w-8 p-0 hover:bg-gray-100"
            >
              <X className="h-4 w-4" />
              <span className="sr-only">Close</span>
            </Button>
          </div>
        </SheetHeader>

        <div className="flex-1 overflow-y-auto py-6">
          {children}
        </div>
      </SheetContent>
    </Sheet>
  );
}