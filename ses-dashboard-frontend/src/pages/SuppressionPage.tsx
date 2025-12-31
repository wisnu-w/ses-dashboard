import { useState, useEffect } from 'react';
import { Search, Plus, Trash2, Eye, AlertTriangle, RefreshCw, Clock, Activity, Upload } from 'lucide-react';
import Layout from '../components/Layout';

interface SyncStatus {
  last_sync: string;
  in_progress: boolean;
  next_sync_in: string;
  aws_enabled: boolean;
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
  const [showAddModal, setShowAddModal] = useState(false);
  const [showBulkModal, setShowBulkModal] = useState(false);
  const [bulkAction, setBulkAction] = useState<'add' | 'remove'>('add');
  const [bulkEmails, setBulkEmails] = useState('');
  const [bulkReason, setBulkReason] = useState('');
  const [selectedEmails, setSelectedEmails] = useState<string[]>([]);
  const [newEmail, setNewEmail] = useState('');
  const [newReason, setNewReason] = useState('');
  const [message, setMessage] = useState('');
  const [syncing, setSyncing] = useState(false);
  const [syncStatus, setSyncStatus] = useState<SyncStatus | null>(null);

  const loadSyncStatus = async () => {
    try {
      const response = await fetch('/api/suppression/sync/status', {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
      });
      
      if (response.ok) {
        const data = await response.json();
        setSyncStatus(data);
      }
    } catch (error) {
      console.error('Failed to load sync status:', error);
    }
  };

  const loadSuppressions = async () => {
    try {
      setLoading(true);
      const params = new URLSearchParams();
      if (search) params.append('search', search);
      
      const response = await fetch(`/api/suppression?${params}`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
      });
      
      if (response.ok) {
        const data = await response.json();
        setSuppressions(data.suppressions || []);
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
        setMessage('Email added to suppression list');
        setNewEmail('');
        setNewReason('');
        setShowAddModal(false);
        loadSuppressions();
      } else {
        const error = await response.json();
        setMessage(error.error || 'Failed to add email');
      }
    } catch (error) {
      setMessage('Failed to add email');
    }
    
    setTimeout(() => setMessage(''), 3000);
  };

  const removeSuppression = async (email: string) => {
    if (!confirm(`Remove ${email} from AWS SES suppression list?`)) return;
    
    try {
      const response = await fetch(`/api/suppression/${encodeURIComponent(email)}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
      });
      
      if (response.ok) {
        setMessage('Email removed from AWS SES suppression list');
        loadSuppressions();
      } else {
        const error = await response.json();
        setMessage(error.error || 'Failed to remove email');
      }
    } catch (error) {
      setMessage('Failed to remove email');
    }
    
    setTimeout(() => setMessage(''), 3000);
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
        setMessage('Sync triggered. Data will be updated in background.');
        loadSyncStatus(); // Refresh status
      } else {
        const error = await response.json();
        setMessage(error.error || 'Failed to trigger sync');
      }
    } catch (error) {
      setMessage('Failed to trigger sync');
    } finally {
      setSyncing(false);
    }
    
    setTimeout(() => setMessage(''), 3000);
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
        setMessage(`Bulk add completed: ${result.success_count} success, ${result.failed_count} failed`);
        setBulkEmails('');
        setBulkReason('');
        setShowBulkModal(false);
        loadSuppressions();
      } else {
        const error = await response.json();
        setMessage(error.error || 'Failed to bulk add emails');
      }
    } catch (error) {
      setMessage('Failed to bulk add emails');
    }
    
    setTimeout(() => setMessage(''), 5000);
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
        setMessage(`Bulk remove completed: ${result.success_count} success, ${result.failed_count} failed`);
        setSelectedEmails([]);
        setBulkEmails('');
        setShowBulkModal(false);
        loadSuppressions();
      } else {
        const error = await response.json();
        setMessage(error.error || 'Failed to bulk remove emails');
      }
    } catch (error) {
      setMessage('Failed to bulk remove emails');
    }
    
    setTimeout(() => setMessage(''), 5000);
  };

  const toggleEmailSelection = (email: string) => {
    setSelectedEmails(prev => 
      prev.includes(email) 
        ? prev.filter(e => e !== email)
        : [...prev, email]
    );
  };

  const selectAllEmails = () => {
    if (selectedEmails.length === filteredSuppressions.length) {
      setSelectedEmails([]);
    } else {
      setSelectedEmails(filteredSuppressions.map(s => s.email));
    }
  };

  useEffect(() => {
    loadSuppressions();
    loadSyncStatus();
    
    // Auto refresh sync status every 30 seconds
    const interval = setInterval(loadSyncStatus, 30000);
    return () => clearInterval(interval);
  }, [search]);

  // Check if AWS is disabled
  const isAWSDisabled = syncStatus ? !syncStatus.aws_enabled : false;

  const filteredSuppressions = suppressions.filter(s => 
    s.email.toLowerCase().includes(search.toLowerCase()) ||
    s.reason.toLowerCase().includes(search.toLowerCase())
  );

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
                  {syncStatus?.next_sync_in || '5 minutes'}
                </p>
                <p className="text-xs text-gray-500">Next Auto Sync</p>
              </div>
            </div>
          </div>
        </div>

        <div className="flex items-center justify-between">
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
          <div className="p-4 rounded-lg bg-blue-50 text-blue-700">
            {message}
          </div>
        )}

        {/* Search */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
            <input
              type="text"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder="Search by email or reason..."
              className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
        </div>

        {/* Suppression List */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50 border-b border-gray-200">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    <input
                      type="checkbox"
                      checked={selectedEmails.length === filteredSuppressions.length && filteredSuppressions.length > 0}
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
                  <tr>
                    <td colSpan={7} className="px-6 py-4 text-center text-gray-500">Loading...</td>
                  </tr>
                ) : filteredSuppressions.length === 0 ? (
                  <tr>
                    <td colSpan={7} className="px-6 py-4 text-center text-gray-500">No suppressed emails found</td>
                  </tr>
                ) : (
                  filteredSuppressions.map((suppression) => (
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