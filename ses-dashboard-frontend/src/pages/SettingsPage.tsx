import { useState, useEffect } from 'react';
import { Settings, Save, TestTube, AlertCircle, CheckCircle } from 'lucide-react';
import Layout from '../components/Layout';

interface AWSSettings {
  enabled: boolean;
  region: string;
  access_key: string;
  secret_key: string;
  sync_interval: number;
}

interface RetentionSettings {
  retention_days: number;
  enabled: boolean;
}

interface TimezoneSettings {
  timezone: string;
}

const SettingsPage = () => {
  const [settings, setSettings] = useState<AWSSettings>({
    enabled: false,
    region: 'us-east-1',
    access_key: '',
    secret_key: '',
    sync_interval: 5
  });
  const [retentionSettings, setRetentionSettings] = useState<RetentionSettings>({
    retention_days: 30,
    enabled: true
  });
  const [timezoneSettings, setTimezoneSettings] = useState<TimezoneSettings>({
    timezone: 'Asia/Jakarta'
  });
  const [initialLoading, setInitialLoading] = useState(true);
  const [loading, setLoading] = useState(false);
  const [testing, setTesting] = useState(false);
  const [message, setMessage] = useState<{ text: string; type: 'success' | 'error' | 'info' } | null>(null);

  const setTimedMessage = (text: string, type: 'success' | 'error' | 'info', timeout = 3000) => {
    setMessage({ text, type });
    setTimeout(() => setMessage(null), timeout);
  };

  const loadSettings = async (showLoading = false) => {
    try {
      if (showLoading) setInitialLoading(true);
      // Load AWS settings
      const awsResponse = await fetch('/api/settings/aws', {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
      });
      
      if (awsResponse.ok) {
        const awsData = await awsResponse.json();
        setSettings(awsData);
      }
      
      // Load retention settings
      const retentionResponse = await fetch('/api/settings/retention', {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
      });
      
      if (retentionResponse.ok) {
        const retentionData = await retentionResponse.json();
        setRetentionSettings(retentionData);
      } else {
        console.error('Retention response error:', retentionResponse.status);
      }
      
      // Load timezone settings
      const timezoneResponse = await fetch('/api/settings/timezone', {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
      });
      
      if (timezoneResponse.ok) {
        const timezoneData = await timezoneResponse.json();
        setTimezoneSettings(timezoneData);
      }
    } catch (error) {
      console.error('Failed to load settings:', error);
      if (showLoading) {
        setTimedMessage('Failed to load settings', 'error', 5000);
      }
    } finally {
      if (showLoading) setInitialLoading(false);
    }
  };

  const saveRetentionSettings = async () => {
    try {
      setLoading(true);
      const response = await fetch('/api/settings/retention', {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(retentionSettings)
      });
      
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Failed to save retention settings');
      }
      
      setTimedMessage('Retention settings saved successfully', 'success');
    } catch (error: any) {
      setTimedMessage(error.message || 'Failed to save retention settings', 'error', 5000);
    } finally {
      setLoading(false);
    }
  };

  const saveTimezoneSettings = async () => {
    try {
      setLoading(true);
      const response = await fetch('/api/settings/timezone', {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(timezoneSettings)
      });
      
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Failed to save timezone settings');
      }
      
      setTimedMessage('Timezone settings saved successfully', 'success');
    } catch (error: any) {
      setTimedMessage(error.message || 'Failed to save timezone settings', 'error', 5000);
    } finally {
      setLoading(false);
    }
  };

  const saveSettings = async () => {
    try {
      setLoading(true);
      const response = await fetch('/api/settings/aws', {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(settings)
      });
      
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Failed to save settings');
      }
      
      setTimedMessage('Settings saved successfully', 'success');
    } catch (error: any) {
      setTimedMessage(error.message || 'Failed to save settings', 'error', 5000);
    } finally {
      setLoading(false);
    }
  };

  const testConnection = async () => {
    try {
      setTesting(true);
      const response = await fetch('/api/settings/aws/test', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(settings)
      });
      
      if (response.ok) {
        setTimedMessage('AWS connection successful', 'success', 5000);
      } else {
        const errorData = await response.json();
        setTimedMessage(errorData.error || 'AWS connection failed', 'error', 5000);
      }
    } catch (error: any) {
      setTimedMessage(error.message || 'Connection test failed', 'error', 5000);
    } finally {
      setTesting(false);
    }
  };

  useEffect(() => {
    loadSettings(true);
    
    // Auto refresh settings setiap 30 detik
    const interval = setInterval(loadSettings, 30000);
    return () => clearInterval(interval);
  }, []);

  if (initialLoading) {
    return (
      <Layout title="Settings">
        <div className="animate-pulse space-y-6">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
              <div className="h-4 bg-gray-200 rounded w-1/3 mb-4"></div>
              <div className="space-y-3">
                <div className="h-10 bg-gray-200 rounded"></div>
                <div className="h-10 bg-gray-200 rounded"></div>
              </div>
            </div>
          ))}
        </div>
      </Layout>
    );
  }

  return (
    <Layout title="Settings">
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">System Settings</h1>
            <p className="text-gray-600 mt-1">Configure AWS integration and advanced features</p>
          </div>
        </div>

        {message && (
          <div className={`p-4 rounded-lg ${
            message.type === 'success'
              ? 'bg-green-50 text-green-700'
              : message.type === 'info'
              ? 'bg-blue-50 text-blue-700'
              : 'bg-red-50 text-red-700'
          }`}>
            <div className="flex items-center">
              {message.type === 'success' ? <CheckCircle className="w-5 h-5 mr-2" /> : <AlertCircle className="w-5 h-5 mr-2" />}
              {message.text}
            </div>
          </div>
        )}

        {/* Timezone Settings */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <div className="flex items-center mb-6">
            <Settings className="w-6 h-6 text-purple-600 mr-3" />
            <h2 className="text-xl font-semibold">Timezone Settings</h2>
          </div>

          <div className="space-y-6">
            <div className="max-w-md">
              <label className="block text-sm font-medium text-gray-700 mb-2">Application Timezone</label>
              <select
                value={timezoneSettings.timezone}
                onChange={(e) => setTimezoneSettings({...timezoneSettings, timezone: e.target.value})}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
              >
                <option value="Asia/Jakarta">Asia/Jakarta (WIB)</option>
                <option value="Asia/Makassar">Asia/Makassar (WITA)</option>
                <option value="Asia/Jayapura">Asia/Jayapura (WIT)</option>
                <option value="UTC">UTC</option>
                <option value="America/New_York">America/New_York (EST)</option>
                <option value="Europe/London">Europe/London (GMT)</option>
                <option value="Asia/Singapore">Asia/Singapore (SGT)</option>
                <option value="Asia/Tokyo">Asia/Tokyo (JST)</option>
              </select>
              <p className="text-xs text-gray-500 mt-1">
                This timezone will be used for displaying dates and times in metrics and monitoring handlers
              </p>
            </div>
          </div>

          <div className="flex justify-end mt-8 pt-6 border-t border-gray-200">
            <button
              onClick={saveTimezoneSettings}
              disabled={loading}
              className="inline-flex items-center px-6 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 disabled:opacity-50"
            >
              <Save className="w-4 h-4 mr-2" />
              {loading ? 'Saving...' : 'Save Timezone Settings'}
            </button>
          </div>
        </div>

        {/* Retention Settings */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <div className="flex items-center mb-6">
            <Settings className="w-6 h-6 text-green-600 mr-3" />
            <h2 className="text-xl font-semibold">Event Log Retention</h2>
          </div>

          <div className="space-y-6">
            <div className="flex items-center">
              <input
                type="checkbox"
                id="retention-enabled"
                checked={retentionSettings.enabled}
                onChange={(e) => setRetentionSettings({...retentionSettings, enabled: e.target.checked})}
                className="h-4 w-4 text-green-600 focus:ring-green-500 border-gray-300 rounded"
              />
              <label htmlFor="retention-enabled" className="ml-2 text-sm font-medium text-gray-700">
                Enable automatic event log cleanup
              </label>
            </div>

            {retentionSettings.enabled && (
              <div className="pl-6 border-l-2 border-green-200">
                <div className="max-w-md">
                  <label className="block text-sm font-medium text-gray-700 mb-2">Retention Period</label>
                  <select
                    value={retentionSettings.retention_days}
                    onChange={(e) => setRetentionSettings({...retentionSettings, retention_days: parseInt(e.target.value)})}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-green-500"
                  >
                    <option value={0}>Never delete (Keep forever)</option>
                    <option value={7}>7 days</option>
                    <option value={14}>14 days</option>
                    <option value={30}>30 days</option>
                    <option value={60}>60 days</option>
                    <option value={90}>90 days</option>
                    <option value={180}>6 months</option>
                    <option value={365}>1 year</option>
                  </select>
                  <p className="text-xs text-gray-500 mt-1">
                    {retentionSettings.retention_days === 0 
                      ? 'Event logs will be kept forever'
                      : `Event logs older than ${retentionSettings.retention_days} days will be automatically deleted`
                    }
                  </p>
                </div>
              </div>
            )}
          </div>

          <div className="flex justify-end mt-8 pt-6 border-t border-gray-200">
            <button
              onClick={saveRetentionSettings}
              disabled={loading}
              className="inline-flex items-center px-6 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 disabled:opacity-50"
            >
              <Save className="w-4 h-4 mr-2" />
              {loading ? 'Saving...' : 'Save Retention Settings'}
            </button>
          </div>
        </div>

        {/* AWS Settings */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <div className="flex items-center mb-6">
            <Settings className="w-6 h-6 text-blue-600 mr-3" />
            <h2 className="text-xl font-semibold">AWS SES Integration</h2>
          </div>

          <div className="space-y-6">
            <div className="flex items-center">
              <input
                type="checkbox"
                id="aws-enabled"
                checked={settings.enabled}
                onChange={(e) => setSettings({...settings, enabled: e.target.checked})}
                className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
              />
              <label htmlFor="aws-enabled" className="ml-2 text-sm font-medium text-gray-700">
                Enable AWS SES Advanced Features
              </label>
            </div>

            {settings.enabled && (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6 pl-6 border-l-2 border-blue-200">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">AWS Region</label>
                  <select
                    value={settings.region}
                    onChange={(e) => setSettings({...settings, region: e.target.value})}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option value="us-east-1">US East (N. Virginia)</option>
                    <option value="us-west-2">US West (Oregon)</option>
                    <option value="eu-west-1">Europe (Ireland)</option>
                    <option value="ap-southeast-1">Asia Pacific (Singapore)</option>
                    <option value="ap-southeast-2">Asia Pacific (Sydney)</option>
                  </select>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">AWS Access Key</label>
                  <input
                    type="text"
                    value={settings.access_key}
                    onChange={(e) => setSettings({...settings, access_key: e.target.value})}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    placeholder="AKIA..."
                  />
                </div>

                <div className="md:col-span-2">
                  <label className="block text-sm font-medium text-gray-700 mb-2">AWS Secret Key</label>
                  <input
                    type="password"
                    value={settings.secret_key}
                    onChange={(e) => setSettings({...settings, secret_key: e.target.value})}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    placeholder="Enter secret key"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">Auto Sync Interval</label>
                  <select
                    value={settings.sync_interval}
                    onChange={(e) => setSettings({...settings, sync_interval: parseInt(e.target.value)})}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option value={1}>1 minute</option>
                    <option value={5}>5 minutes</option>
                    <option value={10}>10 minutes</option>
                    <option value={15}>15 minutes</option>
                    <option value={30}>30 minutes</option>
                    <option value={60}>1 hour</option>
                  </select>
                  <p className="text-xs text-gray-500 mt-1">
                    How often to sync suppression list from AWS
                  </p>
                </div>

                <div className="md:col-span-2 flex space-x-3">
                  <button
                    onClick={testConnection}
                    disabled={testing || !settings.access_key || !settings.secret_key}
                    className="inline-flex items-center px-4 py-2 bg-yellow-600 text-white rounded-lg hover:bg-yellow-700 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <TestTube className="w-4 h-4 mr-2" />
                    {testing ? 'Testing...' : 'Test Connection'}
                  </button>
                </div>
              </div>
            )}
          </div>

          <div className="flex justify-end mt-8 pt-6 border-t border-gray-200">
            <button
              onClick={saveSettings}
              disabled={loading}
              className="inline-flex items-center px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
            >
              <Save className="w-4 h-4 mr-2" />
              {loading ? 'Saving...' : 'Save Settings'}
            </button>
          </div>
        </div>
      </div>
    </Layout>
  );
};

export default SettingsPage;
