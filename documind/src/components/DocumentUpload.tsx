import React, { useState } from 'react';
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Upload, FileText, CheckCircle, AlertTriangle } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { ScrollArea } from "@/components/ui/scroll-area";
import ErrorMessage from "@/components/ui/error-message";
import DocumentChat from './DocumentChat';
import { analyzeContract, AnalyzeResult } from '@/api/analysisService';

const DocumentUpload: React.FC = () => {
  const [isUploading, setIsUploading] = useState(false);
  const [isAnalyzing, setIsAnalyzing] = useState(false);
  const [analysisResult, setAnalysisResult] = useState<AnalyzeResult | null>(null);
  const [uploadedFile, setUploadedFile] = useState<File | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [modalTitle, setModalTitle] = useState('');
  const [modalContent, setModalContent] = useState<string[]>([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isReUploading, setIsReUploading] = useState(false);

  const showDetails = (title: string, content: string[] | undefined) => {
    if (content && content.length > 0) {
        setModalTitle(title);
        setModalContent(content);
        setIsModalOpen(true);
    }
  };

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

  const handleReUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    setUploadedFile(file);
    setIsReUploading(true);
    setError(null);

    try {
      const result = await analyzeContract(file);
      setAnalysisResult(result);
    } catch (err: any) {
      setError(err?.response?.data?.message || 'Đã xảy ra lỗi khi phân tích hợp đồng.');
    } finally {
      setIsReUploading(false);
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
        <Dialog open={isModalOpen} onOpenChange={setIsModalOpen}>
          <DialogContent className="sm:max-w-[625px]">
            <DialogHeader>
              <DialogTitle>{modalTitle}</DialogTitle>
            </DialogHeader>
            <ScrollArea className="max-h-[60vh] pr-4">
              <ul className="list-disc pl-5 space-y-2">
                {modalContent.map((item, index) => (
                  <li key={index} className="text-slate-700">{item}</li>
                ))}
              </ul>
            </ScrollArea>
          </DialogContent>
        </Dialog>
        <Card className="border-l-4 border-l-green-500 bg-green-50">
          <CardContent className="p-6">
            <div className="flex items-start gap-4">
              <CheckCircle className="w-6 h-6 text-green-600 flex-shrink-0 mt-1" />
              <div className="flex-1">
                <h3 className="text-lg font-semibold text-green-900 mb-2">
                  Analysis Complete: {uploadedFile?.name}
                </h3>
                <div className="grid md:grid-cols-2 gap-4 mb-4">
                  <div className="text-center p-3 bg-white rounded-lg cursor-pointer hover:bg-slate-50 transition-colors" onClick={() => showDetails('Điều khoản chính', analysisResult.key_clauses)}>
                    <div className="text-2xl font-bold text-slate-900">{analysisResult.key_clauses?.length ?? 0}</div>
                    <div className="text-sm text-slate-600">Điều khoản chính</div>
                  </div>
                  <div className="text-center p-3 bg-white rounded-lg cursor-pointer hover:bg-slate-50 transition-colors" onClick={() => showDetails('Rủi ro tiềm ẩn', analysisResult.potential_risks)}>
                    <div className="text-2xl font-bold text-slate-900">{analysisResult.potential_risks?.length ?? 0}</div>
                    <div className="text-sm text-slate-600">Rủi ro tiềm ẩn</div>
                  </div>
                </div>
                <p className="text-slate-700 mb-4">{analysisResult.summary}</p>
                <div className="flex gap-2">
                  <input
                    id="re-upload-input"
                    type="file"
                    className="hidden"
                    onChange={handleReUpload}
                    accept=".pdf,.doc,.docx"
                  />
                  <Button
                    onClick={() => document.getElementById('re-upload-input')?.click()}
                    variant="outline"
                    size="sm"
                    disabled={isReUploading}
                  >
                    {isReUploading ? 'Đang tải lên...' : 'Tải lên tài liệu khác'}
                  </Button>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        <DocumentChat
          key={analysisResult.file_hash}
          documentSummary={analysisResult.summary}
          keyClauses={analysisResult.key_clauses}
          potentialRisks={analysisResult.potential_risks}
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
              {error && (
                <ErrorMessage message={error} />
              )}
            </>
          )}
        </div>
      </CardContent>
    </Card>
  );
};

export default DocumentUpload;
