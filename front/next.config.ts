import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
  reactStrictMode: true,

  async rewrites() {
    const oryUrl = process.env.ORY_SDK_URL || "http://kratos:4433";

    return [
      {
        source: "/self-service/:path*",
        destination: `${oryUrl}/self-service/:path*`,
      },
      {
        source: "/api/kratos/:path*",
        destination: `${oryUrl}/:path*`,
      },
    ];
  },
};

export default nextConfig;
