import { Tab } from '@headlessui/react';
import { ExclamationTriangleIcon } from '@heroicons/react/24/outline';

function classNames(...classes) {
  return classes.filter(Boolean).join(' ');
}

export default function ResultsDisplay({ analysisResult }) {
  if (!analysisResult) return null;

  const { summary, key_clauses = [], potential_risks = [] } = analysisResult;

  return (
    <div className="mt-8">
      <Tab.Group>
        <Tab.List className="flex space-x-2 rounded-xl bg-slate-200 p-1 mb-4">
          <Tab
            className={({ selected }) =>
              classNames(
                'w-full py-2.5 text-sm leading-5 font-medium rounded-lg',
                selected
                  ? 'bg-white shadow text-blue-700'
                  : 'text-slate-700 hover:bg-white/[0.60]'
              )
            }
          >
            Tóm tắt
          </Tab>
          <Tab
            className={({ selected }) =>
              classNames(
                'w-full py-2.5 text-sm leading-5 font-medium rounded-lg',
                selected
                  ? 'bg-white shadow text-blue-700'
                  : 'text-slate-700 hover:bg-white/[0.60]'
              )
            }
          >
            Điều khoản chính
          </Tab>
          <Tab
            className={({ selected }) =>
              classNames(
                'w-full py-2.5 text-sm leading-5 font-medium rounded-lg',
                selected
                  ? 'bg-white shadow text-blue-700'
                  : 'text-slate-700 hover:bg-white/[0.60]'
              )
            }
          >
            Rủi ro tiềm ẩn
          </Tab>
        </Tab.List>
        <Tab.Panels>
          <Tab.Panel>
            <div className="bg-slate-50 p-6 rounded-lg">
              <h2 className="font-semibold text-xl text-slate-800 mb-3">Tóm tắt chung</h2>
              <p className="text-slate-700 leading-relaxed">{summary}</p>
            </div>
          </Tab.Panel>
          <Tab.Panel>
            <div className="mb-6">
              <h2 className="font-semibold text-xl text-slate-800 mb-3">Các điều khoản chính</h2>
              <ul className="space-y-3">
                {key_clauses.map((clause, idx) => (
                  <li key={idx} className="text-slate-700 bg-white p-4 rounded-lg shadow-sm">
                    {clause}
                  </li>
                ))}
              </ul>
            </div>
          </Tab.Panel>
          <Tab.Panel>
            <div className="bg-red-50 border-l-4 border-red-400 p-6 rounded-lg">
              <h2 className="font-semibold text-xl text-red-800 mb-3">Các rủi ro tiềm ẩn</h2>
              <ul className="space-y-3">
                {potential_risks.map((risk, idx) => (
                  <li key={idx} className="text-red-900 flex items-start bg-white/50 p-4 rounded-lg">
                    <ExclamationTriangleIcon className="w-5 h-5 mr-3 text-red-500 flex-shrink-0 mt-1" />
                    <span>{risk}</span>
                  </li>
                ))}
              </ul>
            </div>
          </Tab.Panel>
        </Tab.Panels>
      </Tab.Group>
    </div>
  );
}