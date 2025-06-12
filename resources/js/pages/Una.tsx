import React from 'react';
// @ts-ignore
import { Head } from '@inertiajs/react';
import AuthLayout from "@/layouts/Auth";

interface UnaProps {
  title: string;
}

const Una: React.FC<UnaProps> = ({ title }) => {
  return (
    <AuthLayout>
      <Head title={title} />
      <div className="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10">
        <h1 className="text-2xl font-bold text-center text-gray-800 mb-6">{title}</h1>
        <p className="text-center text-gray-600">
          You do not have permission to access this page. 
          Please <a href="/" className="font-medium text-indigo-600 hover:text-indigo-500">login</a> to continue.
        </p>
      </div>
    </AuthLayout>
  );
};

export default Una;
