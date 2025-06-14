import React, { useState } from 'react';
import FileUploader from './components/FileUploader';
import LoadingSpinner from './components/LoadingSpinner';
import ResultsDisplay from './components/ResultsDisplay';
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
      setError('Đã xảy ra lỗi khi phân tích hợp đồng.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-100 py-8 px-2 md:px-0">
      <div className="max-w-2xl mx-auto">
        <h1 className="text-2xl md:text-3xl font-bold text-center mb-8 text-gray-800">Công cụ Phân tích & Tóm tắt Hợp đồng</h1>
        <FileUploader onFileSelected={handleFileSelected} isLoading={isLoading} />
        {isLoading && <LoadingSpinner />}
        {error && <div className="text-red-600 text-center mt-4">{error}</div>}
        <ResultsDisplay analysisResult={analysisResult} />
      </div>
    </div>
  );
}
