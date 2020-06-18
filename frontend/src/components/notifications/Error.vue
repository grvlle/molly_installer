<template>
<section>
  <transition name="fade" mode="out-in">
        <div v-if="this.$store.state.showErrorNotification"
          class="alert fade alert-simple alert-danger"
          role="alert"
          data-brk-library="component__alert"
        >
          <i class="start-icon far fa-times-circle faa-pulse animated"></i>
          <strong class="font__weight-semibold">{{
            this.$store.state.errorTitle
          }}</strong>
          {{ this.$store.state.errorMsg }}
        </div>
  </transition>
</section>
</template>

<script>
export default {};
</script>

<style>
.alert > .start-icon {
  margin-right: 0;
  min-width: 20px;
}

.alert > .start-icon {
  margin-right: 5px;
}

.alert-simple.alert-danger {
  margin-top: 3em;
  position: fixed;
  margin-top: 2em;
  border: 1px solid rgba(241, 6, 6, 0.81);
  background-color: rgba(220, 17, 1, 0.16);
  box-shadow: 0px 0px 2px #ff0303;
  color: #ff0303;
  text-shadow: 2px 1px #00040a;
  padding: 1em;
  border-radius: 0.3rem;
  transition: 0.5s;
  cursor: pointer;
  transform: translate(-50%, -50%);
}

.alert-danger:hover {
  background-color: rgba(220, 17, 1, 0.33);
  transition: 0.5s;
}

.danger {
  font-size: 18px;
  color: #ff0303;
  text-shadow: none;
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
