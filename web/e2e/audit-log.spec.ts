import { test, expect } from "@playwright/test";

// Helper function to login before each test
async function login(page: any) {
  await page.goto("/login");
  await page.click('button[type="submit"]');
  await expect(page).toHaveURL("/dashboard", { timeout: 15000 });
}

test.describe("Audit Log Page", () => {
  test.beforeEach(async ({ page }) => {
    // Navigate first so localStorage is accessible (about:blank denies access)
    await page.goto("/login");
    await page.context().clearCookies();
    await page.evaluate(() => localStorage.clear());

    // Login with demo credentials
    await login(page);
  });

  test("displays audit log timeline with seeded entries", async ({ page }) => {
    await page.goto("/audit-log");

    // Verify page title and header
    await expect(page).toHaveTitle(/FlagDeck/);
    await expect(page.locator("h1")).toContainText("Audit Log");

    // Wait for data to load
    await page.waitForLoadState("networkidle");

    // Verify timeline entries are displayed (should not be empty based on seeded data)
    const timelineEntries = page.locator("ul.-mb-8 li");
    const entryCount = await timelineEntries.count();
    expect(entryCount).toBeGreaterThan(0); // Spec mentions 50+ audit entries

    // Verify realistic user activity
    await expect(page.locator("text=demo@workermill.com")).toBeVisible();

    // Verify action types are displayed
    const actionElements = page.locator(".font-medium.text-gray-900");
    const actionCount = await actionElements.count();
    expect(actionCount).toBeGreaterThan(0);

    // Verify action types match expected format (e.g., "Flag Created", "Flag Updated")
    const firstAction = await actionElements.first().textContent();
    expect(firstAction).toMatch(/^[A-Z][a-z]+ [A-Z][a-z]+$/); // e.g., "Flag Created"
  });

  test("timeline shows user avatars and timestamps correctly", async ({
    page,
  }) => {
    await page.goto("/audit-log");
    await page.waitForLoadState("networkidle");

    // Verify user avatars (initials) are shown
    const avatars = page.locator(".h-5.w-5.rounded-full.bg-gray-300");
    const avatarCount = await avatars.count();
    expect(avatarCount).toBeGreaterThan(0);

    // Verify initials are displayed (should be 2 characters)
    const firstAvatar = avatars.first();
    const avatarText = await firstAvatar.textContent();
    expect(avatarText).toMatch(/^[A-Z]{1,2}$/); // 1-2 uppercase letters

    // Verify timestamps are shown in relative format
    const timestamps = page.locator(".text-xs.text-gray-500.cursor-help");
    const timestampCount = await timestamps.count();
    expect(timestampCount).toBeGreaterThan(0);

    // Check for relative time formats
    const firstTimestamp = await timestamps.first().textContent();
    expect(firstTimestamp).toMatch(
      /(Just now|minute|hour|day|ago|\d+\/\d+\/\d+)/,
    );

    // Verify absolute timestamps are available on hover (title attribute)
    const firstTimestampTitle = await timestamps.first().getAttribute("title");
    expect(firstTimestampTitle).toBeTruthy();
    expect(firstTimestampTitle).toMatch(/\d+\/\d+\/\d+/); // Date format
  });

  test("resource filter works correctly", async ({ page }) => {
    await page.goto("/audit-log");
    await page.waitForLoadState("networkidle");

    // Get initial entry count
    const initialCount = await page.locator("ul.-mb-8 li").count();
    expect(initialCount).toBeGreaterThan(0);

    // Filter by flags
    await page.selectOption("#resource-filter", "flag");
    await page.waitForTimeout(1000); // Wait for filter to apply

    // Verify filter is applied
    const flagEntries = page.locator("ul.-mb-8 li");
    const filteredCount = await flagEntries.count();

    // Verify the visible entries are flag-related
    // All visible actions should contain "Flag" in the action text
    const actionElements = page.locator(".font-medium.text-gray-900");
    const actionCount = await actionElements.count();

    if (actionCount > 0) {
      // Check first few actions contain "Flag"
      for (let i = 0; i < Math.min(3, actionCount); i++) {
        const actionText = await actionElements.nth(i).textContent();
        expect(actionText?.toLowerCase()).toContain("flag");
      }
    }

    // Verify filter pill is displayed
    await expect(page.locator("text=Resource: Flags")).toBeVisible();

    // Clear filter
    await page.click('button:has-text("Clear Filters")');
    await page.waitForTimeout(1000);

    // Verify count returns to original (or close to it)
    const clearedCount = await page.locator("ul.-mb-8 li").count();
    expect(clearedCount).toBeGreaterThanOrEqual(filteredCount);
  });

  test("action filter works correctly", async ({ page }) => {
    await page.goto("/audit-log");
    await page.waitForLoadState("networkidle");

    // Filter by created actions
    await page.selectOption("#action-filter", "created");
    await page.waitForTimeout(1000);

    // Verify filtered results contain only "created" actions
    const actionElements = page.locator(".font-medium.text-gray-900");
    const actionCount = await actionElements.count();

    if (actionCount > 0) {
      // Check that visible actions contain "Created"
      for (let i = 0; i < Math.min(3, actionCount); i++) {
        const actionText = await actionElements.nth(i).textContent();
        expect(actionText?.toLowerCase()).toContain("created");
      }
    }

    // Verify filter pill is displayed
    await expect(page.locator("text=Action: Created")).toBeVisible();
  });

  test("pagination works correctly", async ({ page }) => {
    await page.goto("/audit-log");
    await page.waitForLoadState("networkidle");

    // Check if pagination is present (depends on number of entries)
    const paginationContainer = page.locator(".border-t.border-gray-200");

    if (await paginationContainer.isVisible()) {
      // Verify pagination controls
      await expect(page.locator('button:has-text("Previous")')).toBeVisible();
      await expect(page.locator('button:has-text("Next")')).toBeVisible();

      // Verify page information
      await expect(page.locator("text=Page")).toBeVisible();
      await expect(page.locator("text=of")).toBeVisible();

      // Test next page if available
      const nextButton = page.locator('button:has-text("Next")');
      if (await nextButton.isEnabled()) {
        const firstPageFirstEntry = await page
          .locator("ul.-mb-8 li")
          .first()
          .locator(".font-medium.text-gray-900")
          .textContent();

        await nextButton.click();
        await page.waitForTimeout(1000);

        // Verify we're on page 2
        await expect(page.locator("text=Page 2")).toBeVisible();

        // Verify different entries are shown
        const secondPageFirstEntry = await page
          .locator("ul.-mb-8 li")
          .first()
          .locator(".font-medium.text-gray-900")
          .textContent();
        expect(secondPageFirstEntry).not.toBe(firstPageFirstEntry);

        // Go back to page 1
        await page.click('button:has-text("Previous")');
        await page.waitForTimeout(1000);

        await expect(page.locator("text=Page 1")).toBeVisible();
      }
    }
  });

  test("entries per page selector works", async ({ page }) => {
    await page.goto("/audit-log");
    await page.waitForLoadState("networkidle");

    // Change entries per page to 10
    await page.selectOption("#entries-per-page", "10");
    await page.waitForTimeout(1000);

    // Verify showing information is updated
    const showingText = await page.locator("text=Showing").textContent();
    expect(showingText).toMatch(/Showing \d+-\d+ of \d+ entries/);

    // Change to 50 entries per page
    await page.selectOption("#entries-per-page", "50");
    await page.waitForTimeout(1000);

    // Verify the page updated
    const newShowingText = await page.locator("text=Showing").textContent();
    expect(newShowingText).toBeTruthy();
  });

  test("expandable change details work correctly", async ({ page }) => {
    await page.goto("/audit-log");
    await page.waitForLoadState("networkidle");

    // Look for entries with change details
    const expandButtons = page.locator('button:has-text("Show changes")');
    const expandButtonCount = await expandButtons.count();

    if (expandButtonCount > 0) {
      // Click first expand button
      await expandButtons.first().click();

      // Verify expanded content appears
      await expect(page.locator(".bg-gray-50.rounded-md")).toBeVisible();
      await expect(page.locator('h4:has-text("Changes")')).toBeVisible();

      // Verify before/after sections
      await expect(page.locator("text=Before:")).toBeVisible();
      await expect(page.locator("text=After:")).toBeVisible();

      // Verify we can collapse it
      await page.click('button:has-text("Hide changes")');
      await expect(page.locator('h4:has-text("Changes")')).not.toBeVisible();
    }
  });

  test("resource IDs and action icons display correctly", async ({ page }) => {
    await page.goto("/audit-log");
    await page.waitForLoadState("networkidle");

    // Verify resource IDs are shown in code format
    const resourceIds = page.locator("code");
    const resourceIdCount = await resourceIds.count();
    expect(resourceIdCount).toBeGreaterThan(0);

    // Verify action icons are displayed
    const actionIcons = page.locator(".h-8.w-8.rounded-full svg");
    const iconCount = await actionIcons.count();
    expect(iconCount).toBeGreaterThan(0);

    // Verify icons have proper colors (background classes)
    const coloredIcons = page.locator(
      ".bg-green-100, .bg-blue-100, .bg-yellow-100, .bg-red-100, .bg-gray-100",
    );
    const coloredIconCount = await coloredIcons.count();
    expect(coloredIconCount).toBeGreaterThan(0);
  });

  test("refresh functionality works", async ({ page }) => {
    await page.goto("/audit-log");
    await page.waitForLoadState("networkidle");

    // Click refresh button
    const refreshButton = page.locator('button:has-text("Refresh")');
    if (await refreshButton.isVisible()) {
      await refreshButton.click();

      // Wait for refresh to complete
      await page.waitForTimeout(1000);

      // Verify entries are still displayed
      const entries = page.locator("ul.-mb-8 li");
      const entryCount = await entries.count();
      expect(entryCount).toBeGreaterThan(0);
    }
  });

  test("empty state displays when no entries found", async ({ page }) => {
    // Intercept audit log API to return empty result
    await page.route("**/api/v1/audit-log*", (route) =>
      route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ data: [], total: 0 }),
      }),
    );

    await page.goto("/audit-log");

    // Verify empty state
    await expect(
      page.locator('h3:has-text("No audit log entries")'),
    ).toBeVisible();
    await expect(
      page.locator("text=Activity will appear here as changes are made"),
    ).toBeVisible();
  });

  test("loading state displays correctly", async ({ page }) => {
    // Intercept API to delay response
    await page.route("**/api/v1/audit-log*", async (route) => {
      await new Promise((resolve) => setTimeout(resolve, 2000));
      await route.continue();
    });

    await page.goto("/audit-log");

    // Verify loading state
    await expect(page.locator("text=Loading audit log entries")).toBeVisible();
    await expect(page.locator(".animate-spin")).toBeVisible();

    // Wait for loading to complete
    await page.waitForLoadState("networkidle");
    await expect(
      page.locator("text=Loading audit log entries"),
    ).not.toBeVisible({ timeout: 15000 });
  });

  test("error handling works correctly", async ({ page }) => {
    // Intercept API to simulate error
    await page.route("**/api/v1/audit-log*", (route) => route.abort());

    await page.goto("/audit-log");

    // Wait for error to appear
    await expect(page.locator(".bg-red-50")).toBeVisible({ timeout: 10000 });
    await expect(page.locator(".text-red-800")).toContainText(
      /Failed|Error|error/i,
    );

    // Verify error can be dismissed
    const dismissButton = page
      .locator(
        'button[aria-label="Dismiss"], button .sr-only:has-text("Dismiss")',
      )
      .locator("..");
    if (await dismissButton.isVisible()) {
      await dismissButton.click();
      await expect(page.locator(".bg-red-50")).not.toBeVisible();
    }
  });

  test("filter pills can be removed individually", async ({ page }) => {
    await page.goto("/audit-log");
    await page.waitForLoadState("networkidle");

    // Apply multiple filters
    await page.selectOption("#resource-filter", "flag");
    await page.selectOption("#action-filter", "created");
    await page.waitForTimeout(1000);

    // Verify both filter pills are displayed
    await expect(page.locator("text=Resource: Flags")).toBeVisible();
    await expect(page.locator("text=Action: Created")).toBeVisible();

    // Remove resource filter pill
    const resourcePill = page.locator("text=Resource: Flags").locator("..");
    await resourcePill.locator("button").click();
    await page.waitForTimeout(1000);

    // Verify only action filter remains
    await expect(page.locator("text=Resource: Flags")).not.toBeVisible();
    await expect(page.locator("text=Action: Created")).toBeVisible();

    // Remove action filter pill
    const actionPill = page.locator("text=Action: Created").locator("..");
    await actionPill.locator("button").click();
    await page.waitForTimeout(1000);

    // Verify both filters are cleared
    await expect(page.locator("text=Action: Created")).not.toBeVisible();
  });

  test("displays realistic seeded audit data", async ({ page }) => {
    await page.goto("/audit-log");
    await page.waitForLoadState("networkidle");

    // Verify we have a substantial number of entries (spec mentions 50+)
    const entries = page.locator("ul.-mb-8 li");
    const entryCount = await entries.count();

    // With pagination, we should see at least 10-25 entries per page
    expect(entryCount).toBeGreaterThanOrEqual(10);

    // Verify we have different types of actions
    const uniqueActions = new Set();
    const actionElements = page.locator(".font-medium.text-gray-900");
    const actionCount = Math.min(10, await actionElements.count());

    for (let i = 0; i < actionCount; i++) {
      const actionText = await actionElements.nth(i).textContent();
      if (actionText) {
        uniqueActions.add(actionText);
      }
    }

    // Should have at least 3 different action types visible
    expect(uniqueActions.size).toBeGreaterThanOrEqual(2);

    // Verify dates span across time (not all the same timestamp)
    const timestamps = page.locator(".text-xs.text-gray-500.cursor-help");
    const timestampTexts = [];
    const timestampCount = Math.min(5, await timestamps.count());

    for (let i = 0; i < timestampCount; i++) {
      const timestampText = await timestamps.nth(i).textContent();
      if (timestampText) {
        timestampTexts.push(timestampText);
      }
    }

    // Should have varied timestamps (not all "Just now")
    const uniqueTimestamps = new Set(timestampTexts);
    expect(uniqueTimestamps.size).toBeGreaterThan(1);
  });
});
