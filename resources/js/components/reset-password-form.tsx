import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
// @ts-ignore
import { useForm, usePage } from '@inertiajs/react';
import React, { useState } from 'react';

export function ResetPasswordForm({ className, ...props }: React.ComponentPropsWithoutRef<"form">) {
    // Get token and email from query string via Inertia usePage
    const { url } = usePage();
    const params = new URLSearchParams(url.split('?')[1] || '');
    const token = params.get('token') || '';
    const email = params.get('email') || '';

    const { data, setData, post, processing, errors, reset } = useForm({
        email: email,
        token: token,
        password: '',
        password_confirmation: '',
    });
    const [localError, setLocalError] = useState<string | null>(null);

    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setLocalError(null);
        if (data.password !== data.password_confirmation) {
            setLocalError('Passwords do not match');
            return;
        }
        post('/reset-password');
    };

    return (
        <form onSubmit={handleSubmit} className={cn("flex flex-col gap-6 w-full max-w-md mx-auto", className)} {...props}>
            <div className="flex flex-col items-center gap-2 text-center">
                <h1 className="text-2xl font-bold">Reset your password</h1>
                <p className="text-balance text-sm text-muted-foreground">
                    Enter your new password below.
                </p>
            </div>
            <div className="grid gap-6">
                <div className="grid gap-2">
                    <Label htmlFor="password">New Password</Label>
                    <Input
                        id="password"
                        type="password"
                        value={data.password}
                        onChange={e => setData('password', e.target.value)}
                        required
                        disabled={processing}
                    />
                    {errors.password && <p className="text-xs text-red-500 mt-1">{errors.password}</p>}
                </div>
                <div className="grid gap-2">
                    <Label htmlFor="password_confirmation">Confirm Password</Label>
                    <Input
                        id="password_confirmation"
                        type="password"
                        value={data.password_confirmation}
                        onChange={e => setData('password_confirmation', e.target.value)}
                        required
                        disabled={processing}
                    />
                    {errors.password_confirmation && <p className="text-xs text-red-500 mt-1">{errors.password_confirmation}</p>}
                    {localError && <p className="text-xs text-red-500 mt-1">{localError}</p>}
                </div>
            </div>
            <Button type="submit" className="w-full" disabled={processing}>
                {processing ? 'Resetting...' : 'Reset Password'}
            </Button>
        </form>
    );
} 