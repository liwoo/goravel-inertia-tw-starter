import AuthLayout from '@/layouts/Auth';
import { LoginForm } from "@/components/login-form";
// @ts-ignore
import { Head, usePage } from "@inertiajs/react";
import type { SharedData } from "@/types/app";
import { useTranslation } from 'react-i18next';

export default function LoginPage() {
    const { appVersion } = usePage<SharedData>().props;
    const { t } = useTranslation('auth');

    return (
        <AuthLayout>
            <Head>
                <title>{t('login.submit')}</title>
            </Head>
            <LoginForm />
            <div className="text-center">
                <p className="text-xs text-muted-foreground">
                    v{appVersion || '0.0.0'}
                </p>
                <p className="text-xs text-muted-foreground mt-1">
                    {t('branding.copyright', { year: new Date().getFullYear() })}
                </p>
            </div>
        </AuthLayout>
    );
}
