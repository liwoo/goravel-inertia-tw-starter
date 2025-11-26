import React, { useState } from 'react';
import { Head } from '@inertiajs/react';
import { ApplicationCreateForm } from './sections/ApplicationCreateForm';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { ArrowLeft } from 'lucide-react';

export default function ApplicationPage() {
    const [isSuccess, setIsSuccess] = useState(false);
    const formRef = React.useRef<any>(null);

    const handleSuccess = () => {
        setIsSuccess(true);
    };

    const handleSubmit = () => {
        if (formRef.current) {
            formRef.current.handleSubmit();
        }
    };

    if (isSuccess) {
        return (
            <div className="min-h-screen bg-background flex items-center justify-center p-4">
                <Head title="Application Submitted" />
                <Card className="w-full max-w-2xl">
                    <CardHeader>
                        <CardTitle className="text-center text-2xl text-green-600">Application Submitted Successfully!</CardTitle>
                        <CardDescription className="text-center">
                            Thank you for your application. We have received your details and will process them shortly.
                        </CardDescription>
                    </CardHeader>
                    <CardContent className="flex justify-center">
                        <Button variant="outline" onClick={() => window.location.href = '/'}>
                            <ArrowLeft className="mr-2 h-4 w-4" />
                            Back to Home
                        </Button>
                    </CardContent>
                </Card>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-background py-12 px-4 sm:px-6 lg:px-8">
            <Head title="Submit Application" />

            <div className="max-w-3xl mx-auto">
                <div className="mb-8 text-center">
                    <h1 className="text-3xl font-bold tracking-tight text-foreground">Submit Application</h1>
                    <p className="mt-2 text-muted-foreground">
                        Please fill out the form below to submit your application.
                    </p>
                </div>

                <Card>
                    <CardContent className="pt-6">
                        <ApplicationCreateForm
                            ref={formRef}
                            onSuccess={handleSuccess}
                            onCancel={() => window.history.back()}
                        />
                        <div className="mt-6 flex justify-end gap-4">
                            <Button variant="outline" onClick={() => window.history.back()}>
                                Cancel
                            </Button>
                            <Button onClick={handleSubmit}>
                                Submit Application
                            </Button>
                        </div>
                    </CardContent>
                </Card>
            </div>
        </div>
    );
}
