/**
 * claude-gate NPM package
 * This file is not meant to be used programmatically.
 * Install this package globally and use the 'claude-gate' command.
 */

module.exports = {
  name: 'claude-gate',
  version: '0.1.0',
  description: 'OAuth proxy for Anthropic Claude API',
  
  // Display a helpful message if someone tries to require() this package
  __esModule: true,
  default: function() {
    console.error(`
claude-gate is a command-line tool and should be installed globally:

  npm install -g claude-gate

Then use it from the command line:

  claude-gate start
  claude-gate auth login
  claude-gate --help

For more information, visit:
https://github.com/ml0-1337/claude-gate
`);
    process.exit(1);
  }
};

// Show message if this file is run directly
if (require.main === module) {
  module.exports.default();
}