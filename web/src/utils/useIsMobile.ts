import { ref, computed, onMounted, onUnmounted } from "vue";

const windowWidth = ref(window.innerWidth);

const updateWidth = () => {
  windowWidth.value = window.innerWidth;
};

/**
 * 判断是否是移动端
 * @param width
 * @returns
 */
export function useIsMobile(width = 768) {
  onMounted(() => {
    window.addEventListener("resize", updateWidth);
    updateWidth();
  });

  onUnmounted(() => {
    window.removeEventListener("resize", updateWidth);
  });

  return computed(() => windowWidth.value <= width);
}
