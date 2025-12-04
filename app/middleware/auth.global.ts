export default defineNuxtRouteMiddleware(async to => {
  const { isLoggedIn, setUserSession, clearUserSession } = useUserSession()

  const res = await useRequestFetch()<{ username: string }>('/server/api/check-session').catch(() => null)

  if(!res){
    clearUserSession()
    if(to.fullPath !== '/') return navigateTo('/')
    return
  }
  setUserSession({ username: res.username })

  if(to.fullPath === '/' && !isLoggedIn.value) return navigateTo('/login')

  if(to.fullPath === '/login' && isLoggedIn.value) return navigateTo('/')
})
