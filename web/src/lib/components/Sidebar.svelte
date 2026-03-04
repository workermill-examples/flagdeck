<script lang="ts">
import { page } from "$app/stores";
import { auth } from "$lib/auth.js";

interface NavItem {
  href: string;
  label: string;
  icon: string;
}

const navItems: NavItem[] = [
  { href: "/", label: "Dashboard", icon: "🏠" },
  { href: "/flags", label: "Flags", icon: "🏁" },
  { href: "/environments", label: "Environments", icon: "🌍" },
  { href: "/segments", label: "Segments", icon: "👥" },
  { href: "/experiments", label: "Experiments", icon: "🧪" },
  { href: "/audit-log", label: "Audit Log", icon: "📋" },
  { href: "/settings", label: "Settings", icon: "⚙️" },
];

async function handleLogout() {
  await auth.logout();
  window.location.href = "/login";
}
</script>

<aside class="w-64 bg-white border-r border-gray-200 flex flex-col h-full">
  <!-- Header with FlagDeck branding -->
  <div class="px-6 py-4 border-b border-gray-200">
    <div class="flex items-center space-x-2">
      <div class="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center text-white font-bold text-sm">
        FD
      </div>
      <h1 class="text-xl font-bold text-gray-900">FlagDeck</h1>
    </div>
  </div>

  <!-- Navigation -->
  <nav class="flex-1 px-4 py-4 space-y-2">
    {#each navItems as item}
      <a
        href={item.href}
        class="flex items-center space-x-3 px-3 py-2 text-sm font-medium rounded-lg transition-colors duration-200 {$page.url.pathname === item.href
          ? 'bg-blue-50 text-blue-700 border-r-2 border-blue-700'
          : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900'}"
      >
        <span class="text-lg">{item.icon}</span>
        <span>{item.label}</span>
      </a>
    {/each}
  </nav>

  <!-- User info and logout -->
  <div class="px-4 py-4 border-t border-gray-200">
    <div class="flex items-center justify-between">
      <div class="flex items-center space-x-3">
        <div class="w-8 h-8 bg-gray-300 rounded-full flex items-center justify-center text-gray-600 text-sm">
          {auth.user?.name?.charAt(0).toUpperCase() || "U"}
        </div>
        <div class="flex flex-col">
          <span class="text-sm font-medium text-gray-900">{auth.user?.name || "User"}</span>
          <span class="text-xs text-gray-500">{auth.user?.email || ""}</span>
        </div>
      </div>
      <button
        onclick={handleLogout}
        class="p-1 text-gray-400 hover:text-gray-600 transition-colors"
        title="Logout"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
        </svg>
      </button>
    </div>
  </div>
</aside>