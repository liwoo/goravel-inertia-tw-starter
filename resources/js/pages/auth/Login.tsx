import AuthLayout from '@/layouts/Auth';
import { LoginForm } from "@/components/login-form";
// @ts-ignore
import { Head } from "@inertiajs/react"; // Assuming this component exists

interface HomeProps {
    version?: string;
    // Add other props your component might receive
}

export default function LoginPage({ version }: HomeProps) {
    return (
        <AuthLayout>
            <Head>
                <title>Login</title>
            </Head>
            <LoginForm />
            <div>
                <h5 className="bg-background/50 p-2 rounded w-auto">
                    {version}
                </h5>
                <h6 className="text-xs text-muted-foreground my-2">Property of Ministry of Trade and Industry and SMEDI. All rights reserved.</h6>
            </div>
        </AuthLayout>
    );
}
