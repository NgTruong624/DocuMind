import React from 'react';
import { PaperClipIcon, XMarkIcon } from '@heroicons/react/24/outline';

export default function FileInfo({ file, onRemove }) {
  if (!file) return null;
  return (
    <div className="flex items-center justify-between bg-slate-50 border border-slate-200 rounded-lg px-4 py-3 mb-4">
      <div className="flex items-center gap-3">
        <PaperClipIcon className="w-6 h-6 text-blue-500" />
        <span className="font-medium text-slate-800 truncate max-w-xs" title={file.name}>{file.name}</span>
        <span className="text-xs text-slate-400 ml-2">{(file.size / 1024).toFixed(1)} KB</span>
      </div>
      {onRemove && (
        <button
          onClick={onRemove}
          className="ml-4 p-1 rounded bg-red-100 hover:bg-red-200 active:bg-red-500 transition-colors"
        >
          <XMarkIcon className="w-5 h-5 text-red-500" />
        </button>
      )}
    </div>
  );
}
