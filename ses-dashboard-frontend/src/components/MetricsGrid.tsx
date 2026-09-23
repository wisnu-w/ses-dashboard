import { TrendingUp, TrendingDown, Mail, Send, CheckCircle, XCircle, Eye, MousePointer, AlertCircle } from 'lucide-react';

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
    <div className="bg-white rounded-2xl shadow-sm border border-slate-200 p-6 hover:shadow-md transition-shadow duration-200 relative overflow-hidden group">
      {/* Subtle background glow effect */}
      <div className={`absolute -right-4 -top-4 w-24 h-24 rounded-full opacity-10 transition-transform group-hover:scale-150 duration-500 ${color.replace('text-', 'bg-').replace('100', '500')}`} />
      
      <div className="flex items-center justify-between relative z-10">
        <div>
          <p className="text-sm font-medium text-slate-500 mb-1">{title}</p>
          <p className="text-3xl font-bold tracking-tight text-slate-800">{value}</p>
          {trend && (
            <div className={`flex items-center mt-3 text-sm font-medium ${trend.isPositive ? 'text-emerald-600' : 'text-rose-600'}`}>
              {trend.isPositive ? <TrendingUp className="w-4 h-4 mr-1.5" /> : <TrendingDown className="w-4 h-4 mr-1.5" />}
              <span>{Math.abs(trend.value)}%</span>
            </div>
          )}
        </div>
        <div className={`p-4 rounded-2xl ${color}`}>
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
      title: 'Total Log Events',
      value: metrics.total_events.toLocaleString(),
      icon: <Mail className="w-6 h-6" />,
      color: 'bg-blue-50 text-blue-600',
    },
    {
      title: 'Total Sent',
      value: metrics.send_count.toLocaleString(),
      icon: <Send className="w-6 h-6" />,
      color: 'bg-sky-50 text-sky-600',
    },
    {
      title: 'Total Delivered',
      value: metrics.delivery_count.toLocaleString(),
      icon: <CheckCircle className="w-6 h-6" />,
      color: 'bg-emerald-50 text-emerald-600',
    },
    {
      title: 'Delivery Rate',
      value: `${metrics.delivery_rate.toFixed(1)}%`,
      icon: <TrendingUp className="w-6 h-6" />,
      color: 'bg-teal-50 text-teal-600',
    },
    {
      title: 'Total Bounced',
      value: metrics.bounce_count.toLocaleString(),
      icon: <XCircle className="w-6 h-6" />,
      color: 'bg-rose-50 text-rose-600',
    },
    {
      title: 'Total Complaints',
      value: metrics.complaint_count.toLocaleString(),
      icon: <AlertCircle className="w-6 h-6" />,
      color: 'bg-amber-50 text-amber-600',
    },
    {
      title: 'Total Opens',
      value: metrics.open_count.toLocaleString(),
      icon: <Eye className="w-6 h-6" />,
      color: 'bg-purple-50 text-purple-600',
    },
    {
      title: 'Total Clicks',
      value: metrics.click_count.toLocaleString(),
      icon: <MousePointer className="w-6 h-6" />,
      color: 'bg-indigo-50 text-indigo-600',
    },
  ];

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
      {cards.map((card, index) => (
        <MetricCard key={index} {...card} />
      ))}
    </div>
  );
};

export default MetricsGrid;
