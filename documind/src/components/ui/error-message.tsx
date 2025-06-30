import React from 'react';
import { AlertTriangle } from 'lucide-react';

interface ErrorMessageProps {
  title?: string;
  message: string;
  className?: string;
}

const ErrorMessage: React.FC<ErrorMessageProps> = ({ 
  title = "Lỗi", 
  message, 
  className = "" 
}) => {
  const isQuotaError = message.includes('quota') || message.includes('429') || message.includes('exceeded');
  
  return (
    <div className={`p-4 bg-red-50 border border-red-200 rounded-lg ${className}`}>
      <div className="flex items-start gap-3">
        <AlertTriangle className="w-4 h-4 text-red-600 flex-shrink-0 mt-0.5" />
        <div className="flex-1">
          <h4 className="text-sm font-medium text-red-900 mb-1">{title}</h4>
          <p className="text-sm text-red-700">{message}</p>
          {isQuotaError && (
            <div className="mt-2 text-xs text-red-600">
              💡 <strong>Gợi ý:</strong> Thử lại sau vài phút hoặc liên hệ admin để nâng cấp quota API.
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default ErrorMessage; 