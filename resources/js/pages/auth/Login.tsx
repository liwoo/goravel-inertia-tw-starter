import AuthLayout from '@/layouts/Auth';
import {LoginForm} from "@/components/login-form";
// @ts-ignore
import {Head} from "@inertiajs/react"; // Assuming this component exists

interface HomeProps {
    version?: string;
    // Add other props your component might receive
}

export default function LoginPage({version}: HomeProps) {
    return (
        <AuthLayout>
            <Head>
                <title>Login</title>
            </Head>
            <LoginForm/>
            <div>
                {version}
            </div>
        </AuthLayout>
    );
}
