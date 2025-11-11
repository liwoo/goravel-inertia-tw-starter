import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Spinner } from "@/components/ui/spinner"
// @ts-ignore
import { Link, useForm, usePage } from '@inertiajs/react';
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
            },
        });
    };

    return (
        <form onSubmit={handleSubmit} className={cn("flex flex-col gap-6 pb-8", className)} {...props}>
            <div className="flex flex-col items-start gap-2">
                <div className="flex items-center gap-3 w-full">
                    <img src="/images/mw-coat.svg" alt="MW Gov Emblem" className="w-1/5" />
                    <div className="flex flex-col items-start">
                        <h3 className="text-xl font-semibold uppercase text-nowrap">Forgot Password</h3>
                        <h4>Get Back Into Your Account</h4>
                    </div>
                </div>
                <p className="text-sm text-muted-foreground">
                    Enter your email and we'll send you a link to reset your password
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
                            <Spinner className="mr-2" />
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