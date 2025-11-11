import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
// @ts-ignore
import { Link, useForm, usePage } from '@inertiajs/react';
import React from 'react';
import { toast } from "sonner";

export function LoginForm({
  className,
  ...props
}: React.ComponentPropsWithoutRef<"form">) {
  const { data, setData, post, processing, errors, reset } = useForm({
    email: '',
    password: '',
  });

  const page = usePage();

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    post('/login', {
      onFinish: () => reset('password'), // Optional: reset password field on finish
      onError: (errors: any) => {
        if (errors.email) {
          toast.error(errors.email);
        }
        if (errors.password) {
          toast.error(errors.password);
        }
        if (errors.general) {
          toast.error(errors.general);
        }
      },
    });
  };

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
            <Link
              href="/forgot-password"
              className="ml-auto text-sm underline-offset-4 hover:underline"
            >
              Forgot your password?
            </Link>
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
        <Button type="submit" className="w-full" disabled={processing}>
          {processing ? 'Logging in...' : 'Login'}
        </Button>
      </div>
    </form>
  )
}