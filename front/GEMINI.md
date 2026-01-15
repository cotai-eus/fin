# Project Overview

This is a Next.js project that uses Ory Elements for user authentication. It is based on the official Ory quickstart guide for Next.js with the App Router.

The project is structured with the source code in the `src` directory and uses `bun` as the package manager.

## Building and Running

### Development

To run the application in development mode, use the following command:

```bash
bun dev
```

This will start the development server on `http://localhost:3000`.

### Building

To build the application for production, use the following command:

```bash
bun build
```

### Starting

To start a production server, use the following command:

```bash
bun start
```

## Development Conventions

- The project uses TypeScript.
- Styling is done with Tailwind CSS.
- ESLint is used for linting.
- The project follows the Next.js App Router structure.
- All source code is located in the `src` directory.

## Key Files

- `next.config.ts`: Next.js configuration file.
- `package.json`: Defines project dependencies and scripts.
- `src/app`: Contains the application's routes.
- `src/app/layout.tsx`: The main layout of the application.
- `src/app/page.tsx`: The main page of the application.
- `public`: Contains static assets.
- `README.md`: Provides instructions on how to set up and run the project.