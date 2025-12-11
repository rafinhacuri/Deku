<script setup lang="ts">
const { t } = useI18n()

const model = defineModel<Date | null | undefined>({ default: undefined })

const startOfMonth = (date: Date): Date => new Date(date.getFullYear(), date.getMonth(), 1)

const sameMonth = (date: Date | null | undefined, year: number, month: number): boolean => !!date && date.getFullYear() === year && date.getMonth() === month

const initialYear = computed(() => model.value?.getFullYear() ?? new Date().getFullYear())
const year = ref<number>(initialYear.value)

watch(() => model.value, value => {
  if (value) year.value = value.getFullYear()
})

const now = new Date()
const maxYear = now.getFullYear()
const maxMonth = now.getMonth()

function isAfterMax(year: number, month: number): boolean {
  return year > maxYear || (year === maxYear && month > maxMonth)
}

const monthFmtLong = new Intl.DateTimeFormat('pt-BR', { month: 'long', year: 'numeric' })

const label = computed(() => {
  if (model.value) return monthFmtLong.format(model.value)
  return t('header.month_placeholder')
})

const months = Array.from({ length: 12 }, (_data, i) => new Date(2000, i, 1).toLocaleDateString('pt-BR', { month: 'short' }))

function selectMonth(year: number, month: number, close?: () => void): void {
  if (isAfterMax(year, month)) return
  const clickStart = startOfMonth(new Date(year, month, 1))
  model.value = clickStart
  if (close) { close() }
}

function monthVariant(year: number, month: number): 'solid' | 'ghost' {
  if (sameMonth(model.value, year, month)) return 'solid' as const
  return 'ghost' as const
}
</script>

<template>
  <UPopover :popper="{ placement: 'bottom-start' }" class="w-full max-w-72">
    <UButton icon="i-heroicons-calendar-days-20-solid" :label="label" size="sm" class="w-full" color="info" />
    <template #content>
      <div class="flex gap-4 p-2">
        <div class="w-64">
          <div class="mb-2 flex items-center justify-between">
            <UButton icon="i-heroicons-chevron-left" variant="ghost"  @click="year--" color="info" />
            <div class="text-sm font-medium">
              {{ year }}
            </div>
            <UButton icon="i-heroicons-chevron-right" variant="ghost" :disabled="year >= maxYear" @click="year++" color="info" />
          </div>
          <div class="grid grid-cols-3 gap-1">
            <UButton  v-for="(mName, m) in months"  :key="'y-' + m"  size="xs"  :variant="monthVariant(year, m)" color="info" class="justify-center"  :disabled="isAfterMax(year, m)"  @click="selectMonth(year, m)">
              {{ mName }}
            </UButton>
          </div>
        </div>
      </div>
    </template>
  </UPopover>
</template>
