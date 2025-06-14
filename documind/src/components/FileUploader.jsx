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
    <div className={`border-2 border-dashed rounded-lg p-6 text-center bg-white shadow-sm transition-colors ${dragActive ? 'border-blue-500 bg-blue-50' : 'border-gray-300'}`}
      onDragEnter={handleDrag} onDragOver={handleDrag} onDragLeave={handleDrag} onDrop={handleDrop}>
      <input
        type="file"
        accept=".pdf,.docx,application/pdf,application/vnd.openxmlformats-officedocument.wordprocessingml.document"
        className="hidden"
        ref={inputRef}
        onChange={handleChange}
        disabled={isLoading}
      />
      <ArrowUpTrayIcon className="mx-auto h-10 w-10 text-blue-500" />
      <p className="mt-2 text-gray-700">Kéo & thả file PDF/DOCX vào đây hoặc <span className="text-blue-600 underline cursor-pointer" onClick={() => inputRef.current.click()}>chọn file</span> từ máy tính.</p>
    </div>
  );
} 