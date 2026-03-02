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

		interface PageData {
			// Page-level data
		}

		interface PageState {
			// Page state for shallow routing
		}

		interface Platform {
			// Platform-specific properties
		}
	}

	// Extend Window interface for client-side globals
	interface Window {
		// Add any client-side globals here
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