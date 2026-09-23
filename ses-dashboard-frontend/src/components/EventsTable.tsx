import React, { useState } from 'react';
import { ChevronLeft, ChevronRight, Mail, Copy, CheckCircle, XCircle, AlertTriangle, Eye, MousePointer, ChevronDown } from 'lucide-react';
import type { Event, MessageGroup, PaginationInfo } from '../types/api';
import { eventsService } from '../services/api';

interface EventsTableProps {
  events: MessageGroup[];
  pagination: PaginationInfo;
  onPageChange: (page: number) => void;
  loading?: boolean;
}

const formatDate = (dateString: string) => {
  if (!dateString) return '—';
  try {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat('id-ID', {
      year: 'numeric', month: 'short', day: 'numeric',
      hour: '2-digit', minute: '2-digit', second: '2-digit'
    }).format(date);
  } catch (e) {
    return dateString;
  }
};

const EventsTable = ({ events, pagination, onPageChange, loading }: EventsTableProps) => {
  const [expandedRows, setExpandedRows] = useState<Record<string, boolean>>({});
  const [copiedCode, setCopiedCode] = useState<string | null>(null);

  const toggleRow = (messageId: string) => {
    setExpandedRows(prev => ({
      ...prev,
      [messageId]: !prev[messageId]
    }));
  };

  const copyToClipboard = (text: string, e: React.MouseEvent) => {
    e.stopPropagation();
    navigator.clipboard.writeText(text);
    setCopiedCode(text);
    setTimeout(() => setCopiedCode(null), 2000);
  };

  if (loading) {
    return (
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-8">
        <div className="animate-pulse space-y-4">
          <div className="h-10 bg-slate-200 rounded w-full"></div>
          {[...Array(5)].map((_, i) => (
            <div key={i} className="h-16 bg-slate-100 rounded w-full"></div>
          ))}
        </div>
      </div>
    );
  }

  if (!events || events.length === 0) {
    return (
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-12 text-center">
        <Mail className="w-12 h-12 text-slate-300 mx-auto mb-4" />
        <h3 className="text-lg font-medium text-slate-900 mb-1">Tidak ada data</h3>
        <p className="text-slate-500">Belum ada log email yang terekam.</p>
      </div>
    );
  }

  return (
    <div className="overflow-x-auto rounded-xl border border-slate-200 bg-white shadow-sm">
      <table className="w-full text-left border-collapse">
        <thead>
          <tr className="bg-slate-50/70 border-b border-slate-200 text-xs font-semibold uppercase tracking-wider text-slate-500">
            <th className="py-4 px-6 w-10"></th>
            <th className="py-4 px-6">Message ID</th>
            <th className="py-4 px-6">Subject</th>
            <th className="py-4 px-6">Status Ringkasan</th>
            <th className="py-4 px-6">Waktu</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-100 text-sm text-slate-700">
          {events.map((event) => {
            const recipients = event.recipients_detail || [];
            return (
              <React.Fragment key={event.message_id}>
                {/* BARIS UTAMA - Putih Bersih & Berjarak */}
                <tr 
                  onClick={() => toggleRow(event.message_id)}
                  className="hover:bg-slate-50/80 transition-colors duration-150 cursor-pointer group"
                >
                  <td className="py-4 px-6 text-center">
                    <ChevronRight 
                      className={`w-4 h-4 text-slate-400 transition-transform duration-200 group-hover:text-slate-600 ${
                        expandedRows[event.message_id] ? 'rotate-90 text-slate-600' : ''
                      }`} 
                    />
                  </td>
                  <td className="py-4 px-6 font-mono text-xs text-slate-500 max-w-xs truncate" title={event.message_id}>
                    {event.message_id}
                  </td>
                  <td className="py-4 px-6 font-medium text-slate-900 max-w-sm truncate" title={event.subject || 'No subject'}>
                    {event.subject || 'No subject'}
                  </td>
                  <td className="py-4 px-6">
                    {/* Mini Badges Ringkasan */}
                    <span className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-50 text-blue-700 border border-blue-200">
                      {recipients.length || 1} Penerima
                    </span>
                  </td>
                  <td className="py-4 px-6 text-xs text-slate-400 whitespace-nowrap">
                    {formatDate(event.last_event_at)}
                  </td>
                </tr>

                {/* PANEL DETAIL - Hierarki Efek Masuk Ke Dalam (Nested) */}
                {expandedRows[event.message_id] && (
                  <tr className="bg-slate-50/60 transition-all">
                    <td colSpan={5} className="py-4 px-8 border-t border-b border-slate-100 shadow-inner">
                      <div className="rounded-lg border border-slate-200 bg-white overflow-hidden">
                        <table className="w-full text-xs">
                          <thead>
                            <tr className="bg-slate-50 border-b border-slate-200 text-slate-500 font-medium">
                              <th className="p-3">Email Penerima</th>
                              <th className="p-3">Tipe</th>
                              <th className="p-3">Status</th>
                              <th className="p-3">Alasan / Diagnostic Code SMTP</th>
                            </tr>
                          </thead>
                          <tbody className="divide-y divide-slate-100 text-slate-600">
                            {recipients.length > 0 ? recipients.map((rcpt, idx) => (
                              <tr key={idx} className="hover:bg-slate-50/40">
                                <td className="p-3 font-medium text-slate-800">{rcpt.email}</td>
                                <td className="p-3 text-slate-400 font-semibold">{rcpt.type || 'To'}</td>
                                <td className="p-3">
                                  {/* Soft Status Badges Interaktif */}
                                  <span className={`inline-flex px-2 py-0.5 rounded-md font-medium ${
                                    rcpt.status === 'Delivery' ? 'bg-emerald-50 text-emerald-700 border border-emerald-200' :
                                    (rcpt.status === 'Bounce' || rcpt.status === 'FAILED') ? 'bg-rose-50 text-rose-700 border border-rose-200' :
                                    rcpt.status === 'Complaint' ? 'bg-amber-50 text-amber-700 border border-amber-200' :
                                    'bg-blue-50 text-blue-700 border border-blue-200'
                                  }`}>
                                    {rcpt.status}
                                  </span>
                                </td>
                                <td className="p-3">
                                  {rcpt.diagnostic_code && rcpt.diagnostic_code !== '-' ? (
                                    <div className="flex items-center gap-2 max-w-md bg-slate-50 p-2 rounded border border-slate-200 font-mono text-[11px] text-slate-600">
                                      <span className="truncate flex-1" title={rcpt.diagnostic_code}>
                                        {rcpt.diagnostic_code}
                                      </span>
                                      <button 
                                        onClick={(e) => copyToClipboard(rcpt.diagnostic_code, e)}
                                        className="p-1 hover:bg-white rounded border border-slate-200 transition-colors text-slate-400 hover:text-slate-600 relative"
                                        title="Copy to clipboard"
                                      >
                                        {copiedCode === rcpt.diagnostic_code ? (
                                          <CheckCircle className="w-3 h-3 text-emerald-500" />
                                        ) : (
                                          <Copy className="w-3 h-3" />
                                        )}
                                      </button>
                                    </div>
                                  ) : (
                                    <span className="text-slate-300">—</span>
                                  )}
                                </td>
                              </tr>
                            )) : (
                              <tr>
                                <td colSpan={4} className="p-3 text-center text-slate-400 italic">
                                  Tidak ada detail (Data legacy: {event.email})
                                </td>
                              </tr>
                            )}
                          </tbody>
                        </table>
                      </div>
                    </td>
                  </tr>
                )}
              </React.Fragment>
            );
          })}
        </tbody>
      </table>
      
      {/* Pagination Footer */}
      <div className="px-6 py-4 border-t border-slate-200 bg-white flex items-center justify-between">
        <div className="text-sm text-slate-500">
          Menampilkan <span className="font-medium text-slate-700">{pagination.total === 0 ? 0 : ((pagination.page - 1) * pagination.limit) + 1}</span> hingga <span className="font-medium text-slate-700">{Math.min(pagination.page * pagination.limit, pagination.total)}</span> dari <span className="font-medium text-slate-700">{pagination.total}</span> hasil
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={() => onPageChange(pagination.page - 1)}
            disabled={!pagination.hasPrev}
            className="p-2 border border-slate-200 rounded-md text-slate-500 hover:bg-slate-50 hover:text-slate-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            <ChevronLeft className="w-4 h-4" />
          </button>
          <span className="px-3 py-1.5 text-sm text-slate-700 bg-slate-50 border border-slate-200 rounded-md font-medium">
            Hal {pagination.page} / {pagination.totalPages}
          </span>
          <button
            onClick={() => onPageChange(pagination.page + 1)}
            disabled={!pagination.hasNext}
            className="p-2 border border-slate-200 rounded-md text-slate-500 hover:bg-slate-50 hover:text-slate-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            <ChevronRight className="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>
  );
};

export default EventsTable;
