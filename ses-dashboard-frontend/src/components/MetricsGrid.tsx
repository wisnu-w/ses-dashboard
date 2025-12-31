import { TrendingUp, TrendingDown, Mail, Send, CheckCircle, XCircle, Eye, MousePointer } from 'lucide-react';

interface MetricCardProps {
  title: string;
  value: number | string;
  icon: React.ReactNode;
  trend?: {
    value: number;
    isPositive: boolean;
  };
  color: string;
}

const MetricCard = ({ title, value, icon, trend, color }: MetricCardProps) => {
  return (
    <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm font-medium text-gray-600">{title}</p>
          <p className="text-2xl font-bold text-gray-900 mt-1">{value}</p>
          {trend && (
            <div className={`flex items-center mt-2 text-sm ${trend.isPositive ? 'text-green-600' : 'text-red-600'}`}>
              {trend.isPositive ? (
                <TrendingUp className="w-4 h-4 mr-1" />
              ) : (
                <TrendingDown className="w-4 h-4 mr-1" />
              )}
              <span>{Math.abs(trend.value)}%</span>
            </div>
          )}
        </div>
        <div className={`p-3 rounded-full ${color}`}>
          {icon}
        </div>
      </div>
    </div>
  );
};

interface MetricsGridProps {
  metrics: {
    total_events: number;
    send_count: number;
    delivery_count: number;
    bounce_count: number;
    complaint_count: number;
    open_count: number;
    click_count: number;
    bounce_rate: number;
    delivery_rate: number;
  };
}

const MetricsGrid = ({ metrics }: MetricsGridProps) => {
  const cards = [
    {
      title: 'Total Events',
      value: metrics.total_events.toLocaleString(),
      icon: <Mail className="w-6 h-6 text-white" />,
      color: 'bg-blue-500',
    },
    {
      title: 'Sent',
      value: metrics.send_count.toLocaleString(),
      icon: <Send className="w-6 h-6 text-white" />,
      color: 'bg-green-500',
    },
    {
      title: 'Delivered',
      value: metrics.delivery_count.toLocaleString(),
      icon: <CheckCircle className="w-6 h-6 text-white" />,
      color: 'bg-emerald-500',
    },
    {
      title: 'Bounced',
      value: metrics.bounce_count.toLocaleString(),
      icon: <XCircle className="w-6 h-6 text-white" />,
      color: 'bg-red-500',
    },
    {
      title: 'Complaints',
      value: metrics.complaint_count.toLocaleString(),
      icon: <XCircle className="w-6 h-6 text-white" />,
      color: 'bg-orange-500',
    },
    {
      title: 'Opens',
      value: metrics.open_count.toLocaleString(),
      icon: <Eye className="w-6 h-6 text-white" />,
      color: 'bg-purple-500',
    },
    {
      title: 'Clicks',
      value: metrics.click_count.toLocaleString(),
      icon: <MousePointer className="w-6 h-6 text-white" />,
      color: 'bg-indigo-500',
    },
    {
      title: 'Delivery Rate',
      value: `${metrics.delivery_rate.toFixed(1)}%`,
      icon: <TrendingUp className="w-6 h-6 text-white" />,
      color: 'bg-teal-500',
    },
  ];

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 items-start">
      {cards.map((card, index) => (
        <MetricCard key={index} {...card} />
      ))}
    </div>
  );
};

export default MetricsGrid;