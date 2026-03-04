import { api } from "./api.js";
import type { AuthResponse, User } from "./types.js";

export interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}

// Create a simple auth store using Svelte 5 runes
function createAuthStore() {
  let state = $state<AuthState>({
    user: null,
    isAuthenticated: false,
    isLoading: true,
  });

  return {
    get user() {
      return state.user;
    },
    get isAuthenticated() {
      return state.isAuthenticated;
    },
    get isLoading() {
      return state.isLoading;
    },

    async login(email: string, password: string): Promise<AuthResponse> {
      state.isLoading = true;
      try {
        const authResponse = await api.login(email, password);
        const user = await api.getCurrentUser();
        state.user = user;
        state.isAuthenticated = true;
        return authResponse;
      } catch (error) {
        state.user = null;
        state.isAuthenticated = false;
        throw error;
      } finally {
        state.isLoading = false;
      }
    },

    async register(
      email: string,
      password: string,
      name: string,
    ): Promise<AuthResponse> {
      state.isLoading = true;
      try {
        const authResponse = await api.register(email, password, name);
        const user = await api.getCurrentUser();
        state.user = user;
        state.isAuthenticated = true;
        return authResponse;
      } catch (error) {
        state.user = null;
        state.isAuthenticated = false;
        throw error;
      } finally {
        state.isLoading = false;
      }
    },

    async logout(): Promise<void> {
      state.isLoading = true;
      try {
        await api.logout();
      } catch (error) {
        console.warn("Logout error:", error);
      } finally {
        state.user = null;
        state.isAuthenticated = false;
        state.isLoading = false;
      }
    },

    async checkAuth(): Promise<void> {
      state.isLoading = true;
      try {
        if (!hasValidToken()) {
          state.user = null;
          state.isAuthenticated = false;
          return;
        }

        const user = await api.getCurrentUser();
        state.user = user;
        state.isAuthenticated = true;
      } catch (error) {
        console.warn("Auth check failed:", error);
        state.user = null;
        state.isAuthenticated = false;
        clearTokens();
      } finally {
        state.isLoading = false;
      }
    },

    setUser(user: User | null): void {
      state.user = user;
      state.isAuthenticated = !!user;
    },
  };
}

export const auth = createAuthStore();

// Token management utilities
export function hasValidToken(): boolean {
  if (typeof window === "undefined") return false;

  const token = localStorage.getItem("access_token");
  const expiresIn = localStorage.getItem("expires_in");

  if (!token || !expiresIn) return false;

  // Simple check - in production, you'd want more sophisticated token validation
  return true;
}

export function clearTokens(): void {
  if (typeof window === "undefined") return;

  localStorage.removeItem("access_token");
  localStorage.removeItem("refresh_token");
  localStorage.removeItem("expires_in");
  localStorage.removeItem("token_type");
}

export function getAccessToken(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem("access_token");
}

export function getRefreshToken(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem("refresh_token");
}

// Check if user is authenticated (synchronous check for guards)
export function isAuthenticated(): boolean {
  return hasValidToken();
}
