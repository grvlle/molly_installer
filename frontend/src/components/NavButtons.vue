<template>
  <div class="container">
    <div class="box">
      <h1>{{ message }}</h1>
      <!-- <center><a @click="getMessage">Press Me!</a></center> -->
      <div align="center">
        <a class="btn_cancel" @click="getMessage">Cancel</a>
      </div>
      <div align="center"><router-link to="/install"><a class="btn">Install ↠</a></router-link></div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      message: " ",
    };
  },
  methods: {
    getMessage: function() {
      var self = this;
      window.backend.basic().then((result) => {
        self.message = result;
      });
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->

<style scoped>
.box {
  display: flex;
  flex-direction: row;
}

.btn {
  float: right;
  color: #494949 !important;
  text-transform: uppercase;
  text-decoration: none;
  background: #00f1b7;
  padding: 20px;
  border: 4px solid #02cc9d !important;
  display: block;
  transition: all 0.4s ease 0s;
  cursor: pointer;
}
.btn:hover {
  color: #fff !important;
  background: #0f41ba;
  border-color: #0bd4ff !important;
  transition: all 0.4s ease 0s;
}

.btn_cancel {
  float: left;
  margin-right: 2em;
  color: #f7f7f7 !important;
  text-transform: uppercase;
  text-decoration: none;
  background: #5a5a5a;
  padding: 20px;
  border: 4px solid #494949 !important;
  display: block;
  transition: all 0.4s ease 0s;
  cursor: pointer;
}
.btn_cancel:hover {
  color: rgb(0, 0, 0) !important;
  background: rgb(75, 75, 75);
  border-color: rgb(131, 131, 131) !important;
  transition: all 0.4s ease 0s;
}

</style>
