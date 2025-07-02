import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
// @ts-ignore
import { useForm, Link } from '@inertiajs/react';
import React, { useState } from 'react';

export function ForgotPasswordForm({
    className,
    ...props
}: React.ComponentPropsWithoutRef<"form">) {
    const [isSubmitted, setIsSubmitted] = useState(false);

    const { data, setData, post, processing, errors, reset } = useForm({
        email: '',
    });

    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        post('/forgot-password', {
            onSuccess: () => {
                setIsSubmitted(true);
                reset('email');
            },
            onError: () => {
                // Handle errors if needed
            }
        });
    };

    if (isSubmitted) {
        return (
            <div className={cn("flex flex-col gap-6", className)}>
                <Card className="w-full max-w-md mx-auto">
                    <CardHeader className="text-center">
                        <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-green-100">
                            <svg
                                className="h-6 w-6 text-green-600"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                    strokeWidth={2}
                                    d="M5 13l4 4L19 7"
                                />
                            </svg>
                        </div>
                        <CardTitle className="text-2xl font-bold">Check your email</CardTitle>
                        <CardDescription>
                            We've sent a password reset link to your email address. Please check your inbox and follow the instructions to reset your password.
                        </CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        <div className="text-center text-sm text-muted-foreground">
                            Didn't receive the email? Check your spam folder or{" "}
                            <button
                                onClick={() => setIsSubmitted(false)}
                                className="text-primary underline underline-offset-4 hover:no-underline"
                            >
                                try again
                            </button>
                        </div>
                        <div className="text-center">
                            <Link
                                href="/login"
                                className="text-sm text-primary underline underline-offset-4 hover:no-underline"
                            >
                                Back to login
                            </Link>
                        </div>
                    </CardContent>
                </Card>
            </div>
        );
    }

    return (
        <form onSubmit={handleSubmit} className={cn("flex flex-col gap-6", className)} {...props}>
            <div className="flex flex-col items-center gap-2 text-center">
                <h1 className="text-2xl font-bold">Forgot your password?</h1>
                <p className="text-balance text-sm text-muted-foreground">
                    Enter your email address and we'll send you a link to reset your password
                </p>
            </div>

            <div className="grid gap-6">
                <div className="grid gap-2">
                    <Label htmlFor="email">Email</Label>
                    <Input
                        id="email"
                        type="email"
                        placeholder="m@example.com"
                        value={data.email}
                        onChange={(e) => setData('email', e.target.value)}
                        required
                        disabled={processing}
                    />
                    {errors.email && (
                        <p className="text-xs text-red-500 mt-1">{errors.email}</p>
                    )}
                </div>

                <Button type="submit" className="w-full" disabled={processing}>
                    {processing ? (
                        <>
                            <svg
                                className="mr-2 h-4 w-4 animate-spin"
                                xmlns="http://www.w3.org/2000/svg"
                                fill="none"
                                viewBox="0 0 24 24"
                            >
                                <circle
                                    className="opacity-25"
                                    cx="12"
                                    cy="12"
                                    r="10"
                                    stroke="currentColor"
                                    strokeWidth="4"
                                />
                                <path
                                    className="opacity-75"
                                    fill="currentColor"
                                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                                />
                            </svg>
                            Sending reset link...
                        </>
                    ) : (
                        'Send reset link'
                    )}
                </Button>
            </div>

            <div className="text-center text-sm">
                Remember your password?{" "}
                <Link
                    href="/login"
                    className="text-primary underline underline-offset-4 hover:no-underline"
                >
                    Back to login
                </Link>
            </div>
        </form>
    );
}