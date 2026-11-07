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

Set `VITE_PB_API` in `.env.local` (see the root [README](../README.md)) before building — the API URL is baked into the bundle at build time:

- Local dev: `VITE_PB_API="http://127.0.0.1:8090"`
- Production: the API domain you deploy to, e.g. `VITE_PB_API="https://api.example.com"`

