import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { cn } from '@/lib/utils';

interface StatsCardProps {
  title: string;
  value: string | number;
  icon?: React.ReactNode;
  description?: string;
  trend?: {
    value: number;
    label: string;
  };
  className?: string;
  iconClassName?: string;
  valueClassName?: string;
}

export function StatsCard({
  title,
  value,
  icon,
  description,
  trend,
  className,
  iconClassName,
  valueClassName,
}: StatsCardProps) {
  return (
    <Card className={cn("bg-gradient-to-br from-primary/5 to-card shadow-xs", className)}>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-base font-medium">{title}</CardTitle>
        {icon && (
          <div className={cn("h-4 w-4 text-muted-foreground", iconClassName)}>
            {icon}
          </div>
        )}
      </CardHeader>
      <CardContent>
        <div className={cn("text-2xl font-bold", valueClassName)}>{value}</div>
        {description && (
          <p className="text-xs text-muted-foreground">{description}</p>
        )}
        {trend && (
          <p className="text-xs text-muted-foreground">
            {trend.value > 0 ? '+' : ''}{trend.value}% {trend.label}
          </p>
        )}
      </CardContent>
    </Card>
  );
}