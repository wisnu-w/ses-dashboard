import axios from 'axios';
import type {
  LoginRequest,
  LoginResponse,
  EventsResponse,
  MetricsResponse,
  DailyMetrics,
  MonthlyMetrics,
  HourlyMetrics,
  User,
  CreateUserRequest
} from '../types/api';


const API_BASE_URL = import.meta.env.VITE_API_URL || '';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add auth token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor to handle 401 errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const authService = {
  login: async (credentials: LoginRequest): Promise<LoginResponse> => {
    const response = await api.post('/api/login', credentials);
    return response.data;
  },

  logout: () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    // Trigger custom event for same-tab token changes
    window.dispatchEvent(new Event('tokenChanged'));
  },

  isAuthenticated: (): boolean => {
    return !!localStorage.getItem('token');
  },

  getCurrentUser: (): LoginResponse['user'] | null => {
    const userStr = localStorage.getItem('user');
    return userStr ? JSON.parse(userStr) : null;
  },

  isAdmin: (): boolean => {
    const user = authService.getCurrentUser();
    return user?.role === 'admin';
  },
};

export const eventsService = {
  getEvents: async (page = 1, limit = 50, search = '', startDate = '', endDate = ''): Promise<EventsResponse> => {
    try {
      const params = new URLSearchParams({
        page: page.toString(),
        limit: limit.toString()
      });
      
      if (search) params.append('search', search);
      if (startDate) params.append('start_date', startDate);
      if (endDate) params.append('end_date', endDate);
      
      const response = await api.get(`/api/events?${params}`);
      return response.data;
    } catch (error: any) {
      console.error('Failed to fetch events:', error);
      // Return empty data structure if API fails
      return {
        events: [],
        pagination: {
          page: 1,
          limit: 50,
          total: 0,
          totalPages: 0,
          hasNext: false,
          hasPrev: false
        }
      };
    }
  },
};

export const metricsService = {
  getMetrics: async (): Promise<MetricsResponse> => {
    const response = await api.get('/api/metrics');
    return response.data;
  },

  getDailyMetrics: async (): Promise<{ daily_metrics: DailyMetrics[] }> => {
    const response = await api.get('/api/metrics/daily');
    return response.data;
  },

  getMonthlyMetrics: async (): Promise<{ monthly_metrics: MonthlyMetrics[] }> => {
    const response = await api.get('/api/metrics/monthly');
    return response.data;
  },

  getHourlyMetrics: async (): Promise<{ hourly_metrics: HourlyMetrics[] }> => {
    const response = await api.get('/api/metrics/hourly');
    return response.data;
  },
};

export const userService = {
  createUser: async (userData: CreateUserRequest): Promise<void> => {
    await api.post('/api/users', userData);
  },

  getUsers: async (): Promise<{ users: User[] }> => {
    const response = await api.get('/api/users');
    return response.data;
  },

  resetPassword: async (userId: number, newPassword: string): Promise<void> => {
    await api.put(`/api/users/${userId}/reset-password`, { new_password: newPassword });
  },

  disableUser: async (userId: number): Promise<void> => {
    await api.put(`/api/users/${userId}/disable`);
  },

  enableUser: async (userId: number): Promise<void> => {
    await api.put(`/api/users/${userId}/enable`);
  },

  deleteUser: async (userId: number): Promise<void> => {
    await api.delete(`/api/users/${userId}`);
  },

  changePassword: async (oldPassword: string, newPassword: string): Promise<void> => {
    await api.put('/api/change-password', {
      old_password: oldPassword,
      new_password: newPassword
    });
  },
};

export default api;