import { ref, onMounted, onUnmounted, computed } from 'vue'

export type Breakpoint = 'xs' | 'sm' | 'md' | 'lg' | 'xl'

const breakpointValues: Record<Breakpoint, number> = {
  xs: 0,
  sm: 768,
  md: 992,
  lg: 1200,
  xl: 1920,
}

export function useBreakpoint() {
  const width = ref(window.innerWidth)

  const updateWidth = () => {
    width.value = window.innerWidth
  }

  onMounted(() => {
    window.addEventListener('resize', updateWidth)
  })

  onUnmounted(() => {
    window.removeEventListener('resize', updateWidth)
  })

  const current = computed<Breakpoint>(() => {
    if (width.value >= breakpointValues.xl) return 'xl'
    if (width.value >= breakpointValues.lg) return 'lg'
    if (width.value >= breakpointValues.md) return 'md'
    if (width.value >= breakpointValues.sm) return 'sm'
    return 'xs'
  })

  const isXs = computed(() => width.value < breakpointValues.sm)
  const isSm = computed(() => width.value >= breakpointValues.sm && width.value < breakpointValues.md)
  const isMd = computed(() => width.value >= breakpointValues.md && width.value < breakpointValues.lg)
  const isLg = computed(() => width.value >= breakpointValues.lg && width.value < breakpointValues.xl)
  const isXl = computed(() => width.value >= breakpointValues.xl)

  const isMobile = computed(() => width.value < breakpointValues.md)
  const isTablet = computed(() => width.value >= breakpointValues.md && width.value < breakpointValues.lg)
  const isDesktop = computed(() => width.value >= breakpointValues.lg)

  return {
    width,
    current,
    isXs,
    isSm,
    isMd,
    isLg,
    isXl,
    isMobile,
    isTablet,
    isDesktop,
  }
}