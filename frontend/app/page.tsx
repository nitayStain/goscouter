'use client';

import { useState } from 'react';
import { Scanner } from '@/components/Scanner';
import { Results } from '@/components/Results';

export default function Home() {
  const [results, setResults] = useState<any>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleScan = async (domain: string) => {
    setLoading(true);
    setError(null);
    setResults(null);

    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/subdomains?domain=${encodeURIComponent(domain)}`
      );

      if (!response.ok) {
        throw new Error(`Failed to scan: ${response.statusText}`);
      }

      const data = await response.json();
      setResults(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-[#1a1625]">
      <div className="container mx-auto px-4 py-16">
        <div className="max-w-5xl mx-auto">
          {/* Header */}
          <div className="mb-16">
            <h1 className="text-4xl font-bold text-white mb-3">
              GoScouter
            </h1>
            <p className="text-slate-400 text-sm">
              Discover subdomains through Certificate Transparency logs
            </p>
          </div>

          {/* Scanner */}
          <Scanner onScan={handleScan} loading={loading} />

          {/* Error Display */}
          {error && (
            <div className="mt-6 p-4 bg-red-950/50 border border-red-800/50 rounded-md">
              <p className="text-red-300 text-sm">{error}</p>
            </div>
          )}

          {/* Results */}
          {results && <Results data={results} />}
        </div>
      </div>
    </div>
  );
}
