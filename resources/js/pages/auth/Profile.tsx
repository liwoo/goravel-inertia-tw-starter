import React from 'react';
import AdminLayout from '@/layouts/Admin';
import { Head, usePage } from '@inertiajs/react';
import { SharedData } from '@/types/app';

// Import section components
import ProfileSection from './sections/ProfileSection';
import SecuritySection from './sections/SecuritySection';
import ActivitySection from './sections/ActivitySection';

const Profile: React.FC = () => {
  const { props } = usePage<SharedData>();
  const user = props.auth?.user;

  return (
    <AdminLayout title="My Account">
      <Head title="My Account" />
      <div className="p-4 md:p-6 lg:p-8">
        {/* Page Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-foreground">Account Settings</h1>
          <p className="text-muted-foreground mt-2">
            Manage your profile, security settings, and view your activity history
          </p>
        </div>

        {/* Two-column layout */}
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 lg:gap-8">
          {/* Left column - Main content (Profile + Security) */}
          <div className="lg:col-span-8 space-y-6">
            <ProfileSection user={user} />
            <SecuritySection userId={user?.id} />
          </div>

          {/* Right column - Activity sidebar */}
          <div className="lg:col-span-4">
            <div className="lg:sticky lg:top-6">
              <div className="max-h-[calc(100vh-8rem)] overflow-y-auto">
                <ActivitySection userId={user?.id} />
              </div>
            </div>
          </div>
        </div>
      </div>
    </AdminLayout>
  );
};

export default Profile;
