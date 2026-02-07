import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Alert, AlertDescription } from "@/components/ui/alert"
// @ts-ignore
import { useForm, usePage, Link, router } from '@inertiajs/react';
import React, { useState } from 'react';
import { toast } from "sonner";
import axios from "@/lib/axios";
import { useTranslation } from 'react-i18next';
import {
  Shield,
  ArrowLeft,
  Loader2,
  AlertCircle,
  Key,
} from "lucide-react";
import type { ApiErrorResponse } from "@/types/account";

type LoginStep = 'credentials' | '2fa';

interface TwoFactorState {
  tempToken: string;
  useBackupCode: boolean;
}

export function LoginForm({
  className,
  ...props
}: React.ComponentPropsWithoutRef<"form">) {
  const { t } = useTranslation('auth');
  const { data, setData, errors, reset } = useForm({
    email: '',
    password: '',
  });

  const page = usePage();

  // Login state
  const [isLoggingIn, setIsLoggingIn] = useState(false);

  // 2FA state
  const [loginStep, setLoginStep] = useState<LoginStep>('credentials');
  const [twoFactorState, setTwoFactorState] = useState<TwoFactorState>({
    tempToken: '',
    useBackupCode: false,
  });
  const [twoFactorCode, setTwoFactorCode] = useState('');
  const [twoFactorError, setTwoFactorError] = useState<string | null>(null);
  const [isVerifying2FA, setIsVerifying2FA] = useState(false);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsLoggingIn(true);

    try {
      // Use axios for login to properly handle 2FA JSON response
      const response = await axios.post('/login', {
        email: data.email,
        password: data.password,
      });

      // Check if 2FA is required (JSON response)
      if (response.data?.requires_2fa && response.data?.temp_token) {
        setTwoFactorState({
          tempToken: response.data.temp_token,
          useBackupCode: false,
        });
        setLoginStep('2fa');
        return;
      }

      // Successful login - use the redirect URL from the response
      const redirectUrl = response.data?.redirect || '/dashboard';
      router.visit(redirectUrl, { replace: true });
    } catch (error: unknown) {
      reset('password');

      if (axios.isAxiosError(error) && error.response) {
        // Handle redirect responses (302/303) - axios follows them automatically
        // but may throw on CORS or other issues
        if (error.response.status === 302 || error.response.status === 303) {
          // Successful login with redirect
          router.visit('/dashboard', { replace: true });
          return;
        }

        const errorData = error.response.data;

        // Check if 2FA required came back as error
        if (errorData?.requires_2fa && errorData?.temp_token) {
          setTwoFactorState({
            tempToken: errorData.temp_token,
            useBackupCode: false,
          });
          setLoginStep('2fa');
          return;
        }

        // Handle validation errors
        if (errorData?.errors) {
          if (errorData.errors.email) {
            toast.error(errorData.errors.email);
          }
          if (errorData.errors.password) {
            toast.error(errorData.errors.password);
          }
        } else if (errorData?.message) {
          toast.error(errorData.message);
        } else {
          toast.error(t('login.invalidCredentials'));
        }
      } else {
        toast.error(t('login.unexpectedError'));
      }
    } finally {
      setIsLoggingIn(false);
    }
  };

  const handle2FASubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setTwoFactorError(null);

    if (!twoFactorCode) {
      setTwoFactorError(t('twoFactor.enterCodeError'));
      return;
    }

    // Validate code format
    const isBackupCode = twoFactorState.useBackupCode || twoFactorCode.length === 8;
    if (!isBackupCode && twoFactorCode.length !== 6) {
      setTwoFactorError(t('twoFactor.invalidCodeFormat'));
      return;
    }

    setIsVerifying2FA(true);

    try {
      // Call the web route for 2FA verification
      const response = await axios.post('/verify-2fa', {
        temp_token: twoFactorState.tempToken,
        code: twoFactorCode,
      });

      // On success, the backend sets the session cookies
      toast.success(t('twoFactor.loginSuccess'));

      // Use Inertia router to navigate to the dashboard
      router.visit('/dashboard', { replace: true });
    } catch (error: unknown) {
      if (axios.isAxiosError(error) && error.response) {
        const errorData = error.response.data as ApiErrorResponse;

        if (error.response.status === 401) {
          setTwoFactorError(
            twoFactorState.useBackupCode
              ? t('twoFactor.invalidBackupCode')
              : t('twoFactor.invalidAuthCode')
          );
        } else if (error.response.status === 410) {
          // Token expired
          setTwoFactorError(t('twoFactor.sessionExpired'));
          setTimeout(() => {
            handleBack();
          }, 2000);
        } else if (errorData.errors?.code) {
          setTwoFactorError(errorData.errors.code[0]);
        } else {
          setTwoFactorError(errorData.message || t('twoFactor.verifyFailed'));
        }
      } else {
        setTwoFactorError(t('twoFactor.unexpectedError'));
      }
    } finally {
      setIsVerifying2FA(false);
    }
  };

  const handleBack = () => {
    setLoginStep('credentials');
    setTwoFactorState({ tempToken: '', useBackupCode: false });
    setTwoFactorCode('');
    setTwoFactorError(null);
    reset('password');
  };

  const toggleBackupCodeMode = () => {
    setTwoFactorState(prev => ({
      ...prev,
      useBackupCode: !prev.useBackupCode,
    }));
    setTwoFactorCode('');
    setTwoFactorError(null);
  };

  // 2FA verification step
  if (loginStep === '2fa') {
    return (
      <form onSubmit={handle2FASubmit} className={cn("flex flex-col gap-6 pb-8", className)} {...props}>
        <div className="flex flex-col items-center gap-4 text-center">
          <div className="p-3 rounded-full bg-purple-100 dark:bg-purple-900">
            <Shield className="h-8 w-8 text-purple-600 dark:text-purple-400" />
          </div>
          <div>
            <h3 className="text-xl font-semibold">{t('twoFactor.title')}</h3>
            <p className="text-sm text-muted-foreground mt-1">
              {twoFactorState.useBackupCode
                ? t('twoFactor.enterBackupCode')
                : t('twoFactor.enterCode')}
            </p>
          </div>
        </div>

        {twoFactorError && (
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>{twoFactorError}</AlertDescription>
          </Alert>
        )}

        <div className="grid gap-6">
          <div className="grid gap-2">
            <Label htmlFor="2fa-code">
              {twoFactorState.useBackupCode ? t('twoFactor.backupCodeLabel') : t('twoFactor.authCodeLabel')}
            </Label>
            <div className="relative">
              <Key className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                id="2fa-code"
                type="text"
                inputMode={twoFactorState.useBackupCode ? "text" : "numeric"}
                pattern={twoFactorState.useBackupCode ? "[A-Za-z0-9]*" : "[0-9]*"}
                maxLength={twoFactorState.useBackupCode ? 8 : 6}
                value={twoFactorCode}
                onChange={(e) => {
                  const value = twoFactorState.useBackupCode
                    ? e.target.value.toUpperCase()
                    : e.target.value.replace(/\D/g, '');
                  setTwoFactorCode(value);
                }}
                placeholder={twoFactorState.useBackupCode ? "XXXXXXXX" : "000000"}
                className="pl-10 text-center text-xl font-mono tracking-widest"
                autoComplete="one-time-code"
                autoFocus
                disabled={isVerifying2FA}
              />
            </div>
            <p className="text-xs text-muted-foreground text-center">
              {twoFactorState.useBackupCode
                ? t('twoFactor.backupCodeHint')
                : t('twoFactor.authCodeHint')}
            </p>
          </div>

          <Button type="submit" className="w-full" disabled={isVerifying2FA || !twoFactorCode}>
            {isVerifying2FA ? (
              <>
                <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                {t('twoFactor.verifying')}
              </>
            ) : (
              t('twoFactor.verify')
            )}
          </Button>

          <div className="relative">
            <div className="absolute inset-0 flex items-center">
              <span className="w-full border-t" />
            </div>
            <div className="relative flex justify-center text-xs uppercase">
              <span className="bg-background px-2 text-muted-foreground">{t('common:labels.or')}</span>
            </div>
          </div>

          <Button
            type="button"
            variant="outline"
            className="w-full"
            onClick={toggleBackupCodeMode}
            disabled={isVerifying2FA}
          >
            {twoFactorState.useBackupCode
              ? t('twoFactor.useAuthenticator')
              : t('twoFactor.useBackupCode')}
          </Button>

          <Button
            type="button"
            variant="ghost"
            className="w-full"
            onClick={handleBack}
            disabled={isVerifying2FA}
          >
            <ArrowLeft className="h-4 w-4 mr-2" />
            {t('twoFactor.backToLogin')}
          </Button>
        </div>
      </form>
    );
  }

  // Normal login form
  return (
    <form onSubmit={handleSubmit} className={cn("flex flex-col gap-6 pb-8", className)} {...props}>
      <div className="flex flex-col items-start gap-2 text-center">
        <div className="flex flex-col items-start w-full">
          <h3 className="text-xl font-semibold">{t('login.title')}</h3>
          <p className="text-balance text-sm text-muted-foreground">
            {t('login.subtitle')}
          </p>
        </div>
      </div>

      {/* Display general errors */}
      {page.props.errors?.general && (
        <div className="text-sm text-red-500 text-center bg-red-50 p-2 rounded">
          {page.props.errors.general}
        </div>
      )}

      <div className="grid gap-6">
        <div className="grid gap-2">
          <Label htmlFor="email">{t('login.email')}</Label>
          <Input
            id="email"
            type="email"
            placeholder={t('login.emailPlaceholder')}
            value={data.email}
            onChange={(e) => setData('email', e.target.value)}
            required
          />
          {(errors.email || page.props.errors?.email) && (
            <p className="text-xs text-red-500 mt-1">
              {errors.email || page.props.errors?.email}
            </p>
          )}
        </div>
        <div className="grid gap-2">
          <div className="flex items-center">
            <Label htmlFor="password">{t('login.password')}</Label>
          </div>
          <Input
            id="password"
            type="password"
            value={data.password}
            onChange={(e) => setData('password', e.target.value)}
            required
          />
          {(errors.password || page.props.errors?.password) && (
            <p className="text-xs text-red-500 mt-1">
              {errors.password || page.props.errors?.password}
            </p>
          )}
        </div>
        <Button type="submit" className="w-full" disabled={isLoggingIn}>
          {isLoggingIn ? t('login.submitting') : t('login.submit')}
        </Button>
        <div className="relative">
          <div className="absolute inset-0 flex items-center">
            <span className="w-full border-t" />
          </div>
          <div className="relative flex justify-center text-xs uppercase">
            <span className="bg-background px-2 text-muted-foreground">{t('common:labels.or')}</span>
          </div>
        </div>
        <Button variant="outline" className="w-full" asChild>
          <Link href="/una">{t('common:labels.register')}</Link>
        </Button>
      </div>
    </form>
  )
}
