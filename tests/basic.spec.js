const { test, expect } = require('@playwright/test');

test('basic health check', async ({ page }) => {
  const response = await page.request.get('/health');
  expect(response.status()).toBe(200);
});