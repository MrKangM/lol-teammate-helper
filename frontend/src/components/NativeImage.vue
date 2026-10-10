<script setup lang="ts">
import { ref, watch } from "vue"
import { loadPositionIcon, loadRankEmblem } from "@/lib/assets"

// Shows an image served by the League client. While it loads, or when the
// client does not have it, the default slot is rendered instead.
const props = defineProps<{ kind: "rank" | "position"; value: string; size?: number }>()

const src = ref("")

watch(
  () => [props.kind, props.value] as const,
  async ([kind, value]) => {
    src.value = ""
    const load = kind === "rank" ? loadRankEmblem : loadPositionIcon
    src.value = await load(value)
  },
  { immediate: true },
)
</script>

<template>
  <img v-if="src" :src="src" alt="" :style="{ width: `${size ?? 32}px`, height: `${size ?? 32}px` }" class="object-contain" />
  <slot v-else />
</template>
