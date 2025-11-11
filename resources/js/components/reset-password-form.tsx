import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Spinner } from "@/components/ui/spinner";
// @ts-ignore
import { useForm, usePage } from '@inertiajs/react';
import React, { useState } from 'react';
import { toast } from "sonner";

export function ResetPasswordForm({ className, ...props }: React.ComponentPropsWithoutRef<"form">) {
    // Get token and email from query string via Inertia usePage
    const page = usePage();
    const { url } = page;
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
    const [shownErrorKeys, setShownErrorKeys] = useState<Set<string>>(new Set());

    // Show toasts for new errors
    const errorsObj = errors as Record<string, string>;
    const currentErrorKeys = Object.keys(errorsObj).filter(key => errorsObj[key]);
    const newErrorKeys = currentErrorKeys.filter(key => !shownErrorKeys.has(key));

    if (newErrorKeys.length > 0 && !processing) {
        newErrorKeys.forEach(key => {
            const errorMessage = errorsObj[key];
            if (errorMessage) {
                toast.error(String(errorMessage));
            }
        });
        setShownErrorKeys(new Set([...shownErrorKeys, ...newErrorKeys]));
    }

    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setLocalError(null);
        setShownErrorKeys(new Set()); // Reset shown errors on new submission

        if (data.password !== data.password_confirmation) {
            setLocalError('Passwords do not match');
            toast.error('Passwords do not match');
            return;
        }

        post('/reset-password', {
            onFinish: () => reset('password', 'password_confirmation'),
        });
    };

    return (
        <form onSubmit={handleSubmit} className={cn("flex flex-col gap-6 pb-8", className)} {...props}>
            <div className="flex flex-col items-start gap-2">
                <div className="flex items-center gap-3 w-full">
                    <img src="/images/mw-coat.svg" alt="MW Gov Emblem" className="w-1/5" />
                    <div className="flex flex-col items-start">
                        <h3 className="text-xl font-semibold uppercase text-nowrap">Reset Password</h3>
                        <h4>Create Your New Password</h4>
                    </div>
                </div>
                <p className="text-sm text-muted-foreground">
                    Enter your new password below
                </p>
            </div>

            {/* Display general errors */}
            {(page.props.errors?.general || page.props.errors?.message) && (
                <div className="text-sm text-red-500 text-center bg-red-50 p-2 rounded">
                    {page.props.errors?.general || page.props.errors?.message}
                </div>
            )}

            {/* Display token errors */}
            {(errors.token || page.props.errors?.token) && (
                <div className="text-sm text-red-500 text-center bg-red-50 p-2 rounded">
                    {errors.token || page.props.errors?.token}
                </div>
            )}

            {/* Display email errors */}
            {(errors.email || page.props.errors?.email) && (
                <div className="text-sm text-red-500 text-center bg-red-50 p-2 rounded">
                    {errors.email || page.props.errors?.email}
                </div>
            )}

            <div className="grid gap-6">
                <div className="grid gap-2">
                    <Label htmlFor="password">New Password</Label>
                    <Input
                        id="password"
                        type="password"
                        value={data.password}
                        onChange={e => setData('password', e.target.value)}
                        required
                    />
                    {(errors.password || page.props.errors?.password) && (
                        <p className="text-xs text-red-500 mt-1">
                            {errors.password || page.props.errors?.password}
                        </p>
                    )}
                </div>
                <div className="grid gap-2">
                    <Label htmlFor="password_confirmation">Confirm Password</Label>
                    <Input
                        id="password_confirmation"
                        type="password"
                        value={data.password_confirmation}
                        onChange={e => setData('password_confirmation', e.target.value)}
                        required
                    />
                    {(errors.password_confirmation || page.props.errors?.password_confirmation) && (
                        <p className="text-xs text-red-500 mt-1">
                            {errors.password_confirmation || page.props.errors?.password_confirmation}
                        </p>
                    )}
                    {localError && <p className="text-xs text-red-500 mt-1">{localError}</p>}
                </div>
            </div>
            <Button type="submit" className="w-full" disabled={processing}>
                {processing ? (
                    <>
                        <Spinner className="mr-2" />
                        Resetting password...
                    </>
                ) : (
                    'Reset Password'
                )}
            </Button>
        </form>
    );
} 