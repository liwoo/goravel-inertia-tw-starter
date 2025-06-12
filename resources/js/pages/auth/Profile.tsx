import React from 'react';
import AdminLayout from "@/layouts/Admin";
// @ts-ignore
import {Head} from "@inertiajs/react";

const Profile = () => {
    return (
        <AdminLayout title="My Profile">
            <Head title="My Profile" />
            <div>Profile</div>
        </AdminLayout>
    );
};

export default Profile;