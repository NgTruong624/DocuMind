import axios from 'axios';

// 1. Định nghĩa một API client với base URL duy nhất
const apiClient = axios.create({
  baseURL: (import.meta as any).env.VITE_API_URL || 'http://localhost:8090/api/v1',
});

// Định nghĩa các kiểu dữ liệu trả về (giữ nguyên)
export interface AnalyzeResult {
  summary: string;
  key_clauses?: string[];
  potential_risks?: string[];
  file_hash?: string;
  [key: string]: any;
}

export interface ContractChatResponse {
  answer: string;
  [key: string]: any;
}

// 2. Sửa lại hàm analyzeContract để dùng apiClient
export async function analyzeContract(file: File): Promise<AnalyzeResult> {
  const formData = new FormData();
  formData.append('file', file);

  try {
    // Gọi đến endpoint "/analyze" (nó sẽ tự động nối với baseURL)
    const response = await apiClient.post<AnalyzeResult>('/analyze', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  } catch (error: any) {
    // Xử lý lỗi (giữ nguyên logic cũ)
    console.error('Error analyzing contract:', error);
    const defaultMessage = 'Đã xảy ra lỗi khi phân tích hợp đồng.';
    throw new Error(error?.response?.data?.message || defaultMessage);
  }
}

// 3. Sửa lại hàm contractChat để dùng apiClient
export interface ContractChatParams {
  fileHash?: string;
  contractText?: string;
  question: string;
}

export async function contractChat({ fileHash, contractText, question }: ContractChatParams): Promise<ContractChatResponse> {
  try {
    // Gọi đến endpoint "/contract-chat"
    const response = await apiClient.post<ContractChatResponse>('/contract-chat', {
      file_hash: fileHash || '',
      contract_text: contractText || '',
      question,
    });
    return response.data;
  } catch (error: any) {
    console.error('Error in contract chat:', error);
    const defaultMessage = 'Đã xảy ra lỗi khi gửi câu hỏi.';
    throw new Error(error?.response?.data?.message || defaultMessage);
  }
}