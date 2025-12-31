import { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { LogOut, Menu, X, BarChart3, Mail, TrendingUp, Users, Settings, Cog, Shield } from 'lucide-react';
import { authService } from '../services/api';
import ChangePassword from './ChangePassword';

interface LayoutProps {
  children: React.ReactNode;
  title: string;
}

const Layout = ({ children, title }: LayoutProps) => {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [desktopSidebarOpen, setDesktopSidebarOpen] = useState(true);
  const [showChangePassword, setShowChangePassword] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();

  const handleLogout = () => {
    authService.logout();
    navigate('/login');
  };

  const menuItems = [
    { name: 'Dashboard', icon: BarChart3, path: '/dashboard' },
    { name: 'Events', icon: Mail, path: '/events' },
    { name: 'Analytics', icon: TrendingUp, path: '/analytics' },
  ];

  const adminMenuItems = [
    { name: 'Users', icon: Users, path: '/admin/users' },
    { name: 'Suppression List', icon: Shield, path: '/admin/suppression' },
    { name: 'Settings', icon: Cog, path: '/admin/settings' },
  ];

  const isCurrentPath = (path: string) => location.pathname === path;
  const currentUser = authService.getCurrentUser();
  const isAdmin = authService.isAdmin();

  return (
    <div className="min-h-screen bg-gray-50 flex">
      {/* Mobile sidebar overlay */}
      {sidebarOpen && (
        <div className="fixed inset-0 z-40 lg:hidden">
          <div 
            className="fixed inset-0 bg-black bg-opacity-50" 
            onClick={() => setSidebarOpen(false)} 
          />
        </div>
      )}

      {/* Desktop Sidebar */}
      {desktopSidebarOpen && (
        <div className="hidden lg:flex lg:w-64 lg:flex-col lg:fixed lg:inset-y-0 z-40 bg-white shadow-xl">
          {/* Sidebar Header */}
          <div className="flex items-center h-16 px-6 bg-gradient-to-r from-blue-600 to-blue-700 border-b border-blue-800">
            <h1 className="text-xl font-bold text-white tracking-wide">SES Dashboard</h1>
          </div>

          {/* Navigation Menu */}
          <nav className="flex-1 px-4 py-6 space-y-1">
            {menuItems.map((item) => {
              const isActive = isCurrentPath(item.path);
              return (
                <button
                  key={item.name}
                  onClick={() => navigate(item.path)}
                  className={`w-full flex items-center px-4 py-3 text-sm font-medium rounded-xl transition-all duration-200 group ${
                    isActive
                      ? 'bg-blue-50 text-blue-700 shadow-sm border-l-4 border-blue-600'
                      : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900 hover:shadow-sm'
                  }`}
                >
                  <item.icon className={`w-5 h-5 mr-3 transition-colors duration-200 ${
                    isActive ? 'text-blue-600' : 'text-gray-400 group-hover:text-gray-600'
                  }`} />
                  <span className="font-medium">{item.name}</span>
                </button>
              );
            })}
            
            {/* Admin Menu */}
            {isAdmin && (
              <>
                <div className="pt-4 pb-2">
                  <div className="px-4 text-xs font-semibold text-gray-500 uppercase tracking-wider">
                    Administration
                  </div>
                </div>
                {adminMenuItems.map((item) => {
                  const isActive = isCurrentPath(item.path);
                  return (
                    <button
                      key={item.name}
                      onClick={() => navigate(item.path)}
                      className={`w-full flex items-center px-4 py-3 text-sm font-medium rounded-xl transition-all duration-200 group ${
                        isActive
                          ? 'bg-blue-50 text-blue-700 shadow-sm border-l-4 border-blue-600'
                          : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900 hover:shadow-sm'
                      }`}
                    >
                      <item.icon className={`w-5 h-5 mr-3 transition-colors duration-200 ${
                        isActive ? 'text-blue-600' : 'text-gray-400 group-hover:text-gray-600'
                      }`} />
                      <span className="font-medium">{item.name}</span>
                    </button>
                  );
                })}
              </>
            )}
          </nav>

          {/* Logout Button */}
          <div className="px-4 py-4 border-t border-gray-100">
            <button
              onClick={() => setShowChangePassword(true)}
              className="w-full flex items-center px-4 py-3 text-sm font-medium text-gray-600 rounded-xl hover:bg-blue-50 hover:text-blue-700 transition-all duration-200 group mb-2"
            >
              <Settings className="w-5 h-5 mr-3 text-gray-400 group-hover:text-blue-600 transition-colors duration-200" />
              <span className="font-medium">Change Password</span>
            </button>
            <button
              onClick={handleLogout}
              className="w-full flex items-center px-4 py-3 text-sm font-medium text-gray-600 rounded-xl hover:bg-red-50 hover:text-red-700 transition-all duration-200 group"
            >
              <LogOut className="w-5 h-5 mr-3 text-gray-400 group-hover:text-red-600 transition-colors duration-200" />
              <span className="font-medium">Logout</span>
            </button>
          </div>

          {/* Copyright */}
          <div className="px-4 py-3 border-t border-gray-100 bg-gray-50">
            <p className="text-xs text-gray-500 text-center">
              © 2025 Wisnu. All rights reserved.
            </p>
          </div>
        </div>
      )}

      {/* Mobile Sidebar */}
      <div className={`fixed inset-y-0 left-0 z-50 w-64 bg-white shadow-xl transform ${
        sidebarOpen ? 'translate-x-0' : '-translate-x-full'
      } transition-transform duration-300 ease-in-out lg:hidden flex flex-col`}>
        
        {/* Sidebar Header */}
        <div className="flex items-center justify-between h-16 px-6 bg-gradient-to-r from-blue-600 to-blue-700 border-b border-blue-800">
          <h1 className="text-xl font-bold text-white tracking-wide">SES Dashboard</h1>
          <button
            onClick={() => setSidebarOpen(false)}
            className="lg:hidden p-1.5 rounded-md text-blue-200 hover:text-white hover:bg-blue-800 transition-colors duration-200"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Navigation Menu */}
        <nav className="flex-1 px-4 py-6 space-y-1">
          {menuItems.map((item) => {
            const isActive = isCurrentPath(item.path);
            return (
              <button
                key={item.name}
                onClick={() => {
                  navigate(item.path);
                  setSidebarOpen(false);
                }}
                className={`w-full flex items-center px-4 py-3 text-sm font-medium rounded-xl transition-all duration-200 group ${
                  isActive
                    ? 'bg-blue-50 text-blue-700 shadow-sm border-l-4 border-blue-600'
                    : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900 hover:shadow-sm'
                }`}
              >
                <item.icon className={`w-5 h-5 mr-3 transition-colors duration-200 ${
                  isActive ? 'text-blue-600' : 'text-gray-400 group-hover:text-gray-600'
                }`} />
                <span className="font-medium">{item.name}</span>
              </button>
            );
          })}
          
          {/* Admin Menu */}
          {isAdmin && (
            <>
              <div className="pt-4 pb-2">
                <div className="px-4 text-xs font-semibold text-gray-500 uppercase tracking-wider">
                  Administration
                </div>
              </div>
              {adminMenuItems.map((item) => {
                const isActive = isCurrentPath(item.path);
                return (
                  <button
                    key={item.name}
                    onClick={() => {
                      navigate(item.path);
                      setSidebarOpen(false);
                    }}
                    className={`w-full flex items-center px-4 py-3 text-sm font-medium rounded-xl transition-all duration-200 group ${
                      isActive
                        ? 'bg-blue-50 text-blue-700 shadow-sm border-l-4 border-blue-600'
                        : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900 hover:shadow-sm'
                    }`}
                  >
                    <item.icon className={`w-5 h-5 mr-3 transition-colors duration-200 ${
                      isActive ? 'text-blue-600' : 'text-gray-400 group-hover:text-gray-600'
                    }`} />
                    <span className="font-medium">{item.name}</span>
                  </button>
                );
              })}
            </>
          )}
        </nav>

        {/* Logout Button */}
        <div className="px-4 py-4 border-t border-gray-100">
          <button
            onClick={() => setShowChangePassword(true)}
            className="w-full flex items-center px-4 py-3 text-sm font-medium text-gray-600 rounded-xl hover:bg-blue-50 hover:text-blue-700 transition-all duration-200 group mb-2"
          >
            <Settings className="w-5 h-5 mr-3 text-gray-400 group-hover:text-blue-600 transition-colors duration-200" />
            <span className="font-medium">Change Password</span>
          </button>
          <button
            onClick={handleLogout}
            className="w-full flex items-center px-4 py-3 text-sm font-medium text-gray-600 rounded-xl hover:bg-red-50 hover:text-red-700 transition-all duration-200 group"
          >
            <LogOut className="w-5 h-5 mr-3 text-gray-400 group-hover:text-red-600 transition-colors duration-200" />
            <span className="font-medium">Logout</span>
          </button>
        </div>

        {/* Copyright */}
        <div className="px-4 py-3 border-t border-gray-100 bg-gray-50">
          <p className="text-xs text-gray-500 text-center">
            © 2025 Wisnu. All rights reserved.
          </p>
        </div>
      </div>

      {/* Main content */}
      <div className={`flex-1 flex flex-col ${desktopSidebarOpen ? 'lg:ml-64' : 'lg:ml-0'} transition-all duration-300`}>
        {/* Top bar */}
        <div className="sticky top-0 z-50 bg-white/95 backdrop-blur-sm border-b border-gray-200 shadow-sm">
          <div className="flex items-center justify-between h-16 px-4 sm:px-6 lg:px-8">
            <div className="flex items-center">
              <button
                onClick={() => setSidebarOpen(true)}
                className="lg:hidden p-2 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-all duration-200"
              >
                <Menu className="w-6 h-6" />
              </button>
              <button
                onClick={() => setDesktopSidebarOpen(!desktopSidebarOpen)}
                className="hidden lg:block p-2 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-all duration-200 mr-3"
                title={desktopSidebarOpen ? 'Hide Menu' : 'Show Menu'}
              >
                <Menu className="w-5 h-5" />
              </button>
              <div className="ml-4 lg:ml-0">
                <h2 className="text-xl font-bold text-gray-900">{title}</h2>
                <div className="h-0.5 w-8 bg-blue-600 rounded-full mt-1"></div>
              </div>
            </div>

            <div className="flex items-center space-x-4">
              <div className="hidden sm:block text-sm text-gray-500">
                Welcome back!
              </div>
              <button
                onClick={handleLogout}
                className="hidden lg:flex items-center px-4 py-2 text-sm font-medium text-gray-600 bg-gray-100 rounded-lg hover:bg-red-50 hover:text-red-700 transition-all duration-200 group"
              >
                <LogOut className="w-4 h-4 mr-2 group-hover:text-red-600 transition-colors duration-200" />
                Logout
              </button>
            </div>
          </div>
        </div>

        {/* Page content */}
        <main className="flex-1 p-6 lg:p-8">
          <div className="max-w-7xl mx-auto">
            {children}
            
            {/* Footer */}
            <div className="text-center py-6 mt-12 border-t border-gray-200">
              <p className="text-sm text-gray-500">
                © 2025 Wisnu. All rights reserved.
              </p>
            </div>
          </div>
        </main>
      </div>
      
      {/* Change Password Modal */}
      {showChangePassword && (
        <ChangePassword onClose={() => setShowChangePassword(false)} />
      )}
    </div>
  );
};

export default Layout;