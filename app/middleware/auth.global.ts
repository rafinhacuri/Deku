export default defineNuxtRouteMiddleware(async to => {
  const { isLoggedIn, setUserSession, clearUserSession } = useUserSession()

  const res = await useRequestFetch()<goRes>('/server/api/check-session').catch(() => null)

  if(!res){
    clearUserSession()
    if(to.fullPath !== '/login') return navigateTo('/login')
    return
  }
  setUserSession({ username: res.message })

  if(to.fullPath === '/' && !isLoggedIn.value) return navigateTo('/login')

  if(to.fullPath === '/login' && isLoggedIn.value) return navigateTo('/')
})
