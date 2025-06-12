import React from 'react';
import { Button } from '@/components/ui/button'; // Import the Shadcn Button

interface HomeProps {
  version?: string;
  // Add other props your component might receive
}

function Home(props: HomeProps) {
  console.log('Home.tsx rendering, props:', props);
  return (
    <div className="p-4">
      <h1 className="text-3xl font-bold text-blue-600">Test Home Works!</h1>
      <p className="mb-4">Version: {props.version || 'N/A'}</p>
      <Button onClick={() => alert('Shadcn Button Clicked!')}>
        Click Me
      </Button>
    </div>
  );
}

export default Home;
