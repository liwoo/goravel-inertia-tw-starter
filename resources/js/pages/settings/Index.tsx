//index page

import React from 'react';
import AdminLayout from "@/layouts/Admin";
// @ts-ignore
import {Head} from "@inertiajs/react";

const SettingsPage: React.FC = () => {
    return (
        <AdminLayout title="Settings">
            <Head title="Settings" />
            <h1>Settings</h1>
        </AdminLayout>
    );
};

export default SettingsPage;