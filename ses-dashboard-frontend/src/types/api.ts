export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: {
    id: number;
    username: string;
    email: string;
    role: string;
  };
}

export interface Event {
  ID: number;
  MessageID: string;
  Email: string;
  Subject: string;
  EventType: string;
  Status: string;
  Reason?: string;
  Source: string;
  Recipients: string;
  EventTimestamp: string;
  BounceType?: string;
  BounceSubType?: string;
  DiagnosticCode?: string;
  ProcessingTimeMillis?: number;
  SmtpResponse?: string;
  RemoteMtaIp?: string;
  ReportingMTA?: string;
  Tags?: string;
}

export interface PaginationInfo {
  page: number;
  limit: number;
  total: number;
  totalPages: number;
  hasNext: boolean;
  hasPrev: boolean;
}

export interface MessageGroup {
  message_id: string;
  email: string;
  subject: string;
  source: string;
  latest_status: string;
  latest_event: string;
  event_types: string[];
  event_count: number;
  first_event_at: string;
  last_event_at: string;
  has_bounce: boolean;
  has_complaint: boolean;
  has_delivery: boolean;
  has_open: boolean;
  has_click: boolean;
}

export interface EventsResponse {
  events: MessageGroup[];
  pagination: PaginationInfo;
}

export interface EventDetailResponse {
  message_id: string;
  events: Event[];
}

export interface MetricsResponse {
  total_events: number;
  send_count: number;
  delivery_count: number;
  bounce_count: number;
  complaint_count: number;
  open_count: number;
  click_count: number;
  bounce_rate: number;
  delivery_rate: number;
}

// API Types for SES Dashboard
export interface DailyMetrics {
  date: string;
  total_events: number;
  send_count: number;
  delivery_count: number;
  bounce_count: number;
  complaint_count: number;
  open_count: number;
  click_count: number;
  bounce_rate: number;
  delivery_rate: number;
}

export interface MonthlyMetrics {
  month: string;
  total_events: number;
  send_count: number;
  delivery_count: number;
  bounce_count: number;
  complaint_count: number;
  open_count: number;
  click_count: number;
  bounce_rate: number;
  delivery_rate: number;
}

export interface HourlyMetrics {
  hour: string;
  total_events: number;
  send_count: number;
  delivery_count: number;
  bounce_count: number;
  complaint_count: number;
  open_count: number;
  click_count: number;
  bounce_rate: number;
  delivery_rate: number;
}

export interface User {
  id: number;
  username: string;
  email: string;
  role: string;
  active: boolean;
}

export interface CreateUserRequest {
  username: string;
  password: string;
  email: string;
  role: string;
}
