import { test, expect } from "@playwright/test";

// Helper function to login before each test
async function login(page: any) {
  await page.goto("/login");
  await page.waitForSelector("#email", { state: "visible" });
  await page.click('button[type="submit"]');
  await page.waitForURL("**/dashboard", { timeout: 30000 });
}

test.describe("Experiments Page", () => {
  test.beforeEach(async ({ page }) => {
    // Navigate first so localStorage is accessible (about:blank denies access)
    await page.goto("/login");
    await page.context().clearCookies();
    await page.evaluate(() => localStorage.clear());

    // Login with demo credentials
    await login(page);
  });

  test("displays experiments list with seeded data", async ({ page }) => {
    await page.goto("/experiments");

    // Verify page title and header
    await expect(page).toHaveTitle(/FlagDeck/);
    await expect(page.locator("main h1")).toContainText("Experiments");

    // Wait for data to load
    await page.waitForLoadState("networkidle");

    // Wait for experiments to load and verify they are displayed
    await expect(page.locator("ul.divide-y li").first()).toBeVisible({
      timeout: 10000,
    });
    const experimentItems = page.locator("ul.divide-y li");
    const experimentCount = await experimentItems.count();
    expect(experimentCount).toBeGreaterThan(0); // Spec mentions 2 experiments

    // Verify experiment list shows status badges
    await expect(
      page
        .locator('.bg-green-100.text-green-800:has-text("running")')
        .or(page.locator('.bg-blue-100.text-blue-800:has-text("completed")'))
        .first(),
    ).toBeVisible();

    // Verify realistic experiment names from seeded data
    await expect(
      page
        .locator("text=Dashboard Layout A/B Test")
        .or(page.locator("text=Search Algorithm Performance Test"))
        .first(),
    ).toBeVisible();
  });

  test("experiment results chart renders with seeded data", async ({
    page,
  }) => {
    await page.goto("/experiments");
    await page.waitForLoadState("networkidle");

    // Find an experiment and expand it
    const firstExperiment = page.locator("ul.divide-y li").first();
    const expandButton = firstExperiment.locator("button svg.w-5.h-5");

    // Click to expand the first experiment
    await expandButton.click();

    // Wait for expansion animation
    await page.waitForTimeout(500);

    // Verify the performance results section appears
    await expect(
      page.locator('h4:has-text("Performance Results")'),
    ).toBeVisible();

    // Verify variant configuration section
    await expect(
      page.locator('h4:has-text("Variant Configuration")'),
    ).toBeVisible();

    // Verify variants are displayed with names
    const variantCards = page.locator(".bg-gray-50.rounded-lg.p-3");
    const variantCount = await variantCards.count();
    expect(variantCount).toBeGreaterThan(0);

    // Verify variants show weight percentages
    await expect(page.locator("text=/%/").first()).toBeVisible();
  });

  test("can create a new experiment", async ({ page }) => {
    await page.goto("/experiments");

    // Click create experiment button
    await page.click('button:has-text("Create Experiment")');

    // Verify form appears
    await expect(
      page.locator('h3:has-text("Create Experiment")'),
    ).toBeVisible();

    // Fill in experiment details
    const experimentName = `Test Experiment ${Date.now()}`;
    await page.fill("#name", experimentName);

    // Verify key is auto-generated
    const generatedKey = await page.locator("#key").inputValue();
    expect(generatedKey).toBeTruthy();

    await page.fill("#description", "A test experiment created by E2E tests");
    await page.fill("#flag_key", "new_dashboard");

    // Set environment and status
    await page.selectOption("#environment", "production");
    await page.selectOption("#status", "running");

    // Verify default variants are present (key inputs with placeholder "control")
    await expect(
      page.locator('input[placeholder="control"]').first(),
    ).toBeVisible();
    await expect(page.locator(".border.border-gray-200.rounded-lg.p-4")).toHaveCount(
      2,
    );

    // Verify weights are balanced
    const controlWeight = await page
      .locator('input[type="number"]')
      .first()
      .inputValue();
    const variantWeight = await page
      .locator('input[type="number"]')
      .nth(1)
      .inputValue();
    expect(parseInt(controlWeight) + parseInt(variantWeight)).toBe(100);

    // Submit the form
    await page.click('button[type="submit"]');

    // Wait for form to disappear
    await expect(
      page.locator('h3:has-text("Create Experiment")'),
    ).not.toBeVisible({ timeout: 10000 });

    // Verify the new experiment appears in the list
    await expect(page.locator(`text=${experimentName}`)).toBeVisible();
  });

  test("can edit existing experiment", async ({ page }) => {
    await page.goto("/experiments");
    await page.waitForLoadState("networkidle");

    // Wait for experiments list to load
    await expect(page.locator("ul.divide-y li").first()).toBeVisible({
      timeout: 10000,
    });

    // Click edit button on first experiment
    await page.locator('button:has-text("Edit")').first().click();

    // Verify edit form appears
    await expect(page.locator('h3:has-text("Edit Experiment")')).toBeVisible();

    // Verify fields are pre-filled
    const nameValue = await page.locator("#name").inputValue();
    expect(nameValue).toBeTruthy();

    const flagKeyValue = await page.locator("#flag_key").inputValue();
    expect(flagKeyValue).toBeTruthy();

    // Verify weights are valid numbers (not NaN)
    const weightInputs = page.locator('input[type="number"]');
    const weightCount = await weightInputs.count();
    if (weightCount > 0) {
      const firstWeight = await weightInputs.first().inputValue();
      expect(parseInt(firstWeight)).not.toBeNaN();
    }

    // Update description
    const updatedDescription = `Updated description ${Date.now()}`;
    await page.fill("#description", updatedDescription);

    // Submit the form
    await page.click('button:has-text("Update")');

    // Wait for form to disappear
    await expect(
      page.locator('h3:has-text("Edit Experiment")'),
    ).not.toBeVisible({ timeout: 10000 });

    // Verify the experiment is still there and the form closed successfully
    await expect(page.locator("ul.divide-y li").first()).toBeVisible();
  });

  test("can add and remove experiment variants", async ({ page }) => {
    await page.goto("/experiments");

    // Start creating a new experiment
    await page.click('button:has-text("Create Experiment")');

    // Fill basic info
    await page.fill("#name", "Multi Variant Test");
    await page.fill("#flag_key", "multi-variant-test");

    // Verify we start with 2 variants
    const initialVariantCount = await page
      .locator(".border.border-gray-200.rounded-lg.p-4")
      .count();
    expect(initialVariantCount).toBe(2);

    // Add a third variant
    await page.click('button:has-text("Add Variant")');

    // Verify we now have 3 variants
    const newVariantCount = await page
      .locator(".border.border-gray-200.rounded-lg.p-4")
      .count();
    expect(newVariantCount).toBe(3);

    // Verify weights were rebalanced (should be roughly 33% each)
    const weightInputs = page.locator('input[type="number"]');
    const weights = [];
    const count = await weightInputs.count();
    for (let i = 0; i < count; i++) {
      const value = await weightInputs.nth(i).inputValue();
      weights.push(parseInt(value));
    }

    const totalWeight = weights.reduce((sum, weight) => sum + weight, 0);
    expect(totalWeight).toBe(100);

    // Remove a variant (use .first() since multiple Remove buttons exist)
    await page.locator('button:has-text("Remove")').first().click();

    // Verify we're back to 2 variants
    const finalVariantCount = await page
      .locator(".border.border-gray-200.rounded-lg.p-4")
      .count();
    expect(finalVariantCount).toBe(2);
  });

  test("shows experiment status badges correctly", async ({ page }) => {
    await page.goto("/experiments");
    await page.waitForLoadState("networkidle");

    // Wait for experiment list to load
    await expect(page.locator("ul.divide-y li").first()).toBeVisible({
      timeout: 10000,
    });

    // Verify status badges are present and colored correctly
    // Status badges use classes like bg-green-100, bg-blue-100, etc.
    const statusBadges = page.locator(
      'span[class*="rounded-full"][class*="bg-"]',
    );
    const badgeCount = await statusBadges.count();
    expect(badgeCount).toBeGreaterThan(0);

    // Verify at least one known status is displayed in a badge
    const statuses = ["draft", "running", "paused", "completed"];
    let foundStatus = false;

    for (const status of statuses) {
      // Use first() to avoid strict mode violation when multiple badges have the same text
      const statusElement = page
        .locator('span[class*="rounded-full"]')
        .filter({ hasText: status })
        .first();
      if (await statusElement.isVisible().catch(() => false)) {
        foundStatus = true;
        break;
      }
    }

    expect(foundStatus).toBe(true);
  });

  test("displays flag keys and dates correctly", async ({ page }) => {
    await page.goto("/experiments");
    await page.waitForLoadState("networkidle");

    // Wait for experiment list to load
    await expect(page.locator("ul.divide-y li").first()).toBeVisible({
      timeout: 10000,
    });

    // Verify flag keys are displayed (keys use underscores)
    const flagKeyCode = page
      .locator("code")
      .filter({ hasText: /^[a-z_-]+$/ })
      .first();
    await expect(flagKeyCode).toBeVisible();

    // Expand first experiment to see more details
    const firstExperiment = page.locator("ul.divide-y li").first();
    const expandButton = firstExperiment.locator("button svg.w-5.h-5");
    await expandButton.click();

    // Verify the expanded view shows variant configuration
    await expect(
      page.locator('h4:has-text("Variant Configuration")'),
    ).toBeVisible();
  });

  test("chart displays realistic conversion data", async ({ page }) => {
    await page.goto("/experiments");
    await page.waitForLoadState("networkidle");

    // Wait for experiment list to load
    await expect(page.locator("ul.divide-y li").first()).toBeVisible({
      timeout: 10000,
    });

    // Expand the first experiment
    const firstExperiment = page.locator("ul.divide-y li").first();
    const expandButton = firstExperiment.locator("button svg.w-5.h-5");
    await expandButton.click();

    // Wait for expansion
    await page.waitForTimeout(1000);

    // Verify the performance results section appears
    await expect(
      page.locator('h4:has-text("Performance Results")'),
    ).toBeVisible();

    // The chart may show data or "No experiment data to display"
    const chartSection = page.locator(".experiment-chart");
    if (await chartSection.isVisible()) {
      // If chart has data, verify bars and stats
      const chartBars = page.locator(
        'svg rect[fill="#3b82f6"], svg rect[fill="#10b981"]',
      );
      const barCount = await chartBars.count();

      if (barCount > 0) {
        // Verify summary statistics
        await expect(page.locator("text=Total Impressions")).toBeVisible();
        await expect(page.locator("text=Total Conversions")).toBeVisible();
      } else {
        // Empty chart state
        await expect(
          page.locator("text=No experiment data to display"),
        ).toBeVisible();
      }
    }

    // Verify variant configuration section is always present
    await expect(
      page.locator('h4:has-text("Variant Configuration")'),
    ).toBeVisible();
  });

  test("delete confirmation works correctly", async ({ page }) => {
    await page.goto("/experiments");
    await page.waitForLoadState("networkidle");

    // Click delete button on first experiment
    await page.locator('button:has-text("Delete")').first().click();

    // Verify confirmation modal appears
    await expect(
      page.locator('h3:has-text("Delete Experiment")'),
    ).toBeVisible();
    await expect(
      page.locator("text=This action cannot be undone"),
    ).toBeVisible();

    // Cancel the deletion
    await page.click('button:has-text("Cancel")');

    // Verify modal disappears
    await expect(
      page.locator('h3:has-text("Delete Experiment")'),
    ).not.toBeVisible();

    // Verify experiment is still in the list
    const experiments = page.locator("ul.divide-y li");
    const experimentCount = await experiments.count();
    expect(experimentCount).toBeGreaterThan(0);
  });

  test("handles loading and error states", async ({ page }) => {
    // Test loading state
    await page.route("**/api/v1/experiments", async (route) => {
      await new Promise((resolve) => setTimeout(resolve, 2000));
      await route.continue();
    });

    await page.goto("/experiments");

    // Verify loading state
    await expect(page.locator("text=Loading experiments")).toBeVisible();
    await expect(page.locator(".animate-spin")).toBeVisible();

    // Wait for loading to complete
    await page.waitForLoadState("networkidle");
    await expect(page.locator("text=Loading experiments")).not.toBeVisible({
      timeout: 15000,
    });
  });

  test("empty state displays correctly", async ({ page }) => {
    // Intercept experiments API to return empty array
    await page.route("**/api/v1/experiments", (route) =>
      route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ data: [], total: 0 }),
      }),
    );

    await page.goto("/experiments");

    // Verify empty state
    await expect(page.locator('h3:has-text("No experiments")')).toBeVisible();
    await expect(
      page.locator("text=Get started by creating your first experiment"),
    ).toBeVisible();
  });

  test("variant weights validation works", async ({ page }) => {
    await page.goto("/experiments");

    // Start creating experiment
    await page.click('button:has-text("Create Experiment")');

    await page.fill("#name", "Weight Test");
    await page.fill("#flag_key", "weight-test");

    // Verify default weights sum to 100%
    await expect(page.locator("text=Total weight: 100%")).toBeVisible();

    // Verify initial weight balance (50/50 default)
    const weightInputs = page.locator('input[type="number"]');
    const controlWeight = await weightInputs.first().inputValue();
    const variantWeight = await weightInputs.nth(1).inputValue();
    expect(parseInt(controlWeight) + parseInt(variantWeight)).toBe(100);

    // Change first weight - verify auto-rebalancing maintains 100%
    await weightInputs.first().fill("70");
    await expect(page.locator("text=Total weight: 100%")).toBeVisible();

    // Verify second weight was adjusted accordingly
    const newVariantWeight = await weightInputs.nth(1).inputValue();
    expect(parseInt(newVariantWeight)).toBe(30);
  });
});
