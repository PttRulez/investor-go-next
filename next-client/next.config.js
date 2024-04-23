// eslint-disable-next-line @typescript-eslint/no-var-requires
const path = require('node:path');

/** @type {import('next').NextConfig} */
const nextConfig = {
  experimental: {
    serverActions: true,
    externalDir: true,
  },
  webpack: (config, { buildId, dev, isServer, defaultLoaders, webpack }) => {
    // Define the root directory where your external files are located
    const externalDir = path.join(__dirname, '../..');

    // Set up an alias to load files from the external directory
    config.resolve.alias['@contracts'] = externalDir;

    return config;
  },
};

module.exports = nextConfig;
