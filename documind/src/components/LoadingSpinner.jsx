import React from 'react';
import { ClipLoader } from 'react-spinners';

export default function LoadingSpinner() {
  return (
    <div className="flex flex-col items-center justify-center py-8">
      <ClipLoader color="#3b82f6" size={48} speedMultiplier={1.2} />
      <span className="text-lg text-slate-600 font-medium mt-4">Đang phân tích...</span>
    </div>
  );
} 