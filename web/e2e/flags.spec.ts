import { test, expect } from "@playwright/test";

// Helper function to login before each test
async function login(page: any) {
  await page.goto("/login");
  await page.click('button[type="submit"]');
  await expect(page).toHaveURL("/dashboard", { timeout: 15000 });
}

test.describe("Flags Page", () => {
  test.beforeEach(async ({ page }) => {
    // Navigate first so localStorage is accessible (about:blank denies access)
    await page.goto("/login");
    await page.context().clearCookies();
    await page.evaluate(() => localStorage.clear());

    // Login with demo credentials
    await login(page);
  });

  test("displays flags list with 10+ seeded flags", async ({ page }) => {
    await page.goto("/flags");

    // Verify page title and header
    await expect(page).toHaveTitle(/Feature Flags - FlagDeck/);
    await expect(page.locator("h1")).toContainText("Feature Flags");

    // Wait for data to load
    await page.waitForLoadState("networkidle");

    // Verify flags table is displayed
    await expect(page.locator("table")).toBeVisible();

    // Count flag rows (excluding header)
    const flagRows = page.locator("tbody tr");
    const flagCount = await flagRows.count();
    expect(flagCount).toBeGreaterThanOrEqual(10); // Spec mentions 10+ flags

    // Verify table headers
    await expect(page.locator('th:has-text("Flag")')).toBeVisible();
    await expect(page.locator('th:has-text("Type")')).toBeVisible();
    await expect(page.locator('th:has-text("Status")')).toBeVisible();
    await expect(
      page.locator('th:has-text("Environment Toggles")'),
    ).toBeVisible();
    await expect(page.locator('th:has-text("Tags")')).toBeVisible();

    // Verify some realistic flag names from the seeded data
    await expect(
      page
        .locator("text=dark-mode")
        .or(
          page
            .locator("text=checkout-flow")
            .or(page.locator("text=search-algorithm")),
        )
        .first(),
    ).toBeVisible();
  });

  test("search functionality filters flags correctly", async ({ page }) => {
    await page.goto("/flags");
    await page.waitForLoadState("networkidle");

    // Get initial flag count
    const initialCount = await page.locator("tbody tr").count();
    expect(initialCount).toBeGreaterThan(0);

    // Search for a specific term
    await page.fill("#search", "dark");

    // Wait for filtering to apply
    await page.waitForTimeout(500);

    // Verify filtered results
    const filteredCount = await page.locator("tbody tr").count();
    expect(filteredCount).toBeLessThanOrEqual(initialCount);

    // Verify the visible flags contain the search term
    const visibleFlags = page.locator("tbody tr");
    const firstFlag = await visibleFlags
      .first()
      .locator(".text-sm.font-medium.text-gray-900, .text-sm.text-gray-500")
      .allTextContents();
    expect(firstFlag.some((text) => text.toLowerCase().includes("dark"))).toBe(
      true,
    );

    // Clear search
    await page.click('button:has-text("Clear Filters")');

    // Verify count returns to original
    const clearedCount = await page.locator("tbody tr").count();
    expect(clearedCount).toBe(initialCount);
  });

  test("type filter works correctly", async ({ page }) => {
    await page.goto("/flags");
    await page.waitForLoadState("networkidle");

    // Filter by boolean type
    await page.selectOption("#type-filter", "boolean");
    await page.waitForTimeout(500);

    // Verify all visible flags are boolean type
    const booleanBadges = page.locator(
      '.bg-blue-100.text-blue-800:has-text("boolean")',
    );
    const visibleRows = page.locator("tbody tr");
    const visibleCount = await visibleRows.count();
    const booleanCount = await booleanBadges.count();

    expect(booleanCount).toBeGreaterThan(0);
    expect(booleanCount).toBe(visibleCount);
  });

  test("can create a new flag", async ({ page }) => {
    await page.goto("/flags");

    // Click create flag button
    await page.click('a[href="/flags/create"]');
    await expect(page).toHaveURL("/flags/create");

    // Verify create page is loaded
    await expect(page.locator("h1")).toContainText("Create Feature Flag");

    // Wait for environments to load
    await page.waitForLoadState("networkidle");

    // Fill in flag details
    const flagName = `Test Flag ${Date.now()}`;
    await page.fill("#name", flagName);

    // Verify key is auto-generated
    const generatedKey = await page.locator("#key").inputValue();
    expect(generatedKey).toBeTruthy();
    expect(generatedKey.toLowerCase().replace(/[^a-z0-9-]/g, "")).toBeTruthy();

    await page.fill("#description", "A test flag created by E2E tests");
    await page.fill("#tags", "test, e2e");

    // Select string type and set value
    await page.selectOption("#type", "string");
    await page.fill("#defaultValue", "test-value");

    // Submit the form
    await page.click('button[type="submit"]');

    // Wait for redirect to flag detail page
    await page.waitForURL(/\/flags\/.*/, { timeout: 10000 });

    // Verify we're on the flag detail page
    await expect(page.locator("h1")).toContainText(flagName);
  });

  test("can toggle flag environments", async ({ page }) => {
    await page.goto("/flags");
    await page.waitForLoadState("networkidle");

    // Find the first flag row
    const firstFlagRow = page.locator("tbody tr").first();
    const flagKey = await firstFlagRow
      .locator(".text-sm.text-gray-500")
      .first()
      .textContent();

    // Find a toggle switch in the first row
    const toggleSwitch = firstFlagRow.locator('input[type="checkbox"]').first();
    const initialState = await toggleSwitch.isChecked();

    // Click the toggle
    await toggleSwitch.click();

    // Wait for the API call to complete
    await page.waitForTimeout(1000);

    // Verify the toggle state changed
    const newState = await toggleSwitch.isChecked();
    expect(newState).not.toBe(initialState);

    // Toggle it back
    await toggleSwitch.click();
    await page.waitForTimeout(1000);

    // Verify it's back to original state
    const finalState = await toggleSwitch.isChecked();
    expect(finalState).toBe(initialState);
  });

  test("can edit flag targeting rules", async ({ page }) => {
    await page.goto("/flags");
    await page.waitForLoadState("networkidle");

    // Click on the first flag's edit button
    await page.click('tbody tr:first-child button:has-text("Edit")');

    // Wait for flag detail page to load
    await page.waitForURL(/\/flags\/.*/, { timeout: 10000 });

    // Verify we're on flag detail page
    await expect(page.locator("h1")).toBeVisible();

    // Look for targeting rule builder or environment configuration
    // The targeting rules might be in a tab or section
    const targetingSection = page
      .locator("text=Targeting")
      .or(page.locator("text=Rules"))
      .or(page.locator("text=Environment"));

    if (await targetingSection.isVisible()) {
      await targetingSection.click();

      // Look for add rule button or rule builder
      const addRuleButton = page
        .locator('button:has-text("Add Rule")')
        .or(page.locator('button:has-text("Add Condition")'));
      if (await addRuleButton.isVisible()) {
        await addRuleButton.click();

        // Verify rule builder interface appeared
        await expect(
          page
            .locator(
              'input[placeholder*="property"], select, input[placeholder*="value"]',
            )
            .first(),
        ).toBeVisible();
      }
    }

    // Test rollout slider if present
    const rolloutSlider = page.locator('input[type="range"]');
    if (await rolloutSlider.isVisible()) {
      await rolloutSlider.fill("75");
    }
  });

  test("flag status badges display correctly", async ({ page }) => {
    await page.goto("/flags");
    await page.waitForLoadState("networkidle");

    // Verify status badges are present
    const activeFlags = page.locator(
      '.bg-green-100.text-green-800:has-text("Active")',
    );
    const inactiveFlags = page.locator(
      '.bg-red-100.text-red-800:has-text("Inactive")',
    );

    // At least some flags should be active (based on seeded data)
    const activeCount = await activeFlags.count();
    expect(activeCount).toBeGreaterThan(0);

    // Verify type badges
    const booleanBadges = page.locator(
      '.bg-blue-100.text-blue-800:has-text("boolean")',
    );
    const stringBadges = page.locator(
      '.bg-green-100.text-green-800:has-text("string")',
    );
    const numberBadges = page.locator(
      '.bg-purple-100.text-purple-800:has-text("number")',
    );

    const totalTypeBadges =
      (await booleanBadges.count()) +
      (await stringBadges.count()) +
      (await numberBadges.count());
    expect(totalTypeBadges).toBeGreaterThan(0);
  });

  test("environment status dots are visible", async ({ page }) => {
    await page.goto("/flags");
    await page.waitForLoadState("networkidle");

    // Verify environment status dots (colored circles) are present
    const statusDots = page.locator(".w-3.h-3.rounded-full");
    const dotsCount = await statusDots.count();

    // Should have dots for all environments across all flags
    // With 3 environments and 10+ flags, expect 30+ dots
    expect(dotsCount).toBeGreaterThan(30);

    // Verify environment keys are shown
    await expect(
      page
        .locator("text=production")
        .or(page.locator("text=staging"))
        .or(page.locator("text=development"))
        .first(),
    ).toBeVisible();
  });

  test("displays flag summary correctly", async ({ page }) => {
    await page.goto("/flags");
    await page.waitForLoadState("networkidle");

    // Verify summary text at bottom
    const summaryText = page
      .locator("text=Showing")
      .or(page.locator("text=of"));
    await expect(summaryText.first()).toBeVisible();

    // Get the total count from summary
    const summaryContent = await page
      .locator(".text-sm.text-gray-500.text-center")
      .textContent();
    expect(summaryContent).toMatch(/\d+.*of.*\d+.*flags/i);
  });

  test("handles error states gracefully", async ({ page }) => {
    // Intercept flags API to simulate error
    await page.route("**/api/v1/flags", (route) => route.abort());

    await page.goto("/flags");

    // Wait for error to appear
    await expect(page.locator(".bg-red-50")).toBeVisible({ timeout: 10000 });
    await expect(page.locator(".text-red-700")).toContainText(
      /Failed|Error|error/i,
    );

    // Verify retry button is available
    await expect(page.locator('button:has-text("Try again")')).toBeVisible();
  });

  test("displays loading state initially", async ({ page }) => {
    // Intercept API to delay response
    await page.route("**/api/v1/flags", async (route) => {
      await new Promise((resolve) => setTimeout(resolve, 2000));
      await route.continue();
    });

    await page.goto("/flags");

    // Verify loading spinner appears
    await expect(page.locator(".animate-spin")).toBeVisible();

    // Wait for loading to complete
    await page.waitForLoadState("networkidle");
    await expect(page.locator(".animate-spin")).not.toBeVisible({
      timeout: 15000,
    });
  });
});
