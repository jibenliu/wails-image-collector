<template>
  <div id="app">
    <Login/>
  </div>
</template>

<script>
import Login from "./components/Login.vue";
import 'element-theme-chalk';

export default {
  name: "app",
  components: {
    Login,
  }
};
</script>
