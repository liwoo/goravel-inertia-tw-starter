---
name: inertia-custom-page
description: Create a custom (non-CRUD) Inertia page controller and React component with i18n. Use for dashboards, analytics, reports, or any page that doesn't follow the standard CRUD pattern.
argument-hint: "[PageName]"
allowed-tools: Bash, Read, Write, Edit, Grep, Glob
---

# Custom Inertia Page (Non-CRUD, i18n-aware)

Create custom page for `$ARGUMENTS`.

## When to Use

Use this for pages that don't follow the standard CRUD pattern:
- Dashboards with aggregated data
- Analytics/reporting pages
- Settings pages
- Multi-service data aggregation
- Custom workflows

## Step 1: Go Page Controller

Create `app/http/controllers/<page_name>/<page_name>_controller.go`:

```go
package pagename

import (
    "github.com/goravel/framework/contracts/http"
    inertiaHelper "smedi-sme-db/app/http/inertia"
    "smedi-sme-db/app/services"
)

type PageNameController struct {
    serviceA *services.ServiceA
    serviceB *services.ServiceB
}

func NewPageNameController() *PageNameController {
    return &PageNameController{
        serviceA: services.NewServiceA(),
        serviceB: services.NewServiceB(),
    }
}

func (c *PageNameController) Show(ctx http.Context) http.Response {
    dataA, _ := c.serviceA.GetSummary()
    dataB, _ := c.serviceB.GetRecentItems()

    return inertiaHelper.Render(ctx, "PageName/Index", map[string]interface{}{
        "summary":     dataA,
        "recentItems": dataB,
        "version":     "1.0",
    })
}
```

## Step 2: Register Route

In `routes/web.go`:

```go
pageNameController := pagename.NewPageNameController()

// Inside authenticated group:
router.Get("/page-name", pageNameController.Show)
```

## Step 3: i18n Translation Namespace

Create `resources/js/locales/en/<page_name>.json`:

```json
{
  "page": {
    "title": "Page Name",
    "headTitle": "Page Name - Admin"
  },
  "sections": {
    "summary": "Summary",
    "recentItems": "Recent Items"
  },
  "stats": {
    "totalCount": "Total",
    "activeCount": "Active",
    "revenue": "Revenue"
  },
  "labels": {
    "notSpecified": "Not specified"
  }
}
```

Register in `resources/js/locales/index.ts`:

```typescript
import pageName from './en/page_name.json';

// Add to ns array and resources.en
ns: [..., 'pageName'],
resources: { en: { ..., pageName } },
```

## Step 4: TypeScript Props Interface

Create `resources/js/types/page_name.ts`:

```typescript
export interface Summary {
    totalCount: number;
    activeCount: number;
    revenue: number;
}

export interface RecentItem {
    id: number;
    title: string;
    createdAt: string;
}

export interface PageNameProps {
    summary: Summary;
    recentItems: RecentItem[];
}
```

## Step 5: React Page Component

Create `resources/js/pages/PageName/Index.tsx`:

```tsx
import React from 'react';
import { Head } from '@inertiajs/react';
import { useTranslation } from 'react-i18next';
import Admin from '@/layouts/Admin';
import { PageNameProps } from '@/types/page_name';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

export default function PageNameIndex({ summary, recentItems }: PageNameProps) {
    const { t } = useTranslation('pageName');

    return (
        <Admin title={t('page.title')}>
            <Head title={t('page.headTitle')} />

            <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
                {/* Stats Cards */}
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4 px-4">
                    <Card>
                        <CardHeader className="pb-2">
                            <CardTitle className="text-sm font-medium text-muted-foreground">
                                {t('stats.totalCount')}
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="text-2xl font-bold">{summary.totalCount}</div>
                        </CardContent>
                    </Card>
                    <Card>
                        <CardHeader className="pb-2">
                            <CardTitle className="text-sm font-medium text-muted-foreground">
                                {t('stats.activeCount')}
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="text-2xl font-bold">{summary.activeCount}</div>
                        </CardContent>
                    </Card>
                    <Card>
                        <CardHeader className="pb-2">
                            <CardTitle className="text-sm font-medium text-muted-foreground">
                                {t('stats.revenue')}
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="text-2xl font-bold">
                                {new Intl.NumberFormat('en-US', {
                                    style: 'currency',
                                    currency: 'USD',
                                }).format(summary.revenue)}
                            </div>
                        </CardContent>
                    </Card>
                </div>

                {/* Content Section */}
                <div className="px-4">
                    <Card>
                        <CardHeader>
                            <CardTitle>{t('sections.recentItems')}</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="space-y-3">
                                {recentItems.map((item) => (
                                    <div key={item.id} className="flex justify-between items-center p-3 rounded-lg bg-muted/50">
                                        <div>
                                            <p className="font-medium">{item.title}</p>
                                            <p className="text-sm text-muted-foreground">
                                                {new Date(item.createdAt).toLocaleDateString()}
                                            </p>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </CardContent>
                    </Card>
                </div>
            </div>
        </Admin>
    );
}
```

## i18n Pattern for Custom Pages

Custom pages use `useTranslation` directly (like forms, not like config functions):

```tsx
const { t } = useTranslation('pageName');

// Page titles
<Admin title={t('page.title')}>
<Head title={t('page.headTitle')} />

// Section headers
<CardTitle>{t('sections.recentItems')}</CardTitle>

// Stats labels
<CardTitle>{t('stats.totalCount')}</CardTitle>
```

## Loading External Data (Axios)

```tsx
import axios from 'axios';
import { useTranslation } from 'react-i18next';
import { useEffect, useState } from 'react';

const { t } = useTranslation('pageName');
const [data, setData] = useState<DataType[]>([]);

useEffect(() => {
    axios.get<DataType[]>('/api/endpoint')
        .then(({ data }) => setData(data))
        .catch(console.error);
}, []);
```

## Key Principles

1. Aggregate data on the **backend** — keep frontend display simple
2. Use proper TypeScript interfaces for all props
3. Wrap in `Admin` layout with i18n-translated `Head` title
4. Use shadcn/ui `Card` components for sections
5. Make responsive with grid `grid-cols-1 md:grid-cols-3`
6. **Create dedicated i18n namespace** for the page
7. All user-visible strings go through `t()` — no hardcoded text

## Reference

See `app/http/controllers/dashboard_controller.go` and `resources/js/pages/dashboard/Index.tsx` for the dashboard pattern.
