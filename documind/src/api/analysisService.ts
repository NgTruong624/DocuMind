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
    
    // Xử lý lỗi cụ thể
    if (error.response?.status === 500) {
      const errorMessage = error.response?.data?.error || error.message;
      
      // Kiểm tra lỗi quota API
      if (errorMessage.includes('quota') || errorMessage.includes('429') || errorMessage.includes('exceeded')) {
        throw new Error('API quota đã hết. Vui lòng thử lại sau hoặc liên hệ admin để nâng cấp quota.');
      }
      
      // Kiểm tra lỗi khác
      if (errorMessage.includes('AI analysis failed')) {
        throw new Error('Lỗi phân tích AI. Vui lòng thử lại sau.');
      }
      
      if (errorMessage.includes('Could not extract text')) {
        throw new Error('Không thể đọc nội dung file. Vui lòng kiểm tra định dạng file.');
      }
      
      if (errorMessage.includes('File upload failed')) {
        throw new Error('Lỗi tải file lên. Vui lòng thử lại.');
      }
      
      // Lỗi chung
      throw new Error(`Lỗi server: ${errorMessage}`);
    }
    
    if (error.response?.status === 400) {
      throw new Error('Định dạng file không được hỗ trợ. Vui lòng sử dụng PDF, DOC hoặc DOCX.');
    }
    
    if (error.code === 'NETWORK_ERROR' || error.message.includes('Network Error')) {
      throw new Error('Lỗi kết nối mạng. Vui lòng kiểm tra kết nối internet và thử lại.');
    }
    
    // Lỗi mặc định
    throw new Error('Đã xảy ra lỗi khi phân tích hợp đồng. Vui lòng thử lại sau.');
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