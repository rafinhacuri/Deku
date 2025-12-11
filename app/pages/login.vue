<script setup lang="ts">
definePageMeta({
  layout: false,
})

const { t } = useI18n()

useHead({ title: t('header.title') })
useSeoMeta({ description: t('header.subtitle') })

const toast = useToast()
const { isLoading, start, finish } = useLoadingIndicator()
const { setUserSession } = useUserSession()

const state = ref<Auth>({ email: '', password: '' })

function checkStrength(str: string): { met: boolean, text: string }[] {
  const requirements = [
    { regex: /.{8,}/, text: t('login.week1') },
    { regex: /\d/, text: t('login.week2') },
    { regex: /[a-z]/, text: t('login.week3') },
    { regex: /[A-Z]/, text: t('login.week4') },
  ]

  return requirements.map(req => ({ met: req.regex.test(str), text: req.text }))
}

const strength = computed(() => checkStrength(state.value.password))
const score = computed(() => strength.value.filter(req => req.met).length)

const color = computed(() => {
  if(score.value === 0) return 'neutral'
  if(score.value <= 1) return 'error'
  if(score.value <= 2) return 'warning'
  if(score.value === 3) return 'warning'
  return 'success'
})

const text = computed(() => {
  if(score.value === 0) return t('login.enter')
  if(score.value <= 2) return t('login.weak')
  if(score.value === 3) return t('login.medium')
  return t('login.strong')
})

const show = ref(false)

const modal = ref(false)

async function login(): Promise<void> {
  start()

  const body = AuthSchema.safeParse(state.value)

  if(!body.success){
    for(const e of body.error.issues) toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const res = await $fetch<AuthResponse>('/server/login', { method: 'post', body: body.data })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  finish({ force: true })
  setUserSession({ username: state.value.email, token: res.token })
  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
  await navigateTo('/')
}

async function register(): Promise<void> {
  start()

  const body = AuthSchema.safeParse(state.value)

  if(!body.success){
    for(const e of body.error.issues) toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const res = await $fetch<goRes>('/server/api/user', { method: 'PUT', body: body.data })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  finish({ force: true })
  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
  modal.value = false
}
</script>

<template>
  <section class="relative flex h-screen items-center justify-center overflow-hidden bg-slate-950">
    <div class="pointer-events-none absolute inset-0 opacity-70">
      <div class="absolute -top-32 -left-16 h-72 w-72 rounded-full bg-emerald-500/20 blur-3xl" />
      <div class="absolute -right-16 -bottom-32 h-72 w-72 rounded-full bg-cyan-500/20 blur-3xl" />
      <div class="absolute inset-0 bg-[radial-gradient(circle_at_top,rgba(15,23,42,0.9),transparent_60%)]" />
    </div>

    <div class="relative z-10 flex w-full max-w-md flex-col items-center gap-8 px-4 py-10">
      <div class="flex flex-col items-center gap-7 text-center">
        <div class="relative flex items-center justify-center">
          <div class="absolute -inset-3 rounded-2xl bg-slate-900/80 blur-xl" />
          <div class="relative flex h-16 w-16 items-center justify-center rounded-2xl border border-emerald-400/40 bg-slate-900/80 shadow-xl shadow-emerald-500/20">
            <NuxtImg src="/logo.png" alt="SanchezDNS Logo" width="42" height="42" class="rounded-xl" />
          </div>
        </div>

        <div class="space-y-1">
          <p class="text-xs font-medium tracking-[0.18em] text-emerald-400/80 uppercase">
            {{ t('header.tagline') }}
          </p>
        </div>
      </div>

      <UCard class="relative w-full overflow-hidden rounded-2xl border border-slate-700/60 bg-slate-900/85 shadow-2xl shadow-emerald-500/15 backdrop-blur-xl">
        <div class="pointer-events-none absolute inset-0 bg-linear-to-br from-emerald-500/15 via-slate-900/40 to-cyan-500/10" />

        <template #header>
          <div class="flex items-center justify-between gap-3">
            <div>
              <h2 class="text-base font-semibold text-slate-50">
                {{ t('login.title') }}
              </h2>
              <p class="mt-1 text-[0.72rem] text-slate-400">
                {{ t('login.enter') }}
              </p>
            </div>

            <div class="flex items-center gap-2 text-[0.65rem] text-slate-400">
              <span class="inline-flex h-6 items-center gap-1 rounded-full border border-emerald-500/40 bg-emerald-500/10 px-2 text-[0.65rem] font-medium text-emerald-200">
                <UIcon name="i-lucide-shield-check" class="h-3 w-3" />
                {{ t('login.strong') }}
              </span>
            </div>
          </div>
        </template>

        <div class="relative z-10">
          <UForm :schema="AuthSchema" :state="state" class="space-y-4" @submit="login">
            <UFormField label="Email" name="email">
              <UInput v-model="state.email" icon="i-lucide-mail" class="w-full" />
            </UFormField>

            <UFormField label="Password" name="password">
              <UInput v-model="state.password" class="w-full" icon="i-lucide-lock" :type="show ? 'text' : 'password'">
                <template #trailing>
                  <UButton color="neutral" variant="link" size="sm" :icon="show ? 'i-lucide-eye-off' : 'i-lucide-eye'" :aria-label="show ? 'Hide password' : 'Show password'" :aria-pressed="show" aria-controls="password" @click="show = !show" />
                </template>
              </UInput>
            </UFormField>
            <div class="flex justify-end">
              <UButton variant="ghost" size="sm" @click="modal = true">
                {{ t('login.register') }}
              </UButton>
            </div>

            <UButton type="submit" class="mt-5 flex w-full justify-center">
              {{ t('login.title') }}
            </UButton>
          </UForm>
        </div>
      </UCard>

      <UModal v-model:open="modal" :title="t('login.create_account')" :description="t('login.create_description')" :ui="{ footer: 'justify-end' }">
        <template #body>
          <UForm :schema="AuthSchema" :state="state" class="space-y-4" @submit="register">
            <UFormField label="Email" name="email">
              <UInput v-model="state.email" icon="i-lucide-mail" class="w-full" />
            </UFormField>

            <div class="space-y-2">
              <UFormField label="Password">
                <UInput v-model="state.password" placeholder="Password" :color="color" :type="show ? 'text' : 'password'" :aria-invalid="score < 4" aria-describedby="password-strength" :ui="{ trailing: 'pe-1' }" class="w-full">
                  <template #trailing>
                    <UButton color="neutral" variant="link" size="sm" :icon="show ? 'i-lucide-eye-off' : 'i-lucide-eye'" :aria-label="show ? 'Hide password' : 'Show password'" :aria-pressed="show" aria-controls="password" @click="show = !show" />
                  </template>
                </UInput>
              </UFormField>
              <UProgress :color="color" :indicator="text" :model-value="score" :max="4" size="sm" />
              <p id="password-strength" class="text-sm font-medium">
                {{ text }}. {{ t('login.must_contain') }}
              </p>
              <ul class="space-y-1">
                <li v-for="(req, index) in strength" :key="index" class="flex items-center gap-0.5" :class="req.met ? 'text-success' : 'text-muted'">
                  <UIcon :name="req.met ? 'i-lucide-circle-check' : 'i-lucide-circle-x'" class="size-4 shrink-0" />
                  <span class="text-xs font-light">
                    {{ req.text }}
                    <span class="sr-only">
                      {{ req.met ? t('login.requeriment_meeted') : t('login.requeriment_not_meeted') }}
                    </span>
                  </span>
                </li>
              </ul>
            </div>
          </UForm>
        </template>

        <template #footer>
          <UButton label="Cancel" :loading="isLoading" variant="outline" @click="modal = false" />
          <UButton label="Confirm" :loading="isLoading" @click="register" />
        </template>
      </UModal>
    </div>
  </section>
</template>
