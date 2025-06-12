import React from 'react';
import { useTheme } from '../context/ThemeContext';

export const ThemeToggle: React.FC = () => {
  const { theme, toggleTheme } = useTheme();

  return (
    <button
      onClick={toggleTheme}
      style={{
        padding: '8px 16px',
        border: '1px solid var(--border-color)',
        backgroundColor: 'var(--background-secondary)',
        color: 'var(--text-primary)',
        cursor: 'pointer',
        borderRadius: '4px',
      }}
    >
      Switch to {theme === 'light' ? 'Dark' : 'Light'} Mode
    </button>
  );
};
