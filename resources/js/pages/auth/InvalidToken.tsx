import AuthLayout from '@/layouts/Auth';
// @ts-ignore
import { Head, Link } from "@inertiajs/react";
import { AlertCircle} from 'lucide-react';
import { Button } from '@/components/ui/button';

export default function InvalidTokenPage() {
    return (
        <AuthLayout>
            <Head>
                <title>Invalid Token</title>
            </Head>
            <div className="flex flex-col gap-6">
                <div className="flex flex-col items-center gap-2 text-center">
                    <div className="flex justify-center mb-4">
                        <AlertCircle className="h-12 w-12 text-destructive" />
                    </div>
                    <h1 className="text-2xl font-bold">Invalid or Expired Token</h1>
                    <p className="text-balance text-sm text-muted-foreground">
                        The password reset link you used is either invalid or has expired.
                    </p>
                </div>

                <div className="grid gap-6">
                    <div className="text-center text-sm text-muted-foreground space-y-1">
                        <p>Password reset tokens expire after 1 hour for security reasons.</p>
                        <p>Please request a new password reset link to continue.</p>
                    </div>

                    <Button asChild className="w-full">
                        <Link href="/forgot-password">
                            Request New Reset Link
                        </Link>
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
            </div>
        </AuthLayout>
    );
} 