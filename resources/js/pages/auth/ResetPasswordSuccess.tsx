import { Button } from '@/components/ui/button';
import AuthLayout from '@/layouts/Auth';
// @ts-ignore
import { Head, Link } from "@inertiajs/react";

export default function ResetPasswordSuccessPage() {
    return (
        <AuthLayout>
            <Head>
                <title>Password Reset Successful</title>
            </Head>
            <div className="flex flex-col gap-6 items-center justify-center min-h-[60vh]">
                <div className="flex flex-col items-center gap-2 text-center w-full max-w-md mx-auto">
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
                    <h1 className="text-2xl font-bold">Password Reset Successful</h1>
                    <p className="text-balance text-sm text-muted-foreground mb-2">
                        Your password has been successfully reset. You can now log in with your new password.
                    </p>
                </div>
                <div className="w-full max-w-md mx-auto">
                    <Button asChild className="w-full">
                        <Link href="/login">
                            Continue to Login
                        </Link>
                    </Button>
                </div>
            </div>
        </AuthLayout>
    );
} 