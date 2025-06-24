#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

function cleanup() {
  console.log('Cleaning up claude-gate...');
  
  const binDir = path.join(__dirname, '..', 'bin');
  
  // Files to remove
  const filesToRemove = [
    path.join(binDir, 'claude-gate'),
    path.join(binDir, 'claude-gate.cmd'),
    path.join(binDir, 'claude-gate.ps1'),
    path.join(binDir, 'claude-gate.exe'),
    path.join(binDir, 'claude-gate-bin')
  ];
  
  let removed = 0;
  
  filesToRemove.forEach(file => {
    try {
      if (fs.existsSync(file)) {
        fs.unlinkSync(file);
        removed++;
        console.log(`Removed: ${path.basename(file)}`);
      }
    } catch (error) {
      console.warn(`Warning: Could not remove ${file}:`, error.message);
    }
  });
  
  // Try to remove bin directory if empty
  try {
    if (fs.existsSync(binDir)) {
      const files = fs.readdirSync(binDir);
      if (files.length === 0) {
        fs.rmdirSync(binDir);
        console.log('Removed empty bin directory');
      }
    }
  } catch (error) {
    // Ignore errors when removing directory
  }
  
  if (removed > 0) {
    console.log('✅ claude-gate uninstalled successfully');
  } else {
    console.log('No files to clean up');
  }
}

// Run cleanup
if (require.main === module) {
  try {
    cleanup();
  } catch (error) {
    console.error('Uninstall error:', error.message);
    // Don't exit with error code as this might prevent npm from completing uninstall
  }
}

module.exports = { cleanup };