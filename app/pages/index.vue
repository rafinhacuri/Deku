<script setup lang="ts">
import type { Salary } from '~~/shared/types/goServer'
import { SalarySchema } from '~/utils/salary.schema'

definePageMeta({
  colorMode: 'dark',
})

const { t, setLocale, setLocaleCookie, locale } = useI18n()
const { clearUserSession } = useUserSession()
const toast = useToast()
const { isLoading, start, finish } = useLoadingIndicator()

useHead({ title: t('header.title') })
useSeoMeta({ description: t('header.subtitle') })
defineOgImageComponent('Home', { title: t('header.title'), subtitle: t('header.subtitle') })

function changeLanguage(){
  if(locale.value === 'en'){
    setLocale('pt')
    setLocaleCookie('pt')
  }
  else {
    setLocale('en')
    setLocaleCookie('en')
  }
}

const isEn = computed(() => locale.value === 'en')

async function logout(){
  clearUserSession()
  await navigateTo('/login')
}

const months = ref([
  t('months.january'),
  t('months.february'),
  t('months.march'),
  t('months.april'),
  t('months.may'),
  t('months.june'),
  t('months.july'),
  t('months.august'),
  t('months.september'),
  t('months.october'),
  t('months.november'),
  t('months.december'),
])

const currentMonthIndex = new Date().getMonth()
const monthsToShow = ref(months.value.slice(0, currentMonthIndex + 1).toReversed())
const selectedMonth = ref(monthsToShow.value[0])

const expensesTypes = ref([
  t('expenses_types.food'),
  t('expenses_types.shopping'),
  t('expenses_types.transportation'),
  t('expenses_types.healthcare'),
  t('expenses_types.housing'),
  t('expenses_types.utilities'),
  t('expenses_types.education'),
  t('expenses_types.entertainment'),
  t('expenses_types.bills'),
  t('expenses_types.miscellaneous'),
  t('expenses_types.investments'),
  t('expenses_types.others'),
])

const salaryTypes = ref([
  t('salary'),
  t('salary_types.freelance'),
  t('salary_types.investments'),
  t('salary_types.benefits'),
  t('salary_types.others'),
])

const salary = ref<Salary[]>([])
const salaryTotal = computed(() => salary.value.reduce((acc, curr) => acc + curr.vl, 0))

const stateSalary = ref<Salary>({ vl: 0, type: '', day: undefined, month: '' })

async function addSalaryEntry(){
  start()
  const now = new Date()
  const monthIndex = selectedMonth.value ? months.value.indexOf(selectedMonth.value) : currentMonthIndex
  now.setMonth(monthIndex)
  now.setDate(1)
  now.setHours(0, 0, 0, 0)

  stateSalary.value.month = now.toISOString()

  const body = SalarySchema.safeParse(stateSalary.value)

  if(!body.success){
    for(const e of body.error.issues) toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const res = await $fetch<goRes>('/server/api/salary', { method: 'PUT', body: body.data })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  // refresh()
  finish({ force: true })
  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
  stateSalary.value = { vl: 0, type: '', day: undefined, month: '' }
}

const expenses = ref<{ value: number, type: string | undefined, description: string | undefined, day: number | undefined, paymentMethod: string, user: string, month: Date }[]>([])
const expensesTotal = computed(() => expenses.value.reduce((acc, curr) => acc + curr.value, 0))

const leftover = computed(() => salaryTotal.value - expensesTotal.value)
</script>

<template>
  <section>
    <UHeader>
      <template #title>
        <NuxtImg src="/logo.png" alt="Deku Logo" height="32" />
        {{ t('header.title') }}
      </template>
      <template #right>
        <div class="flex items-center space-x-3">
          <USelectMenu v-model="selectedMonth" :items="monthsToShow" size="sm" :placeholder="t('header.month_placeholder')" />
          <UButton :loading="isLoading" variant="outline" size="sm" :icon="isEn ? 'i-circle-flags-us' : 'i-circle-flags-br'" :color="isEn ? 'info' : 'primary'" @click="changeLanguage" />
          <UButton :loading="isLoading" variant="outline" color="neutral" size="sm" icon="i-lucide-log-out" @click="logout()" />
        </div>
      </template>
    </UHeader>

    <UContainer v-if="selectedMonth" class="my-8">
      <div class="grid grid-cols-1 gap-4 md:grid-cols-3">
        <UCard class="relative overflow-hidden rounded-2xl border border-slate-600/40 bg-slate-900/80 backdrop-blur-xl">
          <div class="pointer-events-none absolute inset-0 bg-linear-to-br from-emerald-400/20 to-transparent" />
          <div class="relative z-10 flex justify-between">
            <div>
              <p class="text-[0.7rem] tracking-widest text-slate-400 uppercase">
                {{ t('salary') }}
              </p>
              <p class="text-xl font-semibold text-emerald-400">
                {{ salaryTotal.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' }) }}
              </p>
              <p class="text-xs text-slate-400">
                {{ t('everything_earned') }}
              </p>
              <UBadge class="mt-2 border border-slate-500/60 bg-slate-800/70 text-white">
                {{ t('salary') }}: {{ salary.length }}
              </UBadge>
            </div>
          </div>
        </UCard>

        <UCard class="relative overflow-hidden rounded-2xl border border-slate-600/40 bg-slate-900/80 backdrop-blur-xl">
          <div class="pointer-events-none absolute inset-0 bg-linear-to-br from-rose-400/20 to-transparent" />

          <div class="relative z-10 flex justify-between">
            <div>
              <p class="text-[0.7rem] tracking-widest text-slate-400 uppercase">
                {{ t('expenses') }}
              </p>

              <p class="text-xl font-semibold text-rose-400">
                {{ expensesTotal.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' }) }}
              </p>

              <p class="text-xs text-slate-400">
                {{ t('everything_spent') }}
              </p>

              <UBadge class="mt-2 border border-slate-500/60 bg-slate-800/70 text-white">
                {{ t('expenses') }}: {{ expenses.length > 0 ? 1 : 0 }}
              </UBadge>
            </div>

            <div class="flex items-center justify-center">
              <div class="relative flex h-14 w-14 items-center justify-center rounded-full border border-dashed border-slate-500/50">
                <span class="absolute top-1 left-4 h-2 w-2 rounded-full bg-white shadow" />
                <span class="text-sm font-medium text-rose-400">
                  {{ salaryTotal > 0 ? Math.round((expensesTotal / salaryTotal) * 100) : 0 }}%
                </span>
              </div>
            </div>
          </div>
        </UCard>

        <UCard class="relative overflow-hidden rounded-2xl border border-slate-600/40 bg-slate-900/80 backdrop-blur-xl">
          <div class="pointer-events-none absolute inset-0 bg-linear-to-br from-cyan-400/20 to-transparent" />

          <div class="relative z-10 flex justify-between">
            <div>
              <p class="text-[0.7rem] tracking-widest text-slate-400 uppercase">
                {{ t('leftover') }}
              </p>

              <p class="text-xl font-semibold" :class="[ leftover >= 0 ? 'text-emerald-300' : 'text-rose-400' ]">
                {{ leftover.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' }) }}
              </p>

              <p class="text-xs text-slate-400">
                {{ leftover >= 0 ? t('positive_balance') : t('negative_balance') }}
              </p>
            </div>

            <div class="flex items-center justify-center">
              <div class="relative flex h-14 w-14 items-center justify-center rounded-full border border-dashed border-slate-500/50">
                <span class="absolute top-1 left-4 h-2 w-2 rounded-full bg-white shadow" />
                <span class="text-sm font-medium" :class="[ leftover >= 0 ? 'text-emerald-300' : 'text-rose-400' ]">
                  {{ salaryTotal > 0 ? Math.round((leftover / salaryTotal) * 100) : 0 }}%
                </span>
              </div>
            </div>
          </div>
        </UCard>
      </div>

      <div class="mt-8 grid grid-cols-1 gap-6 lg:grid-cols-2">
        <UCard class="relative overflow-hidden rounded-2xl border border-slate-600/40 bg-slate-900/80 shadow-xl backdrop-blur-xl">
          <div class="pointer-events-none absolute inset-0 bg-linear-to-br from-emerald-500/10 via-emerald-400/5 to-transparent" />

          <div class="relative z-10 space-y-5">
            <div class="flex items-start justify-between gap-3">
              <div>
                <h2 class="text-lg font-semibold tracking-wide text-white">
                  {{ t('salary_section.title') }}
                </h2>
                <p class="mt-1 text-xs text-slate-400">
                  {{ t('everything_earned') }}
                </p>
              </div>

              <div class="flex flex-col items-end gap-1">
                <UBadge class="border border-emerald-400/40 bg-emerald-500/10 text-emerald-300">
                  {{ t('salary') }}
                </UBadge>
                <p class="text-xs text-slate-400">
                  {{ selectedMonth }}
                </p>
              </div>
            </div>

            <div class="grid grid-cols-1 gap-3 sm:grid-cols-3">
              <div class="space-y-1">
                <p class="text-[0.68rem] font-medium tracking-wide text-slate-400 uppercase">
                  {{ t('salary_section.value') }} {{ stateSalary.vl }}
                </p>
                <UInputNumber v-model="stateSalary.vl" :min="0" :step="0.01" locale="pt-BR" :format-options="{ style: 'currency', currency: 'BRL', currencyDisplay: 'symbol', currencySign: 'accounting', minimumFractionDigits: 2, maximumFractionDigits: 2 }" class="w-full" size="md" />
              </div>

              <div class="space-y-1">
                <p class="text-[0.68rem] font-medium tracking-wide text-slate-400 uppercase">
                  {{ t('salary_section.type') }}
                </p>
                <USelectMenu v-model="stateSalary.type" :items="salaryTypes" size="md" :placeholder="t('salary_section.type')" class="w-full" />
              </div>

              <div class="space-y-1">
                <p class="text-[0.68rem] font-medium tracking-wide text-slate-400 uppercase">
                  {{ t('salary_section.day') }}
                </p>
                <UInputNumber v-model="stateSalary.day" :min="1" :max="31" :step="1" locale="pt-BR" size="md" class="w-full" />
              </div>
            </div>

            <UButton :loading="isLoading" class="mt-2 flex w-full items-center justify-center gap-2" color="primary" icon="i-lucide-plus" variant="solid" @click="addSalaryEntry">
              <span>{{ t('salary_section.add_salary') }}</span>
            </UButton>
          </div>
        </UCard>

        <UCard class="relative overflow-hidden rounded-2xl border border-slate-600/40 bg-slate-900/80 shadow-xl backdrop-blur-xl">
          <div class="pointer-events-none absolute inset-0 bg-linear-to-br from-cyan-400/15 via-cyan-300/5 to-transparent" />

          <div class="relative z-10 space-y-4">
            <div class="flex items-start justify-between gap-3">
              <div>
                <h2 class="text-lg font-semibold tracking-wide text-white">
                  {{ t('salary_section.history') }}
                </h2>
                <p class="mt-1 text-xs text-slate-400">
                  {{ t('everything_earned') }}
                </p>
              </div>

              <UBadge class="mt-1 border border-cyan-400/40 bg-cyan-500/10 text-cyan-200">
                {{ t('salary') }}: {{ salary.length }}
              </UBadge>
            </div>

            <div class="max-h-72 space-y-3 overflow-y-auto pr-2">
              <div v-if="salary.length === 0" class="flex flex-col items-center justify-center gap-2 rounded-xl border border-dashed border-slate-600/40 bg-slate-800/40 py-10 text-center">
                <div class="flex h-10 w-10 items-center justify-center rounded-full border border-slate-600/60 bg-slate-900/80">
                  <UIcon name="i-lucide-wallet" class="h-5 w-5 text-slate-300" />
                </div>
                <p class="text-sm font-medium text-slate-200">
                  {{ t('salary_section.empty_title') }}
                </p>
                <p class="text-xs text-slate-400">
                  {{ t('salary_section.empty_subtitle') }}
                </p>
              </div>

              <div v-else>
                <div v-for="(item, index) in salary" :key="index" class="flex items-center justify-between gap-3 rounded-xl border border-slate-600/40 bg-slate-900/70 px-3 py-3 transition hover:border-cyan-400/60 hover:bg-slate-900">
                  <div class="flex items-center gap-3">
                    <div class="flex h-9 w-9 items-center justify-center rounded-full bg-cyan-400/10 ring-1 ring-cyan-400/40">
                      <UIcon name="i-lucide-coins" class="h-4 w-4 text-cyan-300" />
                    </div>

                    <div class="flex flex-col">
                      <span class="text-sm font-semibold text-white">
                        {{ item.vl.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' }) }}
                      </span>

                      <div class="mt-1 flex flex-wrap items-center gap-2 text-xs">
                        <UBadge class="border border-slate-500/60 bg-slate-800/80 text-slate-100">
                          {{ item.type || t('salary_types.others') }}
                        </UBadge>

                        <span class="text-slate-400">
                          <span v-if="item.day">
                            {{ item.day.toString().padStart(2, '0') }}/{{ (new Date(item.month).getMonth() + 1).toString().padStart(2, '0') }}/{{ new Date(item.month).getFullYear() }}
                          </span>
                          <span v-else>
                            {{ (new Date(item.month).getMonth() + 1).toString().padStart(2, '0') }}/{{ new Date(item.month).getFullYear() }}
                          </span>
                        </span>
                      </div>
                    </div>
                  </div>

                  <div class="flex items-center gap-2">
                    <UButton :loading="isLoading" color="info" variant="ghost" icon="i-lucide-pencil" size="xs" />
                    <UButton :loading="isLoading" color="error" variant="ghost" icon="i-lucide-trash" size="xs" />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </UCard>
      </div>
    </UContainer>
  </section>
</template>
