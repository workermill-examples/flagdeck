<script lang="ts">
	import '../app.css';

	// Using Svelte 5 runes syntax
	let count = $state(0);
	let message = $state('Welcome to FlagDeck');

	// Derived state using $derived rune
	let displayMessage = $derived(`${message} - Built by WorkerMill`);

	// Effect using $effect rune
	$effect(() => {
		// Count changed effect - logging disabled per ESLint rules
		// In production, this would trigger analytics or other side effects
		if (count > 10) {
			message = 'You really like clicking!';
		} else {
			message = 'Welcome to FlagDeck';
		}
	});

	function increment() {
		count++;
	}

	function decrement() {
		count--;
	}

	function reset() {
		count = 0;
	}
</script>

<svelte:head>
	<title>FlagDeck - Feature Flag Management Platform</title>
	<meta name="description" content="FlagDeck is a powerful feature flag management platform built by WorkerMill" />
</svelte:head>

<main class="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
	<div class="container-fluid py-12">
		<div class="max-w-4xl mx-auto">
			<!-- Header -->
			<header class="text-center mb-12">
				<h1 class="text-5xl font-bold text-gray-900 mb-4">
					FlagDeck
				</h1>
				<p class="text-xl text-gray-600">
					Feature Flag Management Platform
				</p>
				<p class="text-sm text-gray-500 mt-2">
					v1.0.0 - Built by WorkerMill
				</p>
			</header>

			<!-- Placeholder Card -->
			<div class="card">
				<div class="card-header">
					<h2 class="text-lg font-semibold text-gray-900">
						Placeholder Page
					</h2>
				</div>
				<div class="card-body space-y-6">
					<div class="text-center">
						<p class="text-gray-600 mb-4">
							{displayMessage}
						</p>
						<p class="text-sm text-gray-500">
							This is a minimal placeholder page demonstrating Svelte 5 runes syntax.
						</p>
					</div>

					<!-- Counter Demo -->
					<div class="border-t pt-6">
						<h3 class="text-md font-semibold text-gray-900 mb-4">
							Counter Demo (Svelte 5 Runes)
						</h3>
						<div class="flex items-center justify-center space-x-4">
							<button
								onclick={decrement}
								class="btn-secondary"
								disabled={count <= 0}
							>
								-
							</button>
							<span class="text-2xl font-bold text-gray-900 min-w-[3rem] text-center">
								{count}
							</span>
							<button
								onclick={increment}
								class="btn-secondary"
							>
								+
							</button>
						</div>
						<div class="text-center mt-4">
							<button
								onclick={reset}
								class="btn-primary"
								disabled={count === 0}
							>
								Reset
							</button>
						</div>
					</div>

					<!-- Project Status -->
					<div class="border-t pt-6">
						<h3 class="text-md font-semibold text-gray-900 mb-4">
							Project Status
						</h3>
						<div class="space-y-2 text-sm">
							<div class="flex items-center space-x-2">
								<span class="w-2 h-2 bg-green-500 rounded-full"></span>
								<span class="text-gray-600">SvelteKit 2.x with Svelte 5</span>
							</div>
							<div class="flex items-center space-x-2">
								<span class="w-2 h-2 bg-green-500 rounded-full"></span>
								<span class="text-gray-600">TailwindCSS 4.x (Beta)</span>
							</div>
							<div class="flex items-center space-x-2">
								<span class="w-2 h-2 bg-green-500 rounded-full"></span>
								<span class="text-gray-600">TypeScript (Strict Mode)</span>
							</div>
							<div class="flex items-center space-x-2">
								<span class="w-2 h-2 bg-green-500 rounded-full"></span>
								<span class="text-gray-600">Vitest + Testing Library</span>
							</div>
							<div class="flex items-center space-x-2">
								<span class="w-2 h-2 bg-green-500 rounded-full"></span>
								<span class="text-gray-600">ESLint 9 (Flat Config)</span>
							</div>
						</div>
					</div>

					<!-- API Health Check -->
					<div class="border-t pt-6">
						<h3 class="text-md font-semibold text-gray-900 mb-4">
							API Endpoint
						</h3>
						<p class="text-sm text-gray-600">
							Health check available at: <code class="bg-gray-100 px-2 py-1 rounded">/api/v1/health</code>
						</p>
					</div>
				</div>
			</div>

			<!-- Footer -->
			<footer class="text-center mt-12 text-sm text-gray-500">
				<p>FlagDeck &copy; 2024 - Feature Flag Management Platform</p>
			</footer>
		</div>
	</div>
</main>

<style>
	/* Component-specific styles if needed */
	code {
		font-family: 'JetBrains Mono', monospace;
	}
</style>