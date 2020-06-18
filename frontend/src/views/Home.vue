<template>
<div class="container">
  <div class="info-text-box">
        <p class="info-text"><center><b>Welcome to the Molly Wallet installer!</b></center><br />
        This will install the Molly Wallet that will let you store your $DAG tokens and interact with the HyperGraph Network. Please note that the following dependencies may also be installed on your system. <br />
        <br />
        <table>
        <tr>
          <td><li><a class="links" href="#">Scoop</a></li></td>
          <td><i> Package management tool for Windows</i></td>
        </tr>
          <tr>
          <td><li><a class="links" href="#">Java</a></li></td>
          <td><i> Neccessary to interact with the wallet SDK</i></td>
          </tr>
          <tr>
          <td><li><a class="links" href="#">7zip</a></li></td>
          <td><i> Needed to extract contents from the package manager</i></td>
        </tr>
        </table>
        <br />
        These packages will also be removed as part of the Unistallation proceedure. So make sure to leverage this installer when you wish to remove Molly wallet from your system.
        </p>
        </div>
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

<style scoped>
.container {
  display: flex; /* establish flex container */
  flex: 3;
  flex-direction: column; /* make main axis vertical */
  justify-content: center; /* center items vertically, in this case */
  align-items: center; /* center items horizontally, in this case */
  }

  td, th {
  text-align: left;
  padding-right: 4em;
}

table {
  margin-left: 1em;
}

  a.links {
    text-decoration: none;
    color : #00f1b7;
    text-transform: uppercase;
    font-size: 1rem;
  }

  a.links:hover {
    color: #90ffe5;
    text-shadow : 0px 0px 10px rgba(255,255,255,0.6), 0px 0px 30px rgba(255,255,255,0.4), 0px 0px 50px rgba(255,255,255,0.3), 0px 0px 180px rgba(255,255,255,0.3);
  }

  .component {
  flex: 1;
  position: fixed;
  margin-top: 9em;
    left: 50%;
  top: 60%;
  transform: translate(-50%, -50%);
}
</style>