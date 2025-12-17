import AuthLayout from '@/layouts/Auth';
import { LoginForm } from "@/components/login-form";
// @ts-ignore
import { Head, usePage } from "@inertiajs/react";
import type { SharedData } from "@/types/app";

export default function LoginPage() {
    const { appVersion } = usePage<SharedData>().props;

    return (
        <AuthLayout>
            <Head>
                <title>Login</title>
            </Head>
            <LoginForm />
            <div className="text-center">
                <p className="text-xs text-muted-foreground">
                    v{appVersion || '0.0.0'}
                </p>
                <p className="text-xs text-muted-foreground mt-1">
                    Property of Ministry of Trade and Industry and SMEDI. All rights reserved.
                </p>
            </div>
        </AuthLayout>
    );
}
