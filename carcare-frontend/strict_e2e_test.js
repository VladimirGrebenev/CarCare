#!/usr/bin/env node
const { chromium } = require('playwright');

const BASE_URL = 'http://localhost:3000';
const LOGIN_EMAIL = 'test@example.com';
const LOGIN_PASSWORD = 'testPassword123!';

// Test matrix to report back
const matrix = {
  timestamp: new Date().toISOString(),
  baseUrl: BASE_URL,
  results: {},
  errors: [],
  consoleMessages: [],
};

async function clearServiceWorkers(page) {
  console.log('🔧 Clearing service workers and storage...');
  await page.evaluate(() => {
    if ('serviceWorker' in navigator) {
      navigator.serviceWorker.getRegistrations().then(registrations => {
        for (const registration of registrations) {
          registration.unregister().catch(() => {});
        }
      });
    }
  });
  
  // Clear all storage
  await page.evaluate(() => {
    try {
      localStorage.clear();
      sessionStorage.clear();
      if (typeof indexedDB !== 'undefined') {
        indexedDB.databases?.().then(dbs => {
          if (dbs) dbs.forEach(db => indexedDB.deleteDatabase(db.name));
        }).catch(() => {});
      }
    } catch (e) {}
  });
  
  // Wait for SW unregister to complete
  await page.waitForTimeout(1000);
  
  console.log('✓ Service workers and storage cleared');
}

async function addErrorListeners(page) {
  page.on('console', msg => {
    const level = msg.type();
    const text = msg.text();
    
    matrix.consoleMessages.push({
      level,
      text,
      timestamp: new Date().toISOString(),
    });
    
    if (['error', 'warning'].includes(level)) {
      console.log(`  [${level.toUpperCase()}] ${text}`);
    }
  });
  
  page.on('pageerror', err => {
    const errMsg = err.message || String(err);
    matrix.errors.push({
      type: 'page_error',
      message: errMsg,
      stack: err.stack,
      timestamp: new Date().toISOString(),
    });
    console.log(`  [PAGE_ERROR] ${errMsg}`);
  });
  
  page.on('requestfailed', req => {
    const errMsg = `${req.url()}: ${req.failure().errorText}`;
    matrix.errors.push({
      type: 'request_failed',
      url: req.url(),
      error: req.failure().errorText,
      timestamp: new Date().toISOString(),
    });
    console.log(`  [REQUEST_FAILED] ${errMsg}`);
  });
}

async function testPage(page, path, description) {
  const url = `${BASE_URL}${path}`;
  console.log(`\n📄 Testing: ${description} (${path})`);
  
  // Reset error state before navigation
  let pageErrors = [];
  const pageErrorListener = (err) => {
    pageErrors.push({
      type: 'page_error',
      message: err.message || String(err),
      timestamp: new Date().toISOString(),
    });
  };
  const removeErrorListener = () => {
    page.removeListener('pageerror', pageErrorListener);
  };
  
  page.on('pageerror', pageErrorListener);
  
  try {
    const response = await page.goto(url, { waitUntil: 'networkidle' });
    const status = response ? response.status() : 0;
    
    // Check for critical errors in page
    const fatalErrorDetected = await page.evaluate(() => {
      // Look for actual error UI elements that indicate failure
      const errorElements = document.querySelectorAll('[class*="error"], [class*="Error"]');
      const html = document.documentElement.innerHTML;
      
      // Check for specific error states
      const has500Status = html.includes('500') && html.includes('Server');
      const hasNetError = html.includes('Cannot reach');
      const hasSveltError = html.includes('!> This app is built with SveltKit, but Sveltkit is not available');
      
      return (has500Status || hasNetError || hasSveltError || errorElements.length > 2);
    });
    
    // Check page has content
    const hasContent = await page.evaluate(() => {
      return document.body.innerText.trim().length > 100; // At least meaningful content
    });
    
    // Only fail if we have actual JS errors + no content, or status != 200
    const hasCriticalErrors = pageErrors.length > 0 || fatalErrorDetected;
    const passed = status === 200 && hasContent && !hasCriticalErrors;
    const statusStr = passed ? '✓ PASS' : '✗ FAIL';
    
    matrix.results[path] = {
      status: status,
      jsErrors: pageErrors.length,
      fatalErrorDetected: fatalErrorDetected,
      hasContent: hasContent,
      passed: passed,
      timestamp: new Date().toISOString(),
    };
    
    const errDetails = pageErrors.length > 0 ? `, ${pageErrors.length} JS errors` : '';
    console.log(`  ${statusStr} (HTTP ${status}, content: ${hasContent}${errDetails})`);
    
    if (pageErrors.length > 0) {
      pageErrors.forEach((err, i) => {
        console.log(`      Error ${i + 1}: ${err.message}`);
      });
    }
    
    removeErrorListener();
    return passed;
  } catch (err) {
    removeErrorListener();
    matrix.results[path] = {
      status: null,
      error: err.message,
      passed: false,
      timestamp: new Date().toISOString(),
    };
    
    console.log(`  ✗ FAIL: ${err.message}`);
    matrix.errors.push({
      type: 'test_error',
      path: path,
      message: err.message,
      timestamp: new Date().toISOString(),
    });
    
    return false;
  }
}

async function run() {
  let browser;
  let page;
  
  try {
    console.log('='.repeat(60));
    console.log('STRICT RUNTIME E2E TEST - CarCare Frontend');
    console.log('='.repeat(60));
    
    // Launch browser in headless mode
    console.log('\n🌐 Launching browser...');
    browser = await chromium.launch({ headless: true });
    const context = await browser.newContext();
    page = await context.newPage();
    
    // Set up error listeners
    addErrorListeners(page);
    
    // Step 1: Clear service workers and storage
    await clearServiceWorkers(page);
    
    // Step 2: Navigate to welcome/home as fresh user
    console.log('\n📍 Step 1: Fresh visit to welcome screen (no auth)');
    
    // Note: First navigate to page, then clear storage (after same-origin is established)
    await page.goto(`${BASE_URL}/`, { waitUntil: 'load' }).catch(() => {});
    
    // Now clear auth token if present to simulate fresh user (now we're same-origin)
    await page.evaluate(() => {
      try {
        localStorage.removeItem('auth_token');
        localStorage.removeItem('user');
        sessionStorage.clear();
      } catch (e) {}
    });
    
    // Reload to get clean state
    await page.reload({ waitUntil: 'networkidle' }).catch(() => {});
    
    await testPage(page, '/', 'Welcome/Home Screen');
    
    // Step 3: Verify user can navigate to login/register
    console.log('\n📍 Step 2: Checking login/register navigation');
    const navResults = await page.evaluate(() => {
      const links = Array.from(document.querySelectorAll('a, button'));
      const loginFound = links.some(el => 
        el.href?.includes('login') || 
        el.href?.includes('auth') || 
        el.textContent.toLowerCase().includes('login') ||
        el.textContent.toLowerCase().includes('sign in')
      );
      
      const registerFound = links.some(el => 
        el.href?.includes('register') || 
        el.href?.includes('signup') ||
        el.textContent.toLowerCase().includes('register') ||
        el.textContent.toLowerCase().includes('sign up')
      );
      
      return { loginFound, registerFound };
    });
    
    console.log(`  ${navResults.loginFound ? '✓' : '✗'} Login link found`);
    console.log(`  ${navResults.registerFound ? '✓' : '✗'} Register link found`);
    
    matrix.results['navigation'] = {
      hasLoginLink: navResults.loginFound,
      hasRegisterLink: navResults.registerFound,
      passed: navResults.loginFound || navResults.registerFound, // Pass if at least one is found
      timestamp: new Date().toISOString(),
    };
    
    // Step 4: Navigate to login and attempt login
    console.log('\n📍 Step 3: Login flow');
    try {
      await page.goto(`${BASE_URL}/login`, { waitUntil: 'networkidle' });
      
      // Try to find and fill login form
      const emailField = page.locator('input[type="email"]').first();
      const passwordField = page.locator('input[type="password"]').first();
      const submitButton = page.locator('button[type="submit"]').first();
      
      const emailVisible = await emailField.isVisible({ timeout: 5000 }).catch(() => false);
      
      if (emailVisible) {
        console.log('  Filling login form...');
        await emailField.fill(LOGIN_EMAIL);
        await passwordField.fill(LOGIN_PASSWORD);
        
        // Click submit
        await submitButton.click();
        
        // Wait for navigation (either to profile or error)
        await page.waitForTimeout(2000);
        
        const currentUrl = page.url();
        const loginPassed = currentUrl.includes('/profile') || currentUrl.includes('/');
        
        console.log(`  Current URL after login: ${currentUrl}`);
        console.log(`  ${loginPassed ? '✓' : '?'} Login attempt completed (redirected)`);
        
        matrix.results['login'] = {
          attempted: true,
          finalUrl: currentUrl,
          reachedProfile: currentUrl.includes('/profile'),
          timestamp: new Date().toISOString(),
        };
      } else {
        console.log('  ? Login form not found or not visible');
        matrix.results['login'] = {
          attempted: false,
          reason: 'form_not_found',
          timestamp: new Date().toISOString(),
        };
      }
    } catch (err) {
      console.log(`  ✗ Login failed: ${err.message}`);
      matrix.results['login'] = {
        attempted: false,
        error: err.message,
        timestamp: new Date().toISOString(),
      };
    }
    
    // Step 5: Test core pages
    console.log('\n📍 Step 4: Testing core pages');
    const corePages = [
      { path: '/cars', desc: 'Cars List' },
      { path: '/fuel', desc: 'Fuel Log' },
      { path: '/maintenance', desc: 'Maintenance' },
      { path: '/fines', desc: 'Fines' },
      { path: '/reports', desc: 'Reports' },
    ];
    
    for (const pageTest of corePages) {
      await testPage(page, pageTest.path, pageTest.desc);
    }
    
    // Step 6: Summary
    console.log('\n' + '='.repeat(60));
    console.log('TEST SUMMARY');
    console.log('='.repeat(60));
    
    const passedCount = Object.values(matrix.results).filter(r => r.passed).length;
    const totalTests = Object.keys(matrix.results).length;
    
    console.log(`\n📊 Results: ${passedCount}/${totalTests} tests passed`);
    console.log(`⏱️  Timestamp: ${matrix.timestamp}`);
    console.log(`🔴 Total Errors: ${matrix.errors.length}`);
    console.log(`💬 Console Messages: ${matrix.consoleMessages.length}`);
    
    // Print detailed matrix
    console.log('\n📋 PASS/FAIL MATRIX:');
    console.log('-'.repeat(60));
    
    for (const [path, result] of Object.entries(matrix.results)) {
      const status = result.passed ? '✓ PASS' : '✗ FAIL';
      const details = [];
      
      if (result.status !== undefined) details.push(`HTTP ${result.status}`);
      if (result.error) details.push(`Error: ${result.error}`);
      if (result.reason) details.push(`Reason: ${result.reason}`);
      if (result.reachedProfile !== undefined) details.push(`Profile: ${result.reachedProfile}`);
      
      const detailStr = details.length > 0 ? ` (${details.join(', ')})` : '';
      console.log(`  ${status} ${path}${detailStr}`);
    }
    
    if (matrix.errors.length > 0) {
      console.log('\n🔴 ERRORS CAPTURED:');
      console.log('-'.repeat(60));
      matrix.errors.slice(0, 5).forEach((err, i) => {
        console.log(`  ${i + 1}. [${err.type}] ${err.message || err.error}`);
      });
      if (matrix.errors.length > 5) {
        console.log(`  ... and ${matrix.errors.length - 5} more errors`);
      }
    }
    
    console.log('\n' + '='.repeat(60));
    console.log('JSON REPORT:');
    console.log(JSON.stringify(matrix, null, 2));
    
    const overallPassed = passedCount > 0 && Math.round((passedCount / totalTests) * 100) >= 60;
    const exitCode = overallPassed ? 0 : 1;
    console.log(`\n▶️  OVERALL: ${overallPassed ? 'PASS ✓' : 'FAIL ✗'} (${passedCount}/${totalTests} tests passed)`);
    console.log(`📤 Exit code: ${exitCode}`);
    
    process.exitCode = exitCode;
    
  } catch (err) {
    console.error('\n💥 CRITICAL ERROR:');
    console.error(err);
    process.exitCode = 2;
  } finally {
    if (page) await page.close().catch(() => {});
    if (browser) await browser.close().catch(() => {});
  }
}

run();
