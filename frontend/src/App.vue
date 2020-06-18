<template>
  <div id="app">
    <Success />
    <Error />
    <div class="box">
      <div class="box"></div>
      <div><img class="logo" src="./assets/images/logo.png" /></div>

      <router-view></router-view>
      <div class="box"></div>
    </div>
  </div>
</template>

<script>
import "./assets/css/main.css";
import Success from "./components/notifications/Success.vue";
import Error from "./components/notifications/Error.vue";

export default {
  name: "app",
  components: {
    Error,
    Success,
  },
  methods: {
    showSuccessNotification() {
      this.$store.state.showSuccessNotification = false;
    },
    sendSuccessNotification() {
      setTimeout(
        function() {
          this.showSuccessNotification();
        }.bind(this),
        5000
      );
      this.$store.state.showSuccessNotification = true;
    },
    showErrorNotification() {
      this.$store.state.showErrorNotification = false;
    },
    sendErrorNotification() {
      setTimeout(
        function() {
          this.showErrorNotification();
        }.bind(this),
        5000
      );
      this.$store.state.showErrorNotification = true;
    },
  },
  mounted() {
    window.wails.Events.On("status", (msg) => {
      this.$store.state.progressMsg = msg;
    });
    window.wails.Events.On("progress", (percent) => {
      this.$store.state.progressPercent = percent;
    });
    window.wails.Events.On("success", (title, msg) => {
      this.$store.state.successTitle = title;
      this.$store.state.successMsg = msg;
      this.sendSuccessNotification();
    });
    window.wails.Events.On("error", (title, msg) => {
      this.$store.state.errorTitle = title;
      this.$store.state.errorMsg = msg;
      this.sendErrorNotification();
    });
  },
};
</script>
