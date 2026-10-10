<script setup lang="ts">
import { computed } from "vue"
import NativeImage from "@/components/NativeImage.vue"

const props = withDefaults(defineProps<{ tierKey?: string; division?: string; size?: number }>(), {
  tierKey: "",
  division: "",
  size: 40,
})

// Fallback medal, only drawn when the client cannot supply its own emblem.
const COLORS: Record<string, [string, string]> = {
  IRON: ["#8a8a8a", "#4d4d4d"],
  BRONZE: ["#c98a56", "#7a4a26"],
  SILVER: ["#c4cdd8", "#7d8896"],
  GOLD: ["#f6cf5f", "#c98f1a"],
  PLATINUM: ["#5fd3c3", "#238a7d"],
  EMERALD: ["#4fdc8f", "#1f8a52"],
  DIAMOND: ["#7fa6ff", "#3a5fc4"],
  MASTER: ["#cf86f0", "#8a3fb0"],
  GRANDMASTER: ["#ff7a73", "#b53a35"],
  CHALLENGER: ["#ffe27a", "#d9922b"],
}
const ROMAN: Record<string, string> = { I: "1", II: "2", III: "3", IV: "4" }

const colors = computed(() => COLORS[props.tierKey])
const mark = computed(() => ROMAN[props.division] ?? "")
const gid = computed(() => `g-${props.tierKey || "none"}`)
</script>

<template>
  <NativeImage kind="rank" :value="tierKey" :size="size">
    <svg :width="size" :height="size" viewBox="0 0 40 40" aria-hidden="true">
      <template v-if="colors">
        <defs>
          <linearGradient :id="gid" x1="0" y1="0" x2="1" y2="1">
            <stop offset="0" :stop-color="colors[0]" />
            <stop offset="1" :stop-color="colors[1]" />
          </linearGradient>
        </defs>
        <path d="M20 2 35 8.5v12c0 8.2-6.2 14.2-15 17.5C11.200 34.700 5 28.700 5 20.500v-12z" :fill="`url(#${gid})`" />
        <path d="M20 6.500 31 11v9.500c0 6-4.400 10.600-11 13.300C13.400 31.100 9 26.500 9 20.500V11z" fill="#000" fill-opacity="0.22" />
        <text v-if="mark" x="20" y="25" text-anchor="middle" font-size="15" font-weight="700" fill="#fff" fill-opacity="0.95">{{ mark }}</text>
        <path v-else d="M20 11.500l2.600 5.400 5.900.8-4.300 4.100 1 5.800L20 24.800l-5.200 2.800 1-5.800-4.300-4.100 5.900-.8z" fill="#fff" fill-opacity="0.92" />
      </template>
      <template v-else>
        <path d="M20 2 35 8.500v12c0 8.200-6.200 14.200-15 17.500C11.200 34.700 5 28.700 5 20.500v-12z" fill="none" stroke="#8a92a5" stroke-opacity="0.45" stroke-width="2" stroke-dasharray="4 3" />
        <text x="20" y="25" text-anchor="middle" font-size="15" font-weight="700" fill="#8a92a5" fill-opacity="0.6">?</text>
      </template>
    </svg>
  </NativeImage>
</template>
