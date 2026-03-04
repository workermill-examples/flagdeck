<script lang="ts">
  import { goto } from "$app/navigation";
  import { auth } from "$lib/auth.js";

  // If user is already authenticated, redirect to dashboard
  $effect(() => {
    if (!auth.isLoading && auth.isAuthenticated) {
      goto("/dashboard", { replaceState: true });
    }
  });

  let hoveredFeature = $state<number | null>(null);
</script>

<svelte:head>
  <title>FlagDeck - Open-Source Feature Flag Platform</title>
  <meta
    name="description"
    content="FlagDeck is an open-source feature flag platform with targeting rules, percentage rollouts, experiments, and multi-environment support."
  />
</svelte:head>

<div class="min-h-screen bg-white">
  <!-- Navigation -->
  <nav class="border-b border-gray-100">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center h-16">
        <div class="flex items-center space-x-3">
          <div
            class="w-9 h-9 bg-blue-600 rounded-lg flex items-center justify-center text-white font-bold text-sm"
          >
            FD
          </div>
          <span class="text-xl font-bold text-gray-900">FlagDeck</span>
        </div>
        <div class="flex items-center space-x-4">
          <a
            href="https://github.com/workermill-examples/flagdeck"
            target="_blank"
            rel="noopener noreferrer"
            class="text-sm text-gray-600 hover:text-gray-900 transition-colors"
          >
            GitHub
          </a>
          <a
            href="/login"
            class="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-lg hover:bg-blue-700 transition-colors"
          >
            Sign in
          </a>
        </div>
      </div>
    </div>
  </nav>

  <!-- Hero Section -->
  <div class="relative overflow-hidden">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-20 pb-16">
      <div class="text-center max-w-3xl mx-auto">
        <div
          class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-blue-50 text-blue-700 mb-6"
        >
          Open-source feature flag platform
        </div>
        <h1 class="text-5xl sm:text-6xl font-extrabold text-gray-900 tracking-tight leading-tight">
          Ship features
          <span class="text-blue-600">with confidence</span>
        </h1>
        <p class="mt-6 text-xl text-gray-600 leading-relaxed max-w-2xl mx-auto">
          FlagDeck gives you complete control over feature rollouts with
          targeting rules, percentage-based rollouts, A/B experiments, and
          multi-environment management.
        </p>
        <div class="mt-10 flex flex-col sm:flex-row items-center justify-center gap-4">
          <a
            href="/login"
            class="w-full sm:w-auto inline-flex items-center justify-center px-8 py-3 text-base font-semibold text-white bg-blue-600 rounded-lg hover:bg-blue-700 transition-colors shadow-sm"
          >
            Try the Demo
            <svg
              class="ml-2 w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M13 7l5 5m0 0l-5 5m5-5H6"
              />
            </svg>
          </a>
          <a
            href="https://github.com/workermill-examples/flagdeck"
            target="_blank"
            rel="noopener noreferrer"
            class="w-full sm:w-auto inline-flex items-center justify-center px-8 py-3 text-base font-semibold text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
          >
            <svg class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 24 24">
              <path
                d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"
              />
            </svg>
            View on GitHub
          </a>
        </div>
        <!-- Demo credentials hint -->
        <p class="mt-6 text-sm text-gray-500">
          Demo credentials:
          <span class="font-mono text-gray-600">demo@workermill.com</span>
          /
          <span class="font-mono text-gray-600">demo1234</span>
        </p>
      </div>
    </div>
  </div>

  <!-- Stats Bar -->
  <div class="border-y border-gray-100 bg-gray-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div class="grid grid-cols-2 md:grid-cols-5 gap-8">
        <div class="text-center">
          <div class="text-3xl font-bold text-gray-900">10</div>
          <div class="mt-1 text-sm text-gray-500">Feature Flags</div>
        </div>
        <div class="text-center">
          <div class="text-3xl font-bold text-gray-900">3</div>
          <div class="mt-1 text-sm text-gray-500">Environments</div>
        </div>
        <div class="text-center">
          <div class="text-3xl font-bold text-gray-900">3</div>
          <div class="mt-1 text-sm text-gray-500">User Segments</div>
        </div>
        <div class="text-center">
          <div class="text-3xl font-bold text-gray-900">5</div>
          <div class="mt-1 text-sm text-gray-500">Targeting Rules</div>
        </div>
        <div class="text-center col-span-2 md:col-span-1">
          <div class="text-3xl font-bold text-gray-900">100%</div>
          <div class="mt-1 text-sm text-gray-500">Rollout Control</div>
        </div>
      </div>
    </div>
  </div>

  <!-- Features Section -->
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
    <div class="text-center mb-16">
      <h2 class="text-3xl font-bold text-gray-900">
        Everything you need for feature management
      </h2>
      <p class="mt-4 text-lg text-gray-600 max-w-2xl mx-auto">
        A complete feature flag platform with the tools your team needs to ship
        safely and iterate faster.
      </p>
    </div>

    <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
      <!-- Feature 1: Feature Flags -->
      <div
        class="relative rounded-xl border border-gray-200 p-6 hover:border-blue-300 hover:shadow-md transition-all duration-200"
        role="article"
        onmouseenter={() => (hoveredFeature = 0)}
        onmouseleave={() => (hoveredFeature = null)}
      >
        <div
          class="w-10 h-10 rounded-lg bg-blue-100 text-blue-600 flex items-center justify-center mb-4"
        >
          <svg
            class="w-5 h-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M3 21v-4m0 0V5a2 2 0 012-2h6.5l1 1H21l-3 6 3 6h-8.5l-1-1H5a2 2 0 00-2 2zm9-13.5V9"
            />
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-gray-900">Feature Flags</h3>
        <p class="mt-2 text-gray-600 text-sm leading-relaxed">
          Create boolean or multivariate flags. Toggle features on and off
          instantly without redeploying your application.
        </p>
      </div>

      <!-- Feature 2: Targeting Rules -->
      <div
        class="relative rounded-xl border border-gray-200 p-6 hover:border-blue-300 hover:shadow-md transition-all duration-200"
        role="article"
        onmouseenter={() => (hoveredFeature = 1)}
        onmouseleave={() => (hoveredFeature = null)}
      >
        <div
          class="w-10 h-10 rounded-lg bg-green-100 text-green-600 flex items-center justify-center mb-4"
        >
          <svg
            class="w-5 h-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"
            />
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-gray-900">Targeting Rules</h3>
        <p class="mt-2 text-gray-600 text-sm leading-relaxed">
          Define precise targeting rules based on user attributes, segments, and
          custom properties. Deliver the right experience to the right users.
        </p>
      </div>

      <!-- Feature 3: Percentage Rollouts -->
      <div
        class="relative rounded-xl border border-gray-200 p-6 hover:border-blue-300 hover:shadow-md transition-all duration-200"
        role="article"
        onmouseenter={() => (hoveredFeature = 2)}
        onmouseleave={() => (hoveredFeature = null)}
      >
        <div
          class="w-10 h-10 rounded-lg bg-purple-100 text-purple-600 flex items-center justify-center mb-4"
        >
          <svg
            class="w-5 h-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"
            />
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-gray-900">Percentage Rollouts</h3>
        <p class="mt-2 text-gray-600 text-sm leading-relaxed">
          Gradually roll out features to a percentage of your user base. Reduce
          risk by incrementally increasing exposure.
        </p>
      </div>

      <!-- Feature 4: Environments -->
      <div
        class="relative rounded-xl border border-gray-200 p-6 hover:border-blue-300 hover:shadow-md transition-all duration-200"
        role="article"
        onmouseenter={() => (hoveredFeature = 3)}
        onmouseleave={() => (hoveredFeature = null)}
      >
        <div
          class="w-10 h-10 rounded-lg bg-orange-100 text-orange-600 flex items-center justify-center mb-4"
        >
          <svg
            class="w-5 h-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9"
            />
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-gray-900">Multi-Environment</h3>
        <p class="mt-2 text-gray-600 text-sm leading-relaxed">
          Manage flags across development, staging, and production environments
          independently. Test safely before going live.
        </p>
      </div>

      <!-- Feature 5: Experiments -->
      <div
        class="relative rounded-xl border border-gray-200 p-6 hover:border-blue-300 hover:shadow-md transition-all duration-200"
        role="article"
        onmouseenter={() => (hoveredFeature = 4)}
        onmouseleave={() => (hoveredFeature = null)}
      >
        <div
          class="w-10 h-10 rounded-lg bg-rose-100 text-rose-600 flex items-center justify-center mb-4"
        >
          <svg
            class="w-5 h-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 10.172V5L8 4z"
            />
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-gray-900">A/B Experiments</h3>
        <p class="mt-2 text-gray-600 text-sm leading-relaxed">
          Run controlled experiments to measure the impact of features. Make
          data-driven decisions with confidence.
        </p>
      </div>

      <!-- Feature 6: Audit Log -->
      <div
        class="relative rounded-xl border border-gray-200 p-6 hover:border-blue-300 hover:shadow-md transition-all duration-200"
        role="article"
        onmouseenter={() => (hoveredFeature = 5)}
        onmouseleave={() => (hoveredFeature = null)}
      >
        <div
          class="w-10 h-10 rounded-lg bg-cyan-100 text-cyan-600 flex items-center justify-center mb-4"
        >
          <svg
            class="w-5 h-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
            />
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-gray-900">Audit Log</h3>
        <p class="mt-2 text-gray-600 text-sm leading-relaxed">
          Complete audit trail of every change. Know who changed what, when, and
          why. Full accountability for your team.
        </p>
      </div>
    </div>
  </div>

  <!-- How It Works -->
  <div class="bg-gray-50 border-y border-gray-100">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
      <div class="text-center mb-16">
        <h2 class="text-3xl font-bold text-gray-900">How it works</h2>
        <p class="mt-4 text-lg text-gray-600">
          Three simple steps to safer feature releases
        </p>
      </div>

      <div class="grid md:grid-cols-3 gap-12">
        <div class="text-center">
          <div
            class="w-12 h-12 rounded-full bg-blue-600 text-white flex items-center justify-center text-lg font-bold mx-auto mb-4"
          >
            1
          </div>
          <h3 class="text-lg font-semibold text-gray-900 mb-2">
            Create a Flag
          </h3>
          <p class="text-gray-600 text-sm">
            Define your feature flag with a key, description, and default value.
            Choose boolean or multivariate types.
          </p>
        </div>
        <div class="text-center">
          <div
            class="w-12 h-12 rounded-full bg-blue-600 text-white flex items-center justify-center text-lg font-bold mx-auto mb-4"
          >
            2
          </div>
          <h3 class="text-lg font-semibold text-gray-900 mb-2">
            Set Up Rules
          </h3>
          <p class="text-gray-600 text-sm">
            Configure targeting rules, user segments, and percentage rollouts.
            Control exactly who sees each feature.
          </p>
        </div>
        <div class="text-center">
          <div
            class="w-12 h-12 rounded-full bg-blue-600 text-white flex items-center justify-center text-lg font-bold mx-auto mb-4"
          >
            3
          </div>
          <h3 class="text-lg font-semibold text-gray-900 mb-2">
            Roll Out Safely
          </h3>
          <p class="text-gray-600 text-sm">
            Gradually increase exposure across environments. Monitor with the
            audit log and roll back instantly if needed.
          </p>
        </div>
      </div>
    </div>
  </div>

  <!-- CTA Section -->
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
    <div
      class="bg-blue-600 rounded-2xl px-8 py-12 text-center text-white shadow-lg"
    >
      <h2 class="text-3xl font-bold">Ready to try FlagDeck?</h2>
      <p class="mt-4 text-lg text-blue-100 max-w-xl mx-auto">
        Explore the demo with pre-configured flags, environments, segments, and
        experiments. No setup required.
      </p>
      <div class="mt-8 flex flex-col sm:flex-row items-center justify-center gap-4">
        <a
          href="/login"
          class="w-full sm:w-auto inline-flex items-center justify-center px-8 py-3 text-base font-semibold text-blue-600 bg-white rounded-lg hover:bg-blue-50 transition-colors"
        >
          Try the Demo
        </a>
        <a
          href="https://github.com/workermill-examples/flagdeck"
          target="_blank"
          rel="noopener noreferrer"
          class="w-full sm:w-auto inline-flex items-center justify-center px-8 py-3 text-base font-semibold text-white border-2 border-white/30 rounded-lg hover:bg-white/10 transition-colors"
        >
          View Source Code
        </a>
      </div>
      <p class="mt-6 text-sm text-blue-200">
        Demo: <span class="font-mono">demo@workermill.com</span> /
        <span class="font-mono">demo1234</span>
      </p>
    </div>
  </div>

  <!-- Footer -->
  <footer class="border-t border-gray-100">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div class="flex flex-col sm:flex-row items-center justify-between gap-4">
        <div class="flex items-center space-x-2">
          <div
            class="w-7 h-7 bg-blue-600 rounded-md flex items-center justify-center text-white font-bold text-xs"
          >
            FD
          </div>
          <span class="text-sm text-gray-500">
            FlagDeck &mdash; Built by
            <a
              href="https://workermill.com"
              target="_blank"
              rel="noopener noreferrer"
              class="text-blue-600 hover:text-blue-800 font-medium"
            >
              WorkerMill
            </a>
          </span>
        </div>
        <div class="flex items-center space-x-6">
          <a
            href="https://github.com/workermill-examples/flagdeck"
            target="_blank"
            rel="noopener noreferrer"
            class="text-sm text-gray-500 hover:text-gray-700 transition-colors"
          >
            GitHub
          </a>
          <a
            href="https://workermill.com"
            target="_blank"
            rel="noopener noreferrer"
            class="text-sm text-gray-500 hover:text-gray-700 transition-colors"
          >
            WorkerMill
          </a>
        </div>
      </div>
    </div>
  </footer>
</div>
