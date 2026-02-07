import React from 'react';
import { TFunction } from 'i18next';
import { CheckCircle, XCircle } from 'lucide-react';
import { STATUS_COLORS } from '@/config/status-colors';

export function getAuthorStatusConfig(t: TFunction) {
  return {
    ACTIVE: {
      label: t('status.active'),
      icon: <CheckCircle className="h-3 w-3" />,
      color: STATUS_COLORS.positive,
    },
    INACTIVE: {
      label: t('status.inactive'),
      icon: <XCircle className="h-3 w-3" />,
      color: STATUS_COLORS.negative,
    },
  };
}
