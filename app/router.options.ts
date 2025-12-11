import type { RouterConfig } from '@nuxt/schema'

const RouteOptions: RouterConfig = {
  scrollBehavior(to, from, savedPosition) {
    const nuxtApp = useNuxtApp()

    if (savedPosition) {
      nuxtApp.hooks.hookOnce('page:finish', () => {
        setTimeout(() => {
          window.scrollTo({ top: savedPosition.top, behavior: 'auto' })
        }, 50)
      })
      return savedPosition
    }

    if (to.hash) {
      setTimeout(() => {
        let heading = document.querySelector<HTMLElement>(`[id="${to.hash.replace('#', '')}"]`)
        if (!heading) heading = document.querySelector<HTMLElement>(`[href$="${to.hash}"]`)
        if (!heading) return
        window.scrollTo({ top: heading.offsetTop, behavior: 'smooth' })
      })
      return
    }

    if (from.path !== to.path) {
      window.scrollTo({ top: 0, behavior: 'smooth' })
      return
    }

    return { top: 0 }
  },
}

export default RouteOptions
