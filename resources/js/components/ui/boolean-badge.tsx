import { Badge } from "./badge";
import { CheckCircle2, XCircle } from "lucide-react";

export const BooleanBadge = ({ value, trueLabel = 'Yes', falseLabel = 'No' }: { value: boolean; trueLabel?: string; falseLabel?: string }) => (
    <Badge variant={value ? 'default' : 'secondary'} className="gap-1">
        {value ? <CheckCircle2 className="h-3 w-3" /> : <XCircle className="h-3 w-3" />}
        {value ? trueLabel : falseLabel}
    </Badge>
);