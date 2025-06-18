import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1/analyze';

export async function analyzeContract(file) {
  try {
    const formData = new FormData();
    formData.append('file', file);
    
    console.log('Uploading file:', {
      name: file.name,
      type: file.type,
      size: `${(file.size / 1024 / 1024).toFixed(2)}MB`
    });

    const response = await axios.post(API_URL, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  } catch (error) {
    console.error('Error analyzing contract:', {
      message: error.message,
      response: error.response?.data,
      status: error.response?.status
    });
    throw error;
  }
}

export async function contractChat({ fileHash, contractText, question }) {
  try {
    const response = await axios.post(
      (import.meta.env.VITE_CHAT_API_URL || 'http://localhost:8080/api/v1/contract-chat'),
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