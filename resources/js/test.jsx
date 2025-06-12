import React from 'react';
import { createRoot } from 'react-dom/client';

// A simple test component
const TestComponent = () => {
  return (
    <div style={{ padding: '20px', fontFamily: 'sans-serif', background: '#f0f0f0', borderRadius: '8px' }}>
      <h1>Test Component</h1>
      <p>If you can see this, React and JSX are working correctly!</p>
    </div>
  );
};

// Wait for the DOM to be fully loaded
document.addEventListener('DOMContentLoaded', () => {
  console.log('Test script loaded!');
  
  // Try to find a test element or create one
  let testElement = document.getElementById('test-root');
  if (!testElement) {
    testElement = document.createElement('div');
    testElement.id = 'test-root';
    document.body.appendChild(testElement);
  }
  
  // Render the test component
  createRoot(testElement).render(<TestComponent />);
});
