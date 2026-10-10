<script setup lang="ts">
import { computed } from "vue"
import NativeImage from "@/components/NativeImage.vue"

const props = withDefaults(defineProps<{ position?: string; size?: number }>(), { position: "", size: 20 })

// Minimal lane glyphs used when the client does not supply its own icons:
// the map's three lanes in a dim colour with the player's lane highlighted.
const DIM = "#4a5266"
const HI = "#cfd6e6"
const pos = computed(() => props.position.toLowerCase())
</script>

<template>
  <NativeImage v-if="position" kind="position" :value="position" :size="size">
    <svg :width="size" :height="size" viewBox="0 0 24 24" fill="none" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
      <path d="M4 20V4h16" :stroke="pos === 'top' ? HI : DIM" stroke-width="2.4" />
      <path d="M20 4v16H4" :stroke="pos === 'bottom' ? HI : DIM" stroke-width="2.4" />
      <path d="M5 19 19 5" :stroke="pos === 'middle' ? HI : DIM" stroke-width="2.4" />
      <template v-if="pos === 'jungle'">
        <path d="M8 8l3 3M13 13l3 3" :stroke="HI" stroke-width="3" />
      </template>
      <template v-if="pos === 'utility'">
        <circle cx="18" cy="18" r="3" :fill="HI" stroke="none" />
      </template>
    </svg>
  </NativeImage>
  <span v-else class="text-[11px] text-muted">?</span>
</template>
