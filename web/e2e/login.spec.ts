import { test, expect } from "@playwright/test";

test.describe("Login Flow", () => {
  test.beforeEach(async ({ page }) => {
    // Navigate first so localStorage is accessible (about:blank denies access)
    await page.goto("/login");
    await page.context().clearCookies();
    await page.evaluate(() => localStorage.clear());
  });

  test("user can login with demo credentials and redirect to dashboard", async ({
    page,
  }) => {
    await page.goto("/login");

    // Verify we're on the login page
    await expect(page).toHaveTitle(/Login - FlagDeck/);
    await expect(page.locator("h2")).toContainText("Sign in to FlagDeck");

    // Verify demo credentials are pre-filled
    await expect(page.locator("#email")).toHaveValue("demo@workermill.com");
    await expect(page.locator("#password")).toHaveValue("demo1234");

    // Submit the login form
    await page.click('button[type="submit"]');

    // Wait for redirect to dashboard
    await expect(page).toHaveURL("/dashboard", { timeout: 15000 });

    // Verify we're on the dashboard
    await expect(page.locator("main h1")).toContainText("Dashboard");

    // Verify authentication was successful by checking for user menu or sidebar
    await expect(
      page.locator('[data-testid="user-menu"]').or(page.locator("nav")).first(),
    ).toBeVisible();
  });

  test("shows error for invalid credentials", async ({ page }) => {
    await page.goto("/login");

    // Clear the pre-filled email and enter invalid credentials
    await page.fill("#email", "wrong@example.com");
    await page.fill("#password", "wrongpassword");

    // Submit the form
    await page.click('button[type="submit"]');

    // Wait for error to appear
    await expect(page.locator(".bg-red-50")).toBeVisible();
    await expect(page.locator(".text-red-800")).toContainText(
      /Login failed|Invalid credentials|Authentication failed/,
    );

    // Verify we're still on login page
    await expect(page).toHaveURL("/login");
  });

  test("login button shows loading state during submission", async ({
    page,
  }) => {
    await page.goto("/login");

    // Intercept login request to simulate slow response
    await page.route("**/auth/login", async (route) => {
      await new Promise((resolve) => setTimeout(resolve, 1000));
      await route.continue();
    });

    // Submit the form
    await page.click('button[type="submit"]');

    // Verify loading state
    await expect(page.locator('button[type="submit"]')).toContainText(
      "Signing in...",
    );
    await expect(page.locator(".animate-spin")).toBeVisible();

    // Verify button is disabled during submission
    await expect(page.locator('button[type="submit"]')).toBeDisabled();
  });

  test("redirects to dashboard if already authenticated", async ({ page }) => {
    // First login to get authentication
    await page.goto("/login");
    await page.click('button[type="submit"]');
    await expect(page).toHaveURL("/dashboard");

    // Now try to visit login page again
    await page.goto("/login");

    // Should be redirected to dashboard
    await expect(page).toHaveURL("/dashboard", { timeout: 15000 });
  });

  test("demo credentials notice is visible", async ({ page }) => {
    await page.goto("/login");

    // Verify demo credentials notice is displayed
    await expect(page.locator(".bg-blue-50")).toBeVisible();
    await expect(page.locator(".text-blue-800")).toContainText(
      "Demo credentials:",
    );
    await expect(page.locator(".text-blue-800")).toContainText(
      "demo@workermill.com",
    );
    await expect(page.locator(".text-blue-800")).toContainText("demo1234");
  });
});
