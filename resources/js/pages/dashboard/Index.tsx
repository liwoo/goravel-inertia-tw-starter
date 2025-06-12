import React from 'react';
// @ts-ignore
import { Head, usePage } from '@inertiajs/react';
import Admin from '@/layouts/Admin';
import type {SharedData, User} from "@/types/app";

type DashboardPageProps = SharedData & {}

const DashboardPage: React.FC = () => {
  const { props } = usePage<DashboardPageProps>();
  const user = props.auth?.user;

  return (
    <Admin>
      <Head>
        <title>{props.pageTitle || 'Dashboard'}</title>
      </Head>

      <div className="overflow-hidden shadow-sm sm:rounded-lg">
        <div className="p-6 border-b border-gray-200">
          {user ? (
            <h1 className="text-2xl font-semibold">
              Welcome back, {user.name}!
            </h1>
          ) : (
            <h1 className="text-2xl font-semibold">
              Welcome to the {props.pageTitle}!
            </h1>
          )}
          <p className="mt-2">
            This is your {props.pageTitle}. You are logged in as {user?.role.toLowerCase()}.
          </p>
          {/* Add more dashboard content here */}
        </div>
      </div>
    </Admin>
  );
};

export default DashboardPage;
