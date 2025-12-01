import React from 'react';
import { Download, Upload, BarChart3 } from 'lucide-react';
import {
  StatsCardConfig,
  PageActionConfig,
  SimpleFilterConfig
} from '@/lib/crud-page-utils';
import { ConfigType } from '@/types/config';

/**
 * Stats card configurations for configs
 */
export const configStatsConfigs: StatsCardConfig[] = [
  // TODO: Configure your stats cards here
];

/**
 * Simple filter configurations for configs
 */
export const configSimpleFilters = (stats: any): SimpleFilterConfig[] => [
  {
    key: 'type-financing',
    label: 'Finance',
    value: ConfigType.Financing,
    badge: stats?.financingCount || 0,
    filterParams: { config_type: ConfigType.Financing }
  },
  {
    key: 'type-improvement',
    label: 'Improvements',
    value: ConfigType.ImprovementAspects,
    badge: stats?.improvementAspectsCount || 0,
    filterParams: { config_type: ConfigType.ImprovementAspects }
  },
  {
    key: 'type-business-categories',
    label: 'Categories',
    value: ConfigType.BusinessCategories,
    badge: stats?.businessCategoriesCount || 0,
    filterParams: { config_type: ConfigType.BusinessCategories }
  },
  {
    key: 'type-industries',
    label: 'Industries',
    value: ConfigType.Industries,
    badge: stats?.industriesCount || 0,
    filterParams: { config_type: ConfigType.Industries }
  },
  {
    key: 'type-sectors',
    label: 'Sectors',
    value: ConfigType.Sectors,
    badge: stats?.sectorsCount || 0,
    filterParams: { config_type: ConfigType.Sectors }
  },
  {
    key: 'type-registration',
    label: 'Status',
    value: ConfigType.RegistrationStatus,
    badge: stats?.registrationStatusCount || 0,
    filterParams: { config_type: ConfigType.RegistrationStatus }
  },
  {
    key: 'type-partners',
    label: 'Partners',
    value: ConfigType.DevelopmentPartners,
    badge: stats?.developmentPartnersCount || 0,
    filterParams: { config_type: ConfigType.DevelopmentPartners }
  },
  {
    key: 'type-product-types',
    label: 'Products',
    value: ConfigType.ProductTypes,
    badge: stats?.productTypesCount || 0,
    filterParams: { config_type: ConfigType.ProductTypes }
  },
  {
    key: 'type-procurement-type',
    label: 'Proc. Type',
    value: ConfigType.ProcurementType,
    badge: stats?.procurementTypeCount || 0,
    filterParams: { config_type: ConfigType.ProcurementType }
  },
  {
    key: 'type-procurement-classification',
    label: 'Proc. Class.',
    value: ConfigType.ProcurementClassification,
    badge: stats?.procurementClassificationCount || 0,
    filterParams: { config_type: ConfigType.ProcurementClassification }
  },
  {
    key: 'type-procured-by',
    label: 'Procurers',
    value: ConfigType.ProcuredBy,
    badge: stats?.procuredByCount || 0,
    filterParams: { config_type: ConfigType.ProcuredBy }
  },
  {
    key: 'type-organization',
    label: 'Orgs',
    value: ConfigType.Organization,
    badge: stats?.organizationCount || 0,
    filterParams: { config_type: ConfigType.Organization }
  },
];

/**
 * Page action configurations for configs
 */
export const getConfigPageActions = (
  permissions: any,
  handlers: {
    onImport?: () => void;
    onExport?: () => void;
  }
): PageActionConfig[] => [
  // TODO: Configure your page actions here
];

/**
 * Bulk action configurations for configs
 */
export const configBulkActions = {
  handleBulkDelete: (configIds: number[]) => {
    const confirmMessage = `Are you sure you want to delete ${configIds.length} item(s)?`;
    if (confirm(confirmMessage)) {
      // TODO: Implement bulk delete
    }
  },
};

import { router } from '@inertiajs/react';
