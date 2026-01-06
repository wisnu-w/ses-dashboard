import { useState, useEffect } from 'react';
import { TrendingUp, Calendar, BarChart3 } from 'lucide-react';
import Layout from '../components/Layout';
import { LineChartComponent, BarChartComponent } from '../components/Charts';
import { metricsService } from '../services/api';
import type { MonthlyMetrics, HourlyMetrics } from '../types/api';

const AnalyticsPage = () => {
  const [monthlyMetrics, setMonthlyMetrics] = useState<MonthlyMetrics[]>([]);
  const [hourlyMetrics, setHourlyMetrics] = useState<HourlyMetrics[]>([]);
  const [loading, setLoading] = useState(true);

  const parseHour = (value: string): number | null => {
    const parts = value.split(' ');
    const timePart = parts.length > 1 ? parts[1] : value;
    const hourPart = timePart.split(':')[0];
    const hourNumber = Number(hourPart);
    return Number.isNaN(hourNumber) ? null : hourNumber;
  };

  const loadData = async () => {
    try {
      setLoading(true);
      const [monthlyData, hourlyData] = await Promise.all([
        metricsService.getMonthlyMetrics(),
        metricsService.getHourlyMetrics(),
      ]);

      setMonthlyMetrics(monthlyData.monthly_metrics);
      setHourlyMetrics(hourlyData.hourly_metrics);
    } catch (error) {
      console.error('Failed to load analytics data:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  if (loading) {
    return (
      <Layout title="Analytics">
        <div className="animate-pulse space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {[...Array(3)].map((_, i) => (
              <div key={i} className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
                <div className="h-4 bg-gray-200 rounded w-1/2 mb-4"></div>
                <div className="h-8 bg-gray-200 rounded w-1/3"></div>
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

  const totalMonthlyEmails = monthlyMetrics.reduce((sum, metric) => sum + metric.send_count, 0);
  const maxMonthlySend = Math.max(...monthlyMetrics.map((metric) => metric.send_count), 0);
  const maxHourlySend = Math.max(...hourlyMetrics.map((metric) => metric.send_count), 0);
  const avgDeliveryRate = monthlyMetrics.length > 0 
    ? monthlyMetrics.reduce((sum, metric) => sum + (metric.delivery_count / Math.max(metric.send_count, 1)), 0) / monthlyMetrics.length * 100
    : 0;
  const peakHourMetric = hourlyMetrics.reduce(
    (max, metric) => (metric.send_count > max.send_count ? metric : max),
    hourlyMetrics[0] || { hour: '', send_count: 0 }
  );
  const peakHourValue = parseHour(peakHourMetric.hour);
  const peakHourLabel = peakHourValue === null ? peakHourMetric.hour : `${peakHourValue}:00`;
  const sortedHourlyMetrics = [...hourlyMetrics].sort((a, b) => b.send_count - a.send_count);

  return (
    <Layout title="Analytics">
      <div className="space-y-6">
        {/* Header */}
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Advanced Analytics</h1>
          <p className="text-gray-600 mt-1">Deep insights into your email performance</p>
        </div>

        {/* Summary Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div className="flex items-center">
              <div className="p-3 rounded-full bg-blue-100">
                <BarChart3 className="w-6 h-6 text-blue-600" />
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium text-gray-600">Total Monthly Emails</p>
                <p className="text-2xl font-bold text-gray-900">{totalMonthlyEmails.toLocaleString()}</p>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div className="flex items-center">
              <div className="p-3 rounded-full bg-green-100">
                <TrendingUp className="w-6 h-6 text-green-600" />
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium text-gray-600">Avg Delivery Rate</p>
                <p className="text-2xl font-bold text-gray-900">{avgDeliveryRate.toFixed(1)}%</p>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div className="flex items-center">
              <div className="p-3 rounded-full bg-purple-100">
                <Calendar className="w-6 h-6 text-purple-600" />
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium text-gray-600">Peak Hour</p>
                <p className="text-2xl font-bold text-gray-900">{peakHourLabel}</p>
              </div>
            </div>
          </div>
        </div>

        {/* Charts */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <LineChartComponent
            data={monthlyMetrics.slice(-12)}
            title="Monthly Email Trends"
            dataKey="send_count"
            xAxisKey="month"
          />

          <BarChartComponent
            data={hourlyMetrics}
            title="Hourly Distribution"
            dataKey="send_count"
            xAxisKey="hour"
          />
        </div>

        {/* Detailed Analytics */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Monthly Performance</h3>
            <div className="space-y-4">
              {monthlyMetrics.slice(-6).map((metric, index) => (
                <div key={index} className="flex items-center justify-between">
                  <span className="text-sm text-gray-600">{metric.month}</span>
                  <div className="flex items-center space-x-4">
                    <span className="text-sm font-medium text-gray-900">
                      {metric.send_count.toLocaleString()}
                    </span>
                    <div className="w-20 bg-gray-200 rounded-full h-2">
                      <div 
                        className="bg-blue-600 h-2 rounded-full" 
                        style={{ 
                          width: `${maxMonthlySend > 0 ? Math.min((metric.send_count / maxMonthlySend) * 100, 100) : 0}%` 
                        }}
                      ></div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>

          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Peak Hours Analysis</h3>
            <div className="space-y-4">
              {sortedHourlyMetrics
                .slice(0, 6)
                .map((metric, index) => {
                  const hourValue = parseHour(metric.hour);
                  const hourLabel = hourValue === null
                    ? metric.hour
                    : `${hourValue}:00 - ${(hourValue + 1) % 24}:00`;
                  return (
                <div key={index} className="flex items-center justify-between">
                  <span className="text-sm text-gray-600">{hourLabel}</span>
                  <div className="flex items-center space-x-4">
                    <span className="text-sm font-medium text-gray-900">
                      {metric.send_count.toLocaleString()}
                    </span>
                    <div className="w-20 bg-gray-200 rounded-full h-2">
                      <div 
                        className="bg-purple-600 h-2 rounded-full" 
                        style={{ 
                          width: `${maxHourlySend > 0 ? Math.min((metric.send_count / maxHourlySend) * 100, 100) : 0}%` 
                        }}
                      ></div>
                    </div>
                  </div>
                </div>
              )})}
            </div>
          </div>
        </div>
      </div>
    </Layout>
  );
};

export default AnalyticsPage;
