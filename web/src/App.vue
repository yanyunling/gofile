<template>
  <router-view></router-view>
  <login />
</template>

<script lang="ts">
import { defineComponent } from "vue";
import Login from "@/views/login/index.vue";
export default defineComponent({
  name: "App",
  components: {
    Login,
  },
});
</script>

<style>
#app {
  position: fixed;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
}
</style>
