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
          toast.error('Invalid credentials');
        }
      } else {
        toast.error('An unexpected error occurred');
      }
    } finally {
      setIsLoggingIn(false);
    }
  };

  const handle2FASubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setTwoFactorError(null);

    if (!twoFactorCode) {
      setTwoFactorError('Please enter your authentication code');
      return;
    }

    // Validate code format
    const isBackupCode = twoFactorState.useBackupCode || twoFactorCode.length === 8;
    if (!isBackupCode && twoFactorCode.length !== 6) {
      setTwoFactorError('Please enter a valid 6-digit code');
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
      toast.success('Login successful');

      // Use Inertia router to navigate to the dashboard
      router.visit('/dashboard', { replace: true });
    } catch (error: unknown) {
      if (axios.isAxiosError(error) && error.response) {
        const errorData = error.response.data as ApiErrorResponse;

        if (error.response.status === 401) {
          setTwoFactorError(
            twoFactorState.useBackupCode
              ? 'Invalid backup code. Please try again.'
              : 'Invalid authentication code. Please try again.'
          );
        } else if (error.response.status === 410) {
          // Token expired
          setTwoFactorError('Your login session has expired. Please start over.');
          setTimeout(() => {
            handleBack();
          }, 2000);
        } else if (errorData.errors?.code) {
          setTwoFactorError(errorData.errors.code[0]);
        } else {
          setTwoFactorError(errorData.message || 'Failed to verify code');
        }
      } else {
        setTwoFactorError('An unexpected error occurred. Please try again.');
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
            <h3 className="text-xl font-semibold">Two-Factor Authentication</h3>
            <p className="text-sm text-muted-foreground mt-1">
              {twoFactorState.useBackupCode
                ? 'Enter one of your backup codes'
                : 'Enter the code from your authenticator app'}
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
              {twoFactorState.useBackupCode ? 'Backup Code' : 'Authentication Code'}
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
                ? 'Enter an 8-character backup code'
                : 'Enter the 6-digit code from your authenticator app'}
            </p>
          </div>

          <Button type="submit" className="w-full" disabled={isVerifying2FA || !twoFactorCode}>
            {isVerifying2FA ? (
              <>
                <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                Verifying...
              </>
            ) : (
              'Verify'
            )}
          </Button>

          <div className="relative">
            <div className="absolute inset-0 flex items-center">
              <span className="w-full border-t" />
            </div>
            <div className="relative flex justify-center text-xs uppercase">
              <span className="bg-background px-2 text-muted-foreground">Or</span>
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
              ? 'Use authenticator app'
              : 'Use a backup code'}
          </Button>

          <Button
            type="button"
            variant="ghost"
            className="w-full"
            onClick={handleBack}
            disabled={isVerifying2FA}
          >
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back to login
          </Button>
        </div>
      </form>
    );
  }

  // Normal login form
  return (
    <form onSubmit={handleSubmit} className={cn("flex flex-col gap-6 pb-8", className)} {...props}>
      <div className="flex flex-col items-start gap-2 text-center">
        <div className="flex items-center gap-3 w-full">
          <img src="/images/mw-coat.svg" alt="MW Gov Emblem" className="w-1/5" />
          <div className="flex flex-col items-start">
            <h3 className="text-xl font-semibold uppercase text-nowrap">National MSME's Database</h3>
            <h4>Management Information System</h4>
          </div>
        </div>
        <p className="text-balance text-sm text-muted-foreground">
          Enter your email below to login to your account
        </p>
      </div>

      {/* Display general errors */}
      {page.props.errors?.general && (
        <div className="text-sm text-red-500 text-center bg-red-50 p-2 rounded">
          {page.props.errors.general}
        </div>
      )}

      <div className="grid gap-6">
        <div className="grid gap-2">
          <Label htmlFor="email">Email</Label>
          <Input
            id="email"
            type="email"
            placeholder="m@example.com"
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
            <Label htmlFor="password">Password</Label>
            {/* <a
              href="#"
              className="ml-auto text-sm underline-offset-4 hover:underline"
            >
              Forgot your password?
            </a> */}
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
          {isLoggingIn ? 'Logging in...' : 'Login'}
        </Button>
        <div className="relative">
          <div className="absolute inset-0 flex items-center">
            <span className="w-full border-t" />
          </div>
          <div className="relative flex justify-center text-xs uppercase">
            <span className="bg-background px-2 text-muted-foreground">Or</span>
          </div>
        </div>
        <Button variant="outline" className="w-full" asChild>
          <Link href="/apply">Apply for Access</Link>
        </Button>
      </div>
    </form>
  )
}
