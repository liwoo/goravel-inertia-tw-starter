import AuthLayout from '@/layouts/Auth';
import { LoginForm } from "@/components/login-form";
// @ts-ignore
import { Head, usePage } from "@inertiajs/react";

interface LoginPageProps {
    version?: string;
    errors?: Record<string, string>;
    [key: string]: any;
}

export default function LoginPage() {
    const { props } = usePage<LoginPageProps>();

    return (
        <AuthLayout>
            <Head>
                <title>Login</title>
            </Head>
            <LoginForm errors={props.errors || {}} />
            {props.version && (
                <div>
                    {props.version}
                </div>
            )}
        </AuthLayout>
    );
}
