---
name: file-downloads
description: Add file download and export functionality (CSV, Excel, PDF, JSON) to a Goravel entity page using the existing client-side export system. Also covers chart PNG capture.
argument-hint: "[EntityName]"
allowed-tools: Read, Write, Edit, Grep, Glob
---

# File Downloads & Exports

**Agent**: Thoko Nkhoma — Senior Frontend Engineer

Add export functionality to `$ARGUMENTS`.

## Prerequisites

The export infrastructure already exists:
- `resources/js/utils/exportUtils.ts` — Format-specific export logic
- `resources/js/types/export.ts` — TypeScript types
- `resources/js/hooks/useExport.ts` — React hook for simpler integration
- `resources/js/components/ExportDialog.tsx` — Field picker + format selector dialog
- `resources/js/components/ui/chart-actions.tsx` — Chart CSV/PNG/fullscreen

NPM dependencies (already installed): `xlsx`, `jspdf`, `jspdf-autotable`, `html2canvas`

## Architecture

All exports are **client-side** — no backend endpoints needed:

```
Export button → ExportDialog (format, fields, stats)
  → exportData() in exportUtils.ts
    ├─ CSV  → escapeCSV + Blob
    ├─ JSON → JSON.stringify + Blob
    ├─ Excel → dynamic import('xlsx') → workbook
    └─ PDF  → dynamic import('jspdf') + autoTable
```

## Step 1: Define Export Fields

In `resources/js/pages/<Entity>/sections/<Entity>PageConfig.tsx` (or in `Index.tsx`):

```typescript
import { ExportField } from '@/types/export';

export function getEntityExportFields(t: TFunction): ExportField[] {
  return [
    { id: 'id', label: t('columns.id') },
    { id: 'name', label: t('columns.name') },
    { id: 'email', label: t('columns.email') },
    { id: 'status', label: t('columns.status') },
    { id: 'createdAt', label: t('columns.added'),
      formatter: (v) => v ? new Date(v).toLocaleDateString() : '' },
    // Array fields need a formatter
    { id: 'tags', label: t('columns.tags'),
      formatter: (v) => Array.isArray(v) ? v.join(', ') : (v || '') },
    // Currency fields need formatting
    { id: 'price', label: t('columns.price'),
      formatter: (v) => v != null ? `$${Number(v).toFixed(2)}` : '' },
  ];
}

export const DEFAULT_EXPORT_FIELDS = ['name', 'email', 'status', 'createdAt'];
```

### Field Formatter Rules

| Field Type | Formatter |
|-----------|-----------|
| Plain text | Not needed |
| Date | `(v) => v ? new Date(v).toLocaleDateString() : ''` |
| Currency | `(v) => v != null ? \`$\${Number(v).toFixed(2)}\` : ''` |
| Boolean | `(v) => v ? 'Yes' : 'No'` |
| Array | `(v) => Array.isArray(v) ? v.join(', ') : (v \|\| '')` |
| Enum/Status | `(v) => t(\`status.\${v?.toLowerCase()}\`) \|\| v` |
| FK relation | `(v, row) => row.authorName \|\| ''` (use display field) |

## Step 2: Add Export Handler to Index.tsx

```typescript
import { ExportOptions, ExportColumn } from '@/types/export';
import { exportData } from '@/utils/exportUtils';
import { ExportDialog } from '@/components/ExportDialog';

// In the component:
const [showExportDialog, setShowExportDialog] = useState(false);
const exportFields = getEntityExportFields(t);

const handleExport = async (options: ExportOptions) => {
  const columns: ExportColumn[] = exportFields
    .filter(f => options.fields.includes(f.id))
    .map(f => ({ id: f.id, label: f.label, formatter: f.formatter, width: 20 }));

  const statistics = options.includeStats && stats ? {
    [t('stats.total')]: stats.total,
    [t('stats.active')]: stats.active,
    [t('stats.inactive')]: stats.inactive,
  } : undefined;

  await exportData(
    { rows: data, columns, statistics, title: t('page.title') },
    options
  );
};
```

## Step 3: Add Export Button to Page Actions

In `<Entity>PageConfig.tsx`, add an export action:

```typescript
export function getEntityPageActions(t: TFunction, permissions: any, handlers: any) {
  return [
    // ... existing actions (Add Entity, etc.)
    {
      label: t('actions.export'),
      icon: Download,
      onClick: handlers.onExport,
      variant: 'outline' as const,
    },
  ];
}
```

Or add directly in `Index.tsx` toolbar.

## Step 4: Add ExportDialog to JSX

```tsx
<ExportDialog
  open={showExportDialog}
  onClose={() => setShowExportDialog(false)}
  onExport={handleExport}
  totalItems={stats?.total || data.length}
  availableFields={exportFields}
  defaultFields={DEFAULT_EXPORT_FIELDS}
  title={t('actions.exportTitle')}
  defaultFilename={`${entitySlug}-export`}
  preparedBy={currentUser?.name}
  showStatsOption={!!stats}
/>
```

## Step 5: Add i18n Keys

Add to the entity's translation namespace:

```json
{
  "actions": {
    "export": "Export",
    "exportTitle": "Export Entities"
  }
}
```

## Alternative: useExport Hook

For simpler integration without manual state:

```typescript
import { useExport } from '@/hooks/useExport';

const {
  isOpen, openDialog, closeDialog, handleExport, isExporting
} = useExport({
  getData: () => data,
  fields: exportFields,
  filename: 'entities-export',
  title: t('page.title'),
  getStatistics: () => stats ? {
    [t('stats.total')]: stats.total,
  } : undefined,
  onSuccess: () => toast.success(t('toast.exported')),
});
```

## Chart Exports (Dashboard Pages)

For dashboard charts, attach `ChartActions` to each card:

```tsx
import { ChartActions } from '@/components/ui/chart-actions';

const chartRef = React.useRef<HTMLDivElement>(null);

<Card>
  <CardHeader className="flex flex-row items-center justify-between">
    <CardTitle>{t('charts.revenue')}</CardTitle>
    <ChartActions
      chartRef={chartRef}
      data={chartData}
      filename="revenue-chart"
      title={t('charts.revenue')}
      renderFullscreen={(w, h) => <RevenueChart width={w} height={h} />}
    />
  </CardHeader>
  <CardContent ref={chartRef}>
    <RevenueChart />
  </CardContent>
</Card>
```

## Format Notes

| Format | Engine | Notes |
|--------|--------|-------|
| CSV | Built-in | RFC 4180 quoting, metadata header, `text/csv` MIME |
| Excel | `xlsx` (dynamic) | Data sheet + optional Statistics sheet, auto column widths |
| PDF | `jspdf` (dynamic) | Auto landscape if >5 cols, blue header, page numbers |
| JSON | Built-in | Metadata + column defs for re-import, pretty-printed |
| PNG | `html2canvas` | Charts only, 2x scale, branded dark gradient header |

## Verify

```bash
# TypeScript compiles
npx tsc --noEmit

# Lint the modified files
npx eslint "resources/js/pages/<Entity>/**/*.tsx" --max-warnings=0
```

## Checklist

- [ ] Export fields defined with appropriate formatters
- [ ] Default export fields subset selected (most useful columns)
- [ ] Export handler wired with column filtering and stats
- [ ] ExportDialog added to JSX with all required props
- [ ] Export button added to page actions or toolbar
- [ ] `defaultFilename` set to descriptive slug
- [ ] `preparedBy` passed from user context
- [ ] i18n keys added for export actions
- [ ] Complex fields have formatters (dates, arrays, currency, booleans)
- [ ] Stats mapping included if page has stats cards

## Reference

- Export utilities: `resources/js/utils/exportUtils.ts`
- Export types: `resources/js/types/export.ts`
- Export hook: `resources/js/hooks/useExport.ts`
- Export dialog: `resources/js/components/ExportDialog.tsx`
- Chart actions: `resources/js/components/ui/chart-actions.tsx`
