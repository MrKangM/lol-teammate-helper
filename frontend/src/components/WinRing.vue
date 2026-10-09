<script setup lang="ts">
import { computed } from "vue"

const props = withDefaults(defineProps<{ rate: number; size?: number; label?: string }>(), { size: 84, label: "胜率" })

const stroke = 7
const radius = computed(() => (props.size - stroke) / 2)
const circumference = computed(() => 2 * Math.PI * radius.value)
const dash = computed(() => `${(Math.min(Math.max(props.rate, 0), 100) / 100) * circumference.value} ${circumference.value}`)
const color = computed(() => (props.rate >= 55 ? "var(--color-win)" : props.rate >= 48 ? "var(--color-gold)" : "var(--color-loss)"))
</script>

<template>
  <div class="relative" :style="{ width: `${size}px`, height: `${size}px` }">
    <svg :width="size" :height="size" class="-rotate-90">
      <circle :cx="size / 2" :cy="size / 2" :r="radius" fill="none" stroke="var(--color-line)" :stroke-width="stroke" />
      <circle
        :cx="size / 2"
        :cy="size / 2"
        :r="radius"
        fill="none"
        :stroke="color"
        :stroke-width="stroke"
        stroke-linecap="round"
        :stroke-dasharray="dash"
      />
    </svg>
    <div class="absolute inset-0 flex flex-col items-center justify-center leading-tight">
      <span class="num text-lg font-bold text-gold-bright">{{ Math.round(rate) }}%</span>
      <span class="text-[10px] text-muted">{{ label }}</span>
    </div>
  </div>
</template>
