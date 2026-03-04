import { test, expect } from "@playwright/test";

// Helper function to login before each test
async function login(page: any) {
  await page.goto("/login");
  await page.waitForSelector("#email", { state: "visible" });
  await page.click('button[type="submit"]');
  await page.waitForURL("**/dashboard", { timeout: 30000 });
}

test.describe("Dashboard Page", () => {
  test.beforeEach(async ({ page }) => {
    // Navigate first so localStorage is accessible (about:blank denies access)
    await page.goto("/login");
    await page.context().clearCookies();
    await page.evaluate(() => localStorage.clear());

    // Login with demo credentials
    await login(page);
  });

  test("displays dashboard title and header", async ({ page }) => {
    await expect(page).toHaveTitle(/Dashboard - FlagDeck/);
    await expect(page.locator("main h1")).toContainText("Dashboard");

    // Verify refresh button is present
    await expect(page.locator('button:has-text("Refresh")')).toBeVisible();
  });

  test("displays stat cards with non-zero seeded values", async ({ page }) => {
    // Wait for data to load
    await page.waitForLoadState("networkidle");

    // Wait for loading to complete (stat cards should not show loading state)
    await expect(page.locator(".animate-pulse").first()).not.toBeVisible({
      timeout: 10000,
    });

    // Verify all stat cards are present
    await expect(
      page.locator("main").locator("text=Total Flags"),
    ).toBeVisible();
    await expect(
      page.locator("main").locator("text=Active Flags"),
    ).toBeVisible();
    await expect(
      page.locator("main").locator("text=Environments"),
    ).toBeVisible();
    await expect(
      page.locator("main").locator("text=Total Experiments"),
    ).toBeVisible();
    await expect(
      page.locator("main").locator("text=Running Experiments"),
    ).toBeVisible();

    // Helper to get stat card value by label
    async function getStatValue(label: string): Promise<number> {
      const card = page
        .locator(".bg-white.rounded-lg")
        .filter({ hasText: label });
      const value = await card.locator(".text-3xl").textContent();
      return parseInt(value || "0");
    }

    // Verify stat cards show non-zero values from seeded data
    expect(await getStatValue("Total Flags")).toBeGreaterThan(0);
    expect(await getStatValue("Active Flags")).toBeGreaterThan(0);
    expect(await getStatValue("Environments")).toBe(3);
    expect(await getStatValue("Total Experiments")).toBeGreaterThan(0);
    expect(await getStatValue("Running Experiments")).toBeGreaterThan(0);
  });

  test("displays recent activity timeline with seeded data", async ({
    page,
  }) => {
    // Wait for data to load
    await page.waitForLoadState("networkidle");

    // Verify Recent Activity section
    await expect(page.locator('h2:has-text("Recent Activity")')).toBeVisible();

    // Wait for audit entries to load
    await expect(page.locator(".animate-pulse").first()).not.toBeVisible({
      timeout: 10000,
    });

    // Verify audit entries are displayed (should not be empty)
    const auditEntries = page.locator(
      '[data-testid="audit-entry"], .flex.items-start.space-x-3',
    );
    await expect(auditEntries.first()).toBeVisible();

    // Verify audit entries show user emails and actions
    await expect(
      page.locator("main").locator("text=demo@workermill.com").first(),
    ).toBeVisible();

    // Verify "View all activity" link
    await expect(page.locator('a:has-text("View all activity")')).toBeVisible();
    await expect(page.locator('a[href="/audit-log"]')).toBeVisible();
  });

  test("displays flag status overview with seeded flags", async ({ page }) => {
    // Wait for data to load
    await page.waitForLoadState("networkidle");

    // Verify Flag Status Overview section
    await expect(
      page.locator('h2:has-text("Flag Status Overview")'),
    ).toBeVisible();

    // Wait for flags to load
    await expect(page.locator(".animate-pulse").first()).not.toBeVisible({
      timeout: 10000,
    });

    // Verify flags are displayed (should not be empty)
    const flagKeys = page.locator(".text-xs.text-gray-500.font-mono");
    const flagCount = await flagKeys.count();
    expect(flagCount).toBeGreaterThan(0);

    // Verify environment status dots are shown
    await expect(page.locator(".w-3.h-3.rounded-full").first()).toBeVisible();

    // Verify "View all flags" link
    await expect(page.locator('a:has-text("View all flags")')).toBeVisible();
    await expect(page.locator('a[href="/flags"]')).toBeVisible();
  });

  test("refresh button reloads dashboard data", async ({ page }) => {
    // Wait for initial load
    await page.waitForLoadState("networkidle");
    await expect(page.locator(".animate-pulse").first()).not.toBeVisible({
      timeout: 10000,
    });

    // Click refresh button
    await page.click('button:has-text("Refresh")');

    // Verify loading state appears briefly
    await expect(
      page.locator('button:has-text("Refresh") svg.animate-spin'),
    ).toBeVisible({ timeout: 5000 });

    // Wait for refresh to complete
    await page.waitForLoadState("networkidle");
    await expect(page.locator(".animate-pulse").first()).not.toBeVisible({
      timeout: 10000,
    });

    // Verify data is still displayed
    await expect(
      page.locator("main").locator("text=Total Flags"),
    ).toBeVisible();
    await expect(page.locator('h2:has-text("Recent Activity")')).toBeVisible();
  });

  test("handles error state gracefully", async ({ page }) => {
    // Intercept API calls to simulate server error
    await page.route("**/api/v1/flags", (route) => route.abort());
    await page.route("**/api/v1/environments", (route) => route.abort());
    await page.route("**/api/v1/experiments", (route) => route.abort());
    await page.route("**/api/v1/audit-log*", (route) => route.abort());

    // Reload the page to trigger the error
    await page.reload();

    // Wait for error to appear
    await expect(page.locator(".bg-red-50")).toBeVisible({ timeout: 10000 });
    await expect(page.locator(".text-red-800")).toContainText(
      /Failed to load|error|Error/i,
    );
  });

  test("navigation links work correctly", async ({ page }) => {
    // Wait for data to load
    await page.waitForLoadState("networkidle");

    // Test "View all activity" link
    await page.click('a[href="/audit-log"]');
    await expect(page).toHaveURL("/audit-log");

    // Go back to dashboard
    await page.goto("/dashboard");

    // Test "View all flags" link
    await page.click('a[href="/flags"]');
    await expect(page).toHaveURL("/flags");
  });

  test("shows loading states initially", async ({ page }) => {
    // Clear auth and login again to see loading state
    await page.context().clearCookies();
    await page.evaluate(() => localStorage.clear());

    // Intercept API calls to delay them
    await page.route("**/api/v1/**", async (route) => {
      await new Promise((resolve) => setTimeout(resolve, 1000));
      await route.continue();
    });

    await login(page);

    // Verify loading states are shown
    await expect(page.locator(".animate-pulse").first()).toBeVisible();
    await expect(
      page.locator('button:has-text("Refresh")[disabled]'),
    ).toBeVisible();

    // Wait for loading to complete
    await page.waitForLoadState("networkidle");
    await expect(page.locator(".animate-pulse").first()).not.toBeVisible({
      timeout: 15000,
    });
  });
});
