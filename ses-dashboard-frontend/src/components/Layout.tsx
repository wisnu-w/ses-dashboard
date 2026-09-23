import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { LayoutDashboard, Mail, Search, Menu, X, LogOut, FileText, Ban, Settings, Users, Layers } from 'lucide-react';
import { authService } from '../services/api';
import ChangePassword from './ChangePassword';

interface LayoutProps {
  children: React.ReactNode;
  title?: string;
}

const Layout = ({ children, title = 'SES Dashboard' }: LayoutProps) => {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [desktopSidebarOpen, setDesktopSidebarOpen] = useState(true);
  const [showChangePassword, setShowChangePassword] = useState(false);
  
  const user = authService.getCurrentUser();
  const navigate = useNavigate();
  const location = useLocation();

  const handleLogout = () => {
    authService.logout();
    navigate('/login');
  };

  useEffect(() => {
    if (window.innerWidth < 1024) {
      setDesktopSidebarOpen(false);
    }
  }, []);

  const menuItems = [
    { name: 'Dashboard', icon: LayoutDashboard, path: '/dashboard' },
    { name: 'Event Logs', icon: Mail, path: '/events' },
    { name: 'Analytics', icon: Search, path: '/analytics' },
    { name: 'Suppression List', icon: Ban, path: '/suppression' },
  ];

  const adminMenuItems = [
    { name: 'Users', icon: Users, path: '/users' },
    { name: 'Settings', icon: Settings, path: '/settings' },
  ];

  const isAdmin = user?.role === 'admin';
  const isCurrentPath = (path: string) => {
    if (path === '/' && location.pathname !== '/') return false;
    return location.pathname.startsWith(path);
  };

  const NavItem = ({ item }: { item: typeof menuItems[0] }) => {
    const active = isCurrentPath(item.path);
    return (
      <button
        onClick={() => {
          navigate(item.path);
          setSidebarOpen(false);
        }}
        className={`w-full flex items-center px-4 py-3 text-sm font-medium rounded-xl transition-all duration-200 group ${
          active
            ? 'bg-blue-50 text-blue-700 font-semibold shadow-sm'
            : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'
        }`}
      >
        <item.icon className={`w-5 h-5 mr-3 transition-colors ${active ? 'text-blue-600' : 'text-slate-400 group-hover:text-slate-600'}`} />
        {item.name}
      </button>
    );
  };

  return (
    <div className="h-screen bg-slate-50 flex overflow-hidden selection:bg-blue-100 selection:text-blue-900 font-sans">
      
      {/* Mobile Backdrop */}
      {sidebarOpen && (
        <div 
          className="fixed inset-0 z-40 bg-slate-900/40 backdrop-blur-sm lg:hidden transition-opacity" 
          onClick={() => setSidebarOpen(false)}
        />
      )}

      {/* Sidebar - Fixed Height without Window Scroll */}
      <aside className={`fixed lg:static inset-y-0 left-0 z-50 w-72 bg-white border-r border-slate-200 flex flex-col h-screen transform transition-transform duration-300 ease-in-out ${
        sidebarOpen || desktopSidebarOpen ? 'translate-x-0' : '-translate-x-full lg:-ml-72 lg:w-0 lg:overflow-hidden'
      }`}>
        <div className="flex items-center h-20 px-8 border-b border-slate-100 shrink-0">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center shadow-sm shadow-blue-200">
              <Layers className="w-5 h-5 text-white" />
            </div>
            <h1 className="text-xl font-bold tracking-tight text-slate-800">SES<span className="text-blue-600 font-medium">Dash</span></h1>
          </div>
          <button onClick={() => setSidebarOpen(false)} className="lg:hidden ml-auto p-2 text-slate-400 hover:bg-slate-100 rounded-lg">
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Scrollable Nav Area */}
        <div className="flex-1 overflow-y-auto px-4 py-6 space-y-8 custom-scrollbar">
          <div>
            <p className="px-4 text-[10px] font-bold text-slate-400 uppercase tracking-wider mb-3">Main Menu</p>
            <nav className="space-y-1">
              {menuItems.map(item => <NavItem key={item.name} item={item} />)}
            </nav>
          </div>

          {isAdmin && (
            <div>
              <p className="px-4 text-[10px] font-bold text-slate-400 uppercase tracking-wider mb-3">Administration</p>
              <nav className="space-y-1">
                {adminMenuItems.map(item => <NavItem key={item.name} item={item} />)}
              </nav>
            </div>
          )}
        </div>

        {/* Fixed Bottom Footer */}
        <div className="p-4 border-t border-slate-100 bg-slate-50/50 shrink-0">
          <div className="flex items-center px-4 py-3 mb-3 rounded-xl bg-white border border-slate-200 shadow-sm">
            <div className="w-9 h-9 rounded-full bg-blue-50 text-blue-600 border border-blue-100 flex items-center justify-center font-bold text-sm">
              {user?.username?.charAt(0).toUpperCase() || 'A'}
            </div>
            <div className="ml-3 truncate">
              <p className="text-sm font-semibold text-slate-800 truncate">{user?.username || 'Admin'}</p>
              <p className="text-xs font-medium text-slate-500 capitalize truncate">{user?.role || 'User'}</p>
            </div>
          </div>
          
          <button onClick={() => setShowChangePassword(true)} className="w-full flex items-center px-4 py-2.5 text-sm font-medium text-slate-600 rounded-lg hover:bg-white hover:shadow-sm border border-transparent hover:border-slate-200 transition-all group">
            <Settings className="w-4 h-4 mr-3 text-slate-400 group-hover:text-slate-600" />
            Password
          </button>
          <button onClick={handleLogout} className="w-full flex items-center px-4 py-2.5 text-sm font-medium text-rose-600 rounded-lg hover:bg-rose-50 hover:shadow-sm border border-transparent hover:border-rose-100 transition-all group mt-1">
            <LogOut className="w-4 h-4 mr-3 text-rose-400 group-hover:text-rose-600" />
            Logout
          </button>
          
          <div className="mt-4 text-center">
            <p className="text-[10px] text-slate-400 font-medium">© 2025 Wisnu. All rights reserved.</p>
          </div>
        </div>
      </aside>

      {/* Main Content */}
      <div className="flex-1 flex flex-col min-w-0 h-screen overflow-hidden">
        {/* Header */}
        <header className="shrink-0 sticky top-0 z-30 bg-white/80 backdrop-blur-md border-b border-slate-200">
          <div className="flex items-center justify-between h-20 px-6 lg:px-10">
            <div className="flex items-center gap-4">
              <button 
                onClick={() => {
                  if (window.innerWidth < 1024) setSidebarOpen(true);
                  else setDesktopSidebarOpen(!desktopSidebarOpen);
                }}
                className="p-2 -ml-2 text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-xl transition-colors"
              >
                <Menu className="w-6 h-6" />
              </button>
              <h2 className="text-2xl font-bold text-slate-800 tracking-tight">{title}</h2>
            </div>
            <div className="hidden sm:flex items-center text-sm font-medium text-slate-500">
              {new Date().toLocaleDateString('en-US', { weekday: 'long', month: 'long', day: 'numeric' })}
            </div>
          </div>
        </header>

        {/* Page Container - Scrollable */}
        <main className="flex-1 overflow-y-auto p-6 lg:p-10 custom-scrollbar">
          <div className="max-w-7xl mx-auto">
            {children}
          </div>
        </main>
      </div>

      {showChangePassword && <ChangePassword onClose={() => setShowChangePassword(false)} />}
    </div>
  );
};

export default Layout;
