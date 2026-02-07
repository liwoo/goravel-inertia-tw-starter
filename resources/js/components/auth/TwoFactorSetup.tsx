/**
 * TwoFactorSetup Component
 *
 * Multi-step dialog for enabling TOTP two-factor authentication.
 *
 * Steps:
 * 1. Password confirmation
 * 2. QR code display with manual secret
 * 3. TOTP code verification
 * 4. Backup codes display
 */

import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent } from '@/components/ui/card';
import { Alert, AlertDescription } from '@/components/ui/alert';
import {
  Shield,
  Key,
  Copy,
  Download,
  CheckCircle,
  Loader2,
  Eye,
  EyeOff,
  AlertCircle,
  Smartphone,
  QrCode,
} from 'lucide-react';
import { toast } from 'sonner';
import { setupTOTP, verifyTOTP } from '@/lib/api/two-factor';
import type {
  TOTPSetupData,
  TOTPFormErrors,
  ApiErrorResponse,
} from '@/types/account';
import axios from '@/lib/axios';

type SetupStep = 'password' | 'qr-code' | 'verify' | 'backup-codes';

interface TwoFactorSetupProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: () => void;
}

export const TwoFactorSetup: React.FC<TwoFactorSetupProps> = ({
  open,
  onOpenChange,
  onSuccess,
}) => {
  const { t } = useTranslation('settings');
  const [step, setStep] = useState<SetupStep>('password');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [verificationCode, setVerificationCode] = useState('');
  const [setupData, setSetupData] = useState<TOTPSetupData | null>(null);
  const [backupCodes, setBackupCodes] = useState<string[]>([]);
  const [errors, setErrors] = useState<TOTPFormErrors>({});
  const [isLoading, setIsLoading] = useState(false);
  const [copiedSecret, setCopiedSecret] = useState(false);
  const [copiedCodes, setCopiedCodes] = useState(false);

  const resetState = () => {
    setStep('password');
    setPassword('');
    setShowPassword(false);
    setVerificationCode('');
    setSetupData(null);
    setBackupCodes([]);
    setErrors({});
    setIsLoading(false);
    setCopiedSecret(false);
    setCopiedCodes(false);
  };

  const handleClose = () => {
    if (step === 'backup-codes') {
      // Don't allow closing until they acknowledge backup codes
      return;
    }
    resetState();
    onOpenChange(false);
  };

  const handlePasswordSubmit = async () => {
    if (!password) {
      setErrors({ password: t('validation.passwordRequired') });
      return;
    }

    setIsLoading(true);
    setErrors({});

    try {
      const response = await setupTOTP(password);
      // Backend wraps response in { success, message, data } structure
      setSetupData(response.data.data);
      setStep('qr-code');
    } catch (error: unknown) {
      if (axios.isAxiosError(error) && error.response) {
        const errorData = error.response.data as ApiErrorResponse;
        if (error.response.status === 401 || error.response.status === 403) {
          setErrors({ password: t('validation.incorrectPassword') });
        } else if (errorData.errors?.password) {
          setErrors({ password: errorData.errors.password[0] });
        } else {
          setErrors({ general: errorData.message || t('twoFactor.failedSetup') });
        }
      } else {
        setErrors({ general: t('errors.unexpectedError') });
      }
    } finally {
      setIsLoading(false);
    }
  };

  const handleVerifyCode = async () => {
    if (!verificationCode || verificationCode.length !== 6) {
      setErrors({ code: t('twoFactor.enterSixDigitCode') });
      return;
    }

    setIsLoading(true);
    setErrors({});

    try {
      const response = await verifyTOTP(verificationCode);
      // Backend wraps response in { success, message, data } structure
      setBackupCodes(response.data.data.backup_codes);
      setStep('backup-codes');
      toast.success(t('twoFactor.enabledSuccess'));
    } catch (error: unknown) {
      if (axios.isAxiosError(error) && error.response) {
        const errorData = error.response.data as ApiErrorResponse;
        if (errorData.errors?.code) {
          setErrors({ code: errorData.errors.code[0] });
        } else {
          setErrors({ code: errorData.message || t('twoFactor.invalidCode') });
        }
      } else {
        setErrors({ code: t('errors.unexpectedError') });
      }
    } finally {
      setIsLoading(false);
    }
  };

  const copyToClipboard = async (text: string, type: 'secret' | 'codes') => {
    try {
      await navigator.clipboard.writeText(text);
      if (type === 'secret') {
        setCopiedSecret(true);
        setTimeout(() => setCopiedSecret(false), 2000);
      } else {
        setCopiedCodes(true);
        setTimeout(() => setCopiedCodes(false), 2000);
      }
      toast.success(type === 'secret' ? t('backupCodes.secretCopied') : t('backupCodes.codesCopied'));
    } catch {
      toast.error(t('backupCodes.copyFailed'));
    }
  };

  const downloadBackupCodes = () => {
    const content = [
      t('backupCodes.downloadHeader'),
      '========================================================',
      '',
      t('backupCodes.downloadInstruction'),
      '',
      ...backupCodes.map((code, i) => `${i + 1}. ${code}`),
      '',
      `Generated: ${new Date().toISOString()}`,
    ].join('\n');

    const blob = new Blob([content], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = '2fa-backup-codes.txt';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
    toast.success(t('backupCodes.downloaded'));
  };

  const handleComplete = () => {
    resetState();
    onOpenChange(false);
    onSuccess();
  };

  const renderPasswordStep = () => (
    <>
      <DialogHeader>
        <div className="flex items-center gap-3">
          <div className="p-2 rounded-lg bg-purple-100 dark:bg-purple-900">
            <Shield className="h-5 w-5 text-purple-600 dark:text-purple-400" />
          </div>
          <div>
            <DialogTitle>{t('twoFactor.enableTitle')}</DialogTitle>
            <DialogDescription>
              {t('twoFactor.enableDescription')}
            </DialogDescription>
          </div>
        </div>
      </DialogHeader>

      <div className="space-y-4 py-4">
        <p className="text-sm text-muted-foreground">
          {t('twoFactor.securityInfo')}
        </p>

        {errors.general && (
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>{errors.general}</AlertDescription>
          </Alert>
        )}

        <div className="space-y-2">
          <Label htmlFor="setup-password">{t('security.currentPassword')}</Label>
          <div className="relative">
            <Input
              id="setup-password"
              type={showPassword ? 'text' : 'password'}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder={t('backupCodes.enterPassword')}
              className={`pr-10 ${errors.password ? 'border-destructive' : ''}`}
              disabled={isLoading}
              onKeyDown={(e) => e.key === 'Enter' && handlePasswordSubmit()}
            />
            <Button
              type="button"
              variant="ghost"
              size="icon"
              className="absolute right-1 top-1/2 -translate-y-1/2 h-8 w-8"
              onClick={() => setShowPassword(!showPassword)}
              tabIndex={-1}
            >
              {showPassword ? (
                <EyeOff className="h-4 w-4 text-muted-foreground" />
              ) : (
                <Eye className="h-4 w-4 text-muted-foreground" />
              )}
            </Button>
          </div>
          {errors.password && (
            <p className="text-sm text-destructive flex items-center gap-1">
              <AlertCircle className="h-3 w-3" />
              {errors.password}
            </p>
          )}
        </div>
      </div>

      <DialogFooter>
        <Button variant="outline" onClick={handleClose} disabled={isLoading}>
          {t('common:actions.cancel')}
        </Button>
        <Button onClick={handlePasswordSubmit} disabled={isLoading || !password}>
          {isLoading ? (
            <>
              <Loader2 className="h-4 w-4 mr-2 animate-spin" />
              {t('twoFactor.verifying')}
            </>
          ) : (
            t('twoFactor.continue')
          )}
        </Button>
      </DialogFooter>
    </>
  );

  const renderQRCodeStep = () => (
    <>
      <DialogHeader>
        <div className="flex items-center gap-3">
          <div className="p-2 rounded-lg bg-purple-100 dark:bg-purple-900">
            <QrCode className="h-5 w-5 text-purple-600 dark:text-purple-400" />
          </div>
          <div>
            <DialogTitle>{t('twoFactor.scanQRCode')}</DialogTitle>
            <DialogDescription>
              {t('twoFactor.scanQRCodeDescription')}
            </DialogDescription>
          </div>
        </div>
      </DialogHeader>

      <div className="space-y-4 py-4">
        <div className="flex items-start gap-2 text-sm text-muted-foreground">
          <Smartphone className="h-4 w-4 mt-0.5 flex-shrink-0" />
          <p>
            {t('twoFactor.scanInstruction')}
          </p>
        </div>

        {setupData && (
          <>
            <Card className="bg-white">
              <CardContent className="flex justify-center p-6">
                <img
                  src={setupData.qr_code}
                  alt={t('twoFactor.qrCodeAlt')}
                  className="w-48 h-48"
                />
              </CardContent>
            </Card>

            <div className="space-y-2">
              <Label className="text-xs text-muted-foreground">
                {t('twoFactor.cantScan')}
              </Label>
              <div className="flex items-center gap-2">
                <code className="flex-1 p-3 bg-muted rounded-md font-mono text-sm break-all">
                  {setupData.secret}
                </code>
                <Button
                  variant="outline"
                  size="icon"
                  onClick={() => copyToClipboard(setupData.secret, 'secret')}
                  className="flex-shrink-0"
                >
                  {copiedSecret ? (
                    <CheckCircle className="h-4 w-4 text-green-600" />
                  ) : (
                    <Copy className="h-4 w-4" />
                  )}
                </Button>
              </div>
            </div>
          </>
        )}
      </div>

      <DialogFooter>
        <Button variant="outline" onClick={handleClose}>
          {t('common:actions.cancel')}
        </Button>
        <Button onClick={() => setStep('verify')}>
          {t('twoFactor.scannedCode')}
        </Button>
      </DialogFooter>
    </>
  );

  const renderVerifyStep = () => (
    <>
      <DialogHeader>
        <div className="flex items-center gap-3">
          <div className="p-2 rounded-lg bg-purple-100 dark:bg-purple-900">
            <Key className="h-5 w-5 text-purple-600 dark:text-purple-400" />
          </div>
          <div>
            <DialogTitle>{t('twoFactor.verifyCode')}</DialogTitle>
            <DialogDescription>
              {t('twoFactor.verifyCodeDescription')}
            </DialogDescription>
          </div>
        </div>
      </DialogHeader>

      <div className="space-y-4 py-4">
        <p className="text-sm text-muted-foreground">
          {t('twoFactor.verifyInstruction')}
        </p>

        <div className="space-y-2">
          <Label htmlFor="verification-code">{t('twoFactor.verificationCode')}</Label>
          <Input
            id="verification-code"
            type="text"
            inputMode="numeric"
            pattern="[0-9]*"
            maxLength={6}
            value={verificationCode}
            onChange={(e) => {
              const value = e.target.value.replace(/\D/g, '');
              setVerificationCode(value);
            }}
            placeholder="000000"
            className={`text-center text-2xl font-mono tracking-widest ${
              errors.code ? 'border-destructive' : ''
            }`}
            disabled={isLoading}
            autoComplete="one-time-code"
            onKeyDown={(e) => e.key === 'Enter' && handleVerifyCode()}
          />
          {errors.code && (
            <p className="text-sm text-destructive flex items-center gap-1">
              <AlertCircle className="h-3 w-3" />
              {errors.code}
            </p>
          )}
        </div>
      </div>

      <DialogFooter>
        <Button variant="outline" onClick={() => setStep('qr-code')} disabled={isLoading}>
          {t('twoFactor.back')}
        </Button>
        <Button
          onClick={handleVerifyCode}
          disabled={isLoading || verificationCode.length !== 6}
        >
          {isLoading ? (
            <>
              <Loader2 className="h-4 w-4 mr-2 animate-spin" />
              {t('twoFactor.verifying')}
            </>
          ) : (
            t('twoFactor.verifyAndEnable')
          )}
        </Button>
      </DialogFooter>
    </>
  );

  const renderBackupCodesStep = () => (
    <>
      <DialogHeader>
        <div className="flex items-center gap-3">
          <div className="p-2 rounded-lg bg-green-100 dark:bg-green-900">
            <CheckCircle className="h-5 w-5 text-green-600 dark:text-green-400" />
          </div>
          <div>
            <DialogTitle>{t('backupCodes.title')}</DialogTitle>
            <DialogDescription>
              {t('backupCodes.description')}
            </DialogDescription>
          </div>
        </div>
      </DialogHeader>

      <div className="space-y-4 py-4">
        <Alert className="border-amber-500 bg-amber-50 dark:bg-amber-950">
          <AlertCircle className="h-4 w-4 text-amber-600" />
          <AlertDescription className="text-amber-700 dark:text-amber-300">
            {t('backupCodes.importantMessage')}
          </AlertDescription>
        </Alert>

        <Card>
          <CardContent className="p-4">
            <div className="grid grid-cols-2 gap-2">
              {backupCodes.map((code, index) => (
                <code
                  key={index}
                  className="p-2 bg-muted rounded text-center font-mono text-sm"
                >
                  {code}
                </code>
              ))}
            </div>
          </CardContent>
        </Card>

        <div className="flex gap-2">
          <Button
            variant="outline"
            className="flex-1"
            onClick={() => copyToClipboard(backupCodes.join('\n'), 'codes')}
          >
            {copiedCodes ? (
              <>
                <CheckCircle className="h-4 w-4 mr-2 text-green-600" />
                {t('backupCodes.copied')}
              </>
            ) : (
              <>
                <Copy className="h-4 w-4 mr-2" />
                {t('backupCodes.copyAll')}
              </>
            )}
          </Button>
          <Button variant="outline" className="flex-1" onClick={downloadBackupCodes}>
            <Download className="h-4 w-4 mr-2" />
            {t('backupCodes.download')}
          </Button>
        </div>
      </div>

      <DialogFooter>
        <Button onClick={handleComplete} className="w-full">
          {t('backupCodes.savedCodes')}
        </Button>
      </DialogFooter>
    </>
  );

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent
        className="sm:max-w-md"
        showCloseButton={step !== 'backup-codes'}
      >
        {step === 'password' && renderPasswordStep()}
        {step === 'qr-code' && renderQRCodeStep()}
        {step === 'verify' && renderVerifyStep()}
        {step === 'backup-codes' && renderBackupCodesStep()}
      </DialogContent>
    </Dialog>
  );
};

export default TwoFactorSetup;
