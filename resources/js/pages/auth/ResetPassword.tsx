import { ResetPasswordForm } from '@/components/reset-password-form';
import AuthLayout from '@/layouts/Auth';
// @ts-ignore
import { Head } from "@inertiajs/react";

export default function ResetPasswordPage() {
    return (
        <AuthLayout>
            <Head>
                <title>Reset Password</title>
            </Head>
            <ResetPasswordForm />
        </AuthLayout>
    );
} 