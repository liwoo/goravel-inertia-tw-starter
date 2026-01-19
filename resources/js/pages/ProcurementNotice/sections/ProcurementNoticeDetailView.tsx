import React, { useEffect, useState } from 'react';
import { Calendar, Building2, FileText, Hash, Users, MapPin, CheckCircle, XCircle, ClipboardList } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { CrudDetailViewProps } from '@/types/crud';
import { ProcurementNotice } from '@/types/procurementnotice';
import axios from 'axios';
import { DetailRow } from '@/components/ui/details-row';
import { Button } from '@/components/ui/button';
import MarkdownEditor from '@uiw/react-markdown-editor';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';

export function ProcurementNoticeDetailView({
    item: notice,
    onEdit,
    onClose,
    canEdit
}: CrudDetailViewProps<ProcurementNotice>) {
    const [smeMap, setSmeMap] = useState<Record<string, string>>({});
    const [isPublishing, setIsPublishing] = useState(false);

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

    const handlePublish = async () => {
        if (!confirm('Are you sure you want to publish this procurement notice?')) return;

        setIsPublishing(true);
        try {
            await axios.put(`/api/procurement-notices/${notice.id}`, {
                ...notice,
                is_published: true
            });
            window.location.reload();
        } catch (error) {
            console.error('Error publishing notice:', error);
            alert('Failed to publish notice');
        } finally {
            setIsPublishing(false);
        }
    };

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-start mb-4">
                <h3 className="text-lg font-semibold text-foreground">Procurement Notice Details</h3>
                {!notice.is_published && canEdit && (
                    <Button
                        onClick={handlePublish}
                        disabled={isPublishing}
                        size="sm"
                        className="bg-green-600 hover:bg-green-700"
                    >
                        {isPublishing ? 'Publishing...' : 'Publish Notice'}
                    </Button>
                )}
            </div>

            <Tabs defaultValue="general" className="w-full">
                <TabsList className="grid w-full grid-cols-2">
                    <TabsTrigger value="general" title="General Information">
                        <FileText className="h-4 w-4 sm:mr-2" />
                        <span className="hidden sm:inline">General Info</span>
                    </TabsTrigger>
                    <TabsTrigger value="application" title="Application Details">
                        <ClipboardList className="h-4 w-4 sm:mr-2" />
                        <span className="hidden sm:inline">Application</span>
                    </TabsTrigger>
                </TabsList>

                <TabsContent value="general" className="space-y-6 mt-4">
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
                            <DetailRow label="Min. Qualifying Score" value={notice.minimum_qualifying_score} />
                            <DetailRow label="Status" render={() => (
                                <Badge variant={notice.is_published ? 'default' : 'secondary'}>
                                    {notice.is_published ? (
                                        <><CheckCircle className="h-4 w-4 text-green-500 mr-1" />Published</>
                                    ) : (
                                        <><XCircle className="h-4 w-4 text-muted-foreground mr-1" />Draft</>
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

                            <DetailRow label="Interested MSMEs" render={() => (
                                <div className="flex flex-wrap gap-2 mt-1">
                                    {notice.interested_smes && notice.interested_smes.length > 0 ? (
                                        notice.interested_smes.map((smeId) => (
                                            <Badge key={smeId} variant="outline">
                                                {smeMap[smeId] || `MSME #${smeId}`}
                                            </Badge>
                                        ))
                                    ) : (
                                        <span className="text-muted-foreground italic">No MSMEs interested yet</span>
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
                </TabsContent>

                <TabsContent value="application" className="space-y-6 mt-4">
                    {/* Details */}
                    <div>
                        <h3 className="text-lg font-semibold mb-4 text-foreground">Details</h3>
                        <div className="grid grid-cols-1 gap-6">
                            <DetailRow label="Details" render={() => (
                                <div className="prose prose-sm max-w-none dark:prose-invert border rounded-md p-4 bg-card">
                                    <MarkdownEditor.Markdown source={notice.details || '-'} style={{ backgroundColor: 'transparent' }} />
                                </div>
                            )} />
                            <DetailRow label="Application Details" value={notice.application_details || '-'} />
                        </div>
                    </div>
                </TabsContent>
            </Tabs>
        </div >
    );
}
