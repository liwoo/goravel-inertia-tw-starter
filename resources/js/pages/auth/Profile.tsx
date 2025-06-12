import React from 'react';
import AdminLayout from "@/layouts/Admin";
// @ts-ignore
import {Head} from "@inertiajs/react";

const Profile = () => {
    return (
        <AdminLayout title="My Profile">
            <Head title="My Profile" />
            <div className="p-4">
                <h1 className="text-xl mt-4">Profile Page Content</h1>
                {/* Add more profile content here */}
            </div>
        </AdminLayout>
    );
};

export default Profile;