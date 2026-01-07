import { useState, useEffect } from 'react';
import { Search, Plus, Trash2, Eye, AlertTriangle, RefreshCw, Clock, Activity, Upload } from 'lucide-react';
import Layout from '../components/Layout';

interface SyncStatus {
  last_sync: string;
  in_progress: boolean;
  next_sync_in: string;
  aws_enabled: boolean;
}

interface SuppressionResponse {
  suppressions: SuppressionEntry[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
  has_next: boolean;
  has_prev: boolean;
}

interface SuppressionEntry {
  id: number;
  email: string;
  suppression_type: string;
  reason: string;
  aws_status: string;
  is_active: boolean;
  created_at: string;
  added_by_name?: string;
}

const SuppressionPage = () => {
  const [suppressions, setSuppressions] = useState<SuppressionEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState('');
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(50);
  const [total, setTotal] = useState(0);
  const [totalPages, setTotalPages] = useState(0);
  const [hasNext, setHasNext] = useState(false);
  const [hasPrev, setHasPrev] = useState(false);
  const [showAddModal, setShowAddModal] = useState(false);
  const [showBulkModal, setShowBulkModal] = useState(false);
  const [bulkAction, setBulkAction] = useState<'add' | 'remove'>('add');
  const [bulkEmails, setBulkEmails] = useState('');
  const [bulkReason, setBulkReason] = useState('');
  const [selectedEmails, setSelectedEmails] = useState<string[]>([]);
  const [newEmail, setNewEmail] = useState('');
  const [newReason, setNewReason] = useState('');
  const [message, setMessage] = useState('');
  const [messageType, setMessageType] = useState<'success' | 'error' | ''>('');
  const [syncing, setSyncing] = useState(false);
  const [syncStatus, setSyncStatus] = useState<SyncStatus | null>(null);
  const [syncIntervalMinutes, setSyncIntervalMinutes] = useState<number | null>(null);

  const showMessage = (text: string, type: 'success' | 'error', timeout = 3000) => {
    setMessage(text);
    setMessageType(type);
    setTimeout(() => {
      setMessage('');
      setMessageType('');
    }, timeout);
  };

  const loadSyncStatus = async () => {
    try {
      const response = await fetch('/api/suppression/sync/status', {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
      });
      
      if (response.ok) {
        const data = await response.json();
        setSyncStatus(data);
      }

      const settingsResponse = await fetch('/api/settings/aws', {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
      });
      if (settingsResponse.ok) {
        const settingsData = await settingsResponse.json();
        if (typeof settingsData.sync_interval === 'number') {
          setSyncIntervalMinutes(settingsData.sync_interval);
        }
      }
    } catch (error) {
      console.error('Failed to load sync status:', error);
    }
  };

  const loadSuppressions = async () => {
    try {
      setLoading(true);
      const params = new URLSearchParams();
      params.append('page', page.toString());
      params.append('limit', limit.toString());
      if (search) params.append('search', search);
      
      const response = await fetch(`/api/suppression?${params}`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
      });
      
      if (response.ok) {
        const data: SuppressionResponse = await response.json();
        setSuppressions(data.suppressions || []);
        setTotal(data.total || 0);
        setTotalPages(data.total_pages || 0);
        setHasNext(data.has_next || false);
        setHasPrev(data.has_prev || false);
      }
    } catch (error) {
      console.error('Failed to load suppressions:', error);
    } finally {
      setLoading(false);
    }
  };

  const addSuppression = async () => {
    if (!newEmail) return;
    
    try {
      const response = await fetch('/api/suppression', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          email: newEmail,
          reason: newReason || 'Manually added'
        })
      });
      
      if (response.ok) {
        showMessage('Email added to suppression list', 'success');
        setNewEmail('');
        setNewReason('');
        setShowAddModal(false);
        loadSuppressions();
      } else {
        const error = await response.json();
        showMessage(error.error || 'Failed to add email', 'error', 5000);
      }
    } catch (error) {
      showMessage('Failed to add email', 'error', 5000);
    }
  };

  const removeSuppression = async (email: string) => {
    if (!confirm(`Remove ${email} from AWS SES suppression list?`)) return;
    
    try {
      const response = await fetch(`/api/suppression/${encodeURIComponent(email)}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
      });
      
      if (response.ok) {
        showMessage('Email removed from AWS SES suppression list', 'success');
        loadSuppressions();
      } else {
        const error = await response.json();
        showMessage(error.error || 'Failed to remove email', 'error', 5000);
      }
    } catch (error) {
      showMessage('Failed to remove email', 'error', 5000);
    }
  };

  const checkAWSStatus = async (email: string) => {
    try {
      const response = await fetch(`/api/suppression/${encodeURIComponent(email)}/status`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
      });
      
      if (response.ok) {
        const status = await response.json();
        alert(`AWS Status: ${status.suppressed ? 'Suppressed' : 'Not Suppressed'}\\nReason: ${status.reason}`);
      } else {
        const error = await response.json();
        alert(`Error: ${error.error}`);
      }
    } catch (error) {
      alert('Failed to check AWS status');
    }
  };

  const syncFromAWS = async () => {
    try {
      setSyncing(true);
      const response = await fetch('/api/suppression/sync', {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
      });
      
      if (response.ok) {
        showMessage('Sync triggered. Data will be updated in background.', 'success');
        loadSyncStatus(); // Refresh status
      } else {
        const error = await response.json();
        showMessage(error.error || 'Failed to trigger sync', 'error', 5000);
      }
    } catch (error) {
      showMessage('Failed to trigger sync', 'error', 5000);
    } finally {
      setSyncing(false);
    }
  };

  const bulkAddEmails = async () => {
    const emails = bulkEmails.split('\n').map(e => e.trim()).filter(e => e);
    if (emails.length === 0) return;
    
    try {
      const response = await fetch('/api/suppression/bulk', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          emails: emails,
          reason: bulkReason || 'Bulk added'
        })
      });
      
      if (response.ok) {
        const result = await response.json();
        showMessage(`Bulk add completed: ${result.success_count} success, ${result.failed_count} failed`, 'success', 5000);
        setBulkEmails('');
        setBulkReason('');
        setShowBulkModal(false);
        loadSuppressions();
      } else {
        const error = await response.json();
        showMessage(error.error || 'Failed to bulk add emails', 'error', 5000);
      }
    } catch (error) {
      showMessage('Failed to bulk add emails', 'error', 5000);
    }
  };

  const bulkRemoveEmails = async () => {
    const emails = selectedEmails.length > 0 ? selectedEmails : 
                  bulkEmails.split('\n').map(e => e.trim()).filter(e => e);
    
    if (emails.length === 0) return;
    
    if (!confirm(`Remove ${emails.length} emails from AWS SES suppression list?`)) return;
    
    try {
      const response = await fetch('/api/suppression/bulk', {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ emails: emails })
      });
      
      if (response.ok) {
        const result = await response.json();
        if (result.failed_count > 0) {
          showMessage(`Bulk remove completed: ${result.success_count} success, ${result.failed_count} failed. Failed emails: ${result.failed_emails.join(', ')}`, 'error', 5000);
        } else {
          showMessage(`Bulk remove completed: ${result.success_count} emails removed successfully`, 'success', 5000);
        }
        setSelectedEmails([]);
        setBulkEmails('');
        setShowBulkModal(false);
        loadSuppressions();
      } else {
        const error = await response.json();
        showMessage(error.error || 'Failed to bulk remove emails', 'error', 5000);
      }
    } catch (error) {
      showMessage('Failed to bulk remove emails', 'error', 5000);
    }
  };

  const toggleEmailSelection = (email: string) => {
    setSelectedEmails(prev => 
      prev.includes(email) 
        ? prev.filter(e => e !== email)
        : [...prev, email]
    );
  };

  const selectAllEmails = () => {
    if (selectedEmails.length === suppressions.length) {
      setSelectedEmails([]);
    } else {
      setSelectedEmails(suppressions.map(s => s.email));
    }
  };

  useEffect(() => {
    setPage(1); // Reset to first page when search changes
  }, [search]);

  useEffect(() => {
    loadSuppressions();
    loadSyncStatus();
    
    // Auto refresh sync status every 30 seconds
    const interval = setInterval(loadSyncStatus, 30000);
    return () => clearInterval(interval);
  }, [page, limit, search]);

  // Check if AWS is disabled
  const isAWSDisabled = syncStatus ? !syncStatus.aws_enabled : false;
  const nextSyncLabel = isAWSDisabled
    ? 'Disabled'
    : syncIntervalMinutes
    ? `${syncIntervalMinutes} minute${syncIntervalMinutes === 1 ? '' : 's'}`
    : (syncStatus?.next_sync_in || '—');

  return (
    <Layout title="Suppression List">
      <div className="space-y-6">
        {/* AWS Disabled Warning */}
        {isAWSDisabled && (
          <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
            <div className="flex items-center">
              <AlertTriangle className="w-5 h-5 text-yellow-600 mr-3" />
              <div>
                <h3 className="text-sm font-medium text-yellow-800">
                  AWS Integration Disabled
                </h3>
                <p className="text-sm text-yellow-700 mt-1">
                  Suppression list sync is disabled. Please configure AWS settings to enable this feature.
                </p>
              </div>
            </div>
          </div>
        )}

        {/* Sync Status Card */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold text-gray-900 flex items-center">
              <Activity className="w-5 h-5 mr-2 text-blue-600" />
              Background Sync Status
            </h2>
            <button
              onClick={loadSyncStatus}
              className="text-gray-400 hover:text-gray-600"
            >
              <RefreshCw className="w-4 h-4" />
            </button>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="flex items-center space-x-3">
              <div className={`w-3 h-3 rounded-full ${
                syncStatus?.in_progress ? 'bg-yellow-400 animate-pulse' : 'bg-green-400'
              }`}></div>
              <div>
                <p className="text-sm font-medium text-gray-900">
                  {syncStatus?.in_progress ? 'Syncing...' : 'Idle'}
                </p>
                <p className="text-xs text-gray-500">Current Status</p>
              </div>
            </div>
            
            <div className="flex items-center space-x-3">
              <Clock className="w-4 h-4 text-gray-400" />
              <div>
                <p className="text-sm font-medium text-gray-900">
                  {syncStatus?.last_sync ? 
                    new Date(syncStatus.last_sync).toLocaleString() : 
                    'Never'
                  }
                </p>
                <p className="text-xs text-gray-500">Last Sync</p>
              </div>
            </div>
            
            <div className="flex items-center space-x-3">
              <RefreshCw className="w-4 h-4 text-gray-400" />
              <div>
                <p className="text-sm font-medium text-gray-900">
                  {nextSyncLabel}
                </p>
                <p className="text-xs text-gray-500">Next Auto Sync</p>
              </div>
            </div>
          </div>
        </div>

        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Suppression List</h1>
            <p className="text-gray-600 mt-1">Manage email suppression list</p>
          </div>
          <div className="flex space-x-3">
            <button
              onClick={syncFromAWS}
              disabled={syncing || syncStatus?.in_progress || isAWSDisabled}
              className="inline-flex items-center px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed"
              title={isAWSDisabled ? "AWS integration is disabled" : ""}
            >
              <RefreshCw className={`w-4 h-4 mr-2 ${(syncing || syncStatus?.in_progress) ? 'animate-spin' : ''}`} />
              {isAWSDisabled ? 'AWS Disabled' : (syncing || syncStatus?.in_progress) ? 'Syncing...' : 'Trigger Sync'}
            </button>
            <button
              onClick={() => { setBulkAction('add'); setShowBulkModal(true); }}
              disabled={isAWSDisabled}
              className="inline-flex items-center px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 disabled:opacity-50 disabled:cursor-not-allowed"
              title={isAWSDisabled ? "AWS integration is disabled" : ""}
            >
              <Upload className="w-4 h-4 mr-2" />
              Bulk Add
            </button>
            {selectedEmails.length > 0 && (
              <button
                onClick={() => { setBulkAction('remove'); setShowBulkModal(true); }}
                disabled={isAWSDisabled}
                className="inline-flex items-center px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed"
                title={isAWSDisabled ? "AWS integration is disabled" : ""}
              >
                <Trash2 className="w-4 h-4 mr-2" />
                Remove Selected ({selectedEmails.length})
              </button>
            )}
            <button
              onClick={() => setShowAddModal(true)}
              disabled={isAWSDisabled}
              className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
              title={isAWSDisabled ? "AWS integration is disabled" : ""}
            >
              <Plus className="w-4 h-4 mr-2" />
              Add Email
            </button>
          </div>
        </div>

        {message && (
          <div className={`p-4 rounded-lg ${
            messageType === 'success' ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-700'
          }`}>
            {message}
          </div>
        )}

        {/* Search and Controls */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
          <div className="flex flex-col sm:flex-row gap-4 items-center justify-between">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
              <input
                type="text"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                placeholder="Search by email, reason, or source..."
                className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              />
            </div>
            <div className="flex items-center space-x-4">
              <div className="flex items-center space-x-2">
                <label className="text-sm text-gray-600">Show:</label>
                <select
                  value={limit}
                  onChange={(e) => {
                    setLimit(Number(e.target.value));
                    setPage(1);
                  }}
                  className="border border-gray-300 rounded px-3 py-1 text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                >
                  <option value={25}>25</option>
                  <option value={50}>50</option>
                  <option value={100}>100</option>
                  <option value={500}>500</option>
                </select>
                <span className="text-sm text-gray-600">per page</span>
              </div>
              <div className="text-sm text-gray-600">
                Showing {suppressions.length > 0 ? ((page - 1) * limit + 1) : 0} to {Math.min(page * limit, total)} of {total.toLocaleString()} entries
              </div>
            </div>
          </div>
        </div>

        {/* Suppression List */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200">
          <div className="overflow-x-auto">
            <table className="min-w-[1100px] w-full">
              <thead className="bg-gray-50 border-b border-gray-200 sticky top-0 z-10">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    <input
                      type="checkbox"
                      checked={selectedEmails.length === suppressions.length && suppressions.length > 0}
                      onChange={selectAllEmails}
                      className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                    />
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Email</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Type</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Reason</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">AWS Status</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Added</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                {loading ? (
                  [...Array(6)].map((_, i) => (
                    <tr key={i} className="animate-pulse">
                      <td className="px-6 py-4">
                        <div className="h-4 w-4 bg-gray-200 rounded" />
                      </td>
                      <td className="px-6 py-4">
                        <div className="h-4 bg-gray-200 rounded w-40" />
                      </td>
                      <td className="px-6 py-4">
                        <div className="h-4 bg-gray-200 rounded w-20" />
                      </td>
                      <td className="px-6 py-4">
                        <div className="h-4 bg-gray-200 rounded w-48" />
                      </td>
                      <td className="px-6 py-4">
                        <div className="h-4 bg-gray-200 rounded w-24" />
                      </td>
                      <td className="px-6 py-4">
                        <div className="h-4 bg-gray-200 rounded w-24" />
                      </td>
                      <td className="px-6 py-4">
                        <div className="h-4 bg-gray-200 rounded w-16" />
                      </td>
                    </tr>
                  ))
                ) : suppressions.length === 0 ? (
                  <tr>
                    <td colSpan={7} className="px-6 py-4 text-center text-gray-500">No suppressed emails found</td>
                  </tr>
                ) : (
                  suppressions.map((suppression) => (
                    <tr key={suppression.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4">
                        <input
                          type="checkbox"
                          checked={selectedEmails.includes(suppression.email)}
                          onChange={() => toggleEmailSelection(suppression.email)}
                          className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                        />
                      </td>
                      <td className="px-6 py-4 text-sm font-medium text-gray-900">{suppression.email}</td>
                      <td className="px-6 py-4 text-sm text-gray-500">
                        <span className={`px-2 py-1 text-xs rounded-full ${
                          suppression.suppression_type === 'bounce' ? 'bg-red-100 text-red-800' :
                          suppression.suppression_type === 'complaint' ? 'bg-orange-100 text-orange-800' :
                          'bg-gray-100 text-gray-800'
                        }`}>
                          {suppression.suppression_type}
                        </span>
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-500">{suppression.reason}</td>
                      <td className="px-6 py-4 text-sm text-gray-500">
                        <span className={`px-2 py-1 text-xs rounded-full ${
                          suppression.aws_status === 'suppressed' ? 'bg-red-100 text-red-800' :
                          suppression.aws_status === 'not_suppressed' ? 'bg-green-100 text-green-800' :
                          'bg-gray-100 text-gray-800'
                        }`}>
                          {suppression.aws_status}
                        </span>
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-500">
                        {new Date(suppression.created_at).toLocaleDateString()}
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-500">
                        <div className="flex space-x-2">
                          <button
                            onClick={() => checkAWSStatus(suppression.email)}
                            className="text-blue-600 hover:text-blue-800"
                            title="Check AWS Status"
                          >
                            <Eye className="w-4 h-4" />
                          </button>
                          <button
                            onClick={() => removeSuppression(suppression.email)}
                            className="text-red-600 hover:text-red-800"
                            title="Remove from list"
                          >
                            <Trash2 className="w-4 h-4" />
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
          
          {/* Pagination */}
          {totalPages > 1 && (
            <div className="px-6 py-4 border-t border-gray-200 flex items-center justify-between">
              <div className="flex items-center space-x-2">
                <button
                  onClick={() => setPage(1)}
                  disabled={!hasPrev}
                  className="px-3 py-1 text-sm border border-gray-300 rounded hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  First
                </button>
                <button
                  onClick={() => setPage(page - 1)}
                  disabled={!hasPrev}
                  className="px-3 py-1 text-sm border border-gray-300 rounded hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Previous
                </button>
              </div>
              
              <div className="flex items-center space-x-2">
                {/* Page numbers */}
                {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                  let pageNum;
                  if (totalPages <= 5) {
                    pageNum = i + 1;
                  } else if (page <= 3) {
                    pageNum = i + 1;
                  } else if (page >= totalPages - 2) {
                    pageNum = totalPages - 4 + i;
                  } else {
                    pageNum = page - 2 + i;
                  }
                  
                  return (
                    <button
                      key={pageNum}
                      onClick={() => setPage(pageNum)}
                      className={`px-3 py-1 text-sm border rounded ${
                        page === pageNum
                          ? 'bg-blue-600 text-white border-blue-600'
                          : 'border-gray-300 hover:bg-gray-50'
                      }`}
                    >
                      {pageNum}
                    </button>
                  );
                })}
              </div>
              
              <div className="flex items-center space-x-2">
                <button
                  onClick={() => setPage(page + 1)}
                  disabled={!hasNext}
                  className="px-3 py-1 text-sm border border-gray-300 rounded hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Next
                </button>
                <button
                  onClick={() => setPage(totalPages)}
                  disabled={!hasNext}
                  className="px-3 py-1 text-sm border border-gray-300 rounded hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Last
                </button>
              </div>
            </div>
          )}
        </div>

        {/* Bulk Action Modal */}
        {showBulkModal && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg p-6 w-full max-w-2xl">
              <h3 className="text-lg font-semibold mb-4">
                {bulkAction === 'add' ? 'Bulk Add Emails to Suppression List' : 'Bulk Remove Emails from Suppression List'}
              </h3>
              
              {bulkAction === 'add' ? (
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Email Addresses (one per line)
                    </label>
                    <textarea
                      value={bulkEmails}
                      onChange={(e) => setBulkEmails(e.target.value)}
                      rows={8}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                      placeholder={`user1@example.com\nuser2@example.com\nuser3@example.com`}
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">Reason</label>
                    <input
                      type="text"
                      value={bulkReason}
                      onChange={(e) => setBulkReason(e.target.value)}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                      placeholder="Reason for bulk suppression"
                    />
                  </div>
                </div>
              ) : (
                <div className="space-y-4">
                  {selectedEmails.length > 0 ? (
                    <div>
                      <p className="text-sm text-gray-600 mb-2">
                        Selected emails to remove ({selectedEmails.length}):
                      </p>
                      <div className="max-h-40 overflow-y-auto bg-gray-50 p-3 rounded border">
                        {selectedEmails.map(email => (
                          <div key={email} className="text-sm text-gray-700">{email}</div>
                        ))}
                      </div>
                    </div>
                  ) : (
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Email Addresses to Remove (one per line)
                      </label>
                      <textarea
                        value={bulkEmails}
                        onChange={(e) => setBulkEmails(e.target.value)}
                        rows={8}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                        placeholder={`user1@example.com\nuser2@example.com\nuser3@example.com`}
                      />
                    </div>
                  )}
                </div>
              )}
              
              <div className="flex justify-end space-x-3 mt-6">
                <button
                  onClick={() => {
                    setShowBulkModal(false);
                    setBulkEmails('');
                    setBulkReason('');
                  }}
                  className="px-4 py-2 text-gray-600 border border-gray-300 rounded-lg hover:bg-gray-50"
                >
                  Cancel
                </button>
                <button
                  onClick={bulkAction === 'add' ? bulkAddEmails : bulkRemoveEmails}
                  className={`px-4 py-2 text-white rounded-lg ${
                    bulkAction === 'add' 
                      ? 'bg-purple-600 hover:bg-purple-700' 
                      : 'bg-red-600 hover:bg-red-700'
                  }`}
                >
                  {bulkAction === 'add' ? 'Add Emails' : 'Remove Emails'}
                </button>
              </div>
            </div>
          </div>
        )}

        {/* Add Email Modal */}
        {showAddModal && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg p-6 w-full max-w-md">
              <h3 className="text-lg font-semibold mb-4">Add Email to Suppression List</h3>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">Email Address</label>
                  <input
                    type="email"
                    value={newEmail}
                    onChange={(e) => setNewEmail(e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    placeholder="user@example.com"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">Reason</label>
                  <input
                    type="text"
                    value={newReason}
                    onChange={(e) => setNewReason(e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    placeholder="Reason for suppression"
                  />
                </div>
              </div>
              <div className="flex justify-end space-x-3 mt-6">
                <button
                  onClick={() => setShowAddModal(false)}
                  className="px-4 py-2 text-gray-600 border border-gray-300 rounded-lg hover:bg-gray-50"
                >
                  Cancel
                </button>
                <button
                  onClick={addSuppression}
                  className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
                >
                  Add Email
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </Layout>
  );
};

export default SuppressionPage;
