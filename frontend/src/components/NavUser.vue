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
} from '@/components/ui/avatar'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from '@/components/ui/sidebar'
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
  "\u9ed1\u94c1": "IRON",
  "\u9752\u94dc": "BRONZE",
  "\u767d\u94f6": "SILVER",
  "\u9ec4\u91d1": "GOLD",
  "\u767d\u91d1": "PLATINUM",
  "\u7fe1\u7fe0": "EMERALD",
  "\u94bb\u77f3": "DIAMOND",
  "\u5927\u5e08": "MASTER",
  "\u5b97\u5e08": "GRANDMASTER",
  "\u738b\u8005": "CHALLENGER",
}

const UNRANKED_LABEL = "\u672a\u5b9a\u7ea7"
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
            <Avatar class="h-8 w-8 rounded-lg">
              <AvatarImage :src="user.avatar" :alt="user.name" />
              <AvatarFallback class="rounded-lg">
                CN
              </AvatarFallback>
            </Avatar>
            <div class="grid flex-1 text-left text-sm leading-tight">
              <span class="truncate font-medium">{{ user.name }}</span>
              <span class="truncate text-xs">
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
            <div class="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
              <Avatar class="h-8 w-8 rounded-lg">
                <AvatarImage :src="user.avatar" :alt="user.name" />
                <AvatarFallback class="rounded-lg">
                  CN
                </AvatarFallback>
              </Avatar>
              <div class="grid flex-1 text-left text-sm leading-tight">
                <span class="truncate font-semibold">{{ user.name }}</span>
                <span class="truncate text-xs">{{ badgeValue }}</span>
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
