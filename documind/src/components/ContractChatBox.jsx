import React, { useState, useRef, useEffect } from 'react';
import { contractChat } from '../api/analysisService';

export default function ContractChatBox({ fileHash, contractText, disabled, resetTrigger }) {
  const [messages, setMessages] = useState([]); // {role: 'user'|'ai', content: string}
  const [input, setInput] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const messagesEndRef = useRef(null);

  // Reset chat khi file mới được upload
  useEffect(() => {
    setMessages([]);
    setInput('');
    setError(null);
    setLoading(false);
  }, [resetTrigger, fileHash, contractText]);

  // Scroll xuống cuối khi có tin nhắn mới
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const handleSend = async () => {
    if (!input.trim()) return;
    setMessages((prev) => [...prev, { role: 'user', content: input }]);
    setLoading(true);
    setError(null);
    const question = input;
    setInput('');
    try {
      // Gọi API thực tế
      const res = await contractChat({ fileHash, contractText, question });
      const aiAnswer = res.answer || 'Không nhận được câu trả lời từ AI.';
      setMessages((prev) => [...prev, { role: 'ai', content: aiAnswer }]);
    } catch (err) {
      let msg = 'Đã xảy ra lỗi khi gửi câu hỏi.';
      if (err.response?.data?.error) msg = err.response.data.error;
      setError(msg);
    } finally {
      setLoading(false);
    }
  };

  const handleInputKeyDown = (e) => {
    if (e.key === 'Enter' && !e.shiftKey && !loading && input.trim()) {
      e.preventDefault();
      handleSend();
    }
  };

  return (
    <div className="mt-8 bg-slate-50 rounded-xl shadow-inner p-6">
      <h2 className="font-semibold text-lg text-slate-800 mb-4">Hỏi đáp về hợp đồng</h2>
      <div className="h-64 overflow-y-auto bg-white rounded-lg border border-slate-200 p-4 mb-4 flex flex-col gap-3">
        {messages.length === 0 && (
          <div className="text-slate-400 text-center mt-12">Hãy đặt câu hỏi về hợp đồng, điều khoản, rủi ro...</div>
        )}
        {messages.map((msg, idx) => (
          <div key={idx} className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}>
            <div
              className={`max-w-[75%] px-4 py-2 rounded-2xl text-sm whitespace-pre-line shadow-sm
                ${msg.role === 'user' ? 'bg-blue-500 text-white rounded-br-md' : 'bg-slate-200 text-slate-800 rounded-bl-md'}`}
            >
              {msg.content}
            </div>
          </div>
        ))}
        <div ref={messagesEndRef} />
      </div>
      {error && <div className="text-red-600 mb-2 text-sm">{error}</div>}
      <div className="flex gap-2 items-end">
        <textarea
          className="flex-1 border border-slate-300 rounded-lg p-2 resize-none min-h-[40px] max-h-24 focus:outline-blue-400"
          placeholder="Nhập câu hỏi về hợp đồng..."
          value={input}
          onChange={e => setInput(e.target.value)}
          onKeyDown={handleInputKeyDown}
          disabled={loading || disabled}
          rows={1}
        />
        <button
          className="bg-blue-500 hover:bg-blue-600 text-white font-semibold px-4 py-2 rounded-lg disabled:opacity-50"
          onClick={handleSend}
          disabled={loading || !input.trim() || disabled}
        >
          Gửi
        </button>
      </div>
      {loading && (
        <div className="flex items-center gap-2 mt-2 text-blue-500 text-sm">
          <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24"><circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" /><path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z" /></svg>
          Đang lấy câu trả lời từ AI...
        </div>
      )}
    </div>
  );
} 