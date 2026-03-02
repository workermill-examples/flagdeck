/**
 * Vitest Test Setup
 * Configures the testing environment for Svelte 5 components
 * Requires: vitest@^2.1.8 with browser resolve conditions
 */
import '@testing-library/jest-dom/vitest';
import { cleanup } from '@testing-library/svelte';
import { afterEach, vi, beforeEach } from 'vitest';

// Clean up after each test
afterEach(() => {
	cleanup();
});

// Mock SvelteKit modules
vi.mock('$app/environment', () => ({
	browser: true,
	dev: true,
	building: false,
	version: 'test'
}));

vi.mock('$app/navigation', () => ({
	goto: vi.fn(),
	invalidate: vi.fn(),
	invalidateAll: vi.fn(),
	preloadCode: vi.fn(),
	preloadData: vi.fn(),
	pushState: vi.fn(),
	replaceState: vi.fn(),
	afterNavigate: vi.fn(),
	beforeNavigate: vi.fn(),
	onNavigate: vi.fn()
}));

vi.mock('$app/stores', () => {
	// Helper function for creating writable stores in tests
	// Prefixed with underscore to indicate it's currently unused but may be needed
	const _writable = (value: unknown) => {
		let subscribers: Array<(value: unknown) => void> = [];
		let currentValue = value;

		return {
			subscribe: (fn: (value: unknown) => void) => {
				subscribers.push(fn);
				fn(currentValue);
				return () => {
					subscribers = subscribers.filter(s => s !== fn);
				};
			},
			set: (value: unknown) => {
				currentValue = value;
				subscribers.forEach(fn => fn(value));
			},
			update: (fn: (value: unknown) => unknown) => {
				currentValue = fn(currentValue);
				subscribers.forEach(s => s(currentValue));
			}
		};
	};

	// Suppress unused variable warning
	void _writable;

	const readable = (value: unknown) => ({
		subscribe: (fn: (value: unknown) => void) => {
			fn(value);
			return () => {};
		}
	});

	// Note: _writable is defined for potential future use in tests
	// Currently only readable stores are exposed
	return {
		page: readable({
			url: new URL('http://localhost:5173'),
			params: {},
			route: { id: '/' },
			status: 200,
			error: null,
			data: {},
			form: null,
			state: {}
		}),
		navigating: readable(null),
		updated: readable(false),
		getStores: vi.fn()
	};
});

// Mock fetch for tests
global.fetch = vi.fn();

// Add custom matchers or global test utilities here
declare global {
	// eslint-disable-next-line @typescript-eslint/no-namespace
	namespace Vi {
		// eslint-disable-next-line @typescript-eslint/no-empty-object-type
		interface Assertion {
			// Custom matchers will be added as needed
		}
	}
}

// Reset mocks before each test
beforeEach(() => {
	vi.clearAllMocks();
});