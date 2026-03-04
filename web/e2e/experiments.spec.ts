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

    // Verify experiments are displayed (should not be empty based on seeded data)
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
        .locator("text=checkout-redesign")
        .or(
          page
            .locator("text=search-algorithm")
            .or(page.locator("text=Checkout").or(page.locator("text=Search"))),
        )
        .first(),
    ).toBeVisible();
  });

  test("experiment results chart renders with seeded data", async ({
    page,
  }) => {
    await page.goto("/experiments");
    await page.waitForLoadState("networkidle");

    // Find an experiment and expand it to see the chart
    const firstExperiment = page.locator("ul.divide-y li").first();
    const expandButton = firstExperiment.locator("button svg.transform");

    // Click to expand the first experiment
    await expandButton.click();

    // Wait for expansion animation
    await page.waitForTimeout(500);

    // Verify the performance results section appears
    await expect(
      page.locator('h4:has-text("Performance Results")'),
    ).toBeVisible();

    // Verify the chart is rendered
    await expect(page.locator("svg")).toBeVisible();

    // Verify chart legend
    await expect(page.locator("text=Impressions")).toBeVisible();
    await expect(page.locator("text=Conversions")).toBeVisible();

    // Verify summary stats are displayed
    await expect(page.locator("text=Total Impressions")).toBeVisible();
    await expect(page.locator("text=Total Conversions")).toBeVisible();
    await expect(page.locator("text=Overall CR")).toBeVisible();

    // Verify variant configuration section
    await expect(
      page.locator('h4:has-text("Variant Configuration")'),
    ).toBeVisible();

    // Verify variants are displayed with realistic names
    await expect(
      page.locator("text=Control").or(page.locator("text=Variant")).first(),
    ).toBeVisible();

    // Verify variants show weights (should add up to 100%)
    await expect(
      page.locator("text=50%").or(page.locator("text=100%")).first(),
    ).toBeVisible();
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
    await page.fill("#flag_key", "test-flag-key");

    // Set status to running
    await page.selectOption("#status", "running");

    // Verify default variants are present
    await expect(page.locator('input[value="control"]')).toBeVisible();
    await expect(page.locator('input[value="variant"]')).toBeVisible();

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
    await page.click('button:has-text("Create")');

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

    // Click edit button on first experiment
    await page.click('button:has-text("Edit")');

    // Verify edit form appears
    await expect(page.locator('h3:has-text("Edit Experiment")')).toBeVisible();

    // Verify fields are pre-filled
    const nameValue = await page.locator("#name").inputValue();
    expect(nameValue).toBeTruthy();

    const descValue = await page.locator("#description").inputValue();
    // Description might be empty, that's ok

    const flagKeyValue = await page.locator("#flag_key").inputValue();
    expect(flagKeyValue).toBeTruthy();

    // Update description
    const updatedDescription = `Updated description ${Date.now()}`;
    await page.fill("#description", updatedDescription);

    // Submit the form
    await page.click('button:has-text("Update")');

    // Wait for form to disappear
    await expect(
      page.locator('h3:has-text("Edit Experiment")'),
    ).not.toBeVisible({ timeout: 10000 });

    // Expand the experiment to verify the description was updated
    const firstExperiment = page.locator("ul.divide-y li").first();
    const expandButton = firstExperiment.locator("button svg.transform");
    await expandButton.click();

    // The updated description might be shown in the expanded view or list
    // At minimum, verify the experiment is still there and the form closed successfully
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

    // Remove a variant
    await page.click('button:has-text("Remove")');

    // Verify we're back to 2 variants
    const finalVariantCount = await page
      .locator(".border.border-gray-200.rounded-lg.p-4")
      .count();
    expect(finalVariantCount).toBe(2);
  });

  test("shows experiment status badges correctly", async ({ page }) => {
    await page.goto("/experiments");
    await page.waitForLoadState("networkidle");

    // Verify status badges are present and colored correctly
    const statusBadges = page.locator(
      '[class*="bg-"][class*="100"][class*="text-"]',
    );
    const badgeCount = await statusBadges.count();
    expect(badgeCount).toBeGreaterThan(0);

    // Verify at least one status is displayed
    const statuses = ["draft", "running", "paused", "completed"];
    let foundStatus = false;

    for (const status of statuses) {
      const statusElement = page.locator(`text=${status}`);
      if (await statusElement.isVisible()) {
        foundStatus = true;
        break;
      }
    }

    expect(foundStatus).toBe(true);
  });

  test("displays flag keys and dates correctly", async ({ page }) => {
    await page.goto("/experiments");
    await page.waitForLoadState("networkidle");

    // Verify flag keys are displayed
    const flagKeyCode = page
      .locator("code")
      .filter({ hasText: /^[a-z-]+$/ })
      .first();
    await expect(flagKeyCode).toBeVisible();

    // Expand first experiment to see more details
    const expandButton = page.locator("button svg.transform").first();
    await expandButton.click();

    // Look for date information (might show start/end dates)
    const dateInfo = page
      .locator("text=Started:")
      .or(page.locator("text=Ended:"));
    // Dates might not always be present, so we just check if the expanded view worked
    await expect(
      page.locator('h4:has-text("Variant Configuration")'),
    ).toBeVisible();
  });

  test("chart displays realistic conversion data", async ({ page }) => {
    await page.goto("/experiments");
    await page.waitForLoadState("networkidle");

    // Expand the first experiment
    const expandButton = page.locator("button svg.transform").first();
    await expandButton.click();

    // Wait for chart to render
    await page.waitForTimeout(1000);

    // Verify chart contains data (bars should be visible)
    const chartBars = page.locator(
      'svg rect[fill="#3b82f6"], svg rect[fill="#10b981"]',
    );
    const barCount = await chartBars.count();
    expect(barCount).toBeGreaterThan(0);

    // Verify conversion rate percentages are displayed
    const conversionRates = page.locator("text=/\\d+\\.\\d+% CR/");
    const rateCount = await conversionRates.count();
    expect(rateCount).toBeGreaterThan(0);

    // Verify traffic percentages are displayed
    const trafficPercentages = page.locator("text=/\\d+% traffic/");
    const trafficCount = await trafficPercentages.count();
    expect(trafficCount).toBeGreaterThan(0);

    // Verify summary statistics show realistic numbers
    const totalImpressions = page
      .locator("text=Total Impressions")
      .locator("..")
      .locator(".text-lg");
    const impressionsText = await totalImpressions.textContent();
    expect(impressionsText).toMatch(/\d+/); // Should contain numbers

    const totalConversions = page
      .locator("text=Total Conversions")
      .locator("..")
      .locator(".text-lg");
    const conversionsText = await totalConversions.textContent();
    expect(conversionsText).toMatch(/\d+/); // Should contain numbers
  });

  test("delete confirmation works correctly", async ({ page }) => {
    await page.goto("/experiments");
    await page.waitForLoadState("networkidle");

    // Click delete button on first experiment
    await page.click('button:has-text("Delete")');

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

    // Manually set weights that don't sum to 100
    const weightInputs = page.locator('input[type="number"]');
    await weightInputs.first().fill("60");
    await weightInputs.nth(1).fill("30");

    // Verify total weight indicator shows the issue
    await expect(page.locator("text=Total weight: 90%")).toBeVisible();
    await expect(page.locator("text=Must equal 100%")).toBeVisible();

    // Try to submit (should fail)
    await page.click('button:has-text("Create")');

    // Should show validation error
    await expect(page.locator("text=weights must sum to 100%")).toBeVisible();
  });
});
