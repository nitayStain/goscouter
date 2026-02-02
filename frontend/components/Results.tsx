'use client';

import { Globe, Shield, Server, Award } from 'lucide-react';

interface SubdomainItem {
  name: string;
  ips: string[];
  ip_owner: string;
  cert_issuer: string;
  cert_expiry: string;
}

interface ResultsData {
  domain: string;
  has_wildcard: boolean;
  count: number;
  items: SubdomainItem[];
}

interface ResultsProps {
  data: ResultsData;
}

export function Results({ data }: ResultsProps) {
  // Calculate statistics
  const uniqueIPs = new Set(
    data.items.flatMap((item) => item.ips)
  ).size;

  const uniqueIssuers = new Set(
    data.items.map((item) => item.cert_issuer).filter(Boolean)
  ).size;

  return (
    <div className="mt-6 space-y-6">
      {/* Statistics Grid */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
        <StatCard
          icon={<Globe className="w-4 h-4" />}
          label="Subdomains"
          value={data.count.toString()}
        />
        <StatCard
          icon={<Shield className="w-4 h-4" />}
          label="Wildcard"
          value={data.has_wildcard ? 'Yes' : 'No'}
        />
        <StatCard
          icon={<Server className="w-4 h-4" />}
          label="Unique IPs"
          value={uniqueIPs.toString()}
        />
        <StatCard
          icon={<Award className="w-4 h-4" />}
          label="Issuers"
          value={uniqueIssuers.toString()}
        />
      </div>

      {/* Subdomains List */}
      <div className="bg-[#231d35] rounded-lg border border-purple-800/30 p-6">
        <h2 className="text-lg font-semibold text-white mb-4">
          Discovered Subdomains
        </h2>

        {data.items.length === 0 ? (
          <p className="text-slate-400 text-center py-8 text-sm">
            No subdomains found for this domain.
          </p>
        ) : (
          <div className="space-y-2 max-h-[600px] overflow-y-auto pr-2 custom-scrollbar">
            {data.items.map((item, index) => (
              <SubdomainCard key={index} item={item} />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

interface StatCardProps {
  icon: React.ReactNode;
  label: string;
  value: string;
}

function StatCard({ icon, label, value }: StatCardProps) {
  return (
    <div className="bg-[#231d35] rounded-lg border border-purple-800/30 p-4">
      <div className="flex items-center gap-2 mb-2 text-purple-400">
        {icon}
        <span className="text-xs font-medium text-slate-400">{label}</span>
      </div>
      <p className="text-2xl font-semibold text-white">{value}</p>
    </div>
  );
}

interface SubdomainCardProps {
  item: SubdomainItem;
}

function SubdomainCard({ item }: SubdomainCardProps) {
  return (
    <div className="bg-[#1a1625] hover:bg-[#201c2e] rounded-md p-4 border border-purple-800/30 hover:border-purple-700/50 transition-colors">
      <h3 className="text-sm font-mono font-medium text-purple-400 mb-3 break-all">
        {item.name}
      </h3>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-2 text-xs">
        {item.ips && item.ips.length > 0 && (
          <InfoItem label="IPs" value={item.ips.join(', ')} />
        )}
        {item.ip_owner && (
          <InfoItem label="Owner" value={item.ip_owner} />
        )}
        {item.cert_issuer && (
          <InfoItem label="Issuer" value={item.cert_issuer} />
        )}
        {item.cert_expiry && (
          <InfoItem label="Expiry" value={item.cert_expiry} />
        )}
      </div>
    </div>
  );
}

interface InfoItemProps {
  label: string;
  value: string;
}

function InfoItem({ label, value }: InfoItemProps) {
  return (
    <div>
      <span className="text-slate-500">{label}:</span>{' '}
      <span className="text-slate-300">{value}</span>
    </div>
  );
}
