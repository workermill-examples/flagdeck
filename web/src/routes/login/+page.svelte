<script lang="ts">
  import { goto } from "$app/navigation";
  import { auth } from "$lib/auth.svelte.js";

  let email = $state("demo@workermill.com");
  let password = $state("demo1234");
  let error = $state("");
  let isSubmitting = $state(false);

  // Redirect if already authenticated
  $effect(() => {
    if (auth.isAuthenticated) {
      goto("/dashboard", { replaceState: true });
    }
  });

  async function handleLogin(event: Event) {
    event.preventDefault();

    if (isSubmitting) return;

    error = "";
    isSubmitting = true;

    try {
      await auth.login(email, password);
      goto("/dashboard", { replaceState: true });
    } catch (err) {
      error = err instanceof Error ? err.message : "Login failed";
    } finally {
      isSubmitting = false;
    }
  }
</script>

<svelte:head>
  <title>Login - FlagDeck</title>
</svelte:head>

<div
  class="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8"
>
  <div class="max-w-md w-full space-y-8">
    <!-- Header -->
    <div class="text-center">
      <a href="/" class="inline-block">
        <div
          class="mx-auto w-16 h-16 bg-blue-600 rounded-xl flex items-center justify-center mb-4"
        >
          <span class="text-white font-bold text-xl">FD</span>
        </div>
      </a>
      <h2 class="text-3xl font-bold text-gray-900">Sign in to FlagDeck</h2>
      <p class="mt-2 text-sm text-gray-600">
        Manage your feature flags with confidence
      </p>
    </div>

    <!-- Login Form -->
    <form class="mt-8 space-y-6" onsubmit={handleLogin}>
      <div class="space-y-4">
        <div>
          <label
            for="email"
            class="block text-sm font-medium text-gray-700"
          >
            Email address
          </label>
          <input
            id="email"
            name="email"
            type="email"
            autocomplete="email"
            required
            bind:value={email}
            disabled={isSubmitting}
            class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-50 disabled:text-gray-500"
            placeholder="Enter your email"
          />
        </div>

        <div>
          <label
            for="password"
            class="block text-sm font-medium text-gray-700"
          >
            Password
          </label>
          <input
            id="password"
            name="password"
            type="password"
            autocomplete="current-password"
            required
            bind:value={password}
            disabled={isSubmitting}
            class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-50 disabled:text-gray-500"
            placeholder="Enter your password"
          />
        </div>
      </div>

      <!-- Error Message -->
      {#if error}
        <div class="bg-red-50 border border-red-200 rounded-lg p-3">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg
                class="h-5 w-5 text-red-400"
                viewBox="0 0 20 20"
                fill="currentColor"
              >
                <path
                  fill-rule="evenodd"
                  d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.28 7.22a.75.75 0 00-1.06 1.06L8.94 10l-1.72 1.72a.75.75 0 101.06 1.06L10 11.06l1.72 1.72a.75.75 0 101.06-1.06L11.06 10l1.72-1.72a.75.75 0 00-1.06-1.06L10 8.94 8.28 7.22z"
                  clip-rule="evenodd"
                />
              </svg>
            </div>
            <div class="ml-3">
              <p class="text-sm text-red-800">{error}</p>
            </div>
          </div>
        </div>
      {/if}

      <!-- Submit Button -->
      <div>
        <button
          type="submit"
          disabled={isSubmitting}
          class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-lg text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          {#if isSubmitting}
            <span class="flex items-center">
              <svg
                class="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
                fill="none"
                viewBox="0 0 24 24"
              >
                <circle
                  class="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  stroke-width="4"
                ></circle>
                <path
                  class="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                ></path>
              </svg>
              Signing in...
            </span>
          {:else}
            Sign in
          {/if}
        </button>
      </div>

      <!-- Demo Credentials Notice -->
      <div class="mt-4 p-3 bg-blue-50 border border-blue-200 rounded-lg">
        <p class="text-sm text-blue-800">
          <strong>Demo credentials:</strong><br />
          Email: demo@workermill.com<br />
          Password: demo1234
        </p>
      </div>
    </form>

    <!-- Back to landing -->
    <div class="text-center">
      <a
        href="/"
        class="text-sm text-gray-500 hover:text-gray-700 transition-colors"
      >
        &larr; Back to home
      </a>
    </div>
  </div>
</div>
