<script lang="ts">
  import { goto } from "$app/navigation";
  import { auth } from "$lib/auth.svelte.js";
  import { onMount } from "svelte";

  $effect(() => {
    if (!auth.isLoading && auth.isAuthenticated) {
      goto("/dashboard", { replaceState: true });
    }
  });

  let mounted = $state(false);
  onMount(() => {
    requestAnimationFrame(() => {
      mounted = true;
    });
  });
</script>

<svelte:head>
  <title>FlagDeck — Open-Source Feature Flag Platform</title>
  <meta
    name="description"
    content="Ship features with confidence. FlagDeck gives you targeting rules, percentage rollouts, A/B experiments, and multi-environment management — all open source."
  />
  <link rel="preconnect" href="https://fonts.googleapis.com" />
  <link
    rel="preconnect"
    href="https://fonts.gstatic.com"
    crossorigin="anonymous"
  />
  <link
    href="https://fonts.googleapis.com/css2?family=Sora:wght@400;500;600;700;800&family=DM+Sans:ital,wght@0,400;0,500;0,600;0,700;1,400&family=JetBrains+Mono:wght@400;500&display=swap"
    rel="stylesheet"
  />
</svelte:head>

<div class="landing-page">
  <!-- ============================================ -->
  <!-- HERO SECTION (dark)                          -->
  <!-- ============================================ -->
  <section class="hero-section relative overflow-hidden">
    <!-- Ambient gradient orbs -->
    <div
      class="pointer-events-none absolute -left-48 top-0 h-[600px] w-[600px] rounded-full opacity-30 blur-[140px]"
      style="background: radial-gradient(circle, #3b82f6 0%, transparent 70%)"
    ></div>
    <div
      class="pointer-events-none absolute -right-32 bottom-0 h-[500px] w-[500px] rounded-full opacity-20 blur-[120px]"
      style="background: radial-gradient(circle, #6366f1 0%, transparent 70%)"
    ></div>
    <div
      class="pointer-events-none absolute left-1/2 top-1/3 h-[300px] w-[800px] -translate-x-1/2 rounded-full opacity-10 blur-[100px]"
      style="background: radial-gradient(ellipse, #06b6d4 0%, transparent 60%)"
    ></div>

    <!-- Navigation -->
    <nav class="relative z-10">
      <div
        class="mx-auto flex max-w-7xl items-center justify-between px-6 py-5"
      >
        <a href="/" class="flex items-center gap-2.5">
          <div
            class="flex h-8 w-8 items-center justify-center rounded-lg font-heading text-sm font-bold text-white"
            style="background: linear-gradient(135deg, #3b82f6, #6366f1)"
          >
            FD
          </div>
          <span class="font-heading text-lg font-bold text-white">FlagDeck</span
          >
        </a>
        <div class="flex items-center gap-5">
          <a
            href="https://github.com/workermill-examples/flagdeck"
            target="_blank"
            rel="noopener noreferrer"
            class="group flex items-center gap-1.5 text-sm text-white/50 transition-colors hover:text-white/80"
          >
            <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 24 24">
              <path
                d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"
              />
            </svg>
            <span class="hidden sm:inline">Source</span>
          </a>
          <a
            href="/login"
            class="rounded-lg border border-white/10 bg-white/5 px-4 py-2 text-sm font-medium text-white backdrop-blur transition-all hover:border-white/20 hover:bg-white/10"
          >
            Try the Demo
          </a>
        </div>
      </div>
    </nav>

    <!-- Hero Content -->
    <div
      class="relative z-10 mx-auto max-w-7xl px-6 pb-24 pt-16 lg:pb-32 lg:pt-24"
    >
      <div class="grid items-center gap-12 lg:grid-cols-2 lg:gap-16">
        <!-- Left: Text -->
        <div>
          <div
            class="transition-all duration-700 {mounted
              ? 'translate-y-0 opacity-100'
              : 'translate-y-6 opacity-0'}"
          >
            <div
              class="mb-6 inline-flex items-center gap-2 rounded-full border border-blue-400/20 bg-blue-500/10 px-3.5 py-1.5 text-sm text-blue-300"
            >
              <span
                class="inline-block h-1.5 w-1.5 animate-pulse rounded-full bg-blue-400"
              ></span>
              Open-source feature flags
            </div>
          </div>

          <h1
            class="font-heading text-4xl font-extrabold leading-[1.1] tracking-tight text-white sm:text-5xl lg:text-[3.5rem] transition-all duration-700 delay-75 {mounted
              ? 'translate-y-0 opacity-100'
              : 'translate-y-6 opacity-0'}"
          >
            Control your rollouts.
            <span class="hero-gradient-text">Ship with confidence.</span>
          </h1>

          <p
            class="mt-6 max-w-lg text-lg leading-relaxed text-white/60 transition-all duration-700 delay-150 {mounted
              ? 'translate-y-0 opacity-100'
              : 'translate-y-6 opacity-0'}"
          >
            Feature flags with targeting rules, percentage-based rollouts, A/B
            experiments, and multi-environment management. Everything you need
            to release safely.
          </p>

          <div
            class="mt-10 flex flex-col gap-3 sm:flex-row sm:items-center transition-all duration-700 delay-200 {mounted
              ? 'translate-y-0 opacity-100'
              : 'translate-y-6 opacity-0'}"
          >
            <a
              href="/login"
              class="hero-cta-primary group inline-flex items-center justify-center gap-2 rounded-lg px-6 py-3 text-sm font-semibold text-white shadow-lg shadow-blue-500/25 transition-all hover:shadow-blue-500/40"
            >
              Try the Demo
              <svg
                class="h-4 w-4 transition-transform group-hover:translate-x-0.5"
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
              class="inline-flex items-center justify-center gap-2 rounded-lg border border-white/10 px-6 py-3 text-sm font-medium text-white/70 transition-all hover:border-white/20 hover:text-white"
            >
              <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 24 24">
                <path
                  d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"
                />
              </svg>
              View on GitHub
            </a>
          </div>

          <p
            class="mt-6 font-code text-xs text-white/30 transition-all duration-700 delay-300 {mounted
              ? 'translate-y-0 opacity-100'
              : 'translate-y-6 opacity-0'}"
          >
            demo@workermill.com &middot; demo1234
          </p>
        </div>

        <!-- Right: Code Preview -->
        <div
          class="transition-all duration-1000 delay-300 {mounted
            ? 'translate-y-0 opacity-100'
            : 'translate-y-8 opacity-0'}"
        >
          <div
            class="code-window rounded-xl border border-white/[0.06] bg-[#0d1117] shadow-2xl"
          >
            <!-- Window chrome -->
            <div
              class="flex items-center gap-2 border-b border-white/[0.06] px-4 py-3"
            >
              <div class="h-2.5 w-2.5 rounded-full bg-[#ff5f57]"></div>
              <div class="h-2.5 w-2.5 rounded-full bg-[#febc2e]"></div>
              <div class="h-2.5 w-2.5 rounded-full bg-[#28c840]"></div>
              <span class="ml-3 font-code text-xs text-white/25"
                >evaluate.ts</span
              >
            </div>
            <!-- Code content -->
            <div class="overflow-x-auto p-5">
              <pre class="font-code text-[13px] leading-relaxed"><span
                  class="text-[#ff7b72]">import</span
                > <span class="text-[#c9d1d9]">FlagDeck</span> <span
                  class="text-[#ff7b72]">from</span
                > <span class="text-[#a5d6ff]">'@flagdeck/sdk'</span>
<!--  -->
<span class="text-[#ff7b72]">const</span> <span class="text-[#c9d1d9]"
                  >client</span
                > <span class="text-[#ff7b72]">=</span> <span
                  class="text-[#ff7b72]">new</span
                > <span class="text-[#d2a8ff]">FlagDeck</span><span
                  class="text-[#c9d1d9]">(&#123;</span
                >
  <span class="text-[#c9d1d9]">apiKey:</span> <span class="text-[#a5d6ff]"
                  >'fd_live_sk_...'</span
                ><span class="text-[#c9d1d9]">,</span>
  <span class="text-[#c9d1d9]">environment:</span> <span class="text-[#a5d6ff]"
                  >'production'</span
                >
<span class="text-[#c9d1d9]">&#125;)</span>
<!--  -->
<span class="text-[#8b949e]">// Evaluate flag for current user</span>
<span class="text-[#ff7b72]">const</span> <span class="text-[#c9d1d9]"
                  >flag</span
                > <span class="text-[#ff7b72]">=</span> <span
                  class="text-[#ff7b72]">await</span
                > <span class="text-[#c9d1d9]">client.</span><span
                  class="text-[#d2a8ff]">evaluate</span
                ><span class="text-[#c9d1d9]">(</span><span
                  class="text-[#a5d6ff]">'new_dashboard'</span
                ><span class="text-[#c9d1d9]">, &#123;</span>
  <span class="text-[#c9d1d9]">userId:</span> <span class="text-[#a5d6ff]"
                  >'usr_29a1k'</span
                ><span class="text-[#c9d1d9]">,</span>
  <span class="text-[#c9d1d9]">attributes: &#123;</span>
    <span class="text-[#c9d1d9]">plan:</span> <span class="text-[#a5d6ff]"
                  >'enterprise'</span
                ><span class="text-[#c9d1d9]">,</span>
    <span class="text-[#c9d1d9]">country:</span> <span class="text-[#a5d6ff]"
                  >'US'</span
                >
  <span class="text-[#c9d1d9]">&#125;</span>
<span class="text-[#c9d1d9]">&#125;)</span>
<!--  -->
<span class="text-[#ff7b72]">if</span> <span class="text-[#c9d1d9]">(flag.</span
                ><span class="text-[#79c0ff]">enabled</span><span
                  class="text-[#c9d1d9]">) &#123;</span
                >
  <span class="text-[#d2a8ff]">showNewDashboard</span><span
                  class="text-[#c9d1d9]">()</span
                >  <span class="text-[#8b949e]"
                  >// ✓ Targeting rule matched</span
                >
<span class="text-[#c9d1d9]">&#125;</span></pre>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Stats bar (still in dark section) -->
    <div class="relative z-10 border-t border-white/[0.04]">
      <div
        class="mx-auto grid max-w-7xl grid-cols-2 gap-6 px-6 py-10 md:grid-cols-4 transition-all duration-700 delay-500 {mounted
          ? 'translate-y-0 opacity-100'
          : 'translate-y-4 opacity-0'}"
      >
        <div class="text-center">
          <div class="font-heading text-3xl font-bold text-white">10+</div>
          <div class="mt-1 text-sm text-white/40">Feature Flags</div>
        </div>
        <div class="text-center">
          <div class="font-heading text-3xl font-bold text-white">3</div>
          <div class="mt-1 text-sm text-white/40">Environments</div>
        </div>
        <div class="text-center">
          <div class="font-heading text-3xl font-bold text-white">A/B</div>
          <div class="mt-1 text-sm text-white/40">Experiments</div>
        </div>
        <div class="text-center">
          <div class="font-heading text-3xl font-bold text-white">50+</div>
          <div class="mt-1 text-sm text-white/40">Audit Entries</div>
        </div>
      </div>
    </div>
  </section>

  <!-- ============================================ -->
  <!-- FEATURES SECTION (light)                     -->
  <!-- ============================================ -->
  <section class="bg-white py-24">
    <div class="mx-auto max-w-7xl px-6">
      <div class="mx-auto max-w-2xl text-center">
        <h2
          class="font-heading text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl"
        >
          Built for modern teams
        </h2>
        <p class="mt-4 text-lg leading-relaxed text-gray-500">
          A complete feature management platform with everything you need to
          ship safely and iterate faster.
        </p>
      </div>

      <div class="mt-16 grid gap-5 md:grid-cols-2 lg:grid-cols-3">
        <!-- Feature Flags -->
        <div class="feature-card group rounded-2xl border border-gray-100 p-6">
          <div
            class="mb-4 flex h-11 w-11 items-center justify-center rounded-xl bg-blue-50 text-blue-600 transition-colors group-hover:bg-blue-100"
          >
            <svg
              class="h-5 w-5"
              fill="none"
              stroke="currentColor"
              stroke-width="1.75"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M3 21v-4m0 0V5a2 2 0 012-2h6.5l1 1H21l-3 6 3 6h-8.5l-1-1H5a2 2 0 00-2 2zm9-13.5V9"
              />
            </svg>
          </div>
          <h3 class="font-heading text-base font-semibold text-gray-900">
            Feature Flags
          </h3>
          <p class="mt-2 text-sm leading-relaxed text-gray-500">
            Boolean and multivariate flags. Toggle features instantly across
            environments without redeploying.
          </p>
          <!-- Mini toggle visual -->
          <div class="mt-5 flex items-center gap-3">
            <div class="flex items-center gap-2">
              <div class="h-5 w-9 rounded-full bg-blue-500 p-0.5">
                <div
                  class="h-4 w-4 translate-x-4 rounded-full bg-white shadow-sm"
                ></div>
              </div>
              <span class="font-code text-xs text-gray-400">new_dashboard</span>
            </div>
            <div class="flex items-center gap-2">
              <div class="h-5 w-9 rounded-full bg-gray-200 p-0.5">
                <div class="h-4 w-4 rounded-full bg-white shadow-sm"></div>
              </div>
              <span class="font-code text-xs text-gray-400">dark_mode</span>
            </div>
          </div>
        </div>

        <!-- Targeting Rules -->
        <div class="feature-card group rounded-2xl border border-gray-100 p-6">
          <div
            class="mb-4 flex h-11 w-11 items-center justify-center rounded-xl bg-emerald-50 text-emerald-600 transition-colors group-hover:bg-emerald-100"
          >
            <svg
              class="h-5 w-5"
              fill="none"
              stroke="currentColor"
              stroke-width="1.75"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"
              />
            </svg>
          </div>
          <h3 class="font-heading text-base font-semibold text-gray-900">
            Targeting Rules
          </h3>
          <p class="mt-2 text-sm leading-relaxed text-gray-500">
            Target by user attributes, segments, and custom properties. Deliver
            the right experience to the right users.
          </p>
          <div class="mt-5 space-y-1.5">
            <div
              class="flex items-center gap-2 rounded-md bg-gray-50 px-2.5 py-1.5"
            >
              <span class="font-code text-xs text-emerald-600">IF</span>
              <span class="font-code text-xs text-gray-500"
                >plan = "enterprise"</span
              >
            </div>
            <div
              class="flex items-center gap-2 rounded-md bg-gray-50 px-2.5 py-1.5"
            >
              <span class="font-code text-xs text-blue-600">AND</span>
              <span class="font-code text-xs text-gray-500"
                >country IN ["US","CA"]</span
              >
            </div>
          </div>
        </div>

        <!-- Percentage Rollouts -->
        <div class="feature-card group rounded-2xl border border-gray-100 p-6">
          <div
            class="mb-4 flex h-11 w-11 items-center justify-center rounded-xl bg-violet-50 text-violet-600 transition-colors group-hover:bg-violet-100"
          >
            <svg
              class="h-5 w-5"
              fill="none"
              stroke="currentColor"
              stroke-width="1.75"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"
              />
            </svg>
          </div>
          <h3 class="font-heading text-base font-semibold text-gray-900">
            Percentage Rollouts
          </h3>
          <p class="mt-2 text-sm leading-relaxed text-gray-500">
            Gradually expose features to your user base. Reduce risk with
            incremental, controlled rollouts.
          </p>
          <!-- Mini rollout bar -->
          <div class="mt-5">
            <div class="flex items-center justify-between text-xs">
              <span class="font-code text-gray-400">Rollout</span>
              <span class="font-heading font-semibold text-violet-600">65%</span
              >
            </div>
            <div
              class="mt-1.5 h-2 w-full overflow-hidden rounded-full bg-gray-100"
            >
              <div
                class="h-full rounded-full"
                style="width: 65%; background: linear-gradient(90deg, #7c3aed, #8b5cf6)"
              ></div>
            </div>
          </div>
        </div>

        <!-- Multi-Environment -->
        <div class="feature-card group rounded-2xl border border-gray-100 p-6">
          <div
            class="mb-4 flex h-11 w-11 items-center justify-center rounded-xl bg-amber-50 text-amber-600 transition-colors group-hover:bg-amber-100"
          >
            <svg
              class="h-5 w-5"
              fill="none"
              stroke="currentColor"
              stroke-width="1.75"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9"
              />
            </svg>
          </div>
          <h3 class="font-heading text-base font-semibold text-gray-900">
            Multi-Environment
          </h3>
          <p class="mt-2 text-sm leading-relaxed text-gray-500">
            Manage flags across development, staging, and production. Test
            safely before going live.
          </p>
          <div class="mt-5 flex items-center gap-2">
            <span
              class="rounded-md bg-emerald-50 px-2 py-1 font-code text-xs font-medium text-emerald-700"
              >production</span
            >
            <span
              class="rounded-md bg-amber-50 px-2 py-1 font-code text-xs font-medium text-amber-700"
              >staging</span
            >
            <span
              class="rounded-md bg-sky-50 px-2 py-1 font-code text-xs font-medium text-sky-700"
              >development</span
            >
          </div>
        </div>

        <!-- A/B Experiments -->
        <div class="feature-card group rounded-2xl border border-gray-100 p-6">
          <div
            class="mb-4 flex h-11 w-11 items-center justify-center rounded-xl bg-rose-50 text-rose-600 transition-colors group-hover:bg-rose-100"
          >
            <svg
              class="h-5 w-5"
              fill="none"
              stroke="currentColor"
              stroke-width="1.75"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 10.172V5L8 4z"
              />
            </svg>
          </div>
          <h3 class="font-heading text-base font-semibold text-gray-900">
            A/B Experiments
          </h3>
          <p class="mt-2 text-sm leading-relaxed text-gray-500">
            Run controlled experiments with variant allocation. Measure impact
            and make data-driven decisions.
          </p>
          <!-- Mini chart visual -->
          <div class="mt-5 flex items-end gap-1.5">
            <div class="flex flex-col items-center gap-1">
              <div class="h-8 w-6 rounded-sm bg-blue-400"></div>
              <span class="font-code text-[10px] text-gray-400">A</span>
            </div>
            <div class="flex flex-col items-center gap-1">
              <div class="h-11 w-6 rounded-sm bg-emerald-400"></div>
              <span class="font-code text-[10px] text-gray-400">B</span>
            </div>
            <div class="ml-2 font-code text-xs text-emerald-600">+12.4%</div>
          </div>
        </div>

        <!-- Audit Log -->
        <div class="feature-card group rounded-2xl border border-gray-100 p-6">
          <div
            class="mb-4 flex h-11 w-11 items-center justify-center rounded-xl bg-cyan-50 text-cyan-600 transition-colors group-hover:bg-cyan-100"
          >
            <svg
              class="h-5 w-5"
              fill="none"
              stroke="currentColor"
              stroke-width="1.75"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
              />
            </svg>
          </div>
          <h3 class="font-heading text-base font-semibold text-gray-900">
            Audit Log
          </h3>
          <p class="mt-2 text-sm leading-relaxed text-gray-500">
            Complete audit trail of every change. Know who changed what, when,
            and why. Full accountability.
          </p>
          <div class="mt-5 space-y-2">
            <div class="flex items-center gap-2 text-xs text-gray-400">
              <div class="h-1.5 w-1.5 rounded-full bg-emerald-400"></div>
              <span class="font-code">Flag created</span>
              <span class="ml-auto">2m ago</span>
            </div>
            <div class="flex items-center gap-2 text-xs text-gray-400">
              <div class="h-1.5 w-1.5 rounded-full bg-blue-400"></div>
              <span class="font-code">Rollout updated</span>
              <span class="ml-auto">5m ago</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- ============================================ -->
  <!-- DASHBOARD PREVIEW (gray)                     -->
  <!-- ============================================ -->
  <section class="bg-gray-50 py-24">
    <div class="mx-auto max-w-7xl px-6">
      <div class="mx-auto max-w-2xl text-center">
        <h2
          class="font-heading text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl"
        >
          See it in action
        </h2>
        <p class="mt-4 text-lg text-gray-500">
          A real product with realistic data — flags, environments, experiments,
          and a full audit trail. Log in and explore.
        </p>
      </div>

      <!-- Browser mockup -->
      <div class="mx-auto mt-14 max-w-5xl">
        <div
          class="dashboard-mockup overflow-hidden rounded-xl border border-gray-200 bg-white shadow-xl"
        >
          <!-- Browser chrome -->
          <div
            class="flex items-center gap-2 border-b border-gray-200 bg-gray-50 px-4 py-2.5"
          >
            <div class="flex items-center gap-1.5">
              <div class="h-2.5 w-2.5 rounded-full bg-red-400"></div>
              <div class="h-2.5 w-2.5 rounded-full bg-amber-400"></div>
              <div class="h-2.5 w-2.5 rounded-full bg-green-400"></div>
            </div>
            <div
              class="ml-4 flex-1 rounded-md bg-white px-3 py-1 text-center font-code text-xs text-gray-400"
            >
              flagdeck-app.workermill.com/dashboard
            </div>
          </div>
          <!-- Dashboard content mockup -->
          <div class="flex">
            <!-- Mini sidebar -->
            <div class="hidden w-44 border-r border-gray-100 p-3 sm:block">
              <div class="mb-4 flex items-center gap-2 px-2">
                <div
                  class="flex h-5 w-5 items-center justify-center rounded bg-blue-600 text-[8px] font-bold text-white"
                >
                  FD
                </div>
                <span class="text-xs font-semibold text-gray-800">FlagDeck</span
                >
              </div>
              <div class="space-y-0.5">
                <div
                  class="rounded-md bg-blue-50 px-2.5 py-1.5 text-xs font-medium text-blue-700"
                >
                  Dashboard
                </div>
                <div class="px-2.5 py-1.5 text-xs text-gray-500">Flags</div>
                <div class="px-2.5 py-1.5 text-xs text-gray-500">
                  Environments
                </div>
                <div class="px-2.5 py-1.5 text-xs text-gray-500">Segments</div>
                <div class="px-2.5 py-1.5 text-xs text-gray-500">
                  Experiments
                </div>
                <div class="px-2.5 py-1.5 text-xs text-gray-500">Audit Log</div>
              </div>
            </div>
            <!-- Main content -->
            <div class="flex-1 p-5">
              <!-- Stats row -->
              <div class="grid grid-cols-2 gap-3 sm:grid-cols-4">
                <div class="rounded-lg border border-gray-100 p-3">
                  <div class="text-xs text-gray-400">Total Flags</div>
                  <div
                    class="mt-0.5 font-heading text-xl font-bold text-gray-900"
                  >
                    10
                  </div>
                </div>
                <div class="rounded-lg border border-gray-100 p-3">
                  <div class="text-xs text-gray-400">Active</div>
                  <div
                    class="mt-0.5 font-heading text-xl font-bold text-emerald-600"
                  >
                    7
                  </div>
                </div>
                <div class="rounded-lg border border-gray-100 p-3">
                  <div class="text-xs text-gray-400">Environments</div>
                  <div
                    class="mt-0.5 font-heading text-xl font-bold text-gray-900"
                  >
                    3
                  </div>
                </div>
                <div class="rounded-lg border border-gray-100 p-3">
                  <div class="text-xs text-gray-400">Experiments</div>
                  <div
                    class="mt-0.5 font-heading text-xl font-bold text-blue-600"
                  >
                    2
                  </div>
                </div>
              </div>
              <!-- Mini table -->
              <div class="mt-4 rounded-lg border border-gray-100">
                <div
                  class="border-b border-gray-50 px-3 py-2 text-xs font-medium text-gray-500"
                >
                  Recent Activity
                </div>
                <div class="divide-y divide-gray-50">
                  <div class="flex items-center gap-3 px-3 py-2">
                    <div class="h-1.5 w-1.5 rounded-full bg-emerald-400"></div>
                    <span class="flex-1 text-xs text-gray-600"
                      >Flag <span class="font-code text-blue-600"
                        >new_dashboard</span
                      > enabled in production</span
                    >
                    <span class="text-[10px] text-gray-400">2m ago</span>
                  </div>
                  <div class="flex items-center gap-3 px-3 py-2">
                    <div class="h-1.5 w-1.5 rounded-full bg-blue-400"></div>
                    <span class="flex-1 text-xs text-gray-600"
                      >Experiment <span class="font-code text-blue-600"
                        >checkout_redesign</span
                      > started</span
                    >
                    <span class="text-[10px] text-gray-400">15m ago</span>
                  </div>
                  <div class="flex items-center gap-3 px-3 py-2">
                    <div class="h-1.5 w-1.5 rounded-full bg-amber-400"></div>
                    <span class="flex-1 text-xs text-gray-600"
                      >Rollout for <span class="font-code text-blue-600"
                        >beta_features</span
                      > set to 65%</span
                    >
                    <span class="text-[10px] text-gray-400">1h ago</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <!-- CTA below mockup -->
        <div class="mt-8 text-center">
          <a
            href="/login"
            class="hero-cta-primary inline-flex items-center gap-2 rounded-lg px-6 py-3 text-sm font-semibold text-white shadow-lg shadow-blue-500/25 transition-all hover:shadow-blue-500/40"
          >
            Explore the Live Demo
            <svg
              class="h-4 w-4"
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
          <p class="mt-3 font-code text-xs text-gray-400">
            demo@workermill.com &middot; demo1234
          </p>
        </div>
      </div>
    </div>
  </section>

  <!-- ============================================ -->
  <!-- BUILT BY AI SECTION (dark)                   -->
  <!-- ============================================ -->
  <section class="built-by-section relative overflow-hidden py-24">
    <!-- Subtle gradient -->
    <div
      class="pointer-events-none absolute left-1/2 top-1/2 h-[400px] w-[600px] -translate-x-1/2 -translate-y-1/2 rounded-full opacity-20 blur-[120px]"
      style="background: radial-gradient(circle, #3b82f6, transparent 70%)"
    ></div>

    <div class="relative z-10 mx-auto max-w-7xl px-6">
      <div class="mx-auto max-w-2xl text-center">
        <div
          class="mb-6 inline-flex items-center gap-2 rounded-full border border-white/10 bg-white/5 px-3.5 py-1.5 text-sm text-white/60"
        >
          <svg
            class="h-4 w-4 text-blue-400"
            fill="currentColor"
            viewBox="0 0 24 24"
          >
            <path d="M13 10V3L4 14h7v7l9-11h-7z" />
          </svg>
          Built entirely by AI
        </div>
        <h2
          class="font-heading text-3xl font-bold tracking-tight text-white sm:text-4xl"
        >
          From zero to production
          <span class="hero-gradient-text">in hours, not months</span>
        </h2>
        <p class="mt-4 text-lg leading-relaxed text-white/50">
          This entire application — API, frontend, database schemas, tests, CI
          pipeline, and this landing page — was built autonomously by
          <a
            href="https://workermill.com"
            target="_blank"
            rel="noopener noreferrer"
            class="text-blue-400 transition-colors hover:text-blue-300"
            >WorkerMill</a
          > AI coding agents.
        </p>
      </div>

      <!-- Build stats -->
      <div
        class="mx-auto mt-14 grid max-w-4xl grid-cols-2 gap-4 sm:grid-cols-5"
      >
        <div
          class="rounded-xl border border-white/[0.06] bg-white/[0.03] p-5 text-center backdrop-blur"
        >
          <div class="font-heading text-2xl font-bold text-white">4.5</div>
          <div class="mt-1 text-xs text-white/40">Hours</div>
        </div>
        <div
          class="rounded-xl border border-white/[0.06] bg-white/[0.03] p-5 text-center backdrop-blur"
        >
          <div class="font-heading text-2xl font-bold text-white">$64</div>
          <div class="mt-1 text-xs text-white/40">Total Cost</div>
        </div>
        <div
          class="rounded-xl border border-white/[0.06] bg-white/[0.03] p-5 text-center backdrop-blur"
        >
          <div class="font-heading text-2xl font-bold text-white">20,857</div>
          <div class="mt-1 text-xs text-white/40">Lines of Code</div>
        </div>
        <div
          class="rounded-xl border border-white/[0.06] bg-white/[0.03] p-5 text-center backdrop-blur"
        >
          <div class="font-heading text-2xl font-bold text-white">97</div>
          <div class="mt-1 text-xs text-white/40">Files</div>
        </div>
        <div
          class="col-span-2 rounded-xl border border-white/[0.06] bg-white/[0.03] p-5 text-center backdrop-blur sm:col-span-1"
        >
          <div class="font-heading text-2xl font-bold text-white">67</div>
          <div class="mt-1 text-xs text-white/40">Commits</div>
        </div>
      </div>

      <!-- Tech stack -->
      <div
        class="mt-12 flex flex-wrap items-center justify-center gap-x-6 gap-y-3"
      >
        <span class="text-xs text-white/25">Built with</span>
        <span
          class="rounded-md border border-white/[0.06] bg-white/[0.03] px-3 py-1.5 font-code text-xs text-white/50"
          >Go 1.24</span
        >
        <span
          class="rounded-md border border-white/[0.06] bg-white/[0.03] px-3 py-1.5 font-code text-xs text-white/50"
          >SvelteKit 2</span
        >
        <span
          class="rounded-md border border-white/[0.06] bg-white/[0.03] px-3 py-1.5 font-code text-xs text-white/50"
          >MongoDB 7</span
        >
        <span
          class="rounded-md border border-white/[0.06] bg-white/[0.03] px-3 py-1.5 font-code text-xs text-white/50"
          >Redis 7</span
        >
        <span
          class="rounded-md border border-white/[0.06] bg-white/[0.03] px-3 py-1.5 font-code text-xs text-white/50"
          >Docker</span
        >
        <span
          class="rounded-md border border-white/[0.06] bg-white/[0.03] px-3 py-1.5 font-code text-xs text-white/50"
          >GitHub Actions</span
        >
      </div>
    </div>
  </section>

  <!-- ============================================ -->
  <!-- CTA SECTION                                  -->
  <!-- ============================================ -->
  <section class="cta-section relative overflow-hidden py-24">
    <div class="relative z-10 mx-auto max-w-3xl px-6 text-center">
      <h2
        class="font-heading text-3xl font-bold tracking-tight text-white sm:text-4xl"
      >
        Ready to explore?
      </h2>
      <p class="mt-4 text-lg text-white/60">
        Log in with the demo credentials and explore a fully-loaded feature flag
        platform. No setup required.
      </p>
      <div
        class="mt-10 flex flex-col items-center justify-center gap-4 sm:flex-row"
      >
        <a
          href="/login"
          class="inline-flex w-full items-center justify-center gap-2 rounded-lg bg-white px-8 py-3.5 text-sm font-semibold text-gray-900 shadow-lg transition-all hover:bg-gray-100 sm:w-auto"
        >
          Try the Demo
          <svg
            class="h-4 w-4"
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
          class="inline-flex w-full items-center justify-center gap-2 rounded-lg border border-white/15 px-8 py-3.5 text-sm font-medium text-white/80 transition-all hover:border-white/25 hover:text-white sm:w-auto"
        >
          View Source Code
        </a>
      </div>
      <div
        class="mt-8 inline-flex items-center gap-3 rounded-lg border border-white/[0.06] bg-white/[0.03] px-5 py-3 backdrop-blur"
      >
        <span class="text-sm text-white/40">Demo login:</span>
        <code class="font-code text-sm text-white/70">demo@workermill.com</code>
        <span class="text-white/20">/</span>
        <code class="font-code text-sm text-white/70">demo1234</code>
      </div>
    </div>
  </section>

  <!-- ============================================ -->
  <!-- FOOTER                                       -->
  <!-- ============================================ -->
  <footer class="footer-section border-t border-white/[0.04]">
    <div
      class="mx-auto flex max-w-7xl flex-col items-center justify-between gap-4 px-6 py-8 sm:flex-row"
    >
      <div class="flex items-center gap-2.5">
        <div
          class="flex h-6 w-6 items-center justify-center rounded-md text-[9px] font-bold text-white"
          style="background: linear-gradient(135deg, #3b82f6, #6366f1)"
        >
          FD
        </div>
        <span class="text-sm text-white/30">
          FlagDeck &mdash; Built by
          <a
            href="https://workermill.com"
            target="_blank"
            rel="noopener noreferrer"
            class="text-blue-400/70 transition-colors hover:text-blue-400"
          >
            WorkerMill
          </a>
        </span>
      </div>
      <div class="flex items-center gap-6">
        <a
          href="https://github.com/workermill-examples/flagdeck"
          target="_blank"
          rel="noopener noreferrer"
          class="text-sm text-white/30 transition-colors hover:text-white/50"
        >
          GitHub
        </a>
        <a
          href="https://workermill.com"
          target="_blank"
          rel="noopener noreferrer"
          class="text-sm text-white/30 transition-colors hover:text-white/50"
        >
          WorkerMill
        </a>
      </div>
    </div>
  </footer>
</div>

<style>
  .landing-page {
    font-family:
      "DM Sans",
      ui-sans-serif,
      system-ui,
      -apple-system,
      sans-serif;
  }

  .font-heading {
    font-family:
      "Sora",
      ui-sans-serif,
      system-ui,
      -apple-system,
      sans-serif;
  }

  .font-code {
    font-family: "JetBrains Mono", ui-monospace, "Cascadia Code", monospace;
  }

  .hero-section {
    background: #06080f;
  }

  .hero-gradient-text {
    background: linear-gradient(135deg, #60a5fa 0%, #818cf8 50%, #a78bfa 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .hero-cta-primary {
    background: linear-gradient(135deg, #3b82f6, #6366f1);
  }

  .hero-cta-primary:hover {
    background: linear-gradient(135deg, #2563eb, #4f46e5);
  }

  .code-window {
    background: #0d1117;
  }

  .feature-card {
    transition: all 0.2s ease;
  }

  .feature-card:hover {
    border-color: #d1d5db;
    box-shadow:
      0 4px 6px -1px rgba(0, 0, 0, 0.05),
      0 2px 4px -2px rgba(0, 0, 0, 0.05);
    transform: translateY(-1px);
  }

  .dashboard-mockup {
    transition: all 0.3s ease;
  }

  .dashboard-mockup:hover {
    box-shadow:
      0 25px 50px -12px rgba(0, 0, 0, 0.15),
      0 0 0 1px rgba(0, 0, 0, 0.05);
    transform: translateY(-2px);
  }

  .built-by-section {
    background: #06080f;
  }

  .cta-section {
    background: linear-gradient(180deg, #06080f 0%, #0c1220 100%);
  }

  .footer-section {
    background: #06080f;
  }
</style>
