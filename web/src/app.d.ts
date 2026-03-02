/// <reference types="@sveltejs/kit" />
/// <reference types="@testing-library/jest-dom" />

// See https://svelte.dev/docs/kit/types#app
declare global {
  namespace App {
    interface Error {
      message: string;
      code?: string;
      details?: unknown;
    }

    interface Locals {
      // Server-side locals
      user?: {
        id: string;
        email: string;
        name: string;
        role: 'admin' | 'member' | 'viewer';
      };
      session?: {
        id: string;
        expiresAt: Date;
      };
    }

    // eslint-disable-next-line @typescript-eslint/no-empty-object-type
    interface PageData {
      // Page-level data will be added as needed
    }

    // eslint-disable-next-line @typescript-eslint/no-empty-object-type
    interface PageState {
      // Page state for shallow routing will be added as needed
    }

    // eslint-disable-next-line @typescript-eslint/no-empty-object-type
    interface Platform {
      // Platform-specific properties will be added as needed
    }
  }

  // Extend Window interface for client-side globals
  // eslint-disable-next-line @typescript-eslint/no-empty-object-type
  interface Window {
    // Client-side globals will be added as needed
  }

  // Type definitions for environment variables
  namespace NodeJS {
    interface ProcessEnv {
      NODE_ENV: 'development' | 'production' | 'test';
      PUBLIC_API_URL?: string;
    }
  }
}

// Ensure this is treated as a module
export {};
