import React from 'react';
import AdminLayout from "@/layouts/Admin";
// @ts-ignore
import { Head, usePage } from "@inertiajs/react";
import { ProfileForm } from "@/components/profile-form";
import type { SharedData, User } from "@/types/app";

interface ProfilePageProps extends SharedData {
    user: User;
}

const Profile = () => {
    const { props } = usePage<ProfilePageProps>();
    const user = props.auth?.user;

    return (
        <AdminLayout title="My Profile">
            <Head title="My Profile" />
            <div className="p-6">
                <ProfileForm user={user} />
            </div>
        </AdminLayout>
    );
};

export default Profile;