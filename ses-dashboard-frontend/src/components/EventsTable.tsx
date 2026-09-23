import React, { useState } from 'react';
import { ChevronLeft, ChevronRight, Mail, Copy, CheckCircle, XCircle, AlertTriangle, Eye, MousePointer, ChevronDown, Send } from 'lucide-react';
import type { Event, MessageGroup, PaginationInfo } from '../types/api';
import { eventsService } from '../services/api';

interface EventsTableProps {
  events: MessageGroup[];
  pagination: PaginationInfo;
  onPageChange: (page: number) => void;
  loading?: boolean;
}

const getEventIcon = (eventType: string) => {
  const type = eventType?.toLowerCase() || '';
  switch (type) {
    case 'send': return <Send className="w-4 h-4 text-blue-500" />;
    case 'delivery': return <CheckCircle className="w-4 h-4 text-emerald-500" />;
    case 'bounce': return <XCircle className="w-4 h-4 text-rose-500" />;
    case 'complaint': return <AlertTriangle className="w-4 h-4 text-amber-500" />;
    case 'open': return <Eye className="w-4 h-4 text-purple-500" />;
    case 'click': return <MousePointer className="w-4 h-4 text-indigo-500" />;
    default: return <Mail className="w-4 h-4 text-gray-500" />;
  }
};

const getEventTypeColor = (eventType: string) => {
  const type = eventType?.toLowerCase() || '';
  switch (type) {
    case 'send': return 'bg-blue-100 text-blue-800';
    case 'delivery': return 'bg-emerald-100 text-emerald-800';
    case 'bounce': return 'bg-rose-100 text-rose-800';
    case 'complaint': return 'bg-amber-100 text-amber-800';
    case 'open': return 'bg-purple-100 text-purple-800';
    case 'click': return 'bg-indigo-100 text-indigo-800';
    default: return 'bg-gray-100 text-gray-800';
  }
};

const getStatusColor = (status: string) => {
  const s = status?.toLowerCase() || '';
  if (s.includes('fail') || s.includes('bounce') || s.includes('error')) return 'bg-rose-50 text-rose-700 border border-rose-200';
  if (s.includes('success') || s.includes('deliver')) return 'bg-emerald-50 text-emerald-700 border border-emerald-200';
  if (s.includes('complaint')) return 'bg-amber-50 text-amber-700 border border-amber-200';
  if (s.includes('pending')) return 'bg-blue-50 text-blue-700 border border-blue-200';
  return 'bg-gray-50 text-gray-700 border border-gray-200';
};

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

  // Timeline Modal states
  const [selectedMessageId, setSelectedMessageId] = useState<string | null>(null);
  const [detailEvents, setDetailEvents] = useState<Event[]>([]);
  const [detailLoading, setDetailLoading] = useState(false);
  const [detailError, setDetailError] = useState<string | null>(null);

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

  const openDetail = async (messageId: string, e: React.MouseEvent) => {
    e.stopPropagation();
    setSelectedMessageId(messageId);
    setDetailLoading(true);
    setDetailError(null);
    try {
      const data = await eventsService.getEventDetail(messageId);
      setDetailEvents(data.events || []);
    } catch (err) {
      setDetailError(err instanceof Error ? err.message : 'Failed to load timeline');
    } finally {
      setDetailLoading(false);
    }
  };

  const closeDetail = () => {
    setSelectedMessageId(null);
    setDetailEvents([]);
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

  const currentGroup = events.find(g => g.message_id === selectedMessageId);

  return (
    <>
      <div className="overflow-x-auto rounded-xl border border-slate-200 bg-white shadow-sm">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="bg-slate-50/70 border-b border-slate-200 text-xs font-semibold uppercase tracking-wider text-slate-500">
              <th className="py-4 px-6 w-10"></th>
              <th className="py-4 px-6">Message ID</th>
              <th className="py-4 px-6">Subject</th>
              <th className="py-4 px-6">Status Ringkasan</th>
              <th className="py-4 px-6">Waktu</th>
              <th className="py-4 px-6 text-right">Aksi</th>
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
                    <td className="py-4 px-6 text-right whitespace-nowrap">
                      <button
                        onClick={(e) => openDetail(event.message_id, e)}
                        className="inline-flex items-center p-1.5 text-slate-400 hover:text-blue-600 hover:bg-blue-50 rounded transition-colors"
                        title="Lihat Timeline"
                      >
                        <Eye className="w-5 h-5" />
                      </button>
                    </td>
                  </tr>

                  {/* PANEL DETAIL - Hierarki Efek Masuk Ke Dalam (Nested) */}
                  {expandedRows[event.message_id] && (
                    <tr className="bg-slate-50/60 transition-all">
                      <td colSpan={6} className="py-4 px-8 border-t border-b border-slate-100 shadow-inner">
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

      {/* MODAL EVENT TIMELINE (Dikembalikan) */}
      {selectedMessageId && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg w-full max-w-5xl max-h-[90vh] overflow-hidden shadow-xl">
            <div className="px-6 py-4 border-b border-gray-200 flex items-start justify-between">
              <div>
                <h3 className="text-lg font-semibold text-gray-900">Event Timeline</h3>
                <p className="text-sm text-gray-500 mt-1 break-all">{selectedMessageId}</p>
                {currentGroup && (
                  <p className="text-sm text-gray-600 mt-1">{currentGroup.email} • {currentGroup.subject || 'No subject'}</p>
                )}
              </div>
              <button onClick={closeDetail} className="text-gray-400 hover:text-gray-600 text-xl font-bold">×</button>
            </div>

            <div className="p-6 overflow-y-auto max-h-[calc(90vh-88px)]">
              {detailLoading ? (
                <div className="animate-pulse space-y-3">
                  {[...Array(4)].map((_, i) => (
                    <div key={i} className="h-16 bg-gray-100 rounded-lg"></div>
                  ))}
                </div>
              ) : detailError ? (
                <div className="bg-red-50 text-red-700 border border-red-200 rounded-lg p-4">{detailError}</div>
              ) : detailEvents.length === 0 ? (
                <div className="text-center text-gray-500 py-8">No event timeline found.</div>
              ) : (
                <div className="space-y-4">
                  {detailEvents.map((event, index) => (
                    <div key={`${event.MessageID}-${event.EventType}-${index}`} className="border border-gray-200 rounded-lg p-4">
                      <div className="flex flex-col gap-3 lg:flex-row lg:items-start lg:justify-between">
                        <div>
                          <div className="flex items-center gap-2 flex-wrap">
                            {getEventIcon(event.EventType)}
                            <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getEventTypeColor(event.EventType)}`}>
                              {event.EventType}
                            </span>
                            <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusColor(event.Status)}`}>
                              {event.Status}
                            </span>
                          </div>
                          <p className="text-sm text-gray-900 mt-3 break-all">{event.Email}</p>
                          <p className="text-sm text-gray-500 mt-1">{event.Source}</p>
                        </div>
                        <div className="text-sm text-gray-500">{formatDate(event.EventTimestamp)}</div>
                      </div>

                      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4 text-sm">
                        <div>
                          <p className="text-xs uppercase tracking-wider text-gray-500 mb-1">Reason</p>
                          <p className="text-gray-900 break-words">{event.Reason || '—'}</p>
                        </div>
                        <div>
                          <p className="text-xs uppercase tracking-wider text-gray-500 mb-1">Bounce Type</p>
                          <p className="text-gray-900">{event.BounceType || '—'}</p>
                        </div>
                        <div>
                          <p className="text-xs uppercase tracking-wider text-gray-500 mb-1">Bounce Sub Type</p>
                          <p className="text-gray-900">{event.BounceSubType || '—'}</p>
                        </div>
                        <div>
                          <p className="text-xs uppercase tracking-wider text-gray-500 mb-1">Reporting MTA</p>
                          <p className="text-gray-900 break-words">{event.ReportingMTA || '—'}</p>
                        </div>
                        <div className="md:col-span-2">
                          <p className="text-xs uppercase tracking-wider text-gray-500 mb-1">Diagnostic Code</p>
                          <p className="text-gray-900 break-words">{event.DiagnosticCode || '—'}</p>
                        </div>
                        <div className="md:col-span-2">
                          <p className="text-xs uppercase tracking-wider text-gray-500 mb-1">SMTP Response</p>
                          <p className="text-gray-900 break-words">{event.SmtpResponse || '—'}</p>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </>
  );
};

export default EventsTable;
