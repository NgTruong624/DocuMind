import React, { useState } from 'react';
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Upload, FileText, CheckCircle, AlertTriangle } from "lucide-react";
import DocumentChat from './DocumentChat';
import { analyzeContract, AnalyzeResult } from '@/api/analysisService';

const DocumentUpload: React.FC = () => {
  const [isUploading, setIsUploading] = useState(false);
  const [isAnalyzing, setIsAnalyzing] = useState(false);
  const [analysisResult, setAnalysisResult] = useState<AnalyzeResult | null>(null);
  const [uploadedFile, setUploadedFile] = useState<File | null>(null);
  const [error, setError] = useState<string | null>(null);

  const handleFileUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;
    setUploadedFile(file);
    setIsUploading(true);
    setError(null);
    try {
      const result = await analyzeContract(file);
      setAnalysisResult(result);
    } catch (err: any) {
      setError(err?.response?.data?.message || 'Đã xảy ra lỗi khi phân tích hợp đồng.');
    } finally {
      setIsUploading(false);
      setIsAnalyzing(false);
    }
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
  };

  const handleDrop = async (e: React.DragEvent) => {
    e.preventDefault();
    const files = e.dataTransfer.files;
    if (files.length > 0) {
      const file = files[0];
      setUploadedFile(file);
      setIsUploading(true);
      setError(null);
      try {
        const result = await analyzeContract(file);
        setAnalysisResult(result);
      } catch (err: any) {
        setError(err?.response?.data?.message || 'Đã xảy ra lỗi khi phân tích hợp đồng.');
      } finally {
        setIsUploading(false);
        setIsAnalyzing(false);
      }
    }
  };

  const resetUpload = () => {
    setUploadedFile(null);
    setAnalysisResult(null);
    setIsUploading(false);
    setIsAnalyzing(false);
    setError(null);
  };

  if (analysisResult) {
    return (
      <div className="space-y-6">
        <Card className="border-l-4 border-l-green-500 bg-green-50">
          <CardContent className="p-6">
            <div className="flex items-start gap-4">
              <CheckCircle className="w-6 h-6 text-green-600 flex-shrink-0 mt-1" />
              <div className="flex-1">
                <h3 className="text-lg font-semibold text-green-900 mb-2">
                  Analysis Complete: {uploadedFile?.name}
                </h3>
                <div className="grid md:grid-cols-3 gap-4 mb-4">
                  <div className="text-center p-3 bg-white rounded-lg">
                    <div className="text-2xl font-bold text-slate-900">{analysisResult.riskScore}</div>
                    <div className="text-sm text-slate-600">Risk Level</div>
                  </div>
                  <div className="text-center p-3 bg-white rounded-lg">
                    <div className="text-2xl font-bold text-slate-900">{analysisResult.keyTerms?.length}</div>
                    <div className="text-sm text-slate-600">Key Terms</div>
                  </div>
                  <div className="text-center p-3 bg-white rounded-lg">
                    <div className="flex items-center justify-center gap-1">
                      {analysisResult.riskScore === 'Low' ? (
                        <CheckCircle className="w-5 h-5 text-green-500" />
                      ) : (
                        <AlertTriangle className="w-5 h-5 text-yellow-500" />
                      )}
                      <span className="text-2xl font-bold text-slate-900">{analysisResult.riskScore}</span>
                    </div>
                    <div className="text-sm text-slate-600">Overall Risk</div>
                  </div>
                </div>
                <p className="text-slate-700 mb-4">{analysisResult.summary}</p>
                <div className="flex gap-2">
                  <Button 
                    onClick={resetUpload}
                    variant="outline" 
                    size="sm"
                  >
                    Upload Another Document
                  </Button>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        <DocumentChat
          documentSummary={analysisResult.summary}
          keyTerms={analysisResult.keyTerms}
          riskScore={analysisResult.riskScore}
          fileHash={analysisResult.file_hash}
        />
      </div>
    );
  }

  return (
    <Card className="max-w-2xl mx-auto">
      <CardContent className="p-8">
        <div
          className="border-2 border-dashed border-slate-300 rounded-lg p-8 text-center hover:border-blue-400 transition-colors"
          onDragOver={handleDragOver}
          onDrop={handleDrop}
        >
          {isUploading ? (
            <div className="space-y-4">
              <div className="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto">
                <Upload className="w-8 h-8 text-blue-600 animate-pulse" />
              </div>
              <div>
                <h3 className="text-lg font-semibold text-slate-900">Uploading...</h3>
                <p className="text-slate-600">Please wait while we process your document</p>
              </div>
            </div>
          ) : isAnalyzing ? (
            <div className="space-y-4">
              <div className="w-16 h-16 bg-indigo-100 rounded-full flex items-center justify-center mx-auto">
                <FileText className="w-8 h-8 text-indigo-600 animate-pulse" />
              </div>
              <div>
                <h3 className="text-lg font-semibold text-slate-900">Analyzing Contract...</h3>
                <p className="text-slate-600">Our AI is reading and analyzing your document</p>
                <div className="w-48 h-2 bg-slate-200 rounded-full mx-auto mt-4">
                  <div className="h-2 bg-gradient-to-r from-blue-500 to-indigo-500 rounded-full animate-pulse" style={{ width: '70%' }}></div>
                </div>
              </div>
            </div>
          ) : (
            <>
              <input
                type="file"
                accept=".pdf,.doc,.docx"
                onChange={handleFileUpload}
                className="hidden"
                id="file-upload-input"
              />
              <label htmlFor="file-upload-input" className="cursor-pointer">
                <div className="flex flex-col items-center justify-center gap-2">
                  <Upload className="w-10 h-10 text-blue-600" />
                  <span className="text-lg font-semibold text-blue-600">Click or drag file to upload</span>
                  <span className="text-sm text-slate-500">PDF, DOC, DOCX supported</span>
                </div>
              </label>
              {error && <div className="text-red-600 mt-4">{error}</div>}
            </>
          )}
        </div>
      </CardContent>
    </Card>
  );
};

export default DocumentUpload;
