'use client';

import { useState, FormEvent } from 'react';
import { Search, Loader2 } from 'lucide-react';

interface ScannerProps {
  onScan: (domain: string) => void;
  loading: boolean;
}

export function Scanner({ onScan, loading }: ScannerProps) {
  const [domain, setDomain] = useState('');

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    if (domain.trim()) {
      onScan(domain.trim());
    }
  };

  return (
    <div className="bg-[#231d35] rounded-lg border border-purple-800/30 p-6">
      <form onSubmit={handleSubmit} className="flex gap-3">
        <div className="flex-1 relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 w-4 h-4" />
          <input
            type="text"
            value={domain}
            onChange={(e) => setDomain(e.target.value)}
            placeholder="Enter domain (e.g., example.com)"
            disabled={loading}
            className="w-full pl-10 pr-4 py-3 bg-[#1a1625] text-white rounded-md border border-purple-800/30 focus:border-purple-600 focus:outline-none focus:ring-2 focus:ring-purple-600/20 disabled:opacity-50 disabled:cursor-not-allowed transition-colors placeholder:text-slate-500 text-sm"
          />
        </div>
        <button
          type="submit"
          disabled={loading || !domain.trim()}
          className="px-6 py-3 bg-purple-600 text-white font-medium rounded-md hover:bg-purple-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex items-center gap-2 text-sm"
        >
          {loading ? (
            <>
              <Loader2 className="w-4 h-4 animate-spin" />
              <span>Scanning...</span>
            </>
          ) : (
            <span>Scan</span>
          )}
        </button>
      </form>
    </div>
  );
}
