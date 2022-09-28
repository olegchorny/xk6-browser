import { check } from 'k6';
import { chromium } from 'k6/x/browser';

export const options = {
  thresholds: {
    checks: ["rate==1.0"]
  }
}

export default function() {
  const browser = chromium.launch({
    headless: __ENV.XK6_HEADLESS ? true : false,
  });
  const context = browser.newContext();
  const page = context.newPage();

  // Goto front page, find login link and click it
  page.goto('https://test.k6.io/', { waitUntil: 'networkidle' });
  Promise.all([
    page.waitForNavigation(),
    page.locator('a[href="/my_messages.php"]').click(),
  ]).then(() => {
    // Enter login credentials and login
    page.locator('input[name="login"]').type('admin');
    page.locator('input[name="password"]').type('123');
    // We expect the form submission to trigger a navigation, so to prevent a
    // race condition, setup a waiter concurrently while waiting for the click
    // to resolve.
    return Promise.all([
      page.waitForNavigation(),
      page.locator('input[type="submit"]').click(),
    ]);
  }).then(() => {
    check(page, {
      'header': page.locator('h2').textContent() == 'Welcome, admin!',
    });
  }).finally(() => {
    page.close();
    browser.close();
  });
}
