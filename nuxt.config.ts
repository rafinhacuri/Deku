export default defineNuxtConfig({
  modules: ['@nuxt/ui', '@nuxtjs/seo', '@vueuse/nuxt', 'nuxt-security', '@nuxtjs/i18n'],
  $development: {
    security: { headers: { crossOriginEmbedderPolicy: 'unsafe-none' } },
  },
  devtools: { enabled: true },
  app: { head: { templateParams: { separator: '•' } } },
  css: ['~/assets/main.css'],
  site: {
    name: 'Deku',
    description: '💰 Web application to control your income',
    indexable: true,
  },
  compatibilityDate: '2025-12-02',
  i18n: {
    defaultLocale: 'en',
    locales: [
      { code: 'en', language: 'en-US', name: 'English (US)', file: 'en-US.json', flag: 'flag:us-4x3' },
      { code: 'pt', language: 'pt-BR', name: 'Português (BR)', file: 'pt-BR.json', flag: 'flag:br-4x3' },
    ],
    skipSettingLocaleOnNavigate: true,
  },
  linkChecker: { enabled: false },
  security: {
    headers: {
      contentSecurityPolicy: { 'default-src': ['\'self\''], 'img-src': ['\'self\'', 'data:', 'blob:'] },
    },
  },
  sitemap: { enabled: false },
})
