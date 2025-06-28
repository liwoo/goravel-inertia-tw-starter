import React from 'react';

export const ThemeToggle: React.FC = () => {
  return (
    <button
      style={{
        padding: '8px 16px',
        border: '1px solid var(--border-color)',
        backgroundColor: 'var(--background-secondary)',
        color: 'var(--text-primary)',
        cursor: 'pointer',
        borderRadius: '4px',
      }}
    >
      Theme Toggle
    </button>
  );
};
