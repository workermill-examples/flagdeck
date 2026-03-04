export interface Flag {
  id: string;
  key: string;
  name: string;
  description: string;
  type: "boolean" | "string" | "number" | "json";
  default_value: unknown;
  is_active: boolean;
  tags: string[];
  environments: Record<string, FlagEnvironment>;
  created_by: string;
  updated_by: string;
  created_at: string;
  updated_at: string;
}

export interface FlagEnvironment {
  enabled: boolean;
  value: unknown;
  rollout_percent: number;
  targeting_rules: TargetingRule[];
}

export interface TargetingRule {
  id: string;
  priority: number;
  conditions: Condition[];
  value: unknown;
}

export interface Condition {
  property: string;
  operator:
    | "equals"
    | "not_equals"
    | "contains"
    | "in"
    | "not_in"
    | "gt"
    | "lt"
    | "gte"
    | "lte"
    | "regex";
  value: unknown;
}

export interface Environment {
  id: string;
  key: string;
  name: string;
  description: string;
  color: string;
  sort_order: number;
  is_active: boolean;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface Segment {
  id: string;
  key: string;
  name: string;
  description: string;
  rules: SegmentRule[];
  created_by: string;
  updated_by: string;
  created_at: string;
  updated_at: string;
}

export interface SegmentRule {
  conditions: Condition[];
}

export interface Experiment {
  id: string;
  key: string;
  name: string;
  description: string;
  flag_key: string;
  status: "draft" | "running" | "paused" | "completed";
  variants: ExperimentVariant[];
  start_date: string | null;
  end_date: string | null;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface ExperimentVariant {
  key: string;
  name: string;
  weight: number;
  value: unknown;
  results?: VariantResults;
}

export interface VariantResults {
  impressions: number;
  conversions: number;
  revenue: number;
}

export interface User {
  id: string;
  email: string;
  name: string;
  role: "admin" | "editor" | "viewer";
  created_at: string;
  updated_at: string;
}

export interface ApiKey {
  id: string;
  name: string;
  key_prefix: string;
  environment: string;
  permissions: string[];
  created_by: string;
  last_used_at: string | null;
  created_at: string;
}

export interface AuditLogEntry {
  id: string;
  action: string;
  resource: string;
  resource_id: string;
  user_id: string;
  user_email: string;
  changes: unknown;
  metadata: unknown;
  timestamp: string;
}

export interface ApiResponse<T> {
  data: T;
  total?: number;
}

export interface ApiError {
  error: {
    code: string;
    message: string;
  };
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  name: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  token_type: string;
}

export interface RefreshRequest {
  refresh_token: string;
}

export interface EvaluationRequest {
  flag_key: string;
  context: Record<string, unknown>;
}

export interface EvaluationResponse {
  key: string;
  value: unknown;
  type: string;
  reason: string;
  rule_id?: string;
  environment: string;
  evaluation_ms: number;
}
