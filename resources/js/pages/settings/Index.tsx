//index page

import React from 'react';
import AdminLayout from "@/layouts/Admin";
// @ts-ignore
import {Head} from "@inertiajs/react";
import {useTranslation} from 'react-i18next';

const SettingsPage: React.FC = () => {
    const {t} = useTranslation('settings');

    return (
        <AdminLayout title={t('title')}>
            <Head title={t('title')} />
            <h1>{t('title')}</h1>
        </AdminLayout>
    );
};

export default SettingsPage;
