import React from 'react';

export default function LoadingSpinner() {
  return (
    <div className="flex flex-col items-center justify-center py-8">
      <div className="animate-spin rounded-full h-12 w-12 border-t-4 border-blue-500 border-solid mb-4"></div>
      <span className="text-lg text-slate-600 font-medium">Analyzing...</span>
    </div>
  );
} 