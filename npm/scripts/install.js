#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const { promisify } = require('util');
const chmod = promisify(fs.chmod);
const mkdir = promisify(fs.mkdir);
const copyFile = promisify(fs.copyFile);

// Platform detection
function getPlatform() {
  const type = process.platform;
  const arch = process.arch;
  
  // Map Node.js platform/arch to our package naming
  const platformMap = {
    'darwin-x64': '@claude-gate/darwin-x64',
    'darwin-arm64': '@claude-gate/darwin-arm64',
    'linux-x64': '@claude-gate/linux-x64',
    'linux-arm64': '@claude-gate/linux-arm64',
    'win32-x64': '@claude-gate/win32-x64',
    'win32-ia32': '@claude-gate/win32-x64', // Use x64 for 32-bit Windows too
  };
  
  const platformKey = `${type}-${arch}`;
  const packageName = platformMap[platformKey];
  
  if (!packageName) {
    throw new Error(
      `Unsupported platform: ${type}-${arch}\n` +
      `Supported platforms: ${Object.keys(platformMap).join(', ')}`
    );
  }
  
  return {
    type,
    arch,
    packageName,
    isWindows: type === 'win32'
  };
}

// Find the platform-specific binary
function findBinary(platform) {
  const possiblePaths = [
    // Installed as dependency
    path.join(__dirname, '..', 'node_modules', platform.packageName),
    // Installed globally
    path.join(__dirname, '..', '..', platform.packageName),
    // Development/testing
    path.join(__dirname, '..', '..', '..', 'npm-packages', platform.packageName.replace('@claude-gate/', ''))
  ];
  
  for (const basePath of possiblePaths) {
    const binaryName = platform.isWindows ? 'bin.exe' : 'bin';
    const binaryPath = path.join(basePath, binaryName);
    
    if (fs.existsSync(binaryPath)) {
      return binaryPath;
    }
  }
  
  return null;
}

// Create wrapper script for better error handling
function createWrapper(binDir, actualBinaryPath, platform) {
  const wrapperPath = path.join(binDir, 'claude-gate');
  
  if (platform.isWindows) {
    // Windows batch file
    const batchContent = `@echo off
"${actualBinaryPath}" %*
`;
    fs.writeFileSync(wrapperPath + '.cmd', batchContent);
    
    // PowerShell script for better compatibility
    const psContent = `& "${actualBinaryPath}" $args
`;
    fs.writeFileSync(wrapperPath + '.ps1', psContent);
  } else {
    // Unix shell script
    const shContent = `#!/bin/sh
exec "${actualBinaryPath}" "$@"
`;
    fs.writeFileSync(wrapperPath, shContent);
    fs.chmodSync(wrapperPath, 0o755);
  }
}

// Main installation logic
async function install() {
  console.log('Installing claude-gate...');
  
  try {
    // Detect platform
    const platform = getPlatform();
    console.log(`Platform detected: ${platform.type}-${platform.arch}`);
    
    // Find the binary
    const sourceBinary = findBinary(platform);
    if (!sourceBinary) {
      console.error(`
ERROR: Could not find platform binary for ${platform.packageName}

This might happen if:
1. The platform package failed to install
2. You're using an unsupported platform
3. Installation was run with --ignore-scripts

Try running:
  npm install ${platform.packageName}

Or download the binary manually from:
  https://github.com/ml0-1337/claude-gate/releases
`);
      process.exit(1);
    }
    
    console.log(`Found binary at: ${sourceBinary}`);
    
    // Create bin directory
    const binDir = path.join(__dirname, '..', 'bin');
    await mkdir(binDir, { recursive: true });
    
    // Copy binary to bin directory
    const targetBinary = path.join(binDir, platform.isWindows ? 'claude-gate.exe' : 'claude-gate-bin');
    await copyFile(sourceBinary, targetBinary);
    
    // Set executable permissions on Unix
    if (!platform.isWindows) {
      await chmod(targetBinary, 0o755);
    }
    
    // Create wrapper script
    createWrapper(binDir, targetBinary, platform);
    
    console.log('✅ claude-gate installed successfully!');
    console.log('Run "claude-gate --help" to get started.');
    
  } catch (error) {
    console.error('Installation failed:', error.message);
    console.error('\nFor manual installation instructions, visit:');
    console.error('https://github.com/ml0-1337/claude-gate#installation');
    process.exit(1);
  }
}

// Check if this is being run directly or as a postinstall script
if (require.main === module) {
  install().catch(error => {
    console.error('Unexpected error:', error);
    process.exit(1);
  });
}

module.exports = { install, getPlatform };