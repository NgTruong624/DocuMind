import React, { useRef, useState } from 'react';
import { ArrowUpTrayIcon } from '@heroicons/react/24/outline';

const MAX_FILE_SIZE = 10 * 1024 * 1024; // 10MB in bytes

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
    // Kiểm tra cả MIME type và extension
    const allowedTypes = [
      'application/pdf',
      'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
      'application/msword'
    ];
    const allowedExtensions = ['.pdf', '.doc', '.docx'];
    
    const fileExtension = '.' + file.name.split('.').pop().toLowerCase();
    const isValidType = allowedTypes.includes(file.type) || allowedExtensions.includes(fileExtension);

    if (!isValidType) {
      alert('Chỉ chấp nhận file PDF hoặc DOCX.');
      return false;
    }

    if (file.size > MAX_FILE_SIZE) {
      alert(`File quá lớn. Kích thước tối đa là ${MAX_FILE_SIZE / 1024 / 1024}MB.`);
      return false;
    }

    return true;
  };

  const handleButtonClick = () => {
    inputRef.current?.click();
  };

  return (
    <div
      className={`border-2 border-dashed rounded-lg p-12 text-center transition-colors ${
        dragActive ? 'border-blue-500 bg-blue-50' : 'border-gray-300 hover:border-blue-500 hover:bg-slate-50'
      }`}
      onDragEnter={handleDrag}
      onDragLeave={handleDrag}
      onDragOver={handleDrag}
      onDrop={handleDrop}
    >
      <ArrowUpTrayIcon className="w-12 h-12 mx-auto text-slate-400" />
      <p className="mt-4 text-slate-500">
        Kéo và thả file vào đây, hoặc{' '}
        <span className="font-semibold text-blue-600 cursor-pointer" onClick={handleButtonClick}>
          chọn file
        </span>
      </p>
      <input
        type="file"
        accept=".pdf,.doc,.docx,application/pdf,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document"
        className="hidden"
        ref={inputRef}
        onChange={handleChange}
        disabled={isLoading}
      />
    </div>
  );
} 