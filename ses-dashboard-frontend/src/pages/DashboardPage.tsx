import { useState, useEffect } from 'react';
import { RefreshCw } from 'lucide-react';
import Layout from '../components/Layout';
import MetricsGrid from '../components/MetricsGrid';
import { LineChartComponent, BarChartComponent } from '../components/Charts';
import { metricsService } from '../services/api';
import type{ MetricsResponse, DailyMetrics } from '../types/api';

const DashboardPage = () => {
  const [metrics, setMetrics] = useState<MetricsResponse | null>(null);
  const [dailyMetrics, setDailyMetrics] = useState<DailyMetrics[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);

  const loadData = async (showRefresh = false) => {
    try {
      if (showRefresh) setRefreshing(true);
      else setLoading(true);

      const [metricsData, dailyData] = await Promise.all([
        metricsService.getMetrics(),
        metricsService.getDailyMetrics(),
      ]);

      setMetrics(metricsData);
      setDailyMetrics(dailyData.daily_metrics);
    } catch (error) {
      console.error('Failed to load dashboard data:', error);
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  const handleRefresh = () => {
    loadData(true);
  };

  if (loading) {
    return (
      <Layout title="Dashboard">
        <div className="animate-pulse">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
            {[...Array(8)].map((_, i) => (
              <div key={i} className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
                <div className="h-4 bg-gray-200 rounded w-1/4 mb-2"></div>
                <div className="h-8 bg-gray-200 rounded w-1/2 mb-4"></div>
                <div className="h-4 bg-gray-200 rounded w-1/3"></div>
              </div>
            ))}
          </div>
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {[...Array(2)].map((_, i) => (
              <div key={i} className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
                <div className="h-4 bg-gray-200 rounded w-1/4 mb-4"></div>
                <div className="h-64 bg-gray-200 rounded"></div>
              </div>
            ))}
          </div>
        </div>
      </Layout>
    );
  }

  return (
    <Layout title="Dashboard">
      <div className="space-y-6">
        {/* Header with refresh button */}
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Email Analytics Overview</h1>
            <p className="text-gray-600 mt-1">Monitor your SES email delivery performance</p>
          </div>
          <button
            onClick={handleRefresh}
            disabled={refreshing}
            className="inline-flex items-center px-4 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <RefreshCw className={`w-4 h-4 mr-2 ${refreshing ? 'animate-spin' : ''}`} />
            Refresh
          </button>
        </div>

        {/* Metrics Grid */}
        {metrics && <MetricsGrid metrics={metrics} />}

        {/* Charts */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 items-start">
          <LineChartComponent
            data={dailyMetrics.slice(-7)} // Last 7 days
            title="Daily Send Volume"
            dataKey="send_count"
            xAxisKey="date"
          />

          <BarChartComponent
            data={dailyMetrics.slice(-7)} // Last 7 days
            title="Daily Delivery Performance"
            dataKey="delivery_count"
            xAxisKey="date"
          />
        </div>

        {/* Additional Insights */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 items-start">
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Delivery Rate Trend</h3>
            <div className="text-3xl font-bold text-green-600">
              {metrics ? `${metrics.delivery_rate.toFixed(1)}%` : '0%'}
            </div>
            <p className="text-gray-600 mt-2">Emails successfully delivered</p>
          </div>

          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Bounce Rate</h3>
            <div className="text-3xl font-bold text-red-600">
              {metrics ? `${metrics.bounce_rate.toFixed(1)}%` : '0%'}
            </div>
            <p className="text-gray-600 mt-2">Emails that bounced</p>
          </div>

          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Engagement</h3>
            <div className="text-3xl font-bold text-purple-600">
              {metrics ? `${((metrics.open_count + metrics.click_count) / Math.max(metrics.delivery_count, 1) * 100).toFixed(1)}%` : '0%'}
            </div>
            <p className="text-gray-600 mt-2">Opens + clicks rate</p>
          </div>
        </div>
      </div>
    </Layout>
  );
};

export default DashboardPage;