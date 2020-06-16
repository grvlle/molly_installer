<template>
<div class="container">
        <p class="box info-text"><b>Welcome to the Molly Wallet installer!</b><br />
        This will install the Molly Wallet that will let you store your $DAG tokens and interact with the HyperGraph Network. 
        </p>
        <div class="component"><NavButtons /></div>
</div>
</template>

<script>
import NavButtons from "../components/NavButtons.vue";
import "../assets/css/main.css";

export default {
    name: 'Home',
    components: {
        NavButtons,
    }

}
</script>

<style>
.container {
  display: flex; /* establish flex container */
  flex: 3;
  flex-direction: column; /* make main axis vertical */
  justify-content: center; /* center items vertically, in this case */
  align-items: center; /* center items horizontally, in this case */
  }

  .component {
  flex: 1;
}
</style>