import React from 'react';
import { LucideIcon } from 'lucide-react';
import { Card, CardContent } from '@/components/ui/card';

interface EmptyStateProps {
  icon?: LucideIcon;
  title: string;
  description?: string;
  className?: string;
}

export const EmptyState: React.FC<EmptyStateProps> = ({
  icon: Icon,
  title,
  description,
  className = ''
}) => {
  return (
    <Card className={className}>
      <CardContent className="flex flex-col items-center justify-center py-12 text-center">
        {Icon && (
          <div className="mb-4 rounded-full bg-muted p-3">
            <Icon className="h-6 w-6 text-muted-foreground" />
          </div>
        )}
        <h3 className="text-lg font-semibold text-foreground mb-2">{title}</h3>
        {description && (
          <p className="text-sm text-muted-foreground max-w-sm">{description}</p>
        )}
      </CardContent>
    </Card>
  );
};
