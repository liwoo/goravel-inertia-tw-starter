import React, { useEffect, useState } from 'react';
import { Calendar, Building2, FileText, Hash, Users, MapPin, CheckCircle, XCircle } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { CrudDetailViewProps } from '@/types/crud';
import { ProcurementNotice } from '@/types/procurementnotice';
import axios from 'axios';
import { DetailRow } from '@/components/ui/details-row';

export function ProcurementNoticeDetailView({
    item: notice,
    onEdit,
    onClose,
    canEdit
}: CrudDetailViewProps<ProcurementNotice>) {
    const [smeMap, setSmeMap] = useState<Record<string, string>>({});

    useEffect(() => {
        const fetchSmes = async () => {
            try {
                const response = await axios.get('/api/smes', {
                    params: {
                        pageSize: 1000,
                    }
                });
                const smes = response.data.data?.data || [];
                const map: Record<string, string> = {};
                smes.forEach((sme: any) => {
                    map[sme.id.toString()] = sme.business_name;
                });
                setSmeMap(map);
            } catch (error) {
                console.error('Error fetching SMEs:', error);
            }
        };

        if (notice.interested_smes && notice.interested_smes.length > 0) {
            fetchSmes();
        }
    }, [notice.interested_smes]);

    const formatDate = (date: string | Date | null) => {
        if (!date) return 'Not specified';
        return new Date(date).toLocaleDateString('en-US', {
            month: 'long',
            day: 'numeric',
            year: 'numeric',
        });
    };

    return (
        <div className="space-y-6">
            <div className="space-y-6">

                {/* Basic Info */}
                <div>
                    <h3 className="text-lg font-semibold mb-4 text-foreground">Basic Information</h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <DetailRow icon={Building2} label="Procured By" value={notice.procured_by} />
                        <DetailRow icon={Building2} label="Organization" value={notice.organization} />
                        <DetailRow label="Period" value={formatDate(notice.open_date) + " - " + formatDate(notice.close_date)} />
                        <DetailRow label="Type" value={notice.procurement_type} />
                        <DetailRow label="Reference Number" value={notice.ref_no} />
                        <DetailRow label="Market Approach" value={notice.market_approach} />
                        <DetailRow label="Invitation Type" value={notice.invitation} />
                        <DetailRow label="Reference Number" value={notice.ref_no} />
                        <DetailRow label="Market Approach" value={notice.market_approach} />
                        <DetailRow label="Invitation Type" value={notice.invitation} />
                        <DetailRow label="Min. Qualifying Score" value={notice.minimum_qualifying_score} />
                    </div>
                </div>

                <Separator />
                {/* Details */}
                <div>
                    <h3 className="text-lg font-semibold mb-4 text-foreground">Details</h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <DetailRow label="Details" value={notice.details || '-'} />
                        <DetailRow label="Application Details" value={notice.application_details || '-'} />
                        <DetailRow label="Status" render={() => (
                            <Badge variant={notice.is_published ? 'default' : 'secondary'}>
                                {notice.is_published ? (
                                    <><CheckCircle className="h-4 w-4 text-green-500" />Published</>
                                ) : (
                                    <><XCircle className="h-4 w-4 text-muted-foreground" />Draft</>
                                )}
                            </Badge>
                        )} />
                    </div>
                </div>

                <Separator />

                <div>
                    <h3 className="text-lg font-semibold mb-4 text-foreground">Classification & Location</h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <DetailRow label="Classification" render={() => (
                            <div className="flex flex-wrap gap-2 mt-1">
                                {notice.classification && notice.classification.length > 0 ? (
                                    notice.classification.map((item, index) => (
                                        <Badge key={index} variant="secondary">
                                            {item}
                                        </Badge>
                                    ))
                                ) : (
                                    <span className="text-muted-foreground italic">None specified</span>
                                )}
                            </div>
                        )} />

                        <DetailRow label="Qualifying Districts" render={() => (
                            <div className="flex flex-wrap gap-2 mt-1">
                                {notice.qualifying_districts && notice.qualifying_districts.length > 0 ? (
                                    notice.qualifying_districts.map((item, index) => (
                                        <Badge key={index} variant="secondary">
                                            {item}
                                        </Badge>
                                    ))
                                ) : (
                                    <span className="text-muted-foreground italic">None specified</span>
                                )}
                            </div>
                        )} />
                    </div>
                </div>

                <Separator />

                <div>
                    <h3 className="text-lg font-semibold mb-4 text-foreground">Participants</h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <DetailRow label="Partners" render={() => (
                            <div className="flex flex-wrap gap-2 mt-1">
                                {notice.partners && notice.partners.length > 0 ? (
                                    notice.partners.map((partner, index) => (
                                        <Badge key={index} variant="secondary">
                                            {partner}
                                        </Badge>
                                    ))
                                ) : (
                                    <span className="text-muted-foreground italic">No partners selected</span>
                                )}
                            </div>
                        )} />

                        <DetailRow label="Interested SMEs" render={() => (
                            <div className="flex flex-wrap gap-2 mt-1">
                                {notice.interested_smes && notice.interested_smes.length > 0 ? (
                                    notice.interested_smes.map((smeId) => (
                                        <Badge key={smeId} variant="outline">
                                            {smeMap[smeId] || `SME #${smeId}`}
                                        </Badge>
                                    ))
                                ) : (
                                    <span className="text-muted-foreground italic">No SMEs selected</span>
                                )}
                            </div>
                        )} />
                    </div>
                </div>

                <Separator />

                {/* Metadata Section */}
                <div>
                    <h3 className="text-lg font-semibold mb-4 text-foreground">Metadata</h3>
                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <p className="text-sm text-muted-foreground">Created</p>
                            <p className="font-medium text-sm text-foreground">{formatDate(notice.createdAt || notice.created_at || null)}</p>
                        </div>
                        <div>
                            <p className="text-sm text-muted-foreground">Last Updated</p>
                            <p className="font-medium text-sm text-foreground">{formatDate(notice.updatedAt || notice.updated_at || null)}</p>
                        </div>
                    </div>
                </div>
            </div>
        </div >
    );
}
