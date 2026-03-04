import { PUBLIC_API_URL } from "$env/static/public";
import type {
  ApiError,
  ApiResponse,
  AuthResponse,
  RefreshRequest,
} from "./types.js";

class ApiClient {
  private baseUrl: string;

  constructor() {
    this.baseUrl = PUBLIC_API_URL || "http://localhost:8080";
  }

  private getAuthToken(): string | null {
    if (typeof window === "undefined") return null;
    return localStorage.getItem("access_token");
  }

  private getRefreshToken(): string | null {
    if (typeof window === "undefined") return null;
    return localStorage.getItem("refresh_token");
  }

  private setTokens(authResponse: AuthResponse): void {
    if (typeof window === "undefined") return;
    localStorage.setItem("access_token", authResponse.access_token);
    localStorage.setItem("refresh_token", authResponse.refresh_token);
    localStorage.setItem("expires_in", authResponse.expires_in.toString());
    localStorage.setItem("token_type", authResponse.token_type);
  }

  private clearTokens(): void {
    if (typeof window === "undefined") return;
    localStorage.removeItem("access_token");
    localStorage.removeItem("refresh_token");
    localStorage.removeItem("expires_in");
    localStorage.removeItem("token_type");
  }

  private async refreshAccessToken(): Promise<boolean> {
    const refreshToken = this.getRefreshToken();
    if (!refreshToken) return false;

    try {
      const response = await fetch(`${this.baseUrl}/auth/refresh`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          refresh_token: refreshToken,
        } satisfies RefreshRequest),
      });

      if (!response.ok) return false;

      const authResponse: AuthResponse = await response.json();
      this.setTokens(authResponse);
      return true;
    } catch {
      return false;
    }
  }

  async request<T>(
    endpoint: string,
    options: RequestInit = {},
    requireAuth: boolean = false,
  ): Promise<T> {
    const url = endpoint.startsWith("http")
      ? endpoint
      : `${this.baseUrl}${endpoint}`;

    const headers: HeadersInit = {
      "Content-Type": "application/json",
      ...options.headers,
    };

    if (requireAuth) {
      const token = this.getAuthToken();
      if (token) {
        headers.Authorization = `Bearer ${token}`;
      }
    }

    let response = await fetch(url, {
      ...options,
      headers,
    });

    // Handle 401 with token refresh
    if (response.status === 401 && requireAuth) {
      const refreshed = await this.refreshAccessToken();
      if (refreshed) {
        // Retry request with new token
        const newToken = this.getAuthToken();
        if (newToken) {
          headers.Authorization = `Bearer ${newToken}`;
        }
        response = await fetch(url, {
          ...options,
          headers,
        });
      } else {
        // Refresh failed, redirect to login
        this.clearTokens();
        if (typeof window !== "undefined") {
          window.location.href = "/login";
        }
        throw new Error("Authentication failed");
      }
    }

    if (!response.ok) {
      try {
        const errorData: ApiError = await response.json();
        throw new Error(errorData.error.message);
      } catch {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }
    }

    return response.json();
  }

  // Auth methods
  async login(email: string, password: string): Promise<AuthResponse> {
    const response = await this.request<AuthResponse>("/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });
    this.setTokens(response);
    return response;
  }

  async register(
    email: string,
    password: string,
    name: string,
  ): Promise<AuthResponse> {
    const response = await this.request<AuthResponse>("/auth/register", {
      method: "POST",
      body: JSON.stringify({ email, password, name }),
    });
    this.setTokens(response);
    return response;
  }

  async logout(): Promise<void> {
    try {
      await this.request("/auth/logout", { method: "POST" }, true);
    } catch {
      // Ignore logout errors
    }
    this.clearTokens();
  }

  async getCurrentUser() {
    return this.request("/auth/me", { method: "GET" }, true);
  }

  // Data methods
  async get<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: "GET" }, true);
  }

  async post<T>(endpoint: string, data: unknown): Promise<T> {
    return this.request<T>(
      endpoint,
      {
        method: "POST",
        body: JSON.stringify(data),
      },
      true,
    );
  }

  async put<T>(endpoint: string, data: unknown): Promise<T> {
    return this.request<T>(
      endpoint,
      {
        method: "PUT",
        body: JSON.stringify(data),
      },
      true,
    );
  }

  async delete<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: "DELETE" }, true);
  }

  // Flag methods
  async getFlags(): Promise<ApiResponse<any[]>> {
    return this.get("/api/v1/flags");
  }

  async getFlag(key: string) {
    return this.get(`/api/v1/flags/${key}`);
  }

  async createFlag(data: unknown) {
    return this.post("/api/v1/flags", data);
  }

  async updateFlag(key: string, data: unknown) {
    return this.put(`/api/v1/flags/${key}`, data);
  }

  async deleteFlag(key: string) {
    return this.delete(`/api/v1/flags/${key}`);
  }

  async toggleFlag(key: string, environment?: string) {
    const body = environment ? { environment } : {};
    return this.post(`/api/v1/flags/${key}/toggle`, body);
  }

  // Environment methods
  async getEnvironments(): Promise<ApiResponse<any[]>> {
    return this.get("/api/v1/environments");
  }

  async createEnvironment(data: unknown) {
    return this.post("/api/v1/environments", data);
  }

  async updateEnvironment(key: string, data: unknown) {
    return this.put(`/api/v1/environments/${key}`, data);
  }

  async deleteEnvironment(key: string) {
    return this.delete(`/api/v1/environments/${key}`);
  }

  // Segment methods
  async getSegments(): Promise<ApiResponse<any[]>> {
    return this.get("/api/v1/segments");
  }

  async getSegment(key: string) {
    return this.get(`/api/v1/segments/${key}`);
  }

  async createSegment(data: unknown) {
    return this.post("/api/v1/segments", data);
  }

  async updateSegment(key: string, data: unknown) {
    return this.put(`/api/v1/segments/${key}`, data);
  }

  async deleteSegment(key: string) {
    return this.delete(`/api/v1/segments/${key}`);
  }

  // Experiment methods
  async getExperiments(): Promise<ApiResponse<any[]>> {
    return this.get("/api/v1/experiments");
  }

  async getExperiment(key: string) {
    return this.get(`/api/v1/experiments/${key}`);
  }

  async createExperiment(data: unknown) {
    return this.post("/api/v1/experiments", data);
  }

  async updateExperiment(key: string, data: unknown) {
    return this.put(`/api/v1/experiments/${key}`, data);
  }

  async deleteExperiment(key: string) {
    return this.delete(`/api/v1/experiments/${key}`);
  }

  // API Key methods
  async getApiKeys(): Promise<ApiResponse<any[]>> {
    return this.get("/api/v1/api-keys");
  }

  async createApiKey(data: unknown) {
    return this.post("/api/v1/api-keys", data);
  }

  async deleteApiKey(id: string) {
    return this.delete(`/api/v1/api-keys/${id}`);
  }

  // Audit Log methods
  async getAuditLog(params?: URLSearchParams): Promise<ApiResponse<any[]>> {
    const queryString = params ? `?${params.toString()}` : "";
    return this.get(`/api/v1/audit-log${queryString}`);
  }
}

export const api = new ApiClient();
