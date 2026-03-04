module.exports = {
  testDir: './tests',
  timeout: 30000,
  use: {
    baseURL: process.env.PUBLIC_API_URL || 'http://localhost:8080',
    headless: true,
  },
  projects: [
    {
      name: 'chromium',
      use: { browserName: 'chromium' },
    },
  ],
};