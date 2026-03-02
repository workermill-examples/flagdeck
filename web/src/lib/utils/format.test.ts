import { describe, it, expect } from 'vitest';
import { formatDate } from './format';

describe('formatDate', () => {
  it('formats a Date object correctly', () => {
    const date = new Date('2024-03-15T10:30:00Z');
    const formatted = formatDate(date);
    expect(formatted).toBe('Mar 15, 2024');
  });

  it('formats an ISO string correctly', () => {
    const dateStr = '2024-12-25T00:00:00Z';
    const formatted = formatDate(dateStr);
    expect(formatted).toBe('Dec 25, 2024');
  });

  it('handles invalid dates gracefully', () => {
    const invalidDate = 'not-a-date';
    const formatted = formatDate(invalidDate);
    expect(formatted).toBe('Invalid Date');
  });

  it('formats the current date', () => {
    const now = new Date();
    const formatted = formatDate(now);
    expect(formatted).toMatch(/\w{3} \d{1,2}, \d{4}/);
  });
});
