import { useMemo, useState } from 'react';
import { ChevronLeft, ChevronRight, Mail, Send, CheckCircle, XCircle, AlertTriangle, Eye, MousePointer } from 'lucide-react';
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
    case 'send':
      return <Send className="w-4 h-4 text-blue-500" />;
    case 'delivery':
      return <CheckCircle className="w-4 h-4 text-green-500" />;
    case 'bounce':
      return <XCircle className="w-4 h-4 text-red-500" />;
    case 'complaint':
      return <AlertTriangle className="w-4 h-4 text-orange-500" />;
    case 'open':
      return <Eye className="w-4 h-4 text-purple-500" />;
    case 'click':
      return <MousePointer className="w-4 h-4 text-indigo-500" />;
    default:
      return <Mail className="w-4 h-4 text-gray-500" />;
  }
};

const getEventTypeColor = (eventType: string) => {
  const type = eventType?.toLowerCase() || '';
  switch (type) {
    case 'send':
      return 'bg-blue-100 text-blue-800';
    case 'delivery':
      return 'bg-green-100 text-green-800';
    case 'bounce':
      return 'bg-red-100 text-red-800';
    case 'complaint':
      return 'bg-orange-100 text-orange-800';
    case 'open':
      return 'bg-purple-100 text-purple-800';
    case 'click':
      return 'bg-indigo-100 text-indigo-800';
    default:
      return 'bg-gray-100 text-gray-800';
  }
};

const getStatusColor = (status: string) => {
  switch (status) {
    case 'SUCCESS':
      return 'bg-green-100 text-green-800';
    case 'FAILED':
      return 'bg-red-100 text-red-800';
    default:
      return 'bg-gray-100 text-gray-800';
  }
};

const formatDate = (dateString: string) => {
  if (!dateString) return 'N/A';
  const date = new Date(dateString);
  if (Number.isNaN(date.getTime())) return 'Invalid Date';
  return date.toLocaleString();
};

const EventsTable = ({ events, pagination, onPageChange, loading = false }: EventsTableProps) => {
  const [selectedMessageId, setSelectedMessageId] = useState<string | null>(null);
  const [detailEvents, setDetailEvents] = useState<Event[]>([]);
  const [detailLoading, setDetailLoading] = useState(false);
  const [detailError, setDetailError] = useState('');

  const currentGroup = useMemo(
    () => events.find((event) => event.message_id === selectedMessageId) ?? null,
    [events, selectedMessageId]
  );

  const openDetail = async (messageId: string) => {
    try {
      setSelectedMessageId(messageId);
      setDetailLoading(true);
      setDetailError('');
      const response = await eventsService.getEventDetail(messageId);
      setDetailEvents(response.events || []);
    } catch (error) {
      console.error('Failed to load event detail:', error);
      setDetailError('Failed to load event timeline');
      setDetailEvents([]);
    } finally {
      setDetailLoading(false);
    }
  };

  const closeDetail = () => {
    setSelectedMessageId(null);
    setDetailEvents([]);
    setDetailError('');
  };

  if (loading) {
    return (
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <div className="animate-pulse">
          <div className="h-4 bg-gray-200 rounded w-1/4 mb-4"></div>
          <div className="space-y-3">
            {[...Array(5)].map((_, i) => (
              <div key={i} className="h-4 bg-gray-200 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  return (
    <>
      <div className="bg-white rounded-lg shadow-sm border border-gray-200">
        <div className="px-6 py-4 border-b border-gray-200 bg-gray-50">
          <h3 className="text-lg font-semibold text-gray-900">Message Threads</h3>
          <p className="text-sm text-gray-600 mt-1">Grouped by SES message ID for easier tracing</p>
        </div>

        <div className="overflow-hidden">
          <table className="w-full table-fixed divide-y divide-gray-200">
            <thead className="bg-gray-50 sticky top-0 z-10">
              <tr>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-[22%]">Message ID</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-[22%]">Email</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-[24%]">Subject</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-[12%]">Status</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-[14%]">Last Event</th>
                <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider w-[6%]">Action</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {events.map((event) => (
                <tr key={event.message_id} className="hover:bg-gray-50 transition-colors duration-150">
                  <td className="px-4 py-4 text-sm text-gray-900">
                    <div className="truncate font-medium" title={event.message_id}>{event.message_id}</div>
                  </td>
                  <td className="px-4 py-4 text-sm text-gray-900">
                    <div className="truncate" title={event.email}>{event.email}</div>
                  </td>
                  <td className="px-4 py-4 text-sm text-gray-900">
                    <div className="truncate" title={event.subject || 'No subject'}>
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
                      onClick={() => openDetail(event.message_id)}
                      className="inline-flex items-center px-2.5 py-1.5 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors duration-200"
                      title="View timeline"
                    >
                      <Eye className="w-4 h-4" />
                    </button>
                  </td>
                </tr>
              ))}
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
