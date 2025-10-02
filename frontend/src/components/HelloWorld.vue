<script setup lang="ts">
import { DonutChart } from "@/components/ui/chart-donut"

const data = [
  { name: "胜场", total: Math.floor(Math.random() * 2000) + 500, predicted: Math.floor(Math.random() * 2000) + 500 },
  { name: "败场", total: Math.floor(Math.random() * 2000) + 500, predicted: Math.floor(Math.random() * 2000) + 500 },
]

const valueFormatter = (tick: number | Date) => typeof tick === "number" ? `$ ${new Intl.NumberFormat("us").format(tick).toString()}` : ""
</script>

<template>
  <DonutChart
      index="name"
      :category="'total'"
      :data="data"
      :value-formatter="valueFormatter"
      :colors="['green', 'red']"
  />
</template>