import AuthLayout from '@/layouts/Auth';
import { ForgotPasswordForm } from "@/components/forgot-password-form";
// @ts-ignore
import { Head } from "@inertiajs/react";

export default function ForgotPasswordPage() {
    return (
        <AuthLayout>
            <Head>
                <title>Forgot Password</title>
            </Head>
            <ForgotPasswordForm />
        </AuthLayout>
    );
}