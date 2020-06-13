<template>
    <div id="app">
    
      <div class="box">
        <div class="box"></div>
        <img class="logo" src="./assets/images/logo.png">
        <router-view></router-view>
        <div class="box"></div>
      
    </div>
  </div>
  
</template>

<script>
// import Home from './views/Home.vue'
import "./assets/css/main.css";


export default {
  name: "app",
  components: {
    // Home,
  },
};
</script>
