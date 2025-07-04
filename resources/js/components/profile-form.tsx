import React from 'react';
// @ts-ignore
import { useForm } from "@inertiajs/react";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";
import type { User } from "@/types/app";
import { cn } from "@/lib/utils";

interface ProfileFormProps {
    user: User | null;
    className?: string;
}

export function ProfileForm({ user, className, ...props }: ProfileFormProps & React.ComponentPropsWithoutRef<"div">) {
    const { data, setData, put, processing, errors, reset } = useForm({
        current_password: '',
        password: '',
        password_confirmation: '',
    });

    const handlePasswordUpdate = (e: React.FormEvent) => {
        e.preventDefault();
        put('/profile/password', {
            onSuccess: () => {
                toast.success('Password updated successfully!');
                reset('current_password', 'password', 'password_confirmation');
            },
            onError: (errors: any) => {
                if (errors.current_password) toast.error(errors.current_password);
                if (errors.password) toast.error(errors.password);
                if (errors.password_confirmation) toast.error(errors.password_confirmation);
            }
        });
    };

    const getInitials = (name: string) => {
        return name
            .split(' ')
            .map(word => word.charAt(0))
            .join('')
            .toUpperCase()
            .slice(0, 2);
    };

    return (
        <div className={cn("w-full", className)} {...props}>
            {/* Header Section */}
            <div className="mb-8">
                <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
                    <div className="flex items-center gap-4">
                        <Avatar className="h-16 w-16">
                            <AvatarFallback className="text-lg font-semibold">
                                {user ? getInitials(user.name) : 'U'}
                            </AvatarFallback>
                        </Avatar>
                        <div>
                            <h1 className="text-2xl font-bold">{user?.name || 'User'}</h1>
                            <p className="text-muted-foreground">{user?.email}</p>
                        </div>
                    </div>
                    <div className="text-sm text-muted-foreground">
                        Member since {new Date().toLocaleDateString('en-US', {
                            month: 'long',
                            year: 'numeric'
                        })}
                    </div>
                </div>
            </div>

            <Separator className="mb-8" />

            {/* Main Content Grid */}
            <div className="space-y-8 lg:space-y-0 lg:grid lg:grid-cols-5 lg:gap-8">
                {/* Account Information - First on Mobile, Left Sidebar on Desktop */}
                <div className="lg:col-span-2 space-y-8">
                    {/* Account Information */}
                    <div className="space-y-6">
                        <div>
                            <h3 className="text-lg font-semibold">Account Information</h3>
                            <p className="text-sm text-muted-foreground mt-1">
                                Your account details and current status.
                            </p>
                        </div>
                        <div className="space-y-4">
                            <div className="flex items-center justify-between">
                                <span className="text-sm font-medium">Name</span>
                                <span className="text-sm text-muted-foreground">
                                    {user?.name || 'N/A'}
                                </span>
                            </div>

                            <div className="flex items-center justify-between">
                                <span className="text-sm font-medium">Email</span>
                                <span className="text-sm text-muted-foreground">
                                    {user?.email || 'N/A'}
                                </span>
                            </div>

                            <div className="flex items-center justify-between">
                                <span className="text-sm font-medium">Status</span>
                                <Badge variant="secondary" className="text-green-700 bg-green-50 dark:text-green-400 dark:bg-green-900/20">
                                    Active
                                </Badge>
                            </div>

                            <div className="flex items-center justify-between">
                                <span className="text-sm font-medium">Role</span>
                                <span className="text-sm text-muted-foreground capitalize">
                                    {user?.role?.toLowerCase() || 'User'}
                                </span>
                            </div>

                            <div className="flex items-center justify-between">
                                <span className="text-sm font-medium">User ID</span>
                                <span className="text-sm font-mono text-muted-foreground">
                                    {user?.id || 'N/A'}
                                </span>
                            </div>

                            <div className="flex items-center justify-between">
                                <span className="text-sm font-medium">Member Since</span>
                                <span className="text-sm text-muted-foreground">
                                    {new Date().toLocaleDateString('en-US', {
                                        month: 'short',
                                        year: 'numeric'
                                    })}
                                </span>
                            </div>
                        </div>
                    </div>

                    {/* Account Actions - Third on Mobile, Bottom of Sidebar on Desktop */}
                    <div className="lg:block hidden">
                        <Separator className="mb-6" />
                        <div className="space-y-6">
                            <div>
                                <h3 className="text-lg font-semibold">Account Actions</h3>
                                <p className="text-sm text-muted-foreground mt-1">
                                    Manage your account settings and data.
                                </p>
                            </div>
                            <Button
                                variant="outline"
                                className="w-full text-red-600 hover:text-red-700 hover:bg-red-50 dark:hover:bg-red-950/20"
                            >
                                Delete Account
                            </Button>
                        </div>
                    </div>
                </div>

                {/* Password Update Form - Second on Mobile, Right Side on Desktop */}
                <div className="lg:col-span-3">
                    <div className="max-w-2xl mx-auto">
                        <form onSubmit={handlePasswordUpdate} className="space-y-8">
                            {/* Security Section */}
                            <div className="space-y-6">
                                <div>
                                    <h3 className="text-lg font-semibold">Change Password</h3>
                                    <p className="text-sm text-muted-foreground mt-1">
                                        Update your account password. Leave fields empty to keep current password.
                                    </p>
                                </div>

                                <div className="space-y-6">
                                    <div className="space-y-2">
                                        <Label htmlFor="current_password">Current Password</Label>
                                        <Input
                                            id="current_password"
                                            type="password"
                                            value={data.current_password}
                                            onChange={(e) => setData('current_password', e.target.value)}
                                            className={errors.current_password ? "border-red-500" : ""}
                                            placeholder="Enter your current password"
                                        />
                                        {errors.current_password && (
                                            <p className="text-xs text-red-500">{errors.current_password}</p>
                                        )}
                                    </div>

                                    <div className="space-y-6">
                                        <div className="space-y-2">
                                            <Label htmlFor="password">New Password</Label>
                                            <Input
                                                id="password"
                                                type="password"
                                                value={data.password}
                                                onChange={(e) => setData('password', e.target.value)}
                                                className={errors.password ? "border-red-500" : ""}
                                                placeholder="Enter new password"
                                            />
                                            {errors.password && (
                                                <p className="text-xs text-red-500">{errors.password}</p>
                                            )}
                                        </div>
                                        <div className="space-y-2">
                                            <Label htmlFor="password_confirmation">Confirm New Password</Label>
                                            <Input
                                                id="password_confirmation"
                                                type="password"
                                                value={data.password_confirmation}
                                                onChange={(e) => setData('password_confirmation', e.target.value)}
                                                className={errors.password_confirmation ? "border-red-500" : ""}
                                                placeholder="Confirm new password"
                                            />
                                            {errors.password_confirmation && (
                                                <p className="text-xs text-red-500">{errors.password_confirmation}</p>
                                            )}
                                        </div>
                                    </div>
                                </div>
                            </div>

                            {/* Form Actions */}
                            <div className="flex flex-col sm:flex-row gap-3">
                                <Button type="submit" disabled={processing} className="sm:w-auto">
                                    {processing ? 'Updating Password...' : 'Update Password'}
                                </Button>
                                <Button
                                    type="button"
                                    variant="outline"
                                    onClick={() => reset()}
                                    disabled={processing}
                                    className="sm:w-auto"
                                >
                                    Cancel
                                </Button>
                            </div>
                        </form>
                    </div>
                </div>

                {/* Account Actions - Third on Mobile Only */}
                <div className="lg:hidden">
                    <div className="space-y-6">
                        <div>
                            <h3 className="text-lg font-semibold">Account Actions</h3>
                            <p className="text-sm text-muted-foreground mt-1">
                                Manage your account settings and data.
                            </p>
                        </div>
                        <Button
                            variant="outline"
                            className="w-full text-red-600 hover:text-red-700 hover:bg-red-50 dark:hover:bg-red-950/20"
                        >
                            Delete Account
                        </Button>
                    </div>
                </div>
            </div>
        </div>
    );
} 