import { useState } from 'react';
import { ChevronLeft, ChevronRight, Mail, Send, CheckCircle, XCircle, AlertTriangle, Eye, MousePointer, ChevronDown, ChevronRight as ChevronRightIcon, Copy } from 'lucide-react';
import type { Event, MessageGroup, PaginationInfo, RecipientDetail } from '../types/api';
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
    return new Intl.DateTimeFormat('en-US', {
      year: 'numeric', month: 'short', day: 'numeric',
      hour: '2-digit', minute: '2-digit', second: '2-digit'
    }).format(date);
  } catch (e) {
    return dateString;
  }
};

const EventsTable = ({ events, pagination, onPageChange, loading }: EventsTableProps) => {
  const [selectedMessageId, setSelectedMessageId] = useState<string | null>(null);
  const [detailEvents, setDetailEvents] = useState<Event[]>([]);
  const [detailLoading, setDetailLoading] = useState(false);
  const [detailError, setDetailError] = useState<string | null>(null);
  
  const [expandedRows, setExpandedRows] = useState<Set<string>>(new Set());

  const toggleRow = (messageId: string) => {
    const newExpanded = new Set(expandedRows);
    if (newExpanded.has(messageId)) {
      newExpanded.delete(messageId);
    } else {
      newExpanded.add(messageId);
    }
    setExpandedRows(newExpanded);
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
  
  const copyToClipboard = (text: string, e: React.MouseEvent) => {
    e.stopPropagation();
    navigator.clipboard.writeText(text);
  };

  if (loading) {
    return (
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-8">
        <div className="animate-pulse space-y-4">
          <div className="h-10 bg-gray-200 rounded w-full"></div>
          {[...Array(5)].map((_, i) => (
            <div key={i} className="h-16 bg-gray-100 rounded w-full"></div>
          ))}
        </div>
      </div>
    );
  }

  if (!events || events.length === 0) {
    return (
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-12 text-center">
        <Mail className="w-12 h-12 text-gray-300 mx-auto mb-4" />
        <h3 className="text-lg font-medium text-gray-900 mb-1">No events found</h3>
        <p className="text-gray-500">There are no email events matching your current filters.</p>
      </div>
    );
  }

  const currentGroup = events.find(g => g.message_id === selectedMessageId);

  return (
    <>
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="w-10 px-4 py-3"></th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Email</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Subject & Source</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Latest Status</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Last Event</th>
                <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {events.map((event) => {
                const isExpanded = expandedRows.has(event.message_id);
                const recipients = event.recipients_detail || [];
                const hasIssues = recipients.some(r => r.status === 'Bounce' || r.status === 'Complaint');
                
                return (
                  <React.Fragment key={event.message_id}>
                    <tr 
                      className="hover:bg-slate-50 transition-colors cursor-pointer group"
                      onClick={() => toggleRow(event.message_id)}
                    >
                      <td className="px-4 py-4 text-gray-400">
                        <ChevronRightIcon className={`w-5 h-5 transform transition-transform duration-200 ${isExpanded ? 'rotate-90' : ''}`} />
                      </td>
                      <td className="px-4 py-4 text-sm text-gray-900">
                        <div className="font-medium truncate max-w-xs" title={event.email}>{event.email}</div>
                        {recipients.length > 1 && (
                          <div className={`mt-1 inline-flex items-center px-2 py-0.5 rounded text-[10px] font-medium border ${hasIssues ? 'bg-orange-50 text-orange-700 border-orange-200' : 'bg-slate-100 text-slate-600 border-slate-200'}`}>
                            {recipients.length} recipients
                          </div>
                        )}
                        <div className="text-xs text-gray-500 mt-1 truncate max-w-xs" title={event.message_id}>{event.message_id}</div>
                      </td>
                      <td className="px-4 py-4 text-sm text-gray-900 max-w-xs">
                        <div className="truncate font-medium" title={event.subject || 'No subject'}>
                          {event.subject || 'No subject'}
                        </div>
                        <div className="text-xs text-gray-500 mt-1 truncate" title={event.source}>{event.source}</div>
                      </td>
                      <td className="px-4 py-4 whitespace-nowrap">
                        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusColor(event.latest_status)}`}>
                          {event.latest_status || 'UNKNOWN'}
                        </span>
                      </td>
                      <td className="px-4 py-4 whitespace-nowrap text-sm text-gray-500">
                        {formatDate(event.last_event_at)}
                      </td>
                      <td className="px-4 py-4 whitespace-nowrap text-right text-sm text-gray-500">
                        <button
                          onClick={(e) => openDetail(event.message_id, e)}
                          className="inline-flex items-center p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded transition-colors"
                          title="View timeline"
                        >
                          <Eye className="w-5 h-5" />
                        </button>
                      </td>
                    </tr>
                    
                    {/* Expanded Row Panel */}
                    {isExpanded && (
                      <tr>
                        <td colSpan={6} className="p-0 border-b-0">
                          <div className="bg-slate-50/70 p-6 border-b border-gray-200 shadow-inner">
                            <h4 className="text-sm font-semibold text-gray-700 mb-3">Recipient Details</h4>
                            <div className="bg-white rounded border border-gray-200 overflow-hidden">
                              <table className="min-w-full divide-y divide-gray-200">
                                <thead className="bg-gray-50">
                                  <tr>
                                    <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Email</th>
                                    <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Type</th>
                                    <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                                    <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Diagnostic Code</th>
                                  </tr>
                                </thead>
                                <tbody className="divide-y divide-gray-100">
                                  {recipients.length > 0 ? recipients.map((r, i) => (
                                    <tr key={i} className="hover:bg-slate-50">
                                      <td className="px-4 py-3 text-sm font-medium text-gray-900">{r.email}</td>
                                      <td className="px-4 py-3 text-sm text-gray-500">{r.type || 'To'}</td>
                                      <td className="px-4 py-3 text-sm">
                                        <span className={`inline-flex items-center px-2 py-0.5 rounded text-[11px] font-medium ${getStatusColor(r.status)}`}>
                                          {r.status}
                                        </span>
                                      </td>
                                      <td className="px-4 py-3 text-sm text-gray-500 max-w-xs">
                                        {r.diagnostic_code && r.diagnostic_code !== '-' ? (
                                          <div className="flex items-center group">
                                            <span className="truncate" title={r.diagnostic_code}>{r.diagnostic_code}</span>
                                            <button 
                                              onClick={(e) => copyToClipboard(r.diagnostic_code, e)} 
                                              className="ml-2 opacity-0 group-hover:opacity-100 p-1 text-gray-400 hover:text-gray-700 hover:bg-gray-100 rounded transition-all"
                                              title="Copy to clipboard"
                                            >
                                              <Copy className="w-3.5 h-3.5" />
                                            </button>
                                          </div>
                                        ) : (
                                          <span className="text-gray-400">—</span>
                                        )}
                                      </td>
                                    </tr>
                                  )) : (
                                    <tr>
                                      <td colSpan={4} className="px-4 py-3 text-sm text-gray-500 text-center italic">
                                        No granular recipient data available (legacy data)
                                      </td>
                                    </tr>
                                  )}
                                </tbody>
                              </table>
                            </div>
                          </div>
                        </td>
                      </tr>
                    )}
                  </React.Fragment>
                );
              })}
            </tbody>
          </table>
        </div>

        <div className="px-6 py-4 border-t border-gray-200 bg-gray-50 flex items-center justify-between">
          <div className="text-sm text-gray-700">
            Showing <span className="font-medium">{pagination.total === 0 ? 0 : ((pagination.page - 1) * pagination.limit) + 1}</span> to <span className="font-medium">{Math.min(pagination.page * pagination.limit, pagination.total)}</span> of <span className="font-medium">{pagination.total}</span> results
          </div>
          <div className="flex items-center space-x-2">
            <button
              onClick={() => onPageChange(pagination.page - 1)}
              disabled={!pagination.hasPrev}
              className="px-3 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed flex items-center transition-colors duration-200"
            >
              <ChevronLeft className="w-4 h-4 mr-1" />
              Previous
            </button>

            <span className="px-3 py-2 text-sm text-gray-700 bg-white border border-gray-300 rounded-md">
              Page {pagination.page} of {pagination.totalPages}
            </span>

            <button
              onClick={() => onPageChange(pagination.page + 1)}
              disabled={!pagination.hasNext}
              className="px-3 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed flex items-center transition-colors duration-200"
            >
              Next
              <ChevronRight className="w-4 h-4 ml-1" />
            </button>
          </div>
        </div>
      </div>

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
              <button onClick={closeDetail} className="text-gray-400 hover:text-gray-600">×</button>
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
