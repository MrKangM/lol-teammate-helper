<script setup lang="ts">
import { computed } from "vue"
import NativeImage from "@/components/NativeImage.vue"

const props = withDefaults(defineProps<{ tierKey?: string; size?: number }>(), { tierKey: "", size: 40 })

// Fallback drawing, only used when the client cannot supply the emblem.
const TIER_COLORS: Record<string, string> = {
  IRON: "#6b6b6b",
  BRONZE: "#a8693a",
  SILVER: "#9aa7b4",
  GOLD: "#e2b23c",
  PLATINUM: "#3fb8a9",
  EMERALD: "#2fbf71",
  DIAMOND: "#5b8def",
  MASTER: "#b05bd6",
  GRANDMASTER: "#d9534f",
  CHALLENGER: "#f5c542",
}
const color = computed(() => TIER_COLORS[props.tierKey] ?? "")
</script>

<template>
  <NativeImage kind="rank" :value="tierKey" :size="size">
    <svg :width="size" :height="size" viewBox="0 0 24 24" aria-hidden="true">
      <path
        d="M12 1.5 21 5.5v7.2c0 4.8-4 8.3-9 9.8-5-1.5-9-5-9-9.8V5.5z"
        :fill="color || 'none'"
        :fill-opacity="color ? 0.25 : 0"
        :stroke="color || 'currentColor'"
        :stroke-opacity="color ? 1 : 0.35"
        stroke-width="1.5"
        stroke-linejoin="round"
        :stroke-dasharray="color ? undefined : '2.5 2.5'"
      />
      <path v-if="color" d="M12 6.5 16 12l-4 5.5L8 12z" :fill="color" />
    </svg>
  </NativeImage>
</template>
