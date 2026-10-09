<script setup lang="ts">
import NativeImage from "@/components/NativeImage.vue"
import { positionLabel } from "@/lib/format"

withDefaults(defineProps<{ position?: string; size?: number }>(), { position: "", size: 22 })
</script>

<template>
  <NativeImage v-if="position" kind="position" :value="position" :size="size">
    <span class="text-[11px] font-medium text-gold-dim">{{ positionLabel(position) }}</span>
  </NativeImage>
  <span v-else class="text-[11px] text-muted">待定</span>
</template>
