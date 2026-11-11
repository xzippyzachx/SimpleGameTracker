# SimpleGameTracker

This template should help get you started developing with Vue 3 in Vite.

## Recommended IDE Setup

[VSCode](https://code.visualstudio.com/) + [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (and disable Vetur).

## Customize configuration

See [Vite Configuration Reference](https://vite.dev/config/).

## Project Setup

```sh
npm install
```

### Compile and Hot-Reload for Development

```sh
npm run dev
```

### Compile and Minify for Production

```sh
npm run build
```

### Deploy to GitHub Pages

```sh
npm run deploy
```

Builds `dist/` and force-pushes it to the `gh-pages` branch. A custom domain (if configured) comes from the `CNAME` file in `public/`, which publishes with every deploy.

### Backend API

The API URL is baked into the bundle at build time, so it is configured per mode instead of being edited by hand:

| Command | Env file | Purpose |
| --- | --- | --- |
| `npm run dev` | `.env.local` | Your machine only (gitignored) |
| `npm run build` | `.env.production` | Production API (gitignored) |

Both files are gitignored. Copy `.env.example` to create them, then uncomment the value you need:

```sh
cp .env.example .env.local        # uncomment VITE_PB_API="http://127.0.0.1:8090"
cp .env.example .env.production   # uncomment VITE_PB_API="https://api.example.com"
```

`.env.local` is loaded in every mode but ranks below `.env.production`, so a production build always uses the production API. A plain `.env` file is deliberately not used: `.env.local` would override it in production builds too.

The build fails with a clear error if the mode's env file is missing or `VITE_PB_API` is unset, so a bundle can never ship with a localhost or `undefined` API base. That is why `.env.example` ships with every value commented out — copying it without editing fails the build instead of silently deploying the localhost API.

