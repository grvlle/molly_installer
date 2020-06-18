<template>
  <section>
    <transition name="fade" mode="out-in">
      <div
        class="alert fade alert-simple alert-success"
        v-if="this.$store.state.showSuccessNotification"
      >
        <i class="start-icon far fa-check-circle faa-tada animated"></i>
        <strong class="font__weight-semibold">{{
          this.$store.state.successTitle
        }}</strong>
        {{ this.$store.state.successMsg }}
      </div>
    </transition>
  </section>
</template>

<script>
export default {};
</script>

<style lang="css" scoped>
body {
  background: #000e29;
}

.alert > .start-icon {
  margin-right: 0;
  min-width: 20px;
}

.alert > .start-icon {
  margin-right: 5px;
}

.greencross {
  font-size: 18px;
  color: #25ff0b;
  text-shadow: none;
}

.alert-simple.alert-success {
  margin-top: 3em;
  position: fixed;
  border: 1px solid rgba(36, 241, 6, 0.46);
  background-color: rgba(7, 149, 66, 0.12156862745098039);
  box-shadow: 0px 0px 2px #259c08;
  color: #0ad406;
  text-shadow: 2px 1px #00040a;
  padding: 1em;
  border-radius: 0.3rem;
  transition: 0.5s;
  cursor: pointer;
  transform: translate(-50%, -50%);
}
.alert-success:hover {
  background-color: rgba(7, 149, 66, 0.35);
  transition: 0.5s;
}

.fa-times {
  -webkit-animation: blink-1 2s infinite both;
  animation: blink-1 2s infinite both;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s;
}
.fade-enter, .fade-leave-to /* .fade-leave-active below version 2.1.8 */ {
  opacity: 0;
}

/**
 * ----------------------------------------
 * animation blink-1
 * ----------------------------------------
 */
@-webkit-keyframes blink-1 {
  0%,
  50%,
  100% {
    opacity: 1;
  }
  25%,
  75% {
    opacity: 0;
  }
}
@keyframes blink-1 {
  0%,
  50%,
  100% {
    opacity: 1;
  }
  25%,
  75% {
    opacity: 0;
  }
}
</style>
