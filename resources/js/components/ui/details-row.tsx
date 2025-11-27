export const DetailRow = ({ icon: Icon, label, value, render }: { icon?: any; label: string; value?: any, render?: () => React.ReactNode }) => (
    <div className="flex items-start gap-3">
        {Icon && <div className="p-2 rounded-lg bg-muted">
            <Icon className="h-4 w-4 text-muted-foreground" />
        </div>}
        <div className="flex-1 space-y-1">
            <p className="text-sm text-muted-foreground">{label}</p>
            {render ? render() : <p className="font-medium text-foreground">{value || '-'}</p>}
        </div>
    </div>
);