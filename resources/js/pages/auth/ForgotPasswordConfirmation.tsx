import AuthLayout from '@/layouts/Auth';
import { ForgotPasswordConfirmation } from "@/components/forgot-password-confirmation";
// @ts-ignore
import { Head } from "@inertiajs/react";

export default function ForgotPasswordConfirmationPage() {
    return (
        <AuthLayout>
            <Head>
                <title>Check Your Email</title>
            </Head>
            <ForgotPasswordConfirmation />
        </AuthLayout>
    );
}
