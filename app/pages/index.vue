<script setup lang="ts">
definePageMeta({
  colorMode: 'dark',
})

const { t, setLocale, setLocaleCookie, locale } = useI18n()
const { clearUserSession } = useUserSession()
const toast = useToast()
const { isLoading, start, finish } = useLoadingIndicator()

useHead({ title: t('header.title') })
useSeoMeta({ description: t('header.subtitle') })

function changeLanguage(): void {
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

async function logout(): Promise<void> {
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

const salaryTypes = ref([
  { label: t('salary'), id: 'salary' },
  { label: t('salary_types.freelance'), id: 'freelance' },
  { label: t('salary_types.investments'), id: 'investments' },
  { label: t('salary_types.benefits'), id: 'benefits' },
  { label: t('salary_types.scholarship'), id: 'scholarship' },
  { label: t('salary_types.allowance'), id: 'allowance' },
  { label: t('salary_types.others'), id: 'others' },
])

const mes = computed(() => {
  const now = new Date()
  const monthIndex = selectedMonth.value ? months.value.indexOf(selectedMonth.value) : currentMonthIndex
  now.setMonth(monthIndex)
  now.setDate(1)
  now.setHours(0, 0, 0, 0)
  return now.toISOString()
})

const { data, refresh } = await useFetch<SalaryResponse[]>('/server/api/salary', { method: 'GET', query: { month: mes } })

const { data: expenses, refresh: refreshExpenses } = await useFetch<ExpenseResponse[]>('/server/api/expense', { method: 'GET', query: { month: mes } })

const page = ref(1)
const itemsPerPage = ref(10)

const paginetedExpenses = computed(() => {
  if(!expenses.value) return []

  const start = (page.value - 1) * itemsPerPage.value
  const end = start + itemsPerPage.value

  return expenses.value.slice(start, end)
})

const salaryTotal = computed(() => data.value ? data.value.reduce((acc, curr) => acc + curr.vl, 0) : 0)

const stateSalary = ref<Salary>({ vl: 0, type: '', day: undefined, month: '' })

async function addSalaryEntry(): Promise<void> {
  start()

  stateSalary.value.month = mes.value

  const body = SalarySchema.safeParse(stateSalary.value)

  if(!body.success){
    for(const e of body.error.issues) toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const res = await $fetch<goRes>('/server/api/salary', { method: 'PUT', body: body.data })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  refresh()
  finish({ force: true })
  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
  stateSalary.value = { vl: 0, type: '', day: undefined, month: '' }
}

async function deleteSalary(id: number): Promise<void> {
  start()

  const res = await $fetch<goRes>('/server/api/salary', { method: 'DELETE', query: { id } })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  refresh()
  finish({ force: true })
  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
}

const { y } = useWindowScroll({ behavior: 'smooth' })

const salary = useTemplateRef<HTMLElement>('salary')
function scrollToSalary(): void {
  if(salary.value){
    y.value = salary.value.scrollHeight
  }
}

const isEditing = ref(false)
const stateEditSalaryId = ref(0)

function initEdit(line: SalaryResponse): void {
  stateSalary.value.vl = line.vl
  stateSalary.value.type = line.type
  stateSalary.value.day = line.day
  stateEditSalaryId.value = line.id
  isEditing.value = true
  scrollToSalary()
}

function cancelEdit(): void {
  stateSalary.value = { vl: 0, type: '', day: undefined, month: '' }
  stateEditSalaryId.value = 0
  isEditing.value = false
}

async function updateSalary(): Promise<void> {
  start()

  stateSalary.value.month = mes.value

  const body = SalarySchema.safeParse(stateSalary.value)

  if(!body.success){
    for(const e of body.error.issues) toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const res = await $fetch<goRes>('/server/api/salary', { method: 'POST', query: { id: stateEditSalaryId.value }, body: body.data })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  refresh()
  finish({ force: true })
  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
  stateSalary.value = { vl: 0, type: '', day: undefined, month: '' }
  stateEditSalaryId.value = 0
  isEditing.value = false
}

const expensesTypes = ref([
  { label: t('expenses_types.food'), id: 'food' },
  { label: t('expenses_types.shopping'), id: 'shopping' },
  { label: t('expenses_types.transportation'), id: 'transportation' },
  { label: t('expenses_types.supermarket'), id: 'supermarket' },
  { label: t('expenses_types.healthcare'), id: 'healthcare' },
  { label: t('expenses_types.housing'), id: 'housing' },
  { label: t('expenses_types.utilities'), id: 'utilities' },
  { label: t('expenses_types.education'), id: 'education' },
  { label: t('expenses_types.entertainment'), id: 'entertainment' },
  { label: t('expenses_types.bills'), id: 'bills' },
  { label: t('expenses_types.miscellaneous'), id: 'miscellaneous' },
  { label: t('expenses_types.investments'), id: 'investments' },
  { label: t('expenses_types.others'), id: 'others' },
])

const expenseTypeColorMap: Record<string, { chip: string, bar: string }> = {
  food: { chip: 'bg-emerald-400/80', bar: 'bg-emerald-400' },
  shopping: { chip: 'bg-fuchsia-400/80', bar: 'bg-fuchsia-400' },
  transportation: { chip: 'bg-sky-400/80', bar: 'bg-sky-400' },
  supermarket: { chip: 'bg-yellow-400/80', bar: 'bg-yellow-400' },
  healthcare: { chip: 'bg-rose-400/80', bar: 'bg-rose-400' },
  housing: { chip: 'bg-amber-400/80', bar: 'bg-amber-400' },
  utilities: { chip: 'bg-cyan-400/80', bar: 'bg-cyan-400' },
  education: { chip: 'bg-indigo-400/80', bar: 'bg-indigo-400' },
  entertainment: { chip: 'bg-pink-400/80', bar: 'bg-pink-400' },
  bills: { chip: 'bg-orange-400/80', bar: 'bg-orange-400' },
  miscellaneous: { chip: 'bg-slate-400/80', bar: 'bg-slate-400' },
  investments: { chip: 'bg-lime-400/80', bar: 'bg-lime-400' },
  others: { chip: 'bg-zinc-400/80', bar: 'bg-zinc-400' },
}

const paymentsMethods = computed(() => {
  if(!data.value) return []

  const typesInMonth = [...new Set(data.value.map(s => s.type))]

  return typesInMonth.map(typeId => {
    const match = salaryTypes.value.find(s => s.id === typeId)
    return { id: typeId, label: match ? match.label : typeId }
  })
})

const expensesTypesWithData = computed(() => expensesTypes.value.filter(type => expenses.value?.some(e => e.type === type.id)))
const expensesTotal = computed(() => expenses.value?.reduce((acc, curr) => acc + curr.vl, 0) ?? 0)

const stateExpense = ref<Expense>({ vl: 0, type: '', day: undefined, month: '', description: undefined, paymentMethod: '' })

async function addExpense(): Promise<void> {
  start()

  stateExpense.value.month = mes.value

  const body = ExpenseSchema.safeParse(stateExpense.value)

  if(!body.success){
    for(const e of body.error.issues) toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const res = await $fetch<goRes>('/server/api/expense', { method: 'PUT', body: body.data })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  refreshExpenses()
  refresh()
  page.value = 1
  stateExpense.value = { vl: 0, type: '', day: undefined, month: '', description: undefined, paymentMethod: '' }
  finish({ force: true })
  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
}

async function deleteExpense(id: number): Promise<void> {
  start()

  const res = await $fetch<goRes>('/server/api/expense', { method: 'DELETE', query: { id } })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  refreshExpenses()
  refresh()
  page.value = 1
  finish({ force: true })
  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
}

const isEditingExpense = ref(false)
const stateEditExpenseId = ref(0)

const expense = useTemplateRef<HTMLElement>('expense')
function scrollToExpense(): void {
  if(expense.value){
    y.value = expense.value.scrollHeight
  }
}

function initEditExpense(line: ExpenseResponse): void {
  stateExpense.value.vl = line.vl
  stateExpense.value.type = line.type
  stateExpense.value.day = line.day
  stateExpense.value.description = line.description
  stateExpense.value.paymentMethod = line.paymentMethod
  stateEditExpenseId.value = line.id
  isEditingExpense.value = true
  scrollToExpense()
}

function cancelEditExpense(): void {
  stateExpense.value = { vl: 0, type: '', day: undefined, month: '', description: undefined, paymentMethod: '' }
  stateEditExpenseId.value = 0
  isEditingExpense.value = false
}

async function updateExpense(): Promise<void> {
  start()

  stateExpense.value.month = mes.value

  const body = ExpenseSchema.safeParse(stateExpense.value)

  if(!body.success){
    for(const e of body.error.issues) toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const res = await $fetch<goRes>('/server/api/expense', { method: 'POST', query: { id: stateEditExpenseId.value }, body: body.data })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  refreshExpenses()
  refresh()
  page.value = 1
  stateExpense.value = { vl: 0, type: '', day: undefined, month: '', description: undefined, paymentMethod: '' }
  stateEditExpenseId.value = 0
  isEditingExpense.value = false
  finish({ force: true })
  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
}

const leftover = computed(() => salaryTotal.value - expensesTotal.value)

const incomeSourcesSummary = computed(() => {
  const salaries = data.value ?? []

  const typeIds = [...new Set(salaries.map(s => s.type))]

  return typeIds.map(typeId => {
    const income = salaries.filter(s => s.type === typeId).reduce((acc, curr) => acc + curr.vl, 0)

    const spent = expenses.value?.filter(e => e.paymentMethod === typeId).reduce((acc, curr) => acc + curr.vl, 0) ?? 0

    const leftoverPerSource = income - spent
    const match = salaryTypes.value.find(s => s.id === typeId)

    return { id: typeId, label: match ? match.label : typeId, income, spent, leftover: leftoverPerSource }
  })
})
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
        <UCard data-aos="zoom-in" class="relative overflow-hidden rounded-2xl border border-slate-600/40 bg-slate-900/80 backdrop-blur-xl">
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
                {{ t('salary') }}: {{ data ? data.length : 0 }}
              </UBadge>
            </div>
          </div>
        </UCard>

        <UCard data-aos="zoom-in" class="relative overflow-hidden rounded-2xl border border-slate-600/40 bg-slate-900/80 backdrop-blur-xl">
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
                {{ t('expenses') }}: {{ expenses ? expenses.length : 0 }}
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

        <UCard data-aos="zoom-in" class="relative overflow-hidden rounded-2xl border border-slate-600/40 bg-slate-900/80 backdrop-blur-xl">
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
        <UCard data-aos="zoom-in" class="relative overflow-hidden rounded-2xl border border-slate-600/40 bg-slate-900/80 shadow-xl backdrop-blur-xl">
          <div class="pointer-events-none absolute inset-0 bg-linear-to-br from-emerald-500/10 via-emerald-400/5 to-transparent" />

          <div ref="salary" class="relative z-10 space-y-5">
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
                  {{ t('salary_section.type') }} {{ stateSalary.type }}
                </p>
                <USelectMenu v-model="stateSalary.type" label-key="label" value-key="id" :items="salaryTypes" size="md" class="w-full" />
              </div>

              <div class="space-y-1">
                <p class="text-[0.68rem] font-medium tracking-wide text-slate-400 uppercase">
                  {{ t('salary_section.day') }}
                </p>
                <UInputNumber v-model="stateSalary.day" :min="1" :max="31" :step="1" locale="pt-BR" size="md" class="w-full" />
              </div>
            </div>

            <div v-if="isEditing" class="grid grid-cols-1 gap-2 sm:grid-cols-2">
              <UButton :loading="isLoading" class="mt-2 flex w-full items-center justify-center gap-2" color="info" icon="i-lucide-edit" variant="solid" :label="t('salary_section.edit_salary')" @click="updateSalary" />
              <UButton :loading="isLoading" class="mt-2 flex w-full items-center justify-center gap-2" color="error" icon="i-lucide-x" variant="solid" :label="t('cancel')" @click="cancelEdit" />
            </div>
            <UButton v-else :loading="isLoading" class="mt-2 flex w-full items-center justify-center gap-2" color="primary" icon="i-lucide-plus" variant="solid" :label="t('salary_section.add_salary')" @click="addSalaryEntry" />
          </div>
        </UCard>

        <UCard v-if="data && data.length" data-aos="zoom-in" class="relative overflow-hidden rounded-2xl border border-slate-600/40 bg-slate-900/80 shadow-xl backdrop-blur-xl">
          <div class="pointer-events-none absolute inset-0 bg-linear-to-br from-rose-500/12 via-rose-400/5 to-transparent" />

          <div class="relative z-10 space-y-5">
            <div ref="expense" class="flex items-start justify-between gap-3">
              <div>
                <h2 class="text-lg font-semibold tracking-wide text-white">
                  {{ t('expenses_section.title') }}
                </h2>
                <p class="mt-1 text-xs text-slate-400">
                  {{ t('expenses_section.subtitle') }}
                </p>
              </div>

              <div class="flex flex-col items-end gap-1">
                <UBadge class="border border-rose-400/40 bg-rose-500/10 text-rose-300">
                  {{ t('expenses') }}
                </UBadge>
                <p class="text-xs text-slate-400">
                  {{ selectedMonth }}
                </p>
              </div>
            </div>

            <div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
              <div class="space-y-1">
                <p class="text-[0.68rem] font-medium tracking-wide text-slate-400 uppercase">
                  {{ t('expenses_section.value') }}
                </p>
                <UInputNumber v-model="stateExpense.vl" :min="0" :step="0.01" locale="pt-BR" :format-options="{ style: 'currency', currency: 'BRL', currencyDisplay: 'symbol', currencySign: 'accounting', minimumFractionDigits: 2, maximumFractionDigits: 2 }" class="w-full" size="md" />
              </div>

              <div class="space-y-1">
                <p class="text-[0.68rem] font-medium tracking-wide text-slate-400 uppercase">
                  {{ t('expenses_section.type') }}
                </p>
                <USelectMenu v-model="stateExpense.type" label-key="label" value-key="id" :items="expensesTypes" size="md" class="w-full" />
              </div>

              <div class="space-y-1">
                <p class="text-[0.68rem] font-medium tracking-wide text-slate-400 uppercase">
                  {{ t('expenses_section.day') }}
                </p>
                <UInputNumber v-model="stateExpense.day" :min="1" :max="31" :step="1" locale="pt-BR" size="md" class="w-full" />
              </div>

              <div class="space-y-1">
                <p class="text-[0.68rem] font-medium tracking-wide text-slate-400 uppercase">
                  {{ t('expenses_section.payment_method') }}
                </p>
                <USelectMenu v-model="stateExpense.paymentMethod" label-key="label" value-key="id" :items="paymentsMethods" size="md" class="w-full" />
              </div>
            </div>

            <div class="space-y-1">
              <p class="text-[0.68rem] font-medium tracking-wide text-slate-400 uppercase">
                {{ t('expenses_section.description') }}
              </p>
              <UTextarea v-model="stateExpense.description" :rows="3" resize="none" class="w-full" />
            </div>

            <div v-if="isEditingExpense" class="grid grid-cols-1 gap-2 sm:grid-cols-2">
              <UButton :loading="isLoading" class="mt-2 flex w-full items-center justify-center gap-2" color="info" icon="i-lucide-edit" variant="solid" :label="t('expenses_section.edit_expense')" @click="updateExpense" />
              <UButton :loading="isLoading" class="mt-2 flex w-full items-center justify-center gap-2" color="error" icon="i-lucide-x" variant="solid" :label="t('cancel')" @click="cancelEditExpense" />
            </div>
            <UButton v-else :loading="isLoading" class="mt-1 flex w-full items-center justify-center gap-2" color="error" icon="i-lucide-minus-circle" variant="solid" :label="t('expenses_section.add_expense')" @click="addExpense" />
          </div>
        </UCard>
      </div>

      <div class="mt-10 space-y-6">
        <div class="grid grid-cols-1 gap-6 xl:grid-cols-5">
          <UCard data-aos="zoom-in" class="relative overflow-hidden rounded-2xl border border-slate-600/40 bg-slate-900/80 shadow-xl backdrop-blur-xl xl:col-span-2">
            <div class="pointer-events-none absolute inset-0 bg-linear-to-br from-emerald-400/12 via-cyan-300/4 to-transparent" />

            <div class="relative z-10 space-y-4">
              <div class="flex items-start justify-between gap-3">
                <div>
                  <h2 class="text-lg font-semibold tracking-wide text-white">
                    {{ t('salary_section.sources_overview') }}
                  </h2>
                  <p class="mt-1 text-xs text-slate-400">
                    {{ t('salary_section.sources_overview_subtitle') }}
                  </p>
                </div>
              </div>

              <div class="space-y-2">
                <div v-if="!incomeSourcesSummary.length" class="flex flex-col items-center justify-center gap-2 rounded-xl border border-dashed border-slate-600/40 bg-slate-800/40 py-8 text-center">
                  <div class="flex h-10 w-10 items-center justify-center rounded-full border border-slate-600/60 bg-slate-900/80">
                    <UIcon name="i-lucide-piggy-bank" class="h-5 w-5 text-slate-300" />
                  </div>
                  <p class="text-sm font-medium text-slate-200">
                    {{ t('salary_section.sources_empty_title') }}
                  </p>
                  <p class="text-xs text-slate-400">
                    {{ t('salary_section.sources_empty_subtitle') }}
                  </p>
                </div>

                <div v-else>
                  <div v-for="source in incomeSourcesSummary" :key="source.id" class="space-y-2 rounded-xl bg-slate-900/70 px-3 py-2 ring-1 ring-slate-700/50">
                    <div class="flex items-center justify-between text-xs">
                      <div class="flex items-center gap-2">
                        <span class="flex h-7 w-7 items-center justify-center rounded-full bg-emerald-400/10 ring-1 ring-emerald-400/40">
                          <UIcon name="i-lucide-coins" class="h-3.5 w-3.5 text-emerald-300" />
                        </span>
                        <div class="flex flex-col">
                          <span class="text-sm font-medium text-slate-100">
                            {{ source.label }}
                          </span>
                          <span class="text-[0.7rem] text-slate-400">
                            {{ t('salary_section.source_income_label') }}
                            {{ source.income.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' }) }}
                          </span>
                        </div>
                      </div>

                      <div class="text-right">
                        <p class="text-[0.65rem] text-slate-400">
                          {{ t('salary_section.source_leftover_label') }}
                        </p>
                        <p class="text-sm font-semibold" :class="source.leftover >= 0 ? 'text-emerald-300' : 'text-rose-400'">
                          {{ source.leftover.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' }) }}
                        </p>
                      </div>
                    </div>

                    <div class="space-y-1 text-[0.7rem] text-slate-300">
                      <div class="flex items-center justify-between">
                        <span>{{ t('salary_section.source_spent_label') }}</span>
                        <span class="text-rose-300">
                          {{ source.spent.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' }) }}
                        </span>
                      </div>

                      <div class="h-2 w-full overflow-hidden rounded-full bg-slate-800">
                        <!-- eslint-disable-next-line vue/no-restricted-v-bind -->
                        <div v-if="source.leftover !== 0" class="h-2 rounded-full" :class="source.leftover > 0 ? 'bg-emerald-400' : 'bg-rose-400'" :style="{ width: source.income ? Math.max( 4, Math.min(100, Math.round(Math.abs(source.leftover / source.income) * 100)), ) + '%' : '0%', }" />
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </UCard>

          <UCard data-aos="zoom-in" class="relative overflow-hidden rounded-2xl border border-slate-600/40 bg-slate-900/80 shadow-xl backdrop-blur-xl xl:col-span-3">
            <div class="pointer-events-none absolute inset-0 bg-linear-to-br from-rose-400/18 via-rose-300/4 to-transparent" />

            <div class="relative z-10 space-y-5">
              <div class="flex flex-wrap items-start justify-between gap-3">
                <div>
                  <h2 class="text-lg font-semibold tracking-wide text-white">
                    {{ t('expenses_section.analysis_title') }}
                  </h2>
                  <p class="mt-1 text-xs text-slate-400">
                    {{ t('expenses_section.analysis_subtitle') }}
                  </p>
                </div>

                <div class="flex items-center gap-3">
                  <div class="flex flex-col text-right text-xs">
                    <span class="text-slate-400">
                      {{ t('expenses_section.total_month') }}
                    </span>
                    <span class="text-sm font-semibold text-rose-300">
                      {{ expensesTotal.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' }) }}
                    </span>
                  </div>

                  <div class="relative flex h-16 w-16 items-center justify-center rounded-full border border-dashed border-rose-400/50 bg-rose-500/5">
                    <div class="flex h-12 w-12 items-center justify-center rounded-full bg-slate-950/80 ring-1 ring-rose-400/40">
                      <span class="text-xs font-semibold text-rose-300">
                        {{ salaryTotal > 0 ? Math.round((expensesTotal / salaryTotal) * 100) : 0 }}%
                      </span>
                    </div>
                    <span class="pointer-events-none absolute right-1 -bottom-1 rounded-full bg-rose-500 px-1.5 py-0.5 text-[0.6rem] font-semibold text-white shadow-lg">
                      {{ t('expenses_section.usage') }}
                    </span>
                  </div>
                </div>
              </div>

              <div class="space-y-3">
                <div class="flex items-center justify-between text-[0.7rem] text-slate-400">
                  <span>{{ t('expenses_section.by_category') }}</span>
                  <span v-if="expenses && expenses.length" class="italic">
                    {{ t('expenses_section.tip_focus') }}
                  </span>
                </div>

                <div v-if="!expenses || !expenses.length" class="flex flex-col items-center justify-center gap-2 rounded-xl border border-dashed border-slate-600/40 bg-slate-800/40 py-10 text-center">
                  <div class="flex h-10 w-10 items-center justify-center rounded-full border border-slate-600/60 bg-slate-900/80">
                    <UIcon name="i-lucide-pie-chart" class="h-5 w-5 text-slate-300" />
                  </div>
                  <p class="text-sm font-medium text-slate-200">
                    {{ t('expenses_section.empty_chart_title') }}
                  </p>
                  <p class="text-xs text-slate-400">
                    {{ t('expenses_section.empty_chart_subtitle') }}
                  </p>
                </div>

                <div v-else class="max-h-80 space-y-2 overflow-y-auto pr-2">
                  <div v-for="type in expensesTypesWithData" :key="type.id" class="space-y-2 rounded-xl bg-slate-900/70 px-3 py-2 ring-1 ring-slate-700/50">
                    <div class="flex items-center justify-between text-xs">
                      <div class="flex items-center gap-2">
                        <span class="h-1.5 w-4 rounded-full" :class="expenseTypeColorMap[type.id]?.chip || 'bg-rose-400/80'" />
                        <span class="font-medium text-slate-100">
                          {{ type.label }}
                        </span>
                      </div>

                      <div class="flex flex-col text-right">
                        <span class="text-[0.7rem] text-slate-400">
                          {{ Math.floor( expenses.filter(e => e.type === type.id).reduce((acc, curr) => acc + curr.vl, 0) / (expensesTotal || 1) * 100,) }}%
                        </span>
                        <span class="text-[0.7rem] font-medium text-slate-200">
                          {{ expenses.filter(e => e.type === type.id).reduce((acc, curr) => acc + curr.vl, 0).toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' }) }}
                        </span>
                      </div>
                    </div>

                    <div class="overflow-hiddxen h-2 w-full rounded-full bg-slate-800">
                      <!-- eslint-disable-next-line vue/no-restricted-v-bind -->
                      <div class="h-2 rounded-full" :class="expenseTypeColorMap[type.id]?.bar || 'bg-rose-400'" :style="{ width: expensesTotal ? Math.max( 4, Math.round( expenses .filter(e => e.type === type.id) .reduce((acc, curr) => acc + curr.vl, 0) / expensesTotal * 100, ), ) + '%' : '0%', }" />
                    </div>

                    <div v-if="paymentsMethods.length" class="mt-1 space-y-0.5">
                      <template v-for="method in paymentsMethods" :key="method.id + '-' + type.id">
                        <div v-if="expenses.some(e => e.type === type.id && e.paymentMethod === method.id)" class="flex items-center justify-between text-[0.65rem] text-slate-400">
                          <span class="flex items-center gap-1">
                            <span class="h-1 w-1 rounded-full bg-slate-500" />
                            <span>{{ method.label || t('expenses_section.no_payment_method') }}</span>
                          </span>
                          <span class="text-[0.65rem] text-slate-200">
                            {{ expenses.filter(e => e.type === type.id && e.paymentMethod === method.id).reduce((acc, curr) => acc + curr.vl, 0).toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' }) }}
                          </span>
                        </div>
                      </template>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </UCard>
        </div>

        <UCard data-aos="zoom-in" class="relative overflow-hidden rounded-2xl border border-slate-600/40 bg-slate-900/80 shadow-xl backdrop-blur-xl">
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
                {{ t('salary') }}: {{ data ? data.length : 0 }}
              </UBadge>
            </div>

            <div class="max-h-72 space-y-3 overflow-y-auto pr-2">
              <div v-if="data && data.length === 0" class="flex flex-col items-center justify-center gap-2 rounded-xl border border-dashed border-slate-600/40 bg-slate-800/40 py-10 text-center">
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
                <div v-for="(item, index) in data" :key="index" class="flex items-center justify-between gap-3 rounded-xl border border-slate-600/40 bg-slate-900/70 px-3 py-3 transition hover:border-cyan-400/60 hover:bg-slate-900">
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
                          {{ salaryTypes.find(type => type.id === item.type)?.label || item.type }}
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
                    <UButton v-if="!isEditing" :loading="isLoading" color="info" variant="ghost" icon="i-lucide-pencil" size="xs" @click="initEdit(item)" />
                    <UButton v-else-if="isEditing && stateEditSalaryId === item.id" :loading="isLoading" color="error" variant="ghost" icon="i-lucide-x" size="xs" @click="cancelEdit" />
                    <UPopover v-if="!isEditing" arrow>
                      <UButton :loading="isLoading" color="error" variant="ghost" icon="i-lucide-trash" size="xs" />

                      <template #content="{close}">
                        <div class="max-w-xs space-y-3 p-3">
                          <p class="text-sm text-slate-200">
                            {{ t('salary_section.delete_confirmation') }}
                          </p>
                          <div class="flex items-center justify-end gap-2">
                            <UButton :loading="isLoading" variant="solid" color="error" size="sm" :label="t('confirm')" @click="deleteSalary(item.id), close()" />
                          </div>
                        </div>
                      </template>
                    </UPopover>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </UCard>

        <UCard data-aos="zoom-in" class="relative overflow-hidden rounded-2xl border border-slate-600/40 bg-slate-900/80 shadow-xl backdrop-blur-xl">
          <div class="pointer-events-none absolute inset-0 bg-linear-to-br from-rose-400/15 via-rose-300/5 to-transparent" />

          <div class="relative z-10 space-y-4">
            <div class="flex items-start justify-between gap-3">
              <div>
                <h2 class="text-lg font-semibold tracking-wide text-white">
                  {{ t('expenses_section.history') }}
                </h2>
                <p class="mt-1 text-xs text-slate-400">
                  {{ t('expenses_section.history_subtitle') }}
                </p>
              </div>

              <div class="flex flex-col items-end gap-1">
                <UBadge class="mt-1 border border-rose-400/40 bg-rose-500/10 text-rose-200">
                  {{ t('expenses') }}: {{ expenses ? expenses.length : 0 }}
                </UBadge>
                <p class="text-xs text-slate-400">
                  {{ expensesTotal.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' }) }}
                </p>
              </div>
            </div>

            <div class="space-y-3">
              <div v-if="!expenses || !expenses.length" class="flex flex-col items-center justify-center gap-2 rounded-xl border border-dashed border-slate-600/40 bg-slate-800/40 py-10 text-center">
                <div class="flex h-10 w-10 items-center justify-center rounded-full border border-slate-600/60 bg-slate-900/80">
                  <UIcon name="i-lucide-credit-card" class="h-5 w-5 text-slate-300" />
                </div>
                <p class="text-sm font-medium text-slate-200">
                  {{ t('expenses_section.empty_title') }}
                </p>
                <p class="text-xs text-slate-400">
                  {{ t('expenses_section.empty_subtitle') }}
                </p>
              </div>

              <div v-else>
                <div v-for="(item, index) in paginetedExpenses" :key="index" class="flex flex-col gap-3 rounded-xl border border-slate-600/40 bg-slate-900/70 px-3 py-3 transition hover:border-rose-400/60 hover:bg-slate-900">
                  <div class="flex items-center justify-between gap-3">
                    <div class="flex items-center gap-3">
                      <div class="flex h-9 w-9 items-center justify-center rounded-full bg-rose-400/10 ring-1 ring-rose-400/40">
                        <UIcon name="i-lucide-shopping-bag" class="h-4 w-4 text-rose-300" />
                      </div>

                      <div class="flex flex-col">
                        <span class="text-sm font-semibold text-white">
                          {{ item.vl.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' }) }}
                        </span>

                        <div class="mt-1 flex flex-wrap items-center gap-2 text-xs">
                          <UBadge class="border text-xs font-medium" :class="expenseTypeColorMap[item.type]?.chip || 'bg-slate-800/80 border-slate-500/60 text-slate-100'">
                            {{ expensesTypes.find(type => type.id === item.type)?.label || t('expenses_types.others') }}
                          </UBadge>

                          <UBadge class="border border-slate-500/60 bg-slate-800/80 text-slate-100">
                            {{ paymentsMethods.find(method => method.id === item.paymentMethod)?.label || t('expenses_section.no_payment_method') }}
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
                      <UButton v-if="!isEditingExpense" :loading="isLoading" color="info" variant="ghost" icon="i-lucide-pencil" size="xs" @click="initEditExpense(item)" />
                      <UButton v-else-if="isEditingExpense && stateEditExpenseId === item.id" :loading="isLoading" color="error" variant="ghost" icon="i-lucide-x" size="xs" @click="cancelEditExpense" />
                      <UPopover v-if="!isEditingExpense" arrow>
                        <UButton :loading="isLoading" color="error" variant="ghost" icon="i-lucide-trash" size="xs" />

                        <template #content="{ close }">
                          <div class="max-w-xs space-y-3 p-3">
                            <p class="text-sm text-slate-200">
                              {{ t('expenses_section.delete_confirmation') }}
                            </p>
                            <div class="flex items-center justify-end gap-2">
                              <UButton :loading="isLoading" variant="solid" color="error" size="sm" :label="t('confirm')" @click="deleteExpense(item.id), close()" />
                            </div>
                          </div>
                        </template>
                      </UPopover>
                    </div>
                  </div>

                  <p v-if="item.description" class="line-clamp-2 text-xs text-slate-300">
                    {{ item.description }}
                  </p>
                </div>
              </div>
            </div>
            <UPagination v-if="expenses && expenses.length > 10" v-model:page="page" :items-per-page="itemsPerPage" :total="expenses.length" />
          </div>
        </UCard>
      </div>
    </UContainer>
  </section>
</template>
