import React from 'react';

export default function ResultsDisplay({ analysisResult }) {
  if (!analysisResult) return null;
  const { summary, keyClauses, potentialRisks } = analysisResult;
  return (
    <div className="space-y-6 mt-6">
      {/* Tóm tắt chung */}
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-lg font-semibold mb-2">Tóm tắt chung</h2>
        <p className="text-gray-800 whitespace-pre-line">{summary || 'Không có dữ liệu.'}</p>
      </div>
      {/* Điều khoản chính */}
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-lg font-semibold mb-2">Các điều khoản chính</h2>
        {keyClauses && keyClauses.length > 0 ? (
          <ol className="list-decimal list-inside space-y-1">
            {keyClauses.map((clause, idx) => (
              <li key={idx} className="text-gray-800">{clause}</li>
            ))}
          </ol>
        ) : <p className="text-gray-500">Không có dữ liệu.</p>}
      </div>
      {/* Rủi ro tiềm ẩn */}
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-lg font-semibold mb-2">Rủi ro tiềm ẩn</h2>
        {potentialRisks && potentialRisks.length > 0 ? (
          <ul className="list-disc list-inside space-y-1">
            {potentialRisks.map((risk, idx) => (
              <li key={idx} className="text-red-600 flex items-center"><span className="mr-2">⚠️</span>{risk}</li>
            ))}
          </ul>
        ) : <p className="text-gray-500">Không có dữ liệu.</p>}
      </div>
    </div>
  );
} 