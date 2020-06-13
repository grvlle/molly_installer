<template>
  <div class="component">
      <p class="text">Downloading Molly Binaries...</p>
    <progress-bar size="medium" font-size="1" text-fg-color="#fff" bar-color="#00F1B7" val="20" text="20%"></progress-bar>
  </div>
</template>

<script src="path/to/vue.js"></script>
<script src="path/to/vue-simple-progress.js"></script>
<script>
import ProgressBar from "vue-simple-progress";

export default {
  name: "Install",
  components: {
    ProgressBar,
  },
};
</script>

<style>
.component {
  flex: 1;
  
  width: 100%;
}
.text {
    margin-top: 4em;
    font-size: 0.7em;
}
</style>
