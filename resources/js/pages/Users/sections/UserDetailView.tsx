import React from 'react';
import { 
  Calendar, 
  Mail, 
  Shield, 
  Clock, 
  User as UserIcon,
  Briefcase,
  CalendarDays,
  CheckCircle,
  XCircle
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Separator } from '@/components/ui/separator';
import { CrudDetailViewProps } from '@/types/crud';
import { User } from '@/types/user';
import { Badge } from '@/components/ui/badge';

export function UserDetailView({ 
  item: user,
  onEdit,
  onClose,
  canEdit
}: CrudDetailViewProps<User>) {
  const getInitials = (name: string) => {
    return name
      .split(' ')
      .map(n => n[0])
      .join('')
      .toUpperCase()
      .slice(0, 2);
  };

  const formatDate = (date: string | Date) => {
    return new Date(date).toLocaleDateString('en-US', {
      month: 'long',
      day: 'numeric',
      year: 'numeric'
    });
  };

  return (
    <div className="space-y-8">
      {/* Profile Picture Section */}
      <div className="flex justify-center">
        <Avatar className="h-32 w-32 ring-4 ring-background shadow-xl">
          <AvatarImage src={`https://api.dicebear.com/7.x/avataaars/svg?seed=${user.email}`} alt={user.name} />
          <AvatarFallback className="text-2xl bg-gradient-to-br from-teal-500 to-teal-600 text-white">
            {getInitials(user.name)}
          </AvatarFallback>
        </Avatar>
      </div>

      {/* User Information Section */}
      <div className="space-y-6">
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">User Information</h3>
          <div className="space-y-4">
            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <UserIcon className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">User ID</p>
                <p className="font-medium text-foreground">#{user.id.toString().padStart(6, '0')}</p>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <UserIcon className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">Full Name</p>
                <p className="font-medium text-foreground">{user.name}</p>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <Mail className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">Email Address</p>
                <p className="font-medium text-foreground">{user.email}</p>
                {user.email_verified && (
                  <Badge variant="secondary" className="mt-1">
                    <CheckCircle className="h-3 w-3 mr-1" />
                    Verified
                  </Badge>
                )}
              </div>
            </div>

          </div>
        </div>

        <Separator className="my-6" />

        {/* Account Details */}
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">Account Details</h3>
          <div className="space-y-4">
            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <Briefcase className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">Roles</p>
                <div className="flex flex-wrap gap-2 mt-1">
                  {user.roles && user.roles.length > 0 ? (
                    user.roles.map((role) => (
                      <Badge key={role.id} variant="secondary" className="bg-teal-100 text-teal-800 dark:bg-teal-900/30 dark:text-teal-400">
                        {role.name}
                      </Badge>
                    ))
                  ) : (
                    <span className="text-sm text-muted-foreground">No roles assigned</span>
                  )}
                </div>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <Clock className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">Account Status</p>
                <Badge variant={user.is_active ? 'default' : 'secondary'} className={user.is_active ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400' : ''}>
                  {user.is_active ? (
                    <>
                      <CheckCircle className="h-3 w-3 mr-1" />
                      Active
                    </>
                  ) : (
                    <>
                      <XCircle className="h-3 w-3 mr-1" />
                      Inactive
                    </>
                  )}
                </Badge>
              </div>
            </div>

            {user.is_super_admin && (
              <div className="flex items-start gap-3">
                <div className="p-2 rounded-lg bg-muted">
                  <Shield className="h-4 w-4 text-muted-foreground" />
                </div>
                <div className="flex-1 space-y-1">
                  <p className="text-sm text-muted-foreground">Admin Status</p>
                  <Badge className="bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400">
                    <Shield className="h-3 w-3 mr-1" />
                    Super Administrator
                  </Badge>
                </div>
              </div>
            )}

            {user.last_login_at && (
              <div className="flex items-start gap-3">
                <div className="p-2 rounded-lg bg-muted">
                  <Clock className="h-4 w-4 text-muted-foreground" />
                </div>
                <div className="flex-1 space-y-1">
                  <p className="text-sm text-muted-foreground">Last Login</p>
                  <p className="font-medium text-foreground">{formatDate(user.last_login_at)}</p>
                </div>
              </div>
            )}

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <CalendarDays className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">Member Since</p>
                <p className="font-medium text-foreground">{formatDate(user.created_at)}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}