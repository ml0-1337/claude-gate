const assert = require('assert');
const path = require('path');
const fs = require('fs');
const { getPlatform } = require('./install.js');

// Test suite for install.js
console.log('Running install.js unit tests...\n');

let passed = 0;
let failed = 0;

function test(name, fn) {
  try {
    fn();
    console.log(`✓ ${name}`);
    passed++;
  } catch (error) {
    console.error(`✗ ${name}`);
    console.error(`  ${error.message}`);
    failed++;
  }
}

// Test getPlatform function
test('getPlatform returns correct platform info for current system', () => {
  const platform = getPlatform();
  assert(platform.type, 'Platform type should be defined');
  assert(platform.arch, 'Platform arch should be defined');
  assert(platform.packageName, 'Package name should be defined');
  assert(typeof platform.isWindows === 'boolean', 'isWindows should be boolean');
});

test('getPlatform maps darwin-x64 correctly', () => {
  const originalPlatform = process.platform;
  const originalArch = process.arch;
  
  Object.defineProperty(process, 'platform', { value: 'darwin', configurable: true });
  Object.defineProperty(process, 'arch', { value: 'x64', configurable: true });
  
  const platform = getPlatform();
  assert.strictEqual(platform.packageName, '@claude-gate/darwin-x64');
  assert.strictEqual(platform.isWindows, false);
  
  Object.defineProperty(process, 'platform', { value: originalPlatform, configurable: true });
  Object.defineProperty(process, 'arch', { value: originalArch, configurable: true });
});

test('getPlatform maps darwin-arm64 correctly', () => {
  const originalPlatform = process.platform;
  const originalArch = process.arch;
  
  Object.defineProperty(process, 'platform', { value: 'darwin', configurable: true });
  Object.defineProperty(process, 'arch', { value: 'arm64', configurable: true });
  
  const platform = getPlatform();
  assert.strictEqual(platform.packageName, '@claude-gate/darwin-arm64');
  assert.strictEqual(platform.isWindows, false);
  
  Object.defineProperty(process, 'platform', { value: originalPlatform, configurable: true });
  Object.defineProperty(process, 'arch', { value: originalArch, configurable: true });
});

test('getPlatform maps win32-x64 correctly', () => {
  const originalPlatform = process.platform;
  const originalArch = process.arch;
  
  Object.defineProperty(process, 'platform', { value: 'win32', configurable: true });
  Object.defineProperty(process, 'arch', { value: 'x64', configurable: true });
  
  const platform = getPlatform();
  assert.strictEqual(platform.packageName, '@claude-gate/win32-x64');
  assert.strictEqual(platform.isWindows, true);
  
  Object.defineProperty(process, 'platform', { value: originalPlatform, configurable: true });
  Object.defineProperty(process, 'arch', { value: originalArch, configurable: true });
});

test('getPlatform maps win32-ia32 to x64 package', () => {
  const originalPlatform = process.platform;
  const originalArch = process.arch;
  
  Object.defineProperty(process, 'platform', { value: 'win32', configurable: true });
  Object.defineProperty(process, 'arch', { value: 'ia32', configurable: true });
  
  const platform = getPlatform();
  assert.strictEqual(platform.packageName, '@claude-gate/win32-x64');
  assert.strictEqual(platform.isWindows, true);
  
  Object.defineProperty(process, 'platform', { value: originalPlatform, configurable: true });
  Object.defineProperty(process, 'arch', { value: originalArch, configurable: true });
});

test('getPlatform throws error for unsupported platform', () => {
  const originalPlatform = process.platform;
  const originalArch = process.arch;
  
  Object.defineProperty(process, 'platform', { value: 'freebsd', configurable: true });
  Object.defineProperty(process, 'arch', { value: 'x64', configurable: true });
  
  assert.throws(() => {
    getPlatform();
  }, /Unsupported platform/);
  
  Object.defineProperty(process, 'platform', { value: originalPlatform, configurable: true });
  Object.defineProperty(process, 'arch', { value: originalArch, configurable: true });
});

test('getPlatform error message includes supported platforms', () => {
  const originalPlatform = process.platform;
  const originalArch = process.arch;
  
  Object.defineProperty(process, 'platform', { value: 'aix', configurable: true });
  Object.defineProperty(process, 'arch', { value: 'ppc64', configurable: true });
  
  try {
    getPlatform();
    assert.fail('Should have thrown error');
  } catch (error) {
    assert(error.message.includes('darwin-x64'), 'Error should list darwin-x64');
    assert(error.message.includes('linux-arm64'), 'Error should list linux-arm64');
    assert(error.message.includes('win32-x64'), 'Error should list win32-x64');
  }
  
  Object.defineProperty(process, 'platform', { value: originalPlatform, configurable: true });
  Object.defineProperty(process, 'arch', { value: originalArch, configurable: true });
});

// Test file paths
test('install script uses correct relative paths', () => {
  const scriptPath = path.join(__dirname, 'install.js');
  const content = fs.readFileSync(scriptPath, 'utf8');
  
  // Check for correct relative paths
  assert(content.includes("path.join(__dirname, '..', 'bin')"), 'Should create bin in parent directory');
  assert(content.includes("path.join(__dirname, '..', 'node_modules'"), 'Should look for node_modules in parent');
});

// Summary
console.log('\n----------------------------------------');
console.log(`Total tests: ${passed + failed}`);
console.log(`Passed: ${passed}`);
console.log(`Failed: ${failed}`);
console.log('----------------------------------------\n');

if (failed > 0) {
  process.exit(1);
}