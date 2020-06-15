<template>
  <div id="app">
    <div class="box">
      <div class="box"></div>
      <img class="logo" src="./assets/images/logo.png" />
      <router-view></router-view>
      <div class="box"></div>
    </div>
  </div>
</template>

<script>
import "./assets/css/main.css";

export default {
  name: "app",
  components: {
  },

  mounted() {
    window.wails.Events.On("status", (msg) => {
      this.$store.state.progressMsg = msg;
    });
    window.wails.Events.On("progress", (percent) => {
      this.$store.state.progressPercent = percent;
    });
  },
};
</script>
