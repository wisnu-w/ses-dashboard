class SyncMonitor {
    constructor() {
        this.apiBase = '/api';
        this.token = localStorage.getItem('jwt_token') || 'your-jwt-token-here';
        this.refreshInterval = null;
        this.init();
    }

    init() {
        this.bindEvents();
        this.startAutoRefresh();
        this.loadSyncStatus();
        this.loadSuppressions();
    }

    bindEvents() {
        document.getElementById('manual-sync-btn').addEventListener('click', () => {
            this.triggerManualSync();
        });

        document.getElementById('refresh-btn').addEventListener('click', () => {
            this.loadSyncStatus();
            this.loadSuppressions();
        });
    }

    startAutoRefresh() {
        // Refresh status setiap 30 detik
        this.refreshInterval = setInterval(() => {
            this.loadSyncStatus();
        }, 30000);
    }

    async makeRequest(url, options = {}) {
        const defaultOptions = {
            headers: {
                'Authorization': `Bearer ${this.token}`,
                'Content-Type': 'application/json',
                ...options.headers
            }
        };

        try {
            const response = await fetch(url, { ...defaultOptions, ...options });
            
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
            
            return await response.json();
        } catch (error) {
            this.addLog(`Error: ${error.message}`, 'error');
            throw error;
        }
    }

    async loadSyncStatus() {
        try {
            const data = await this.makeRequest(`${this.apiBase}/suppression/sync/status`);
            this.updateSyncStatus(data);
        } catch (error) {
            this.updateSyncStatusError();
        }
    }

    updateSyncStatus(data) {
        const statusElement = document.getElementById('sync-status');
        const lastSyncElement = document.getElementById('last-sync');
        
        // Update status
        if (data.in_progress) {
            statusElement.textContent = 'Syncing...';
            statusElement.className = 'status-badge syncing';
        } else {
            statusElement.textContent = 'Idle';
            statusElement.className = 'status-badge idle';
        }

        // Update last sync time
        if (data.last_sync && data.last_sync !== '0001-01-01T00:00:00Z') {
            const lastSync = new Date(data.last_sync);
            lastSyncElement.textContent = this.formatDateTime(lastSync);
        } else {
            lastSyncElement.textContent = 'Never';
        }

        // Enable/disable manual sync button
        const manualSyncBtn = document.getElementById('manual-sync-btn');
        manualSyncBtn.disabled = data.in_progress;
    }

    updateSyncStatusError() {
        const statusElement = document.getElementById('sync-status');
        statusElement.textContent = 'Error';
        statusElement.className = 'status-badge error';
    }

    async triggerManualSync() {
        try {
            const manualSyncBtn = document.getElementById('manual-sync-btn');
            manualSyncBtn.disabled = true;
            manualSyncBtn.textContent = 'Syncing...';

            await this.makeRequest(`${this.apiBase}/suppression/sync`, {
                method: 'POST'
            });

            this.addLog('Manual sync triggered successfully', 'success');
            
            // Refresh status setelah 2 detik
            setTimeout(() => {
                this.loadSyncStatus();
                this.loadSuppressions();
            }, 2000);

        } catch (error) {
            this.addLog(`Failed to trigger sync: ${error.message}`, 'error');
        } finally {
            setTimeout(() => {
                const manualSyncBtn = document.getElementById('manual-sync-btn');
                manualSyncBtn.disabled = false;
                manualSyncBtn.textContent = 'Trigger Manual Sync';
            }, 3000);
        }
    }

    async loadSuppressions() {
        try {
            const data = await this.makeRequest(`${this.apiBase}/suppression?limit=10`);
            this.updateSuppressionsTable(data.suppressions || []);
        } catch (error) {
            this.updateSuppressionsError();
        }
    }

    updateSuppressionsTable(suppressions) {
        const tbody = document.getElementById('suppressions-tbody');
        
        if (suppressions.length === 0) {
            tbody.innerHTML = '<tr><td colspan="4" class="loading">No suppressions found</td></tr>';
            return;
        }

        tbody.innerHTML = suppressions.map(suppression => `
            <tr>
                <td>${suppression.email || 'N/A'}</td>
                <td>${suppression.reason || 'N/A'}</td>
                <td>${suppression.suppression_type || 'AWS'}</td>
                <td>${suppression.created_at ? this.formatDateTime(new Date(suppression.created_at)) : 'N/A'}</td>
            </tr>
        `).join('');
    }

    updateSuppressionsError() {
        const tbody = document.getElementById('suppressions-tbody');
        tbody.innerHTML = '<tr><td colspan="4" class="loading">Error loading suppressions</td></tr>';
    }

    addLog(message, type = 'info') {
        const logsContainer = document.getElementById('logs-container');
        const logEntry = document.createElement('div');
        logEntry.className = `log-entry ${type}`;
        logEntry.textContent = `[${this.formatTime(new Date())}] ${message}`;
        
        logsContainer.insertBefore(logEntry, logsContainer.firstChild);
        
        // Keep only last 50 logs
        const logs = logsContainer.querySelectorAll('.log-entry');
        if (logs.length > 50) {
            logs[logs.length - 1].remove();
        }
    }

    formatDateTime(date) {
        return date.toLocaleString('id-ID', {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit'
        });
    }

    formatTime(date) {
        return date.toLocaleTimeString('id-ID', {
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit'
        });
    }
}

// Initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    new SyncMonitor();
});