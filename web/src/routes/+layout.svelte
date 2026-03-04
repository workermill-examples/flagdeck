<script lang="ts">
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import { auth } from "$lib/auth.svelte.js";
  import Sidebar from "$lib/components/Sidebar.svelte";
  import "../../app.css";

  interface Props {
    children: import("svelte").Snippet;
  }

  let { children }: Props = $props();

  // Pages that never show the sidebar (public or transitional)
  let isFullWidthPage = $derived(
    $page.url.pathname === "/login" || $page.url.pathname === "/",
  );

  // Landing page has its own background
  let isLandingPage = $derived($page.url.pathname === "/");

  // Auth guard - redirect to login if not authenticated and not on a public page
  onMount(async () => {
    await auth.checkAuth();

    if (!auth.isAuthenticated && !isFullWidthPage) {
      goto("/login", { replaceState: true });
    }
  });

  // Reactive redirect when auth state changes
  $effect(() => {
    if (!auth.isLoading && !auth.isAuthenticated && !isFullWidthPage) {
      goto("/login", { replaceState: true });
    }
  });
</script>

{#if auth.isLoading && !isFullWidthPage}
  <!-- Loading state while checking authentication -->
  <div class="min-h-screen flex items-center justify-center bg-gray-50">
    <div class="text-center">
      <div
        class="w-8 h-8 mx-auto mb-4 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"
      ></div>
      <p class="text-gray-600">Loading...</p>
    </div>
  </div>
{:else if isFullWidthPage}
  <!-- Full-width pages: landing page, login (no sidebar) -->
  <main class="min-h-screen {isLandingPage ? '' : 'bg-gray-50'}">
    {@render children()}
  </main>
{:else if auth.isAuthenticated}
  <!-- Authenticated layout with sidebar -->
  <div class="min-h-screen bg-gray-50 flex">
    <Sidebar />
    <main class="flex-1 overflow-hidden">
      {@render children()}
    </main>
  </div>
{:else}
  <!-- Fallback - should not normally be seen -->
  <div class="min-h-screen flex items-center justify-center bg-gray-50">
    <div class="text-center">
      <p class="text-gray-600">Redirecting to login...</p>
    </div>
  </div>
{/if}
