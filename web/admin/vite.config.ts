import { defineConfig, type Plugin } from 'vite'
import { tanstackStart } from '@tanstack/react-start/plugin/vite'
import viteReact from '@vitejs/plugin-react'

function faviconRedirect(): Plugin {
  const redirect = (
    req: { url?: string },
    res: { writeHead: (code: number, headers: Record<string, string>) => void; end: () => void },
    next: () => void,
  ) => {
    if (req.url?.split('?')[0] === '/favicon.ico') {
      res.writeHead(302, { Location: '/favicon.svg' })
      res.end()
      return
    }
    next()
  }

  return {
    name: 'favicon-redirect',
    configureServer(server) {
      server.middlewares.use(redirect)
    },
    configurePreviewServer(server) {
      server.middlewares.use(redirect)
    },
  }
}

export default defineConfig({
  server: {
    port: 5174,
  },
  resolve: {
    tsconfigPaths: true,
  },
  plugins: [
    faviconRedirect(),
    tanstackStart(),
    viteReact(),
  ],
})
