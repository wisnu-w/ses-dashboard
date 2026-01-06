import { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import LoginPage from './pages/LoginPage';
import DashboardPage from './pages/DashboardPage';
import EventsPage from './pages/EventsPage';
import AnalyticsPage from './pages/AnalyticsPage';
import UsersPage from './pages/UsersPage';
import SettingsPage from './pages/SettingsPage';
import SuppressionPage from './pages/SuppressionPage';
import { authService } from './services/api';

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean | null>(null);
  const [loading, setLoading] = useState(true);

  const checkAuth = () => {
    try {
      const authenticated = authService.isAuthenticated();
      setIsAuthenticated(authenticated);
      setLoading(false);
    } catch (error) {
      console.error('Auth check error:', error);
      setIsAuthenticated(false);
      setLoading(false);
    }
  };

  useEffect(() => {
    checkAuth();

    // Listen for storage changes (when token is added/removed)
    const handleStorageChange = () => {
      checkAuth();
    };

    window.addEventListener('storage', handleStorageChange);
    
    // Custom event for same-tab token changes
    window.addEventListener('tokenChanged', handleStorageChange);

    return () => {
      window.removeEventListener('storage', handleStorageChange);
      window.removeEventListener('tokenChanged', handleStorageChange);
    };
  }, []);

  if (loading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <Router>
      <Routes>
        <Route
          path="/login"
          element={
            isAuthenticated ? <Navigate to="/dashboard" replace /> : <LoginPage onLoginSuccess={checkAuth} />
          }
        />
        <Route
          path="/dashboard"
          element={
            isAuthenticated ? <DashboardPage /> : <Navigate to="/login" replace />
          }
        />
        <Route
          path="/events"
          element={
            isAuthenticated ? <EventsPage /> : <Navigate to="/login" replace />
          }
        />
        <Route
          path="/analytics"
          element={
            isAuthenticated ? <AnalyticsPage /> : <Navigate to="/login" replace />
          }
        />
        <Route
          path="/admin/users"
          element={
            isAuthenticated && authService.isAdmin() ? <UsersPage /> : <Navigate to="/dashboard" replace />
          }
        />
        <Route
          path="/admin/suppression"
          element={
            isAuthenticated && authService.isAdmin() ? <SuppressionPage /> : <Navigate to="/dashboard" replace />
          }
        />
        <Route
          path="/admin/settings"
          element={
            isAuthenticated && authService.isAdmin() ? <SettingsPage /> : <Navigate to="/dashboard" replace />
          }
        />
        <Route
          path="/"
          element={<Navigate to={isAuthenticated ? "/dashboard" : "/login"} replace />}
        />
      </Routes>
    </Router>
  );
}

export default App;
