import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
// @ts-ignore
import { useForm, usePage } from "@inertiajs/react";
import React, { useState, useMemo } from "react";
import { toast } from "sonner";
import { Check, X, Eye, EyeOff } from "lucide-react";

interface PasswordRequirement {
  label: string;
  test: (password: string) => boolean;
}

export function LoginForm({
  className,
  ...props
}: React.ComponentPropsWithoutRef<"form">) {
  const { data, setData, post, processing, errors, reset } = useForm({
    email: "",
    password: "",
  });

  const [showPasswordRequirements, setShowPasswordRequirements] =
    useState(false);
  const [showPassword, setShowPassword] = useState(false);

  const page = usePage();

  // Password validation requirements
  const passwordRequirements: PasswordRequirement[] = [
    {
      label: "Minimum of 8 characters",
      test: (password) => password.length >= 8,
    },
    {
      label: "1 Capital letter",
      test: (password) => /[A-Z]/.test(password),
    },
    {
      label: "One symbol",
      test: (password) =>
        /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(password),
    },
    {
      label: "One number",
      test: (password) => /[0-9]/.test(password),
    },
    {
      label: "No spaces",
      test: (password) => !/\s/.test(password),
    },
  ];

  // Check which requirements are met
  const passwordValidation = useMemo(() => {
    return passwordRequirements.map((req) => ({
      ...req,
      met: req.test(data.password),
    }));
  }, [data.password]);

  // Check if all requirements are met
  const isPasswordValid = useMemo(() => {
    return passwordValidation.every((req) => req.met);
  }, [passwordValidation]);

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    post("/login", {
      onFinish: () => reset("password"), // Optional: reset password field on finish
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
    <form
      onSubmit={handleSubmit}
      className={cn("flex flex-col gap-6 pb-8", className)}
      {...props}
    >
      <div className="flex flex-col items-start gap-2 text-center">
        <div className="flex items-center gap-3 w-full">
          <img
            src="/images/mw-coat.svg"
            alt="MW Gov Emblem"
            className="w-1/5"
          />
          <div className="flex flex-col items-start">
            <h3 className="text-xl font-semibold uppercase text-nowrap">
              National MSME's Database
            </h3>
            <h4>Management Information System</h4>
          </div>
        </div>
        {/*<p className="text-balance text-sm text-muted-foreground">
          Enter your email below to login to your account
        </p>*/}
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
            onChange={(e) => setData("email", e.target.value)}
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
          <div className="relative">
            <Input
              id="password"
              type={showPassword ? "text" : "password"}
              value={data.password}
              onChange={(e) => setData("password", e.target.value)}
              onFocus={() => setShowPasswordRequirements(true)}
              onBlur={() => setShowPasswordRequirements(false)}
              className={cn(
                "pr-10",
                data.password &&
                  !isPasswordValid &&
                  "border-red-300 focus:border-red-500 outline-red-500",
                data.password &&
                  isPasswordValid &&
                  "border-green-300 focus:border-green-500 outline-green-500",
              )}
              required
            />
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none"
            >
              {showPassword ? (
                <EyeOff className="w-4 h-4" />
              ) : (
                <Eye className="w-4 h-4" />
              )}
            </button>
          </div>

          {/* Password Requirements */}
          {showPasswordRequirements && data.password && (
            <div className="mt-2 p-3 bg-backgroun/50 rounded-md border">
              <p className="text-xs font-medium text-gray-700 mb-2">
                Password Requirements:
              </p>
              <div className="space-y-1">
                {passwordValidation.map((req, index) => (
                  <div key={index} className="flex items-center gap-2 text-xs">
                    {req.met ? (
                      <Check className="w-3 h-3 text-green-500" />
                    ) : (
                      <X className="w-3 h-3 text-red-500" />
                    )}
                    <span
                      className={cn(
                        req.met ? "text-green-700" : "text-red-700",
                      )}
                    >
                      {req.label}
                    </span>
                  </div>
                ))}
              </div>
            </div>
          )}

          {(errors.password || page.props.errors?.password) && (
            <p className="text-xs text-red-500 mt-1">
              {errors.password || page.props.errors?.password}
            </p>
          )}
        </div>
        <Button 
          type="submit" 
          className="w-full" 
          disabled={processing || !isPasswordValid || !data.email || !data.password}
        >
          {processing ? "Logging in..." : "Login"}
        </Button>
      </div>
    </form>
  );
}
