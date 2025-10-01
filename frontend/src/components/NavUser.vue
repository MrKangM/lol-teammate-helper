<script setup lang="ts">
import { computed } from "vue"
import {
  BadgeCheck,
  Bell,
  ChevronsUpDown,
  CreditCard,
  LogOut,
  Sparkles,
} from "lucide-vue-next"

import {
  Avatar,
  AvatarFallback,
  AvatarImage,
} from "@/components/ui/avatar"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "@/components/ui/sidebar"
import { Badge } from "@/components/ui/badge"

const props = defineProps<{
  user: {
    name: string
    email: string
    avatar: string
  }
}>()

const { isMobile } = useSidebar()

const tierColorMap = {
  IRON: "border-transparent bg-stone-600 text-white",
  BRONZE: "border-transparent bg-amber-800 text-white",
  SILVER: "border-transparent bg-slate-300 text-slate-900",
  GOLD: "border-transparent bg-amber-500 text-white",
  PLATINUM: "border-transparent bg-teal-500 text-white",
  EMERALD: "border-transparent bg-emerald-500 text-white",
  DIAMOND: "border-transparent bg-sky-500 text-white",
  MASTER: "border-transparent bg-purple-600 text-white",
  GRANDMASTER: "border-transparent bg-rose-600 text-white",
  CHALLENGER: "border-transparent bg-indigo-600 text-white",
} as const

type TierKey = keyof typeof tierColorMap

const chineseToTierMap: Record<string, TierKey> = {
  "黑铁": "IRON",
  "青铜": "BRONZE",
  "白银": "SILVER",
  "黄金": "GOLD",
  "白金": "PLATINUM",
  "翡翠": "EMERALD",
  "钻石": "DIAMOND",
  "大师": "MASTER",
  "宗师": "GRANDMASTER",
  "王者": "CHALLENGER",
}

const UNRANKED_LABEL = "未定级"
const UNRANKED_CLASS = "border-transparent bg-muted text-muted-foreground"
const DEFAULT_BADGE_CLASS = "border-transparent bg-secondary text-secondary-foreground"

const tierKeys = Object.keys(tierColorMap) as TierKey[]

const resolveTierKey = (raw: string): TierKey | undefined => {
  const value = raw.trim()
  if (!value) return undefined

  const upper = value.toUpperCase()
  for (const tier of tierKeys) {
    if (upper.includes(tier)) {
      return tier
    }
  }

  for (const [label, tier] of Object.entries(chineseToTierMap)) {
    if (value.includes(label)) {
      return tier
    }
  }

  return undefined
}

const badgeValue = computed(() => {
  const raw = props.user.email ?? ""
  const normalized = raw.trim()
  return normalized.length > 0 ? normalized : UNRANKED_LABEL
})

const badgeClass = computed(() => {
  const value = badgeValue.value
  if (value === UNRANKED_LABEL) {
    return UNRANKED_CLASS
  }

  const tier = resolveTierKey(value)
  if (tier) {
    return tierColorMap[tier]
  }

  return DEFAULT_BADGE_CLASS
})
</script>

<template>
  <SidebarMenu>
    <SidebarMenuItem>
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <SidebarMenuButton
            size="lg"
            class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
          >
            <Avatar class="h-12 w-12 rounded-full border border-white/10">
              <AvatarImage :src="user.avatar" :alt="user.name" />
              <AvatarFallback class="rounded-full">
                CN
              </AvatarFallback>
            </Avatar>
            <div class="flex flex-1 flex-col text-left text-sm leading-tight">
              <span class="truncate text-base font-semibold">{{ user.name }}</span>
              <span class="mt-1 inline-flex">
                <Badge :class="badgeClass">{{ badgeValue }}</Badge>
              </span>
            </div>
            <ChevronsUpDown class="ml-auto size-4" />
          </SidebarMenuButton>
        </DropdownMenuTrigger>
        <DropdownMenuContent
          class="w-[--reka-dropdown-menu-trigger-width] min-w-56 rounded-lg"
          :side="isMobile ? 'bottom' : 'right'"
          align="end"
          :side-offset="4"
        >
          <DropdownMenuLabel class="p-0 font-normal">
            <div class="flex items-center gap-3 px-1 py-1.5 text-left text-sm">
              <Avatar class="h-12 w-12 rounded-full border border-white/10">
                <AvatarImage :src="user.avatar" :alt="user.name" />
                <AvatarFallback class="rounded-full">
                  CN
                </AvatarFallback>
              </Avatar>
              <div class="flex flex-1 flex-col text-left leading-tight">
                <span class="truncate text-base font-semibold">{{ user.name }}</span>
                <span class="mt-1 inline-flex text-xs">
                  <Badge :class="badgeClass">{{ badgeValue }}</Badge>
                </span>
              </div>
            </div>
          </DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuGroup>
            <DropdownMenuItem>
              <Sparkles />
              Upgrade to Pro
            </DropdownMenuItem>
          </DropdownMenuGroup>
          <DropdownMenuSeparator />
          <DropdownMenuGroup>
            <DropdownMenuItem>
              <BadgeCheck />
              Account
            </DropdownMenuItem>
            <DropdownMenuItem>
              <CreditCard />
              Billing
            </DropdownMenuItem>
            <DropdownMenuItem>
              <Bell />
              Notifications
            </DropdownMenuItem>
          </DropdownMenuGroup>
          <DropdownMenuSeparator />
          <DropdownMenuItem>
            <LogOut />
            Log out
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </SidebarMenuItem>
  </SidebarMenu>
</template>
