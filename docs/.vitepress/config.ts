import { defineConfig } from 'vitepress'
import { description, version } from '../../package.json'

export default defineConfig({
  cleanUrls: true,
  description:
    description || 'A modern personal finance dashboard to track income, expenses, and insights with a beautiful UI.',
  head: [
    ['meta', { name: 'theme-color', content: '#00c850' }],
    ['meta', { name: 'og:type', content: 'website' }],
    ['meta', { name: 'og:locale', content: 'en-US' }],
    ['meta', { name: 'og:site_name', content: 'Deku' }],
    ['link', { rel: 'icon', href: '/logo.png', type: 'image/png' }],
  ],
  lang: 'en-US',
  lastUpdated: true,
  sitemap: {
    hostname: 'https://deku.curi.dev.br',
  },
  themeConfig: {
    siteTitle: 'Deku',
    logo: '/logo.png',
    search: {
      provider: 'local',
    },
    socialLinks: [
      { icon: 'github', link: 'https://github.com/rafinhacuri/deku' },
    ],
    nav: [
      { text: 'Home', link: '/' },
      {
        text: `V${version}`,
        items: [
          { text: 'Changelog', link: 'https://github.com/rafinhacuri/deku/releases', target: '_blank' },
          { text: 'Report a Bug', link: 'https://github.com/rafinhacuri/deku/issues', target: '_blank' },
          { text: 'Sponsor', link: 'https://github.com/sponsors/rafinhacuri', target: '_blank' },
        ],
      },
    ],
    sidebar: [
      {
        text: '📘 Introduction',
        items: [
          { text: 'Why Deku', link: '/reason' },
          { text: 'Setup', link: '/setup' },
        ],
      },
      {
        text: '📄 Resources',
        items: [
          {
            text: 'License',
            link: 'https://github.com/rafinhacuri/deku/blob/main/LICENSE',
            target: '_blank',
          },
        ],
      },
    ],
    editLink: {
      pattern: 'https://github.com/rafinhacuri/deku/edit/main/docs/:path',
      text: 'Suggest an edit on GitHub',
    },
    docFooter: {
      prev: '← Previous',
      next: 'Next →',
    },
    footer: {
      message: 'Open‑source project licensed under MIT.',
      copyright: '© 2025 Rafael Curi — Deku',
    },
  },
  title: 'Deku',
})
