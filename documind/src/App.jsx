import React, { useState } from 'react';
import FileUploader from './components/FileUploader';
import LoadingSpinner from './components/LoadingSpinner';
import ResultsDisplay from './components/ResultsDisplay';
import FileInfo from './components/FileInfo';
import { analyzeContract } from './api/analysisService';

export default function App() {
  const [file, setFile] = useState(null);
  const [analysisResult, setAnalysisResult] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleFileSelected = async (selectedFile) => {
    setFile(selectedFile);
    setAnalysisResult(null);
    setError(null);
    setIsLoading(true);
    try {
      const result = await analyzeContract(selectedFile);
      setAnalysisResult(result);
    } catch (err) {
      let errorMessage = 'Đã xảy ra lỗi khi phân tích hợp đồng.';
      
      if (err.response?.status === 413) {
        errorMessage = 'File quá lớn. Vui lòng chọn file nhỏ hơn.';
      } else if (err.response?.status === 415) {
        errorMessage = 'Định dạng file không được hỗ trợ. Vui lòng chọn file PDF hoặc DOCX.';
      } else if (err.response?.status === 500) {
        errorMessage = 'File quá lớn. Vui lòng chọn file nhỏ hơn.';
      }
      
      console.error('Analysis error:', err);
      setError(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  const handleRemoveFile = () => {
    setFile(null);
    setAnalysisResult(null);
    setError(null);
  };

  return (
    <div className="min-h-screen bg-slate-100 flex items-center justify-center">
      <div className="w-full max-w-3xl mx-auto bg-white rounded-xl shadow-lg p-8 mt-8 mb-8">
        <h1 className="text-2xl md:text-3xl font-bold text-center mb-8 text-gray-800">
          Công cụ Phân tích & Tóm tắt Hợp đồng
        </h1>
        <FileUploader onFileSelected={handleFileSelected} isLoading={isLoading} />
        {file && <FileInfo file={file} onRemove={handleRemoveFile} />}
        {isLoading && <LoadingSpinner />}
        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded-lg text-center mt-4">
            {error}
          </div>
        )}
        <ResultsDisplay analysisResult={analysisResult} />
      </div>
    </div>
  );
}
