import Vue from 'vue';
import Vuex from 'vuex';

Vue.use(Vuex);
export const store = new Vuex.Store({
    state: {
        progressMsg: "Installing Molly Wallet...",
        progressPercent: "0",
    },
    })