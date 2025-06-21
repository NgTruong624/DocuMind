// Thêm khai báo cho ImportMetaEnv nếu chưa có (có thể đặt ở src/vite-env.d.ts)
// declare interface ImportMeta {
//   env: ImportMetaEnv;
// }
// declare interface ImportMetaEnv {
//   VITE_API_URL?: string;
//   VITE_CHAT_API_URL?: string;
// }

import axios from 'axios';

const API_URL = (import.meta as any).env.VITE_API_URL || 'http://localhost:8080/api/v1/analyze';

// Định nghĩa kiểu dữ liệu trả về từ API phân tích hợp đồng
export interface AnalyzeResult {
  summary: string;
  key_clauses?: string[];
  potential_risks?: string[];
  file_hash?: string;
  [key: string]: any; // fallback cho các trường khác nếu có
}

// Định nghĩa kiểu dữ liệu trả về từ API chat
export interface ContractChatResponse {
  answer: string;
  [key: string]: any;
}

export async function analyzeContract(file: File): Promise<AnalyzeResult> {
  try {
    const formData = new FormData();
    formData.append('file', file);
    
    console.log('Uploading file:', {
      name: file.name,
      type: file.type,
      size: `${(file.size / 1024 / 1024).toFixed(2)}MB`
    });

    const response = await axios.post<AnalyzeResult>(API_URL, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  } catch (error: any) {
    console.error('Error analyzing contract:', {
      message: error.message,
      response: error.response?.data,
      status: error.response?.status
    });
    throw error;
  }
}

export interface ContractChatParams {
  fileHash?: string;
  contractText?: string;
  question: string;
}

export async function contractChat({ fileHash, contractText, question }: ContractChatParams): Promise<ContractChatResponse> {
  try {
    const response = await axios.post<ContractChatResponse>(
      ((import.meta as any).env.VITE_CHAT_API_URL || 'http://localhost:8080/api/v1/contract-chat'),
      {
        file_hash: fileHash || '',
        contract_text: contractText || '',
        question,
      }
    );
    return response.data;
  } catch (error) {
    throw error;
  }
} 