import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
// @ts-ignore
import { useForm, usePage, Link } from '@inertiajs/react';
import React from 'react';
import { toast } from "sonner";

export function ForgotPasswordForm({
    className,
    ...props
}: React.ComponentPropsWithoutRef<"form">) {
    const { data, setData, post, processing, errors, reset } = useForm({
        email: '',
    });
    const page = usePage();
    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        post('/forgot-password', {
            onFinish: () => reset('email'), // Optional: reset email field on finish
            onError: (errors: any) => {
                if (errors.email) {
                    toast.error(errors.email);
                }
                if (errors.general) {
                    toast.error(errors.general);
                }
            }
        });
    };

    return (
        <form onSubmit={handleSubmit} className={cn("flex flex-col gap-6", className)} {...props}>
            <div className="flex flex-col items-center gap-2 text-center">
                <h1 className="text-2xl font-bold">Forgot your password?</h1>
                <p className="text-balance text-sm text-muted-foreground">
                    Enter your email address and we'll send you a link to reset your password
                </p>
            </div>
            {/* Display general errors */}
            {page.props.errors?.general && (
                <div className="text-sm text-red-500 text-center bg-red-50 p-2 rounded">
                    {page.props.errors.general}
                </div>
            )}
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
                    {(errors.email || page.props.errors?.email) && (
                        <p className="text-xs text-red-500 mt-1">
                            {errors.email || page.props.errors?.email}
                        </p>
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