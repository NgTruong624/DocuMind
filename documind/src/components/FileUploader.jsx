import React, { useRef, useState } from 'react';
import { ArrowUpTrayIcon } from '@heroicons/react/24/outline';

export default function FileUploader({ onFileSelected, isLoading }) {
  const [dragActive, setDragActive] = useState(false);
  const inputRef = useRef();

  const handleDrag = (e) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === 'dragenter' || e.type === 'dragover') setDragActive(true);
    else if (e.type === 'dragleave') setDragActive(false);
  };

  const handleDrop = (e) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);
    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      const file = e.dataTransfer.files[0];
      if (validateFile(file)) onFileSelected(file);
    }
  };

  const handleChange = (e) => {
    const file = e.target.files[0];
    if (file && validateFile(file)) onFileSelected(file);
  };

  const validateFile = (file) => {
    const allowed = ['application/pdf', 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'];
    if (!allowed.includes(file.type)) {
      alert('Chỉ chấp nhận file PDF hoặc DOCX.');
      return false;
    }
    return true;
  };

  return (
    <div className="border-2 border-dashed border-gray-300 rounded-lg p-12 text-center transition-colors hover:border-blue-500 hover:bg-slate-50">
      <ArrowUpTrayIcon className="w-12 h-12 mx-auto text-slate-400" />
      <p className="mt-4 text-slate-500">
        Drag and drop a file here, or <span className="font-semibold text-blue-600 cursor-pointer">choose file</span>
      </p>
      <input
        type="file"
        accept=".pdf,.docx,application/pdf,application/vnd.openxmlformats-officedocument.wordprocessingml.document"
        className="hidden"
        ref={inputRef}
        onChange={handleChange}
        disabled={isLoading}
      />
    </div>
  );
} 