import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { LogIn, Mail, Lock, Layers } from 'lucide-react';
import { authService } from '../services/api';

interface LoginPageProps {
  onLoginSuccess?: () => void;
}

const LoginPage = ({ onLoginSuccess }: LoginPageProps) => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const response = await authService.login({ username, password });
      localStorage.setItem('token', response.token);
      localStorage.setItem('user', JSON.stringify(response.user));
      
      window.dispatchEvent(new Event('tokenChanged'));
      
      if (onLoginSuccess) {
        onLoginSuccess();
      }
      
      navigate('/dashboard');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Login failed. Please check your credentials.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-slate-50 flex items-center justify-center p-4 font-sans selection:bg-blue-100 selection:text-blue-900">
      {/* Background Decorative Elements */}
      <div className="absolute top-0 w-full h-96 bg-blue-600 skew-y-[-4deg] origin-top-left -z-10 shadow-xl opacity-90" />
      <div className="absolute top-0 w-full h-96 bg-blue-700 skew-y-[-6deg] origin-top-left -z-20 opacity-50" />
      
      <div className="max-w-md w-full relative z-10">
        
        {/* Logo/Header */}
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center w-20 h-20 bg-white rounded-2xl mb-6 shadow-xl shadow-blue-900/20 ring-1 ring-slate-900/5">
            <Layers className="w-10 h-10 text-blue-600" />
          </div>
          <h1 className="text-4xl font-extrabold text-white tracking-tight mb-2">SES Dashboard</h1>
          <p className="text-blue-100 font-medium">Monitor your email infrastructure</p>
        </div>

        {/* Login Card */}
        <div className="bg-white rounded-3xl shadow-2xl border border-slate-100 p-8 sm:p-10 relative overflow-hidden">
          
          <h2 className="text-xl font-bold text-slate-800 mb-6">Welcome Back</h2>

          <form onSubmit={handleSubmit} className="space-y-5">
            {error && (
              <div className="bg-rose-50 border border-rose-100 text-rose-600 px-4 py-3 rounded-xl flex items-center text-sm font-medium">
                <div className="w-1.5 h-1.5 bg-rose-500 rounded-full mr-3 shrink-0" />
                {error}
              </div>
            )}

            <div>
              <label htmlFor="username" className="block text-sm font-semibold text-slate-700 mb-2">
                Username
              </label>
              <div className="relative group">
                <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none transition-colors">
                  <Mail className="h-5 w-5 text-slate-400 group-focus-within:text-blue-500 transition-colors" />
                </div>
                <input
                  id="username"
                  type="text"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  className="block w-full pl-12 pr-4 py-3.5 border border-slate-200 rounded-xl focus:ring-4 focus:ring-blue-500/10 focus:border-blue-500 transition-all bg-slate-50 focus:bg-white text-slate-900 font-medium placeholder:text-slate-400 placeholder:font-normal outline-none"
                  placeholder="admin"
                  required
                />
              </div>
            </div>

            <div>
              <label htmlFor="password" className="block text-sm font-semibold text-slate-700 mb-2">
                Password
              </label>
              <div className="relative group">
                <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                  <Lock className="h-5 w-5 text-slate-400 group-focus-within:text-blue-500 transition-colors" />
                </div>
                <input
                  id="password"
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="block w-full pl-12 pr-4 py-3.5 border border-slate-200 rounded-xl focus:ring-4 focus:ring-blue-500/10 focus:border-blue-500 transition-all bg-slate-50 focus:bg-white text-slate-900 font-medium placeholder:text-slate-400 placeholder:font-normal outline-none"
                  placeholder="••••••••"
                  required
                />
              </div>
            </div>

            <button
              type="submit"
              disabled={loading}
              className="w-full bg-blue-600 text-white py-3.5 px-4 rounded-xl hover:bg-blue-700 active:scale-[0.98] focus:ring-4 focus:ring-blue-500/20 disabled:opacity-70 disabled:cursor-not-allowed flex items-center justify-center font-bold transition-all shadow-lg shadow-blue-600/30 mt-4"
            >
              {loading ? (
                <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white" />
              ) : (
                <>
                  Sign In
                  <LogIn className="w-5 h-5 ml-2 opacity-80" />
                </>
              )}
            </button>
          </form>

          {/* Demo Hint */}
          <div className="mt-8 pt-6 border-t border-slate-100">
            <div className="bg-slate-50 rounded-xl p-4 border border-slate-100 flex items-center gap-3">
              <div className="bg-blue-100 p-2 rounded-lg shrink-0">
                <Lock className="w-4 h-4 text-blue-600" />
              </div>
              <div className="text-sm">
                <p className="text-slate-500">Default Demo Credentials</p>
                <p className="font-semibold text-slate-700">admin / password</p>
              </div>
            </div>
          </div>
          
        </div>
        
        {/* Footer */}
        <div className="text-center mt-10">
          <p className="text-sm text-slate-500 font-medium">
            © 2025 Wisnu. All rights reserved.
          </p>
        </div>

      </div>
    </div>
  );
};

export default LoginPage;
