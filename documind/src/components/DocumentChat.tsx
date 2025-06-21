import React, { useState, useRef, useEffect } from 'react';
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Send, Bot, User, FileText } from "lucide-react";
import { contractChat } from '@/api/analysisService';

interface Message {
  id: string;
  content: string;
  isUser: boolean;
  timestamp: Date;
}

interface DocumentChatProps {
  documentSummary: string;
  keyClauses?: string[];
  potentialRisks?: string[];
  fileHash?: string;
  onClose?: () => void;
}

const DocumentChat: React.FC<DocumentChatProps> = ({ 
  documentSummary, 
  keyClauses = [], 
  potentialRisks = [],
  fileHash,
  onClose 
}) => {
  const [messages, setMessages] = useState<Message[]>([
    {
      id: '1',
      content: `Xin chào! Tôi đã phân tích hợp đồng của bạn. Dưới đây là tóm tắt nhanh:\n\n**Tóm tắt:** ${documentSummary}\n**Số điều khoản chính:** ${keyClauses.length}\n**Số rủi ro tiềm ẩn:** ${potentialRisks.length}\n\nBạn có thể hỏi tôi bất kỳ câu hỏi nào về hợp đồng của bạn, bao gồm các điều khoản cụ thể, rủi ro hoặc các thuật ngữ bạn muốn tôi giải thích.`,
      isUser: false,
      timestamp: new Date()
    }
  ]);
  const [inputValue, setInputValue] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const handleSendMessage = async () => {
    if (!inputValue.trim()) return;
    setMessages(prev => [
      ...prev,
      {
        id: (prev.length + 1).toString(),
        content: inputValue,
        isUser: true,
        timestamp: new Date()
      }
    ]);
    setIsLoading(true);
    setError(null);
    const question = inputValue;
    setInputValue('');
    try {
      const res = await contractChat({ fileHash, contractText: documentSummary, question });
      setMessages(prev => [
        ...prev,
        {
          id: (prev.length + 2).toString(),
          content: res.answer || 'Không có câu trả lời từ AI.',
          isUser: false,
          timestamp: new Date()
        }
      ]);
    } catch (err: any) {
      setError(err?.response?.data?.message || 'Đã xảy ra lỗi khi gửi câu hỏi.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey && !isLoading && inputValue.trim()) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  return (
    <Card className="w-full max-w-4xl mx-auto mt-8 shadow-lg">
      <CardHeader className="bg-gradient-to-r from-blue-50 to-indigo-50 border-b">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-gradient-to-r from-blue-500 to-indigo-500 rounded-full flex items-center justify-center">
              <FileText className="w-5 h-5 text-white" />
            </div>
            <div>
              <CardTitle className="text-xl">Trợ lý hợp đồng</CardTitle>
              <p className="text-sm text-slate-600">Đặt câu hỏi về hợp đồng đã được phân tích của bạn</p>
            </div>
          </div>
          {onClose && (
            <Button variant="ghost" size="sm" onClick={onClose}>
              ×
            </Button>
          )}
        </div>
      </CardHeader>
      
      <CardContent className="p-0">
        <div className="h-96 overflow-y-auto p-4 space-y-4">
          {messages.map((message) => (
            <div
              key={message.id}
              className={`flex gap-3 ${message.isUser ? 'justify-end' : 'justify-start'}`}
            >
              {!message.isUser && (
                <div className="w-8 h-8 bg-gradient-to-r from-blue-100 to-indigo-100 rounded-full flex items-center justify-center flex-shrink-0">
                  <Bot className="w-4 h-4 text-blue-600" />
                </div>
              )}
              
              <div
                className={`max-w-[80%] p-3 rounded-lg ${
                  message.isUser
                    ? 'bg-blue-600 text-white ml-auto'
                    : 'bg-slate-100 text-slate-900'
                }`}
              >
                <div className="whitespace-pre-wrap">{message.content}</div>
                <div
                  className={`text-xs mt-1 ${
                    message.isUser ? 'text-blue-100' : 'text-slate-500'
                  }`}
                >
                  {message.timestamp.toLocaleTimeString([], { 
                    hour: '2-digit', 
                    minute: '2-digit' 
                  })}
                </div>
              </div>
              
              {message.isUser && (
                <div className="w-8 h-8 bg-gradient-to-r from-slate-100 to-slate-200 rounded-full flex items-center justify-center flex-shrink-0">
                  <User className="w-4 h-4 text-slate-600" />
                </div>
              )}
            </div>
          ))}
          
          {isLoading && (
            <div className="flex gap-3 justify-start">
              <div className="w-8 h-8 bg-gradient-to-r from-blue-100 to-indigo-100 rounded-full flex items-center justify-center">
                <Bot className="w-4 h-4 text-blue-600" />
              </div>
              <div className="bg-slate-100 p-3 rounded-lg">
                <div className="flex space-x-1">
                  <div className="w-2 h-2 bg-slate-400 rounded-full animate-bounce"></div>
                  <div className="w-2 h-2 bg-slate-400 rounded-full animate-bounce" style={{ animationDelay: '0.1s' }}></div>
                  <div className="w-2 h-2 bg-slate-400 rounded-full animate-bounce" style={{ animationDelay: '0.2s' }}></div>
                </div>
              </div>
            </div>
          )}
          
          <div ref={messagesEndRef} />
        </div>
        
        {error && <div className="text-red-600 mb-2 text-sm">{error}</div>}
        
        <div className="p-4 border-t bg-white">
          <div className="flex gap-2">
            <Input
              value={inputValue}
              onChange={(e) => setInputValue(e.target.value)}
              onKeyPress={handleKeyPress}
              placeholder="Hỏi về rủi ro, điều khoản, hoặc bất cứ điều gì khác..."
              className="flex-1"
              disabled={isLoading}
            />
            <Button 
              onClick={handleSendMessage} 
              disabled={!inputValue.trim() || isLoading}
              className="bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700"
            >
              <Send className="w-4 h-4" />
            </Button>
          </div>
          <p className="text-xs text-slate-500 mt-2">
            Nhấn Enter để gửi • Trợ lý AI này dựa trên phân tích hợp đồng của bạn
          </p>
        </div>
      </CardContent>
    </Card>
  );
};

export default DocumentChat;
