export default createGlobalState(() => {
  const userInfo = ref({ username: '' })
  const sessionToken = useCookie('session')

  const isLoggedIn = computed(() => userInfo.value.username !== '')
  const user = computed(() => ({
    username: userInfo.value.username,
    token: sessionToken.value,
  }))

  const setUserSession = ({ username, token }: { username: string, token?: string }) => {
    userInfo.value.username = username

    if(token) sessionToken.value = token
  }

  const clearUserSession = () => {
    userInfo.value = { username: '' }
    sessionToken.value = undefined
  }

  return { user, isLoggedIn, setUserSession, clearUserSession }
})
